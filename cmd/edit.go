package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type editOptions struct{}

const (
	layoutDateTime = "2 Jan 2006 15:04"
	layoutDate     = "2 Jan 2006"
)

func NewEditCmd() *cobra.Command {
	_ = &editOptions{}

	cmd := &cobra.Command{
		Use:     "edit <id>",
		Aliases: []string{"e"},
		Short:   "Edit todo",
		Long:    "Edit todo content in external program.",
		Args:    cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ tdx edit 1`),
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			IDs, err := stringsToInts(args)
			if err != nil {
				return err
			}

			var item *vdir.Item
			for _, id := range IDs {
				i, err := vd.ItemById(id)
				if err != nil {
					return err
				}
				item = i
			}

			return runEdit(item)
		},
	}

	return cmd
}

func runEdit(item *vdir.Item) error {
	tmp, err := os.CreateTemp("", "tdx")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	editorBin, err := exec.LookPath(editor)
	if err != nil {
		return errors.New("Set the VISUAL env variable to edit todos")
	}

	vtodo, err := item.Vtodo()
	if err != nil {
		return err
	}

	template := generateTemplate(*vtodo)
	_, err = tmp.Write([]byte(template))
	if err != nil {
		return err
	}

	cmd := exec.Command(editorBin, tmp.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: true}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	newProps, err := parseTemplate(tmp)
	if err != nil {
		return err
	}

	for p, newVal := range newProps {
		curP := vtodo.Props.Get(p)
		newP := ical.NewProp(p)

		if curP == nil && newVal == "" {
			continue
		}

		if curP != nil && curP.Value == "" && newVal == "" {
			continue
		}

		if curP != nil && curP.Value != "" && newVal == "" {
			vtodo.Props.Del(p)
			continue
		}

		switch p {
		case ical.PropDue:
			if t, _ := time.Parse(layoutDateTime, newVal); !t.IsZero() {
				newP.SetDateTime(t)
			} else if t, _ := time.Parse(layoutDate, newVal); !t.IsZero() {
				newP.SetDateTime(t)
			} else if due, _, err := parseDate(newVal); err == nil {
				newP.SetDateTime(due)
			}
		case ical.PropPriority:
			prioMap := map[string]vdir.ToDoPriority{
				"!!!": vdir.PriorityHigh,
				"!!":  vdir.PriorityMedium,
				"!":   vdir.PriorityLow,
			}
			p := prioMap[newVal]
			newP.Value = fmt.Sprint(p)
		case ical.PropStatus:
			statusMap := map[string]vdir.ToDoStatus{
				"":    vdir.StatusNeedsAction,
				"[ ]": vdir.StatusNeedsAction,
				"[-]": vdir.StatusCancelled,
				"[x]": vdir.StatusCompleted,
				"[X]": vdir.StatusCompleted,
			}
			if s, exists := statusMap[newVal]; exists {
				newP.Value = fmt.Sprint(s)
			} else {
				newP.Value = newVal
			}
		default:
			newP.Value = newVal
		}

		vtodo.Props.Set(newP)
	}

	if err := item.WriteFile(); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	f, err := item.Format(vdir.FormatDescription, vdir.FormatMultiline)
	if err != nil {
		return err
	}
	fmt.Print(f)

	return nil
}

func parseTemplate(f *os.File) (map[string]string, error) {
	props := make(map[string]string)

	_, err := f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" || strings.HasPrefix(s.Text(), "-") {
			break
		}
		s := strings.SplitN(s.Text(), ":", 2)
		k := strings.Trim(s[0], " ")
		v := strings.Trim(s[1], " ")
		props[k] = v
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return props, nil
}

func generateTemplate(vtodo ical.Component) string {
	sb := strings.Builder{}

	props := map[string]string{
		ical.PropSummary:     "",
		ical.PropDescription: "",
		ical.PropStatus:      "",
		ical.PropPriority:    "",
		ical.PropDue:         "",
		ical.PropLocation:    "",
	}

	prioMap := map[vdir.ToDoPriority]string{
		vdir.PriorityHigh:   "!!!",
		vdir.PriorityMedium: "!!",
		vdir.PriorityLow:    "!",
	}

	statusMap := map[vdir.ToDoStatus]string{
		vdir.StatusNeedsAction: "[ ]",
		vdir.StatusCancelled:   "[-]",
		vdir.StatusCompleted:   "[x]",
	}

	for name, prop := range vtodo.Props {
		p := prop[0]
		if _, exists := props[name]; exists {
			var val string
			switch name {
			case ical.PropDue:
				date, _ := p.DateTime(time.Local)
				if !date.IsZero() {
					val = date.Format(layoutDateTime)
				}
			case ical.PropPriority:
				curPrio, _ := strconv.Atoi(strings.Trim(p.Value, " "))
				val = prioMap[vdir.ToDoPriority(curPrio)]
			case ical.PropStatus:
				if v, exists := statusMap[vdir.ToDoStatus(strings.Trim(p.Value, " "))]; exists {
					val = v
				} else {
					val = p.Value
				}
			default:
				val = p.Value
			}
			props[p.Name] = val
		}
	}

	order := []string{
		ical.PropSummary,
		ical.PropDescription,
		ical.PropStatus,
		ical.PropPriority,
		ical.PropDue,
		ical.PropLocation,
	}

	for _, p := range order {
		sb.WriteString(fmt.Sprintf("%s: %s\n", p, props[p]))
	}

	sb.WriteString(templateHelp)

	return sb.String()
}

const templateHelp = `
--------------------

Edit todo fields above. Here's a cheatsheet.

DUE accepts following formats:
- 2 Jan 2006 15:04
- 2 Jan 2006
- natural date; same as <add>: see 'tdx add -h'

STATUS:
- [ ]
- [x]
- [-] (cancelled)

PRIORITY:
- !!!
- !!
- !
- [empty]
`
