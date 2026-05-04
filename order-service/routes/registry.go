package routes

import (
	"order-service/clients"
	controllers "order-service/controllers/http"
	routes "order-service/routes/order"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.ControllerRegistryInterface
	client     clients.ClientRegistryInterface
	group      *gin.RouterGroup
}

type RouteRegistryInterface interface {
	Serve()
}

func NewRouteRegistry(controller controllers.ControllerRegistryInterface, client clients.ClientRegistryInterface, group *gin.RouterGroup) RouteRegistryInterface {
	return &Registry{controller: controller, client: client, group: group}
}

func (r *Registry) orderRoute() routes.OrderRouteInterface {
	return routes.NewOrderRoute(r.controller, r.client, r.group)
}

func (r *Registry) Serve() {
	r.orderRoute().Run()
}
