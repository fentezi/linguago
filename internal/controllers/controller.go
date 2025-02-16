package controllers

import "github.com/fentezi/translator/internal/services"

type Controller struct {
	service *services.Service
}

func New(service *services.Service) *Controller {
	return &Controller{
		service: service,
	}
}
