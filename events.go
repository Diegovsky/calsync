package calsync

import (
	"log"
	"os"
	"time"

	"google.golang.org/api/calendar/v3"

	cp "github.com/otiai10/copy"
)

var inputFile string = "events.cal"

func GetEvents() Document {
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	b, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file %s: %v", inputFile, err)
	}
	if err := cp.Copy(inputFile, inputFile+".bak"); err != nil {
		log.Fatalf("Failed to copy: %v", err)
	}
	events := Parse(string(b))
	return events
}

func SaveEvents(events Document) {
	os.WriteFile(inputFile, []byte(events.ToDSL()), 0666)
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d Date) Less(other Date) bool {
	return d.Year < other.Year ||
		(d.Year == other.Year && d.Month < other.Month) ||
		(d.Year == other.Year && d.Month == other.Month && d.Day < other.Day)
}

func (date Date) NextDay() Date {
	date.Day += 1
	return date
}

const DATE_FMT = "2006-01-02"

func (date Date) ToDateString() string {
	return time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.Local).Format(DATE_FMT)
}

// Accepts both YYYY-MM-DD and RFC3339
func NewDateFromString(input string) (date Date, err error) {
	timeDate, err := time.Parse(DATE_FMT, input)
	if err != nil {
		timeDate, err = time.Parse(time.RFC3339, input)
		if err != nil {
			return
		}
	}
	date.Day = timeDate.Day()
	date.Month = int(timeDate.Month())
	date.Year = timeDate.Year()
	return
}

func NewDateFromEventDateTime(eventDate *calendar.EventDateTime) (date Date, err error) {
	date, err = NewDateFromString(eventDate.Date)
	if err == nil {
		return
	}
	date, err = NewDateFromString(eventDate.DateTime)
	return
}

type Event struct {
	When  Date
	Title string
	Tag   string
	Id    string
}
