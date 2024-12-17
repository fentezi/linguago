package server

import (
	"github.com/fentezi/translator/internal/controllers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func (s *Server) Start() error {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", s.Controller.IndexHTML)
	e.GET("/saved-words", s.Controller.WordHTML)
	e.POST("/translations", s.Controller.TranslateWord)
	e.GET("/words", s.Controller.GetAllWords)
	e.DELETE("/words/:word", s.Controller.DeleteWord)
	e.POST("/add", s.Controller.AddTranslate)

	e.Static("/static", "static")

	return e.Start(":8080")
}
