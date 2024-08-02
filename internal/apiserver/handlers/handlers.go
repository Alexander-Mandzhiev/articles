package handlers

import (
	"articles/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/articles", h.Create)
		api.GET("/articles", h.GetAll)
		api.GET("/articles/:id", h.GetOne)
		api.PUT("/articles/:id", h.Update)
		api.DELETE("articles/:id", h.Delete)
	}

	return router
}
