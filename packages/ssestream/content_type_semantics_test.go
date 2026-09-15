package ssestream

import (
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestRegisterDecoderDoesNotFoldProtocolDefinedOrExtensionValues(t *testing.T) {
	for name, test := range map[string]struct {
		base       string
		registered string
		response   string
	}{
		"multipart signed protocol-defined micalg": {
			base: "multipart/signed",
			// RFC 1847 delegates micalg value semantics to the selected protocol.
			// This extension protocol intentionally treats V1 and v1 as distinct.
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

func TestRegisterDecoderFoldsExternalBodyModeForExtendedAccessTypes(t *testing.T) {
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
			registered: "message/external-body; access-type*=\"ISO-8859-1''TFTP\"; mode=NETASCII",
			response:   "Message/External-Body; access-type*=\"iso-8859-1''tftp\"; mode=netascii",
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
		"utf32be value": {
			registered: "message/external-body; access-type*=UTF-32BE''%00%00%00F%00%00%00T%00%00%00P; mode=IMAGE",
			response:   "Message/External-Body; access-type*=utf-32be''%00%00%00f%00%00%00t%00%00%00p; mode=image",
		},
		"utf32le value": {
			registered: "message/external-body; access-type*=UTF-32LE''F%00%00%00T%00%00%00P%00%00%00; mode=IMAGE",
			response:   "Message/External-Body; access-type*=utf-32le''f%00%00%00t%00%00%00p%00%00%00; mode=image",
		},
		"utf32 value with bom": {
			registered: "message/external-body; access-type*=UTF-32''%00%00%FE%FF%00%00%00F%00%00%00T%00%00%00P; mode=IMAGE",
			response:   "Message/External-Body; access-type*=utf-32''%00%00%fe%ff%00%00%00f%00%00%00t%00%00%00p; mode=image",
		},
		"utf32be iana alias": {
			registered: "message/external-body; access-type*=csUTF32BE''%00%00%00F%00%00%00T%00%00%00P; mode=IMAGE",
			response:   "Message/External-Body; access-type*=csutf32be''%00%00%00f%00%00%00t%00%00%00p; mode=image",
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

func TestRegisterDecoderUsesExtendedExternalBodyAccessTypeOverPlainFallback(t *testing.T) {
	const base = "message/external-body"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"plain then extended": {
			registered: "access-type=X-TEST; access-type*=UTF-8''FTP",
			response:   "access-type=x-test; access-type*=utf-8''ftp",
		},
		"extended then plain": {
			registered: "access-type*=UTF-8''FTP; access-type=X-TEST",
			response:   "access-type*=utf-8''ftp; access-type=x-test",
		},
	} {
		t.Run(name, func(t *testing.T) {
			registered := base + "; " + test.registered + "; mode=IMAGE"
			response := base + "; " + test.response + "; mode=image"
			want := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })

			decoder := NewDecoder(&http.Response{
				Header: http.Header{"Content-Type": {response}},
				Body:   io.NopCloser(strings.NewReader("")),
			})
			if decoder != want {
				t.Fatalf("decoder = %T, want extended access-type registration", decoder)
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
	for name, test := range map[string]struct {
		accessType string
		mode       string
	}{
		"ftp image":       {accessType: "FTP", mode: "IMAGE"},
		"ftp local":       {accessType: "FTP", mode: "LOCAL8"},
		"anon ftp ebcdic": {accessType: "ANON-FTP", mode: "EBCDIC"},
		"tftp netascii":   {accessType: "TFTP", mode: "NETASCII"},
		"tftp octet":      {accessType: "TFTP", mode: "OCTET"},
		"tftp mail":       {accessType: "TFTP", mode: "MAIL"},
	} {
		t.Run(name, func(t *testing.T) {
			registered := "message/external-body; access-type=" + test.accessType + "; mode=" + test.mode
			want := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(registered))
			})

			decoder := NewDecoder(&http.Response{
				Header: http.Header{
					"Content-Type": {"Message/External-Body; access-type=" + strings.ToLower(test.accessType) + "; mode=" + strings.ToLower(test.mode)},
				},
				Body: io.NopCloser(strings.NewReader("")),
			})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderRejectsMalformedUnquotedExtendedValueCollision(t *testing.T) {
	registered := `multipart/signed; protocol*="UTF-8''application/pgp-signature"`
	malformed := `multipart/signed; protocol*=UTF-8''application/pgp-signature`
	want := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })

	if registeredKey, malformedKey := decoderContentTypeKey(registered), decoderContentTypeKey(malformed); registeredKey == malformedKey {
		t.Fatalf("malformed unquoted extended value shares decoder key %q", registeredKey)
	}
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {malformed}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder == want {
		t.Fatal("malformed unquoted extended value selected parameter-specific decoder")
	}
}

func TestRegisterDecoderFoldsTextCalendarMIMEParameters(t *testing.T) {
	for name, test := range map[string]struct{ registered, response string }{
		"component": {
			registered: "text/calendar; component=VEVENT",
			response:   "Text/Calendar; component=vevent",
		},
		"method": {
			registered: "text/calendar; method=REQUEST",
			response:   "Text/Calendar; method=request",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderDecodesMixedRFC2231ContinuationSegments(t *testing.T) {
	registered := "text/plain; charset*0*=UTF-16BE''%00U; charset*1=TF-8"
	response := "Text/Plain; charset*0*=utf-16be''%00u; charset*1=tf-8"
	want := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != want {
		t.Fatalf("decoder = %T, want registered decoder", decoder)
	}
}

func TestRegisterDecoderPreservesRFC2231LanguageIdentity(t *testing.T) {
	en := "text/plain; format*=UTF-8'en'FLOWED"
	fr := "text/plain; format*=UTF-8'fr'FLOWED"
	enResponse := "Text/Plain; format*=utf-8'EN'flowed"
	wantEN := &testDecoder{}
	wantFR := &testDecoder{}
	RegisterDecoder(en, func(io.ReadCloser) Decoder { return wantEN })
	RegisterDecoder(fr, func(io.ReadCloser) Decoder { return wantFR })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(en))
		delete(decoderTypes, decoderContentTypeKey(fr))
	})
	if enKey, frKey := decoderContentTypeKey(en), decoderContentTypeKey(fr); enKey == frKey {
		t.Fatalf("distinct language tags share decoder key %q", enKey)
	}
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {enResponse}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantEN {
		t.Fatalf("decoder = %T, want English registration", decoder)
	}
}

func TestRegisterDecoderDoesNotFoldInvalidUnicodeIntoMIMEProtocol(t *testing.T) {
	registered := "multipart/signed; protocol*=UTF-8''application%2FK"
	malformed := "multipart/signed; protocol*=UTF-8''application%2F%E2%84%AA"
	want := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })
	if registeredKey, malformedKey := decoderContentTypeKey(registered), decoderContentTypeKey(malformed); registeredKey == malformedKey {
		t.Fatalf("invalid Unicode protocol shares valid MIME decoder key %q", registeredKey)
	}
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {malformed}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder == want {
		t.Fatal("invalid Unicode protocol selected valid parameter-specific decoder")
	}
}

func TestRegisterDecoderFoldsPrivateCharsetTokens(t *testing.T) {
	for name, test := range map[string]struct{ registered, response string }{
		"ordinary": {
			registered: "text/plain; charset=X-OPENAI-TEST",
			response:   "Text/Plain; charset=x-openai-test",
		},
		"extended value": {
			registered: "text/plain; charset*=UTF-8''X-OPENAI-TEST",
			response:   "Text/Plain; charset*=utf-8''x-openai-test",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldInvalidRFC1766DigitSubtag(t *testing.T) {
	const base = "text/plain"
	registered := "text/plain; format*=UTF-8'x-1'FLOWED"
	response := "Text/Plain; format*=utf-8'x-1'flowed"
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for invalid RFC 1766 language tag", decoder)
	}
}

func TestRegisterDecoderDoesNotTreatRFC2231LeadingZeroSectionsAsContinuations(t *testing.T) {
	registered := "text/plain; format*00=FLOW; format*01=ED"
	response := "Text/Plain; format*00=flow; format*01=ed"
	wantSpecific := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })

	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("invalid leading-zero RFC 2231 sections share decoder key %q", registeredKey)
	}
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder == wantSpecific {
		t.Fatal("invalid leading-zero RFC 2231 sections selected parameter-specific decoder")
	}
}

func TestRegisterDecoderFoldsExternalBodyExpirationDateTokens(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"ordinary": {
			registered: `message/external-body; access-type=FTP; expiration="Fri, 14 Jun 2024 12:00:00 GMT"`,
			response:   `Message/External-Body; access-type=ftp; expiration="fri, 14 jun 2024 12:00:00 gmt"`,
		},
		"extended": {
			registered: `message/external-body; access-type=FTP; expiration*=UTF-8''Fri%2C%2014%20Jun%202024%2012%3A00%3A00%20GMT`,
			response:   `Message/External-Body; access-type=ftp; expiration*=utf-8''fri%2c%2014%20jun%202024%2012%3a00%3a00%20gmt`,
		},
		"military zone ordinary": {
			registered: `message/external-body; access-type=FTP; expiration="Fri, 14 Jun 2024 12:00:00 Z"`,
			response:   `Message/External-Body; access-type=ftp; expiration="fri, 14 jun 2024 12:00:00 z"`,
		},
		"military zone extended": {
			registered: `message/external-body; access-type=FTP; expiration*=UTF-8''Fri%2C%2014%20Jun%202024%2012%3A00%3A00%20A`,
			response:   `Message/External-Body; access-type=ftp; expiration*=utf-8''fri%2c%2014%20jun%202024%2012%3a00%3a00%20a`,
		},
		"military zone with trailing comment": {
			registered: `message/external-body; access-type=FTP; expiration="14 Jun 2024 12:00:00 A (NOTE)"`,
			response:   `Message/External-Body; access-type=ftp; expiration="14 jun 2024 12:00:00 a (note)"`,
		},
		"trailing comment containing comma": {
			registered: `message/external-body; access-type=FTP; expiration="14 Jun 2024 12:00:00 GMT (NOTE, EXTRA)"`,
			response:   `Message/External-Body; access-type=ftp; expiration="14 jun 2024 12:00:00 gmt (note, extra)"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldInvalidExternalBodyExpiration(t *testing.T) {
	const base = "message/external-body"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"not a date": {
			registered: `message/external-body; access-type=FTP; expiration="NOT-A-DATE"`,
			response:   `Message/External-Body; access-type=ftp; expiration="not-a-date"`,
		},
		"wrong weekday": {
			registered: `message/external-body; access-type=FTP; expiration="Thu, 14 Jun 2024 12:00:00 GMT"`,
			response:   `Message/External-Body; access-type=ftp; expiration="thu, 14 jun 2024 12:00:00 gmt"`,
		},
		"unused military zone": {
			registered: `message/external-body; access-type=FTP; expiration="Fri, 14 Jun 2024 12:00:00 J"`,
			response:   `Message/External-Body; access-type=ftp; expiration="fri, 14 jun 2024 12:00:00 j"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for invalid expiration date", decoder)
			}
		})
	}
}

func TestRegisterDecoderFoldsAdditionalStandardMIMEParameters(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"h264 svc profile-level-id": {
			registered: "video/H264-SVC; profile-level-id=42E01F",
			response:   "Video/H264-SVC; profile-level-id=42e01f",
		},
		"h264 svc extended profile-level-id": {
			registered: "video/H264-SVC; profile-level-id*=UTF-8''42E01F",
			response:   "Video/H264-SVC; profile-level-id*=utf-8''42e01f",
		},
		"text directory profile": {
			registered: "text/directory; profile=VCARD",
			response:   "Text/Directory; profile=vcard",
		},
		"text directory extended profile": {
			registered: "text/directory; profile*=UTF-8''VCARD",
			response:   "Text/Directory; profile*=utf-8''vcard",
		},
		"xop type": {
			registered: `application/xop+xml; type="APPLICATION/SOAP+XML"`,
			response:   `Application/Xop+Xml; type="application/soap+xml"`,
		},
		"xop extended type": {
			registered: "application/xop+xml; type*=UTF-8''APPLICATION%2FSOAP%2BXML",
			response:   "Application/Xop+Xml; type*=utf-8''application%2fsoap%2bxml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderXOPTypePreservesNestedParameterValueCase(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="APPLICATION/SOAP+XML; ACTION=\"https://example.com/V1\""`
	equivalent := `Application/Xop+Xml; type="application/soap+xml; action=\"https://example.com/V1\""`
	distinct := `Application/Xop+Xml; type="application/soap+xml; action=\"https://example.com/v1\""`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent nested media type":    {contentType: equivalent, want: wantSpecific},
		"distinct nested parameter value": {contentType: distinct, want: wantBare},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.contentType}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want %T", decoder, test.want)
			}
		})
	}
}

func TestDecoderContentTypeKeyCachesDuplicateSemanticValidation(t *testing.T) {
	contentType := "message/external-body; access-type=FTP" + strings.Repeat("; expiration=x", (256<<10)/14)
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(contentType)
	})
	if allocs > 32 {
		t.Fatalf("duplicate semantic parameter allocations = %.0f, want <= 32", allocs)
	}
}

func TestDecoderContentTypeKeyUsesCompactContinuationState(t *testing.T) {
	var contentType strings.Builder
	contentType.Grow(256 << 10)
	contentType.WriteString("text/plain; charset*0=x")
	for i := 1; contentType.Len() < 256<<10; i++ {
		contentType.WriteString("; charset*")
		contentType.WriteString(strconv.Itoa(i))
		contentType.WriteString("=x")
	}
	value := contentType.String()
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(value)
	})
	if allocs > 128 {
		t.Fatalf("RFC 2231 continuation allocations = %.0f, want <= 128", allocs)
	}
}

func TestDecoderContentTypeKeyReusesCharsetDecoderAcrossAlternatingRFC2231Runs(t *testing.T) {
	var contentType strings.Builder
	contentType.Grow(256 << 10)
	contentType.WriteString("text/plain; charset*0*=UTF-8''x")
	for i := 1; contentType.Len() < 256<<10; i++ {
		contentType.WriteString("; charset*")
		contentType.WriteString(strconv.Itoa(i))
		if i%2 == 0 {
			contentType.WriteByte('*')
		}
		contentType.WriteString("=x")
	}
	value := contentType.String()
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(value)
	})
	if allocs > 256 {
		t.Fatalf("alternating RFC 2231 charset-run allocations = %.0f, want <= 256", allocs)
	}
}

func TestDecoderContentTypeKeyWritesExtendedValuesWithoutPerParameterAllocations(t *testing.T) {
	var contentType strings.Builder
	contentType.Grow(256 << 10)
	contentType.WriteString("application/x-test")
	for i := 0; contentType.Len() < 256<<10; i++ {
		contentType.WriteString("; a")
		contentType.WriteString(strconv.Itoa(i))
		contentType.WriteString("*=UTF-8''x")
	}
	value := contentType.String()
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(value)
	})
	if allocs > 32 {
		t.Fatalf("extended parameter allocations = %.0f, want <= 32", allocs)
	}
}

func TestDecoderContentTypeKeyAvoidsHexExpansionForUndecodableExtendedValue(t *testing.T) {
	payload := strings.Repeat("A", 1<<20)
	contentType := "text/plain; charset*=" + payload
	key := decoderContentTypeKey(contentType)
	if len(key) > len(contentType)+128 {
		t.Fatalf("decoder key expanded from %d to %d bytes", len(contentType), len(key))
	}
}

func TestDecoderContentTypeKeyAvoidsHexExpansionForDecodedExtendedValue(t *testing.T) {
	payload := strings.Repeat("A", 1<<20)
	contentType := "text/plain; charset*=UTF-8''" + payload
	key := decoderContentTypeKey(contentType)
	if len(key) > len(contentType)+128 {
		t.Fatalf("decoded extended decoder key expanded from %d to %d bytes", len(contentType), len(key))
	}
}

func TestDecoderContentTypeKeyEmitsDecodedExtendedValueOnceAcrossFallbacks(t *testing.T) {
	payload := strings.Repeat("A", 4<<10)
	var contentType strings.Builder
	contentType.WriteString("text/plain; charset*=UTF-8''")
	contentType.WriteString(payload)
	for i := 0; i < 128; i++ {
		contentType.WriteString("; charset=x")
	}
	value := contentType.String()
	key := decoderContentTypeKey(value)
	if len(key) > len(value)*4 {
		t.Fatalf("decoded fallback decoder key expanded from %d to %d bytes", len(value), len(key))
	}
}

func TestParsePlainMediaParameterValueAvoidsSyntheticMediaTypeAllocations(t *testing.T) {
	value := `"` + strings.Repeat("A", 256<<10) + `"`
	var decoded string
	var ok bool
	allocs := testing.AllocsPerRun(5, func() {
		decoded, ok = parsePlainMediaParameterValue(value)
	})
	if !ok || len(decoded) != len(value)-2 {
		t.Fatal("large quoted parameter did not parse")
	}
	if allocs > 1 {
		t.Fatalf("plain parameter allocations = %.0f, want <= 1", allocs)
	}
}

func TestParsePlainMediaParameterValueMatchesMIMEParser(t *testing.T) {
	for _, value := range []string{
		"token",
		`"quoted;value"`,
		`"escaped\\\"quote"`,
		`"C:\\dev\\go\\foo.txt"`,
		`""`,
		"BAD/VALUE",
		`"unterminated`,
		"\"line\nbreak\"",
	} {
		t.Run(value, func(t *testing.T) {
			_, params, err := mime.ParseMediaType("application/octet-stream; x=" + value)
			want, wantOK := params["x"]
			if err != nil {
				wantOK = false
			}
			got, gotOK := parsePlainMediaParameterValue(value)
			if gotOK != wantOK || got != want {
				t.Fatalf("parsePlainMediaParameterValue(%q) = %q, %v; want %q, %v", value, got, gotOK, want, wantOK)
			}
		})
	}
}

func TestWriteEncodedDecoderKeyValueIsLengthDelimited(t *testing.T) {
	var first strings.Builder
	writeEncodedDecoderKeyValue(&first, 'r', "a")
	first.WriteString("bc")

	var second strings.Builder
	writeEncodedDecoderKeyValue(&second, 'r', "ab")
	second.WriteString("c")

	if first.String() == second.String() {
		t.Fatalf("length-delimited values collided at %q", first.String())
	}
}

func TestWriteDecodedDecoderKeyValueIsLengthDelimited(t *testing.T) {
	var first strings.Builder
	writeDecodedDecoderKeyValue(&first, "a", "bc")

	var second strings.Builder
	writeDecodedDecoderKeyValue(&second, "ab", "c")

	if first.String() == second.String() {
		t.Fatalf("length-delimited decoded fields collided at %q", first.String())
	}
}

func TestDecoderContentTypeKeyReusesBufferForTinyParameters(t *testing.T) {
	contentType := "text/event-stream" + strings.Repeat(";a=", (256<<10)/3)
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(contentType)
	})
	if allocs > 16 {
		t.Fatalf("decoderContentTypeKey allocations = %.0f, want <= 16", allocs)
	}
}

func TestNormalizeXOPTypeCanonicalizesOnce(t *testing.T) {
	value := largeXOPNestedMediaType(256 << 10)
	canonicalAllocs := testing.AllocsPerRun(1, func() {
		if _, ok := canonicalMIMEMediaTypeValue(value); !ok {
			panic("canonicalization failed")
		}
	})
	normalizedAllocs := testing.AllocsPerRun(1, func() {
		if _, ok := normalizeCaseInsensitiveMediaParameterValue("application/xop+xml", "type", value, ""); !ok {
			panic("normalization failed")
		}
	})
	if normalizedAllocs > canonicalAllocs*1.25+128 {
		t.Fatalf("XOP normalization allocations = %.0f, single canonicalization = %.0f", normalizedAllocs, canonicalAllocs)
	}
}

func BenchmarkNormalizeLargeXOPType(b *testing.B) {
	value := largeXOPNestedMediaType(2 << 20)
	b.ReportAllocs()
	b.SetBytes(int64(len(value)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := normalizeCaseInsensitiveMediaParameterValue("application/xop+xml", "type", value, ""); !ok {
			b.Fatal("normalization failed")
		}
	}
}

func largeXOPNestedMediaType(minBytes int) string {
	var nested strings.Builder
	nested.Grow(minBytes)
	nested.WriteString("application/soap+xml")
	for i := 0; nested.Len() < minBytes; i++ {
		nested.WriteString("; p")
		nested.WriteString(strconv.Itoa(i))
		nested.WriteString("=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
	return nested.String()
}

func TestRegisterDecoderNormalizesEqualDuplicateCaseInsensitiveParameters(t *testing.T) {
	const base = "application/x-openai-go-equal-duplicates"
	registered := base + "; charset=UTF-8; charset=UTF-8"
	response := "Application/X-OpenAI-Go-Equal-Duplicates; charset=utf-8; charset=utf-8"
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantSpecific {
		t.Fatalf("decoder = %T, want equal-duplicate registration", decoder)
	}
}

func TestRegisterDecoderNormalizesEqualDuplicateExtendedParameters(t *testing.T) {
	const base = "text/plain"
	registered := `text/plain; charset*=UTF-8''UTF-8; charset*="UTF-8''UTF-8"`
	response := `Text/Plain; charset*=utf-8''utf-8; charset*="utf-8''utf-8"`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantSpecific {
		t.Fatalf("decoder = %T, want equal extended-duplicate registration", decoder)
	}
}

func TestRegisterDecoderPreservesConflictingDuplicateExtendedParameters(t *testing.T) {
	registered := `text/plain; charset*=UTF-8''UTF-8; charset*=UTF-8''ISO-8859-1`
	response := `Text/Plain; charset*=utf-8''utf-8; charset*=utf-8''iso-8859-1`
	if decoderContentTypeKey(registered) == decoderContentTypeKey(response) {
		t.Fatal("conflicting extended duplicate values collapsed")
	}
}

func TestRegisterDecoderCanonicalizesContextualParametersInsideXOPType(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format=FLOWED"`
	response := `Application/Xop+Xml; type="Text/Plain; format=flowed"`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantSpecific {
		t.Fatalf("decoder = %T, want nested-context registration", decoder)
	}

	sensitiveUpper := `application/xop+xml; type="application/soap+xml; action=\"https://example.com/V1\""`
	sensitiveLower := `application/xop+xml; type="application/soap+xml; action=\"https://example.com/v1\""`
	if decoderContentTypeKey(sensitiveUpper) == decoderContentTypeKey(sensitiveLower) {
		t.Fatal("nested case-sensitive action values collapsed")
	}
}

func TestRegisterDecoderXOPTypeRejectsGappedRFC2231Continuations(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format*0*=UTF-8''FLOWED; format*2*=X"`
	response := `Application/Xop+Xml; type="text/plain; format*0*=UTF-8''FLOWED; format*2*=Y"`
	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("gapped XOP continuations share decoder key %q", registeredKey)
	}

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for gapped XOP continuation", decoder)
	}
}

func TestRegisterDecoderXOPTypeAcceptsContiguousRFC2231Continuations(t *testing.T) {
	registered := `application/xop+xml; type="text/plain; format*0*=UTF-8''FLO; format*1*=WED"`
	response := `Application/Xop+Xml; type="Text/Plain; format*0*=utf-8''flo; format*1*=wed"`
	want := &testDecoder{}
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
	t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != want {
		t.Fatalf("decoder = %T, want registered decoder for contiguous XOP continuation", decoder)
	}
}

func TestExternalBodyAccessTypeAvoidsSecondFullMediaTypeParse(t *testing.T) {
	var contentType strings.Builder
	contentType.WriteString("message/external-body; access-type=FTP")
	for i := 0; i < 1000; i++ {
		contentType.WriteString(";x")
		contentType.WriteString(strconv.Itoa(i))
		contentType.WriteString("=1")
	}
	value := contentType.String()
	baseline := testing.AllocsPerRun(3, func() {
		_, _, _ = mime.ParseMediaType(value)
	})
	got := testing.AllocsPerRun(3, func() {
		_, _ = decoderContentTypes(value)
	})
	if got > baseline+32 {
		t.Fatalf("decoderContentTypes allocations = %.0f, single ParseMediaType = %.0f; want no second full parse", got, baseline)
	}
}

func TestRegisterDecoderFoldsH264ProfileLevelID(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"ordinary": {
			registered: "video/H264; profile-level-id=42E01F",
			response:   "Video/H264; profile-level-id=42e01f",
		},
		"extended": {
			registered: "video/H264; profile-level-id*=UTF-8''42E01F",
			response:   "Video/H264; profile-level-id*=utf-8''42e01f",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderFoldsH264ReceiveLevelHex(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"h264 max-recv-level": {
			registered: "video/H264; max-recv-level=000D",
			response:   "Video/H264; max-recv-level=000d",
		},
		"h264 extended max-recv-level": {
			registered: "video/H264; max-recv-level*=UTF-8''000D",
			response:   "Video/H264; max-recv-level*=utf-8''000d",
		},
		"h264 svc max-recv-level": {
			registered: "video/H264-SVC; max-recv-level=000D",
			response:   "Video/H264-SVC; max-recv-level=000d",
		},
		"h264 svc max-recv-base-level": {
			registered: "video/H264-SVC; max-recv-base-level=000D",
			response:   "Video/H264-SVC; max-recv-base-level=000d",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldInvalidOrCaseSensitiveH264Values(t *testing.T) {
	const base = "video/H264"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"invalid profile-level-id": {
			registered: "video/H264; profile-level-id=42E01",
			response:   "video/H264; profile-level-id=42e01",
		},
		"invalid two-digit max-recv-level": {
			registered: "video/H264; max-recv-level=0D",
			response:   "video/H264; max-recv-level=0d",
		},
		"sprop parameter sets": {
			registered: "video/H264; sprop-parameter-sets=QUJD",
			response:   "video/H264; sprop-parameter-sets=qujd",
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder", decoder)
			}
		})
	}
}

func TestRegisterDecoderRejectsMalformedPlainFallbackWithExtendedValue(t *testing.T) {
	for name, test := range map[string]struct {
		base       string
		registered string
		response   string
	}{
		"charset": {
			base:       "text/plain",
			registered: "text/plain; charset=UTF-8; charset*=UTF-8''UTF-8",
			response:   "text/plain; charset=BAD/VALUE; charset*=utf-8''utf-8",
		},
		"multipart protocol": {
			base:       "multipart/signed",
			registered: `multipart/signed; protocol="application/pgp-signature"; protocol*=UTF-8''application%2Fpgp-signature`,
			response:   `multipart/signed; protocol=application/pgp-signature; protocol*=utf-8''application%2fpgp-signature`,
		},
		"quoted fallback with newline": {
			base:       "text/plain",
			registered: "text/plain; charset=UTF-8; charset*=UTF-8''UTF-8",
			response:   "text/plain; charset=\"BAD\nVALUE\"; charset*=utf-8''utf-8",
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

			if registeredKey, responseKey := decoderContentTypeKey(test.registered), decoderContentTypeKey(test.response); registeredKey == responseKey {
				t.Fatalf("malformed plain fallback shares decoder key %q", registeredKey)
			}
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder == wantSpecific {
				t.Fatal("malformed plain fallback selected parameter-specific decoder")
			}
		})
	}
}

func TestRegisterDecoderRejectsInvalidUnquotedEqualDuplicateParameter(t *testing.T) {
	const base = "multipart/signed"
	registered := `multipart/signed; protocol="application/pgp-signature"; protocol="application/pgp-signature"`
	response := `multipart/signed; protocol="application/pgp-signature"; protocol=application/pgp-signature`

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("invalid unquoted duplicate shares decoder key %q", registeredKey)
	}
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder == wantSpecific {
		t.Fatal("invalid unquoted duplicate selected parameter-specific decoder")
	}
}

func TestRegisterDecoderUsesPlainFallbackWhenSingleExtendedValueCannotDecode(t *testing.T) {
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"invalid percent encoding": {
			registered: "text/plain; charset=UTF-8; charset*=utf-8''%zz",
			response:   "Text/Plain; charset=utf-8; charset*=UTF-8''%ZZ",
		},
		"missing extended metadata": {
			registered: "text/plain; charset=UTF-8; charset*=BROKEN",
			response:   "Text/Plain; charset=utf-8; charset*=broken",
		},
		"external body access type": {
			registered: "message/external-body; access-type=FTP; access-type*=utf-8''%zz; mode=IMAGE",
			response:   "Message/External-Body; access-type=ftp; access-type*=UTF-8''%ZZ; mode=image",
		},
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(test.registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != want {
				t.Fatalf("decoder = %T, want registered decoder via plain fallback", decoder)
			}
		})
	}
}

func TestRegisterDecoderDoesNotUsePlainFallbackForInvalidExtendedSyntax(t *testing.T) {
	registered := "text/plain; charset=UTF-8; charset*=UTF-8''UTF-8"
	for name, response := range map[string]string{
		"continuation gap": "Text/Plain; charset=utf-8; charset*0*=utf-8''UT; charset*2*=F-8",
		"unescaped slash":  "Text/Plain; charset=utf-8; charset*=UTF-8''BAD/VALUE",
	} {
		t.Run(name, func(t *testing.T) {
			want := &testDecoder{}
			RegisterDecoder(registered, func(io.ReadCloser) Decoder { return want })
			t.Cleanup(func() { delete(decoderTypes, decoderContentTypeKey(registered)) })
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder == want {
				t.Fatal("invalid extended syntax selected parameter-specific decoder through plain fallback")
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldUnknownTextPlainOptionValues(t *testing.T) {
	const base = "text/plain"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"format ordinary": {
			registered: "text/plain; format=x-openai-v1",
			response:   "Text/Plain; format=X-OPENAI-V1",
		},
		"delsp ordinary": {
			registered: "text/plain; delsp=x-openai-v1",
			response:   "Text/Plain; delsp=X-OPENAI-V1",
		},
		"format extended": {
			registered: "text/plain; format*=UTF-8''x-openai-v1",
			response:   "Text/Plain; format*=utf-8''X-OPENAI-V1",
		},
		"delsp extended": {
			registered: "text/plain; delsp*=UTF-8''x-openai-v1",
			response:   "Text/Plain; delsp*=utf-8''X-OPENAI-V1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for unknown text/plain option", decoder)
			}
		})
	}
}

func TestRegisterDecoderDoesNotFoldUnknownExternalBodyOptionValues(t *testing.T) {
	const base = "message/external-body"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"permission ordinary": {
			registered: "message/external-body; access-type=FTP; permission=x-openai-v1",
			response:   "Message/External-Body; access-type=ftp; permission=X-OPENAI-V1",
		},
		"permission extended": {
			registered: "message/external-body; access-type=FTP; permission*=UTF-8''x-openai-v1",
			response:   "Message/External-Body; access-type=ftp; permission*=utf-8''X-OPENAI-V1",
		},
		"ftp mode ordinary": {
			registered: "message/external-body; access-type=FTP; mode=x-openai-v1",
			response:   "Message/External-Body; access-type=ftp; mode=X-OPENAI-V1",
		},
		"tftp mode extended": {
			registered: "message/external-body; access-type=TFTP; mode*=UTF-8''x-openai-v1",
			response:   "Message/External-Body; access-type=tftp; mode*=utf-8''X-OPENAI-V1",
		},
		"tftp image invalid": {
			registered: "message/external-body; access-type=TFTP; mode=IMAGE",
			response:   "Message/External-Body; access-type=tftp; mode=image",
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for unknown external-body option", decoder)
			}
		})
	}
}

func TestDecoderContentTypeKeyPreservesPaddedBaseWithSemanticNormalization(t *testing.T) {
	padding := strings.Repeat(" ", 4096)
	extra := strings.Repeat("; x=v", 512)
	upper := "Message/External-Body" + padding + "; access-type=FTP; mode=IMAGE" + extra
	lower := "message/external-body" + padding + "; access-type=ftp; mode=image" + extra
	upperKey := decoderContentTypeKey(upper)
	lowerKey := decoderContentTypeKey(lower)
	if upperKey != lowerKey {
		t.Fatalf("equivalent padded content types produced different keys")
	}
	if !strings.HasPrefix(upperKey, "message/external-body"+padding+";") {
		t.Fatal("decoder key did not preserve padded raw media-type prefix")
	}
}

func TestDecoderContentTypeKeyAvoidsPerQuotedDuplicateAllocations(t *testing.T) {
	contentType := "text/plain" + strings.Repeat(`;charset="x"`, (256<<10)/12)
	allocs := testing.AllocsPerRun(3, func() {
		_ = decoderContentTypeKey(contentType)
	})
	if allocs > 32 {
		t.Fatalf("quoted duplicate parameter allocations = %.0f, want <= 32", allocs)
	}
}

func TestRegisterDecoderXOPTypePreservesRFC2231MetadataIdentity(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format*=UTF-8'en'FLOWED"`
	equivalent := `Application/Xop+Xml; type="Text/Plain; format*=utf-8'EN'flowed"`
	distinctLanguage := `Application/Xop+Xml; type="text/plain; format*=UTF-8'fr'FLOWED"`
	plain := `Application/Xop+Xml; type="text/plain; format=flowed"`

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent metadata":  {contentType: equivalent, want: wantSpecific},
		"distinct language":    {contentType: distinctLanguage, want: wantBare},
		"plain representation": {contentType: plain, want: wantBare},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.contentType}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want %T", decoder, test.want)
			}
		})
	}

	ordered := `application/xop+xml; type="text/plain; charset*=UTF-8''UTF-8; format*=UTF-8'en'FLOWED"`
	reordered := `Application/Xop+Xml; type="Text/Plain; format*=utf-8'EN'flowed; charset*=utf-8''utf-8"`
	if decoderContentTypeKey(ordered) != decoderContentTypeKey(reordered) {
		t.Fatal("reordered equivalent nested RFC 2231 parameters produced different keys")
	}
}

func TestRegisterDecoderNormalizesEqualDuplicateContinuationSections(t *testing.T) {
	const base = "text/plain"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"ordinary duplicate": {
			registered: `text/plain; charset*0*=UTF-8''UT; charset*0*=UTF-8''UT; charset*1*=F-8`,
			response:   `Text/Plain; charset*0*=utf-8''ut; charset*0*=utf-8''ut; charset*1*=f-8`,
		},
		"quoted duplicate": {
			registered: `text/plain; charset*0*=UTF-8''UT; charset*0*="UTF-8''UT"; charset*1*=F-8`,
			response:   `Text/Plain; charset*0*=utf-8''ut; charset*0*="utf-8''ut"; charset*1*=f-8`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantSpecific {
				t.Fatalf("decoder = %T, want duplicate-continuation decoder", decoder)
			}
		})
	}

	conflictingUpper := `text/plain; charset*0*=UTF-8''UT; charset*0*=UTF-8''XX; charset*1*=F-8`
	conflictingLower := `Text/Plain; charset*0*=utf-8''ut; charset*0*=utf-8''xx; charset*1*=f-8`
	if decoderContentTypeKey(conflictingUpper) == decoderContentTypeKey(conflictingLower) {
		t.Fatal("conflicting duplicate continuation sections collapsed")
	}
}

func TestRegisterDecoderAllowsEqualDuplicatePlainFallbackWithExtendedValue(t *testing.T) {
	const base = "text/plain"
	registered := `text/plain; charset=UTF-8; charset=UTF-8; charset*=UTF-8''UTF-8`
	response := `Text/Plain; charset=utf-8; charset=utf-8; charset*=utf-8''utf-8`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantSpecific {
		t.Fatalf("decoder = %T, want equal-fallback extended decoder", decoder)
	}
}

func TestRegisterDecoderXOPTypePreservesUnsupportedRFC2231ValueIdentity(t *testing.T) {
	for name, test := range map[string]struct {
		first  string
		second string
	}{
		"iso-8859-1": {
			first:  `application/xop+xml; type="text/plain; format*=ISO-8859-1'en'FLOWED"`,
			second: `application/xop+xml; type="text/plain; format*=ISO-8859-1'en'FIXED"`,
		},
		"blank charset": {
			first:  `application/xop+xml; type="text/plain; format*='en'FLOWED"`,
			second: `application/xop+xml; type="text/plain; format*='en'FIXED"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if decoderContentTypeKey(test.first) == decoderContentTypeKey(test.second) {
				t.Fatal("distinct nested RFC 2231 values collapsed")
			}
		})
	}
}

func TestRegisterDecoderXOPTypeSeparatesUnsupportedRFC2231FallbackFields(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; x*=\"''+\""`
	response := "application/xop+xml; type=\"text/plain; x*=\\\"" + string([]byte{0xb2}) + "''\\\"\""

	if decoderContentTypeKey(registered) == decoderContentTypeKey(response) {
		t.Fatal("distinct unsupported nested RFC 2231 charset/data fields collapsed")
	}

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
		t.Fatalf("decoder = %T, want bare decoder for distinct unsupported RFC 2231 identity", decoder)
	}
}

func TestRegisterDecoderXOPTypeNormalizesUnsupportedRFC2231MetadataCase(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format*=ISO-8859-1'en'FLOWED"`
	response := `Application/Xop+Xml; type="Text/Plain; format*=iso-8859-1'EN'FLOWED"`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantSpecific {
		t.Fatalf("decoder = %T, want unsupported-RFC2231 registration", decoder)
	}
}

func TestRegisterDecoderXOPTypeRejectsMalformedRFC2231ContinuationPayload(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format*0*=UTF-8''FLOWED; format*1*=\"\""`
	equivalent := `Application/Xop+Xml; type="Text/Plain; format*0*=utf-8''flowed; format*1*=\"\""`
	malformed := `Application/Xop+Xml; type="text/plain; format*0*=utf-8''flowed; format*1*=%ZZ"`
	malformedQuoted := `Application/Xop+Xml; type="text/plain; format*0*=utf-8''flowed; format*1*=\"%ZZ\""`

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent continuation":  {contentType: equivalent, want: wantSpecific},
		"malformed percent escape": {contentType: malformed, want: wantBare},
		"malformed quoted escape":  {contentType: malformedQuoted, want: wantBare},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.contentType}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want %T", decoder, test.want)
			}
		})
	}
}

func TestRegisterDecoderXOPTypeNormalizesUnsupportedRFC2231Continuation(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain; format*0*=ISO-8859-1'en'FLOW; format*1*=ED"`
	equivalent := `Application/Xop+Xml; type="Text/Plain; format*0*=iso-8859-1'EN'flow; format*1*=ed"`
	distinct := `Application/Xop+Xml; type="text/plain; format*0*=iso-8859-1'en'fix; format*1*=ed"`

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent continuation": {contentType: equivalent, want: wantSpecific},
		"distinct decoded value":  {contentType: distinct, want: wantBare},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.contentType}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want %T", decoder, test.want)
			}
		})
	}
}

func TestRegisterDecoderPreservesNonSpecialBackslashesInQuotedXOPType(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="text/plain"`
	backslashed := `Application/Xop+Xml; type="text/pla\in"`
	escapedSpecial := `Application/Xop+Xml; type="text\/plain"`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})
	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {backslashed}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for non-special quoted backslash", decoder)
	}
	if decoderContentTypeKey(registered) == decoderContentTypeKey(backslashed) {
		t.Fatal("non-special quoted backslash collapsed onto canonical media type")
	}
	escapedDecoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {escapedSpecial}}, Body: io.NopCloser(strings.NewReader(""))})
	if escapedDecoder != wantSpecific {
		t.Fatalf("decoder = %T, want specific decoder for escaped MIME tspecial", escapedDecoder)
	}
}

func TestRegisterDecoderEqualPlainFallbacksDoNotHideExtendedValue(t *testing.T) {
	const base = "text/plain"
	registered := `text/plain; charset=UTF-8; charset=UTF-8; charset*=UTF-8''US-ASCII`
	equivalent := `Text/Plain; charset=utf-8; charset=utf-8; charset*=utf-8''us-ascii`
	distinct := `Text/Plain; charset=utf-8; charset=utf-8; charset*=utf-8''UTF-8`
	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})
	for name, test := range map[string]struct {
		contentType string
		want        Decoder
	}{
		"equivalent extended value": {contentType: equivalent, want: wantSpecific},
		"distinct extended value":   {contentType: distinct, want: wantBare},
	} {
		t.Run(name, func(t *testing.T) {
			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.contentType}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != test.want {
				t.Fatalf("decoder = %T, want %T", decoder, test.want)
			}
		})
	}
}

func TestRegisterDecoderXOPTypeRejectsMixedSingleAndContinuedRFC2231Values(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="application/soap+xml; action*=UTF-8''V1; action*0*=UTF-8''LEFT; action*1*=ONE"`
	response := `Application/Xop+Xml; type="application/soap+xml; action*=utf-8''V1; action*0*=utf-8''RIGHT; action*1*=TWO"`

	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("mixed single/continued XOP RFC 2231 values share decoder key %q", registeredKey)
	}

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for mixed single/continued XOP RFC 2231 value", decoder)
	}
}

func TestRegisterDecoderXOPTypeSeparatesUnsupportedRFC2231LaterContinuationData(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="application/soap+xml; action*0*=X-UNKNOWN'en'V; action*1*=ONE"`
	response := `Application/Xop+Xml; type="application/soap+xml; action*0*=x-unknown'EN'V; action*1*=TWO"`

	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("unsupported XOP RFC 2231 continuations with distinct later data share decoder key %q", registeredKey)
	}

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for distinct unsupported RFC 2231 continuation data", decoder)
	}
}

func TestRegisterDecoderXOPTypeRejectsMixedRFC2231SectionEncodings(t *testing.T) {
	const base = "application/xop+xml"
	registered := `application/xop+xml; type="application/soap+xml; action*0=V1; action*0*=UTF-8''LEFT; action*1=Z"`
	response := `Application/Xop+Xml; type="application/soap+xml; action*0=V1; action*0*=utf-8''RIGHT; action*1=Z"`

	if registeredKey, responseKey := decoderContentTypeKey(registered), decoderContentTypeKey(response); registeredKey == responseKey {
		t.Fatalf("mixed encoded/unencoded XOP RFC 2231 sections share decoder key %q", registeredKey)
	}

	wantBare := &testDecoder{}
	wantSpecific := &testDecoder{}
	RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
	RegisterDecoder(registered, func(io.ReadCloser) Decoder { return wantSpecific })
	t.Cleanup(func() {
		delete(decoderTypes, decoderContentTypeKey(base))
		delete(decoderTypes, decoderContentTypeKey(registered))
	})

	decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {response}}, Body: io.NopCloser(strings.NewReader(""))})
	if decoder != wantBare {
		t.Fatalf("decoder = %T, want bare decoder for mixed encoded/unencoded RFC 2231 section", decoder)
	}
}

func TestRegisterDecoderXOPTypeRejectsInvalidRFC2231SectionNames(t *testing.T) {
	const base = "application/xop+xml"
	for name, test := range map[string]struct {
		registered string
		response   string
	}{
		"encoded leading zero": {
			registered: `application/xop+xml; type="application/soap+xml; action*00*=LEFT"`,
			response:   `Application/Xop+Xml; type="application/soap+xml; action*00*=RIGHT"`,
		},
		"unencoded leading zero": {
			registered: `application/xop+xml; type="application/soap+xml; action*00=LEFT"`,
			response:   `Application/Xop+Xml; type="application/soap+xml; action*00=RIGHT"`,
		},
		"nonnumeric section": {
			registered: `application/xop+xml; type="application/soap+xml; action*x*=LEFT"`,
			response:   `Application/Xop+Xml; type="application/soap+xml; action*x*=RIGHT"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if registeredKey, responseKey := decoderContentTypeKey(test.registered), decoderContentTypeKey(test.response); registeredKey == responseKey {
				t.Fatalf("invalid RFC 2231 section names share decoder key %q", registeredKey)
			}

			wantBare := &testDecoder{}
			wantSpecific := &testDecoder{}
			RegisterDecoder(base, func(io.ReadCloser) Decoder { return wantBare })
			RegisterDecoder(test.registered, func(io.ReadCloser) Decoder { return wantSpecific })
			t.Cleanup(func() {
				delete(decoderTypes, decoderContentTypeKey(base))
				delete(decoderTypes, decoderContentTypeKey(test.registered))
			})

			decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {test.response}}, Body: io.NopCloser(strings.NewReader(""))})
			if decoder != wantBare {
				t.Fatalf("decoder = %T, want bare decoder for invalid RFC 2231 section name", decoder)
			}
		})
	}
}
