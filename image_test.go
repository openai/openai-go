// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
)

func TestImageNewVariationWithOptionalParams(t *testing.T) {
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
	_, err := client.Images.NewVariation(context.TODO(), openai.ImageNewVariationParams{
		Image:          io.Reader(bytes.NewBuffer([]byte("some file contents"))),
		Model:          openai.ImageModelDallE2,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageNewVariationParamsResponseFormatURL,
		Size:           openai.ImageNewVariationParamsSize256x256,
		User:           openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestImageEditWithOptionalParams(t *testing.T) {
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
	_, err := client.Images.Edit(context.TODO(), openai.ImageEditParams{
		Image:          io.Reader(bytes.NewBuffer([]byte("some file contents"))),
		Prompt:         "A cute baby sea otter wearing a beret",
		Mask:           io.Reader(bytes.NewBuffer([]byte("some file contents"))),
		Model:          openai.ImageModelDallE2,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageEditParamsResponseFormatURL,
		Size:           openai.ImageEditParamsSize256x256,
		User:           openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestImageGenerateWithOptionalParams(t *testing.T) {
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
	_, err := client.Images.Generate(context.TODO(), openai.ImageGenerateParams{
		Prompt:         "A cute baby sea otter",
		Model:          openai.ImageModelDallE2,
		N:              openai.Int(1),
		Quality:        openai.ImageGenerateParamsQualityStandard,
		ResponseFormat: openai.ImageGenerateParamsResponseFormatURL,
		Size:           openai.ImageGenerateParamsSize256x256,
		Style:          openai.ImageGenerateParamsStyleVivid,
		User:           openai.String("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
