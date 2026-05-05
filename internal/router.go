package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"re-partners/internal/handler"
	internalmiddleware "re-partners/internal/middleware"
	"re-partners/internal/repository"
	"re-partners/internal/service"
)

type Server struct {
	log            *zap.SugaredLogger
	packHandler    *handler.PackHandler
	orderHandler   *handler.OrderHandler
	packService    service.PackService
	packRepository repository.Repository
}

func NewServer(log *zap.SugaredLogger, pool *pgxpool.Pool) *Server {
	packRepository := repository.NewPackRepository(pool)
	packService := service.NewPackSizeService(packRepository, log)
	packHandler := handler.NewPackHandler(log, packService)
	orderHandler := handler.NewOrderHandler(log, packService)

	return &Server{
		log:            log,
		packHandler:    packHandler,
		orderHandler:   orderHandler,
		packService:    packService,
		packRepository: packRepository,
	}
}

func (s *Server) SetupRoutes() *echo.Echo {
	e := echo.New()
	e.Use(internalmiddleware.ZapRequestLogger(s.log))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	api := e.Group("/api")

	pack := api.Group("/packs")
	pack.GET("", s.packHandler.GetPackSizes)
	pack.POST("", s.packHandler.AddPackSize)
	pack.DELETE("/:id", s.packHandler.DeletePackSize)

	order := api.Group("/orders")
	order.POST("/calculate", s.orderHandler.Calculate)

	return e
}
