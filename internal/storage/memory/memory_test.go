package memory

import (
	"testing"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestAvgLoadSystem(t *testing.T) {
	testCases := []struct {
		desc   string
		sl     []model.LoadSystem
		period int32
		want   float64
	}{
		{
			desc: "test 1",
			sl: []model.LoadSystem{
				{QueryTime: time.Now().Add(-time.Second * time.Duration(6)),
					SystemLoadValue: 1.4},
				{QueryTime: time.Now(), SystemLoadValue: 1.2},
			},
			period: 15,
			want:   1.3,
		},
		{
			desc: "test 2",
			sl: []model.LoadSystem{
				{QueryTime: time.Now(), SystemLoadValue: 1.2},
			},
			period: 15,
			want:   1.2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := avgLoadSystem(tC.sl, tC.period)
			assert.Empty(t, err)
			assert.Equal(t, tC.want, got.SystemLoadValue)
		})
	}
}

func TestAvgLoadCPU(t *testing.T) {
	testCases := []struct {
		desc   string
		sl     []model.LoadCPU
		period int32
		want   model.LoadCPU
	}{
		{
			desc: "test 1",
			sl: []model.LoadCPU{
				{QueryTime: time.Now().Add(-time.Second * time.Duration(6)),
					UserMode: 1.3, SystemMode: 0.6, Idle: 96.7},
				{QueryTime: time.Now(), UserMode: 8.5, SystemMode: 3.0, Idle: 87.0},
			},
			period: 15,
			want:   model.LoadCPU{UserMode: 4.9, SystemMode: 1.8, Idle: 91.85},
		},
		{
			desc: "test 2",
			sl: []model.LoadCPU{
				{QueryTime: time.Now().Add(-time.Second * time.Duration(20)),
					UserMode: 1.3, SystemMode: 0.6, Idle: 96.7},
				{QueryTime: time.Now().Add(-time.Second * time.Duration(10)),
					UserMode: 8.7, SystemMode: 3.1, Idle: 86.7},
				{QueryTime: time.Now(), UserMode: 8.5, SystemMode: 3.0, Idle: 87.0},
			},
			period: 15,
			want:   model.LoadCPU{UserMode: 8.6, SystemMode: 3.05, Idle: 86.85},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := avgLoadCPU(tC.sl, tC.period)
			assert.Empty(t, err)
			assert.Equal(t, tC.want.UserMode, got.UserMode)
			assert.Equal(t, tC.want.SystemMode, got.SystemMode)
			assert.Equal(t, tC.want.Idle, got.Idle)
		})
	}
}
