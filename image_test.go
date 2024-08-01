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
		Image:          openai.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		Model:          openai.F(openai.ImageNewVariationParamsModelDallE2),
		N:              openai.F(int64(1)),
		ResponseFormat: openai.F(openai.ImageNewVariationParamsResponseFormatURL),
		Size:           openai.F(openai.ImageNewVariationParamsSize1024x1024),
		User:           openai.F("user-1234"),
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
		Image:          openai.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		Prompt:         openai.F("A cute baby sea otter wearing a beret"),
		Mask:           openai.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		Model:          openai.F(openai.ImageEditParamsModelDallE2),
		N:              openai.F(int64(1)),
		ResponseFormat: openai.F(openai.ImageEditParamsResponseFormatURL),
		Size:           openai.F(openai.ImageEditParamsSize1024x1024),
		User:           openai.F("user-1234"),
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
		Prompt:         openai.F("A cute baby sea otter"),
		Model:          openai.F(openai.ImageGenerateParamsModelDallE2),
		N:              openai.F(int64(1)),
		Quality:        openai.F(openai.ImageGenerateParamsQualityStandard),
		ResponseFormat: openai.F(openai.ImageGenerateParamsResponseFormatURL),
		Size:           openai.F(openai.ImageGenerateParamsSize1024x1024),
		Style:          openai.F(openai.ImageGenerateParamsStyleVivid),
		User:           openai.F("user-1234"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
