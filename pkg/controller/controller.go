package controller

import (
	"server/pkg/service"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Auth(c *gin.Context)
	Message(c *gin.Context)
}

type controller struct {
	Service service.Service
}

func NewController() Controller {
	return &controller{
		Service: service.NewService(),
	}
}
