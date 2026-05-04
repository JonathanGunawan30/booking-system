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

type FieldClient struct {
	client config.ClientConfigInterface
}

type FieldClientInterface interface {
	GetFieldByUUID(ctx context.Context, uuid uuid.UUID) (*FieldData, error)
	UpdateStatus(ctx context.Context, request *dto.UpdateFieldScheduleStatusRequest) error
}

func NewFieldClient(client config.ClientConfigInterface) FieldClientInterface {
	return &FieldClient{client: client}
}

func (f *FieldClient) GetFieldByUUID(ctx context.Context, uuid uuid.UUID) (*FieldData, error) {
	unixTime := time.Now().UTC().Format(time.RFC3339)
	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		f.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	url := fmt.Sprintf("%s/api/v1/field/%s", f.client.BaseURL(), uuid.String())

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

	var response FieldResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("field response: %s", response.Message)
	}

	return &response.Data, nil
}

func (f *FieldClient) UpdateStatus(ctx context.Context, request *dto.UpdateFieldScheduleStatusRequest) error {
	unixTime := time.Now().UTC().Format(time.RFC3339)
	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		f.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)

	url := fmt.Sprintf("%s/api/v1/field/status", f.client.BaseURL())

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set(constants.XServiceName, config2.AppConfig.AppName)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XRequestAt, unixTime)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response FieldResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("field response: %s", response.Message)
	}

	return nil

}
