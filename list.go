package main

import "strings"

type List struct {
	Name  string
	ToDos map[ToDoUID]ToDo
}

func NewList() *List {
	return &List{}
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
