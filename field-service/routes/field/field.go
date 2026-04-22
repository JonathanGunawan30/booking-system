package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middleware"

	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.ControllerRegistryInterface
	group      *gin.RouterGroup
	client     clients.ClientRegistryInterface
}

type FieldRouteInterface interface {
	Run()
}

func NewFieldRoute(controller controllers.ControllerRegistryInterface, group *gin.RouterGroup, client clients.ClientRegistryInterface) FieldRouteInterface {
	return &FieldRoute{controller: controller, group: group, client: client}
}

func (f *FieldRoute) Run() {
	group := f.group.Group("/field")
	group.GET("", middleware.AuthenticateWithoutToken(), f.controller.GetField().GetAllWithoutPagination)
	group.GET("/:uuid", middleware.AuthenticateWithoutToken(), f.controller.GetField().GetByUUID)

	group.Use(middleware.Authenticate())
	group.GET("/pagination", middleware.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetField().GetAllWithPagination)
	group.POST("", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetField().Create)
	group.PUT("/:uuid", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetField().Update)
	group.DELETE("/:uuid", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetField().Delete)
}
