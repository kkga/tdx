package cmd

import (
	"flag"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List todos",
		usageLine: "[options]",
	}}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	return c
}

type ListCmd struct {
	Cmd
	json bool
}

func (c *ListCmd) Run() error {
	// decode todos into map

	// sb := strings.Builder{}

	// for _, t := range c.ToDos {
	// 	if t.Status != todo.ToDoStatusCompleted {
	// 		sb.WriteString(t.String())
	// 		sb.WriteString("\n")
	// 	}
	// }
	// fmt.Println(strings.TrimSpace(sb.String()))

	// if c.json {
	// 	fmt.Println("list -json called")
	// } else {
	// 	fmt.Println("list called")
	// }

	return nil
}
