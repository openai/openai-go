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
		"external body utf16 extension mode": {
			base:       "message/external-body",
			registered: "message/external-body; access-type*=UTF-16BE''%00X%00-%00T%00E%00S%00T; mode=V1",
			response:   "message/external-body; access-type*=utf-16be''%00x%00-%00t%00e%00s%00t; mode=v1",
		},
		"external body unsupported utf32 mode": {
			base:       "message/external-body",
			registered: "message/external-body; access-type*=UTF-32BE''%00%00%00F%00%00%00T%00%00%00P; mode=IMAGE",
			response:   "message/external-body; access-type*=utf-32be''%00%00%00f%00%00%00t%00%00%00p; mode=image",
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
		"invalid charset syntax":  `access-type*="BAD CHAR''FTP"`,
		"invalid language syntax": `access-type*="UTF-8'BAD LANG'FTP"`,
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
		"empty charset and language": {
			registered: "message/external-body; access-type*=''FTP; mode=IMAGE",
			response:   "Message/External-Body; access-type*=''ftp; mode=image",
		},
		"valid language tag": {
			registered: "message/external-body; access-type*=UTF-8'en-US'FTP; mode=IMAGE",
			response:   "Message/External-Body; access-type*=utf-8'EN-us'ftp; mode=image",
		},
		"registered charset alias with period": {
			registered: "message/external-body; access-type*=ANSI_X3.4-1968''FTP; mode=IMAGE",
			response:   "Message/External-Body; access-type*=ansi_x3.4-1968''ftp; mode=image",
		},
		"utf16be value": {
			registered: "message/external-body; access-type*=UTF-16BE''%00F%00T%00P; mode=IMAGE",
			response:   "Message/External-Body; access-type*=utf-16be''%00f%00t%00p; mode=image",
		},
		"utf16be continued value": {
			registered: "message/external-body; access-type*0*=UTF-16BE''%00F%00; access-type*1*=T%00P; mode=IMAGE",
			response:   "Message/External-Body; access-type*0*=utf-16be''%00f%00; access-type*1*=t%00p; mode=image",
		},
		"ebcdic value": {
			registered: "message/external-body; access-type*=IBM037''%C6%E3%D7; mode=IMAGE",
			response:   "Message/External-Body; access-type*=ibm037''%86%A3%97; mode=image",
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

func TestRegisterDecoderExternalBodyAccessTypeKeysDoNotCollide(t *testing.T) {
	tests := map[string]struct {
		first  string
		second string
	}{
		"decoded extended delimiter": {
			first:  "message/external-body;access-type*=UTF-8''FTP%3Bmode%3DIMAGE",
			second: "message/external-body;access-type*=UTF-8''FTP;mode=image",
		},
		"quoted ordinary delimiter": {
			first:  `message/external-body;access-type="FTP;mode=IMAGE"`,
			second: "message/external-body;access-type=FTP;mode=image",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			firstDecoder := &testDecoder{}
			secondDecoder := &testDecoder{}
			RegisterDecoder(test.first, func(io.ReadCloser) Decoder { return firstDecoder })
			RegisterDecoder(test.second, func(io.ReadCloser) Decoder { return secondDecoder })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(test.first))
				delete(decoderTypes, decoderContentTypeKey(test.second))
			})

			if firstKey, secondKey := decoderContentTypeKey(test.first), decoderContentTypeKey(test.second); firstKey == secondKey {
				t.Fatalf("distinct Content-Type values share decoder key %q", firstKey)
			}

			for label, response := range map[string]struct {
				contentType string
				want        Decoder
			}{
				"first":  {contentType: test.first, want: firstDecoder},
				"second": {contentType: test.second, want: secondDecoder},
			} {
				t.Run(label, func(t *testing.T) {
					decoder := NewDecoder(&http.Response{
						Header: http.Header{"Content-Type": {response.contentType}},
						Body:   io.NopCloser(strings.NewReader("")),
					})
					if decoder != response.want {
						t.Fatalf("decoder = %T, want independently registered decoder", decoder)
					}
				})
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldExternalBodyModeForAmbiguousAccessType(t *testing.T) {
	const base = "message/external-body"
	for name, accessTypes := range map[string]string{
		"plain then extended": "access-type=X-TEST; access-type*=UTF-8''FTP",
		"extended then plain": "access-type*=UTF-8''FTP; access-type=X-TEST",
	} {
		t.Run(name, func(t *testing.T) {
			registered := base + "; " + accessTypes + "; mode=IMAGE"
			response := base + "; " + accessTypes + "; mode=image"
			wantSpecific := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(registered))
			})

			if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
				t.Fatalf("ambiguous duplicate access-type folded mode into shared key %q", registeredKey)
			}
			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {response}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder == wantSpecific {
				t.Fatal("ambiguous duplicate access-type selected parameter-specific decoder")
			}
		})
	}
}

func TestRegisterDecoderFoldsCaseInsensitiveExtendedValuesAcrossCharsets(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"charset": {
			registered: `text/plain; charset*=IBM037''%E4%E3%C6%60%F8`,
			response:   `text/plain; charset*=ibm037''%A4%A3%86%60%F8`,
		},
		"external permission": {
			registered: `message/external-body; permission*=IBM037''%D9%C5%C1%C4`,
			response:   `message/external-body; permission*=ibm037''%99%85%81%84`,
		},
		"external mode": {
			registered: `message/external-body; access-type=FTP; mode*=IBM037''%C9%D4%C1%C7%C5`,
			response:   `message/external-body; access-type=ftp; mode*=ibm037''%89%94%81%87%85`,
		},
		"multipart encrypted protocol": {
			registered: `multipart/encrypted; protocol*=IBM037''%C1%D7%D7%D3%C9%C3%C1%E3%C9%D6%D5%61%E3%C5%E2%E3`,
			response:   `multipart/encrypted; protocol*=ibm037''%81%97%97%93%89%83%81%A3%89%96%95%61%A3%85%A2%A3`,
		},
		"multipart signed protocol": {
			registered: `multipart/signed; protocol*=IBM037''%C1%D7%D7%D3%C9%C3%C1%E3%C9%D6%D5%61%E3%C5%E2%E3`,
			response:   `multipart/signed; protocol*=ibm037''%81%97%97%93%89%83%81%A3%89%96%95%61%A3%85%A2%A3`,
		},
		"multipart report type": {
			registered: `multipart/report; report-type*=IBM037''%C4%C5%D3%C9%E5%C5%D9%E8%60%E2%E3%C1%E3%E4%E2`,
			response:   `multipart/report; report-type*=ibm037''%84%85%93%89%A5%85%99%A8%60%A2%A3%81%A3%A4%A2`,
		},
		"multipart related type": {
			registered: `multipart/related; type*=IBM037''%C1%D7%D7%D3%C9%C3%C1%E3%C9%D6%D5%61%E3%C5%E2%E3`,
			response:   `multipart/related; type*=ibm037''%81%97%97%93%89%83%81%A3%89%96%95%61%A3%85%A2%A3`,
		},
		"text plain format": {
			registered: `text/plain; format*=IBM037''%C6%D3%D6%E6%C5%C4`,
			response:   `text/plain; format*=ibm037''%86%93%96%A6%85%84`,
		},
		"text plain delsp": {
			registered: `text/plain; delsp*=IBM037''%E8%C5%E2`,
			response:   `text/plain; delsp*=ibm037''%A8%85%A2`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
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
