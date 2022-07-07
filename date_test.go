package timefmt_test

import (
	"testing"
	"time"

	"github.com/FallenTaters/timefmt"
)

var times = []struct {
	name string
	t    time.Time
}{
	{`zero time`, time.Time{}},
	{`negative time (out of range values)`, time.Date(-1, -1, -1, 0, 0, 0, 0, time.UTC)},
	{`far future time`, time.Date(100_000, time.December, 25, 0, 0, 0, 0, time.UTC)},
	{`time included`, time.Date(2020, time.February, 30, 12, 13, 14, 151617, time.Local)},
	{`now`, time.Now()},
}

func TestNewDate(t *testing.T) {
	t.Run(`normal date`, func(t *testing.T) {
		if d := timefmt.NewDate[timefmt.ISO8601Date](2020, time.December, 13); d.String() != `2020-12-13` {
			t.Error(d)
		}
	})

	t.Run(`values out of range -> should be normalized`, func(t *testing.T) {
		if d := timefmt.NewDate[timefmt.ISO8601Date](-1, -1, -1); d.String() != `-0002-10-30` {
			t.Error(d)
		}
	})
}

func TestDateFrom(t *testing.T) {
	for _, c := range times {
		t.Run(c.name, func(t *testing.T) {
			if d := timefmt.DateFrom[timefmt.ISO8601Date](c.t); d.String() != c.t.Format((timefmt.ISO8601Date{}).DateFormat()) {
				t.Errorf(`date %s should be %q but is %q`, c.name, d.String(), c.t.Format((timefmt.ISO8601Date{}).DateFormat()))
			}
		})
	}
}

func TestDateGoString(t *testing.T) {
	expected := `timefmt.NewDate[timefmt.ISO8601Date](1, time.January, 1)`
	if actual := timefmt.DateFrom[timefmt.ISO8601Date](times[0].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.NewDate[timefmt.ISO8601Date](-2, time.October, 30)`
	if actual := timefmt.DateFrom[timefmt.ISO8601Date](times[1].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.NewDate[timefmt.ISO8601Date](100000, time.December, 25)`
	if actual := timefmt.DateFrom[timefmt.ISO8601Date](times[2].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.NewDate[timefmt.ISO8601Date](2020, time.March, 1)`
	if actual := timefmt.DateFrom[timefmt.ISO8601Date](times[3].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}
}

func TestDateComparable(t *testing.T) {
	actual, expected := timefmt.DateFrom[timefmt.ISO8601Date](times[3].t), timefmt.NewDate[timefmt.ISO8601Date](2020, 3, 1)
	if expected != actual {
		t.Error(actual, expected)
	}
}

func TestDateScan(t *testing.T) {
	cases := []struct {
		name      string
		input     any
		expectErr string
		expected  timefmt.Date[timefmt.ISO8601Date]
	}{
		{
			name:     `string date`,
			input:    `2006-10-12`,
			expected: timefmt.NewDate[timefmt.ISO8601Date](2006, 10, 12),
		},
		{
			name:     `byte slice`,
			input:    []byte(`2006-10-12`),
			expected: timefmt.NewDate[timefmt.ISO8601Date](2006, 10, 12),
		},
		{
			name:     `time.Time`,
			input:    time.Date(2006, 10, 12, 1, 2, 3, 4, time.Local),
			expected: timefmt.NewDate[timefmt.ISO8601Date](2006, 10, 12),
		},
		{
			name:      `wrong type`,
			input:     123345,
			expectErr: `scan failed: cannot unmarshal variable of type int into timefmt.Date[timefmt.ISO8601Date]`,
		},
		{
			name:      `bad format`,
			input:     `213-4=vdf&`,
			expectErr: `scan failed: cannot unmarshal 213-4=vdf& into timefmt.Date[timefmt.ISO8601Date]: parsing time "213-4=vdf&" as "2006-01-02": cannot parse "4=vdf&" as "2006"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var d timefmt.Date[timefmt.ISO8601Date]
			err := d.Scan(c.input)
			if c.expectErr != `` {
				if err == nil {
					t.Errorf(`err should be %q, but its nil`, c.expectErr)
				}

				if err != nil && err.Error() != c.expectErr {
					t.Errorf(`err should be %q but it is %q`, c.expectErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDateValue(t *testing.T) {
	d := timefmt.DateFrom[timefmt.ISO8601Date](times[3].t)
	actual, err := d.Value()
	if err != nil {
		t.Error(err)
	}
	expected := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	if actual != expected {
		t.Errorf(`should be %#v, but its %#v`, expected, actual)
	}
}

func TestDateUnmarshalText(t *testing.T) {
	t.Run(`successful unmarshal`, func(t *testing.T) {
		var actual timefmt.Date[timefmt.ISO8601Date]
		err := actual.UnmarshalText([]byte(`2020-03-01`))
		if err != nil {
			t.Error(err)
		}
		expected := timefmt.NewDate[timefmt.ISO8601Date](2020, 3, 1)
		if actual != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, actual)
		}
	})

	t.Run(`failed unmarshal`, func(t *testing.T) {
		var d timefmt.Date[timefmt.ISO8601Date]
		err := d.UnmarshalText([]byte(`20%0-0.3-01`))
		if err == nil {
			t.Error(`err should be non-nil, but is nil`)
			return
		}
		expected := `parsing time "20%0-0.3-01" as "2006-01-02": cannot parse "-0.3-01" as "2006"`
		if err.Error() != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, err.Error())
		}
	})
}

func TestDateMarshalText(t *testing.T) {
	d := timefmt.NewDate[timefmt.ISO8601Date](2020, 3, 1)
	actual, err := d.MarshalText()
	if err != nil {
		t.Error(err)
	}
	expected := `2020-03-01`
	if string(actual) != expected {
		t.Errorf(`should be %s but is %s`, expected, actual)
	}
}
