package vdir

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/emersion/go-ical"
)

func TestInit(t *testing.T) {
	cwd, _ := os.Getwd()
	vdpath := path.Join(cwd, "testdata/vdir/corrupted/")
	var tests = []struct {
		path string
	}{
		{vdpath},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vd := Vdir{}
			if err := vd.Init(tt.path); err != nil {
				t.Fatal(err)
			}
			for col, items := range vd {
				for _, item := range items {
					summary := item.Ical.Children[0].Props.Get(ical.PropSummary)
					fmt.Printf("%s: %+s\n", col, summary)
				}
			}
		})
	}
}

func TestUpdateDB(t *testing.T) {
	t.Run("", func(t *testing.T) {
		err := updateDB()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestViewDB(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, err := viewDB()
		if err != nil {
			t.Fatal(err)
		}
	})
}
