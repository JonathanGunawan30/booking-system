package routes

import (
	"user-service/controllers"
	"user-service/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.UserControllerRegistryInterface
	group      *gin.RouterGroup
}

type UserRouteInterface interface {
	Run()
}

func NewUserRoute(controller controllers.UserControllerRegistryInterface, group *gin.RouterGroup) UserRouteInterface {
	return &UserRoute{controller: controller, group: group}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middleware.Authenticate(), u.controller.GetUserController().GetUserLogin)
	group.GET("/user/:uuid", middleware.AuthenticateWithoutToken(), u.controller.GetUserController().GetUserByUUID)
	group.POST("/login", u.controller.GetUserController().Login)
	group.POST("/register", u.controller.GetUserController().Register)
	group.PUT("/:uuid", middleware.Authenticate(), u.controller.GetUserController().Update)
}
