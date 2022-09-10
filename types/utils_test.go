package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampFromTime(t *testing.T) {
	x := time.Now()
	y, err := TimestampFromTime(x).Time()
	assert.NoError(t, err)
	assert.Equal(t, x.Unix(), y.Unix())
}

func TestStringPointerify(t *testing.T) {
	x := "x"
	assert.Equal(t, &x, StringPointerify(x))
	x = ""
	assert.Nil(t, StringPointerify(x))
}

func expectsBytes(t *testing.T, s Size, bytes int) {
	t.Helper()
	x, err := s.Bytes()
	assert.NoError(t, err)
	assert.Equal(t, bytes, x)
}

func TestKilobytes(t *testing.T) {
	expectsBytes(t, Kilobytes(1), 1024)
}

func TestMegabytes(t *testing.T) {
	expectsBytes(t, Megabytes(1), 1024*1024)
}

func TestGigabytes(t *testing.T) {
	expectsBytes(t, Gigabytes(1), 1024*1024*1024)
}

func TestBytes(t *testing.T) {
	expectsBytes(t, Bytes(2), 2)
}
