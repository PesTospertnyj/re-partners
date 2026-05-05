package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"re-partners/internal/dto"
	"re-partners/internal/service"
	"strconv"
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

func (h *PackHandler) GetPackSizes(c echo.Context) error {
	packSizes, err := h.packService.GetPackSizes(c.Request().Context())
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, packSizes)
}

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
