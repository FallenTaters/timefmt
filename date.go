package timefmt

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateFormat interface {
	DateFormat() string
}

type Date[F DateFormat] struct {
	t time.Time
}

func (d Date[F]) format() string {
	var f F
	return f.DateFormat()
}

func NewDate[T DateFormat](year int, month time.Month, day int) Date[T] {
	return Date[T]{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func DateFrom[T DateFormat](t time.Time) Date[T] {
	return Date[T]{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

func (d Date[F]) Time() time.Time {
	return d.t
}

func (d Date[F]) String() string {
	return d.t.Format(d.format())
}

func (d Date[F]) GoString() string {
	return fmt.Sprintf(`timefmt.NewDate[%T](%d, time.%s, %d)`, *new(F), d.t.Year(), d.t.Month().String(), d.t.Day())
}

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
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into timefmt.Date[%T]`, ErrScan, src, *new(F))
	}

	if err != nil {
		return fmt.Errorf(`%w: cannot unmarshal %v into timefmt.Date[%T]: %s`, ErrScan, src, *new(F), err.Error())
	}

	return nil
}

func (d Date[F]) Value() (driver.Value, error) {
	return d.Time(), nil
}

func (d *Date[F]) UnmarshalText(data []byte) error {
	tim, err := time.Parse(d.format(), string(data))
	if err == nil {
		*d = Date[F]{tim}
	}

	return err
}

func (d Date[F]) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
