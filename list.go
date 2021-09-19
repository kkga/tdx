package main

import (
	"strings"

	"github.com/emersion/go-ical"
)

type List struct {
	Name  string
	ToDos map[ToDoUID]ToDo
}

func NewList() *List {
	return &List{}
}

func (l *List) Init(name string, todos []ical.Component) error {
	for _, todo := range todos {
		t := NewToDo()

		err := t.ParseComponent(todo)
		if err != nil {
			return err
		}

		uid, err := todo.Props.Get(ical.PropUID).Text()
		if err != nil {
			return err
		}

		l.ToDos[ToDoUID(uid)] = *t
	}
	return nil
}

func (l *List) String() string {
	sb := strings.Builder{}

	for _, t := range l.ToDos {
		if t.Status != ToDoStatusCompleted {
			sb.WriteString(t.String())
			sb.WriteString("\n")
		}
	}
	return strings.TrimSpace(sb.String())
}
