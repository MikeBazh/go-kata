package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ptflp/geotask/module/courierfacade/controller"
)

type Router struct {
	courier *controller.CourierController
}

func NewRouter(courier *controller.CourierController) *Router {
	return &Router{courier: courier}
}

func (r *Router) CourierAPI(router *gin.RouterGroup) {
	//router.GET("/api") // прописать роуты для courier API
	// Роут для движения курьера вверх
	router.GET("/ws", r.courier.Websocket)
	router.GET("/status", r.courier.GetStatus)
}

func (r *Router) Swagger(router *gin.RouterGroup) {
	router.GET("/swagger", swaggerUI)
}
