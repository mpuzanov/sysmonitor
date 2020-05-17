//+build linux darwin

package parser_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/errors"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParserSystemLoad(t *testing.T) {

	testCases := []struct {
		desc string
		in   string
		want float64
	}{
		{
			desc: "test 1",
			in:   "top - 19:30:09 up  4:30,  1 user,  load average: 1,02, 0,95, 0,80",
			want: 1.02,
		},
		{
			desc: "test 2",
			in:   "top - 19:30:09 up  4:30,  1 user,  load average: 1.02, 0.95, 0.80",
			want: 1.02,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.SystemLoad(tC.in)
			assert.Empty(t, err)
			assert.Equal(t, tC.want, got.SystemLoadValue)
		})
	}
}

func TestParserLoadCPU(t *testing.T) {

	testCases := []struct {
		desc string
		in   string
		want model.LoadCPU
		err  error
	}{
		{
			desc: "test 1",
			in:   "%Cpu(s):  1,3 us,  0,6 sy,  0,0 ni, 96,7 id,  1,3 wa,  0,0 hi,  0,1 si,  0,0 st",
			want: model.LoadCPU{UserMode: 1.3, SystemMode: 0.6, Idle: 96.7},
			err:  nil,
		},
		{
			desc: "test 2",
			in:   "%Cpu(s):  1.3 us,  0.6 sy,  0.0 ni, 96.7 id,  1.3 wa,  0.0 hi,  0.1 si,  0.0 st",
			want: model.LoadCPU{UserMode: 1.3, SystemMode: 0.6, Idle: 96.7},
			err:  nil,
		},
		{
			desc: "test error",
			in:   "%Cpu(s):  1.3 us,  0.6 sy,  0.0 ni, ",
			want: model.LoadCPU{UserMode: 1.3, SystemMode: 0.6, Idle: 96.7},
			err:  errors.ErrParserReadInfoCPU,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.LoadCPU(tC.in)
			//assert.Empty(t, err)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.Equal(t, tC.want.UserMode, got.UserMode)
				assert.Equal(t, tC.want.SystemMode, got.SystemMode)
				assert.Equal(t, tC.want.Idle, got.Idle)
			}
		})
	}
}
