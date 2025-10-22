package calsync

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"

	cp "github.com/otiai10/copy"
)

var inputFile string = "events.cal"

func GetEvents() []Event {
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
	currentYear := time.Now().Year()
	for i := range events {
		ev := &events[i]
		if ev.When.Year == 0 {
			ev.When.Year = currentYear
		}
	}
	return events
}

func SaveEvents(events []Event) {
	os.WriteFile(inputFile, []byte(EventsToDSL(events)), 0666)
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

func sortByDate(events []Event) {
	sort.Slice(events, func(a, b int) bool {
		return events[a].When.Less(events[b].When)
	})

}

func EventsToDSL(events []Event) string {
	var sections OrderMap[string, []Event] = NewOrderMap[string, []Event]()
	for _, ev := range events {
		sections.Update(ev.Tag, func(events []Event) []Event { return append(events, ev) })
	}
	var builder strings.Builder
	for k, section := range sections.Items() {

		if k != "" {
			if builder.Len() != 0 {
				builder.WriteRune('\n')
			}
			fmt.Fprintf(&builder, "[%s]\n", k)
		}
		sortByDate(section)
		for _, event := range section {
			event.toDSL(&builder)
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

func (self Event) toDSL(builder *strings.Builder) {
	fmt.Fprintf(builder, `%d/%d`, self.When.Day, self.When.Month)
	if self.When.Year != 0 {
		fmt.Fprintf(builder, "/%d", self.When.Year)
	}
	fmt.Fprintf(builder, " - %s", self.Title)
	if self.Id != "" {
		fmt.Fprintf(builder, ` @%s`, self.Id)
	}
}
