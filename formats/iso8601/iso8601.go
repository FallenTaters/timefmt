package iso8601

import "github.com/FallenTaters/timefmt"

// Date follows the format YYYY-MM-DD and is a wrapper for timefmt.Date[iso8601.Format]
type Date struct {
	timefmt.Date[Format]
}

// Format is YYYY-MM-DD
type Format struct{}

// DateFormat implements timefmt.DateFormat
func (Format) DateFormat() string {
	return `2006-01-02`
}
