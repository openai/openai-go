// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package apiform

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/packages/param"
)

func TestEncodedMultipartParts(t *testing.T) {
	const offer = "v=0\r\ns=Unicode π\r\n"
	encodings := map[string]PartEncoding{
		"offer":    {ContentType: "application/sdp"},
		"settings": {ContentType: "application/json", JSON: true},
	}
	for _, settings := range []any{nil, map[string]any{"future": []any{false, nil, "π"}}} {
		value := struct {
			Offer    string         `json:"offer" api:"required"`
			Settings map[string]any `json:"settings,omitzero"`
			Ordinary string         `json:"ordinary,omitzero"`
		}{Offer: offer, Settings: map[string]any{"old": true}, Ordinary: "value"}
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		if err := MarshalEncodedRoot(value, writer, map[string]any{"settings": settings}, encodings); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		reader := multipart.NewReader(&buffer, writer.Boundary())
		count := 0
		for {
			part, err := reader.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}
			if part.FileName() != "" {
				t.Fatal("encoded fields must not be files")
			}
			switch part.FormName() {
			case "offer":
				if string(data) != offer || part.Header.Get("Content-Type") != "application/sdp" {
					t.Fatal("incorrect SDP part")
				}
			case "settings":
				var decoded any
				if err := json.Unmarshal(data, &decoded); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(decoded, settings) || part.Header.Get("Content-Type") != "application/json" {
					t.Fatal("incorrect JSON part")
				}
			case "ordinary":
				if string(data) != "value" {
					t.Fatal("ordinary form field changed")
				}
			default:
				t.Fatalf("unexpected form field %q", part.FormName())
			}
			count++
		}
		if count != 3 {
			t.Fatalf("expected three parts, got %d", count)
		}
	}
}

func TestEncodedMultipartOmitOverrides(t *testing.T) {
	value := struct {
		Offer    string         `json:"offer" api:"required"`
		Settings map[string]any `json:"settings,omitzero"`
		Ordinary string         `json:"ordinary,omitzero"`
	}{Offer: "v=0\r\n", Settings: map[string]any{"model": "test-model"}, Ordinary: "value"}
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	err := MarshalEncodedRoot(value, writer,
		map[string]any{"offer": param.Omit, "settings": param.Omit, "ordinary": param.Omit},
		map[string]PartEncoding{"offer": {ContentType: "application/sdp"}, "settings": {ContentType: "application/json", JSON: true}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if closeErr := writer.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	_, err = multipart.NewReader(&buffer, writer.Boundary()).NextPart()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("omitted fields were sent: %v", err)
	}
}
