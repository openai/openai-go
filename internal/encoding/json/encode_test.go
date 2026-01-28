package json

import (
	"bytes"
	"testing"
)

// Inner implements MarshalJSON to trigger the optimized code path
type benchInner struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (b benchInner) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{b.Name, b.Value})
}

// Nested structure with multiple MarshalJSON calls
type benchNested struct {
	Inner benchInner `json:"inner"`
	Items []int      `json:"items"`
}

func (b benchNested) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Inner benchInner `json:"inner"`
		Items []int      `json:"items"`
	}{b.Inner, b.Items})
}

// Deeply nested to amplify the effect
type benchDeep struct {
	Level1 benchNested `json:"level1"`
	Level2 benchNested `json:"level2"`
	Data   string      `json:"data"`
}

func (b benchDeep) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Level1 benchNested `json:"level1"`
		Level2 benchNested `json:"level2"`
		Data   string      `json:"data"`
	}{b.Level1, b.Level2, b.Data})
}

func BenchmarkMarshalNestedMarshalJSON(b *testing.B) {
	data := benchDeep{
		Level1: benchNested{
			Inner: benchInner{Name: "test1", Value: 100},
			Items: []int{1, 2, 3, 4, 5},
		},
		Level2: benchNested{
			Inner: benchInner{Name: "test2", Value: 200},
			Items: []int{6, 7, 8, 9, 10},
		},
		Data: "some test data here",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Slice of nested structs - common real-world pattern
func BenchmarkMarshalSliceOfNestedMarshalJSON(b *testing.B) {
	data := make([]benchDeep, 50)
	for i := range data {
		data[i] = benchDeep{
			Level1: benchNested{
				Inner: benchInner{Name: "test1", Value: i},
				Items: []int{1, 2, 3, 4, 5},
			},
			Level2: benchNested{
				Inner: benchInner{Name: "test2", Value: i * 2},
				Items: []int{6, 7, 8, 9, 10},
			},
			Data: "some test data here that is a bit longer to simulate real payloads",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test that HTML escaping is disabled for nested MarshalJSON calls
type htmlTestInner struct {
	Content string `json:"content"`
}

func (h htmlTestInner) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Content string `json:"content"`
	}{h.Content})
}

type htmlTestOuter struct {
	Inner htmlTestInner `json:"inner"`
}

func (h htmlTestOuter) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Inner htmlTestInner `json:"inner"`
	}{h.Inner})
}

func TestNoHTMLEscape(t *testing.T) {
	type testCase struct {
		name string
		// encodeFunc in each test case captures the various ways in which this package
		// can produce JSON output.
		encodeFunc        func(data any) ([]byte, error)
		expectEscapedHTML bool
	}
	tests := []testCase{
		{
			name:              "Marshal",
			encodeFunc:        Marshal,
			expectEscapedHTML: false,
		},
		{
			name: "Encoder",
			encodeFunc: func(data any) ([]byte, error) {
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				if err := enc.Encode(data); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			expectEscapedHTML: false,
		},
		{
			name: "Encoder with SetEscapeHTML(true)",
			encodeFunc: func(data any) ([]byte, error) {
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				enc.SetEscapeHTML(true)
				if err := enc.Encode(data); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			// Even though the encoder is configured to escape HTML,
			// types that implement MarshalJSON do not receive the
			// encoder's options, so the HTML escaping is not propagated
			// to them. Rather, their implementations of MarshalJSON
			// determine whether to escape HTML.
			expectEscapedHTML: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// If HTML escaping is disabled, this byte slice should be identically present in the test output.
			htmlContentBytes := []byte("<div>&amp;</div><script>alert('xss')</script>")

			data := htmlTestOuter{
				Inner: htmlTestInner{
					Content: string(htmlContentBytes),
				},
			}

			result, err := test.encodeFunc(data)
			if err != nil {
				t.Fatalf("%s failed: %v", test.name, err)
			}

			// remove newlines from the result for comparison
			result = bytes.ReplaceAll(result, []byte("\n"), []byte(""))

			if !test.expectEscapedHTML && !bytes.Contains(result, htmlContentBytes) {
				t.Errorf("%s expected unescaped HTML in output, got: %s", test.name, result)
			}

			if test.expectEscapedHTML {
				// appendCompact will escape HTML characters if instructed.
				// we use it to compute what a payload with escaped HTML would look like,
				// so that we can assert whether the result has been correctly escaped if required.
				compactedResult := []byte{}
				compactedResult, err = appendCompact(compactedResult, result, true)
				if err != nil {
					t.Fatalf("%s failed: %v", test.name, err)
				}

				if !bytes.Equal(compactedResult, result) {
					t.Errorf("%s expected escaped HTML in output, got: %s", test.name, result)
				}
			}
		})
	}
}
