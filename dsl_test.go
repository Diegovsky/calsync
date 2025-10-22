package calsync

import (
	"testing"
)

func assertEq(a any, b any, format string, t *testing.T) {
	if a != b {
		t.Errorf(format, a, b)
	}
}

func expectParse(input string, expectedEvents []Event, t *testing.T) {
	events := Parse(input)
	assertEq(len(expectedEvents), len(events), "Expected '%d' events, got '%d'", t)
	for i := range events {
		expected := expectedEvents[i]
		event := events[i]
		assertEq(expected, event, "Expected\n'%#v'\ngot\n'%#v'", t)
	}
}

func expectReparse(events []Event, t *testing.T) {
	gened := EventsToDSL(events)
	expectParse(gened, events, t)
	if t.Failed() {
		t.Errorf("Generated string:\n%s\n\n", gened)
	}
}

func TestNoId(t *testing.T) {
	expectParse(`
				 30/12 - a

				 [Tag]
				 1/12 - EventNameWithNoSpaces
				 10/12 - Event Name With A lot of Spaces


				 `, []Event{
		{
			Title: "a",
			When: Date{
				Day:   30,
				Month: 12,
			},
		},
		{
			Title: "EventNameWithNoSpaces",
			When: Date{
				Day:   1,
				Month: 12,
			},
			Tag: "Tag",
		},
		{
			Title: "Event Name With A lot of Spaces",
			When: Date{
				Day:   10,
				Month: 12,
			},
			Tag: "Tag",
		},
	}, t)
}

func TestWithId(t *testing.T) {
	expectParse(`10/10 - Event @id1
				 11/11 - Event @id2
				 `, []Event{
		{
			Title: "Event",
			Id:    "id1",
			When: Date{
				Day:   10,
				Month: 10,
			},
		}, {
			Title: "Event",
			Id:    "id2",
			When: Date{
				Day:   11,
				Month: 11,
			},
		},
	}, t)
}

func TestSpaces(t *testing.T) {
	expectParse(`

		10/12 - Event Name


		
				 11/12 - Event Name
				 `, []Event{
		{
			Title: "Event Name",
			When: Date{
				Day:   10,
				Month: 12,
			},
		}, {
			Title: "Event Name",
			When: Date{
				Day:   11,
				Month: 12,
			},
		},
	}, t)
}

func TestComments(t *testing.T) {
	expectParse(`

		10/12 - Event Name # Comment


		# comment # comment
				 11/12 - Event Name #@293039
				 `, []Event{
		{
			Title: "Event Name",
			When: Date{
				Day:   10,
				Month: 12,
			},
		}, {
			Title: "Event Name",
			When: Date{
				Day:   11,
				Month: 12,
			},
		},
	}, t)
}
func TestReparseSimple(t *testing.T) {
	expectReparse([]Event{
		{
			Title: "Event Name",
			When: Date{
				Day:   10,
				Month: 12,
			},
		}, {
			Title: "Event Name",
			When: Date{
				Day:   11,
				Month: 12,
			},
		},
	}, t)
}

func TestReparseEmpty(t *testing.T) {
	expectReparse([]Event{}, t)
}

func TestReparseWithId(t *testing.T) {
	expectReparse([]Event{
		{
			Title: "Event Name",
			Id:    "32492304898234nakjnsas",
			When: Date{
				Day:   10,
				Month: 12,
			},
		},
	}, t)

	expectReparse([]Event{
		{
			Title: "Event Name",
			Id:    "1",
			When: Date{
				Day:   10,
				Month: 12,
			},
		},
		{
			Title: "Event Name",
			Id:    "2832030",
			When: Date{
				Day:   10,
				Month: 12,
			},
		},

		{
			Title: "Event Name",
			Tag:   "a",
			When: Date{
				Day:   10,
				Month: 12,
			},
		}, {
			Title: "Event Name2",
			Tag:   "a",
			When: Date{
				Day:   11,
				Month: 12,
			},
		},
	}, t)
}
