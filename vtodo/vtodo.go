package vtodo

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
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

func Format(comp *ical.Component) (string, error) {
	sb := strings.Builder{}

	colorStatusDone := color.New(color.FgGreen).SprintFunc()
	colorStatusUndone := color.New(color.FgBlue).SprintFunc()
	colorPrioHigh := color.New(color.FgHiRed, color.Bold).SprintFunc()
	colorPrioMedium := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	colorDesc := color.New(color.Faint).SprintFunc()
	// colorDate := color.New(color.FgYellow).SprintFunc()

	if comp.Name != ical.CompToDo {
		return "", fmt.Errorf("Not VTODO component: %v", comp)
	}

	var (
		status      string
		summary     string
		description string
		prio        string
		// due         string
	)

	for name, prop := range comp.Props {
		p := prop[0]

		switch name {

		case ical.PropStatus:
			if p.Value == StatusCompleted {
				status = colorStatusDone("[x]")
			} else {
				status = colorStatusUndone("[ ]")
			}

		case ical.PropSummary:
			summary = p.Value

		case ical.PropDescription:
			description = colorDesc(fmt.Sprintf("(%s)", p.Value))

		case ical.PropPriority:
			v, err := strconv.Atoi(p.Value)
			if err != nil {
				return "", err
			}
			switch {
			case v == PriorityHigh:
				prio = colorPrioHigh("!!!")
			case v > PriorityHigh && v <= PriorityMedium:
				prio = colorPrioMedium("!!")
			case v > PriorityMedium:
				prio = colorPrioMedium("!")
			}

		}
	}

	sb.WriteString(fmt.Sprintf("%s", status))
	if prio != "" {
		sb.WriteString(fmt.Sprintf(" %s", prio))
	}
	sb.WriteString(fmt.Sprintf(" %s", summary))
	if description != "" {
		sb.WriteString(fmt.Sprintf(" %s", description))
	}

	return sb.String(), nil
}
