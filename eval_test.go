// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func TestEvalNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evals.New(context.TODO(), openai.EvalNewParams{
		DataSourceConfig: openai.EvalNewParamsDataSourceConfigUnion{
			OfCustom: &openai.EvalNewParamsDataSourceConfigCustom{
				ItemSchema: map[string]interface{}{
					"0":   "bar",
					"1":   "bar",
					"2":   "bar",
					"3":   "bar",
					"4":   "bar",
					"5":   "bar",
					"6":   "bar",
					"7":   "bar",
					"8":   "bar",
					"9":   "bar",
					"10":  "bar",
					"11":  "bar",
					"12":  "bar",
					"13":  "bar",
					"14":  "bar",
					"15":  "bar",
					"16":  "bar",
					"17":  "bar",
					"18":  "bar",
					"19":  "bar",
					"20":  "bar",
					"21":  "bar",
					"22":  "bar",
					"23":  "bar",
					"24":  "bar",
					"25":  "bar",
					"26":  "bar",
					"27":  "bar",
					"28":  "bar",
					"29":  "bar",
					"30":  "bar",
					"31":  "bar",
					"32":  "bar",
					"33":  "bar",
					"34":  "bar",
					"35":  "bar",
					"36":  "bar",
					"37":  "bar",
					"38":  "bar",
					"39":  "bar",
					"40":  "bar",
					"41":  "bar",
					"42":  "bar",
					"43":  "bar",
					"44":  "bar",
					"45":  "bar",
					"46":  "bar",
					"47":  "bar",
					"48":  "bar",
					"49":  "bar",
					"50":  "bar",
					"51":  "bar",
					"52":  "bar",
					"53":  "bar",
					"54":  "bar",
					"55":  "bar",
					"56":  "bar",
					"57":  "bar",
					"58":  "bar",
					"59":  "bar",
					"60":  "bar",
					"61":  "bar",
					"62":  "bar",
					"63":  "bar",
					"64":  "bar",
					"65":  "bar",
					"66":  "bar",
					"67":  "bar",
					"68":  "bar",
					"69":  "bar",
					"70":  "bar",
					"71":  "bar",
					"72":  "bar",
					"73":  "bar",
					"74":  "bar",
					"75":  "bar",
					"76":  "bar",
					"77":  "bar",
					"78":  "bar",
					"79":  "bar",
					"80":  "bar",
					"81":  "bar",
					"82":  "bar",
					"83":  "bar",
					"84":  "bar",
					"85":  "bar",
					"86":  "bar",
					"87":  "bar",
					"88":  "bar",
					"89":  "bar",
					"90":  "bar",
					"91":  "bar",
					"92":  "bar",
					"93":  "bar",
					"94":  "bar",
					"95":  "bar",
					"96":  "bar",
					"97":  "bar",
					"98":  "bar",
					"99":  "bar",
					"100": "bar",
					"101": "bar",
					"102": "bar",
					"103": "bar",
					"104": "bar",
					"105": "bar",
					"106": "bar",
					"107": "bar",
					"108": "bar",
					"109": "bar",
					"110": "bar",
					"111": "bar",
					"112": "bar",
					"113": "bar",
					"114": "bar",
					"115": "bar",
					"116": "bar",
					"117": "bar",
					"118": "bar",
					"119": "bar",
					"120": "bar",
					"121": "bar",
					"122": "bar",
					"123": "bar",
					"124": "bar",
					"125": "bar",
					"126": "bar",
					"127": "bar",
					"128": "bar",
					"129": "bar",
					"130": "bar",
					"131": "bar",
					"132": "bar",
					"133": "bar",
					"134": "bar",
					"135": "bar",
					"136": "bar",
					"137": "bar",
					"138": "bar",
					"139": "bar",
				},
				IncludeSampleSchema: openai.Bool(true),
			},
		},
		TestingCriteria: []openai.EvalNewParamsTestingCriterionUnion{{
			OfLabelModel: &openai.EvalNewParamsTestingCriterionLabelModel{
				Input: []openai.EvalNewParamsTestingCriterionLabelModelInputUnion{{
					OfSimpleInputMessage: &openai.EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage{
						Content: "content",
						Role:    "role",
					},
				}},
				Labels:        []string{"string"},
				Model:         "model",
				Name:          "name",
				PassingLabels: []string{"string"},
			},
		}},
		Metadata: shared.MetadataParam{
			"foo": "string",
		},
		Name:            openai.String("name"),
		ShareWithOpenAI: openai.Bool(true),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestEvalGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evals.Get(context.TODO(), "eval_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestEvalUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evals.Update(
		context.TODO(),
		"eval_id",
		openai.EvalUpdateParams{
			Metadata: shared.MetadataParam{
				"foo": "string",
			},
			Name: openai.String("name"),
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestEvalListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evals.List(context.TODO(), openai.EvalListParams{
		After:   openai.String("after"),
		Limit:   openai.Int(0),
		Order:   openai.EvalListParamsOrderAsc,
		OrderBy: openai.EvalListParamsOrderByCreatedAt,
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestEvalDelete(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Evals.Delete(context.TODO(), "eval_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
