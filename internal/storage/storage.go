package storage

import (
	"github.com/mpuzanov/sysmonitor/internal/repository"
	"github.com/mpuzanov/sysmonitor/internal/storage/memory"
)

// NewStorage create storage for sysmonitor
func NewStorage() repository.Storage {
	return memory.NewSystemStore()
}
