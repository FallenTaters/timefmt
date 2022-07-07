package timefmt

import "time"

// YearMonthDayTime is YYYY-MM-DD HH:MM:SS
type YearMonthDayTime struct{}

// TimeFormat implements timefmt.TimeFormat
func (YearMonthDayTime) TimeFormat() string {
	return `2006-01-02 15:04:05`
}

// MonthDayYear is MM/DD/YYYY, as is common in North America.
type MonthDayYear struct{}

// DateFormat implements timefmt.DateFormat
func (MonthDayYear) DateFormat() string {
	return `01/02/2006`
}

// DayMonthYear is DD/MM/YYYY, as is common in Europe.
type DayMonthYear struct{}

// DateFormat implements timefmt.DateFormat
func (DayMonthYear) DateFormat() string {
	return `02/01/2006`
}

// ISO8601Date is YYYY-MM-DD
type ISO8601Date struct{}

// DateFormat implements timefmt.DateFormat
func (ISO8601Date) DateFormat() string {
	return `2006-01-02`
}

// RFC3339 is a timefmt.TimeFormat for time.RFC3339.
// It is sometimes preferable to simply use time.Time,
// which uses time.RFC3339 and time.RFC3339Nano for json marshal/unmarshal.
type RFC3339 struct{}

// TimeFormat implements timefmt.TimeFormat
func (RFC3339) TimeFormat() string {
	return time.RFC3339
}
