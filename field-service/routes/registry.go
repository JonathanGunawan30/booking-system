package routes

import (
	"field-service/clients"
	"field-service/controllers"
	fieldRoute "field-service/routes/field"
	fieldScheduleRoute "field-service/routes/fieldschedule"
	timeRoute "field-service/routes/time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.ControllerRegistryInterface
	route      *gin.RouterGroup
	client     clients.ClientRegistryInterface
}

type RegistryInterface interface {
	Serve()
}

func NewRouteRegistry(router *gin.RouterGroup, controller controllers.ControllerRegistryInterface, client clients.ClientRegistryInterface) RegistryInterface {
	return &Registry{controller: controller, route: router, client: client}
}

func (r *Registry) fieldRoute() fieldRoute.FieldRouteInterface {
	return fieldRoute.NewFieldRoute(r.controller, r.route, r.client)
}

func (r *Registry) fieldScheduleRoute() fieldScheduleRoute.FieldScheduleRouteInterface {
	return fieldScheduleRoute.NewFieldRoute(r.controller, r.route, r.client)
}

func (r *Registry) timeRoute() timeRoute.TimeRouteInterface {
	return timeRoute.NewFieldRoute(r.controller, r.route, r.client)
}

func (r *Registry) Serve() {
	r.fieldRoute().Run()
	r.fieldScheduleRoute().Run()
	r.timeRoute().Run()
}
