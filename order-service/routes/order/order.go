package routes

import (
	"order-service/clients"
	"order-service/constants"
	controllers "order-service/controllers/http"
	middleware "order-service/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controller controllers.ControllerRegistryInterface
	client     clients.ClientRegistryInterface
	group      *gin.RouterGroup
}

type OrderRouteInterface interface {
	Run()
}

func NewOrderRoute(controller controllers.ControllerRegistryInterface, client clients.ClientRegistryInterface, group *gin.RouterGroup) OrderRouteInterface {
	return &OrderRoute{controller: controller, client: client, group: group}
}

func (o *OrderRoute) Run() {
	group := o.group.Group("/order")
	group.Use(middleware.Authenticate())

	allRoles := group.Group("")
	allRoles.Use(middleware.CheckRole([]string{constants.Admin, constants.Customer}, o.client))
	{
		allRoles.GET("", o.controller.GetOrder().GetAllWithPagination)
		allRoles.GET("/:uuid", o.controller.GetOrder().GetByUUID)
	}

	customerOnly := group.Group("")
	customerOnly.Use(middleware.CheckRole([]string{constants.Customer}, o.client))
	{
		customerOnly.GET("/user", o.controller.GetOrder().GetOrderByUserID)
		customerOnly.POST("", o.controller.GetOrder().Create)
	}
}
