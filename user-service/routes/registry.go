package routes

import (
	"user-service/controllers"
	_ "user-service/docs"
	routes "user-service/routes/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.userRoute().Run()
}

func (r *Registry) userRoute() routes.UserRouteInterface {
	return routes.NewUserRoute(r.controller, r.group)
}
