package storage

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
