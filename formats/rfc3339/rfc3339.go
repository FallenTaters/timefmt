package rfc3339

import (
	"time"

	"github.com/FallenTaters/timefmt"
)

// Time follows time.RFC3339 format and is a shorthand variant for timefmt.Time[rfc3339.Format]
type Time struct {
	timefmt.Time[Format]
}

// Format is a timefmt.TimeFormat for time.RFC3339.
// It is sometimes preferable to simply use time.Time,
// which uses time.RFC3339 and time.RFC3339Nano for json marshal/unmarshal.
type Format struct{}

// TimeFormat implements timefmt.TimeFormat
func (Format) TimeFormat() string {
	return time.RFC3339
}
