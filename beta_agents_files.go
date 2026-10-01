package openai

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"maps"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
)

const betaAgentFileLimit int64 = 50 * 1024 * 1024

// BetaAgentPreparedFiles holds ordinary environment inputs and Files API objects
// created by this helper. The caller owns their lifecycle; cleanup is explicit.
// This beta helper is experimental.
type BetaAgentPreparedFiles struct {
	Files   []HostedEnvironmentFileParamUnion
	Uploads []FileObject
}

// BetaAgentFilePreparationError retains successful uploads if a later operation
// fails. Files are never deleted automatically. This beta helper is experimental.
type BetaAgentFilePreparationError struct {
	Prepared *BetaAgentPreparedFiles
	Cause    error
}

func (e *BetaAgentFilePreparationError) Error() string { return "cannot prepare beta agent files" }
func (e *BetaAgentFilePreparationError) Unwrap() error { return e.Cause }

type betaAgentLocalFile struct {
	destination, source string
	info                os.FileInfo
}

// Prepare uploads selected local files for a future hosted session. files maps
// absolute /workspace destinations to local paths. All selections are checked
// before upload, including symlinks, collisions, 50-file and 50 MiB total limits.
func (r *BetaAgentEnvironmentFileService) Prepare(ctx context.Context, files map[string]string, opts ...option.RequestOption) (*BetaAgentPreparedFiles, error) {
	selected, err := betaAgentPrepareSelection(files, true)
	if err != nil {
		return nil, err
	}
	return r.prepareFiles(ctx, selected, opts...)
}

// PrepareDirectory prepares a one-time selection of regular files. include is an
// explicit list of filepath.Match patterns relative to directory (not a mount or
// synchronization rule). Symlinks are not followed. The caller owns upload cleanup.
func (r *BetaAgentEnvironmentFileService) PrepareDirectory(ctx context.Context, directory, destination string, include []string, opts ...option.RequestOption) (*BetaAgentPreparedFiles, error) {
	if len(include) == 0 {
		return nil, errors.New("include patterns must be explicit")
	}
	for _, pattern := range include {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return nil, err
		}
	}
	directory = filepath.Clean(directory)
	info, err := os.Lstat(directory)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("directory preparation requires a non-symlink directory")
	}
	directory, err = filepath.EvalSymlinks(directory)
	if err != nil {
		return nil, err
	}
	directory, err = filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	files := make(map[string]string)
	err = filepath.WalkDir(directory, func(source string, entry fs.DirEntry, walkErr error) error {
		if contextErr := ctx.Err(); contextErr != nil {
			return contextErr
		}
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, relErr := filepath.Rel(directory, source)
		if relErr != nil {
			return relErr
		}
		for _, pattern := range include {
			matched, _ := filepath.Match(pattern, relative)
			if matched {
				files[destination+"/"+filepath.ToSlash(relative)] = source
				if len(files) > 50 {
					return errors.New("initial hosted files exceed 50 files")
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.Prepare(ctx, files, opts...)
}

// Upload stages one local file into a live environment and returns its owned
// Files API upload along with the live file metadata. A staging failure retains
// the uploaded file in BetaAgentFilePreparationError; cleanup remains explicit.
func (r *BetaAgentEnvironmentFileService) Upload(ctx context.Context, environmentID, source, destination string, opts ...option.RequestOption) (*BetaAgentStagedFile, error) {
	if environmentID == "" {
		return nil, errors.New("staging requires an environment ID")
	}
	selected, err := betaAgentPrepareSelection(map[string]string{destination: source}, false)
	if err != nil {
		return nil, err
	}
	prepared, err := r.prepareFiles(ctx, selected, opts...)
	if err != nil {
		return nil, err
	}
	service := *r
	capture, require := agentStreamResponseGuard[EnvironmentFile]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	file, err := service.New(ctx, environmentID, BetaAgentEnvironmentFileNewParams{HostedEnvironmentFileParam: prepared.Files[0]}, append(slices.Clone(opts), require)...)
	if err != nil {
		return nil, &BetaAgentFilePreparationError{Prepared: prepared, Cause: err}
	}
	return &BetaAgentStagedFile{File: *file, Upload: prepared.Uploads[0], Input: prepared.Files[0]}, nil
}

// BetaAgentStagedFile keeps the live environment file and owned Files API upload
// separate. This beta helper is experimental.
type BetaAgentStagedFile struct {
	File   EnvironmentFile
	Upload FileObject
	Input  HostedEnvironmentFileParamUnion
}

func (r *BetaAgentEnvironmentFileService) prepareFiles(ctx context.Context, selected []betaAgentLocalFile, opts ...option.RequestOption) (*BetaAgentPreparedFiles, error) {
	if len(selected) > 1 {
		cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, "files", nil, nil, append(slices.Clone(r.Options), opts...)...)
		if err != nil {
			return nil, err
		}
		if cfg.Request.Header.Get("Idempotency-Key") != "" {
			return nil, errors.New("multi-file preparation cannot reuse one Idempotency-Key")
		}
	}
	prepared := &BetaAgentPreparedFiles{Files: []HostedEnvironmentFileParamUnion{}, Uploads: []FileObject{}}
	files := NewFileService(r.Options...)
	capture, require := agentStreamResponseGuard[FileObject]()
	files.Options = append([]option.RequestOption{capture}, files.Options...)
	for _, item := range selected {
		source, err := os.Open(item.source)
		if err != nil {
			return prepared, &BetaAgentFilePreparationError{Prepared: prepared, Cause: err}
		}
		info, err := source.Stat()
		if err != nil || !os.SameFile(item.info, info) || info.Size() != item.info.Size() {
			_ = source.Close()
			if err == nil {
				err = errors.New("source changed after file preflight")
			}
			return prepared, &BetaAgentFilePreparationError{Prepared: prepared, Cause: err}
		}
		upload, uploadErr := files.New(ctx, FileNewParams{File: File(io.NewSectionReader(source, 0, info.Size()), filepath.Base(item.source), "application/octet-stream"), Purpose: FilePurposeUserData}, append(slices.Clone(opts), require)...)
		closeErr := source.Close()
		if uploadErr != nil {
			return prepared, &BetaAgentFilePreparationError{Prepared: prepared, Cause: uploadErr}
		}
		prepared.Uploads = append(prepared.Uploads, *upload)
		prepared.Files = append(prepared.Files, HostedEnvironmentFileParamOfParamFileID(upload.ID, item.destination))
		if closeErr != nil {
			return prepared, &BetaAgentFilePreparationError{Prepared: prepared, Cause: closeErr}
		}
	}
	return prepared, nil
}

func betaAgentPrepareSelection(files map[string]string, initial bool) ([]betaAgentLocalFile, error) {
	if initial && len(files) > 50 {
		return nil, errors.New("initial hosted files exceed 50 files")
	}
	selected := make([]betaAgentLocalFile, 0, len(files))
	var total int64
	destinations := slices.Sorted(maps.Keys(files))
	for _, destination := range destinations {
		if err := betaAgentFileDestination(destination); err != nil {
			return nil, err
		}
		for parent := path.Dir(destination); parent != "/workspace"; parent = path.Dir(parent) {
			if _, exists := files[parent]; exists {
				return nil, errors.New("file destinations overlap")
			}
		}
		source, err := filepath.Abs(files[destination])
		if err != nil {
			return nil, err
		}
		info, err := os.Lstat(source)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, errors.New("selected source must be a regular file")
		}
		if info.Size() > betaAgentFileLimit {
			return nil, errors.New("hosted file exceeds 50 MiB")
		}
		total += info.Size()
		if initial && total > betaAgentFileLimit {
			return nil, errors.New("initial hosted files exceed 50 MiB total")
		}
		selected = append(selected, betaAgentLocalFile{destination: destination, source: source, info: info})
	}
	return selected, nil
}
func betaAgentFileDestination(destination string) error {
	if !strings.HasPrefix(destination, "/workspace/") || strings.ContainsAny(destination, "\x00\\") || path.Clean(destination) != destination || destination == "/workspace/outputs" {
		return errors.New("destination must name a file inside /workspace using a clean absolute POSIX path")
	}
	first := strings.Split(strings.TrimPrefix(destination, "/workspace/"), "/")[0]
	if first == ".codex" || first == ".managed-agents" || strings.HasPrefix(first, ".managed-agents-") {
		return errors.New("destination uses a reserved environment path")
	}
	return nil
}

// BetaAgentResultArtifacts selects immutable artifacts belonging to one completed
// result. It does not store anything on the result. This beta helper is experimental.
type BetaAgentResultArtifacts struct {
	service           *BetaAgentSessionArtifactService
	sessionID, turnID string
}

// ForResult scopes artifact lookup to this result's session and root turn without
// making requests or retaining the result's message payloads.
func (r *BetaAgentSessionArtifactService) ForResult(result *BetaAgentTurnResult) *BetaAgentResultArtifacts {
	artifacts := &BetaAgentResultArtifacts{service: r}
	if result != nil {
		artifacts.sessionID = result.SessionID()
		artifacts.turnID = result.TurnID()
	}
	return artifacts
}

// Find retrieves the artifact with the exact published path for the selected turn.
// It follows every page and reports missing or ambiguous matches.
func (r *BetaAgentResultArtifacts) Find(ctx context.Context, path string, opts ...option.RequestOption) (*SessionArtifact, error) {
	if r.sessionID == "" || r.turnID == "" {
		return nil, errors.New("artifact lookup requires a selected result turn")
	}
	var found *SessionArtifact
	service := *r.service
	capture, require := agentStreamResponseGuard[pagination.CursorPage[SessionArtifact]]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	items := service.ListAutoPaging(ctx, r.sessionID, BetaAgentSessionArtifactListParams{}, append(slices.Clone(opts), require)...)
	for items.Next() {
		item := items.Current()
		if item.TurnID == r.turnID && item.Path == path {
			if found != nil && found.ID != item.ID {
				return nil, errors.New("multiple artifacts match the result turn and requested path")
			}
			found = &item
		}
	}
	if err := items.Err(); err != nil {
		return nil, err
	}
	if found == nil {
		return nil, errors.New("no artifact matches the result turn and requested path")
	}
	return found, nil
}

// Download streams an exact turn/path match into the caller's writer, returning
// immutable artifact metadata. It never derives a local destination from a remote
// path. The caller owns the writer and its cleanup.
func (r *BetaAgentResultArtifacts) Download(ctx context.Context, path string, to io.Writer, opts ...option.RequestOption) (*SessionArtifact, error) {
	if to == nil {
		return nil, errors.New("artifact download requires a destination writer")
	}
	artifact, err := r.Find(ctx, path, opts...)
	if err != nil {
		return nil, err
	}
	service := *r.service
	capture, require := agentStreamResponseGuard[http.Response]()
	service.Options = append([]option.RequestOption{capture}, service.Options...)
	response, err := service.Content(ctx, r.sessionID, artifact.ID, append(slices.Clone(opts), require)...)
	if err != nil {
		return nil, err
	}
	_, copyErr := io.Copy(to, response.Body)
	closeErr := response.Body.Close()
	if copyErr != nil {
		return artifact, copyErr
	}
	return artifact, closeErr
}
