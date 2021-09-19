package main

import (
	"fmt"
	"strings"
	"time"
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
