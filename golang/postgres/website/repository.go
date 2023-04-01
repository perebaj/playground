package website

import (
	"context"
	"errors"
)

var (
	ErrDuplicate    = errors.New("Record already exists")
	ErrNotExist     = errors.New("Record does not exist")
	ErrUpdateFailed = errors.New("Update failed")
	ErrDeleteFailed = errors.New("Delete failed")
)

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, website Website) (*Website, error)
	All(ctx context.Context) ([]Website, error)
	GetByName(ctx context.Context, name string) (*Website, error)
	Update(ctx context.Context, id int64, update Website) (*Website, error)
	Delete(ctx context.Context, id int64) error
}
