package ssestream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterDecoderDoesNotFoldProtocolDefinedOrExtensionValues(t *testing.T) {
	for name, test := range map[string]struct {
		base       string
		registered string
		response   string
	}{
		"multipart signed micalg": {
			base:       "multipart/signed",
			registered: `multipart/signed; protocol="application/x-test-signature"; micalg=V1`,
			response:   `multipart/signed; protocol="application/x-test-signature"; micalg=v1`,
		},
		"text csv header": {
			base:       "text/csv",
			registered: "text/csv; header=PRESENT",
			response:   "text/csv; header=present",
		},
		"external body extension mode": {
			base:       "message/external-body",
			registered: "message/external-body; access-type=X-TEST; mode=V1",
			response:   "message/external-body; access-type=x-test; mode=v1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(test.base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(test.base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {test.response}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for distinct parameter value", decoder)
			}
		})
	}
}

func TestRegisterDecoderFoldsExternalBodyModeForStandardAccessTypes(t *testing.T) {
	for _, accessType := range []string{"FTP", "ANON-FTP", "TFTP"} {
		t.Run(accessType, func(t *testing.T) {
			registered := "message/external-body; access-type=" + accessType + "; mode=IMAGE"
			want := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{
					"Content-Type": {"Message/External-Body; access-type=" + strings.ToLower(accessType) + "; mode=image"},
				},
				Body: io.NopCloser(strings.NewReader("")),
			})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}
