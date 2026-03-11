package calsync

import (
	"fmt"
	"iter"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DocumentItem struct {
	// Same as the Event.tag, except it is also used to put comments where they belong
	Tag     string
	Event   *Event
	Comment string
}

type Document struct {
	Items []DocumentItem
}

// Returns all the events contained in the calendar file
func (self Document) IterEvents() iter.Seq[*Event] {
	return func(yield func(*Event) bool) {
		for _, item := range self.Items {
			if item.Event != nil {
				if !yield(item.Event) {
					return
				}
			}
		}
	}
}

func stripParts(arr []string) {
	for i, part := range arr {
		arr[i] = strings.TrimSpace(part)
	}
}

func takeUntil(self string, chars string) (string, string) {
	newStart := strings.IndexAny(self, chars)
	if newStart < 0 {
		return self, ""
	}
	return strings.TrimSpace(self[:newStart]), strings.TrimSpace(self[newStart:])
}

func between(self string, starter string, ender string) (string, string) {
	rest, self := takeUntil(self, starter)
	if self == "" {
		return "", rest
	}
	self = self[1:]
	content, rest := takeUntil(self, ender)
	if rest != "" {
		rest = rest[1:]
	}
	return content, rest
}

func Parse(input string) Document {
	var currentTag = ""
	var items []DocumentItem
	for _, line := range regexp.MustCompile("\n+").Split(input, -1) {
		line = strings.TrimSpace(line)
		// skip empty lines
		if line == "" {
			continue
		}
		var item DocumentItem
		line, comment := takeUntil(line, "#")

		if comment != "" {
			item.Comment = comment[1:]
		}

		if line != "" {
			item.Tag = currentTag

			if strings.HasPrefix(line, "[") {
				subject, _ := between(line, "[", "]")
				currentTag = subject
				continue
			}
			parts := strings.SplitN(line, "-", 2)
			if len(parts) != 2 {
				continue
			}
			stripParts(parts)

			date_str := strings.Split(parts[0], "/")
			day, _ := strconv.Atoi(date_str[0])
			month, _ := strconv.Atoi(date_str[1])
			year := 0
			if len(date_str) == 3 {
				year, _ = strconv.Atoi(date_str[2])
			}

			date := Date{Day: day, Month: month, Year: year}

			title, rest := takeUntil(parts[1], "@[")
			id, rest := between(rest, "@", " ")

			item.Event = &Event{
				When:  date,
				Title: title,
				Tag:   currentTag,
				Id:    id,
			}
		}

		items = append(items, item)

	}
	return Document{Items: items}
}

func (self Document) replaceWithCurrentYear() {
	currentYear := time.Now().Year()
	for ev := range self.IterEvents() {
		if ev.When.Year == 0 {
			ev.When.Year = currentYear
		}
	}

}

func (events Document) ToDSL() string {
	var sections OrderMap[string, []DocumentItem] = NewOrderMap[string, []DocumentItem]()
	for _, item := range events.Items {
		if item.Event != nil {
			sections.Update(item.Tag, func(items []DocumentItem) []DocumentItem { return append(items, item) })
		}
	}
	var builder strings.Builder
	for k, section := range sections.Items() {
		if k != "" {
			if builder.Len() != 0 {
				builder.WriteRune('\n')
			}
			fmt.Fprintf(&builder, "[%s]\n", k)
		}
		sort.Slice(section, func(a, b int) bool {
			if eventA, eventB := section[a].Event, section[b].Event; eventA != nil && eventB != nil {
				return eventA.When.Less(eventB.When)

			}

			// Else, use comparison by index
			return a < b
		})

		for _, item := range section {
			if event := item.Event; event != nil {
				event.toDSL(&builder)
				builder.WriteRune('\n')
			}
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
