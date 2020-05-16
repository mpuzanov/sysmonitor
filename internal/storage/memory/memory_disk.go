package memory

import (
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// SaveLoadDisk Сохраняем текущую статистику по дискам
func (s *Store) SaveLoadDisk(data *model.LoadDisk) {
	s.m.Lock()
	defer s.m.Unlock()
	s.dbDisk = *data
}

// GetInfoDisk Возвращаем текущую статистику по дискам
func (s *Store) GetInfoDisk() model.LoadDisk {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.dbDisk
}
