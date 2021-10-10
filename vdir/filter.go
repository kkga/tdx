package vdir

import (
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

type filter interface {
	Items() []*Item
	Keep(item Item, i interface{}) (bool, error)
}

func (s ByStatus) Items() []*Item { return s }
func (s ByStatus) Keep(item Item, i interface{}) (bool, error) {
	status := i.(ToDoStatus)
	if status.String() == StatusAny.String() {
		return true, nil
	}

	vt, err := item.Vtodo()
	if err != nil {
		return false, err
	}

	st, err := vt.Props.Text(ical.PropStatus)
	if err != nil {
		return false, err
	}

	return status.String() == ToDoStatus(st).String(), nil
}

func (t ByTag) Items() []*Item { return t }
func (t ByTag) Keep(item Item, i interface{}) (bool, error) {
	tag := i.(Tag)
	if tag.String() == "" {
		return true, nil
	}

	tags, err := item.Tags()
	if err != nil {
		return false, err
	}
	for _, t := range tags {
		if tag.String() == t.String() {
			return true, nil
		}
	}
	return false, nil
}

func (d ByDue) Items() []*Item { return d }
func (d ByDue) Keep(item Item, i interface{}) (bool, error) {
	dueDays := i.(int)
	if dueDays == 0 {
		return true, nil
	}

	now := time.Now()
	inDueDays := now.AddDate(0, 0, dueDays)

	vt, err := item.Vtodo()
	if err != nil {
		return false, err
	}

	due, err := vt.Props.DateTime(ical.PropDue, time.Local)
	if err != nil {
		return false, err
	}

	if !due.IsZero() && due.Before(inDueDays) {
		return true, nil
	}
	return false, nil
}

func (d ByText) Items() []*Item { return d }
func (d ByText) Keep(item Item, i interface{}) (bool, error) {
	text := i.(string)
	if text == "" {
		return true, nil
	}

	vt, err := item.Vtodo()
	if err != nil {
		return false, err
	}

	summary, err := vt.Props.Text(ical.PropSummary)
	if err != nil {
		return false, err
	}

	if strings.Contains(strings.ToLower(summary), strings.ToLower(text)) {
		return true, nil
	}
	return false, nil
}

func Filter(f filter, i interface{}) (filtered []*Item, err error) {
	for _, item := range f.Items() {
		keep, err := f.Keep(*item, i)
		if err != nil {
			return filtered, err
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	return
}
