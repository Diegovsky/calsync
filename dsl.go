package calsync

import (
	"regexp"
	"strconv"
	"strings"
)

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

func Parse(input string) (events []Event) {
	var currentTag = ""
	for _, line := range regexp.MustCompile("\n+").Split(input, -1) {
		line = strings.TrimSpace(line)
		line, _ := takeUntil(line, "#")
		if line == "" {
			continue
		}
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

		events = append(events, Event{
			When:  date,
			Title: title,
			Tag:   currentTag,
			Id:    id,
		})
	}
	return
}
