package types

import (
	"bytes"
	"encoding/json"
	"reflect"
	"unsafe"
)

// --- CREDIT: https://github.com/valyala/fasthttp#tricks-with-byte-buffers (MIT licensed) ---

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// --- END CREDIT ---

// This function turns a value into a stringified version of itself. It uses JSON to do this.
func stringifyValue(x any) string {
	// Get the type name with reflect.
	typeName := reflect.TypeOf(x).Name()

	// Create a bytes buffer.
	var b bytes.Buffer
	_, _ = b.Write([]byte(typeName))
	_ = b.WriteByte('(')

	// Use a JSON writer to write the value to the buffer.
	err := json.NewEncoder(&b).Encode(x)
	if err != nil {
		return "<json marshal error, please report this to the hop-go repository>"
	}

	// Return a stringified version of the buffer with the last value as a bracket.
	ret := b.Bytes()
	ret[len(ret)-1] = ')'
	return b2s(ret)
}
