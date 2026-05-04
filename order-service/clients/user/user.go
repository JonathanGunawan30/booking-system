package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/clients/config"
	"order-service/common/util"
	config2 "order-service/config"
	"order-service/constants"
	errConstant "order-service/constants/error"
	"time"

	"github.com/google/uuid"
)

type UserClient struct {
	client config.ClientConfigInterface
}

type UserClientInterface interface {
	GetUserByToken(ctx context.Context) (*UserData, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error)
}

func NewUserClient(client config.ClientConfigInterface) UserClientInterface {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().UTC().Format(time.RFC3339)

	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)

	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	url := fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constants.Authorization, bearerToken)
	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.XServiceName, config2.AppConfig.AppName)
	req.Header.Set(constants.XRequestAt, unixTime)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errConstant.ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errConstant.ErrInternalServerError
	}

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
	unixTime := time.Now().UTC().Format(time.RFC3339)

	generateAPIKey := fmt.Sprintf("%s:%s:%s",
		config2.AppConfig.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)

	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	url := fmt.Sprintf("%s/api/v1/auth/%s", u.client.BaseURL(), uuid)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constants.XApiKey, apiKey)
	req.Header.Set(constants.Authorization, bearerToken)
	req.Header.Set(constants.XServiceName, config2.AppConfig.AppName)
	req.Header.Set(constants.XRequestAt, unixTime)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errConstant.ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errConstant.ErrInternalServerError
	}

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response.Data, nil
}
