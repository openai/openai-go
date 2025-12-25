package respjson_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3/internal/apijson"
	rj "github.com/Nordlys-Labs/openai-go/v3/packages/respjson"
)

type UnionOfStringIntOrObject struct {
	OfString string    `json:",inline"`
	OfInt    int       `json:",inline"`
	Type     string    `json:"type"`
	Function SubFields `json:"function"`
	JSON     struct {
		OfString rj.Field
		OfInt    rj.Field
		Type     rj.Field
		Function rj.Field
		raw      string
	} `json:"-"`
}

func (u UnionOfStringIntOrObject) RawJSON() string { return u.JSON.raw }
func (r *UnionOfStringIntOrObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SubFields struct {
	OfBool bool   `json:",inline"`
	Name   string `json:"name,required"`
	JSON   struct {
		OfBool      rj.Field
		Name        rj.Field
		ExtraFields map[string]rj.Field
		raw         string
	} `json:"-"`
}

func (r SubFields) RawJSON() string { return r.JSON.raw }
func (r *SubFields) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func TestUnmarshalUnionString(t *testing.T) {
	rawJSON := `"123"`
	testUnmarshalUnion(t, rawJSON, func(res UnionOfStringIntOrObject) map[string]error {
		return map[string]error{
			"rawJSON": checkEqual(res.RawJSON(), rawJSON),

			"string":          checkEqual(res.OfString, "123"),
			"int":             checkEqual(res.OfInt, 0),
			"$.type":          checkEqual(res.Type, ""),
			"$.function.name": checkEqual(res.Function.Name, ""),

			"string.meta":          checkMeta(res.JSON.OfString, rawJSON, shouldBePresent),
			"int.meta":             checkMeta(res.JSON.OfInt, "", shouldBeNullish),
			"$.type.meta":          checkMeta(res.JSON.Type, "", shouldBeNullish),
			"$.function.meta":      checkMeta(res.Function.JSON.Name, "", shouldBeNullish),
			"$.function.name.meta": checkMeta(res.Function.JSON.Name, "", shouldBeNullish),
		}
	})
}

func TestUnmarshalUnionInt(t *testing.T) {
	rawJSON := `123`
	testUnmarshalUnion(t, rawJSON, func(res UnionOfStringIntOrObject) map[string]error {
		return map[string]error{
			"rawJSON": checkEqual(res.RawJSON(), rawJSON),

			"string":          checkEqual(res.OfString, ""),
			"int":             checkEqual(res.OfInt, 123),
			"$.type":          checkEqual(res.Type, ""),
			"$.function.name": checkEqual(res.Function.Name, ""),
			"$.function.bool": checkEqual(res.Function.OfBool, false),

			"string.meta":          checkMeta(res.JSON.OfString, "", shouldBeNullish),
			"int.meta":             checkMeta(res.JSON.OfInt, rawJSON, shouldBePresent),
			"$.type.meta":          checkMeta(res.JSON.Type, "", shouldBeNullish),
			"$.function.meta":      checkMeta(res.Function.JSON.Name, "", shouldBeNullish),
			"$.function.name.meta": checkMeta(res.Function.JSON.Name, "", shouldBeNullish),
		}
	})

	testUnmarshalUnion(t, `0`, func(res UnionOfStringIntOrObject) map[string]error {
		return map[string]error{
			"rawJSON": checkEqual(res.RawJSON(), "0"),
			"string":  checkEqual(res.OfString, ""),

			"int":         checkEqual(res.OfInt, 0),
			"int.meta":    checkMeta(res.JSON.OfInt, "0", shouldBePresent),
			"string.meta": checkMeta(res.JSON.OfString, "", shouldBeNullish),
		}
	})
}

func TestUnmarshalUnionObject(t *testing.T) {
	rawJSON := `{"type":"auto","function":{"name":"test_fn"}}`
	testUnmarshalUnion(t, rawJSON, func(res UnionOfStringIntOrObject) map[string]error {
		return map[string]error{
			"rawJSON": checkEqual(res.RawJSON(), rawJSON),

			"string":          checkEqual(res.OfString, ""),
			"int":             checkEqual(res.OfInt, 0),
			"$.type":          checkEqual(res.Type, "auto"),
			"$.function.name": checkEqual(res.Function.Name, "test_fn"),
			"$.function.bool": checkEqual(res.Function.OfBool, false),

			"string.meta":          checkMeta(res.JSON.OfString, "", shouldBeNullish),
			"int.meta":             checkMeta(res.JSON.OfInt, "", shouldBeNullish),
			"$.type.meta":          checkMeta(res.JSON.Type, `"auto"`, shouldBePresent),
			"$.function.meta":      checkMeta(res.JSON.Function, `{"name":"test_fn"}`, shouldBePresent),
			"$.function.name.meta": checkMeta(res.Function.JSON.Name, `"test_fn"`, shouldBePresent),
			"$.function.bool.meta": checkMeta(res.Function.JSON.OfBool, "", shouldBeNullish),
		}
	})
}

func TestUnmarshalUnionObjectWithInlineSubUnion(t *testing.T) {
	rawJSON := `{"type":"auto","function":true}`
	testUnmarshalUnion(t, rawJSON, func(res UnionOfStringIntOrObject) map[string]error {
		return map[string]error{
			"rawJSON": checkEqual(res.RawJSON(), rawJSON),

			"string":     checkEqual(res.OfString, ""),
			"int":        checkEqual(res.OfInt, 0),
			"$.type":     checkEqual(res.Type, "auto"),
			"$.function": checkEqual(res.Function.OfBool, true),

			"string.meta":          checkMeta(res.JSON.OfString, "", shouldBeNullish),
			"int.meta":             checkMeta(res.JSON.OfInt, "", shouldBeNullish),
			"$.type.meta":          checkMeta(res.JSON.Type, `"auto"`, shouldBePresent),
			"$.function.meta":      checkMeta(res.JSON.Function, `true`, shouldBePresent),
			"$.function.name.meta": checkMeta(res.Function.JSON.Name, "", shouldBeNullish),
			"$.function.bool.meta": checkMeta(res.Function.JSON.OfBool, `true`, shouldBePresent),
		}
	})
}

/*********/
/* UTILS */
/*********/

func testUnmarshalUnion[T any](t *testing.T, raw string, check testChecks[T]) {
	var res T
	err := json.Unmarshal([]byte(raw), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal %v", err.Error())
	}

	for label, fail := range check(res) {
		if fail != nil {
			t.Errorf("failed check %v: %v", label, fail.Error())
		}
	}
}

func checkEqual[T any](got, expected T) error {
	if reflect.DeepEqual(got, expected) {
		return nil
	}
	return fmt.Errorf("not equal: got %v, expected %v", got, expected)
}

type metaStatus int

const (
	shouldBePresent metaStatus = iota
	shouldBeNullish
	shouldBeInvalid
)

type testChecks[T any] func(T) map[string]error

func checkMeta(got rj.Field, raw string, stat metaStatus) error {
	switch stat {
	case shouldBePresent:
		if !got.Valid() {
			return fmt.Errorf("expected field to be present, but got nullish")
		}
		if got.Raw() != raw {
			return fmt.Errorf("expected field to be present with raw value %v, but got %v", raw, got.Raw())
		}
	case shouldBeNullish:
		if got.Valid() {
			return fmt.Errorf("expected field to be nullish, but got %v", got.Raw())
		}
		if got.Raw() != rj.Omitted && got.Raw() != rj.Null {
			return fmt.Errorf("expected field to be nullish, but got %v", got.Raw())
		}
	case shouldBeInvalid:
		if !got.Valid() || got.Raw() == "" {
			return fmt.Errorf("expected field to be invalid, but got valid value %v", got.Raw())
		}
		if got.Raw() != raw {
			return fmt.Errorf("expected field to be invalid, but got valid value %v", got.Raw())
		}
	default:
		return fmt.Errorf("unknown metaStatus: %v", stat)
	}
	return nil
}
