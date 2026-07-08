package ssestream

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestStreamSkipsEventsWithoutData(t *testing.T) {
	tests := map[string]struct {
		body string
		want []string
	}{
		"comment": {
			body: ": OPENROUTER PROCESSING\n\ndata: {\"value\":\"first\"}\n\ndata: [DONE]\n\n",
			want: []string{"first"},
		},
		"retry directive": {
			body: "retry: 3000\n\ndata: {\"value\":\"first\"}\n\ndata: [DONE]\n\n",
			want: []string{"first"},
		},
		"multiple empty events": {
			body: "data: {\"value\":\"first\"}\n\n: keep-alive\n\nretry: 3000\n\ndata: {\"value\":\"second\"}\n\ndata: [DONE]\n\n",
			want: []string{"first", "second"},
		},
		"only empty events": {
			body: ": keep-alive\n\nretry: 3000\n\n",
			want: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := &http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(test.body)),
			}
			stream := NewStream[struct {
				Value string `json:"value"`
			}](NewDecoder(res), nil)

			var got []string
			for stream.Next() {
				got = append(got, stream.Current().Value)
			}

			if err := stream.Err(); err != nil {
				t.Fatalf("stream ended with error: %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("received values %v, want %v", got, test.want)
			}
		})
	}
}
