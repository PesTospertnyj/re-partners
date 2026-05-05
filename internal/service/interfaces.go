package service

import (
	"context"
	"re-partners/internal/dto"
)

type PackService interface {
	GetPackSizes(context.Context) ([]dto.PackSize, error)
	AddPackSize(context.Context, dto.PackSize) error
	DeletePackSize(context.Context, int) error
}
