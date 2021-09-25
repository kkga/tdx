package vdir

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
	"github.com/hako/durafmt"
)

type ToDoStatus string
type ToDoPriority int

const (
	StatusCompleted   ToDoStatus   = "COMPLETED"
	StatusNeedsAction ToDoStatus   = "NEEDS-ACTION"
	StatusCancelled   ToDoStatus   = "CANCELLED"
	StatusAny         ToDoStatus   = "ANY"
	PriorityHigh      ToDoPriority = 1
	PriorityMedium    ToDoPriority = 5
	PriorityLow       ToDoPriority = 6
)

type FormatOption int

const (
	FormatInline FormatOption = iota
	FormatNoDescription
)

// Item represents an iCalendar item with a unique id
type Item struct {
	Id   int
	Path string
	Ical *ical.Calendar
}

// Init initializes an Item with a decoded ical data from path
func (i *Item) Init(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	i.Path = path

	dec := ical.NewDecoder(file)

	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// filter only items that contain vtodo
		for _, comp := range cal.Children {
			if comp.Name == ical.CompToDo {
				i.Ical = cal
				return nil
			}
		}
	}
	return nil
}

// Vtodo returns a pointer to inner todo ical component
func (i *Item) Vtodo() (*ical.Component, error) {
	for _, comp := range i.Ical.Children {
		if comp.Name == ical.CompToDo {
			return comp, nil
		}
	}
	return nil, fmt.Errorf("Vtodo not found: %s", i.Ical.Name)
}

// Format returns a string representation of an item ready for output
func (i *Item) Format(options ...FormatOption) (string, error) {
	colorStatusDone := color.New(color.Faint).SprintFunc()
	colorStatusUndone := color.New(color.FgBlue).SprintFunc()
	colorPrioHigh := color.New(color.FgHiRed, color.Bold).SprintFunc()
	colorPrioMedium := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	colorDesc := color.New(color.Faint, color.Italic).SprintFunc()
	var vtodo *ical.Component
	for _, comp := range i.Ical.Children {
		if comp.Name == ical.CompToDo {
			vtodo = comp
		}
	}

	if vtodo.Name != ical.CompToDo {
		return "", fmt.Errorf("Not VTODO component: %v", vtodo)
	}

	var (
		status      string
		summary     string
		description string
		prio        string
		due         string
		repeat      string
	)

	for name, prop := range vtodo.Props {
		p := prop[0]

		switch name {

		case ical.PropStatus:
			if ToDoStatus(p.Value) == StatusCompleted {
				status = colorStatusDone("[x]")
			} else {
				status = colorStatusUndone("[ ]")
			}

		case ical.PropSummary:
			summary = p.Value

		case ical.PropDescription:
			description = colorDesc(fmt.Sprintf("%s", p.Value))

		case ical.PropRecurrenceRule:
			c := color.New(color.FgGreen).SprintFunc()
			repeat = c("⟳")

		case ical.PropDue:
			d, err := p.DateTime(time.Local)
			if err != nil {
				return "", err
			}
			now := time.Now()
			diff := d.Sub(now)
			diff = diff.Round(1 * time.Minute)

			var prefix string
			var humanDate string
			var colorizer = color.New(color.Reset).SprintFunc()

			if diff.Hours() < 24 {
				if now.Day() == d.Day() {
					colorizer = color.New(color.FgGreen).SprintFunc()
					prefix = ""
					humanDate = "today"
				} else if math.Signbit(diff.Hours()) {
					colorizer = color.New(color.FgRed).SprintFunc()
					prefix = "overdue "
					humanDate = "yesterday"
				} else {
					colorizer = color.New(color.FgGreen).SprintFunc()
					prefix = ""
					humanDate = "tomorrow"
				}
			} else {
				if math.Signbit(diff.Hours()) {
					prefix = "overdue "
					colorizer = color.New(color.FgRed).SprintFunc()
				} else {
					prefix = "in "
					colorizer = color.New(color.Faint).SprintFunc()
				}

				humanDate = durafmt.ParseShort(diff).String()
			}

			due = colorizer(fmt.Sprintf("(%s%s)", prefix, humanDate))

		case ical.PropPriority:
			v, err := strconv.Atoi(p.Value)
			if err != nil {
				return "", err
			}
			switch {
			case ToDoPriority(v) == PriorityHigh:
				prio = colorPrioHigh("!!!")
			case ToDoPriority(v) > PriorityHigh && ToDoPriority(v) <= PriorityMedium:
				prio = colorPrioMedium("!!")
			case ToDoPriority(v) > PriorityMedium:
				prio = colorPrioMedium("!")
			}

		}
	}

	checkOpt := func(o FormatOption) bool {
		for _, opt := range options {
			if opt == o {
				return true
			}
		}
		return false
	}

	todoSb := strings.Builder{}
	metaSb := strings.Builder{}

	if i.Id != 0 {
		todoSb.WriteString(fmt.Sprintf("%2d", i.Id))
	} else {
		todoSb.WriteString("  ")
	}

	todoSb.WriteString(fmt.Sprintf(" %s", status))

	if prio != "" {
		todoSb.WriteString(fmt.Sprintf(" %s", prio))
	}
	if repeat != "" {
		todoSb.WriteString(fmt.Sprintf(" %s", repeat))
	}

	todoSb.WriteString(fmt.Sprintf(" %s", summary))

	if due != "" {
		metaSb.WriteString(due)
	}

	if description != "" {
		if due != "" {
			metaSb.WriteString(fmt.Sprintf("%s %s", colorDesc("|"), description))
		} else {
			metaSb.WriteString(fmt.Sprintf("%s", description))
		}
	}

	if metaSb.String() != "" {
		var meta string
		if checkOpt(FormatInline) {
			meta = fmt.Sprintf(" %s\n", metaSb.String())
		} else {
			meta = fmt.Sprintf("\n       %s %s\n", colorDesc("↳"), metaSb.String())
		}
		todoSb.WriteString(meta)
	} else {
		todoSb.WriteString("\n")
	}

	return todoSb.String(), nil
}

// WriteFile encodes ical data and writes to file at Item.Path
func (i *Item) WriteFile() error {
	var buf bytes.Buffer
	err := ical.NewEncoder(&buf).Encode(i.Ical)
	if err != nil {
		return err
	}

	// TODO: update modified date prop

	f, err := os.Create(i.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// GenerateUID returns a random string containing timestamp and hostname
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
