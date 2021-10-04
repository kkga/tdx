// The cmd package implements a command-line interface.
package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kkga/tdx/vdir"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/ru"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
	Alias() []string
}

type Cmd struct {
	fs        *flag.FlagSet
	alias     []string
	short     string
	long      string
	usageLine string

	conf Config
	args []string

	vdir vdir.Vdir
}

type Config struct {
	Path     string `required:"true"`    // Path to vdir
	ListOpts string `split_words:"true"` // Default options for list command
	AddOpts  string `split_words:"true"` // Default options for add command
}

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	var conf Config
	err := envconfig.Process("TDX", &conf)
	if err != nil {
		return err
	}

	c.conf = conf
	c.args = args
	c.fs.Usage = c.usage

	c.vdir = vdir.Vdir{}
	if err := c.vdir.Init(c.conf.Path); err != nil {
		return err
	}

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) checkListFlag(list string, required bool, cmd Runner) error {
	if list == "" && required {
		return fmt.Errorf("List flag required. See 'tdx %s -h'", cmd.Name())
	} else if list != "" {
		names := []string{}
		for col := range c.vdir {
			names = append(names, col.Name)
			if col.Name == list {
				return nil
			}
		}
		return fmt.Errorf("List does not exist: %q\nAvailable lists: %s", list, strings.Join(names, ", "))
	} else {
		return nil
	}
}

func (c *Cmd) usage() {
	fmt.Println(c.short)
	fmt.Println()

	fmt.Println("USAGE")
	fmt.Printf("  tdx %s %s\n\n", c.fs.Name(), c.usageLine)

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
	}

	if c.long != "" {
		fmt.Println()
		fmt.Println(c.long)
	}

}

func (c *Cmd) argsToIDs() (IDs []int, err error) {
	if len(c.fs.Args()) == 0 {
		return IDs, errors.New("Specify one or multiple IDs")
	}

	for _, s := range c.fs.Args() {
		id, err := strconv.Atoi(s)
		if err != nil {
			return IDs, fmt.Errorf("Invalid todo ID: %q", s)
		}
		IDs = append(IDs, id)
	}
	return
}

func promptConfirm(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}
	r := bufio.NewReader(os.Stdin)

	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)

		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func parseDate(s string) (t time.Time, text string, err error) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(ru.All...)
	w.Add(common.All...)

	now := time.Now()

	r, err := w.Parse(s, now)
	if err != nil {
		return t, text, err
	}
	if r == nil {
		return t, text, errors.New("No date found")
	}

	// strip clock from time if it's the same as now (i.e. not specified)
	rH, rM, rS := r.Time.Clock()
	nH, nM, nS := now.Clock()
	if time.Date(0, 0, 0, rH, rM, rS, 0, time.Local).Equal(time.Date(0, 0, 0, nH, nM, nS, 0, time.Local)) {
		y, m, d := r.Time.Date()
		t = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	} else {
		t = r.Time
	}

	text = r.Text

	fmt.Println(
		"found time:",
		t.Format("2 Jan 2006 15:04:05"),
		"mentioned in:",
		s[r.Index:r.Index+len(r.Text)],
	)

	return
}
