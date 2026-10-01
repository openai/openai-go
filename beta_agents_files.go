package openai

import (
	"context"
	"errors"
	"io"
	"maps"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
)

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
// before upload for regular source files and conflicting destinations.
// File-count, size, and destination-length limits are enforced by the API.
// Local paths must be application-controlled and stable during preparation.
func (r *BetaAgentEnvironmentFileService) Prepare(ctx context.Context, files map[string]string, opts ...option.RequestOption) (*BetaAgentPreparedFiles, error) {
	selected, err := betaAgentPrepareSelection(ctx, files)
	if err != nil {
		return nil, err
	}
	return r.prepareFiles(ctx, selected, opts...)
}

// PrepareDirectory prepares a one-time selection of regular files. include is an
// explicit list of filepath.Glob patterns relative to directory (not a mount or
// synchronization rule). Glob skips missing or unreadable entries. Selected
// symlink entries are not followed. The directory must be application-controlled
// and stable during preparation; this helper is not a sandbox for untrusted paths
// or concurrent filesystem writers.
// The caller owns upload cleanup.
func (r *BetaAgentEnvironmentFileService) PrepareDirectory(ctx context.Context, directory, destination string, include []string, opts ...option.RequestOption) (*BetaAgentPreparedFiles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if directory == "" {
		return nil, errors.New("directory must be explicit")
	}
	if len(include) == 0 {
		return nil, errors.New("include patterns must be explicit")
	}
	for _, pattern := range include {
		// Match can stop validating after an unmatched star-delimited chunk.
		// For this syntax-only probe, ? preserves validity and forces a full scan.
		if _, err := filepath.Match(strings.ReplaceAll(pattern, "*", "?"), ""); err != nil {
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
	// Keep the application-owned directory literal when composing the glob.
	volume := filepath.VolumeName(directory)
	rootPattern := directory[len(volume):]
	if runtime.GOOS != "windows" {
		rootPattern = strings.ReplaceAll(rootPattern, `\`, `\\`)
	}
	rootPattern = volume + strings.NewReplacer("[", "[[]", "*", "[*]", "?", "[?]").Replace(rootPattern)
	files := make(map[string]string)
	for _, pattern := range include {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !filepath.IsLocal(pattern) {
			return nil, errors.New("include patterns must stay inside the selected directory")
		}
		selected, selectErr := filepath.Glob(rootPattern + string(filepath.Separator) + pattern)
		if selectErr != nil {
			return nil, selectErr
		}
		for _, source := range selected {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			relative, relErr := filepath.Rel(directory, source)
			if relErr != nil {
				return nil, relErr
			}
			if !filepath.IsLocal(relative) {
				return nil, errors.New("selected source is outside the directory")
			}
			symlinkParent := false
			for parent := filepath.Dir(relative); parent != "."; parent = filepath.Dir(parent) {
				info, parentErr := os.Lstat(filepath.Join(directory, parent))
				if parentErr != nil {
					return nil, parentErr
				}
				if info.Mode()&os.ModeSymlink != 0 {
					symlinkParent = true
					break
				}
			}
			if symlinkParent {
				continue
			}
			info, statErr := os.Lstat(source)
			if statErr != nil {
				return nil, statErr
			}
			if !info.IsDir() {
				files[destination+"/"+filepath.ToSlash(relative)] = source
			}
		}
	}
	return r.Prepare(ctx, files, opts...)
}

// Upload stages one local file into a live environment and returns its owned
// Files API upload along with the live file metadata. A staging failure retains
// the uploaded file in BetaAgentFilePreparationError; cleanup remains explicit.
// The local source path must be application-controlled and stable during upload.
func (r *BetaAgentEnvironmentFileService) Upload(ctx context.Context, environmentID, source, destination string, opts ...option.RequestOption) (*BetaAgentStagedFile, error) {
	if environmentID == "" {
		return nil, errors.New("staging requires an environment ID")
	}
	selected, err := betaAgentPrepareSelection(ctx, map[string]string{destination: source})
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
	if err == nil && file == nil {
		err = errors.New("staging received an empty file response")
	}
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
		if err := ctx.Err(); err != nil {
			return prepared, &BetaAgentFilePreparationError{Prepared: prepared, Cause: err}
		}
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
		if uploadErr == nil && upload == nil {
			uploadErr = errors.New("upload received an empty file response")
		}
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

func betaAgentPrepareSelection(ctx context.Context, files map[string]string) ([]betaAgentLocalFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	selected := make([]betaAgentLocalFile, 0, len(files))
	destinations := slices.Sorted(maps.Keys(files))
	for _, destination := range destinations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
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
		selected = append(selected, betaAgentLocalFile{destination: destination, source: source, info: info})
	}
	return selected, nil
}
func betaAgentFileDestination(destination string) error {
	if !utf8.ValidString(destination) {
		return errors.New("destination must be valid UTF-8")
	}
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
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		return nil, err
	}
	_, copyErr := io.Copy(to, response.Body)
	closeErr := response.Body.Close()
	if copyErr != nil {
		return artifact, copyErr
	}
	return artifact, closeErr
}
