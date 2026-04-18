package routes

import (
	"user-service/controllers"
	routes "user-service/routes/user"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.UserControllerRegistryInterface
	group      *gin.RouterGroup
}

type RouteRegistryInterface interface {
	Serve()
}

func NewRouteRegistry(controller controllers.UserControllerRegistryInterface, group *gin.RouterGroup) RouteRegistryInterface {
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}

func (r *Registry) userRoute() routes.UserRouteInterface {
	return routes.NewUserRoute(r.controller, r.group)
}
