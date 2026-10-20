package openai_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
)

func TestErrorDiagnosticsTransport(t *testing.T) {
	for _, route := range []string{"responses", "beta", "videos"} {
		for _, shape := range []string{"values", "nulls", "omitted"} {
			t.Run(route+"/"+shape, func(t *testing.T) {
				fields := `,"type":"synthetic_error","param":"input[0]"`
				if shape == "nulls" {
					fields = `,"type":null,"param":null`
				}
				if shape == "omitted" {
					fields = ""
				}
				path := "/v1/responses/resp_synthetic"
				if route == "videos" {
					path = "/v1/videos/video_synthetic"
				}
				client := imageErrorClient(t, "GET", path, `{"error":{"code":"server_error","message":"Synthetic failure"`+fields+`}}`, func(r *http.Request) {
					want := ""
					if route == "beta" {
						want = "true"
					}
					if r.URL.Query().Get("beta") != want {
						t.Fatal("unexpected beta query")
					}
				})
				var kind, parameter, code, message string
				var typeField, paramField respjson.Field
				switch route {
				case "responses":
					result, err := client.Responses.Get(context.Background(), "resp_synthetic", responses.ResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					e := result.Error
					kind, parameter, code, message = e.Type, e.Param, string(e.Code), e.Message
					typeField, paramField = e.JSON.Type, e.JSON.Param
				case "beta":
					result, err := client.Beta.Responses.Get(context.Background(), "resp_synthetic", openai.BetaResponseGetParams{})
					if err != nil {
						t.Fatal(err)
					}
					e := result.Error
					kind, parameter, code, message = e.Type, e.Param, string(e.Code), e.Message
					typeField, paramField = e.JSON.Type, e.JSON.Param
				case "videos":
					result, err := client.Videos.Get(context.Background(), "video_synthetic")
					if err != nil {
						t.Fatal(err)
					}
					e := result.Error
					kind, parameter, code, message = e.Type, e.Param, e.Code, e.Message
					typeField, paramField = e.JSON.Type, e.JSON.Param
				}
				wantType, wantParam, rawType, rawParam := "", "", "", ""
				if shape == "values" {
					wantType, wantParam, rawType, rawParam = "synthetic_error", "input[0]", `"synthetic_error"`, `"input[0]"`
				}
				if shape == "nulls" {
					rawType, rawParam = "null", "null"
				}
				if kind != wantType {
					t.Errorf("type = %q, want %q", kind, wantType)
				}
				if parameter != wantParam {
					t.Errorf("param = %q, want %q", parameter, wantParam)
				}
				if typeField.Raw() != rawType || paramField.Raw() != rawParam {
					t.Error("lost null/omission distinction")
				}
				if code != "server_error" || message != "Synthetic failure" {
					t.Error("existing error fields changed")
				}
			})
		}
	}
}
