package cmd

import (
	"flag"
)

func NewAddCmd() *AddCmd {
	c := &AddCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("add", flag.ExitOnError),
		alias:     []string{"a"},
		shortDesc: "Add todo",
		usageLine: "[options]",
	}}
	return c
}

type AddCmd struct {
	Cmd
}

func (c *AddCmd) Run() error {

	// if len(os.Args) > 1 && os.Args[1] == "new" {
	// 	todo := &todo.ToDo{}
	// 	todoBuf, _ := todo.Encode()

	// 	f, _ := os.Create(fmt.Sprintf("%s/ctdo-testing-%d.ics", calDir, time.Now().UnixNano()))
	// 	defer f.Close()
	// 	w := bufio.NewWriter(f)
	// 	_, _ = w.Write(todoBuf.Bytes())
	// 	w.Flush()
	// } else {
	// 	files, err := ioutil.ReadDir(calDir)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	todos := decodeFiles(files)
	// 	list := todo.NewList()
	// 	list.Init("new list", todos)
	// 	fmt.Println(list.String())
	// }
	return nil
}
