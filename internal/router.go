package internal

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "re-partners/docs"
	"re-partners/internal/handler"
	internalmiddleware "re-partners/internal/middleware"
	"re-partners/internal/repository"
	"re-partners/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	log             *zap.SugaredLogger
	packHandler     *handler.PackHandler
	orderHandler    *handler.OrderHandler
	templateHandler *handler.TemplateHandler
	packService     service.PackService
	packRepository  repository.Repository
}

func NewServer(log *zap.SugaredLogger, pool *pgxpool.Pool) *Server {
	packRepository := repository.NewPackRepository(pool)
	packService := service.NewPackSizeService(packRepository, log)
	packHandler := handler.NewPackHandler(log, packService)
	orderHandler := handler.NewOrderHandler(log, packService)
	templateHandler, err := handler.NewTemplateHandler(log, packService)
	if err != nil {
		log.Fatalw("parse templates failed", "error", err)
	}

	return &Server{
		log:             log,
		packHandler:     packHandler,
		orderHandler:    orderHandler,
		templateHandler: templateHandler,
		packService:     packService,
		packRepository:  packRepository,
	}
}

//	@title			RE Partners API
//	@version		1.0
//	@description	REST API for pack sizes and order pack calculation.
//	@schemes		http

// @host		localhost:8080
// @BasePath	/api
func (s *Server) SetupRoutes() *echo.Echo {
	e := echo.New()
	e.Use(internalmiddleware.ZapRequestLogger(s.log))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/", s.templateHandler.RenderUI)

	api := e.Group("/api")
	api.GET("*", echoSwagger.WrapHandler)

	pack := api.Group("/packs")
	pack.GET("", s.packHandler.GetPackSizes)
	pack.POST("", s.packHandler.AddPackSize)
	pack.DELETE("/:id", s.packHandler.DeletePackSize)

	order := api.Group("/orders")
	order.POST("/calculate", s.orderHandler.Calculate)

	return e
}
