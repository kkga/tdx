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
			collections, err := vd.Collections()
			if err != nil {
				t.Fatal(err)
			}
			for col := range collections {
				for _, item := range collections[col] {
					fmt.Printf("%s: %d,  %+v\n", col, item.id, item.ical.Children[0].Props.Get(ical.PropSummary))
				}
			}
		})
	}
}
