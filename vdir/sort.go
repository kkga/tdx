package vdir

import (
	"time"

	"github.com/emersion/go-ical"
)

type ByPriority []*Item
type ByDue []*Item

func (p ByPriority) Len() int      { return len(p) }
func (p ByPriority) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ByPriority) Less(i, j int) bool {
	vt1, _ := p[i].Vtodo()
	vt2, _ := p[j].Vtodo()
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

func (d ByDue) Len() int      { return len(d) }
func (d ByDue) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (d ByDue) Less(i, j int) bool {
	vt1, _ := d[i].Vtodo()
	vt2, _ := d[j].Vtodo()
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
