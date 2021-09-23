package vdir

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
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

// Item represents an iCalendar item with a unique id
type Item struct {
	Id   int
	Path string
	Ical *ical.Calendar
}

// Init initializes an Item with a decoded ical from path
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

// Format returns a string representation of an item
func (i *Item) Format() (string, error) {
	colorStatusDone := color.New(color.Faint).SprintFunc()
	colorStatusUndone := color.New(color.FgBlue).SprintFunc()
	colorPrioHigh := color.New(color.FgHiRed, color.Bold).SprintFunc()
	colorPrioMedium := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	colorDesc := color.New(color.Faint, color.Italic).SprintFunc()
	colorDate := color.New(color.Faint).SprintFunc()

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

		case ical.PropDue:
			d, _ := p.DateTime(time.Local)
			due = colorDate(d.Format("02-Jan-2006"))

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

	todoSb := strings.Builder{}
	metaSb := strings.Builder{}

	todoSb.WriteString(fmt.Sprintf("%2d %s", i.Id, status))

	if prio != "" {
		todoSb.WriteString(fmt.Sprintf(" %s", prio))
	}

	todoSb.WriteString(fmt.Sprintf(" %s", summary))

	if due != "" || description != "" {
		if due != "" {
			metaSb.WriteString(fmt.Sprintf("%s", due))
		} else if due != "" && description != "" {
			metaSb.WriteString(fmt.Sprintf("%s / %s", due, description))
		} else if description != "" {
			metaSb.WriteString(fmt.Sprintf("%s", description))
		}
	}

	if metaSb.String() != "" {
		meta := fmt.Sprintf("\n       %s %s", colorDesc("↳"), metaSb.String())
		todoSb.WriteString(meta)
	}

	return todoSb.String(), nil
}

// WriteFile enodes ical data and writes to file
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
