package vdir

import (
	"fmt"
	"testing"
)

func TestCollections(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"/home/kkga/.local/share/calendars/migadu/"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vd, err := NewVdirRoot(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			collections, _ := vd.Collections()
			fmt.Printf("%v", collections)
			// for _, d := range collections {
			// 	fmt.Println(d.Name())
			// }
		})
	}
}

func TestItems(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"/home/kkga/.local/share/calendars/migadu/"},
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
				fmt.Println(c.Name)
				fmt.Printf("%+v\n", items)
				fmt.Println("-------")
			}
		})
	}
}
