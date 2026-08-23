package ssestream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterDecoderPreservesCaseSensitiveParameterValues(t *testing.T) {
	const (
		mediaType   = "application/x-openai-go-test-registration-case"
		profileV1   = mediaType + "; profile=\"https://example.com/V1\""
		profileV1Lo = mediaType + "; profile=\"https://example.com/v1\""
	)
	wantDefault := &testDecoder{}
	wantProfile := &testDecoder{}
	RegisterDecoder(mediaType, func(io.ReadCloser) Decoder { return wantDefault })
	RegisterDecoder(profileV1, func(io.ReadCloser) Decoder { return wantProfile })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(mediaType))
		delete(decoderTypes, decoderContentTypeKey(profileV1))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"registered uppercase profile": {
			contentType: "Application/X-OpenAI-Go-Test-Registration-Case; profile=\"https://example.com/V1\"",
			want:        wantProfile,
		},
		"distinct lowercase profile": {
			contentType: profileV1Lo,
			want:        wantDefault,
		},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {test.contentType}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderNormalizesCaseInsensitiveParameterComponents(t *testing.T) {
	const registered = "application/x-openai-go-test-registration-components; Profile=\"https://example.com/V1\"; Charset=UTF-8"
	want := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{
		Header: http.Header{
			"Content-Type": {"Application/X-OpenAI-Go-Test-Registration-Components; profile=\"https://example.com/V1\"; charset=utf-8"},
		},
		Body: io.NopCloser(strings.NewReader("")),
	})
	if decoder != want {
		t.Fatalf("decoder = %T, want registered decoder", decoder)
	}
}

func TestRegisterDecoderNormalizesExternalBodyCaseInsensitiveValues(t *testing.T) {
	for name, test := range map[string]struct {
		parameter string
		registered string
		response   string
	}{
		"access type": {
			parameter:  "Access-Type",
			registered: "LOCAL-FILE",
			response:   "local-file",
		},
		"permission": {
			parameter:  "Permission",
			registered: "READ-WRITE",
			response:   "read-write",
		},
		"mode": {
			parameter:  "Mode",
			registered: "IMAGE",
			response:   "image",
		},
	} {
		t.Run(name, func(t *testing.T) {
			registered := "message/external-body; " + test.parameter + "=" + test.registered
			want := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{
					"Content-Type": {"Message/External-Body; " + strings.ToLower(test.parameter) + "=" + test.response},
				},
				Body: io.NopCloser(strings.NewReader("")),
			})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderExtendedParameterPreservesUnescapedValueCase(t *testing.T) {
	const (
		mediaType = "application/x-openai-go-test-registration-extended"
		variantV1 = mediaType + "; Variant*=ISO-8859-1'EN'caf%E9V1"
		variantv1 = mediaType + "; variant*=iso-8859-1'en'caf%e9v1"
	)
	wantDefault := &testDecoder{}
	wantVariant := &testDecoder{}
	RegisterDecoder(mediaType, func(io.ReadCloser) Decoder { return wantDefault })
	RegisterDecoder(variantV1, func(io.ReadCloser) Decoder { return wantVariant })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(mediaType))
		delete(decoderTypes, decoderContentTypeKey(variantV1))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent encoding case": {
			contentType: mediaType + "; variant*=iso-8859-1'en'caf%e9V1",
			want:        wantVariant,
		},
		"distinct unescaped value case": {
			contentType: variantv1,
			want:        wantDefault,
		},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {test.contentType}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderUnencodedContinuationPreservesValueCase(t *testing.T) {
	const (
		mediaType = "application/x-openai-go-test-registration-continuation"
		upper     = mediaType + "; title*0=V%AB"
		lower     = mediaType + "; title*0=V%ab"
	)
	wantDefault := &testDecoder{}
	wantUpper := &testDecoder{}
	RegisterDecoder(mediaType, func(io.ReadCloser) Decoder { return wantDefault })
	RegisterDecoder(upper, func(io.ReadCloser) Decoder { return wantUpper })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(mediaType))
		delete(decoderTypes, decoderContentTypeKey(upper))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"registered uppercase escape text": {
			contentType: upper,
			want:        wantUpper,
		},
		"distinct lowercase escape text": {
			contentType: lower,
			want:        wantDefault,
		},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {test.contentType}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestDecoderContentTypeKeyDoesNotSplitQuotedSemicolons(t *testing.T) {
	const contentType = "Application/X-OpenAI-Go-Test-Quoted; Profile=\"https://example.com/a;b?x*=V1\"; Charset=UTF-8"
	got := decoderContentTypeKey(contentType)
	want := "application/x-openai-go-test-quoted; profile=\"https://example.com/a;b?x*=V1\"; charset=utf-8"
	if got != want {
		t.Fatalf("decoder content type key = %q, want %q", got, want)
	}
}
