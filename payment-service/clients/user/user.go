package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"payment-service/clients/config"
	"payment-service/common/util"
	config2 "payment-service/config"
	"payment-service/constants"
	"time"
)

type UserClient struct {
	client config.ClientConfigInterface
}

type UserClientInterface interface {
	GetUserByToken(ctx context.Context) (*UserData, error)
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

	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("X-Service-Name", config2.AppConfig.AppName)
	req.Header.Set("X-Request-At", unixTime)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}
