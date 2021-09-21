package vdir

import (
	"fmt"
	"testing"
)

func TestCollections(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{"/home/kkga/.local/share/calendars/"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			collections, _ := Collections(tt.path)
			fmt.Printf("%+v", collections)
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
		{"/home/kkga/.local/share/calendars/"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			collections, _ := Collections(tt.path)
			// items, err := collections[1].Items()
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// fmt.Println(items)

			for _, c := range collections {
				items, err := c.Items()
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(c.Name)
				fmt.Printf("%v\n", items)
				fmt.Println("-------")
			}
		})
	}
}
