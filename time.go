package timefmt

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// ErrScan is returned when timefmt.Time.Scan() or timefmt.Date.Scan() fails
var ErrScan = errors.New(`scan failed`)

type TimeFormat interface {
	TimeFormat() string
}

type Time[F TimeFormat] struct {
	t time.Time
}

func (t Time[F]) format() string {
	var f F
	return f.TimeFormat()
}

func TimeFrom[T TimeFormat](t time.Time) Time[T] {
	return Time[T]{t}
}

func (t Time[F]) Time() time.Time {
	return t.t
}

func (t Time[F]) String() string {
	return t.t.Format(t.format())
}

func (t Time[F]) GoString() string {
	return fmt.Sprintf(`timefmt.TimeFrom[%T](%#v)`, *new(F), t.t)
}

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
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into timefmt.Time[%T]`, ErrScan, src, *new(F))
	}

	if err != nil {
		return fmt.Errorf(`%w: cannot unmarshal %v into timefmt.Time[%T]: %s`, ErrScan, src, *new(F), err.Error())
	}

	return nil
}

func (t Time[F]) Value() (driver.Value, error) {
	return t.t, nil
}

func (t *Time[F]) UnmarshalText(data []byte) error {
	tim, err := time.Parse(t.format(), string(data))
	if err == nil {
		*t = Time[F]{tim}
	}

	return err
}

func (t Time[F]) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
