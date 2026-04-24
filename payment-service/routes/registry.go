package routes

import (
	"payment-service/clients"
	controllers "payment-service/controllers/payment"
	"payment-service/routes/payment"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.RegistryControllerInterface
	group      *gin.RouterGroup
	client     clients.ClientRegistryInterface
}

type RouteRegistryInterface interface {
	Serve()
}

func NewRouteRegistry(group *gin.RouterGroup, controller controllers.RegistryControllerInterface, client clients.ClientRegistryInterface) RouteRegistryInterface {
	return &Registry{client: client, group: group, controller: controller}
}

func (r *Registry) Serve() {
	payment.NewPaymentRoute(r.group, r.controller, r.client).Run()
}
