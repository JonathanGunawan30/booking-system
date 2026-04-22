package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middleware"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.ControllerRegistryInterface
	group      *gin.RouterGroup
	client     clients.ClientRegistryInterface
}

type FieldScheduleRouteInterface interface {
	Run()
}

func NewFieldRoute(controller controllers.ControllerRegistryInterface, group *gin.RouterGroup, client clients.ClientRegistryInterface) FieldScheduleRouteInterface {
	return &FieldScheduleRoute{controller: controller, group: group, client: client}
}

func (f *FieldScheduleRoute) Run() {
	public := f.group.Group("/field/schedule")
	public.GET("/lists/:uuid", middleware.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().GetAllByFieldIDAndDate)
	public.PATCH("", middleware.AuthenticateWithoutToken(), f.controller.GetFieldSchedule().UpdateStatus)

	private := f.group.Group("/field/schedule")
	private.Use(middleware.Authenticate())
	private.GET("/pagination", middleware.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetFieldSchedule().GetAllWithPagination)
	private.GET("/:uuid", middleware.CheckRole([]string{constants.Admin, constants.Customer}, f.client), f.controller.GetFieldSchedule().GetByUUID)
	private.POST("", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Create)
	private.POST("/one-month", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().GenerateFieldScheduleForOneMonth)
	private.PUT("/:uuid", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Update)
	private.DELETE("/:uuid", middleware.CheckRole([]string{constants.Admin}, f.client), f.controller.GetFieldSchedule().Delete)
}
