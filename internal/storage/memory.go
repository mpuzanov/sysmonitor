package storage

import (
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
	//log.Println("now:", now, "timeStart", timeStart)

	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if timeStart.Before(s[i].QueryTime) {
			res.SystemLoadValue += s[i].SystemLoadValue
			count++
		} else {
			break
		}
	}
	if count > 1 {
		res.SystemLoadValue = utils.RandToFixed(res.SystemLoadValue/float64(count), 2)
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
	return avgLoadCPU(s.dbCPU, period)
}

//avgLoadCPU получить среднее значение показателей за период
func avgLoadCPU(s []model.LoadCPU, period int32) (*model.LoadCPU, error) {
	res := model.LoadCPU{}
	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))
	//log.Println("now:", now, "timeStart", timeStart)

	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if timeStart.Before(s[i].QueryTime) {
			res.UserMode += s[i].UserMode
			res.SystemMode += s[i].SystemMode
			res.Idle += s[i].Idle
			count++
		} else {
			break
		}
	}
	if count > 1 {
		res.UserMode = utils.RandToFixed(res.UserMode/float64(count), 2)
		res.SystemMode = utils.RandToFixed(res.SystemMode/float64(count), 2)
		res.Idle = utils.RandToFixed(res.Idle/float64(count), 2)
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
		}
	}
	return &res, nil
}
