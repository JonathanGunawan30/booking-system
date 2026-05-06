package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/clients/config"
	"order-service/common/util"
	config2 "order-service/config"
	"order-service/constants"
	"order-service/domain/dto"
	"time"

	"github.com/google/uuid"
)

type PaymentClient struct {
	client config.ClientConfigInterface
}

type PaymentClientInterface interface {
	GetPaymentByUUID(ctx context.Context, uuid uuid.UUID) (*PaymentData, error)
	CreatePaymentLink(ctx context.Context, req *dto.PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client config.ClientConfigInterface) PaymentClientInterface {
	return &PaymentClient{client: client}
}

func (p *PaymentClient) GetPaymentByUUID(ctx context.Context, uuid uuid.UUID) (*PaymentData, error) {
	unixTime := time.Now().UTC().Format(time.RFC3339)
	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		p.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	url := fmt.Sprintf("%s/api/v1/payments/%s", p.client.BaseURL(), uuid.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constants.Authorization, bearerToken)
	req.Header.Set(constants.XServiceName, config2.AppConfig.AppName)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XRequestAt, unixTime)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment response: %s", response.Message)
	}

	return &response.Data, nil
}

func (p *PaymentClient) CreatePaymentLink(ctx context.Context, req *dto.PaymentRequest) (*PaymentData, error) {
	unixTime := time.Now().UTC().Format(time.RFC3339)
	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		p.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/payments", p.client.BaseURL())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set(constants.Authorization, bearerToken)
	httpReq.Header.Set(constants.XServiceName, config2.AppConfig.AppName)
	httpReq.Header.Set(constants.XApiKey, apiKey)
	httpReq.Header.Set(constants.XRequestAt, unixTime)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("payment response: %s", response.Message)
	}

	return &response.Data, nil
}
