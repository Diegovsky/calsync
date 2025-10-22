package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	. "calsync"

	"google.golang.org/api/calendar/v3"
)

func say(format string, args ...any) {
	fmt.Printf(strings.ReplaceAll(format, "{}", "%#v")+"\n", args...)
}

func ensureCalendar(srv *calendar.Service) string {
	const CALENDAR_SUMMARY_ID = "Calsync Events"
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Failed to list calendars: %v", err)
	}
	for _, cal := range calendars.Items {
		if cal.Summary == CALENDAR_SUMMARY_ID {
			return cal.Id
		}
	}
	cal, err := srv.Calendars.Insert(&calendar.Calendar{
		Summary: CALENDAR_SUMMARY_ID,
	}).Do()
	if err != nil {
		log.Fatalf("Failed to create own calendar: %v", err)
	}
	return cal.Id
}

func editOrDelete(calendarId string, ids map[string]*Event, srv *calendar.Service) error {
	say("Updating existing events...")
	upEvents, err := srv.Events.List(calendarId).Do()
	if err != nil {
		return errors.New("Failed to list our calendar events")
	}

	var toDelete []*calendar.Event
	var wg sync.WaitGroup
	for _, apiEvent := range upEvents.Items {
		if ev, ok := ids[apiEvent.Id]; ok {
			var dirty = false
			if apiEvent.Summary != ev.Title {
				apiEvent.Summary = ev.Title
				dirty = true
			}
			date, err := NewDateFromEventDateTime(apiEvent.Start)
			if err != nil {
				say("Got invalid date from api {}", err)
				continue
			}
			if date != ev.When {
				apiEvent.Start.Date = ev.When.ToDateString()
				apiEvent.End.Date = ev.When.NextDay().ToDateString()
				dirty = true
			}

			if !dirty {
				continue
			}
			wg.Go(func() {

				say("Editing %s", ev.Title)
				if _, err := srv.Events.Patch(calendarId, ev.Id, apiEvent).Do(); err != nil {
					say("Error while updating event %s: %w", apiEvent.Summary, err)
				}
			})
		} else {
			// Non existent ID, should be deleted
			toDelete = append(toDelete, apiEvent)
			event := apiEvent
			wg.Go(func() {
				if err := srv.Events.Delete(calendarId, event.Id).Do(); err != nil {
					say("Failed to delete event %s", event.Summary)
				}
				say("Deleted %s", event.Summary)
			})
		}

	}
	wg.Wait()

	return nil
}

func main() {
	events := GetEvents()
	ctx := context.Background()
	srv := GetServer(ctx)
	calendarId := ensureCalendar(srv)
	say("Got calendar")
	ids := make(map[string]*Event)
	var wg sync.WaitGroup
	var lock sync.Mutex
	putEvent := func(event *Event) {
		lock.Lock()
		ids[event.Id] = event
		lock.Unlock()
	}
	for i := range events {
		ev := &events[i]
		// say("Got event %#v", ev)
		if ev.Id == "" {
			result := srv.Events.Insert(calendarId, &calendar.Event{
				Summary: ev.Title,
				Start: &calendar.EventDateTime{
					Date: ev.When.ToDateString(),
				},
				End: &calendar.EventDateTime{
					Date: ev.When.NextDay().ToDateString(),
				},
			})
			wg.Go(func() {
				e, err := result.Do()
				if err != nil {
					log.Fatalf("Failed adding event '%s': %v [%v]", ev.Title, err, e)
				}
				ev.Id = e.Id
				say("Created event %s", ev.Title)
				putEvent(ev)
			})
		} else {
			putEvent(ev)
		}
	}
	wg.Wait()

	if err := editOrDelete(calendarId, ids, srv); err != nil {
		say("Error while syncing existing events: ", err)
	}
	SaveEvents(events)
}
