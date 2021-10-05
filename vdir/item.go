package vdir

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
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
	StatusCompleted   ToDoStatus = "COMPLETED"
	StatusNeedsAction ToDoStatus = "NEEDS-ACTION"
	StatusCancelled   ToDoStatus = "CANCELLED"
	StatusInProcess   ToDoStatus = "IN-PROCESS"
	StatusAny         ToDoStatus = "ANY"

	PriorityHigh   ToDoPriority = 1
	PriorityMedium ToDoPriority = 5
	PriorityLow    ToDoPriority = 6
)

const HashtagRe = "\\B#\\w+"

type FormatOption int

const (
	FormatMultiline FormatOption = iota
	FormatDescription
)

type FormatFullOption int

const (
	FormatFullRaw FormatFullOption = iota
)

// Item represents an iCalendar item with a unique id
type Item struct {
	Id   int
	Path string
	Ical *ical.Calendar
}

// Tag represents a hashtag label in todo summary
type Tag string

// DecodeError represents and error occured during ical decoding
type DecodeError struct {
	Path string
	Err  error
}

func (d *DecodeError) Error() string {
	return fmt.Sprintf("%s\npath: %s", d.Err, d.Path)
}

// String returns a lowercased tag string
func (t Tag) String() string {
	return strings.ToLower(string(t))
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
			return &DecodeError{
				i.Path,
				err,
			}
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
	return nil, fmt.Errorf("Vtodo not found: %q", i.Ical.Name)
}

// FormatFull returns a full detailed info about an item
func (i *Item) FormatFull(options ...FormatFullOption) (string, error) {
	sb := strings.Builder{}

	checkOpt := func(o FormatFullOption) bool {
		for _, opt := range options {
			if opt == o {
				return true
			}
		}
		return false
	}

	if checkOpt(FormatFullRaw) {
		j, err := json.MarshalIndent(i, "", " ")
		if err != nil {
			return "", err
		}
		sb.WriteString(string(j))
	} else {
		vtodo, err := i.Vtodo()
		if err != nil {
			return "", nil
		}
		sb.WriteString(fmt.Sprintf("ID: %d\n", i.Id))
		for name, prop := range vtodo.Props {
			p := prop[0]
			if date, _ := p.DateTime(time.Local); !date.IsZero() {
				sb.WriteString(fmt.Sprintf("%s: %s\n", name, date.Format("2 Jan 2006 15:04")))
			} else {
				sb.WriteString(fmt.Sprintf("%s: %s\n", name, p.Value))
			}
		}
	}

	return sb.String(), nil
}

// Format returns a readable representation of an item ready for output
func (i *Item) Format(options ...FormatOption) (string, error) {

	vtodo, err := i.Vtodo()
	if err != nil {
		return "", err
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
			colDone := color.New(color.Faint).SprintFunc()
			colUndone := color.New(color.FgBlue).SprintFunc()
			if ToDoStatus(p.Value) == StatusCompleted {
				status = colDone("[x]")
			} else if ToDoStatus(p.Value) == StatusCancelled {
				status = colDone("[-]")
			} else {
				status = colUndone("[ ]")
			}

		case ical.PropSummary:
			summary = p.Value

			tags, err := i.Tags()
			if err != nil {
				return "", err
			}
			if len(tags) > 0 {
				c := color.New(color.FgBlue).SprintFunc()
				for _, t := range tags {
					summary = strings.ReplaceAll(summary, string(t), c(t))
				}
			}

		case ical.PropDescription:
			col := color.New(color.Faint, color.Italic).SprintFunc()
			description = col(fmt.Sprintf("%s", p.Value))

		case ical.PropRecurrenceRule:
			c := color.New(color.FgGreen).SprintFunc()
			repeat = c("⟳")

		case ical.PropDue:
			d, _ := p.DateTime(time.Local)
			if d.IsZero() {
				continue
			}
			now := time.Now()
			diff := d.Sub(now)
			diff = diff.Round(1 * time.Minute)

			var prefix string
			var humanDate string
			var col = color.New(color.Reset).SprintFunc()

			if math.Abs(diff.Hours()) < 24 {
				if now.Day() == d.Day() {
					col = color.New(color.FgGreen).SprintFunc()
					prefix = ""
					humanDate = "today"
				} else if math.Signbit(diff.Hours()) {
					col = color.New(color.FgRed).SprintFunc()
					prefix = "overdue "
					humanDate = "yesterday"
				} else {
					col = color.New(color.FgGreen).SprintFunc()
					prefix = ""
					humanDate = "tomorrow"
				}
			} else {
				if math.Signbit(diff.Hours()) {
					prefix = "overdue "
					col = color.New(color.FgRed).SprintFunc()
				} else {
					prefix = "in "
					col = color.New(color.Faint).SprintFunc()
				}

				humanDate = durafmt.ParseShort(diff).String()
			}

			humanDate = strings.TrimPrefix(humanDate, "-")
			due = col(fmt.Sprintf("(%s%s)", prefix, humanDate))

		case ical.PropPriority:
			colHigh := color.New(color.FgHiRed, color.Bold).SprintFunc()
			colMedium := color.New(color.FgHiYellow, color.Bold).SprintFunc()
			v, err := strconv.Atoi(p.Value)
			if err != nil {
				return "", err
			}
			switch {
			case ToDoPriority(v) == PriorityHigh:
				prio = colHigh("!!!")
			case ToDoPriority(v) > PriorityHigh && ToDoPriority(v) <= PriorityMedium:
				prio = colMedium("!!")
			case ToDoPriority(v) > PriorityMedium:
				prio = colMedium("!")
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

	if checkOpt(FormatDescription) && description != "" {
		colFaint := color.New(color.Faint).SprintFunc()
		if due != "" {
			metaSb.WriteString(fmt.Sprintf(" %s %s", colFaint("|"), description))
		} else {
			metaSb.WriteString(fmt.Sprintf("%s", description))
		}
	}

	if metaSb.String() != "" {
		var meta string
		if checkOpt(FormatMultiline) {
			meta = fmt.Sprintf("\n       %s\n", metaSb.String())
		} else {
			meta = fmt.Sprintf(" %s\n", metaSb.String())
		}
		todoSb.WriteString(meta)
	} else {
		todoSb.WriteString("\n")
	}

	return todoSb.String(), nil
}

// WriteFile encodes ical data and writes to file at Item.Path
func (i *Item) WriteFile() error {
	if i.Path == "" {
		return fmt.Errorf("Can not write Item without Path: %v", i)
	}

	// check and set topmost calendar object props
	requiredIcalProps := []string{
		ical.PropProductID,
		ical.PropVersion,
	}
	for _, p := range requiredIcalProps {
		prop := i.Ical.Props.Get(p)
		if prop == nil || prop.Value == "" {
			switch p {
			case ical.PropProductID:
				i.Ical.Props.SetText(ical.PropProductID, "-//KKGA.ME//NONSGML tdx//EN")
			case ical.PropVersion:
				i.Ical.Props.SetText(ical.PropVersion, "2.0")
			}
		}
	}

	vtodo, err := i.Vtodo()
	if err != nil {
		return err
	}

	// update modified date
	t := time.Now()
	vtodo.Props.SetDateTime(ical.PropLastModified, t)

	// update sequence
	seq := vtodo.Props.Get(ical.PropSequence)
	seqProp := ical.NewProp(ical.PropSequence)
	if seq == nil || seq.Value == "" {
		seqProp.Value = "0"
		vtodo.Props.Set(seqProp)
	} else {
		v, _ := seq.Int()
		nextSeq := fmt.Sprintf("%d", v+1)
		seqProp.Value = nextSeq
	}
	vtodo.Props.Set(seqProp)

	// check and set required vtodo props
	requiredVtodoProps := []string{
		ical.PropUID,
		ical.PropCreated,
		ical.PropDateTimeStamp,
	}
	for _, p := range requiredVtodoProps {
		prop := vtodo.Props.Get(p)
		if prop == nil || prop.Value == "" {
			switch p {
			case ical.PropCreated:
				vtodo.Props.SetDateTime(ical.PropCreated, t)
			case ical.PropDateTimeStamp:
				vtodo.Props.SetDateTime(ical.PropDateTimeStamp, t)
			case ical.PropUID:
				uid := GenerateUID()
				vtodo.Props.SetText(ical.PropUID, uid)
			}
		}
	}

	isStatusString := func(s string) bool {
		switch ToDoStatus(s) {
		case StatusNeedsAction, StatusInProcess, StatusCompleted, StatusCancelled:
			return true
		default:
			return false
		}
	}

	// check status
	st := vtodo.Props.Get(ical.PropStatus)
	if st == nil {
		stProp := ical.NewProp(ical.PropStatus)
		stProp.Value = string(StatusNeedsAction)
		vtodo.Props.Set(stProp)
	} else if !isStatusString(st.Value) {
		st.Value = string(StatusNeedsAction)
	}

	var buf bytes.Buffer
	err = ical.NewEncoder(&buf).Encode(i.Ical)
	if err != nil {
		return err
	}

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

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Tags returns a slice of hashtag strings parsed from summary and description
func (i *Item) Tags() (tags []Tag, err error) {
	re := regexp.MustCompile(HashtagRe)

	vt, err := i.Vtodo()
	if err != nil {
		return
	}
	summary, err := vt.Props.Text(ical.PropSummary)
	if err != nil {
		return
	}
	description, err := vt.Props.Text(ical.PropDescription)
	if err != nil {
		return
	}

	st := re.FindAllString(summary, -1)
	dt := re.FindAllString(description, -1)

	tagExists := func(tags []Tag, t Tag) bool {
		for _, tag := range tags {
			if tag == t {
				return true
			}
		}
		return false
	}

	for _, t := range st {
		if !tagExists(tags, Tag(t)) {
			tags = append(tags, Tag(t))
		}
	}
	for _, t := range dt {
		if !tagExists(tags, Tag(t)) {
			tags = append(tags, Tag(t))
		}
	}
	return
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
