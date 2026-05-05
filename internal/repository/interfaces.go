package repository

import (
	"context"
	"re-partners/internal/model"
)

type Repository interface {
	GetPackSizes(context.Context) ([]model.PackSize, error)
	AddPackSize(context.Context, int) error
	DeletePackSize(context.Context, int) error
}
