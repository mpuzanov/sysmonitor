package memory

import (
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// SaveLoadSystem Сохраняем текущую статистику системы
func (s *Store) SaveLoadSystem(data *model.LoadSystem) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.dbSys = append(s.dbSys, *data)
	return nil
}

// GetAvgLoadSystem Возврат среднего значения загрузки системы за period
func (s *Store) GetAvgLoadSystem(period int32) (*model.LoadSystem, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return avgLoadSystem(s.dbSys, period)
}

//avgLoadSystem получить среднее значение показателей за период
func avgLoadSystem(s []model.LoadSystem, period int32) (*model.LoadSystem, error) {
	res := model.LoadSystem{}
	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))
	//fmt.Println("now:", now, "timeStart", timeStart)

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
		res.QueryTime = now
		res.SystemLoadValue = RoundToFixed(res.SystemLoadValue/float64(count), 2)
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
			res.QueryTime = now
		}
	}
	return &res, nil
}
