package timefmt

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateFormat needs to be implemented for a custom format F for Date[F].
type DateFormat interface {
	DateFormat() string
}

// Date contains a UTC time.Time stripped of time (hour, minute etc.) and timezone information.
// The specified format is used for unmarshalling and marshalling.
type Date[F DateFormat] struct {
	t time.Time
}

// NewDate makes a new Date[F] from the specified year, month and day.
func NewDate[F DateFormat](year int, month time.Month, day int) Date[F] {
	return Date[F]{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateFrom strips the time and zone component from the time and returns it as a Date[F].
func DateFrom[F DateFormat](t time.Time) Date[F] {
	return Date[F]{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// ParseDate parses the input in format F and strips the time component.
func ParseDate[F DateFormat](input string) (Date[F], error) {
	var f F
	tim, err := time.Parse(f.DateFormat(), input)
	return DateFrom[F](tim), err
}

func (d Date[F]) format() string {
	var f F
	return f.DateFormat()
}

// Time returns the date as a UTC time.Time.
func (d Date[F]) Time() time.Time {
	return d.t
}

// String implements fmt.Stringer and formats the date according to F.
func (d Date[F]) String() string {
	return d.t.Format(d.format())
}

// GoString implements fmt.GoStringer.
func (d Date[F]) GoString() string {
	var f F
	return fmt.Sprintf(`timefmt.NewDate[%T](%d, time.%s, %d)`, f, d.t.Year(), d.t.Month().String(), d.t.Day())
}

// Scan implements sql.Scanner and supports both string/bytes (in format F) and time.Time.
func (d *Date[F]) Scan(src any) error {
	var err error

	switch v := src.(type) {
	case []byte:
		err = d.UnmarshalText(v)
	case string:
		err = d.UnmarshalText([]byte(v))
	case time.Time:
		*d = DateFrom[F](v)
	default:
		var f F
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into timefmt.Date[%T]`, ErrScan, src, f)
	}

	if err != nil {
		var f F
		return fmt.Errorf(`%w: cannot unmarshal %v into timefmt.Date[%T]: %s`, ErrScan, src, f, err.Error())
	}

	return nil
}

// Value implements driver.Valuer and returns a time.Time.
func (d Date[F]) Value() (driver.Value, error) {
	return d.Time(), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and parses the date and strips the time components.
func (d *Date[F]) UnmarshalText(data []byte) error {
	tim, err := time.Parse(d.format(), string(data))
	if err == nil {
		*d = DateFrom[F](tim)
	}

	return err //nolint:wrapcheck
}

// MarshalText implements encoding.TextMarshaler and uses F to format the date.
func (d Date[F]) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
