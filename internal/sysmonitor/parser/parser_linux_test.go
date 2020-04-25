//+build linux

package parser_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
)

func TestParserSystemInfo(t *testing.T) {

	testCases := []struct {
		desc string
		in   string
		want float64
	}{
		{
			desc: "",
			in:   "top - 19:30:09 up  4:30,  1 user,  load average: 1,02, 0,95, 0,80",
			want: 1.02,
		},
		{
			desc: "",
			in:   "top - 19:30:09 up  4:30,  1 user,  load average: 1.02, 0.95, 0.80",
			want: 1.02,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.ParserSystemInfo(tC.in)
			if err != nil {
				t.Errorf("%s, got=%v, expected=%v, error: %v", tC.desc, got, tC.want, err)
			}
		})
	}
}
