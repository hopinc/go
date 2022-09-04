package types

import (
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
