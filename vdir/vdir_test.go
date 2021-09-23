package vdir

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/emersion/go-ical"
)

var (
	cwd, _ = os.Getwd()
	// testVdirPath = path.Join(cwd, "test_data/vdir")
	testVdirPath = path.Join("/home/kkga/.local/share/calendars")
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
			collections, err := vd.Collections()
			if err != nil {
				t.Fatal(err)
			}
			for col, items := range collections {
				for _, item := range items {
					summary := item.Ical.Children[0].Props.Get(ical.PropSummary)
					fmt.Printf("%s: %+s\n", col, summary)
				}
			}
		})
	}
}
