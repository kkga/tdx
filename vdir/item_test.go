package vdir

import (
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTags(t *testing.T) {
	cwd, _ := os.Getwd()
	vdpath := path.Join(cwd, "testdata/vdir/with_tags/")
	var tests = []struct {
		path string
		want []Tag
	}{
		{vdpath, []Tag{"#Quebec", "#go", "#sway", "#Later"}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			vd := Vdir{}
			if err := vd.Init(tt.path); err != nil {
				t.Fatal(err)
			}
			got := []Tag{}
			for _, items := range vd {
				for _, item := range items {
					tags, err := item.Tags()
					if err != nil {
						t.Fatal(err)
					}
					got = append(got, tags...)
				}
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Tags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
