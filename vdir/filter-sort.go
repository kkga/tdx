package vdir

import (
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

type (
	ByText         []*Item
	ByPriority     []*Item
	ByDue          []*Item
	ByStatus       []*Item
	ByCreated      []*Item
	ByTags         []*Item
	ByTagsExcluded []*Item
)

// Filter

type filter interface {
	Items() []*Item
	Keep(item Item, i interface{}) (bool, error)
}

func (f ByStatus) Items() []*Item { return f }
func (f ByStatus) Keep(item Item, i interface{}) (bool, error) {
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

func (f ByTags) Items() []*Item { return f }
func (f ByTags) Keep(item Item, i interface{}) (bool, error) {
	tags := i.([]Tag)
	if len(tags) == 0 {
		return true, nil
	}
	for _, tag := range tags {
		hasTag, err := item.HasTag(tag)
		if err != nil {
			return false, err
		}
		if hasTag {
			return true, nil
		}
	}
	return false, nil
}

func (f ByTagsExcluded) Items() []*Item { return f }
func (f ByTagsExcluded) Keep(item Item, i interface{}) (bool, error) {
	tags := i.([]Tag)
	if len(tags) == 0 {
		return true, nil
	}

	for _, tag := range tags {
		hasTag, err := item.HasTag(tag)
		if err != nil {
			return false, err
		}
		if hasTag {
			return false, nil
		}
	}
	return true, nil
}

func (f ByDue) Items() []*Item { return f }
func (f ByDue) Keep(item Item, i interface{}) (bool, error) {
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

func (f ByText) Items() []*Item { return f }
func (f ByText) Keep(item Item, i interface{}) (bool, error) {
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

// Sort

func (s ByPriority) Len() int      { return len(s) }
func (s ByPriority) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByPriority) Less(i, j int) bool {
	vt1, _ := s[i].Vtodo()
	vt2, _ := s[j].Vtodo()
	prio1 := vt1.Props.Get(ical.PropPriority)
	prio2 := vt2.Props.Get(ical.PropPriority)

	var prio1Val int
	var prio2Val int

	if prio1 != nil {
		prio1Val, _ = prio1.Int()
	}
	if prio2 != nil {
		prio2Val, _ = prio2.Int()
	}

	if prio1Val == 0 {
		return false
	} else if prio2Val == 0 {
		return true
	} else {
		return prio1Val < prio2Val
	}
}

func (s ByDue) Len() int      { return len(s) }
func (s ByDue) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByDue) Less(i, j int) bool {
	vt1, _ := s[i].Vtodo()
	vt2, _ := s[j].Vtodo()
	p1 := vt1.Props.Get(ical.PropDue)
	p2 := vt2.Props.Get(ical.PropDue)

	var v1 time.Time
	var v2 time.Time

	if p1 != nil {
		v1, _ = p1.DateTime(time.UTC)
	}
	if p2 != nil {
		v2, _ = p2.DateTime(time.UTC)
	}

	if v1.IsZero() {
		return false
	} else if v2.IsZero() {
		return true
	} else {
		return v1.Before(v2)
	}
}

func (s ByStatus) Len() int      { return len(s) }
func (s ByStatus) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByStatus) Less(i, j int) bool {
	vt1, _ := s[i].Vtodo()
	vt2, _ := s[j].Vtodo()
	p1 := vt1.Props.Get(ical.PropStatus)
	p2 := vt2.Props.Get(ical.PropStatus)

	var v1 string
	var v2 string

	if p1 != nil {
		v1, _ = p1.Text()
	}
	if p2 != nil {
		v2, _ = p2.Text()
	}

	if ToDoStatus(v1) == StatusCompleted || ToDoStatus(v1) == StatusCancelled {
		return false
	} else if ToDoStatus(v2) == StatusCompleted || ToDoStatus(v2) == StatusCancelled {
		return true
	} else {
		return false
	}
}

func (s ByCreated) Len() int      { return len(s) }
func (s ByCreated) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByCreated) Less(i, j int) bool {
	vt1, _ := s[i].Vtodo()
	vt2, _ := s[j].Vtodo()
	p1 := vt1.Props.Get(ical.PropCreated)
	p2 := vt2.Props.Get(ical.PropCreated)

	var v1 time.Time
	var v2 time.Time

	if p1 != nil {
		v1, _ = p1.DateTime(time.UTC)
	}
	if p2 != nil {
		v2, _ = p2.DateTime(time.UTC)
	}

	if v1.IsZero() {
		return true
	} else if v2.IsZero() {
		return false
	} else {
		return v1.After(v2)
	}
}
