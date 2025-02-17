package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/fentezi/translator/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Controller controllers.Controller
}

func New(controller controllers.Controller) *Server {
	return &Server{
		Controller: controller,
	}
}

func (s *Server) Start(log *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			LogMethod:   true,
			LogRemoteIP: true,
			LogLatency:  true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {

					log.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("method", v.Method),
						slog.String("remote_ip", v.RemoteIP),
						slog.Duration("latency", time.Duration(v.Latency.Seconds())),
					)
				} else {
					log.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("method", v.Method),
						slog.String("err", v.Error.Error()),
					)
				}
				return nil
			},
		}))

	api := e.Group("/api/v1")
	{
		api.GET("/words", s.Controller.GetWords)
		api.POST("/words", s.Controller.CreateWord)
		api.DELETE("/words/:word_id", s.Controller.DeleteWord)
		api.GET("/words/:word_id/audio", s.Controller.GetAudio)
		api.POST("/translations", s.Controller.TranslateWord)

	}

	return e
}
