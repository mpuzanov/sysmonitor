package repository

import (
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Storage интерфейс для работы с DB
type Storage interface {
	StorageLoadSystem
	StorageLoadCPU
	StorageLoadDisk
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