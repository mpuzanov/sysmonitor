package memory

import (
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// SaveLoadCPU Сохраняем текущую статистику по CPU
func (s *Store) SaveLoadCPU(data *model.LoadCPU) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.dbCPU = append(s.dbCPU, *data)
	return nil
}

// GetAvgLoadCPU Возврат среднего значения загрузки системы за period
func (s *Store) GetAvgLoadCPU(period int32) (*model.LoadCPU, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return avgLoadCPU(s.dbCPU, period)
}

//avgLoadCPU получить среднее значение показателей за период
func avgLoadCPU(s []model.LoadCPU, period int32) (*model.LoadCPU, error) {
	res := model.LoadCPU{}
	now := time.Now().Local()
	timeStart := now.Add(-time.Second * time.Duration(period))
	//fmt.Println("now:", now, "timeStart", timeStart)

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
		res.QueryTime = now
		res.UserMode = RoundToFixed(res.UserMode/float64(count), 2)
		res.SystemMode = RoundToFixed(res.SystemMode/float64(count), 2)
		res.Idle = RoundToFixed(res.Idle/float64(count), 2)
	} else {
		// берём последнее значение
		if len(s) > 0 {
			res = s[len(s)-1]
			res.QueryTime = now
		}
	}
	return &res, nil
}
