package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middleware"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.ControllerRegistryInterface
	group      *gin.RouterGroup
	client     clients.ClientRegistryInterface
}

type TimeRouteInterface interface {
	Run()
}

func NewFieldRoute(controller controllers.ControllerRegistryInterface, group *gin.RouterGroup, client clients.ClientRegistryInterface) TimeRouteInterface {
	return &TimeRoute{controller: controller, group: group, client: client}
}

func (t *TimeRoute) Run() {
	group := t.group.Group("/time")
	group.Use(middleware.Authenticate(), middleware.CheckRole([]string{constants.Admin}, t.client))
	group.GET("", t.controller.GetTime().GetAll)
	group.GET("/:uuid", t.controller.GetTime().GetByUUID)
	group.POST("", t.controller.GetTime().Create)
}
