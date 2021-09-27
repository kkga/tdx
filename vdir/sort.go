package vdir

import "github.com/emersion/go-ical"

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
