package storage

import (
	"github.com/mpuzanov/sysmonitor/internal/interfaces"
)

// NewStorage create storage for calendar
func NewStorage() interfaces.Storage {

	db := NewSystemStore()

	return db
}
