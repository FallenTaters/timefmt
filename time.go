/*
Package timefmt contains the types timefmt.Date and timefmt.Time,
which wrap time.Time in order to do marshal and unmarshal in a certain format using type parameters.

The types also implement sql.Scanner and driver.Valuer, but Value() always returns a time.Time, never a []byte/string.

Custom formats can be made by defining a struct that implements either timefmt.TimeFormat or timefmt.DateFormat
for timefmt.Time and timefmt.Date, respectively.
*/
package timefmt

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// ErrScan is returned when timefmt.Time.Scan() or timefmt.Date.Scan() fails.
var ErrScan = errors.New(`scan failed`)

// TimeFormat needs to be implemented for a custom format F for Time[F].
type TimeFormat interface {
	TimeFormat() string
}

// Time wraps time.Time. The specified format is used for unmarshalling and marshalling.
type Time[F TimeFormat] struct {
	t time.Time
}

// NewTime works like time.Date() but returns a Time[F].
func NewTime[F TimeFormat](year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time[F] {
	return Time[F]{time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// TimeFrom wraps the time in a Time[F].
func TimeFrom[F TimeFormat](t time.Time) Time[F] {
	return Time[F]{t}
}

// ParseTime parses the input in format F.
func ParseTime[F TimeFormat](input string) (Time[F], error) {
	var f F
	tim, err := time.Parse(f.TimeFormat(), input)
	return Time[F]{tim}, err //nolint:wrapcheck
}

func (t Time[F]) format() string {
	var f F
	return f.TimeFormat()
}

// Time returns the time.Time wrapped in Time[F].
func (t Time[F]) Time() time.Time {
	return t.t
}

// String implements fmt.Stringer and formats the time.Time according to F.
func (t Time[F]) String() string {
	return t.t.Format(t.format())
}

// GoString implements fmt.GoStringer.
func (t Time[F]) GoString() string {
	var f F
	return fmt.Sprintf(`timefmt.TimeFrom[%T](%#v)`, f, t.t)
}

// Scan imoplements sql.Scanner and supports both string/bytes (in format F) and time.Time.
func (t *Time[F]) Scan(src any) error {
	var err error

	switch v := src.(type) {
	case []byte:
		err = t.UnmarshalText(v)
	case string:
		err = t.UnmarshalText([]byte(v))
	case time.Time:
		*t = TimeFrom[F](v)
	default:
		var f F
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into timefmt.Time[%T]`, ErrScan, src, f)
	}

	if err != nil {
		var f F
		return fmt.Errorf(`%w: cannot unmarshal %v into timefmt.Time[%T]: %s`, ErrScan, src, f, err.Error())
	}

	return nil
}

// Value implements driver.Valuer and returns a time.Time.
func (t Time[F]) Value() (driver.Value, error) {
	return t.t, nil
}

// UnmarshalText implements encoding.TextUnmarshaler and parses the time using F.
func (t *Time[F]) UnmarshalText(data []byte) error {
	tim, err := time.Parse(t.format(), string(data))
	if err == nil {
		*t = Time[F]{tim}
	}

	return err //nolint:wrapcheck
}

// MarshalText implements encoding.TextMarshaler and uses F to format the time.
func (t Time[F]) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
