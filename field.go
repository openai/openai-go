package openai

import (
	"github.com/openai/openai-go/packages/param"
	"io"
	"time"
)

func String(s string) param.String {
	fld := param.NeverOmitted[param.String]()
	fld.V = s
	return fld
}

func Int(i int64) param.Int {
	fld := param.NeverOmitted[param.Int]()
	fld.V = i
	return fld
}

func Bool(b bool) param.Bool {
	fld := param.NeverOmitted[param.Bool]()
	fld.V = b
	return fld
}

func Float(f float64) param.Float {
	fld := param.NeverOmitted[param.Float]()
	fld.V = f
	return fld
}

func Datetime(t time.Time) param.Datetime {
	fld := param.NeverOmitted[param.Datetime]()
	fld.V = t
	return fld
}

func Date(t time.Time) param.Date {
	fld := param.NeverOmitted[param.Date]()
	fld.V = t
	return fld
}

func Ptr[T any](v T) *T { return &v }

func File(rdr io.Reader, filename string, contentType string) file {
	return file{rdr, filename, contentType}
}

type file struct {
	io.Reader
	name        string
	contentType string
}

func (f file) Filename() string {
	if f.name != "" {
		return f.name
	} else if named, ok := f.Reader.(interface{ Name() string }); ok {
		return named.Name()
	}
	return ""
}

func (f file) ContentType() string {
	return f.contentType
}
