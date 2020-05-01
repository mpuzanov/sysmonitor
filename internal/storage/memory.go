package storage

import (
	"log"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/utils"
)

// Store структура хранения информации по системе
type Store struct {
	dbSys []model.LoadSystem
	dbCPU []model.LoadCPU
}

// NewSystemStore Возвращаем новое хранилище
func NewSystemStore() *Store {
	return &Store{
		dbSys: make([]model.LoadSystem, 0),
		dbCPU: make([]model.LoadCPU, 0),
	}
}

// SaveLoadSystem Вводим текущее показание системы
func (s *Store) SaveLoadSystem(data *model.LoadSystem) error {
	s.dbSys = append(s.dbSys, *data)
	return nil
}

// SaveLoadCPU Вводим текущее показание CPU
func (s *Store) SaveLoadCPU(data *model.LoadCPU) error {
	s.dbCPU = append(s.dbCPU, *data)
	return nil
}

// GetAvgLoadSystem Возврат среднего значения загрузки системы за period
func (s *Store) GetAvgLoadSystem(period int32) (*model.LoadSystem, error) {
	return avgLoadSystem(s.dbSys, period)
}

//avgLoadSystem получить среднее значение показателей за период
func avgLoadSystem(s []model.LoadSystem, period int32) (*model.LoadSystem, error) {
	res := model.LoadSystem{}
	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))
	log.Println("now      ", now)
	log.Println("timeStart", timeStart)

	slv := 0.0
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		log.Println(s[i])
		if timeStart.Before(s[i].QueryTime) {
			slv += s[i].SystemLoadValue
			count++
		}
	}
	if count > 1 {
		res.SystemLoadValue = utils.RandToFixed(slv/float64(count), 2)
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
		}
	}

	return &res, nil
}

// GetAvgLoadCPU Возврат среднего значения загрузки системы за period
func (s *Store) GetAvgLoadCPU(period int32) (*model.LoadCPU, error) {
	res := model.LoadCPU{}
	// TODO получить среднее значение показателей за период

	// пока берём последнее значение
	if len(s.dbCPU) > 0 {
		res = s.dbCPU[len(s.dbCPU)-1]
	}
	return &res, nil
}
