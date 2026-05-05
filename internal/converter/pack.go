package converter

import (
	"re-partners/internal/dto"
	"re-partners/internal/model"
)

func PackSizeToDTO(pack model.PackSize) dto.PackSize {
	return dto.PackSize{
		ID:   pack.ID,
		Size: pack.Size,
	}
}

func PackSizeToModel(pack dto.PackSize) model.PackSize {
	return model.PackSize{
		ID:   pack.ID,
		Size: pack.Size,
	}
}
