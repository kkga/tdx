package main

import (
	"fmt"
	"strings"
	"time"
)

type Todo struct {
	Status      TodoStatus
	Summary     string
	Description string
	Due         time.Time
}

type TodoStatus string

const (
	TodoCompleted   TodoStatus = "COMPLETED"
	TodoNeedsACtion TodoStatus = "NEEDS-ACTION"
)

func (t Todo) String() string {
	sb := strings.Builder{}

	if t.Status == TodoCompleted {
		sb.WriteString("[x]")
	} else if t.Status == TodoNeedsACtion {
		sb.WriteString("[ ]")
	}

	if t.Summary != "" {
		sb.WriteString(" ")
		sb.WriteString(t.Summary)
	}

	if !t.Due.IsZero() {
		sb.WriteString(fmt.Sprintf(" (%s)", t.Due.String()))
	}

	if t.Description != "" {
		sb.WriteString("\n    ↳ ")
		sb.WriteString(t.Description)
	}

	return sb.String()
}
