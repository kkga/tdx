// The cmd package implements a command-line interface.
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kkga/tdx/vdir"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/ru"
)

func checkList(vd vdir.Vdir, list string, required bool) error {
	if list == "" && required {
		return errors.New("List flag required. See 'tdx %s -h'")
	} else if list != "" {
		names := []string{}
		for col := range vd {
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

func stringsToInts(ss []string) (ints []int, err error) {
	// if len(ss) == 0 {
	// 	return ints, errors.New("Specify one or multiple IDs")
	// }

	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return ints, fmt.Errorf("Invalid arg: %q", s)
		}
		ints = append(ints, i)
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

func containsString(ss []string, s string) bool {
	for _, a := range ss {
		if a == s {
			return true
		}
	}
	return false
}
