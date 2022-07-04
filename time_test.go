package timefmt_test

import (
	"testing"
	"time"

	"github.com/FallenTaters/timefmt"
	"github.com/FallenTaters/timefmt/formats/rfc3339"
)

func TestTimeFrom(t *testing.T) {
	for _, c := range times {
		t.Run(c.name, func(t *testing.T) {
			if d := timefmt.TimeFrom[rfc3339.Format](c.t); d.String() != c.t.Format((rfc3339.Format{}).TimeFormat()) {
				t.Errorf(`date %s should be %q but is %q`, c.name, d.String(), c.t.Format((rfc3339.Format{}).TimeFormat()))
			}
		})
	}
}

func TestTimeGoString(t *testing.T) {
	expected := `timefmt.TimeFrom[rfc3339.Format](time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))`
	if actual := timefmt.TimeFrom[rfc3339.Format](times[0].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.TimeFrom[rfc3339.Format](time.Date(-2, time.October, 30, 0, 0, 0, 0, time.UTC))`
	if actual := timefmt.TimeFrom[rfc3339.Format](times[1].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.TimeFrom[rfc3339.Format](time.Date(100000, time.December, 25, 0, 0, 0, 0, time.UTC))`
	if actual := timefmt.TimeFrom[rfc3339.Format](times[2].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `timefmt.TimeFrom[rfc3339.Format](time.Date(2020, time.March, 1, 12, 13, 14, 151617, time.Local))`
	if actual := timefmt.TimeFrom[rfc3339.Format](times[3].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}
}

func TestTimeComparable(t *testing.T) {
	actual, expected := timefmt.TimeFrom[rfc3339.Format](times[3].t), timefmt.TimeFrom[rfc3339.Format](times[3].t)
	if expected != actual {
		t.Error(actual, expected)
	}
}

func TestTimeScan(t *testing.T) {
	cases := []struct {
		name      string
		input     any
		expectErr string
		expected  timefmt.Time[rfc3339.Format]
	}{
		{
			name:     `string date`,
			input:    `2006-10-12T01:02:03Z`,
			expected: timefmt.TimeFrom[rfc3339.Format](time.Date(2006, 10, 12, 1, 2, 3, 0, time.UTC)),
		},
		{
			name:     `byte slice`,
			input:    []byte(`2006-10-12T01:02:03Z`),
			expected: timefmt.TimeFrom[rfc3339.Format](time.Date(2006, 10, 12, 1, 2, 3, 0, time.UTC)),
		},
		{
			name:     `time.Time`,
			input:    time.Date(2006, 10, 12, 1, 2, 3, 0, time.UTC),
			expected: timefmt.TimeFrom[rfc3339.Format](time.Date(2006, 10, 12, 1, 2, 3, 0, time.UTC)),
		},
		{
			name:      `wrong type`,
			input:     123345,
			expectErr: `scan failed: cannot unmarshal variable of type int into timefmt.Time[rfc3339.Format]`,
		},
		{
			name:      `bad format`,
			input:     `213-4=vdf&`,
			expectErr: `scan failed: cannot unmarshal 213-4=vdf& into timefmt.Time[rfc3339.Format]: parsing time "213-4=vdf&" as "2006-01-02T15:04:05Z07:00": cannot parse "4=vdf&" as "2006"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var d timefmt.Time[rfc3339.Format]
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

func TestTimeValue(t *testing.T) {
	d := timefmt.TimeFrom[rfc3339.Format](times[3].t)
	actual, err := d.Value()
	if err != nil {
		t.Error(err)
	}
	expected := time.Date(2020, 3, 1, 12, 13, 14, 151617, time.Local)
	if actual != expected {
		t.Errorf(`should be %#v, but its %#v`, expected, actual)
	}
}

func TestTimeUnmarshalText(t *testing.T) {
	t.Run(`successful unmarshal`, func(t *testing.T) {
		var actual timefmt.Time[rfc3339.Format]
		err := actual.UnmarshalText([]byte(`2020-03-01T01:02:03Z`))
		if err != nil {
			t.Error(err)
		}
		expected := timefmt.TimeFrom[rfc3339.Format](time.Date(2020, 3, 1, 1, 2, 3, 0, time.UTC))
		if actual != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, actual)
		}
	})

	t.Run(`failed unmarshal`, func(t *testing.T) {
		var d timefmt.Time[rfc3339.Format]
		err := d.UnmarshalText([]byte(`20%0-0.3-01`))
		if err == nil {
			t.Error(`err should be non-nil, but is nil`)
			return
		}
		expected := `parsing time "20%0-0.3-01" as "2006-01-02T15:04:05Z07:00": cannot parse "-0.3-01" as "2006"`
		if err.Error() != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, err.Error())
		}
	})
}

func TestTimeMarshalText(t *testing.T) {
	d := timefmt.TimeFrom[rfc3339.Format](time.Date(2020, 3, 1, 1, 2, 3, 0, time.UTC))
	actual, err := d.MarshalText()
	if err != nil {
		t.Error(err)
	}
	expected := `2020-03-01T01:02:03Z`
	if string(actual) != expected {
		t.Errorf(`should be %s but is %s`, expected, actual)
	}
}
