package handler

import (
	"net/http"
	"re-partners/internal/dto"
	"re-partners/internal/service"
	"re-partners/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type OrderHandler struct {
	log         *zap.SugaredLogger
	packService service.PackService
}

func NewOrderHandler(log *zap.SugaredLogger, packService service.PackService) *OrderHandler {
	return &OrderHandler{
		log:         log,
		packService: packService,
	}
}

// Calculate calculates minimum packs for requested quantity.
//
//	@Summary		Calculate order packs
//	@Description	Returns minimum total items and pack breakdown for requested order size.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.Calculate	true	"Order calculation request"
//	@Success		200		{object}	dto.CalculateResult
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/orders/calculate [post]
func (h *OrderHandler) Calculate(c echo.Context) error {
	calculateDto := dto.Calculate{}
	if err := c.Bind(&calculateDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sizes, err := h.packService.GetPackSizes(c.Request().Context())
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ascSizes := make([]int, len(sizes))
	for i, s := range sizes {
		ascSizes[i] = s.Size
	}

	minTotalItems := utils.MinimumTotalItems(calculateDto.ItemsOrdered, ascSizes)
	minPacks := utils.MinimumPackBreakdown(minTotalItems, ascSizes)

	return c.JSON(http.StatusOK, dto.CalculateResult{
		ItemsOrdered: calculateDto.ItemsOrdered,
		TotalItems:   minTotalItems,
		Packs:        minPacks,
	})
}
