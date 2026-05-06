package idnow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// IDnowClient encapsulates the IDnow API interactions.
type IDnowClient struct {
	APIToken string
	BaseURL  string
	Customer string
	Client   *http.Client
}

// NewIDnowClient initializes a new IDnowClient with your IDnow credentials.
func NewIDnowClient(apiToken, baseURL, customer string) *IDnowClient {
	return &IDnowClient{
		APIToken: apiToken,
		BaseURL:  baseURL,
		Customer: customer,
		Client:   &http.Client{},
	}
}

// -----------------------------------------------------------------------------
// 1) ARCHIVE IDENT
// -----------------------------------------------------------------------------

// ArchiveIdent archives an IDnow Ident by its transaction number.
//
// Endpoint: POST /api/v1/{customer}/identifications/{transactionNumber}/archive
// This call does not return a body on success (2xx). If the ident is already archived or deleted,
// the server may return an error.
func (c *IDnowClient) ArchiveIdent(ctx context.Context, transactionNumber string) error {
	url := fmt.Sprintf("%s/api/v1/%s/identifications/%s/archive",
		c.BaseURL, c.Customer, transactionNumber,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("creating archive request: %w", err)
	}

	req.Header.Set("X-API-LOGIN-TOKEN", c.APIToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("archive ident request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("archive ident failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// -----------------------------------------------------------------------------
// 2) CREATE DOCUMENT DEFINITION
// -----------------------------------------------------------------------------

// CreateDocumentDefinition creates a "document definition" on the IDnow system.
// Endpoint: POST /api/v1/{customer}/documentdefinitions
//
// Pass in "name", "identifier", "mimeType" in the request body.
// Returns no body on success (201). If not 2xx, returns an error.
func (c *IDnowClient) CreateDocumentDefinition(ctx context.Context, name, identifier, mimeType string) error {
	url := fmt.Sprintf("%s/api/v1/%s/documentdefinitions", c.BaseURL, c.Customer)

	payload := map[string]string{
		"name":       name,
		"identifier": identifier,
		"mimeType":   mimeType,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal createDocumentDefinition body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("creating createDocumentDefinition request: %w", err)
	}

	req.Header.Set("X-API-LOGIN-TOKEN", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("createDocumentDefinition request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("createDocumentDefinition failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// -----------------------------------------------------------------------------
// 3) GET IDENT
// -----------------------------------------------------------------------------

// IdentData represents a simplified structure of IDnow ident data as returned by
// GET /api/v1/{customer}/identifications/{transactionNumber}.
type IdentData struct {
	Result         string `json:"result,omitempty"`
	TransactionNbr string `json:"transactionnumber,omitempty"`
	// Add more fields as needed from IDnow's JSON response
}

// GetIdent retrieves details about an existing ident.
//
// Endpoint: GET /api/v1/{customer}/identifications/{transactionNumber}
func (c *IDnowClient) GetIdent(ctx context.Context, transactionNumber string) (*IdentData, error) {
	url := fmt.Sprintf("%s/api/v1/%s/identifications/%s", c.BaseURL, c.Customer, transactionNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating getIdent request: %w", err)
	}
	req.Header.Set("X-API-LOGIN-TOKEN", c.APIToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getIdent request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("getIdent failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var ident IdentData
	if err := json.NewDecoder(resp.Body).Decode(&ident); err != nil {
		return nil, fmt.Errorf("failed to decode getIdent response: %w", err)
	}

	return &ident, nil
}

// -----------------------------------------------------------------------------
// 4) CREATE AN IDENT FOR VIDEOIDENT, EID, OR ESIGN
// -----------------------------------------------------------------------------

// The following three methods demonstrate how you might create a new ident for
//  - VideoIdent ("VIDEO_IDENT"),
//  - eID ("EID"),
//  - eSign ("ESIGN")
// by calling the same endpoint:
// POST /api/v1/{customer}/identifications/{transactionNumber}/start
//
// These methods differ in the "profile" or parameters you send.

// CreateIdentRequest holds typical fields IDnow might expect
// when creating a new ident. The example below includes a "profile" field
// for eID, VideoIdent, or eSign usage, plus user data like firstname, lastname, etc.
type CreateIdentRequest struct {
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	Email       string `json:"email,omitempty"`
	Birthday    string `json:"birthday,omitempty"`
	Birthplace  string `json:"birthplace,omitempty"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	MobilePhone string `json:"mobilephone,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Custom1     string `json:"custom1,omitempty"`
	Custom2     string `json:"custom2,omitempty"`
	Custom3     string `json:"custom3,omitempty"`
	Profile     string `json:"profile,omitempty"` // e.g. "VIDEO_IDENT", "EID", "ESIGN"
	// etc. per your actual IDnow contract
}

type CreateIdentResponse struct {
	ID string `json:"id,omitempty"`
	// e.g. "TST-FXWF" or something assigned by IDnow
	// you can add more fields from IDnow's response
}

// CreateVideoIdent calls POST /api/v1/{customer}/identifications/{transactionNumber}/start
// with "profile = VIDEO_IDENT" plus user data.
func (c *IDnowClient) CreateVideoIdent(ctx context.Context, transactionNumber string, reqData CreateIdentRequest) (*CreateIdentResponse, error) {
	reqData.Profile = "VIDEO_IDENT"
	return c.createIdent(ctx, transactionNumber, reqData)
}

// CreateEIDIdent calls the same endpoint with "profile = EID".
func (c *IDnowClient) CreateEIDIdent(ctx context.Context, transactionNumber string, reqData CreateIdentRequest) (*CreateIdentResponse, error) {
	reqData.Profile = "EID"
	return c.createIdent(ctx, transactionNumber, reqData)
}

// CreateESignIdent calls the same endpoint with "profile = ESIGN".
func (c *IDnowClient) CreateESignIdent(ctx context.Context, transactionNumber string, reqData CreateIdentRequest) (*CreateIdentResponse, error) {
	reqData.Profile = "ESIGN"
	return c.createIdent(ctx, transactionNumber, reqData)
}

// createIdent is an internal helper that calls
// POST /api/v1/{customer}/identifications/{transactionNumber}/start with the request data
func (c *IDnowClient) createIdent(ctx context.Context, transactionNumber string, reqData CreateIdentRequest) (*CreateIdentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/%s/identifications/%s/start",
		c.BaseURL, c.Customer, transactionNumber,
	)

	jsonBytes, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal createIdent request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create ident request: %w", err)
	}

	req.Header.Set("X-API-LOGIN-TOKEN", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create ident request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create ident failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var result CreateIdentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode create ident response failed: %w", err)
	}

	return &result, nil
}
