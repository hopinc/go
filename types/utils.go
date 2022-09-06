package types

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/relvacode/iso8601"
)

// Timestamp is used to define a timestamp within the Hop API. Hop uses the ISO 8601 standard.
type Timestamp string

// TimestampFromTime takes a time.Time and turns it into a Hop compatible timestamp.
func TimestampFromTime(t time.Time) Timestamp {
	// ISO 8601 is not always valid RFC3339, but RFC3339 is always valid ISO 8601.
	return (Timestamp)(t.Format(time.RFC3339))
}

// Time turns the timestamp into a time.Time object.
func (t Timestamp) Time() (time.Time, error) {
	return iso8601.ParseString(string(t))
}

// StringPointerify is used to turn the normal Go string behaviour into a string pointer.
// If the string is blank, a nil pointer will be returned. If not, a pointer to the string will be instead.
func StringPointerify(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Size is used to define a memory/storage size and allow for easy parsing of it.
type Size string

type strIntPair struct {
	str string
	int int
}

var multipliers = []strIntPair{
	{"kb", 1024}, {"mb", 1024 * 1024},
	{"gb", 1024 * 1024 * 1024}, {"b", 1},
}

// Bytes returns the size in bytes.
func (s Size) Bytes() (int, error) {
	for _, p := range multipliers {
		if strings.HasSuffix((string)(s), p.str) {
			// Parse the rest of the content and multiply it by the multiplier.
			i, err := strconv.Atoi((string)(s)[:len(s)-len(p.str)])
			if err != nil {
				return 0, err
			}
			return i * p.int, nil
		}
	}
	return 0, errors.New("invalid size")
}

// Kilobytes returns a size type for the number of kilobytes specified.
func Kilobytes(kb int) Size {
	return Size(strconv.Itoa(kb) + "kb")
}

// Megabytes returns a size type for the number of megabytes specified.
func Megabytes(mb int) Size {
	return Size(strconv.Itoa(mb) + "mb")
}

// Gigabytes returns a size type for the number of gigabytes specified.
func Gigabytes(gb int) Size {
	return Size(strconv.Itoa(gb) + "gb")
}

// Bytes returns a size type for the number of bytes specified.
func Bytes(b int) Size {
	return Size(strconv.Itoa(b) + "b")
}
