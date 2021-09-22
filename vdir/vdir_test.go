package vdir

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/emersion/go-ical"
)

var (
	cwd, _       = os.Getwd()
	testVdirPath = path.Join(cwd, "test_data/vdir")
)

func TestCollections(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{testVdirPath},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vd, err := NewVdirRoot(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			collections, _ := vd.Collections()
			fmt.Printf("%v\n", collections)
		})
	}
}

func TestItems(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{testVdirPath},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vd, err := NewVdirRoot(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			collections, _ := vd.Collections()
			for _, c := range collections {
				items, err := c.Items()
				if err != nil {
					t.Fatal(err)
				}
				for id, item := range items {
					fmt.Printf("%d: %+v\n", id, item.Children[0].Props.Get(ical.PropSummary))
				}
			}
		})
	}
}
