package api

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode"
)

// StripControl removes characters that can manipulate a terminal or spoof text:
// control characters (ESC, CR, C0/C1, DEL), format characters such as bidi
// overrides (U+202E) and zero-width spaces, and line/paragraph separators.
// Newlines, tabs and the zero-width joiner used by emoji are kept.
func StripControl(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\t' || r == '‍':
			return r
		case unicode.IsControl(r), unicode.Is(unicode.Cf, r), r == ' ', r == ' ':
			return -1
		}
		return r
	}, s)
}

// decode unmarshals server JSON and strips control characters from every
// string field, so server-supplied text is safe wherever it is printed.
func decode(data []byte, out any) error {
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	sanitize(reflect.ValueOf(out))
	return nil
}

func sanitize(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			sanitize(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				sanitize(v.Field(i))
			}
		}
	case reflect.Slice, reflect.Array:
		// Skip byte slices such as json.RawMessage.
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return
		}
		for i := 0; i < v.Len(); i++ {
			sanitize(v.Index(i))
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(StripControl(v.String()))
		}
	}
}
