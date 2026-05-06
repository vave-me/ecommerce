package dhl

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a DHL Express MyDHL API client
type Client struct {
	username    string
	password    string
	apiURL      string
	accountNumber string
	httpClient  *http.Client
	isTest      bool
}

// NewClient creates a new DHL Express API client
func NewClient(username, password, apiURL, accountNumber string, isTest bool) *Client {
	if apiURL == "" {
		if isTest {
			apiURL = "https://express.api.dhl.com/mydhlapi/test"
		} else {
			apiURL = "https://express.api.dhl.com/mydhlapi"
		}
	}
	
	return &Client{
		username:    username,
		password:    password,
		apiURL:      apiURL,
		accountNumber: accountNumber,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		isTest: isTest,
	}
}

// CreateShipmentRequest represents DHL Express shipment creation request
type CreateShipmentRequest struct {
	PlannedShippingDateAndTime string                  `json:"plannedShippingDateAndTime"`
	Pickup                     *Pickup                 `json:"pickup,omitempty"`
	ProductCode               string                  `json:"productCode"`
	LocalProductCode          string                  `json:"localProductCode,omitempty"`
	GetRateEstimates          bool                    `json:"getRateEstimates"`
	Accounts                  []Account               `json:"accounts"`
	OutputImageProperties     *OutputImageProperties  `json:"outputImageProperties,omitempty"`
	CustomerDetails           CustomerDetails         `json:"customerDetails"`
	Content                   Content                 `json:"content"`
}

// Pickup represents pickup details
type Pickup struct {
	IsRequested bool   `json:"isRequested"`
	PickupDetails *PickupDetails `json:"pickupDetails,omitempty"`
}

// PickupDetails contains pickup scheduling information
type PickupDetails struct {
	PostalAddress            PostalAddress `json:"postalAddress"`
	ContactInformation       ContactInfo   `json:"contactInformation"`
	PickupRequestorDetails   ContactInfo   `json:"pickupRequestorDetails"`
	TypeCode                 string        `json:"typeCode"`
	PickupDateTime           string        `json:"pickupDateTime"`
	PickupDateTimeGMTOffset  string        `json:"pickupDateTimeGMTOffset"`
	LocationCloseTime        string        `json:"locationCloseTime"`
	SpecialInstructions      []ValueOnly   `json:"specialInstructions,omitempty"`
	Remark                   string        `json:"remark,omitempty"`
}

// Account represents DHL account information
type Account struct {
	TypeCode string `json:"typeCode"`
	Number   string `json:"number"`
}

// OutputImageProperties defines label output properties
type OutputImageProperties struct {
	PrinterDPI              int                      `json:"printerDPI,omitempty"`
	CustomerBarcodes        []CustomerBarcode        `json:"customerBarcodes,omitempty"`
	CustomerLogos           []CustomerLogo           `json:"customerLogos,omitempty"`
	EncoderType             string                   `json:"encoderType,omitempty"`
	ImageOptions            *ImageOptions            `json:"imageOptions,omitempty"`
}

// CustomerBarcode represents custom barcode on label
type CustomerBarcode struct {
	Content  string `json:"content"`
	TextBelowBarcode string `json:"textBelowBarcode"`
	SymbologyCode string `json:"symbologyCode"`
}

// CustomerLogo represents custom logo on label
type CustomerLogo struct {
	FileFormat string `json:"fileFormat"`
	Content    string `json:"content"`
}

// ImageOptions for label generation
type ImageOptions struct {
	TypeCode                 string `json:"typeCode"`
	TemplateName             string `json:"templateName,omitempty"`
	IsRequested              bool   `json:"isRequested"`
	HideAccountNumber        bool   `json:"hideAccountNumber,omitempty"`
	NumberOfCopies           int    `json:"numberOfCopies,omitempty"`
}

// CustomerDetails contains shipper and receiver information
type CustomerDetails struct {
	ShipperDetails  ShipperReceiverDetails `json:"shipperDetails"`
	ReceiverDetails ShipperReceiverDetails `json:"receiverDetails"`
	BuyerDetails    *ShipperReceiverDetails `json:"buyerDetails,omitempty"`
}

// ShipperReceiverDetails represents contact and address details
type ShipperReceiverDetails struct {
	PostalAddress       PostalAddress `json:"postalAddress"`
	ContactInformation  ContactInfo   `json:"contactInformation"`
	RegistrationNumbers []RegistrationNumber `json:"registrationNumbers,omitempty"`
	BankDetails         []BankDetail  `json:"bankDetails,omitempty"`
	TypeCode            string        `json:"typeCode,omitempty"`
}

// PostalAddress represents a physical address
type PostalAddress struct {
	PostalCode   string `json:"postalCode,omitempty"`
	CityName     string `json:"cityName"`
	CountryCode  string `json:"countryCode"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	CountryName  string `json:"countryName,omitempty"`
	StateOrProvinceCode string `json:"stateOrProvinceCode,omitempty"`
}

// ContactInfo represents contact information
type ContactInfo struct {
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone"`
	MobilePhone string `json:"mobilePhone,omitempty"`
	CompanyName string `json:"companyName"`
	FullName    string `json:"fullName"`
}

// RegistrationNumber for tax/business registration
type RegistrationNumber struct {
	TypeCode string `json:"typeCode"`
	Number   string `json:"number"`
	IssuerCountryCode string `json:"issuerCountryCode"`
}

// BankDetail represents bank account information
type BankDetail struct {
	Name        string `json:"name"`
	SettlementLocalCurrency string `json:"settlementLocalCurrency"`
	SettlementForeignCurrency string `json:"settlementForeignCurrency,omitempty"`
}

// Content represents shipment content details
type Content struct {
	Packages               []Package              `json:"packages"`
	IsCustomsDeclarable    bool                   `json:"isCustomsDeclarable"`
	DeclaredValue          float64                `json:"declaredValue,omitempty"`
	DeclaredValueCurrency  string                 `json:"declaredValueCurrency,omitempty"`
	ExportDeclaration      *ExportDeclaration     `json:"exportDeclaration,omitempty"`
	Description            string                 `json:"description"`
	USFilingTypeValue      string                 `json:"USFilingTypeValue,omitempty"`
	Incoterm               string                 `json:"incoterm,omitempty"`
	UnitOfMeasurement      string                 `json:"unitOfMeasurement"`
}

// Package represents a single package in the shipment
type Package struct {
	TypeCode            string               `json:"typeCode,omitempty"`
	Weight              float64              `json:"weight"`
	Dimensions          Dimensions           `json:"dimensions"`
	CustomerReferences  []CustomerReference  `json:"customerReferences,omitempty"`
	Identifiers         []PackageIdentifier  `json:"identifiers,omitempty"`
	Description         string               `json:"description,omitempty"`
	LabelBarcodes       []LabelBarcode       `json:"labelBarcodes,omitempty"`
	LabelText           []LabelText          `json:"labelText,omitempty"`
	LabelDescription    string               `json:"labelDescription,omitempty"`
}

// Dimensions represents package dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// CustomerReference for package tracking
type CustomerReference struct {
	Value    string `json:"value"`
	TypeCode string `json:"typeCode,omitempty"`
}

// PackageIdentifier for unique package identification
type PackageIdentifier struct {
	TypeCode    string `json:"typeCode"`
	Value       string `json:"value"`
	DataIdentifier string `json:"dataIdentifier,omitempty"`
}

// LabelBarcode for custom barcodes on package label
type LabelBarcode struct {
	Position        string `json:"position"`
	SymbologyCode   string `json:"symbologyCode"`
	Content         string `json:"content"`
	TextBelowBarcode string `json:"textBelowBarcode,omitempty"`
}

// LabelText for custom text on package label
type LabelText struct {
	Position string `json:"position"`
	Caption  string `json:"caption"`
	Value    string `json:"value"`
}

// ExportDeclaration for customs
type ExportDeclaration struct {
	LineItems                    []LineItem    `json:"lineItems,omitempty"`
	Invoice                      *Invoice      `json:"invoice,omitempty"`
	Remarks                      []ValueOnly   `json:"remarks,omitempty"`
	AdditionalCharges            []AdditionalCharge `json:"additionalCharges,omitempty"`
	DestinationPortName          string        `json:"destinationPortName,omitempty"`
	PlaceOfIncoterm              string        `json:"placeOfIncoterm,omitempty"`
	PayerVATNumber               string        `json:"payerVATNumber,omitempty"`
	RecipientReference           string        `json:"recipientReference,omitempty"`
	Exporter                     *Exporter     `json:"exporter,omitempty"`
	PackageMarks                 string        `json:"packageMarks,omitempty"`
	ExportReason                 string        `json:"exportReason,omitempty"`
	ExportReasonType             string        `json:"exportReasonType,omitempty"`
	Licenses                     []License     `json:"licenses,omitempty"`
	ShipmentType                 string        `json:"shipmentType,omitempty"`
	CustomsDocuments             []CustomsDocument `json:"customsDocuments,omitempty"`
}

// LineItem for export declaration
type LineItem struct {
	Number               int                  `json:"number"`
	Description          string               `json:"description"`
	Price                float64              `json:"price"`
	Quantity             *Quantity            `json:"quantity"`
	CommodityCodes       []CommodityCode      `json:"commodityCodes,omitempty"`
	ExportReasonType     string               `json:"exportReasonType,omitempty"`
	ManufacturerCountry  string               `json:"manufacturerCountry"`
	ExportControlClassificationNumber string    `json:"exportControlClassificationNumber,omitempty"`
	Weight               *Weight              `json:"weight"`
	IsTaxesPaid          bool                 `json:"isTaxesPaid,omitempty"`
	CustomerReferences   []CustomerReference  `json:"customerReferences,omitempty"`
	CustomsDocuments     []CustomsDocument    `json:"customsDocuments,omitempty"`
}

// ValueOnly for simple value fields
type ValueOnly struct {
	Value string `json:"value"`
}

// AdditionalCharge for customs
type AdditionalCharge struct {
	Value    float64 `json:"value"`
	Caption  string  `json:"caption,omitempty"`
	TypeCode string  `json:"typeCode,omitempty"`
}

// Invoice details
type Invoice struct {
	Number          string    `json:"number"`
	Date            string    `json:"date"`
	SignatureName   string    `json:"signatureName,omitempty"`
	SignatureTitle  string    `json:"signatureTitle,omitempty"`
	SignatureImage  string    `json:"signatureImage,omitempty"`
	Instructions    []string  `json:"instructions,omitempty"`
	TotalNetWeight  float64   `json:"totalNetWeight,omitempty"`
	TotalGrossWeight float64  `json:"totalGrossWeight,omitempty"`
	CustomerReferences []CustomerReference `json:"customerReferences,omitempty"`
	TermsOfPayment  string    `json:"termsOfPayment,omitempty"`
}

// Exporter information
type Exporter struct {
	ID               string           `json:"id,omitempty"`
	Code             string           `json:"code,omitempty"`
}

// License for export
type License struct {
	TypeCode string `json:"typeCode"`
	Value    string `json:"value"`
}

// CustomsDocument reference
type CustomsDocument struct {
	TypeCode string `json:"typeCode"`
	Value    string `json:"value"`
}

// Quantity information
type Quantity struct {
	Value             int    `json:"value"`
	UnitOfMeasurement string `json:"unitOfMeasurement"`
}

// CommodityCode for customs
type CommodityCode struct {
	TypeCode string `json:"typeCode"`
	Value    string `json:"value"`
}

// Weight information
type Weight struct {
	NetValue   float64 `json:"netValue"`
	GrossValue float64 `json:"grossValue"`
}

// CreateShipmentResponse from DHL API
type CreateShipmentResponse struct {
	URL              string            `json:"url,omitempty"`
	ShipmentTrackingNumber string       `json:"shipmentTrackingNumber"`
	CancelPickupURL  string            `json:"cancelPickupUrl,omitempty"`
	TrackingURL      string            `json:"trackingUrl"`
	DispatchConfirmationNumber string   `json:"dispatchConfirmationNumber,omitempty"`
	Packages         []PackageResult   `json:"packages,omitempty"`
	Documents        []Document        `json:"documents,omitempty"`
	OnDemandDeliveryURL string         `json:"onDemandDeliveryURL,omitempty"`
	ShipmentDetails  []ShipmentDetail  `json:"shipmentDetails,omitempty"`
	EstimatedDeliveryDate *EstimatedDeliveryDate `json:"estimatedDeliveryDate,omitempty"`
	Warnings         []string          `json:"warnings,omitempty"`
}

// PackageResult for created package
type PackageResult struct {
	ReferenceNumber int        `json:"referenceNumber"`
	TrackingNumber  string     `json:"trackingNumber"`
	TrackingURL     string     `json:"trackingUrl"`
	Identifiers     []PackageIdentifier `json:"identifiers,omitempty"`
	Documents       []Document `json:"documents,omitempty"`
}

// Document represents a shipping document
type Document struct {
	ImageFormat string `json:"imageFormat"`
	Content     string `json:"content"`
	TypeCode    string `json:"typeCode,omitempty"`
	PDF         *PDFDocument `json:"pdf,omitempty"`
}

// PDFDocument with rendering options
type PDFDocument struct {
	PrinterDPI int    `json:"printerDPI,omitempty"`
	Content    string `json:"content"`
}

// ShipmentDetail with charges
type ShipmentDetail struct {
	ServiceHandlingFeatureCodes []string  `json:"serviceHandlingFeatureCodes,omitempty"`
	VolumetricWeight           float64   `json:"volumetricWeight,omitempty"`
	BillingCode                string    `json:"billingCode,omitempty"`
	ServiceContentCode         string    `json:"serviceContentCode,omitempty"`
	CustomerDetails            *CustomerDetails `json:"customerDetails,omitempty"`
}

// EstimatedDeliveryDate information
type EstimatedDeliveryDate struct {
	EstimatedDeliveryDate string `json:"estimatedDeliveryDate,omitempty"`
	EstimatedDeliveryType string `json:"estimatedDeliveryType,omitempty"`
}

// TrackingRequest for shipment tracking
type TrackingRequest struct {
	TrackingNumber       string `json:"trackingNumber"`
	Service              string `json:"service,omitempty"`
	RequesterCountryCode string `json:"requesterCountryCode,omitempty"`
	OriginCountryCode    string `json:"originCountryCode,omitempty"`
	DestinationCountryCode string `json:"destinationCountryCode,omitempty"`
}

// TrackingResponse from DHL API
type TrackingResponse struct {
	Shipments []TrackedShipment `json:"shipments"`
}

// TrackedShipment information
type TrackedShipment struct {
	ShipmentTrackingNumber string           `json:"shipmentTrackingNumber"`
	Status                 string           `json:"status"`
	ShipmentTimestamp      string           `json:"shipmentTimestamp"`
	ProductCode            string           `json:"productCode"`
	Description            string           `json:"description"`
	ShipperDetails         ShipperReceiverDetails `json:"shipperDetails"`
	ReceiverDetails        ShipperReceiverDetails `json:"receiverDetails"`
	TotalWeight            float64          `json:"totalWeight"`
	UnitOfMeasurements     string           `json:"unitOfMeasurements"`
	EstimatedDeliveryDate  string           `json:"estimatedDeliveryDate,omitempty"`
	Events                 []Event          `json:"events"`
	NumberOfPieces         int              `json:"numberOfPieces"`
	PieceIds               []string         `json:"pieceIds,omitempty"`
	OriginServiceArea      ServiceArea      `json:"originServiceArea"`
	DestinationServiceArea ServiceArea      `json:"destinationServiceArea"`
}

// Event in shipment tracking
type Event struct {
	Date           string      `json:"date"`
	Time           string      `json:"time"`
	TypeCode       string      `json:"typeCode"`
	Description    string      `json:"description"`
	ServiceArea    ServiceArea `json:"serviceArea,omitempty"`
	SignedBy       string      `json:"signedBy,omitempty"`
}

// ServiceArea information
type ServiceArea struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// RateRequest for calculating shipping rates
type RateRequest struct {
	CustomerDetails            CustomerDetails  `json:"customerDetails"`
	Accounts                   []Account        `json:"accounts"`
	ProductCode                string           `json:"productCode,omitempty"`
	LocalProductCode           string           `json:"localProductCode,omitempty"`
	ValueAddedServices         []ValueAddedService `json:"valueAddedServices,omitempty"`
	ProductsAndServices        []ProductsAndServices `json:"productsAndServices,omitempty"`
	PayerCountryCode           string           `json:"payerCountryCode,omitempty"`
	PlannedShippingDateAndTime string           `json:"plannedShippingDateAndTime"`
	UnitOfMeasurement          string           `json:"unitOfMeasurement"`
	IsCustomsDeclarable        bool             `json:"isCustomsDeclarable"`
	MonetaryAmount             []MonetaryAmount `json:"monetaryAmount,omitempty"`
	RequestAllValueAddedServices bool           `json:"requestAllValueAddedServices,omitempty"`
	EstimatedDeliveryDate      *EstimatedDeliveryDateRequest `json:"estimatedDeliveryDate,omitempty"`
	GetAdditionalInformation   []AdditionalInformationRequest `json:"getAdditionalInformation,omitempty"`
	ReturnStandardProductsOnly bool             `json:"returnStandardProductsOnly,omitempty"`
	NextBusinessDay            bool             `json:"nextBusinessDay,omitempty"`
	ProductTypeCode            string           `json:"productTypeCode,omitempty"`
	Packages                   []Package        `json:"packages"`
}

// ValueAddedService for rate calculation
type ValueAddedService struct {
	ServiceCode string `json:"serviceCode"`
	Value       float64 `json:"value,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

// ProductsAndServices filter
type ProductsAndServices struct {
	TypeCode string `json:"typeCode"`
	ServiceCode string `json:"serviceCode,omitempty"`
	LocalServiceCode string `json:"localServiceCode,omitempty"`
}

// MonetaryAmount for rate calculation
type MonetaryAmount struct {
	TypeCode string  `json:"typeCode"`
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// EstimatedDeliveryDateRequest options
type EstimatedDeliveryDateRequest struct {
	IsRequested bool   `json:"isRequested"`
	TypeCode    string `json:"typeCode,omitempty"`
}

// AdditionalInformationRequest types
type AdditionalInformationRequest struct {
	TypeCode    string `json:"typeCode"`
	IsRequested bool   `json:"isRequested"`
}

// RateResponse from DHL API
type RateResponse struct {
	Products []Product `json:"products"`
}

// Product with rate information
type Product struct {
	ProductName         string              `json:"productName"`
	ProductCode         string              `json:"productCode"`
	LocalProductCode    string              `json:"localProductCode,omitempty"`
	LocalProductCountryCode string          `json:"localProductCountryCode,omitempty"`
	NetworkTypeCode     string              `json:"networkTypeCode,omitempty"`
	IsCustomerAgreement bool                `json:"isCustomerAgreement,omitempty"`
	Weight              Weight              `json:"weight"`
	TotalPrice          []TotalPrice        `json:"totalPrice"`
	TotalPriceBreakdown []PriceBreakdown    `json:"totalPriceBreakdown,omitempty"`
	DetailedPriceBreakdown []DetailedPriceBreakdown `json:"detailedPriceBreakdown,omitempty"`
	ServiceCodeMutuallyExclusiveGroups []MutuallyExclusiveGroup `json:"serviceCodeMutuallyExclusiveGroups,omitempty"`
	ServiceCodeDependencyRuleGroups []DependencyRuleGroup `json:"serviceCodeDependencyRuleGroups,omitempty"`
	PickupCapabilities  *PickupCapabilities `json:"pickupCapabilities,omitempty"`
	DeliveryCapabilities *DeliveryCapabilities `json:"deliveryCapabilities,omitempty"`
}

// TotalPrice information
type TotalPrice struct {
	CurrencyType string  `json:"currencyType,omitempty"`
	PriceCurrency string `json:"priceCurrency"`
	Price        float64 `json:"price"`
}

// PriceBreakdown for rates
type PriceBreakdown struct {
	PriceType    string      `json:"priceType,omitempty"`
	TypeCode     string      `json:"typeCode"`
	Price        float64     `json:"price"`
	Rate         float64     `json:"rate,omitempty"`
	BasePrice    float64     `json:"basePrice,omitempty"`
}

// DetailedPriceBreakdown with more detail
type DetailedPriceBreakdown struct {
	Breakdown []PriceBreakdown `json:"breakdown"`
}

// MutuallyExclusiveGroup of services
type MutuallyExclusiveGroup struct {
	ServiceCodes []string `json:"serviceCodes"`
	Description  string   `json:"description,omitempty"`
}

// DependencyRuleGroup for services
type DependencyRuleGroup struct {
	DependentServiceCode string   `json:"dependentServiceCode"`
	DependencyRuleGroup  []DependencyRule `json:"dependencyRuleGroup"`
}

// DependencyRule for services
type DependencyRule struct {
	RequiredServiceCode string `json:"requiredServiceCode"`
	MinRequiredNumberOfOptions int `json:"minRequiredNumberOfOptions,omitempty"`
}

// PickupCapabilities information
type PickupCapabilities struct {
	NextBusinessDay          bool     `json:"nextBusinessDay,omitempty"`
	LocalCutoffDateAndTime   string   `json:"localCutoffDateAndTime,omitempty"`
	GMTCutoffTime            string   `json:"GMTCutoffTime,omitempty"`
	PickupEarliest           string   `json:"pickupEarliest,omitempty"`
	PickupLatest             string   `json:"pickupLatest,omitempty"`
	OriginServiceAreaCode    string   `json:"originServiceAreaCode,omitempty"`
	OriginFacilityAreaCode   string   `json:"originFacilityAreaCode,omitempty"`
	PickupAdditionalDays     int      `json:"pickupAdditionalDays,omitempty"`
	PickupDayOfWeek          int      `json:"pickupDayOfWeek,omitempty"`
}

// DeliveryCapabilities information
type DeliveryCapabilities struct {
	EstimatedDeliveryDateAndTime string `json:"estimatedDeliveryDateAndTime,omitempty"`
	DestinationServiceAreaCode   string `json:"destinationServiceAreaCode,omitempty"`
	DestinationFacilityAreaCode  string `json:"destinationFacilityAreaCode,omitempty"`
	DeliveryAdditionalDays       int    `json:"deliveryAdditionalDays,omitempty"`
	DeliveryDayOfWeek            int    `json:"deliveryDayOfWeek,omitempty"`
	TotalTransitDays             int    `json:"totalTransitDays,omitempty"`
}

// CreateShipment creates a new shipment
func (c *Client) CreateShipment(ctx context.Context, data map[string]interface{}) (*CreateShipmentResponse, error) {
	// For backward compatibility, convert the map to proper request structure
	// In production, you should use the CreateShipmentRequest struct directly
	
	req := CreateShipmentRequest{
		PlannedShippingDateAndTime: time.Now().Format("2006-01-02T15:04:05 MST"),
		ProductCode: "P", // DHL Express Worldwide
		GetRateEstimates: false,
		Accounts: []Account{{
			TypeCode: "shipper",
			Number:   c.accountNumber,
		}},
		OutputImageProperties: &OutputImageProperties{
			ImageOptions: &ImageOptions{
				TypeCode: "PDF",
				IsRequested: true,
			},
		},
		CustomerDetails: CustomerDetails{
			ShipperDetails: ShipperReceiverDetails{
				PostalAddress: PostalAddress{
					AddressLine1: data["sender"].(map[string]interface{})["address"].(string),
					CityName:     "New York",
					PostalCode:   "10001",
					CountryCode:  "US",
				},
				ContactInformation: ContactInfo{
					Email:       "redacted-email@example.com",
					Phone:       "+1234567890",
					CompanyName: data["sender"].(map[string]interface{})["name"].(string),
					FullName:    data["sender"].(map[string]interface{})["name"].(string),
				},
			},
			ReceiverDetails: ShipperReceiverDetails{
				PostalAddress: PostalAddress{
					AddressLine1: data["receiver"].(map[string]interface{})["address"].(string),
					CityName:     "Los Angeles",
					PostalCode:   "90001",
					CountryCode:  "US",
				},
				ContactInformation: ContactInfo{
					Email:       "redacted-email@example.com",
					Phone:       "+1234567890",
					CompanyName: data["receiver"].(map[string]interface{})["name"].(string),
					FullName:    data["receiver"].(map[string]interface{})["name"].(string),
				},
			},
		},
		Content: Content{
			Packages: []Package{{
				Weight: parseWeight(data["package"].(map[string]interface{})["weight"].(string)),
				Dimensions: Dimensions{
					Length: data["package"].(map[string]interface{})["dimensions"].(map[string]float64)["length"],
					Width:  data["package"].(map[string]interface{})["dimensions"].(map[string]float64)["width"],
					Height: data["package"].(map[string]interface{})["dimensions"].(map[string]float64)["height"],
				},
			}},
			IsCustomsDeclarable: false,
			Description:         "Shipment",
			UnitOfMeasurement:   "metric",
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"/shipments", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var shipmentResp CreateShipmentResponse
	if err := json.Unmarshal(respBody, &shipmentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// For development/testing, generate mock data if response is empty
	if c.isTest && shipmentResp.ShipmentTrackingNumber == "" {
		shipmentResp = CreateShipmentResponse{
			ShipmentTrackingNumber: fmt.Sprintf("1234567890%d", time.Now().Unix()%1000),
			TrackingURL: "https://www.dhl.com/track",
			Documents: []Document{{
				ImageFormat: "PDF",
				Content:     base64.StdEncoding.EncodeToString([]byte("Mock PDF Label")),
				TypeCode:    "label",
			}},
			EstimatedDeliveryDate: &EstimatedDeliveryDate{
				EstimatedDeliveryDate: time.Now().Add(3 * 24 * time.Hour).Format("2006-01-02"),
				EstimatedDeliveryType: "QDDC",
			},
		}
	}

	return &shipmentResp, nil
}

// TrackShipment retrieves tracking information
func (c *Client) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("%s/shipments/%s/tracking", c.apiURL, trackingNumber), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var trackingResp TrackingResponse
	if err := json.Unmarshal(respBody, &trackingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// For development/testing, generate mock data if response is empty
	if c.isTest && len(trackingResp.Shipments) == 0 {
		trackingResp = TrackingResponse{
			Shipments: []TrackedShipment{{
				ShipmentTrackingNumber: trackingNumber,
				Status: "transit",
				ShipmentTimestamp: time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05"),
				ProductCode: "P",
				Description: "EXPRESS WORLDWIDE",
				TotalWeight: 1.0,
				UnitOfMeasurements: "metric",
				EstimatedDeliveryDate: time.Now().Add(2 * 24 * time.Hour).Format("2006-01-02"),
				Events: []Event{
					{
						Date:        time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
						Time:        "10:00:00",
						TypeCode:    "PU",
						Description: "Shipment picked up",
						ServiceArea: ServiceArea{
							Code:        "NYC",
							Description: "New York City",
						},
					},
					{
						Date:        time.Now().Format("2006-01-02"),
						Time:        "14:30:00",
						TypeCode:    "PL",
						Description: "Processed at DHL facility",
						ServiceArea: ServiceArea{
							Code:        "ORD",
							Description: "Chicago",
						},
					},
				},
				NumberOfPieces: 1,
			}},
		}
	}

	return &trackingResp, nil
}

// GetRates calculates shipping rates
func (c *Client) GetRates(ctx context.Context, req RateRequest) (*RateResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"/rates", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var rateResp RateResponse
	if err := json.Unmarshal(respBody, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &rateResp, nil
}

// CancelShipment cancels a shipment
func (c *Client) CancelShipment(ctx context.Context, trackingNumber string, reason string) error {
	req := map[string]string{
		"reason": reason,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", 
		fmt.Sprintf("%s/shipments/%s", c.apiURL, trackingNumber), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GetLabel retrieves shipping label
func (c *Client) GetLabel(ctx context.Context, trackingNumber string) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("%s/shipments/%s/get-image?pickupYearAndMonth=%s&encodingFormat=pdf&allInOnePDF=false&compressedPackage=false", 
			c.apiURL, trackingNumber, time.Now().Format("2006-01")), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// SchedulePickup schedules a pickup for shipments
func (c *Client) SchedulePickup(ctx context.Context, req PickupRequest) (*PickupResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"/pickups", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var pickupResp PickupResponse
	if err := json.Unmarshal(respBody, &pickupResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &pickupResp, nil
}

// PickupRequest represents a pickup scheduling request
type PickupRequest struct {
	PlannedPickupDateAndTime string                  `json:"plannedPickupDateAndTime"`
	Accounts                 []Account               `json:"accounts"`
	CustomerDetails          PickupCustomerDetails   `json:"customerDetails"`
	ShipmentDetails          []PickupShipmentDetail  `json:"shipmentDetails"`
	CloseTime                string                  `json:"closeTime,omitempty"`
	Location                 string                  `json:"location,omitempty"`
	LocationType             string                  `json:"locationType,omitempty"`
	SpecialInstructions      []ValueOnly             `json:"specialInstructions,omitempty"`
	Remark                   string                  `json:"remark,omitempty"`
}

// PickupCustomerDetails for pickup request
type PickupCustomerDetails struct {
	ShipperDetails   ShipperReceiverDetails `json:"shipperDetails"`
	ReceiverDetails  ShipperReceiverDetails `json:"receiverDetails"`
	BookingRequestorDetails *ShipperReceiverDetails `json:"bookingRequestorDetails,omitempty"`
	PickupDetails    *ShipperReceiverDetails `json:"pickupDetails,omitempty"`
}

// PickupShipmentDetail for pickup
type PickupShipmentDetail struct {
	ProductCode              string    `json:"productCode"`
	LocalProductCode         string    `json:"localProductCode,omitempty"`
	Accounts                 []Account `json:"accounts,omitempty"`
	ValueAddedServices       []ValueAddedService `json:"valueAddedServices,omitempty"`
	IsCustomsDeclarable      bool      `json:"isCustomsDeclarable"`
	DeclaredValue            float64   `json:"declaredValue,omitempty"`
	DeclaredValueCurrency    string    `json:"declaredValueCurrency,omitempty"`
	UnitOfMeasurement        string    `json:"unitOfMeasurement"`
	ShipmentTrackingNumber   string    `json:"shipmentTrackingNumber,omitempty"`
	Packages                 []Package `json:"packages"`
}

// PickupResponse from DHL API
type PickupResponse struct {
	DispatchConfirmationNumbers []string `json:"dispatchConfirmationNumbers"`
	ReadyByTime                 string   `json:"readyByTime,omitempty"`
	NextPickupDate              string   `json:"nextPickupDate,omitempty"`
	Warnings                    []string `json:"warnings,omitempty"`
}

// setHeaders sets common headers for DHL API requests
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Message-Reference", fmt.Sprintf("MSG_%d", time.Now().Unix()))
	req.Header.Set("Message-Reference-Date", time.Now().Format("2006-01-02T15:04:05-07:00"))
}

// parseWeight converts string weight to float64
func parseWeight(weight string) float64 {
	var w float64
	fmt.Sscanf(weight, "%f", &w)
	if w == 0 {
		w = 1.0 // Default weight
	}
	return w
}

// Helper function to create a simple shipment (convenience method)
func (c *Client) CreateSimpleShipment(ctx context.Context, shipper, receiver ShipperReceiverDetails, 
	packages []Package, serviceCode string) (*CreateShipmentResponse, error) {
	
	req := CreateShipmentRequest{
		PlannedShippingDateAndTime: time.Now().Format("2006-01-02T15:04:05 MST"),
		ProductCode: serviceCode,
		GetRateEstimates: false,
		Accounts: []Account{{
			TypeCode: "shipper",
			Number:   c.accountNumber,
		}},
		OutputImageProperties: &OutputImageProperties{
			ImageOptions: &ImageOptions{
				TypeCode: "PDF",
				IsRequested: true,
			},
		},
		CustomerDetails: CustomerDetails{
			ShipperDetails:  shipper,
			ReceiverDetails: receiver,
		},
		Content: Content{
			Packages:            packages,
			IsCustomsDeclarable: false,
			Description:         "Shipment",
			UnitOfMeasurement:   "metric",
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"/shipments", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var shipmentResp CreateShipmentResponse
	if err := json.Unmarshal(respBody, &shipmentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &shipmentResp, nil
}