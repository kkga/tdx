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

	"github.com/kelseyhightower/envconfig"
	"github.com/kkga/tdx/vdir"
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

	vdir       vdir.Vdir
	collection *vdir.Collection

	listFlag     string
	listRequired bool

	conf Config
}

type Config struct {
	Path          string `required:"true"`
	DefaultList   string `split_words:"true"`
	DefaultStatus string `split_words:"true" default:"NEEDS-ACTION"`
	DefaultSort   string `split_words:"true" default:"PRIO"`
	Color         bool   `default:"true"`
}

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	var Conf Config
	err := envconfig.Process("TDX", &Conf)
	if err != nil {
		return err
	}
	c.conf = Conf

	c.listFlag = c.conf.DefaultList
	c.fs.Usage = c.usage

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	c.vdir = vdir.Vdir{}
	if err := c.vdir.Init(c.conf.Path); err != nil {
		return err
	}

	if c.listRequired && c.listFlag == "" {
		return errors.New("List flag required. See 'tdx <command> -h'")
	} else if c.listFlag != "" {
		names := []string{}
		for col := range c.vdir {
			names = append(names, col.Name)
			if col.Name == c.listFlag {
				c.collection = col
			}
		}
		if c.collection == nil {
			return fmt.Errorf("List does not exist: %q\nAvailable lists: %s", c.listFlag, strings.Join(names, ", "))
		}
	}

	return nil
}

func (c *Cmd) usage() {
	fmt.Println(c.short)
	fmt.Println()

	fmt.Println("USAGE")
	fmt.Printf("  tdx %s %s\n\n", c.fs.Name(), c.usageLine)

	if c.long != "" {
		fmt.Println(c.long)
		fmt.Println()
	}

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
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
