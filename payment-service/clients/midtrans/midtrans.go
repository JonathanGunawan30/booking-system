package midtrans

import (
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type MidtransClient struct {
	ServerKey          string
	Production         bool
	SuccessCallbackURL string
}

type MidtransClientInterface interface {
	CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error)
}

func NewMidtransClient(serverKey string, production bool, successCallbackURL string) MidtransClientInterface {
	return &MidtransClient{
		ServerKey:          serverKey,
		Production:         production,
		SuccessCallbackURL: successCallbackURL,
	}
}

func (m *MidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error) {
	var (
		snapClient snap.Client
		production = midtrans.Sandbox
	)

	expiryDateTime := request.ExpiredAt
	currentTime := time.Now()
	duration := expiryDateTime.Sub(currentTime)
	if duration <= 0 {
		logrus.Errorf("expiredAt must be greater than current time")
		return nil, errPayment.ErrExpiredAtInvalid
	}

	expiryUnit := "minute"
	expiryDuration := int64(duration.Minutes())

	if duration.Hours() >= 24 {
		expiryUnit = "day"
		expiryDuration = int64(duration.Hours() / 24)
	}

	if expiryDuration <= 0 && expiryUnit == "minute" {
		expiryDuration = 1
	}

	if m.Production {
		production = midtrans.Production
	}

	snapClient.New(m.ServerKey, production)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderID,
			GrossAmt: int64(request.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerDetail.Name,
			Email: request.CustomerDetail.Email,
			Phone: request.CustomerDetail.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    request.ItemDetails[0].ID,
				Price: int64(request.ItemDetails[0].Amount),
				Qty:   int32(request.ItemDetails[0].Quantity),
				Name:  request.ItemDetails[0].Name,
			},
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     expiryUnit,
			Duration: expiryDuration,
		},
		Callbacks: &snap.Callbacks{
			Finish: m.SuccessCallbackURL,
		},
	}

	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		logrus.Errorf("failed to create payment link: %v", err)
		return nil, err
	}

	return &MidtransData{
		RedirectUrl: response.RedirectURL,
		Token:       response.Token,
	}, nil

}
