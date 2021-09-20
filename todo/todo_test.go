package todo

import (
	"io"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/emersion/go-ical"
	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	var tests = []struct {
		filepath string
		want     ToDo
	}{
		{
			"test_data/20070313T123432Z-456553@example.com.ics",
			ToDo{
				Status:  ToDoStatusNeedsAction,
				Summary: "Submit Quebec Income Tax Return for 2006",
				UID:     "20070313T123432Z-456553@example.com",
				Due:     time.Date(2007, time.May, 01, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			dir, _ := os.Getwd()
			f, _ := os.Open(path.Join(dir, tt.filepath))
			defer f.Close()

			var vtodo ical.Component
			dec := ical.NewDecoder(f)

			for {
				cal, err := dec.Decode()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
				for _, comp := range cal.Children {
					if comp.Name == ical.CompToDo {
						vtodo = *comp
					}
				}

			}

			got := NewToDo()
			got.Init(vtodo)
			if diff := cmp.Diff(tt.want, *got); diff != "" {
				t.Errorf("Init() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
