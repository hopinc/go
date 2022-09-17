package types

import (
	"crypto/rand"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_b2s(t *testing.T) {
	// Allocate 10MB of memory.
	b := make([]byte, 10*1024*1024)

	// Fill it with random data.
	_, err := rand.Read(b)
	require.NoError(t, err)

	// Convert it to a string with b2s and the normal copying Go method.
	s1 := string(b)
	s2 := b2s(b)

	// Make sure both strings are equal.
	assert.Equal(t, s1, s2)
}

func Test_s2b(t *testing.T) {
	// Allocate 10MB of memory.
	b := make([]byte, 10*1024*1024)

	// Fill it with random data.
	_, err := rand.Read(b)
	require.NoError(t, err)

	// Convert it to a string which we will use as the origin.
	s := string(b)

	// Convert the string to a byte slice with s2b.
	b1 := s2b(s)

	// Make sure the byte slices are equal.
	assert.Equal(t, b, b1)
}

type pog struct {
	Hello string `json:"hello"`
}

func Test_stringifyValue(t *testing.T) {
	// Create a object.
	p := pog{Hello: "world"}

	// Stringify the value.
	s := stringifyValue(p)

	// GC the application to make sure the value is not garbage collected.
	runtime.GC()

	// Make sure the string is correct.
	assert.Equal(t, `pog({"hello":"world"})`, s)
}
