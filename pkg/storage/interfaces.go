package storage

import (
	"context"
)

// источник данных
type IStorage interface {
	// закрытие базы
	Close(ctx context.Context) error
}
