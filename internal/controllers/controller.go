package controllers

import (
	"github.com/fentezi/translator/internal/services"
	"github.com/fentezi/translator/pkg/vld"
)

type Controller struct {
	service   *services.Service
	validator *vld.Validator
}

func New(service *services.Service) *Controller {
	return &Controller{
		service:   service,
		validator: vld.New(),
	}
}
