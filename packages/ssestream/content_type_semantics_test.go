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
		"external body extended extension mode": {
			base:       "message/external-body",
			registered: "message/external-body; access-type*=ISO-8859-1''X-TEST; mode=V1",
			response:   "message/external-body; access-type*=iso-8859-1''x-test; mode=v1",
		},
		"smime type": {
			base:       "application/pkcs7-mime",
			registered: "application/pkcs7-mime; smime-type=SIGNED-DATA",
			response:   "application/pkcs7-mime; smime-type=signed-data",
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

func TestRegisterDecoderDoesNotFoldExternalBodyModeForMalformedContinuation(t *testing.T) {
	const base = "message/external-body"
	for name, accessType := range map[string]string{
		"malformed section zero":  "access-type*0*=UTF-8''%ZZ; access-type*1*=FTP",
		"malformed later section": "access-type*0*=UTF-8''FTP; access-type*1*=%ZZ",
		"missing section":         "access-type*0*=UTF-8''FTP; access-type*2*=X",
		"missing metadata":        "access-type*0*=BROKEN; access-type*1*=FTP",
		"empty charset":           "access-type*0*=''FTP; access-type*1*=X",
	} {
		t.Run(name, func(t *testing.T) {
			registered := base + "; " + accessType + "; mode=IMAGE"
			response := base + "; " + accessType + "; mode=image"
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {response}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for malformed continuation", decoder)
			}
		})
	}
}

func TestRegisterDecoderFoldsExternalBodyModeWithUnsupportedExtendedCharset(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"single extended value": {
			registered: "message/external-body; access-type*=ISO-8859-1''FTP; mode=IMAGE",
			response:   "Message/External-Body; access-type*=iso-8859-1''ftp; mode=image",
		},
		"continued extended value": {
			registered: "message/external-body; access-type*0*=ISO-8859-1''ANON%2D; access-type*1*=FTP; mode=IMAGE",
			response:   "Message/External-Body; access-type*0*=iso-8859-1''anon%2d; access-type*1*=ftp; mode=image",
		},
		"quoted extended value": {
			registered: "message/external-body; access-type*=\"ISO-8859-1''TFTP\"; mode=IMAGE",
			response:   "Message/External-Body; access-type*=\"iso-8859-1''tftp\"; mode=image",
		},
		"mode before access type": {
			registered: "message/external-body; mode=IMAGE; access-type*=ISO-8859-1''FTP",
			response:   "Message/External-Body; mode=image; access-type*=iso-8859-1''ftp",
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
