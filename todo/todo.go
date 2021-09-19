package todo

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

type ToDo struct {
	UID         string
	Status      ToDoStatus
	Priority    ToDoPriority
	Summary     string
	Description string
	Tags        []string
	Due         time.Time
}

type (
	ToDoStatus   string
	ToDoPriority int
	ToDoUID      string
)

const (
	ToDoStatusCompleted   ToDoStatus   = "COMPLETED"
	ToDoStatusNeedsAction ToDoStatus   = "NEEDS-ACTION"
	ToDoPriorityLow       ToDoPriority = 9
	ToDoPriorityMedium    ToDoPriority = 14
	ToDoPriorityHigh      ToDoPriority = 1
)

func NewToDo() *ToDo {
	return &ToDo{}
}

func (t *ToDo) Init(todo ical.Component) error {
	for p := range todo.Props {
		switch p {
		case ical.PropUID:
			uid, err := todo.Props.Get(ical.PropUID).Text()
			if err != nil {
				return err
			}
			t.UID = uid
		case ical.PropStatus:
			s, err := todo.Props.Get(ical.PropStatus).Text()
			if err != nil {
				return err
			}
			t.Status = ToDoStatus(s)
		case ical.PropSummary:
			s, err := todo.Props.Get(ical.PropSummary).Text()
			if err != nil {
				return err
			}
			t.Summary = s
		case ical.PropDescription:
			s, err := todo.Props.Get(ical.PropDescription).Text()
			if err != nil {
				return err
			}
			t.Description = s
		case ical.PropDue:
			time, err := todo.Props.Get(ical.PropDue).DateTime(t.Due.Location())
			if err != nil {
				return err
			}
			t.Due = time
		case ical.PropPriority:
			prio, err := todo.Props.Get(ical.PropPriority).Int()
			if err != nil {
				return err
			}
			t.Priority = ToDoPriority(prio)
		}
	}
	return nil
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
		case ToDoPriorityLow:
			prio = "!"
		case ToDoPriorityMedium:
			prio = "!!"
		case ToDoPriorityHigh:
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
