package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/fentezi/translator/internal/controllers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Controller controllers.Controller
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func NewServer(controller controllers.Controller) *Server {
	return &Server{
		Controller: controller,
	}
}

func (s *Server) Start(log *slog.Logger) error {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

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

	api := e.Group("/api")
	{
		api.POST("/add", s.Controller.AddTranslate)
		api.DELETE("/words/:word", s.Controller.DeleteWord)
		api.GET("/words/:word", s.Controller.GetAudio)
		api.GET("/words", s.Controller.GetAllWords)
		api.POST("/translations", s.Controller.TranslateWord)

	}

	e.GET("/", s.Controller.IndexHTML)
	e.GET("/words", s.Controller.WordHTML)

	e.Static("/static", "static")

	return e.Start(":8080")
}
