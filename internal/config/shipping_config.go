// payments/internal/config/payments_config.go
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
	"time"
)

type DHLShipmentResponse struct {
	ShipmentID string `json:"shipmentId"`
	Status     string `json:"status"`
	LabelURL   string `json:"labelUrl"`
}
type DHLClient struct {
	client       *resty.Client
	clientID     string
	clientSecret string
	apiEndpoint  string
	accessToken  string
	tokenExpiry  time.Time
}

type ShippingConfig struct {
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
	AccessToken  string `envconfig:"ACCESS_TOKEN" required:"true"`
	APIEndpoint  string `envconfig:"API_ENDPOINT" required:"https://api-sandbox.dhl.com/parcel/de/shipping/v2/"` // capitalized and tagged
}

func InitShippingConfig() (shippingCfg ShippingConfig, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "staging" {
		if err = dotenv.Load(dotenv.EnvironmentFiles(env)); err != nil {
			return
		}
	}
	err = envconfig.Process("", &shippingCfg)
	return
}

func NewDHLClient(clientID, clientSecret, apiEndpoint string) DHLClient {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	return DHLClient{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
		apiEndpoint:  apiEndpoint,
	}
}

func (c *DHLClient) authenticate() error {
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return nil // Token is still valid
	}

	authURL := fmt.Sprintf("%s/oauth/token", c.apiEndpoint)
	resp, err := c.client.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		Post(authURL)

	if err != nil {
		return fmt.Errorf("error during authentication: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("authentication failed, status: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("error decoding authentication response: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return fmt.Errorf("no access_token in authentication response")
	}

	expiresIn, ok := result["expires_in"].(float64)
	if !ok {
		return fmt.Errorf("no expires_in in authentication response")
	}

	c.accessToken = token
	c.tokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	c.client.SetAuthToken(c.accessToken)
	return nil
}

func (c *DHLClient) CreateShipment(ctx context.Context, shipmentData interface{}) (DHLShipmentResponse, error) {
	if err := c.authenticate(); err != nil {
		return DHLShipmentResponse{}, fmt.Errorf("failed to authenticate DHL client: %w", err)
	}

	resp, err := c.client.R().
		SetBody(shipmentData).
		Post(fmt.Sprintf("%s/shipments", c.apiEndpoint))

	if err != nil {
		return DHLShipmentResponse{}, err
	}

	if resp.IsError() {
		return DHLShipmentResponse{}, fmt.Errorf("DHL API error: %s", resp.String())
	}

	var result DHLShipmentResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return DHLShipmentResponse{}, fmt.Errorf("failed to parse DHL response: %w", err)
	}

	return result, nil
}
