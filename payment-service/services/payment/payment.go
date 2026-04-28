package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"payment-service/clients/midtrans"
	"payment-service/common/cloudflare"
	"payment-service/common/util"
	config2 "payment-service/config"
	"payment-service/constants"
	errPayment "payment-service/constants/error/payment"
	"payment-service/controllers/kafka"
	"payment-service/domain/dto"
	"payment-service/domain/models"
	"payment-service/repositories"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentService struct {
	repository repositories.RepositoryRegistryInterface
	r2         cloudflare.R2Client
	kafka      kafka.KafkaRegistryInterface
	midtrans   midtrans.MidtransClientInterface
}

type PaymentServiceInterface interface {
	GetAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error)
	Create(ctx context.Context, request *dto.PaymentRequest) (*dto.PaymentResponse, error)
	WebHook(ctx context.Context, hook *dto.WebHook) error
}

func NewPaymentService(repository repositories.RepositoryRegistryInterface, r2 cloudflare.R2Client, kafka kafka.KafkaRegistryInterface, midtrans midtrans.MidtransClientInterface) PaymentServiceInterface {
	return &PaymentService{repository: repository, r2: r2, kafka: kafka, midtrans: midtrans}
}

func (p *PaymentService) GetAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	payments, total, err := p.repository.GetPayment().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	paymentResult := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		paymentResult = append(paymentResult, dto.PaymentResponse{
			UUID:          payment.UUID,
			TransactionID: payment.TransactionID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Status:        payment.Status.GetStatusString(),
			PaymentLink:   payment.PaymentLink,
			InvoiceLink:   payment.InvoiceLink,
			VANumber:      payment.VANumber,
			Bank:          payment.Bank,
			Acquirer:      payment.Acquirer,
			Description:   payment.Description,
			ExpiredAt:     payment.ExpiredAt,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Count: total,
		Limit: param.Limit,
		Data:  paymentResult,
	}

	response := util.GeneratePagination(paginationParam)
	return &response, nil
}

func (p *PaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {
	payment, err := p.repository.GetPayment().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		TransactionID: payment.TransactionID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		PaymentLink:   payment.PaymentLink,
		InvoiceLink:   payment.InvoiceLink,
		VANumber:      payment.VANumber,
		Bank:          payment.Bank,
		Acquirer:      payment.Acquirer,
		Description:   payment.Description,
		PaidAt:        payment.PaidAt,
		ExpiredAt:     payment.ExpiredAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (p *PaymentService) Create(ctx context.Context, request *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var create *models.Payment

	err := p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		if !request.ExpiredAt.After(time.Now()) {
			return errPayment.ErrExpiredAtInvalid
		}

		midtransData, txErr := p.midtrans.CreatePaymentLink(request)
		if txErr != nil {
			return txErr
		}

		paymentRequest := &dto.PaymentRequest{
			OrderID:     request.OrderID,
			Amount:      request.Amount,
			PaymentLink: midtransData.RedirectUrl,
			ExpiredAt:   request.ExpiredAt,
			Description: request.Description,
		}

		create, txErr = p.repository.GetPayment().Create(ctx, tx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		txErr = p.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: create.ID,
			Status:    create.Status.GetStatusString(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	response := &dto.PaymentResponse{
		UUID:        create.UUID,
		OrderID:     create.OrderID,
		Amount:      create.Amount,
		Status:      create.Status.GetStatusString(),
		PaymentLink: create.PaymentLink,
		Description: create.Description,
	}

	return response, nil
}

func (p *PaymentService) generatePDF(request *dto.InvoiceRequest) ([]byte, error) {
	htmlTemplatePath := "template/invoice.html"
	htmlTemplate, err := os.ReadFile(htmlTemplatePath)
	if err != nil {
		return nil, err
	}

	var data = map[string]any{}
	jsonData, _ := json.Marshal(request)
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	pdf, err := util.GeneratePDFFromHTML(string(htmlTemplate), data)
	if err != nil {
		return nil, err
	}

	return pdf, nil
}

func (p *PaymentService) uploadToR2(ctx context.Context, invoice string, pdf []byte) (string, error) {
	invoiceNumberReplace := strings.ToLower(strings.ReplaceAll(invoice, "/", ""))
	fileName := fmt.Sprintf("%s.pdf", invoiceNumberReplace)
	url, err := p.r2.Upload(fileName, bytes.NewReader(pdf), "application/pdf")
	if err != nil {
		return "", err
	}
	return url, nil
}

func (p *PaymentService) randomNumber() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := random.Intn(900000) + 100000
	return number
}

func (p *PaymentService) mapTransactionStatusToEvent(status constants.PaymentStatusString) string {
	lowerStatus := strings.ToLower(strings.TrimSpace(string(status)))
	switch constants.PaymentStatusString(lowerStatus) {
	case constants.PendingString:
		return "PENDING"
	case constants.SettlementString:
		return "SETTLEMENT"
	case constants.ExpiredString:
		return "EXPIRE"
	default:
		return strings.ToUpper(lowerStatus)
	}
}

func (p *PaymentService) produceToKafka(req *dto.WebHook, payment *models.Payment, paidAt *time.Time) error {
	if payment == nil {
		return fmt.Errorf("payment data is nil")
	}
	event := dto.KafkaEvent{
		Name: p.mapTransactionStatusToEvent(req.TransactionStatus),
	}

	metaData := dto.KafkaMetaData{
		Sender:    "payment-service",
		SendingAt: time.Now().Format(time.RFC3339),
	}

	var expiredAt time.Time
	if payment.ExpiredAt != nil {
		expiredAt = *payment.ExpiredAt
	}

	body := dto.KafkaBody{
		Type: "JSON",
		Data: &dto.KafkaData{
			OrderID:   payment.OrderID,
			PaymentID: payment.UUID,
			Status:    string(req.TransactionStatus),
			ExpiredAt: expiredAt,
			PaidAt:    paidAt,
		},
	}

	kafkaMsg := dto.KafkaMessage{
		Event:    event,
		MetaData: metaData,
		Body:     body,
	}

	topic := config2.AppConfig.Kafka.Topic
	kafkaMessageJSON, _ := json.Marshal(kafkaMsg)
	err := p.kafka.GetKafkaProducer().ProduceMessage(topic, kafkaMessageJSON)
	return err
}

func (p *PaymentService) WebHook(ctx context.Context, req *dto.WebHook) error {
	var (
		paidAt        *time.Time
		paymentUpdate *models.Payment
	)

	err := p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		_, txErr := p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		if req.TransactionStatus == constants.SettlementString {
			now := time.Now()
			paidAt = &now
		}

		status := req.TransactionStatus.GetStatusInt()

		var vaNumber, bank string
		if len(req.VANumbers) > 0 {
			vaNumber = req.VANumbers[0].VANumber
			bank = req.VANumbers[0].Bank
		}

		_, txErr = p.repository.GetPayment().Update(ctx, tx, &dto.UpdatePaymentRequest{
			TransactionID: &req.TransactionID,
			Status:        &status,
			PaidAt:        paidAt,
			VANumber:      &vaNumber,
			Bank:          &bank,
			Acquirer:      req.Acquirer,
		}, req.OrderID.String())

		if txErr != nil {
			return txErr
		}

		paymentUpdate, txErr = p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = p.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: paymentUpdate.ID,
			Status:    paymentUpdate.Status.GetStatusString(),
		})

		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return err
	}

	paymentUpdate, err = p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
	if err != nil {
		logrus.Errorf("[WebHook] failed to find payment for invoice: %v", err)
		return err
	}

	if strings.ToLower(string(req.TransactionStatus)) == strings.ToLower(string(constants.SettlementString)) {
		datePaid := paidAt.Format(time.DateOnly)
		invoiceNumber := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format(time.DateOnly), p.randomNumber())
		total := util.RupiahFormat(&paymentUpdate.Amount)
		invoiceRequest := &dto.InvoiceRequest{
			InvoiceNumber: invoiceNumber,
			Data: dto.InvoiceData{
				PaymentDetail: dto.InvoicePaymentDetail{
					PaymentMethod: req.PaymentType,
					BankName:      strings.ToUpper(util.GetValueOrDefault(paymentUpdate.Bank)),
					VANumber:      util.GetValueOrDefault(paymentUpdate.VANumber),
					Date:          datePaid,
					IsPaid:        true,
				},
				Items: []dto.InvoiceItem{
					{
						Description: util.GetValueOrDefault(paymentUpdate.Description),
						Price:       total,
					},
				},
				Total: total,
			},
		}
		pdf, invErr := p.generatePDF(invoiceRequest)
		if invErr != nil {
			logrus.Errorf("[WebHook] failed to generate PDF: %v", invErr)
		} else {
			invoiceLink, uploadErr := p.uploadToR2(ctx, invoiceNumber, pdf)
			if uploadErr != nil {
				logrus.Errorf("[WebHook] failed to upload PDF to R2: %v", uploadErr)
			} else {
				err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
					_, txErr := p.repository.GetPayment().Update(ctx, tx, &dto.UpdatePaymentRequest{
						InvoiceLink: &invoiceLink,
					}, req.OrderID.String())
					return txErr
				})
				if err != nil {
					logrus.Errorf("[WebHook] failed to update invoice link in DB: %v", err)
				} else {
					logrus.Infof("[WebHook] invoice generated and uploaded: %s", invoiceLink)
				}
			}
		}
	}

	paymentUpdate, err = p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
	if err != nil {
		logrus.Errorf("[WebHook] failed to find payment for final kafka: %v", err)
		return err
	}

	err = p.produceToKafka(req, paymentUpdate, paidAt)
	return err
}
