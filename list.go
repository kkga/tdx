package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
)

type List struct {
	Name  string
	ToDos []ToDo
}

func NewList() *List {
	return &List{}
}

func (l *List) Init(name string, todos []ical.Component) error {
	l.ToDos = make([]ToDo, 0, len(todos))

	for _, todo := range todos {
		if todo.Name != ical.CompToDo {
			return errors.New(fmt.Sprintf("Not VTODO component: %v", todo))
		}
		t := NewToDo()

		err := t.Init(todo)
		if err != nil {
			return err
		}

		l.ToDos = append(l.ToDos, *t)
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
