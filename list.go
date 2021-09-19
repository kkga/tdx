package main

type List struct {
	Name  string
	ToDos map[string]ToDo
}

func NewList() *List {
	return &List{}
}
