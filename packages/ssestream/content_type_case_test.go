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
		parameter  string
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

func TestRegisterDecoderNormalizesLogicalCaseInsensitiveParameters(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"text plain format": {
			registered: "text/plain; Format=FLOWED",
			response:   "Text/Plain; format=flowed",
		},
		"text plain delsp": {
			registered: "text/plain; DelSP=YES",
			response:   "Text/Plain; delsp=yes",
		},
		"encoded charset": {
			registered: "text/plain; Charset*=US-ASCII'EN'%55TF-8",
			response:   "Text/Plain; charset*=us-ascii'en'%75tf-8",
		},
		"quoted encoded charset": {
			registered: "text/plain; Charset*=\"US-ASCII'EN'UTF%2D8\"",
			response:   "Text/Plain; charset*=\"us-ascii'en'utf%2d8\"",
		},
		"multipart related type": {
			registered: "multipart/related; Type=\"Application/X-Test\"",
			response:   "Multipart/Related; type=\"application/x-test\"",
		},
		"multipart signed protocol": {
			registered: "multipart/signed; Protocol=\"Application/PGP-Signature\"",
			response:   "Multipart/Signed; protocol=\"application/pgp-signature\"",
		},
		"multipart signed micalg": {
			registered: "multipart/signed; Micalg=PGP-SHA256",
			response:   "Multipart/Signed; micalg=pgp-sha256",
		},
		"multipart encrypted protocol": {
			registered: "multipart/encrypted; Protocol=\"Application/PGP-Encrypted\"",
			response:   "Multipart/Encrypted; protocol=\"application/pgp-encrypted\"",
		},
		"multipart report type": {
			registered: "multipart/report; Report-Type=DELIVERY-STATUS",
			response:   "Multipart/Report; report-type=delivery-status",
		},
		"text csv header": {
			registered: "text/csv; Header=PRESENT",
			response:   "Text/CSV; header=present",
		},
		"unencoded format continuation": {
			registered: "text/plain; Format*0=FLO; Format*1=WED",
			response:   "Text/Plain; format*0=flo; format*1=wed",
		},
		"encoded access type continuation": {
			registered: "message/external-body; Access-Type*0*=US-ASCII''LOCAL-; Access-Type*1*=FILE",
			response:   "Message/External-Body; access-type*0*=us-ascii''local-; access-type*1*=file",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {test.response}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestDecoderContentTypeKeyKeepsContextualValuesCaseSensitive(t *testing.T) {
	for name, test := range map[string]struct {
		contentType string
		want        string
	}{
		"format outside text plain": {
			contentType: "text/html; Format=FLOWED",
			want:        "text/html; format=FLOWED",
		},
		"report type outside multipart report": {
			contentType: "multipart/mixed; Report-Type=DELIVERY-STATUS",
			want:        "multipart/mixed; report-type=DELIVERY-STATUS",
		},
		"header outside text csv": {
			contentType: "text/plain; Header=PRESENT",
			want:        "text/plain; header=PRESENT",
		},
		"micalg outside multipart signed": {
			contentType: "multipart/encrypted; Micalg=PGP-SHA256",
			want:        "multipart/encrypted; micalg=PGP-SHA256",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := decoderContentTypeKey(test.contentType); got != test.want {
				t.Fatalf("decoder content type key = %q, want %q", got, test.want)
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

func TestRegisterDecoderEncodedContinuationPreservesLaterSegmentValueCase(t *testing.T) {
	const (
		mediaType = "application/x-openai-go-test-registration-encoded-continuation"
		titleV1   = mediaType + "; title*0*=us-ascii''prefix; title*1*=Bob%2D's'V1"
		titlev1   = mediaType + "; title*0*=us-ascii''prefix; title*1*=bob%2d's'V1"
	)
	wantDefault := &testDecoder{}
	wantTitle := &testDecoder{}
	RegisterDecoder(mediaType, func(io.ReadCloser) Decoder { return wantDefault })
	RegisterDecoder(titleV1, func(io.ReadCloser) Decoder { return wantTitle })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(mediaType))
		delete(decoderTypes, decoderContentTypeKey(titleV1))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent percent encoding case": {
			contentType: mediaType + "; title*0*=us-ascii''prefix; title*1*=Bob%2d's'V1",
			want:        wantTitle,
		},
		"distinct later segment data case": {
			contentType: titlev1,
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
