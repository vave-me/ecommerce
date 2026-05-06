package merchant

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/content/v2.1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type MerchantCenterClient struct {
	Service    *content.APIService
	merchantID uint64
}

func NewMerchantCenterClient(ctx context.Context, merchantID uint64, serviceAccountJSONPath string) (*MerchantCenterClient, error) {
	// If no service account path provided, create a disabled client
	if serviceAccountJSONPath == "" || serviceAccountJSONPath == "/path/to/some.json" {
		return &MerchantCenterClient{
			Service:    nil,
			merchantID: merchantID,
		}, nil
	}

	jsonKey, err := os.ReadFile(serviceAccountJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(jsonKey, content.ContentScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account JSON: %w", err)
	}

	httpClient := config.Client(ctx)
	svc, err := content.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create content API service: %w", err)
	}

	return &MerchantCenterClient{
		Service:    svc,
		merchantID: merchantID,
	}, nil
}

// Expose the lower-level methods for the application to use:
func (m *MerchantCenterClient) InsertProduct(ctx context.Context, prod *content.Product) error {
	if m.Service == nil {
		// Service disabled - no operation
		return nil
	}
	_, err := m.Service.Products.Insert(m.merchantID, prod).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("insert product failed: %w", err)
	}
	return nil
}

func (m *MerchantCenterClient) UpdateProduct(ctx context.Context, prod *content.Product) error {
	if m.Service == nil {
		// Service disabled - no operation
		return nil
	}
	_, err := m.Service.Products.Update(m.merchantID, prod.Id, prod).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("update product failed: %w", err)
	}
	return nil
}

func (m *MerchantCenterClient) GetProduct(ctx context.Context, productID string) (*content.Product, error) {
	if m.Service == nil {
		// Service disabled - return not found
		return nil, &googleapi.Error{
			Code:    404,
			Message: "Merchant service disabled - no service account configured",
		}
	}
	resp, err := m.Service.Products.Get(m.merchantID, productID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("get product failed: %w", err)
	}
	return resp, nil
}

func (m *MerchantCenterClient) DeleteProduct(ctx context.Context, productID string) error {
	if m.Service == nil {
		// Service disabled - no operation
		return nil
	}
	err := m.Service.Products.Delete(m.merchantID, productID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("delete product failed: %w", err)
	}
	return nil
}

func (m *MerchantCenterClient) ListProducts(ctx context.Context, pageSize int64, pageToken string) ([]*content.Product, string, error) {
	if m.Service == nil {
		// Service disabled - return empty list
		return []*content.Product{}, "", nil
	}
	call := m.Service.Products.List(m.merchantID)
	if pageSize > 0 {
		call.MaxResults(pageSize)
	}
	if pageToken != "" {
		call.PageToken(pageToken)
	}

	resp, err := call.Context(ctx).Do()
	if err != nil {
		return nil, "", fmt.Errorf("list products failed: %w", err)
	}
	return resp.Resources, resp.NextPageToken, nil
}

// We'll keep an isNotFoundErr helper to handle upsert logic in the application layer
func (m *MerchantCenterClient) IsNotFoundErr(err error) bool {
	var gErr *googleapi.Error
	if errors.As(err, &gErr) {
		return gErr.Code == http.StatusNotFound
	}
	return false
}

// MerchantID returns the merchant ID
func (m *MerchantCenterClient) MerchantID() uint64 {
	return m.merchantID
}

// IsEnabled returns true if the merchant service is configured and enabled
func (m *MerchantCenterClient) IsEnabled() bool {
	return m.Service != nil
}
