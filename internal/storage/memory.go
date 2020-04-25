package storage

import (
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Store структура хранения информации по системе
type Store struct {
	dbSys []model.LoadSystem
	dbCPU []model.LoadCPU
}

// NewSystemStore Возвращаем новое хранилище
func NewSystemStore() *Store {
	return &Store{dbSys: make([]model.LoadSystem, 0),
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
	res := model.LoadSystem{}
	// TODO получить среднее значение показателей за период
	// пока берём последнее значение
	if len(s.dbSys) > 0 {
		res = s.dbSys[len(s.dbSys)-1]
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
