package timefmt

/*
Further formats can be found in date/format/formats/* packages.
*/

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
