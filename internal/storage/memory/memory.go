package memory

import (
	"math"
	"sync"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Store структура хранения информации по системе
type Store struct {
	m            *sync.RWMutex
	dbSys        []model.LoadSystem
	dbCPU        []model.LoadCPU
	dbDisk       model.LoadDisk
	dbTalkersNet []model.TalkersNet
	dbNetStat    []model.NetworkStatistics
}

// NewSystemStore Возвращаем новое хранилище
func NewSystemStore() *Store {
	return &Store{
		m:            &sync.RWMutex{},
		dbSys:        make([]model.LoadSystem, 0),
		dbCPU:        make([]model.LoadCPU, 0),
		dbDisk:       model.LoadDisk{},
		dbTalkersNet: make([]model.TalkersNet, 0),
		dbNetStat:    make([]model.NetworkStatistics, 0),
	}
}

// RoundToFixed округление с заданным значением точности
func RoundToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}
