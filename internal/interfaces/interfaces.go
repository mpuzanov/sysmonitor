package interfaces

import (
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// Storage интерфейс для работы с DB
type Storage interface {
	StorageLoadSystem
	StorageLoadCPU
}

// StorageLoadSystem интерфейс для работы с LoadSystem
type StorageLoadSystem interface {
	SaveLoadSystem(data *model.LoadSystem) error
	GetAvgLoadSystem(period int32) (*model.LoadSystem, error)
}

// StorageLoadCPU интерфейс для работы с CPU
type StorageLoadCPU interface {
	SaveLoadCPU(data *model.LoadCPU) error
	GetAvgLoadCPU(period int32) (*model.LoadCPU, error)
}
