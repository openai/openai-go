package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type emptyPathMarshaler struct {
	calls *int
	err   error
}

func (v emptyPathMarshaler) MarshalJSON() ([]byte, error) {
	*v.calls++
	if v.err != nil {
		return nil, v.err
	}
	return []byte(`"encoded"`), nil
}

func TestJSONOptionsEmptyPathContract(t *testing.T) {
	for _, deleteValue := range []bool{false, true} {
		t.Run(map[bool]string{false: "set", true: "delete"}[deleteValue], func(t *testing.T) {
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://api.openai.com/v1"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
					t.Fatal("empty-path option must fail before sending a request")
					return nil, errors.New("unexpected request")
				}},
			}))
			opt := option.WithJSONSet("", "replacement")
			if deleteValue {
				opt = option.WithJSONDel("")
			}
			err := client.Post(context.Background(), "/test", json.RawMessage(`{"":"keep","z":true}`), nil, opt)
			if err == nil || !strings.Contains(err.Error(), "path cannot be empty") {
				t.Fatalf("expected legacy empty-path error, got %v", err)
			}
		})
	}
}

func TestJSONSetEmptyPathMarshalOrder(t *testing.T) {
	sentinel := errors.New("custom marshal failure")
	for _, marshalErr := range []error{nil, sentinel} {
		calls := 0
		client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://api.openai.com/v1"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				t.Fatal("invalid option must not send a request")
				return nil, errors.New("unexpected request")
			}},
		}))
		err := client.Post(context.Background(), "/test", json.RawMessage(`{}`), nil,
			option.WithJSONSet("", emptyPathMarshaler{&calls, marshalErr}))
		if calls != 1 || err == nil {
			t.Fatalf("marshal count=%d, error=%v", calls, err)
		}
		if marshalErr != nil {
			if !errors.Is(err, sentinel) {
				t.Fatalf("custom marshal error must take precedence: %v", err)
			}
		} else if !strings.Contains(err.Error(), "path cannot be empty") {
			t.Fatalf("expected empty-path error after successful marshaling: %v", err)
		}
	}
}

func TestJSONSetRejectsNonfiniteValues(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://api.openai.com/v1"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				t.Fatal("nonfinite values must not be sent as quoted strings")
				return nil, errors.New("unexpected request")
			}},
		}))
		err := client.Post(context.Background(), "/test", json.RawMessage(`{}`), nil, option.WithJSONSet("value", value))
		var unsupported *json.UnsupportedValueError
		if !errors.As(err, &unsupported) {
			t.Fatalf("expected unsupported JSON value error, got %v", err)
		}
	}
}

func TestJSONSetLiteralPunctuationKey(t *testing.T) {
	var body string
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://api.openai.com/v1"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{
		Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			data, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			body = string(data)
			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		}},
	}))
	err := client.Post(context.Background(), "/test", json.RawMessage(`{"a.b":"old","keep":true}`), nil, option.WithJSONSet(`a\.b`, "new"))
	if err != nil || body != `{"a.b":"new","keep":true}` {
		t.Fatalf("literal key changed: %s / %v", body, err)
	}
}
