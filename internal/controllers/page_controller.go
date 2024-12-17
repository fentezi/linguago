package controllers

import "github.com/labstack/echo/v4"

func (h *Controller) IndexHTML(c echo.Context) error {
	return c.File("./templates/index.html")
}

func (h *Controller) WordHTML(c echo.Context) error {
	return c.File("./templates/word.html")
}
