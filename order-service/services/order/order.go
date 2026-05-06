package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clients4 "order-service/clients/payment"
	clients2 "order-service/clients/user"
	error2 "order-service/common/error"
	"order-service/common/util"
	"order-service/constants"
	errConstant "order-service/constants/error"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

		payment, err := o.client.GetPayment().GetPaymentByUUID(ctx, order.PaymentID)
		if err != nil {
			logrus.Warnf("[GetAllWithPagination] failed to get payment for order %s: %v", order.UUID, err)
		}

		var paymentLink string
		var invoiceLink *string
		if payment != nil {
			paymentLink = payment.PaymentLink
			invoiceLink = payment.InvoiceLink
		}

		orderFields, err := o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		schedules := make([]dto.FieldData, 0, len(orderFields))
		for _, of := range orderFields {
			schedule, err := o.client.GetField().GetFieldByUUID(ctx, of.FieldScheduleID)
			if err != nil {
				logrus.Warnf("[GetAllWithPagination] failed to get schedule details for %s: %v", of.FieldScheduleID, err)
				continue
			}
			schedules = append(schedules, *schedule)
		}

		orderResults = append(orderResults, dto.OrderResponse{
			UUID:        order.UUID,
			Code:        order.Code,
			UserName:    user.Name,
			Amount:      order.Amount,
			Status:      order.Status.GetStatusString(),
			PaymentLink: paymentLink,
			InvoiceLink: invoiceLink,
			Schedules:   schedules,
			OrderDate:   order.Date,
			CreatedAt:   *order.CreatedAt,
			UpdatedAt:   *order.UpdatedAt,
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

	payment, err := o.client.GetPayment().GetPaymentByUUID(ctx, order.PaymentID)
	if err != nil {
		logrus.Warnf("[GetOrderByUUID] failed to get payment for order %s: %v", order.UUID, err)
	}

	var paymentLink string
	var invoiceLink *string
	if payment != nil {
		paymentLink = payment.PaymentLink
		invoiceLink = payment.InvoiceLink
	}

	orderFields, err := o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	schedules := make([]dto.FieldData, 0, len(orderFields))
	for _, of := range orderFields {
		schedule, err := o.client.GetField().GetFieldByUUID(ctx, of.FieldScheduleID)
		if err != nil {
			logrus.Warnf("[GetOrderByUUID] failed to get schedule details for %s: %v", of.FieldScheduleID, err)
			continue
		}
		schedules = append(schedules, *schedule)
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		UserName:    user.Name,
		Amount:      order.Amount,
		Status:      order.Status.GetStatusString(),
		PaymentLink: paymentLink,
		InvoiceLink: invoiceLink,
		Schedules:   schedules,
		OrderDate:   order.Date,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}

	return &response, nil
}

func (o *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	user, ok := ctx.Value(constants.User).(*clients2.UserData)
	if !ok || user == nil {
		return nil, errConstant.ErrUnauthorized
	}

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

		orderFields, err := o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		schedules := make([]dto.FieldData, 0, len(orderFields))
		for _, of := range orderFields {
			schedule, err := o.client.GetField().GetFieldByUUID(ctx, of.FieldScheduleID)
			if err != nil {
				logrus.Warnf("[GetOrderByUserID] failed to get schedule details for %s: %v", of.FieldScheduleID, err)
				continue
			}
			schedules = append(schedules, *schedule)
		}

		orderResults = append(orderResults, dto.OrderByUserIDResponse{
			Code:        order.Code,
			Amount:      util.RupiahFormat(&order.Amount),
			Status:      order.Status.GetStatusString(),
			OrderDate:   order.Date,
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
			Schedules:   schedules,
		})
	}

	return orderResults, nil
}

func (o *OrderService) Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	user, ok := ctx.Value(constants.User).(*clients2.UserData)
	if !ok || user == nil {
		return nil, errConstant.ErrUnauthorized
	}

	var (
		orderCreated    *models.Order
		field           *dto.FieldData
		totalAmount     = float64(0)
		paymentResponse *clients4.PaymentData
		err             error
	)

	orderFieldSchedules := make([]models.OrderField, 0, len(req.FieldScheduleIDs))
	for _, fieldID := range req.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldID)
		field, err = o.client.GetField().GetFieldByUUID(ctx, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedStatus {
			return nil, error2.WrapError(errOrder.ErrFieldAlreadyBooked)
		}
	}

	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
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

	if orderCreated == nil {
		return nil, error2.WrapError(errOrder.ErrOrderFailed)
	}

	if paymentResponse == nil {
		return nil, error2.WrapError(errOrder.ErrPaymentLinkCreationFailed)
	}

	var createdAt, updatedAt time.Time
	if orderCreated.CreatedAt != nil {
		createdAt = *orderCreated.CreatedAt
	}
	if orderCreated.UpdatedAt != nil {
		updatedAt = *orderCreated.UpdatedAt
	}

	return &dto.OrderResponse{
		UUID:        orderCreated.UUID,
		Code:        orderCreated.Code,
		UserName:    user.Name,
		Amount:      orderCreated.Amount,
		Status:      orderCreated.Status.GetStatusString(),
		OrderDate:   orderCreated.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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
		orderData           *models.Order
		orderFieldSchedules []models.OrderField
	)
	
	orderData, err = o.repository.GetOrder().FindByUUID(ctx, req.OrderID.String())
	if err != nil {
		return err
	}

	status, body := o.mapPaymentStatusToOrder(req)
	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repository.GetOrder().Update(ctx, tx, body, req.OrderID)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: orderData.ID,
		})
		if txErr != nil {
			return txErr
		}

		if req.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repository.GetOrderField().FindByOrderID(ctx, orderData.ID)
			if txErr != nil {
				return txErr
			}

			fieldScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				fieldScheduleIDs = append(fieldScheduleIDs, item.FieldScheduleID.String())
			}

			logrus.Infof("[HandlePayment] Attempting to update field schedule status for IDs: %v", fieldScheduleIDs)
			txErr = o.client.GetField().UpdateStatus(ctx, &dto.UpdateFieldScheduleStatusRequest{
				FieldScheduleIDs: fieldScheduleIDs,
				Status:           200,
			})
			if txErr != nil {
				logrus.Errorf("[HandlePayment] Failed to update field status: %v", txErr)
				return txErr
			}
			logrus.Infof("[HandlePayment] Successfully sent update status request to field-service")
			}

		return nil
	})

	return err
}
