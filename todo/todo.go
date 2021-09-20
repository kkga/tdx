package todo

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
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

type icalToDo struct {
	*ical.Component
}

func GenerateUID() string {
	sb := strings.Builder{}

	randStr := func(n int) string {
		var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

		s := make([]rune, n)
		for i := range s {
			s[i] = letters[rand.Intn(len(letters))]
		}
		return string(s)
	}

	sb.WriteString(fmt.Sprint(time.Now().UnixNano()))
	sb.WriteString(fmt.Sprintf("-%s", randStr(8)))
	if hostname, _ := os.Hostname(); hostname != "" {
		sb.WriteString(fmt.Sprintf("@%s", hostname))
	}

	return sb.String()
}

func (t ToDo) Encode() (bytes.Buffer, error) {
	icalToDo := &icalToDo{ical.NewComponent(ical.CompToDo)}
	icalToDo.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())

	if t.UID == "" {
		icalToDo.Props.SetText(ical.PropUID, GenerateUID())
	} else {
		icalToDo.Props.SetText(ical.PropUID, t.UID)
	}
	icalToDo.Props.SetText(ical.PropSummary, "testing encode")

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, icalToDo.Component)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return buf, err
	}
	return buf, nil
}
