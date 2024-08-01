package apijson

import (
	"reflect"

	"github.com/tidwall/gjson"
)

type UnionVariant struct {
	TypeFilter         gjson.Type
	DiscriminatorValue interface{}
	Type               reflect.Type
}

var unionRegistry = map[reflect.Type]unionEntry{}

type unionEntry struct {
	discriminatorKey string
	variants         []UnionVariant
}

func RegisterUnion(typ reflect.Type, discriminator string, variants ...UnionVariant) {
	unionRegistry[typ] = unionEntry{
		discriminatorKey: discriminator,
		variants:         variants,
	}
}
