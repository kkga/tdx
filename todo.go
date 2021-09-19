package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

type ToDo struct {
	Status      ToDoStatus
	Summary     string
	Description string
	Due         time.Time
	Priority    int
}

type ToDoStatus string
type ToDoPriority int

const (
	ToDoStatusCompleted   ToDoStatus   = "COMPLETED"
	ToDoStatusNeedsAction ToDoStatus   = "NEEDS-ACTION"
	ToDoPriorityLow       ToDoPriority = 9
	ToDoPriorityMedium    ToDoPriority = 14
	ToDoPriorityHigh      ToDoPriority = 1
)

func NewToDo(todo ical.Component) (*ToDo, error) {
	t := &ToDo{}

	for p := range todo.Props {
		switch p {
		case ical.PropStatus:
			s, err := todo.Props.Get(ical.PropStatus).Text()
			if err != nil {
				return nil, err
			}
			t.Status = ToDoStatus(s)
		case ical.PropSummary:
			s, err := todo.Props.Get(ical.PropSummary).Text()
			if err != nil {
				return nil, err
			}
			t.Summary = s
		case ical.PropDescription:
			s, err := todo.Props.Get(ical.PropDescription).Text()
			if err != nil {
				return nil, err
			}
			t.Description = s
		case ical.PropDue:
			time, err := todo.Props.Get(ical.PropDue).DateTime(t.Due.Location())
			if err != nil {
				return nil, err
			}
			t.Due = time
		case ical.PropPriority:
			prio, err := todo.Props.Get(ical.PropPriority).Int()
			if err != nil {
				return nil, err
			}
			t.Priority = prio
		}
	}
	return t, nil
}

func (t ToDo) String() string {
	sb := strings.Builder{}

	if t.Status == ToDoStatusCompleted {
		sb.WriteString("[x]")
	} else if t.Status == ToDoStatusNeedsAction {
		sb.WriteString("[ ]")
	}

	if t.Priority != 0 {
		var prio string
		switch t.Priority {
		case int(ToDoPriorityLow):
			prio = "!"
		case int(ToDoPriorityMedium):
			prio = "!!"
		case int(ToDoPriorityHigh):
			prio = "!!!"
		}
		sb.WriteString(fmt.Sprintf(" (%s)", prio))
	}

	if t.Summary != "" {
		sb.WriteString(" ")
		sb.WriteString(t.Summary)
	}

	if !t.Due.IsZero() {
		date := t.Due.Local().Format(time.RFC822)
		sb.WriteString(fmt.Sprintf(" (%s)", date))
	}

	if t.Description != "" {
		sb.WriteString("\n    ↳ ")
		sb.WriteString(t.Description)
	}

	return sb.String()
}
