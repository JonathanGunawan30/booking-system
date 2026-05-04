package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clients3 "order-service/clients/field"
	clients4 "order-service/clients/payment"
	clients2 "order-service/clients/user"
	"order-service/common/util"
	"order-service/constants"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repository repositories.RegistryRepositoryInterface
	client     clients.ClientRegistryInterface
}

type OrderServiceInterface interface {
	GetAllWithPagination(ctx context.Context, req *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetOrderByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error)
	GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(ctx context.Context, req *dto.PaymentData) error
}

func NewOrderService(repository repositories.RegistryRepositoryInterface, client clients.ClientRegistryInterface) OrderServiceInterface {
	return &OrderService{repository: repository, client: client}
}

func (o *OrderService) GetAllWithPagination(ctx context.Context, req *dto.OrderRequestParam) (*util.PaginationResult, error) {
	orders, total, err := o.repository.GetOrder().FindAllWithPagination(ctx, req)
	if err != nil {
		return nil, err
	}

	orderResults := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		user, err := o.client.GetUser().GetUserByUUID(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		orderResults = append(orderResults, dto.OrderResponse{
			UUID:      order.UUID,
			Code:      order.Code,
			UserName:  user.Name,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			OrderDate: order.Date,
			CreatedAt: *order.CreatedAt,
			UpdatedAt: *order.UpdatedAt,
		})
	}

	paginationPayload := util.PaginationParam{
		Page:  req.Page,
		Limit: req.Limit,
		Count: total,
		Data:  orderResults,
	}

	response := util.GeneratePagination(paginationPayload)
	return &response, nil
}

func (o *OrderService) GetOrderByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error) {
	order, err := o.repository.GetOrder().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user, err := o.client.GetUser().GetUserByUUID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:      order.UUID,
		Code:      order.Code,
		UserName:  user.Name,
		Amount:    order.Amount,
		Status:    order.Status.GetStatusString(),
		OrderDate: order.Date,
		CreatedAt: *order.CreatedAt,
		UpdatedAt: *order.UpdatedAt,
	}

	return &response, nil
}

func (o *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	user := ctx.Value(constants.User).(clients2.UserData)

	orderData, err := o.repository.GetOrder().FindByUserID(ctx, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderResults := make([]dto.OrderByUserIDResponse, 0, len(orderData))
	for _, order := range orderData {

		payment, err := o.client.GetPayment().GetPaymentByUUID(ctx, order.PaymentID)
		if err != nil {
			return nil, err
		}

		orderResults = append(orderResults, dto.OrderByUserIDResponse{
			Code:        order.Code,
			Amount:      util.RupiahFormat(&order.Amount),
			Status:      order.Status.GetStatusString(),
			OrderDate:   order.Date,
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}

	return orderResults, nil
}

func (o *OrderService) Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		orderCreated    *models.Order
		field           *clients3.FieldData
		user            = ctx.Value(constants.User).(*clients2.UserData)
		totalAmount     = float64(0)
		paymentResponse *clients4.PaymentData
	)

	orderFieldSchedules := make([]models.OrderField, 0, len(req.FieldScheduleIDs))
	for _, fieldID := range req.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldID)
		field, err := o.client.GetField().GetFieldByUUID(ctx, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedStatus {
			return nil, errOrder.ErrFieldAlreadyBooked
		}
	}

	err := o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		order, txErr := o.repository.GetOrder().Create(ctx, tx, &models.Order{
			UserID: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})
		if txErr != nil {
			return txErr
		}

		orderCreated = order

		for _, fieldID := range req.FieldScheduleIDs {
			uuidParsed := uuid.MustParse(fieldID)
			orderFieldSchedules = append(orderFieldSchedules, models.OrderField{
				OrderID:         order.ID,
				FieldScheduleID: uuidParsed,
			})
		}

		txErr = o.repository.GetOrderField().Create(ctx, tx, orderFieldSchedules)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(time.Hour * 1)
		description := fmt.Sprintf("Payment Order %s", field.FieldName)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(ctx, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAt,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: dto.CustomerDetail{
				Name:        user.Name,
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetails{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().Update(ctx, tx, &models.Order{
			PaymentID: paymentResponse.UUID,
		}, order.UUID)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		UUID:        orderCreated.UUID,
		Code:        orderCreated.Code,
		UserName:    user.Name,
		Amount:      orderCreated.Amount,
		Status:      orderCreated.Status.GetStatusString(),
		OrderDate:   orderCreated.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   *orderCreated.CreatedAt,
		UpdatedAt:   *orderCreated.UpdatedAt,
	}, nil
}

func (o *OrderService) mapPaymentStatusToOrder(req *dto.PaymentData) (constants.OrderStatus, *models.Order) {
	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch req.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: req.PaymentID,
			PaidAt:    req.PaidAt,
			Status:    status,
		}
	case constants.ExpiredPaymentStatus:
		status = constants.Expired
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.PendingPayment
		order = &models.Order{
			IsPaid:    false,
			PaymentID: req.PaymentID,
			Status:    status,
		}
	}
	return status, order
}

func (o *OrderService) HandlePayment(ctx context.Context, req *dto.PaymentData) error {
	var (
		err, txErr          error
		order               *models.Order
		orderFieldSchedules []models.OrderField
	)
	status, body := o.mapPaymentStatusToOrder(req)
	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repository.GetOrder().Update(ctx, tx, body, req.OrderID)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		if req.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
			if txErr != nil {
				return txErr
			}

			fieldScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				fieldScheduleIDs = append(fieldScheduleIDs, item.FieldScheduleID.String())
			}

			txErr = o.client.GetField().UpdateStatus(ctx, &dto.UpdateFieldScheduleStatusRequest{
				FieldScheduleIDs: fieldScheduleIDs,
			})
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})

	return err
}
