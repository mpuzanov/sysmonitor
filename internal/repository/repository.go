package repository

import (
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Storage интерфейс для работы с DB
type Storage interface {
	StorageLoadSystem
	StorageLoadCPU
	StorageLoadDisk
	StorageTalkersNet
	StorageNetworkStatistics
}

// StorageLoadSystem интерфейс для работы с LoadSystem
type StorageLoadSystem interface {
	//SaveLoadSystem сохранение информации
	SaveLoadSystem(data *model.LoadSystem) error
	//GetAvgLoadSystem получение средней информации за период
	GetAvgLoadSystem(period int32) (*model.LoadSystem, error)
}

// StorageLoadCPU интерфейс для работы с CPU
type StorageLoadCPU interface {
	//SaveLoadCPU сохранение информации по CPU
	SaveLoadCPU(data *model.LoadCPU) error
	//GetAvgLoadCPU получение средней информации за период
	GetAvgLoadCPU(period int32) (*model.LoadCPU, error)
}

// StorageLoadDisk интерфейс для работы с дисками
type StorageLoadDisk interface {
	//SaveLoadDisk сохранение информации по дискам
	SaveLoadDisk(data *model.LoadDisk)
	//GetInfoDisk получение информации по дискам
	GetInfoDisk() *model.LoadDisk
}

// StorageTalkersNet интерфейс для работы с top talkers сети
type StorageTalkersNet interface {
	//SaveTalkersNet сохранение информации
	SaveTalkersNet(data *model.TalkersNet) error
	//GetAvgTalkersNet получение средней информации за период
	GetAvgTalkersNet(period int32) (*model.TalkersNet, error)
}

// StorageNetworkStatistics интерфейс для работы со статистикой по сетевым соединениям
type StorageNetworkStatistics interface {
	//SaveTalkersNet сохранение информации
	SaveNetworkStatistics(data *model.NetworkStatistics) error
	//GetAvgTalkersNet получение средней информации за период
	GetAvgNetworkStatistics(period int32) (*model.NetworkStatistics, error)
}
