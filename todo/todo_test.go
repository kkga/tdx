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
		{
			"test_data/e0a2e171-8314-4a36-8236-6851b05ff95b.ics",
			ToDo{
				Status:  ToDoStatusNeedsAction,
				Summary: "update sway timer to use timedown instead of swayidle #Later",
				UID:     "1a0a65abf057464db48d0a9183adb5db@void",
			},
		},
		{
			"test_data/48949a4f8b8149ec97f5a2715ac7c86b@void.ics",
			ToDo{
				Status:  ToDoStatusNeedsAction,
				Summary: "https://dave.cheney.net/practical-go/",
				UID:     "48949a4f8b8149ec97f5a2715ac7c86b@void",
				Due:     time.Date(2021, time.September, 21, 18, 06, 24, 0, time.UTC),
			},
		},
		{
			"test_data/20070514T103211Z-123404@example.com.ics",
			ToDo{
				Status:   ToDoStatusCompleted,
				Priority: ToDoPriorityHigh,
				Summary:  "Submit Revised Internet-Draft",
				UID:      "20070514T103211Z-123404@example.com",
				Due:      time.Date(2007, time.July, 9, 13, 0, 0, 0, time.UTC),
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
