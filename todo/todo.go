package todo

import (
	"time"

	"github.com/emersion/go-ical"
)

type ToDo struct {
	UID         string
	Status      ToDoStatus
	Priority    ToDoPriority
	Summary     string
	Description string
	Tags        []string
	Due         time.Time
	OtherProps  ical.Props
}

type (
	ToDoStatus   string
	ToDoPriority int
	ToDoUID      string
)

// func NewToDo() *ToDo {
// 	return &ToDo{}
// }

// TODO: probably don't need this, use the ical.Component directly
// func (t *ToDo) Init(todo ical.Component) error {
// 	props := todo.Props
// 	for p := range props {
// 		switch p {
// 		case ical.PropUID:
// 			uid, err := todo.Props.Get(ical.PropUID).Text()
// 			if err != nil {
// 				return err
// 			}
// 			t.UID = uid
// 			delete(props, p)
// 		case ical.PropStatus:
// 			s, err := todo.Props.Get(ical.PropStatus).Text()
// 			if err != nil {
// 				return err
// 			}
// 			t.Status = ToDoStatus(s)
// 			delete(props, p)
// 		case ical.PropSummary:
// 			s, err := todo.Props.Get(ical.PropSummary).Text()
// 			if err != nil {
// 				return err
// 			}
// 			t.Summary = s
// 			delete(props, p)
// 		case ical.PropDescription:
// 			s, err := todo.Props.Get(ical.PropDescription).Text()
// 			if err != nil {
// 				return err
// 			}
// 			t.Description = s
// 			delete(props, p)
// 		case ical.PropDue:
// 			time, err := todo.Props.Get(ical.PropDue).DateTime(t.Due.Location())
// 			if err != nil {
// 				return err
// 			}
// 			t.Due = time
// 			delete(props, p)
// 		case ical.PropPriority:
// 			prio, err := todo.Props.Get(ical.PropPriority).Int()
// 			if err != nil {
// 				return err
// 			}
// 			t.Priority = ToDoPriority(prio)
// 			delete(props, p)
// 		}
// 	}
// 	t.OtherProps = props
// 	return nil
// }
