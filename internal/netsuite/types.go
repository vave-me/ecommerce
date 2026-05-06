package netsuite

import (
	"time"
)

// NetSuiteWebhook represents an incoming webhook from NetSuite
type NetSuiteWebhook struct {
	EventID    string                 `json:"eventId"`
	EventType  string                 `json:"eventType"`
	RecordType string                 `json:"recordType"`
	RecordID   string                 `json:"recordId"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	User       string                 `json:"user"`
	Context    map[string]interface{} `json:"context"`
}

// NetSuiteItem represents an item/product in NetSuite
type NetSuiteItem struct {
	InternalID           string  `json:"internalId"`
	ItemID               string  `json:"itemId"`
	DisplayName          string  `json:"displayName"`
	SalesDescription     string  `json:"salesDescription"`
	PurchaseDescription  string  `json:"purchaseDescription"`
	BasePrice            float64 `json:"basePrice"`
	Cost                 float64 `json:"cost"`
	CostingMethod        string  `json:"costingMethod"`
	ItemType             string  `json:"itemType"`
	IsInactive           bool    `json:"isInactive"`
	Weight               float64 `json:"weight"`
	WeightUnit           string  `json:"weightUnit"`
	UnitsType            string  `json:"unitsType"`
	StockUnit            string  `json:"stockUnit"`
	PurchaseUnit         string  `json:"purchaseUnit"`
	SaleUnit             string  `json:"saleUnit"`
	Class                string  `json:"class"`
	Department           string  `json:"department"`
	Location             string  `json:"location"`
	PreferredVendor      string  `json:"preferredVendor"`
	VendorName           string  `json:"vendorName"`
	UPCCode              string  `json:"upcCode"`
	ManufacturerPartNum  string  `json:"manufacturerPartNum"`
	CreatedDate          string  `json:"createdDate"`
	LastModifiedDate     string  `json:"lastModifiedDate"`
}

// NetSuiteInventoryBalance represents inventory balance
type NetSuiteInventoryBalance struct {
	Item              string  `json:"item"`
	ItemID            string  `json:"itemId"`
	Location          string  `json:"location"`
	LocationName      string  `json:"locationName"`
	QuantityOnHand    float64 `json:"quantityOnHand"`
	QuantityAvailable float64 `json:"quantityAvailable"`
	QuantityOnOrder   float64 `json:"quantityOnOrder"`
	QuantityInTransit float64 `json:"quantityInTransit"`
	QuantityCommitted float64 `json:"quantityCommitted"`
	QuantityBackOrder float64 `json:"quantityBackOrder"`
}

// NetSuiteCustomer represents a customer in NetSuite
type NetSuiteCustomer struct {
	InternalID        string  `json:"internalId"`
	EntityID          string  `json:"entityId"`
	CompanyName       string  `json:"companyName"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	Fax               string  `json:"fax"`
	URL               string  `json:"url"`
	Category          string  `json:"category"`
	IsInactive        bool    `json:"isInactive"`
	Balance           float64 `json:"balance"`
	CreditLimit       float64 `json:"creditLimit"`
	Terms             string  `json:"terms"`
	Currency          string  `json:"currency"`
	PriceLevel        string  `json:"priceLevel"`
	TaxItem           string  `json:"taxItem"`
	TaxRegistration   string  `json:"taxRegistration"`
	DefaultAddress    string  `json:"defaultAddress"`
	CreatedDate       string  `json:"createdDate"`
	LastModifiedDate  string  `json:"lastModifiedDate"`
}

// NetSuiteAddress represents an address
type NetSuiteAddress struct {
	InternalID      string `json:"internalId"`
	DefaultBilling  bool   `json:"defaultBilling"`
	DefaultShipping bool   `json:"defaultShipping"`
	Label           string `json:"label"`
	Attention       string `json:"attention"`
	Addressee       string `json:"addressee"`
	Phone           string `json:"phone"`
	Addr1           string `json:"addr1"`
	Addr2           string `json:"addr2"`
	City            string `json:"city"`
	State           string `json:"state"`
	Zip             string `json:"zip"`
	Country         string `json:"country"`
}

// NetSuiteSalesOrder represents a sales order
type NetSuiteSalesOrder struct {
	InternalID        string                   `json:"internalId"`
	TranID            string                   `json:"tranId"`
	TranDate          string                   `json:"tranDate"`
	Entity            string                   `json:"entity"`
	EntityName        string                   `json:"entityName"`
	Status            string                   `json:"status"`
	Total             float64                  `json:"total"`
	SubTotal          float64                  `json:"subTotal"`
	TaxTotal          float64                  `json:"taxTotal"`
	ShippingCost      float64                  `json:"shippingCost"`
	HandlingCost      float64                  `json:"handlingCost"`
	DiscountTotal     float64                  `json:"discountTotal"`
	Currency          string                   `json:"currency"`
	ExchangeRate      float64                  `json:"exchangeRate"`
	Terms             string                   `json:"terms"`
	SalesRep          string                   `json:"salesRep"`
	LeadSource        string                   `json:"leadSource"`
	ShipMethod        string                   `json:"shipMethod"`
	ShipDate          string                   `json:"shipDate"`
	EstShipDate       string                   `json:"estShipDate"`
	Department        string                   `json:"department"`
	Class             string                   `json:"class"`
	Location          string                   `json:"location"`
	Subsidiary        string                   `json:"subsidiary"`
	CreatedDate       string                   `json:"createdDate"`
	LastModifiedDate  string                   `json:"lastModifiedDate"`
	Email             string                   `json:"email"`
	Message           string                   `json:"message"`
	OtherRefNum       string                   `json:"otherRefNum"`
	Memo              string                   `json:"memo"`
	Item              []NetSuiteSalesOrderLine `json:"item"`
}

// NetSuiteSalesOrderLine represents a line item in a sales order
type NetSuiteSalesOrderLine struct {
	Line              int     `json:"line"`
	Item              string  `json:"item"`
	ItemName          string  `json:"itemName"`
	Description       string  `json:"description"`
	Quantity          float64 `json:"quantity"`
	Units             string  `json:"units"`
	Price             string  `json:"price"`
	Rate              float64 `json:"rate"`
	Amount            float64 `json:"amount"`
	TaxCode           string  `json:"taxCode"`
	TaxRate1          float64 `json:"taxRate1"`
	TaxAmount         float64 `json:"taxAmount"`
	GrossAmt          float64 `json:"grossAmt"`
	Department        string  `json:"department"`
	Class             string  `json:"class"`
	Location          string  `json:"location"`
	IsClosed          bool    `json:"isClosed"`
	OrderPriority     float64 `json:"orderPriority"`
	EstimatedShipDate string  `json:"estimatedShipDate"`
	ShipMethod        string  `json:"shipMethod"`
	CommitInventory   string  `json:"commitInventory"`
}

// NetSuiteItemPrice represents pricing information
type NetSuiteItemPrice struct {
	InternalID  string  `json:"internalId"`
	Item        string  `json:"item"`
	ItemName    string  `json:"itemName"`
	PriceLevel  string  `json:"priceLevel"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	PriceType   string  `json:"priceType"`
	Discount    float64 `json:"discount"`
}

// NetSuitePriceLevel represents a price level
type NetSuitePriceLevel struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	IsInactive       bool   `json:"isInactive"`
	IsOnline         bool   `json:"isOnline"`
	UpdateExisting   bool   `json:"updateExisting"`
}

// NetSuiteVendor represents a vendor/supplier
type NetSuiteVendor struct {
	InternalID       string  `json:"internalId"`
	EntityID         string  `json:"entityId"`
	CompanyName      string  `json:"companyName"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Fax              string  `json:"fax"`
	URL              string  `json:"url"`
	Category         string  `json:"category"`
	IsInactive       bool    `json:"isInactive"`
	Balance          float64 `json:"balance"`
	CreditLimit      float64 `json:"creditLimit"`
	Terms            string  `json:"terms"`
	Currency         string  `json:"currency"`
	Is1099Eligible   bool    `json:"is1099Eligible"`
	TaxIdNum         string  `json:"taxIdNum"`
	DefaultAddress   string  `json:"defaultAddress"`
	CreatedDate      string  `json:"createdDate"`
	LastModifiedDate string  `json:"lastModifiedDate"`
}

// NetSuitePurchaseOrder represents a purchase order
type NetSuitePurchaseOrder struct {
	InternalID       string  `json:"internalId"`
	TranID           string  `json:"tranId"`
	TranDate         string  `json:"tranDate"`
	Entity           string  `json:"entity"`
	EntityName       string  `json:"entityName"`
	Status           string  `json:"status"`
	Total            float64 `json:"total"`
	Currency         string  `json:"currency"`
	ExchangeRate     float64 `json:"exchangeRate"`
	Terms            string  `json:"terms"`
	DueDate          string  `json:"dueDate"`
	Department       string  `json:"department"`
	Class            string  `json:"class"`
	Location         string  `json:"location"`
	Subsidiary       string  `json:"subsidiary"`
	CreatedDate      string  `json:"createdDate"`
	LastModifiedDate string  `json:"lastModifiedDate"`
	Memo             string  `json:"memo"`
}

// NetSuiteInventoryAdjustment represents an inventory adjustment
type NetSuiteInventoryAdjustment struct {
	InternalID       string  `json:"internalId"`
	TranID           string  `json:"tranId"`
	TranDate         string  `json:"tranDate"`
	Account          string  `json:"account"`
	AdjLocation      string  `json:"adjLocation"`
	Customer         string  `json:"customer"`
	Department       string  `json:"department"`
	Class            string  `json:"class"`
	EstimatedValue   float64 `json:"estimatedValue"`
	Subsidiary       string  `json:"subsidiary"`
	CreatedDate      string  `json:"createdDate"`
	LastModifiedDate string  `json:"lastModifiedDate"`
	Memo             string  `json:"memo"`
}

// NetSuiteLocation represents a warehouse/store location
type NetSuiteLocation struct {
	InternalID           string `json:"internalId"`
	Name                 string `json:"name"`
	IsInactive           bool   `json:"isInactive"`
	MakeInventoryAvailable bool `json:"makeInventoryAvailable"`
	MakeInventoryAvailableStore bool `json:"makeInventoryAvailableStore"`
	UsesBins             bool   `json:"usesBins"`
	IsIncludeInSupplyPlanning bool `json:"isIncludeInSupplyPlanning"`
	Subsidiary           string `json:"subsidiary"`
	MainAddress          NetSuiteAddress `json:"mainAddress"`
}

// SuiteQL Response structures
type SuiteQLResponse struct {
	Items      []map[string]interface{} `json:"items"`
	TotalCount int                      `json:"totalCount"`
	Offset     int                      `json:"offset"`
	HasMore    bool                     `json:"hasMore"`
}

// RESTlet request/response structures
type RESTletRequest struct {
	Operation string                 `json:"operation"`
	Query     string                 `json:"query,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type RESTletResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   *RESTletError          `json:"error,omitempty"`
}

type RESTletError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NetSuiteTaxItem represents a tax item
type NetSuiteTaxItem struct {
	InternalID       string  `json:"internalId"`
	ItemID           string  `json:"itemId"`
	DisplayName      string  `json:"displayName"`
	Description      string  `json:"description"`
	Rate             float64 `json:"rate"`
	IsInactive       bool    `json:"isInactive"`
	TaxType          string  `json:"taxType"`
	TaxAgency        string  `json:"taxAgency"`
	PurchaseAccount  string  `json:"purchaseAccount"`
	SaleAccount      string  `json:"saleAccount"`
}

// NetSuiteSubsidiary represents a subsidiary
type NetSuiteSubsidiary struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	IsInactive       bool   `json:"isInactive"`
	IsElimination    bool   `json:"isElimination"`
	LegalName        string `json:"legalName"`
	Currency         string `json:"currency"`
	Email            string `json:"email"`
	URL              string `json:"url"`
	MainAddress      NetSuiteAddress `json:"mainAddress"`
	ShippingAddress  NetSuiteAddress `json:"shippingAddress"`
	ReturnAddress    NetSuiteAddress `json:"returnAddress"`
	FederalIdNumber  string `json:"federalIdNumber"`
}

// NetSuiteCurrency represents a currency
type NetSuiteCurrency struct {
	InternalID       string  `json:"internalId"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	IsInactive       bool    `json:"isInactive"`
	IsBaseCurrency   bool    `json:"isBaseCurrency"`
	ExchangeRate     float64 `json:"exchangeRate"`
	CurrencyPrecision string `json:"currencyPrecision"`
}