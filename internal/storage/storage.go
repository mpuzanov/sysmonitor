package storage

import (
	"github.com/mpuzanov/sysmonitor/internal/storage/memory"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Storage интерфейс для работы с DB
type Storage interface {
	LoadSystem
	LoadCPU
	LoadDisk
	TalkersNet
	NetworkStatistics
}

// NewMemoryStorage create storage for sysmonitor
func NewMemoryStorage() Storage {
	return memory.NewSystemStore()
}

// LoadSystem интерфейс для работы с LoadSystem
type LoadSystem interface {
	//SaveLoadSystem сохранение информации
	SaveLoadSystem(data *model.LoadSystem) error
	//GetAvgLoadSystem получение средней информации за период
	GetAvgLoadSystem(period int32) (model.LoadSystem, error)
}

// LoadCPU интерфейс для работы с CPU
type LoadCPU interface {
	//SaveLoadCPU сохранение информации по CPU
	SaveLoadCPU(data *model.LoadCPU) error
	//GetAvgLoadCPU получение средней информации за период
	GetAvgLoadCPU(period int32) (model.LoadCPU, error)
}

// LoadDisk интерфейс для работы с дисками
type LoadDisk interface {
	//SaveLoadDisk сохранение информации по дискам
	SaveLoadDisk(data *model.LoadDisk)
	//GetInfoDisk получение информации по дискам
	GetInfoDisk() model.LoadDisk
}

// TalkersNet интерфейс для работы с top talkers сети
type TalkersNet interface {
	//SaveTalkersNet сохранение информации
	SaveTalkersNet(data *model.TalkersNet) error
	//GetAvgTalkersNet получение средней информации за период
	GetAvgTalkersNet(period int32) (model.TalkersNet, error)
}

// NetworkStatistics интерфейс для работы со статистикой по сетевым соединениям
type NetworkStatistics interface {
	//SaveTalkersNet сохранение информации
	SaveNetworkStatistics(data *model.NetworkStatistics) error
	//GetAvgTalkersNet получение средней информации за период
	GetAvgNetworkStatistics(period int32) (model.NetworkStatistics, error)
}
