package calsync

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func (self DocumentItem) Format(f fmt.State, verb rune) {
	spew.Fdump(f, self)
}

func assertEq(a any, b any, format string, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf(format, a, b)
		panic("cringe")
	}
}
func expectParse(input string, expectedEvents []Event, t *testing.T) {
	expectParseDocument(input, makeDocument(expectedEvents), true, t)
}

func expectParseItems(input string, items []DocumentItem, t *testing.T) {
	expectParseDocument(input, Document{Items: items}, false, t)
}

func expectParseDocument(input string, expected Document, ignoreComments bool, t *testing.T) {
	gotten := Parse(input)

	j := 0
	for i := range gotten.Items {
		gotten := gotten.Items[i]
		expected := expected.Items[j]

		if ignoreComments {
			gotten.Comment = ""
			expected.Comment = ""

			if gotten.Event == nil {
				continue
			}
		}

		assertEq(gotten, expected, "Expected:\n%#v\nGot:\n\n%#v", t)
		j += 1
	}
}

func makeDocument(events []Event) (es Document) {
	for _, e := range events {
		es.Items = append(es.Items, DocumentItem{Event: &e, Tag: e.Tag})
	}
	return
}

func expectReparse(events []Event, t *testing.T) {
	gened := makeDocument(events).ToDSL()
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

func TestWithInlineComments(t *testing.T) {
	expectParseItems(`
			10/12 - Uncommentated event
			10/13 - Commentated event # Important
		`, []DocumentItem{
		{
			Event: &Event{
				When: Date{
					Day:   10,
					Month: 12,
				},
				Title: "Uncommentated event",
			},
			Comment: "",
		},
		{
			Event: &Event{
				When: Date{
					Day:   10,
					Month: 13,
				},
				Title: "Commentated event",
			},
			Comment: " Important",
		},
	}, t)
}

func TestWithComments(t *testing.T) {
	expectParseItems(`
			10/12 - Uncommentated event
			# comment
 			10/13 - Commentated event
		`, []DocumentItem{
		{
			Event: &Event{
				When: Date{
					Day:   10,
					Month: 12,
				},
				Title: "Uncommentated event",
			},
		},
		{
			Comment: " comment",
		},
		{
			Event: &Event{
				When: Date{
					Day:   10,
					Month: 13,
				},
				Title: "Commentated event",
			},
		},
	}, t)
}
