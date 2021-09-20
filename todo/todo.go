package todo

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
)

type ToDo struct {
	UID         string
	Status      ToDoStatus
	Priority    ToDoPriority
	Summary     string
	Description string
	Tags        []string
	Due         time.Time
	OtherProps  ical.Props
}

type (
	ToDoStatus   string
	ToDoPriority int
	ToDoUID      string
)

const (
	ToDoStatusCompleted   ToDoStatus   = "COMPLETED"
	ToDoStatusNeedsAction ToDoStatus   = "NEEDS-ACTION"
	ToDoPriorityHigh      ToDoPriority = 1
	ToDoPriorityMedium    ToDoPriority = 5
	ToDoPriorityLow       ToDoPriority = 6
)

func NewToDo() *ToDo {
	return &ToDo{}
}

func (t *ToDo) Init(todo ical.Component) error {
	props := todo.Props
	for p := range props {
		switch p {
		case ical.PropUID:
			uid, err := todo.Props.Get(ical.PropUID).Text()
			if err != nil {
				return err
			}
			t.UID = uid
			delete(props, p)
		case ical.PropStatus:
			s, err := todo.Props.Get(ical.PropStatus).Text()
			if err != nil {
				return err
			}
			t.Status = ToDoStatus(s)
			delete(props, p)
		case ical.PropSummary:
			s, err := todo.Props.Get(ical.PropSummary).Text()
			if err != nil {
				return err
			}
			t.Summary = s
			delete(props, p)
		case ical.PropDescription:
			s, err := todo.Props.Get(ical.PropDescription).Text()
			if err != nil {
				return err
			}
			t.Description = s
			delete(props, p)
		case ical.PropDue:
			time, err := todo.Props.Get(ical.PropDue).DateTime(t.Due.Location())
			if err != nil {
				return err
			}
			t.Due = time
			delete(props, p)
		case ical.PropPriority:
			prio, err := todo.Props.Get(ical.PropPriority).Int()
			if err != nil {
				return err
			}
			t.Priority = ToDoPriority(prio)
			delete(props, p)
		}
	}
	t.OtherProps = props
	return nil
}

func (t ToDo) String() string {
	sb := strings.Builder{}
	colorPrio := color.New(color.FgRed, color.Bold).SprintFunc()
	colorDate := color.New(color.FgYellow).SprintFunc()
	colorDesc := color.New(color.Faint).SprintFunc()

	if t.Status == ToDoStatusCompleted {
		sb.WriteString("[x]")
	} else if t.Status == ToDoStatusNeedsAction {
		sb.WriteString("[ ]")
	}

	if t.Priority != 0 {
		var prio string
		switch {
		case t.Priority == ToDoPriorityHigh:
			prio = "!!!"
		case t.Priority > ToDoPriorityHigh && t.Priority <= ToDoPriorityMedium:
			prio = "!!"
		case t.Priority > ToDoPriorityMedium:
			prio = "!"
		}
		sb.WriteString(fmt.Sprintf(" %s", colorPrio(prio)))
	}

	if !t.Due.IsZero() {
		date := t.Due.Local().Format("Jan-06")
		sb.WriteString(fmt.Sprintf(" %s", colorDate(date)))
	}

	if t.Summary != "" {
		sb.WriteString(" ")
		sb.WriteString(t.Summary)
	}

	if t.Description != "" {
		sb.WriteString(colorDesc("\n    ↳ "))
		sb.WriteString(colorDesc(t.Description))
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
