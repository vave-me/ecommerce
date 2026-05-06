package dynamic365

import (
	"time"
)

// D365Notification represents a webhook notification from Dynamics 365
type D365Notification struct {
	ValidationToken string                   `json:"validationToken,omitempty"`
	Value           []D365ChangeNotification `json:"value"`
}

// D365ChangeNotification represents a single change notification
type D365ChangeNotification struct {
	ID                             string                 `json:"id"`
	SubscriptionID                 string                 `json:"subscriptionId"`
	SubscriptionExpirationDateTime time.Time              `json:"subscriptionExpirationDateTime"`
	ChangeType                     string                 `json:"changeType"`
	ClientState                    string                 `json:"clientState"`
	Resource                       string                 `json:"resource"`
	ResourceData                   map[string]interface{} `json:"resourceData,omitempty"`
	TenantID                       string                 `json:"tenantId"`
}

// D365Company represents a company in Business Central
type D365Company struct {
	ID                string `json:"id"`
	SystemVersion     string `json:"systemVersion"`
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	BusinessProfileID string `json:"businessProfileId"`
}

// D365Item represents an item/product in Business Central
type D365Item struct {
	ID                   string    `json:"id"`
	Number               string    `json:"number"`
	DisplayName          string    `json:"displayName"`
	Type                 string    `json:"type"`
	ItemCategoryID       string    `json:"itemCategoryId"`
	ItemCategoryCode     string    `json:"itemCategoryCode"`
	Blocked              bool      `json:"blocked"`
	BaseUnitOfMeasure    string    `json:"baseUnitOfMeasure"`
	GTIN                 string    `json:"gtin"`
	Description          string    `json:"description"`
	UnitPrice            float64   `json:"unitPrice"`
	UnitCost             float64   `json:"unitCost"`
	Inventory            float64   `json:"inventory"`
	GrossWeight          float64   `json:"grossWeight"`
	NetWeight            float64   `json:"netWeight"`
	TaxGroupID           string    `json:"taxGroupId"`
	TaxGroupCode         string    `json:"taxGroupCode"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}

// D365Customer represents a customer in Business Central
type D365Customer struct {
	ID                    string    `json:"id"`
	Number                string    `json:"number"`
	DisplayName           string    `json:"displayName"`
	Type                  string    `json:"type"`
	AddressLine1          string    `json:"addressLine1"`
	AddressLine2          string    `json:"addressLine2"`
	City                  string    `json:"city"`
	State                 string    `json:"state"`
	Country               string    `json:"country"`
	PostalCode            string    `json:"postalCode"`
	PhoneNumber           string    `json:"phoneNumber"`
	Email                 string    `json:"email"`
	Website               string    `json:"website"`
	TaxRegistrationNumber string    `json:"taxRegistrationNumber"`
	CurrencyID            string    `json:"currencyId"`
	CurrencyCode          string    `json:"currencyCode"`
	PaymentTermsID        string    `json:"paymentTermsId"`
	PaymentMethodID       string    `json:"paymentMethodId"`
	Blocked               string    `json:"blocked"`
	Balance               float64   `json:"balance"`
	CreditLimit           float64   `json:"creditLimit"`
	LastModifiedDateTime  time.Time `json:"lastModifiedDateTime"`
}

// D365SalesOrder represents a sales order in Business Central
type D365SalesOrder struct {
	ID                      string               `json:"id"`
	Number                  string               `json:"number"`
	OrderDate               string               `json:"orderDate"`
	PostingDate             string               `json:"postingDate"`
	CustomerID              string               `json:"customerId"`
	CustomerNumber          string               `json:"customerNumber"`
	CustomerName            string               `json:"customerName"`
	BillToCustomerID        string               `json:"billToCustomerId"`
	BillToCustomerNumber    string               `json:"billToCustomerNumber"`
	ShipToName              string               `json:"shipToName"`
	ShipToContact           string               `json:"shipToContact"`
	ShipToAddressLine1      string               `json:"shipToAddressLine1"`
	ShipToAddressLine2      string               `json:"shipToAddressLine2"`
	ShipToCity              string               `json:"shipToCity"`
	ShipToCountry           string               `json:"shipToCountry"`
	ShipToState             string               `json:"shipToState"`
	ShipToPostCode          string               `json:"shipToPostCode"`
	CurrencyID              string               `json:"currencyId"`
	CurrencyCode            string               `json:"currencyCode"`
	PaymentTermsID          string               `json:"paymentTermsId"`
	ShipmentMethodID        string               `json:"shipmentMethodId"`
	SalespersonCode         string               `json:"salesperson"`
	PartialShipping         bool                 `json:"partialShipping"`
	RequestedDeliveryDate   string               `json:"requestedDeliveryDate"`
	TotalAmountExcludingTax float64              `json:"totalAmountExcludingTax"`
	TotalTaxAmount          float64              `json:"totalTaxAmount"`
	TotalAmountIncludingTax float64              `json:"totalAmountIncludingTax"`
	FullyShipped            bool                 `json:"fullyShipped"`
	Status                  string               `json:"status"`
	LastModifiedDateTime    time.Time            `json:"lastModifiedDateTime"`
	ExternalDocumentNumber  string               `json:"externalDocumentNumber"`
	SalesOrderLines         []D365SalesOrderLine `json:"salesOrderLines,omitempty"`
}

// D365SalesOrderLine represents a line item in a sales order
type D365SalesOrderLine struct {
	ID                       string  `json:"id"`
	DocumentID               string  `json:"documentId"`
	Sequence                 int     `json:"sequence"`
	ItemID                   string  `json:"itemId"`
	ItemNumber               string  `json:"itemNumber"`
	AccountID                string  `json:"accountId"`
	LineType                 string  `json:"lineType"`
	Description              string  `json:"description"`
	Quantity                 float64 `json:"quantity"`
	UnitOfMeasureID          string  `json:"unitOfMeasureId"`
	UnitOfMeasureCode        string  `json:"unitOfMeasureCode"`
	UnitPrice                float64 `json:"unitPrice"`
	DiscountAmount           float64 `json:"discountAmount"`
	DiscountPercent          float64 `json:"discountPercent"`
	DiscountAppliedBeforeTax bool    `json:"discountAppliedBeforeTax"`
	AmountExcludingTax       float64 `json:"amountExcludingTax"`
	TaxCode                  string  `json:"taxCode"`
	TaxPercent               float64 `json:"taxPercent"`
	TotalTaxAmount           float64 `json:"totalTaxAmount"`
	AmountIncludingTax       float64 `json:"amountIncludingTax"`
	NetAmount                float64 `json:"netAmount"`
	NetAmountIncludingTax    float64 `json:"netAmountIncludingTax"`
	LineDetailsNumber        string  `json:"lineDetailsNumber"`
	ShipmentDate             string  `json:"shipmentDate"`
}

// D365SalesPrice represents a sales price in Business Central
type D365SalesPrice struct {
	ID                   string    `json:"id"`
	ItemNo               string    `json:"itemNo"`
	ItemDescription      string    `json:"itemDescription"`
	SalesType            string    `json:"salesType"`
	SalesCode            string    `json:"salesCode"`
	StartingDate         string    `json:"startingDate"`
	EndingDate           string    `json:"endingDate"`
	CurrencyCode         string    `json:"currencyCode"`
	VariantCode          string    `json:"variantCode"`
	UnitOfMeasureCode    string    `json:"unitOfMeasureCode"`
	MinimumQuantity      float64   `json:"minimumQuantity"`
	UnitPrice            float64   `json:"unitPrice"`
	PriceIncludesVAT     bool      `json:"priceIncludesVAT"`
	AllowInvoiceDisc     bool      `json:"allowInvoiceDisc"`
	AllowLineDisc        bool      `json:"allowLineDisc"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}

// D365ItemLedgerEntry represents inventory movement
type D365ItemLedgerEntry struct {
	ID                   string    `json:"id"`
	EntryNo              int       `json:"entryNo"`
	ItemNo               string    `json:"itemNo"`
	PostingDate          string    `json:"postingDate"`
	EntryType            string    `json:"entryType"`
	SourceNo             string    `json:"sourceNo"`
	DocumentNo           string    `json:"documentNo"`
	Description          string    `json:"description"`
	LocationCode         string    `json:"locationCode"`
	Quantity             float64   `json:"quantity"`
	RemainingQuantity    float64   `json:"remainingQuantity"`
	Open                 bool      `json:"open"`
	Positive             bool      `json:"positive"`
	SourceType           string    `json:"sourceType"`
	CostAmountActual     float64   `json:"costAmountActual"`
	CostAmountExpected   float64   `json:"costAmountExpected"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}

// D365Location represents a warehouse location
type D365Location struct {
	ID                       string    `json:"id"`
	Code                     string    `json:"code"`
	Name                     string    `json:"name"`
	Contact                  string    `json:"contact"`
	AddressLine1             string    `json:"addressLine1"`
	AddressLine2             string    `json:"addressLine2"`
	City                     string    `json:"city"`
	PostCode                 string    `json:"postCode"`
	CountryRegionCode        string    `json:"countryRegionCode"`
	PhoneNo                  string    `json:"phoneNo"`
	Email                    string    `json:"email"`
	HomePage                 string    `json:"homePage"`
	UseAsInTransit           bool      `json:"useAsInTransit"`
	RequirePutAway           bool      `json:"requirePutAway"`
	RequirePick              bool      `json:"requirePick"`
	CrossDockDueDate         string    `json:"crossDockDueDateCalculation"`
	UseCrossDocking          bool      `json:"useCrossDocking"`
	RequireReceive           bool      `json:"requireReceive"`
	RequireShipment          bool      `json:"requireShipment"`
	BinMandatory             bool      `json:"binMandatory"`
	DirectedPutAwayAndPick   bool      `json:"directedPutAwayAndPick"`
	DefaultBinSelection      string    `json:"defaultBinSelection"`
	OutboundWhseHandlingTime string    `json:"outboundWhseHandlingTime"`
	InboundWhseHandlingTime  string    `json:"inboundWhseHandlingTime"`
	LastModifiedDateTime     time.Time `json:"lastModifiedDateTime"`
}

// D365Vendor represents a vendor/supplier
type D365Vendor struct {
	ID                    string    `json:"id"`
	Number                string    `json:"number"`
	DisplayName           string    `json:"displayName"`
	AddressLine1          string    `json:"addressLine1"`
	AddressLine2          string    `json:"addressLine2"`
	City                  string    `json:"city"`
	State                 string    `json:"state"`
	Country               string    `json:"country"`
	PostalCode            string    `json:"postalCode"`
	PhoneNumber           string    `json:"phoneNumber"`
	Email                 string    `json:"email"`
	Website               string    `json:"website"`
	TaxRegistrationNumber string    `json:"taxRegistrationNumber"`
	CurrencyID            string    `json:"currencyId"`
	CurrencyCode          string    `json:"currencyCode"`
	PaymentTermsID        string    `json:"paymentTermsId"`
	PaymentMethodID       string    `json:"paymentMethodId"`
	Blocked               string    `json:"blocked"`
	Balance               float64   `json:"balance"`
	LastModifiedDateTime  time.Time `json:"lastModifiedDateTime"`
}

// D365PurchaseOrder represents a purchase order
type D365PurchaseOrder struct {
	ID                      string    `json:"id"`
	Number                  string    `json:"number"`
	OrderDate               string    `json:"orderDate"`
	PostingDate             string    `json:"postingDate"`
	VendorID                string    `json:"vendorId"`
	VendorNumber            string    `json:"vendorNumber"`
	VendorName              string    `json:"vendorName"`
	PayToVendorID           string    `json:"payToVendorId"`
	PayToVendorNumber       string    `json:"payToVendorNumber"`
	ShipToName              string    `json:"shipToName"`
	ShipToContact           string    `json:"shipToContact"`
	BuyFromAddressLine1     string    `json:"buyFromAddressLine1"`
	BuyFromAddressLine2     string    `json:"buyFromAddressLine2"`
	BuyFromCity             string    `json:"buyFromCity"`
	BuyFromCountry          string    `json:"buyFromCountry"`
	BuyFromState            string    `json:"buyFromState"`
	BuyFromPostCode         string    `json:"buyFromPostCode"`
	ShipToAddressLine1      string    `json:"shipToAddressLine1"`
	ShipToAddressLine2      string    `json:"shipToAddressLine2"`
	ShipToCity              string    `json:"shipToCity"`
	ShipToCountry           string    `json:"shipToCountry"`
	ShipToState             string    `json:"shipToState"`
	ShipToPostCode          string    `json:"shipToPostCode"`
	CurrencyID              string    `json:"currencyId"`
	CurrencyCode            string    `json:"currencyCode"`
	Status                  string    `json:"status"`
	TotalAmountExcludingTax float64   `json:"totalAmountExcludingTax"`
	TotalTaxAmount          float64   `json:"totalTaxAmount"`
	TotalAmountIncludingTax float64   `json:"totalAmountIncludingTax"`
	FullyReceived           bool      `json:"fullyReceived"`
	LastModifiedDateTime    time.Time `json:"lastModifiedDateTime"`
}

// Helper types for API responses

// D365APIResponse represents a standard API response with pagination
type D365APIResponse struct {
	ODataContext string        `json:"@odata.context"`
	Value        []interface{} `json:"value"`
	NextLink     string        `json:"@odata.nextLink,omitempty"`
}

// D365Error represents an API error response
type D365Error struct {
	Error D365ErrorDetails `json:"error"`
}

// D365ErrorDetails contains error details
type D365ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OAuth2TokenResponse represents the OAuth2 token response
type OAuth2TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

// Subscription represents a webhook subscription
type D365Subscription struct {
	ID                        string    `json:"id"`
	Resource                  string    `json:"resource"`
	ChangeType                string    `json:"changeType"`
	ClientState               string    `json:"clientState"`
	NotificationURL           string    `json:"notificationUrl"`
	ExpirationDateTime        time.Time `json:"expirationDateTime"`
	CreatorID                 string    `json:"creatorId"`
	LatestSupportedTLSVersion string    `json:"latestSupportedTlsVersion"`
	EncryptionCertificate     string    `json:"encryptionCertificate"`
	EncryptionCertificateID   string    `json:"encryptionCertificateId"`
	IncludeResourceData       bool      `json:"includeResourceData"`
}
