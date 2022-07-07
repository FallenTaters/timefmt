package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/FallenTaters/timefmt"
)

type customFormat struct{}

func (customFormat) TimeFormat() string {
	return `2006/01/02 03:04 PM`
}

func (customFormat) DateFormat() string {
	return `02 Jan 2006`
}

type Payload struct {
	Date timefmt.Date[customFormat] `json:"date"`
	Time timefmt.Time[customFormat] `json:"time"`
}

func main() {
	d := timefmt.NewDate[customFormat](2022, 7, 4)
	fmt.Println(d) // 04 Jul 2022

	t := timefmt.TimeFrom[customFormat](time.Date(2022, 7, 4, 16, 19, 59, 1_000_000_000, time.UTC))
	fmt.Println(t) // 2022/07/04 04:20 PM

	// json
	jsonText := `{"date":"04 Jul 2022","time":"2022/07/04 04:20 PM"}`
	var p Payload
	err := json.Unmarshal([]byte(jsonText), &p)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data) == jsonText) // true

	// dates of different times are comparable
	d1 := time.Date(2022, 7, 4, 0, 0, 0, 0, time.UTC)
	d2 := d1.Add(time.Hour + time.Minute + time.Second)
	fmt.Println(timefmt.DateFrom[customFormat](d1) == timefmt.DateFrom[customFormat](d2)) // true
}

// Some useful formats are predefined
type predefined struct {
	ISO8601      timefmt.Date[timefmt.ISO8601Date]  // 2006-01-02
	DayMonthYear timefmt.Date[timefmt.DayMonthYear] // 02/01/2006
	MonthDayYear timefmt.Date[timefmt.MonthDayYear] // 01/02/2006

	RFC3339          timefmt.Time[timefmt.RFC3339]          // 2006-01-02T15:04:05Z07:00
	YearMonthDayTime timefmt.Time[timefmt.YearMonthDayTime] // 2006-01-02 15:04:05
}
