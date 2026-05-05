package handler

import (
	"html/template"
	"net/http"
	"re-partners/internal/dto"
	"re-partners/internal/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	log         *zap.SugaredLogger
	packService service.PackService
	tmpl        *template.Template
}

type UIData struct {
	Packs []dto.PackSize
}

func NewTemplateHandler(log *zap.SugaredLogger, packService service.PackService) (*TemplateHandler, error) {
	tmpl, err := template.ParseGlob("internal/template/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &TemplateHandler{
		log:         log,
		packService: packService,
		tmpl:        tmpl,
	}, nil
}

func (h *TemplateHandler) RenderUI(c echo.Context) error {
	packs, err := h.packService.GetPackSizes(c.Request().Context())
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := UIData{Packs: packs}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	if err := h.tmpl.ExecuteTemplate(c.Response().Writer, "base", data); err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
