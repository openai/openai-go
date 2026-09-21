package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// executeAudioTextRequest preserves the default audio response types while
// accepting plaintext in their Text field. Explicit response destinations retain
// the request layer's behavior and body ownership.
func executeAudioTextRequest(ctx context.Context, path string, body any, dst any, setText func(string), opts ...option.RequestOption) error {
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, dst, opts...)
	if err != nil {
		return err
	}
	if cfg.ResponseBodyInto != dst {
		return cfg.Execute()
	}

	var contents []byte
	cfg.ResponseBodyInto = &contents
	if cfg.ResponseInto == nil {
		cfg.ResponseInto = new(*http.Response)
	}
	if err := cfg.Execute(); err != nil {
		return err
	}

	contentType := (*cfg.ResponseInto).Header.Get("Content-Type")
	mediaType, _, mediaErr := mime.ParseMediaType(contentType)
	if mediaErr == nil && mediaType == "text/plain" {
		setText(string(contents))
		return nil
	}
	if !strings.Contains(mediaType, "application/json") && !strings.HasSuffix(mediaType, "+json") {
		return fmt.Errorf("expected destination type of 'string' or '[]byte' for responses with content-type '%s' that is not 'application/json'", contentType)
	}
	if err := json.NewDecoder(bytes.NewReader(contents)).Decode(dst); err != nil {
		return fmt.Errorf("error parsing response json: %w", err)
	}
	return nil
}
