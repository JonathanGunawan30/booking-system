package payment

import (
	"payment-service/clients"
	"payment-service/constants"
	controllers "payment-service/controllers/payment"
	middleware "payment-service/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller controllers.RegistryControllerInterface
	client     clients.ClientRegistryInterface
	group      *gin.RouterGroup
}

type PaymentRouteInterface interface {
	Run()
}

func NewPaymentRoute(group *gin.RouterGroup, controller controllers.RegistryControllerInterface, client clients.ClientRegistryInterface) PaymentRouteInterface {
	return &PaymentRoute{controller: controller, client: client, group: group}
}

func (p *PaymentRoute) Run() {
	public := p.group.Group("/payments")
	{
		public.POST("/webhook", p.controller.GetPayment().WebHook)
	}

	protected := p.group.Group("/payments")
	protected.Use(middleware.Authenticate())
	{
		protected.GET("", middleware.CheckRole([]string{constants.Admin, constants.Customer}, p.client), p.controller.GetPayment().GetAllWithPagination)
		protected.GET("/:uuid", middleware.CheckRole([]string{constants.Admin, constants.Customer}, p.client), p.controller.GetPayment().GetByUUID)
		protected.POST("", middleware.CheckRole([]string{constants.Admin}, p.client), p.controller.GetPayment().Create)
	}
}
