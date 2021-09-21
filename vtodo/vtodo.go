package vtodo

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

const (
	StatusCompleted   = "COMPLETED"
	StatusNeedsAction = "NEEDS-ACTION"
	PriorityHigh      = 1
	PriorityMedium    = 5
	PriorityLow       = 6
)

// generateUID returns a random string containing timestamp and hostname
func GenerateUID() string {
	sb := strings.Builder{}

	time := time.Now().UnixNano()

	randStr := func(n int) string {
		rs := rand.NewSource(time)
		r := rand.New(rs)
		var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

		s := make([]rune, n)
		for i := range s {
			s[i] = letters[r.Intn(len(letters))]
		}
		return string(s)
	}

	sb.WriteString(fmt.Sprint(time))
	sb.WriteString(fmt.Sprintf("-%s", randStr(8)))
	if hostname, _ := os.Hostname(); hostname != "" {
		sb.WriteString(fmt.Sprintf("@%s", hostname))
	}

	return sb.String()
}

// encode adds vtodo into a new Calendar and returns a buffer ready for writing
func Encode(vtodo *ical.Component) (*bytes.Buffer, error) {
	cal := ical.NewCalendar()
	// TODO move this data somewhere
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, vtodo)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return &buf, err
	}
	return &buf, nil
}

// TODO: refactor this for using the ical.Component directly
func ToString(vtodo ical.Component) string {
	// sb := strings.Builder{}
	// colorPrio := color.New(color.FgRed, color.Bold).SprintFunc()
	// colorDate := color.New(color.FgYellow).SprintFunc()
	// colorDesc := color.New(color.Faint).SprintFunc()

	// if t.Status == ToDoStatusCompleted {
	// 	sb.WriteString("[x]")
	// } else if t.Status == ToDoStatusNeedsAction {
	// 	sb.WriteString("[ ]")
	// }

	// if t.Priority != 0 {
	// 	var prio string
	// 	switch {
	// 	case t.Priority == ToDoPriorityHigh:
	// 		prio = "!!!"
	// 	case t.Priority > ToDoPriorityHigh && t.Priority <= ToDoPriorityMedium:
	// 		prio = "!!"
	// 	case t.Priority > ToDoPriorityMedium:
	// 		prio = "!"
	// 	}
	// 	sb.WriteString(fmt.Sprintf(" %s", colorPrio(prio)))
	// }

	// if !t.Due.IsZero() {
	// 	date := t.Due.Local().Format("Jan-06")
	// 	sb.WriteString(fmt.Sprintf(" %s", colorDate(date)))
	// }

	// if t.Summary != "" {
	// 	sb.WriteString(" ")
	// 	sb.WriteString(t.Summary)
	// }

	// if t.Description != "" {
	// 	sb.WriteString(colorDesc("\n    ↳ "))
	// 	sb.WriteString(colorDesc(t.Description))
	// }

	// return sb.String()
	return ""
}
