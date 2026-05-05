package handler

import (
	"net/http"
	"re-partners/internal/dto"
	"re-partners/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PackHandler struct {
	log         *zap.SugaredLogger
	packService service.PackService
}

func NewPackHandler(log *zap.SugaredLogger, packService service.PackService) *PackHandler {
	return &PackHandler{
		log:         log,
		packService: packService,
	}
}

// GetPackSizes returns all available pack sizes.
//
//	@Summary		List pack sizes
//	@Description	Returns pack sizes ordered ascending.
//	@Tags			packs
//	@Produce		json
//	@Success		200	{array}		dto.PackSize
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/packs [get]
func (h *PackHandler) GetPackSizes(c echo.Context) error {
	packSizes, err := h.packService.GetPackSizes(c.Request().Context())
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, packSizes)
}

// AddPackSize creates new pack size.
//
//	@Summary		Add pack size
//	@Description	Creates a pack size used for order calculations.
//	@Tags			packs
//	@Accept			json
//	@Param			request	body	dto.AddPackSize	true	"Pack size request"
//	@Success		201
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/packs [post]
func (h *PackHandler) AddPackSize(c echo.Context) error {
	addPackSizeDto := dto.AddPackSize{}
	if err := c.Bind(&addPackSizeDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.packService.AddPackSize(c.Request().Context(), dto.PackSize{Size: addPackSizeDto.Size}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// DeletePackSize deletes pack size by id.
//
//	@Summary		Delete pack size
//	@Description	Deletes pack size by identifier.
//	@Tags			packs
//	@Param			id	path	int	true	"Pack size id"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/packs/{id} [delete]
func (h *PackHandler) DeletePackSize(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.packService.DeletePackSize(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
