package service

import (
	"context"
	"go.uber.org/zap"
	"re-partners/internal/converter"
	"re-partners/internal/dto"
	"re-partners/internal/repository"
)

type packSizeService struct {
	packRepository repository.Repository
	log            *zap.SugaredLogger
}

func NewPackSizeService(packRepository repository.Repository, log *zap.SugaredLogger) PackService {
	return &packSizeService{
		packRepository: packRepository,
		log:            log,
	}
}

func (s *packSizeService) GetPackSizes(ctx context.Context) ([]dto.PackSize, error) {
	sizes, err := s.packRepository.GetPackSizes(ctx)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	packSizes := make([]dto.PackSize, len(sizes))
	for i, size := range sizes {
		packSizes[i] = converter.PackSizeToDTO(size)
	}

	return packSizes, nil
}

func (s *packSizeService) AddPackSize(ctx context.Context, packSize dto.PackSize) error {
	size := converter.PackSizeToModel(packSize)
	if err := s.packRepository.AddPackSize(ctx, size.Size); err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *packSizeService) DeletePackSize(ctx context.Context, id int) error {
	if err := s.packRepository.DeletePackSize(ctx, id); err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}
