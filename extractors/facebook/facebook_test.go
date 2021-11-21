package facebook

import (
	"testing"

	"github.com/sodapanda/annie/extractors/types"
	"github.com/sodapanda/annie/test"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:     "https://www.facebook.com/groups/314070194112/permalink/10155168902769113/",
				Title:   "Ukrainian Scientists Worldwide Public Group | Facebook",
				Size:    336975453,
				Quality: "hd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.args.URL, types.Options{})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
	}
}
