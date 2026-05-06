package services

import (
	"middleman/managers/internal/models"
	"strings"
)

// LLMOperationSchema provides complete schema information for LLM communication
type LLMOperationSchema struct {
	EntityType       models.EntityType        `json:"entity_type"`
	TableName        string                   `json:"table_name"`
	Description      string                   `json:"description"`
	Fields           []FieldDefinition        `json:"fields"`
	Operations       []OperationDefinition    `json:"operations"`
	Relationships    []RelationshipDefinition `json:"relationships"`
	SearchableFields []string                 `json:"searchable_fields"`
	SortableFields   []string                 `json:"sortable_fields"`
	FilterableFields []string                 `json:"filterable_fields"`
	RequiredFields   []string                 `json:"required_fields"`
	Examples         []OperationExample       `json:"examples"`
}

// FieldDefinition describes a field available for operations
type FieldDefinition struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Searchable  bool     `json:"searchable"`
	Sortable    bool     `json:"sortable"`
	Filterable  bool     `json:"filterable"`
	Constraints []string `json:"constraints"`
	Examples    []string `json:"examples"`
	RelatedTo   string   `json:"related_to,omitempty"`
}

// OperationDefinition describes an available repository operation
type OperationDefinition struct {
	Name                    string          `json:"name"`
	Method                  string          `json:"method"`
	Description             string          `json:"description"`
	Parameters              []ParameterInfo `json:"parameters"`
	ReturnType              string          `json:"return_type"`
	Examples                []string        `json:"examples"`
	RequiredFields          []string        `json:"required_fields"`
	OptionalFields          []string        `json:"optional_fields"`
	NaturalLanguagePatterns []string        `json:"natural_language_patterns"`
}

// ParameterInfo describes operation parameters
type ParameterInfo struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Required     bool        `json:"required"`
	Description  string      `json:"description"`
	DefaultValue interface{} `json:"default_value,omitempty"`
}

// RelationshipDefinition describes entity relationships
type RelationshipDefinition struct {
	Name         string `json:"name"`
	Type         string `json:"type"` // one-to-one, one-to-many, many-to-many
	TargetEntity string `json:"target_entity"`
	ForeignKey   string `json:"foreign_key"`
	Description  string `json:"description"`
}

// OperationExample provides usage examples for LLM understanding
type OperationExample struct {
	Description    string                 `json:"description"`
	NaturalQuery   string                 `json:"natural_query"`
	Operation      string                 `json:"operation"`
	Parameters     map[string]interface{} `json:"parameters"`
	ExpectedResult string                 `json:"expected_result"`
}

// LLMSchemaRegistry maintains complete schema information for all entities
type LLMSchemaRegistry struct {
	schemas map[models.EntityType]*LLMOperationSchema
}

// NewLLMSchemaRegistry creates a new schema registry with complete entity definitions
func NewLLMSchemaRegistry() *LLMSchemaRegistry {
	registry := &LLMSchemaRegistry{
		schemas: make(map[models.EntityType]*LLMOperationSchema),
	}
	registry.initializeProductSchema()
	registry.initializeUserSchema()
	registry.initializeDealSchema()
	// Removed vehicle, property, and job schemas per requirements
	// registry.initializeVehicleSchema()
	// registry.initializePropertySchema()
	// registry.initializeJobSchema()
	registry.initializeOrderSchema()
	registry.initializePaymentSchema()
	registry.initializeOfferSchema()
	registry.initializeReviewSchema()
	registry.initializeCommentSchema()
	registry.initializeNotificationSchema()
	registry.initializeNewsletterSchema()
	registry.initializeBasketSchema()
	registry.initializeCategorySchema()
	registry.initializeMetricSchema()
	registry.initializeMessagesSchema()
	registry.initializeWishlistSchema()
	registry.initializeFollowingSchema()
	registry.initializeActivitySchema()
	registry.initializeMediaSchema()
	registry.initializeServiceSchema()
	registry.initializeShippingSchema()
	registry.initializeSupportSchema()
	registry.initializeGeocodingSchema()
	registry.initializeVariantSchema()

	return registry
}

// GetSchema returns schema for specific entity type
func (r *LLMSchemaRegistry) GetSchema(EntityType models.EntityType) *LLMOperationSchema {
	return r.schemas[EntityType]
}

// GetAllSchemas returns all available schemas
func (r *LLMSchemaRegistry) GetAllSchemas() map[models.EntityType]*LLMOperationSchema {
	return r.schemas
}

// GetOperationsByEntityType returns all operations for an entity type
func (r *LLMSchemaRegistry) GetOperationsByEntityType(entityType models.EntityType) []OperationDefinition {
	schema := r.GetSchema(entityType)
	if schema == nil {
		return []OperationDefinition{}
	}
	return schema.Operations
}

// GetFieldsByEntityType returns all fields for an entity type
func (r *LLMSchemaRegistry) GetFieldsByEntityType(entityType models.EntityType) []FieldDefinition {
	schema := r.GetSchema(entityType)
	if schema == nil {
		return []FieldDefinition{}
	}
	return schema.Fields
}

// GenerateNaturalLanguageHelp generates comprehensive help text for LLM
func (r *LLMSchemaRegistry) GenerateNaturalLanguageHelp() string {
	var help strings.Builder

	help.WriteString("# COMPREHENSIVE REPOSITORY OPERATIONS GUIDE\n\n")
	help.WriteString("You have access to complete database operations across all entity types. ")
	help.WriteString("Each entity supports multiple operations with full field access.\n\n")

	for entityType, schema := range r.schemas {
		help.WriteString(generateEntityHelp(entityType, schema))
	}

	return help.String()
}

// Helper method to generate help for each entity
func generateEntityHelp(entityType models.EntityType, schema *LLMOperationSchema) string {
	var help strings.Builder

	help.WriteString("## " + strings.ToUpper(string(entityType)) + "\n")
	help.WriteString(schema.Description + "\n\n")

	// Available Operations
	help.WriteString("### Available Operations:\n")
	for _, op := range schema.Operations {
		help.WriteString("- **" + op.Name + "**: " + op.Description + "\n")
		if len(op.NaturalLanguagePatterns) > 0 {
			help.WriteString("  Natural language: " + strings.Join(op.NaturalLanguagePatterns, ", ") + "\n")
		}
	}
	help.WriteString("\n")

	// Available Fields
	help.WriteString("### Available Fields:\n")
	searchableFields := []string{}
	filterableFields := []string{}
	for _, field := range schema.Fields {
		help.WriteString("- **" + field.Name + "** (" + field.Type + "): " + field.Description)
		if field.Required {
			help.WriteString(" [REQUIRED]")
		}
		if field.Searchable {
			searchableFields = append(searchableFields, field.Name)
		}
		if field.Filterable {
			filterableFields = append(filterableFields, field.Name)
		}
		help.WriteString("\n")
	}

	if len(searchableFields) > 0 {
		help.WriteString("\n**Searchable fields**: " + strings.Join(searchableFields, ", ") + "\n")
	}
	if len(filterableFields) > 0 {
		help.WriteString("**Filterable fields**: " + strings.Join(filterableFields, ", ") + "\n")
	}

	// Examples
	if len(schema.Examples) > 0 {
		help.WriteString("\n### Examples:\n")
		for _, example := range schema.Examples {
			help.WriteString("- \"" + example.NaturalQuery + "\" → " + example.Operation + "\n")
		}
	}

	help.WriteString("\n")
	return help.String()
}

// initializeProductSchema defines complete product schema
func (r *LLMSchemaRegistry) initializeProductSchema() {
	r.schemas[models.ProductType] = &LLMOperationSchema{
		EntityType:  models.ProductType,
		TableName:   "products",
		Description: "Product catalog management with complete CRUD operations, pricing, inventory, and search capabilities",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Description: "Unique product identifier", Required: true, Searchable: true},
			{Name: "name", Type: "string", Description: "Product name/title", Required: true, Searchable: true, Filterable: true},
			{Name: "description", Type: "string", Description: "Product description", Searchable: true},
			{Name: "base_price", Type: "int64", Description: "Base price in cents", Required: true, Filterable: true, Sortable: true},
			{Name: "user_seller_id", Type: "string", Description: "Seller user ID", Required: true, Filterable: true},
			{Name: "category_id", Type: "string", Description: "Product category ID", Filterable: true, RelatedTo: "categories"},
			{Name: "category_slug", Type: "string", Description: "Product category slug", Filterable: true},
			{Name: "brand", Type: "string", Description: "Product brand", Filterable: true, Searchable: true},
			{Name: "condition", Type: "string", Description: "Product condition (new, used, refurbished)", Filterable: true},
			{Name: "model", Type: "string", Description: "Product model", Filterable: true, Searchable: true},
			{Name: "tags", Type: "[]string", Description: "Product tags for categorization", Searchable: true},
			{Name: "manage_stock", Type: "bool", Description: "Whether stock is managed", Filterable: true},
			{Name: "stock", Type: "int64", Description: "Available stock quantity", Filterable: true, Sortable: true},
			{Name: "sku", Type: "string", Description: "Stock keeping unit", Searchable: true},
			{Name: "status", Type: "string", Description: "Product status (active, archived, sold)", Filterable: true},
			{Name: "negotiable", Type: "bool", Description: "Whether price is negotiable", Filterable: true},
			{Name: "user_type", Type: "string", Description: "Type of seller", Filterable: true},
			{Name: "middleman_service", Type: "bool", Description: "Whether middleman service is used", Filterable: true},
			{Name: "shipping_cost", Type: "int64", Description: "Shipping cost in cents", Filterable: true},
			{Name: "has_variants", Type: "bool", Description: "Whether product has variants", Filterable: true},
			{Name: "weight", Type: "int64", Description: "Product weight in grams", Filterable: true},
			{Name: "height", Type: "int64", Description: "Product height in mm", Filterable: true},
			{Name: "width", Type: "int64", Description: "Product width in mm", Filterable: true},
			{Name: "depth", Type: "int64", Description: "Product depth in mm", Filterable: true},
			{Name: "lat", Type: "float64", Description: "Latitude for location-based search"},
			{Name: "lng", Type: "float64", Description: "Longitude for location-based search"},
			{Name: "thumbnail", Type: "string", Description: "Product thumbnail URL"},
		},
		Operations: []OperationDefinition{
			{
				Name:        "search_with_term",
				Method:      "SearchWithTerm",
				Description: "Search products by name/term with full-text search",
				Parameters: []ParameterInfo{
					{Name: "name", Type: "string", Required: true, Description: "Search term for product name"},
				},
				ReturnType:              "[]*models.Product",
				NaturalLanguagePatterns: []string{"search for products", "find products containing", "look for products"},
				Examples:                []string{"search for 'laptop'", "find products with 'iPhone'"},
			},
			{
				Name:        "search_with_filters",
				Method:      "SearchWithFilters",
				Description: "Advanced search with multiple filters including price, category, location",
				Parameters: []ParameterInfo{
					{Name: "name", Type: "string", Description: "Product name filter"},
					{Name: "category", Type: "string", Description: "Category filter"},
					{Name: "min_price", Type: "int64", Description: "Minimum price in cents"},
					{Name: "max_price", Type: "int64", Description: "Maximum price in cents"},
					{Name: "brand", Type: "string", Description: "Brand filter"},
					{Name: "condition", Type: "string", Description: "Condition filter"},
					{Name: "lat", Type: "float64", Description: "Latitude for location search"},
					{Name: "lng", Type: "float64", Description: "Longitude for location search"},
					{Name: "radius", Type: "int64", Description: "Search radius in meters"},
					{Name: "page", Type: "int64", Description: "Page number for pagination"},
					{Name: "page_size", Type: "int64", Description: "Items per page"},
					{Name: "sort_by", Type: "string", Description: "Sort field (price, name, created_at)"},
					{Name: "sort_order", Type: "string", Description: "Sort order (asc, desc)"},
				},
				ReturnType:              "[]*models.Product",
				NaturalLanguagePatterns: []string{"filter products by", "find products with", "search products where"},
				Examples:                []string{"find laptops under $1000", "search products near me under $500"},
			},
			{
				Name:        "find",
				Method:      "Find",
				Description: "Get specific product by ID",
				Parameters: []ParameterInfo{
					{Name: "product_id", Type: "string", Required: true, Description: "Product ID to retrieve"},
				},
				ReturnType:              "*models.Product",
				NaturalLanguagePatterns: []string{"get product", "find product by id", "retrieve product"},
				Examples:                []string{"get product with id 'abc123'"},
			},
			{
				Name:           "add",
				Method:         "Add",
				Description:    "Create new product with full details",
				RequiredFields: []string{"name", "description", "base_price", "user_seller_id", "category_id"},
				Parameters: []ParameterInfo{
					{Name: "name", Type: "string", Required: true, Description: "Product name"},
					{Name: "description", Type: "string", Required: true, Description: "Product description"},
					{Name: "base_price", Type: "int64", Required: true, Description: "Base price in cents"},
					{Name: "user_seller_id", Type: "string", Required: true, Description: "Seller ID"},
					{Name: "category_id", Type: "string", Required: true, Description: "Category ID"},
					{Name: "brand", Type: "string", Description: "Product brand"},
					{Name: "condition", Type: "string", Description: "Product condition"},
					{Name: "stock", Type: "int64", Description: "Initial stock quantity"},
				},
				ReturnType:              "error",
				NaturalLanguagePatterns: []string{"create product", "add new product", "list product"},
				Examples:                []string{"create a laptop product for $999"},
			},
			{
				Name:        "update",
				Method:      "Update",
				Description: "Update product price",
				Parameters: []ParameterInfo{
					{Name: "product_id", Type: "string", Required: true, Description: "Product ID"},
					{Name: "new_base_price", Type: "int64", Required: true, Description: "New price in cents"},
				},
				ReturnType:              "error",
				NaturalLanguagePatterns: []string{"update product price", "change price", "set new price"},
			},
			{
				Name:        "remove",
				Method:      "Remove",
				Description: "Delete product from catalog",
				Parameters: []ParameterInfo{
					{Name: "product_id", Type: "string", Required: true, Description: "Product ID to delete"},
				},
				ReturnType:              "error",
				NaturalLanguagePatterns: []string{"delete product", "remove product", "delete listing"},
			},
		},
		SearchableFields: []string{"name", "description", "brand", "model", "tags", "sku"},
		SortableFields:   []string{"base_price", "stock", "created_at", "name"},
		FilterableFields: []string{"category_id", "brand", "condition", "status", "user_type", "negotiable"},
		RequiredFields:   []string{"name", "description", "base_price", "user_seller_id", "category_id"},
		Examples: []OperationExample{
			{
				Description:  "Search for electronics under $500",
				NaturalQuery: "find electronics under $500",
				Operation:    "search_with_filters",
				Parameters: map[string]interface{}{
					"category":  "electronics",
					"max_price": 50000,
				},
				ExpectedResult: "List of electronic products under $500",
			},
			{
				Description:  "Get specific product details",
				NaturalQuery: "show me product abc123",
				Operation:    "find",
				Parameters: map[string]interface{}{
					"product_id": "abc123",
				},
				ExpectedResult: "Complete product details for abc123",
			},
		},
	}
}

// initializeUserSchema defines complete user schema
func (r *LLMSchemaRegistry) initializeUserSchema() {
	r.schemas[models.UserEntityType] = &LLMOperationSchema{
		EntityType:  models.UserEntityType,
		TableName:   "users",
		Description: "User management with authentication, profiles, and account operations",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Description: "Unique user identifier", Required: true, Searchable: true},
			{Name: "email", Type: "string", Description: "User email address", Required: true, Searchable: true},
			{Name: "username", Type: "string", Description: "Unique username", Required: true, Searchable: true},
			{Name: "first_name", Type: "string", Description: "User first name", Searchable: true},
			{Name: "last_name", Type: "string", Description: "User last name", Searchable: true},
			{Name: "location", Type: "string", Description: "User location/address", Searchable: true},
			{Name: "lat", Type: "float32", Description: "User latitude"},
			{Name: "lng", Type: "float32", Description: "User longitude"},
			{Name: "thumbnail", Type: "string", Description: "Profile picture URL"},
			{Name: "language", Type: "string", Description: "Preferred language"},
			{Name: "bio", Type: "string", Description: "User biography", Searchable: true},
			{Name: "privacy", Type: "string", Description: "Privacy settings"},
			{Name: "background", Type: "string", Description: "Background image URL"},
		},
		Operations: []OperationDefinition{
			{
				Name:        "find",
				Method:      "Find",
				Description: "Get user by ID with complete profile",
				Parameters: []ParameterInfo{
					{Name: "user_id", Type: "string", Required: true, Description: "User ID to retrieve"},
				},
				ReturnType:              "*models.User",
				NaturalLanguagePatterns: []string{"get user", "find user", "show user profile"},
			},
			{
				Name:           "create_user",
				Method:         "CreateUser",
				Description:    "Register new user account",
				RequiredFields: []string{"email", "password", "username", "first_name", "last_name"},
				Parameters: []ParameterInfo{
					{Name: "email", Type: "string", Required: true, Description: "User email"},
					{Name: "password", Type: "string", Required: true, Description: "User password"},
					{Name: "username", Type: "string", Required: true, Description: "Unique username"},
					{Name: "first_name", Type: "string", Required: true, Description: "First name"},
					{Name: "last_name", Type: "string", Required: true, Description: "Last name"},
					{Name: "location", Type: "string", Description: "User location"},
					{Name: "language", Type: "string", Description: "Preferred language"},
				},
				ReturnType:              "string",
				NaturalLanguagePatterns: []string{"create user", "register user", "sign up user"},
			},
			{
				Name:        "update_user",
				Method:      "UpdateUser",
				Description: "Update user profile information",
				Parameters: []ParameterInfo{
					{Name: "id", Type: "string", Required: true, Description: "User ID"},
					{Name: "username", Type: "string", Description: "New username"},
					{Name: "first_name", Type: "string", Description: "New first name"},
					{Name: "last_name", Type: "string", Description: "New last name"},
					{Name: "bio", Type: "string", Description: "User biography"},
					{Name: "location", Type: "string", Description: "User location"},
				},
				ReturnType:              "string",
				NaturalLanguagePatterns: []string{"update user", "modify profile", "change user details"},
			},
			{
				Name:        "login_user",
				Method:      "LoginUser",
				Description: "Authenticate user with email and password",
				Parameters: []ParameterInfo{
					{Name: "email", Type: "string", Required: true, Description: "User email"},
					{Name: "password", Type: "string", Required: true, Description: "User password"},
				},
				ReturnType:              "*models.LoginResponse",
				NaturalLanguagePatterns: []string{"login user", "authenticate user", "sign in"},
			},
		},
		SearchableFields: []string{"email", "username", "first_name", "last_name", "location", "bio"},
		SortableFields:   []string{"username", "first_name", "last_name", "created_at"},
		FilterableFields: []string{"language", "privacy"},
		RequiredFields:   []string{"email", "username", "first_name", "last_name"},
		Examples: []OperationExample{
			{
				Description:  "Find user by ID",
				NaturalQuery: "get user with id user123",
				Operation:    "find",
				Parameters: map[string]interface{}{
					"user_id": "user123",
				},
				ExpectedResult: "Complete user profile information",
			},
		},
	}
}

// Additional entity schema initialization methods would follow the same pattern...
// For brevity, I'm showing the pattern with Product and User schemas

// Deal schema initialization
func (r *LLMSchemaRegistry) initializeDealSchema() {
	r.schemas[models.DealType] = &LLMOperationSchema{
		EntityType:  models.DealType,
		TableName:   "deals",
		Description: "Manages special deals, discounts, and promotional offers on products and services",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique deal identifier"},
			{Name: "title", Type: "string", Required: true, Description: "Deal title or name"},
			{Name: "description", Type: "string", Required: true, Description: "Detailed deal description"},
			{Name: "discount_percentage", Type: "float", Required: false, Description: "Percentage discount (0-100)"},
			{Name: "discount_amount", Type: "integer", Required: false, Description: "Fixed discount amount in cents"},
			{Name: "original_price", Type: "integer", Required: true, Description: "Original price in cents"},
			{Name: "deal_price", Type: "integer", Required: true, Description: "Discounted price in cents"},
			{Name: "product_id", Type: "string", Required: true, Description: "Associated product ID"},
			{Name: "category_id", Type: "string", Required: false, Description: "Product category ID"},
			{Name: "start_date", Type: "datetime", Required: true, Description: "Deal start date and time"},
			{Name: "end_date", Type: "datetime", Required: true, Description: "Deal expiration date and time"},
			{Name: "max_redemptions", Type: "integer", Required: false, Description: "Maximum number of times deal can be used"},
			{Name: "current_redemptions", Type: "integer", Required: false, Description: "Current number of redemptions"},
			{Name: "is_active", Type: "boolean", Required: true, Description: "Whether deal is currently active"},
			{Name: "deal_type", Type: "string", Required: true, Description: "Type of deal (flash, seasonal, clearance, etc.)"},
			{Name: "terms_conditions", Type: "string", Required: false, Description: "Deal terms and conditions"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Deal creation timestamp"},
			{Name: "updated_at", Type: "datetime", Required: false, Description: "Last update timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find deal by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search deals by title or description", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter deals by criteria", Parameters: []ParameterInfo{
				{Name: "deal_type", Type: "string", Required: false},
				{Name: "category_id", Type: "string", Required: false},
				{Name: "is_active", Type: "boolean", Required: false},
				{Name: "min_discount", Type: "float", Required: false},
				{Name: "max_discount", Type: "float", Required: false},
			}},
			{Name: "add", Description: "Create new deal", Parameters: []ParameterInfo{
				{Name: "title", Type: "string", Required: true},
				{Name: "description", Type: "string", Required: true},
				{Name: "product_id", Type: "string", Required: true},
				{Name: "discount_percentage", Type: "float", Required: false},
				{Name: "discount_amount", Type: "integer", Required: false},
				{Name: "start_date", Type: "datetime", Required: true},
				{Name: "end_date", Type: "datetime", Required: true},
			}},
			{Name: "update", Description: "Update existing deal", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "remove", Description: "Delete deal", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.ProductType), Type: "belongs_to", Description: "Deal belongs to a product"},
			{TargetEntity: string(models.CategoryType), Type: "belongs_to", Description: "Deal belongs to a category"},
		},
		SearchableFields: []string{"title", "description", "deal_type"},
		SortableFields:   []string{"title", "deal_price", "discount_percentage", "start_date", "end_date", "created_at"},
		FilterableFields: []string{"deal_type", "category_id", "is_active", "product_id"},
	}
}

// Vehicle schema initialization
func (r *LLMSchemaRegistry) initializeVehicleSchema() {
	r.schemas[models.VehicleType] = &LLMOperationSchema{
		EntityType:  models.VehicleType,
		TableName:   "vehicles",
		Description: "Manages vehicle listings including cars, motorcycles, trucks, and other automotive items",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique vehicle identifier"},
			{Name: "make", Type: "string", Required: true, Description: "Vehicle manufacturer (e.g., Toyota, Ford)"},
			{Name: "model", Type: "string", Required: true, Description: "Vehicle model name"},
			{Name: "year", Type: "integer", Required: true, Description: "Manufacturing year"},
			{Name: "vin", Type: "string", Required: false, Description: "Vehicle Identification Number"},
			{Name: "mileage", Type: "integer", Required: false, Description: "Current mileage in kilometers"},
			{Name: "fuel_type", Type: "string", Required: false, Description: "Fuel type (gasoline, diesel, electric, hybrid)"},
			{Name: "transmission", Type: "string", Required: false, Description: "Transmission type (manual, automatic, CVT)"},
			{Name: "body_type", Type: "string", Required: false, Description: "Body style (sedan, SUV, hatchback, truck)"},
			{Name: "exterior_color", Type: "string", Required: false, Description: "Exterior paint color"},
			{Name: "interior_color", Type: "string", Required: false, Description: "Interior color scheme"},
			{Name: "engine_size", Type: "float", Required: false, Description: "Engine displacement in liters"},
			{Name: "horsepower", Type: "integer", Required: false, Description: "Engine horsepower"},
			{Name: "condition", Type: "string", Required: true, Description: "Vehicle condition (new, used, certified)"},
			{Name: "price", Type: "integer", Required: true, Description: "Listed price in cents"},
			{Name: "seller_id", Type: "string", Required: true, Description: "Seller user ID"},
			{Name: "location", Type: "string", Required: false, Description: "Vehicle location"},
			{Name: "features", Type: "array", Required: false, Description: "List of vehicle features"},
			{Name: "images", Type: "array", Required: false, Description: "Vehicle image URLs"},
			{Name: "is_available", Type: "boolean", Required: true, Description: "Whether vehicle is available for sale"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Listing creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find vehicle by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search vehicles with complex criteria like make/model/year ranges/specifications - use this for most vehicle search queries", Parameters: []ParameterInfo{
				{Name: "make", Type: "string", Required: false},
				{Name: "model", Type: "string", Required: false},
				{Name: "min_year", Type: "integer", Required: false},
				{Name: "max_year", Type: "integer", Required: false},
				{Name: "min_price", Type: "integer", Required: false},
				{Name: "max_price", Type: "integer", Required: false},
				{Name: "fuel_type", Type: "string", Required: false},
				{Name: "transmission_type", Type: "string", Required: false},
				{Name: "max_mileage", Type: "integer", Required: false},
			}},
			{Name: "filter", Description: "DEPRECATED: Use search operation instead for vehicle queries", Parameters: []ParameterInfo{
				{Name: "make", Type: "string", Required: false},
				{Name: "model", Type: "string", Required: false},
				{Name: "min_year", Type: "integer", Required: false},
				{Name: "max_year", Type: "integer", Required: false},
				{Name: "min_price", Type: "integer", Required: false},
				{Name: "max_price", Type: "integer", Required: false},
				{Name: "fuel_type", Type: "string", Required: false},
				{Name: "transmission", Type: "string", Required: false},
				{Name: "body_type", Type: "string", Required: false},
				{Name: "max_mileage", Type: "integer", Required: false},
			}},
			{Name: "add", Description: "Create vehicle listing", Parameters: []ParameterInfo{
				{Name: "make", Type: "string", Required: true},
				{Name: "model", Type: "string", Required: true},
				{Name: "year", Type: "integer", Required: true},
				{Name: "price", Type: "integer", Required: true},
				{Name: "condition", Type: "string", Required: true},
				{Name: "seller_id", Type: "string", Required: true},
			}},
			{Name: "update", Description: "Update vehicle listing", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "remove", Description: "Delete vehicle listing", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Vehicle belongs to a seller"},
		},
		SearchableFields: []string{"make", "model", "features"},
		SortableFields:   []string{"make", "model", "year", "price", "mileage", "created_at"},
		FilterableFields: []string{"make", "model", "year", "fuel_type", "transmission", "body_type", "condition"},
	}
}

// Order schema initialization
func (r *LLMSchemaRegistry) initializeOrderSchema() {
	r.schemas[models.EntityTypeOrder] = &LLMOperationSchema{
		EntityType:  models.EntityTypeOrder,
		TableName:   "orders",
		Description: "Manages customer orders and purchase transactions",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique order identifier"},
			{Name: "order_number", Type: "string", Required: true, Description: "Human-readable order number"},
			{Name: "buyer_id", Type: "string", Required: true, Description: "Customer user ID"},
			{Name: "seller_id", Type: "string", Required: true, Description: "Seller user ID"},
			{Name: "total_amount", Type: "integer", Required: true, Description: "Total order amount in cents"},
			{Name: "tax_amount", Type: "integer", Required: false, Description: "Tax amount in cents"},
			{Name: "shipping_amount", Type: "integer", Required: false, Description: "Shipping cost in cents"},
			{Name: "discount_amount", Type: "integer", Required: false, Description: "Discount amount in cents"},
			{Name: "status", Type: "string", Required: true, Description: "Order status (pending, confirmed, shipped, delivered, cancelled)"},
			{Name: "payment_status", Type: "string", Required: true, Description: "Payment status (pending, paid, refunded, failed)"},
			{Name: "shipping_address", Type: "object", Required: true, Description: "Shipping address details"},
			{Name: "billing_address", Type: "object", Required: false, Description: "Billing address details"},
			{Name: "items", Type: "array", Required: true, Description: "List of ordered items"},
			{Name: "tracking_number", Type: "string", Required: false, Description: "Shipping tracking number"},
			{Name: "notes", Type: "string", Required: false, Description: "Order notes or special instructions"},
			{Name: "estimated_delivery", Type: "datetime", Required: false, Description: "Estimated delivery date"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Order creation timestamp"},
			{Name: "updated_at", Type: "datetime", Required: false, Description: "Last update timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find order by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search orders by number or customer", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter orders by criteria", Parameters: []ParameterInfo{
				{Name: "buyer_id", Type: "string", Required: false},
				{Name: "seller_id", Type: "string", Required: false},
				{Name: "status", Type: "string", Required: false},
				{Name: "payment_status", Type: "string", Required: false},
				{Name: "min_amount", Type: "integer", Required: false},
				{Name: "max_amount", Type: "integer", Required: false},
				{Name: "start_date", Type: "datetime", Required: false},
				{Name: "end_date", Type: "datetime", Required: false},
			}},
			{Name: "add", Description: "Create new order", Parameters: []ParameterInfo{
				{Name: "buyer_id", Type: "string", Required: true},
				{Name: "seller_id", Type: "string", Required: true},
				{Name: "items", Type: "array", Required: true},
				{Name: "shipping_address", Type: "object", Required: true},
			}},
			{Name: "update", Description: "Update order status", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "remove", Description: "Cancel order", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Order belongs to buyer and seller"},
			{TargetEntity: string(models.ProductType), Type: "has_many", Description: "Order contains multiple products"},
		},
		SearchableFields: []string{"order_number", "tracking_number", "notes"},
		SortableFields:   []string{"total_amount", "created_at", "estimated_delivery"},
		FilterableFields: []string{"status", "payment_status", "buyer_id", "seller_id"},
	}
}

// Payment schema initialization
func (r *LLMSchemaRegistry) initializePaymentSchema() {
	r.schemas[models.PaymentEntityType] = &LLMOperationSchema{
		EntityType:  models.PaymentEntityType,
		TableName:   "payments",
		Description: "Manages payment transactions and financial records",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique payment identifier"},
			{Name: "transaction_id", Type: "string", Required: true, Description: "External transaction identifier"},
			{Name: "order_id", Type: "string", Required: true, Description: "Associated order ID"},
			{Name: "payer_id", Type: "string", Required: true, Description: "User ID making payment"},
			{Name: "payee_id", Type: "string", Required: true, Description: "User ID receiving payment"},
			{Name: "amount", Type: "integer", Required: true, Description: "Payment amount in cents"},
			{Name: "fee_amount", Type: "integer", Required: false, Description: "Processing fee in cents"},
			{Name: "net_amount", Type: "integer", Required: true, Description: "Net amount after fees in cents"},
			{Name: "currency", Type: "string", Required: true, Description: "Payment currency code"},
			{Name: "payment_method", Type: "string", Required: true, Description: "Payment method (card, bank, wallet)"},
			{Name: "status", Type: "string", Required: true, Description: "Payment status (pending, completed, failed, refunded)"},
			{Name: "gateway", Type: "string", Required: true, Description: "Payment gateway used"},
			{Name: "gateway_response", Type: "object", Required: false, Description: "Gateway response data"},
			{Name: "refund_amount", Type: "integer", Required: false, Description: "Refunded amount in cents"},
			{Name: "refund_reason", Type: "string", Required: false, Description: "Reason for refund"},
			{Name: "processed_at", Type: "datetime", Required: false, Description: "Payment processing timestamp"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Payment record creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find payment by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search payments by transaction ID", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter payments by criteria", Parameters: []ParameterInfo{
				{Name: "order_id", Type: "string", Required: false},
				{Name: "payer_id", Type: "string", Required: false},
				{Name: "payee_id", Type: "string", Required: false},
				{Name: "status", Type: "string", Required: false},
				{Name: "payment_method", Type: "string", Required: false},
				{Name: "gateway", Type: "string", Required: false},
				{Name: "min_amount", Type: "integer", Required: false},
				{Name: "max_amount", Type: "integer", Required: false},
			}},
			{Name: "add", Description: "Process new payment", Parameters: []ParameterInfo{
				{Name: "order_id", Type: "string", Required: true},
				{Name: "amount", Type: "integer", Required: true},
				{Name: "payment_method", Type: "string", Required: true},
				{Name: "payer_id", Type: "string", Required: true},
			}},
			{Name: "update", Description: "Update payment status", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "refund", Description: "Process payment refund", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "refund_amount", Type: "integer", Required: true},
				{Name: "refund_reason", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.EntityTypeOrder), Type: "belongs_to", Description: "Payment belongs to an order"},
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Payment involves payer and payee"},
		},
		SearchableFields: []string{"transaction_id", "gateway_response"},
		SortableFields:   []string{"amount", "processed_at", "created_at"},
		FilterableFields: []string{"status", "payment_method", "gateway", "order_id", "payer_id", "payee_id"},
	}
}

// Offer schema initialization
func (r *LLMSchemaRegistry) initializeOfferSchema() {
	r.schemas[models.OfferEntityType] = &LLMOperationSchema{
		EntityType:  models.OfferEntityType,
		TableName:   "offers",
		Description: "Manages purchase offers and negotiations between buyers and sellers",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique offer identifier"},
			{Name: "product_id", Type: "string", Required: true, Description: "Product being offered on"},
			{Name: "buyer_id", Type: "string", Required: true, Description: "User ID making the offer"},
			{Name: "seller_id", Type: "string", Required: true, Description: "User ID receiving the offer"},
			{Name: "offered_price", Type: "integer", Required: true, Description: "Offered price in cents"},
			{Name: "original_price", Type: "integer", Required: true, Description: "Original listed price in cents"},
			{Name: "message", Type: "string", Required: false, Description: "Offer message or negotiation notes"},
			{Name: "status", Type: "string", Required: true, Description: "Offer status (pending, accepted, rejected, countered, expired)"},
			{Name: "counter_price", Type: "integer", Required: false, Description: "Counter offer price in cents"},
			{Name: "counter_message", Type: "string", Required: false, Description: "Counter offer message"},
			{Name: "expires_at", Type: "datetime", Required: false, Description: "Offer expiration timestamp"},
			{Name: "accepted_at", Type: "datetime", Required: false, Description: "Offer acceptance timestamp"},
			{Name: "rejected_at", Type: "datetime", Required: false, Description: "Offer rejection timestamp"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Offer creation timestamp"},
			{Name: "updated_at", Type: "datetime", Required: false, Description: "Last update timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find offer by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search offers by product or participants", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter offers by criteria", Parameters: []ParameterInfo{
				{Name: "product_id", Type: "string", Required: false},
				{Name: "buyer_id", Type: "string", Required: false},
				{Name: "seller_id", Type: "string", Required: false},
				{Name: "status", Type: "string", Required: false},
				{Name: "min_price", Type: "integer", Required: false},
				{Name: "max_price", Type: "integer", Required: false},
			}},
			{Name: "add", Description: "Make new offer", Parameters: []ParameterInfo{
				{Name: "product_id", Type: "string", Required: true},
				{Name: "buyer_id", Type: "string", Required: true},
				{Name: "offered_price", Type: "integer", Required: true},
				{Name: "message", Type: "string", Required: false},
			}},
			{Name: "accept", Description: "Accept offer", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "reject", Description: "Reject offer", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "counter", Description: "Make counter offer", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "counter_price", Type: "integer", Required: true},
				{Name: "counter_message", Type: "string", Required: false},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.ProductType), Type: "belongs_to", Description: "Offer belongs to a product"},
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Offer involves buyer and seller"},
		},
		SearchableFields: []string{"message", "counter_message"},
		SortableFields:   []string{"offered_price", "created_at", "expires_at"},
		FilterableFields: []string{"status", "product_id", "buyer_id", "seller_id"},
	}
}

// Review schema initialization
func (r *LLMSchemaRegistry) initializeReviewSchema() {
	r.schemas[models.ReviewType] = &LLMOperationSchema{
		EntityType:  models.ReviewType,
		TableName:   "reviews",
		Description: "Manages product and seller reviews and ratings",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique review identifier"},
			{Name: "reviewer_id", Type: "string", Required: true, Description: "User ID writing the review"},
			{Name: "reviewee_id", Type: "string", Required: false, Description: "User ID being reviewed (for seller reviews)"},
			{Name: "product_id", Type: "string", Required: false, Description: "Product ID being reviewed"},
			{Name: "order_id", Type: "string", Required: false, Description: "Associated order ID"},
			{Name: "rating", Type: "integer", Required: true, Description: "Rating score (1-5 stars)"},
			{Name: "title", Type: "string", Required: false, Description: "Review title or summary"},
			{Name: "content", Type: "string", Required: true, Description: "Detailed review content"},
			{Name: "review_type", Type: "string", Required: true, Description: "Type of review (product, seller, service)"},
			{Name: "verified_purchase", Type: "boolean", Required: false, Description: "Whether review is from verified purchase"},
			{Name: "helpful_count", Type: "integer", Required: false, Description: "Number of helpful votes"},
			{Name: "images", Type: "array", Required: false, Description: "Review image URLs"},
			{Name: "response", Type: "string", Required: false, Description: "Seller or admin response"},
			{Name: "response_date", Type: "datetime", Required: false, Description: "Response timestamp"},
			{Name: "is_featured", Type: "boolean", Required: false, Description: "Whether review is featured"},
			{Name: "is_verified", Type: "boolean", Required: false, Description: "Whether review is verified by admin"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Review creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find review by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search reviews by content", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter reviews by criteria", Parameters: []ParameterInfo{
				{Name: "product_id", Type: "string", Required: false},
				{Name: "reviewer_id", Type: "string", Required: false},
				{Name: "reviewee_id", Type: "string", Required: false},
				{Name: "review_type", Type: "string", Required: false},
				{Name: "min_rating", Type: "integer", Required: false},
				{Name: "max_rating", Type: "integer", Required: false},
				{Name: "verified_purchase", Type: "boolean", Required: false},
			}},
			{Name: "add", Description: "Create new review", Parameters: []ParameterInfo{
				{Name: "reviewer_id", Type: "string", Required: true},
				{Name: "rating", Type: "integer", Required: true},
				{Name: "content", Type: "string", Required: true},
				{Name: "review_type", Type: "string", Required: true},
			}},
			{Name: "update", Description: "Update review content", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "remove", Description: "Delete review", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "respond", Description: "Add response to review", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "response", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Review belongs to reviewer and reviewee"},
			{TargetEntity: string(models.ProductType), Type: "belongs_to", Description: "Review belongs to a product"},
			{TargetEntity: string(models.EntityTypeOrder), Type: "belongs_to", Description: "Review may belong to an order"},
		},
		SearchableFields: []string{"title", "content", "response"},
		SortableFields:   []string{"rating", "helpful_count", "created_at"},
		FilterableFields: []string{"review_type", "rating", "verified_purchase", "is_featured", "is_verified"},
	}
}

// Comment schema initialization
func (r *LLMSchemaRegistry) initializeCommentSchema() {
	r.schemas[models.CommentType] = &LLMOperationSchema{
		EntityType:  models.CommentType,
		TableName:   "comments",
		Description: "Manages comments on products, reviews, and other content",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique comment identifier"},
			{Name: "content", Type: "string", Required: true, Description: "Comment text content"},
			{Name: "author_id", Type: "string", Required: true, Description: "User ID of comment author"},
			{Name: "parent_type", Type: "string", Required: true, Description: "Type of parent entity (product, review, etc.)"},
			{Name: "parent_id", Type: "string", Required: true, Description: "ID of parent entity"},
			{Name: "parent_comment_id", Type: "string", Required: false, Description: "ID of parent comment (for replies)"},
			{Name: "is_verified", Type: "boolean", Required: false, Description: "Whether comment is verified"},
			{Name: "helpful_count", Type: "integer", Required: false, Description: "Number of helpful votes"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Comment creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find comment by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search comments by content", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
			{Name: "add", Description: "Create comment", Parameters: []ParameterInfo{
				{Name: "content", Type: "string", Required: true},
				{Name: "author_id", Type: "string", Required: true},
				{Name: "parent_type", Type: "string", Required: true},
				{Name: "parent_id", Type: "string", Required: true},
			}},
		},
		SearchableFields: []string{"content"},
		SortableFields:   []string{"created_at", "helpful_count"},
		FilterableFields: []string{"parent_type", "author_id", "is_verified"},
	}
}

// Notification schema initialization
func (r *LLMSchemaRegistry) initializeNotificationSchema() {
	r.schemas[models.NotificationEntityType] = &LLMOperationSchema{
		EntityType:  models.NotificationEntityType,
		TableName:   "notifications",
		Description: "Manages user notifications and alerts",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique notification identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "Target user ID"},
			{Name: "type", Type: "string", Required: true, Description: "Notification type"},
			{Name: "title", Type: "string", Required: true, Description: "Notification title"},
			{Name: "message", Type: "string", Required: true, Description: "Notification message"},
			{Name: "is_read", Type: "boolean", Required: false, Description: "Whether notification is read"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find notification by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "filter", Description: "Filter notifications", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: false},
				{Name: "type", Type: "string", Required: false},
				{Name: "is_read", Type: "boolean", Required: false},
			}},
		},
		SearchableFields: []string{"title", "message"},
		FilterableFields: []string{"type", "is_read", "user_id"},
	}
}

// Category schema initialization
func (r *LLMSchemaRegistry) initializeCategorySchema() {
	r.schemas[models.CategoryType] = &LLMOperationSchema{
		EntityType:  models.CategoryType,
		TableName:   "categories",
		Description: "Manages product and content categories",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique category identifier"},
			{Name: "name", Type: "string", Required: true, Description: "Category name"},
			{Name: "slug", Type: "string", Required: true, Description: "URL-friendly category slug"},
			{Name: "description", Type: "string", Required: false, Description: "Category description"},
			{Name: "parent_id", Type: "string", Required: false, Description: "Parent category ID"},
			{Name: "is_active", Type: "boolean", Required: true, Description: "Whether category is active"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find category by ID", Parameters: []ParameterInfo{{Name: "id", Type: "string", Required: true}}},
			{Name: "search", Description: "Search categories", Parameters: []ParameterInfo{{Name: "search_term", Type: "string", Required: true}}},
		},
		SearchableFields: []string{"name", "description"},
		FilterableFields: []string{"is_active", "parent_id"},
	}
}

// Newsletter gRPC repository schema
func (r *LLMSchemaRegistry) initializeNewsletterSchema() {
	r.schemas[models.NewsletterEntityType] = &LLMOperationSchema{
		EntityType:  models.NewsletterEntityType,
		TableName:   "newsletter_service",
		Description: "Newsletter subscription and email campaign management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique newsletter identifier"},
			{Name: "email", Type: "string", Required: true, Description: "Subscriber email address"},
			{Name: "user_id", Type: "string", Required: false, Description: "Associated user ID"},
			{Name: "subscription_status", Type: "string", Required: true, Description: "Subscription status (active, unsubscribed, bounced)"},
			{Name: "categories", Type: "array", Required: false, Description: "Newsletter categories subscribed to"},
			{Name: "preferences", Type: "object", Required: false, Description: "Email preferences and settings"},
			{Name: "subscribed_at", Type: "datetime", Required: false, Description: "Subscription timestamp"},
			{Name: "unsubscribed_at", Type: "datetime", Required: false, Description: "Unsubscription timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "subscribe", Description: "Subscribe email to newsletter", Parameters: []ParameterInfo{
				{Name: "email", Type: "string", Required: true},
				{Name: "categories", Type: "array", Required: false},
			}},
			{Name: "unsubscribe", Description: "Unsubscribe from newsletter", Parameters: []ParameterInfo{
				{Name: "email", Type: "string", Required: true},
			}},
			{Name: "find", Description: "Find subscription by email or ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: false},
				{Name: "email", Type: "string", Required: false},
			}},
			{Name: "filter", Description: "Filter subscriptions by status", Parameters: []ParameterInfo{
				{Name: "status", Type: "string", Required: false},
				{Name: "category", Type: "string", Required: false},
			}},
		},
		SearchableFields: []string{"email"},
		FilterableFields: []string{"subscription_status", "categories"},
	}
}

// Basket gRPC repository schema
func (r *LLMSchemaRegistry) initializeBasketSchema() {
	r.schemas[models.BasketType] = &LLMOperationSchema{
		EntityType:  models.BasketType,
		TableName:   "basket_service",
		Description: "Shopping cart and basket management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique basket identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "Owner user ID"},
			{Name: "session_id", Type: "string", Required: false, Description: "Session identifier for guest users"},
			{Name: "items", Type: "array", Required: true, Description: "List of basket items"},
			{Name: "total_amount", Type: "integer", Required: false, Description: "Total basket value in cents"},
			{Name: "item_count", Type: "integer", Required: false, Description: "Number of items in basket"},
			{Name: "currency", Type: "string", Required: false, Description: "Currency code"},
			{Name: "expires_at", Type: "datetime", Required: false, Description: "Basket expiration time"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Basket creation timestamp"},
			{Name: "updated_at", Type: "datetime", Required: false, Description: "Last update timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "find", Description: "Find basket by user or session", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: false},
				{Name: "session_id", Type: "string", Required: false},
			}},
			{Name: "add_item", Description: "Add item to basket", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "product_id", Type: "string", Required: true},
				{Name: "quantity", Type: "integer", Required: true},
			}},
			{Name: "remove_item", Description: "Remove item from basket", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "product_id", Type: "string", Required: true},
			}},
			{Name: "clear", Description: "Clear entire basket", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "calculate_total", Description: "Calculate basket total", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Basket belongs to a user"},
			{TargetEntity: string(models.ProductType), Type: "has_many", Description: "Basket contains multiple products"},
		},
		FilterableFields: []string{"user_id", "session_id"},
	}
}

// Metric gRPC repository schema
func (r *LLMSchemaRegistry) initializeMetricSchema() {
	r.schemas[models.MetricEntityType] = &LLMOperationSchema{
		EntityType:  models.MetricEntityType,
		TableName:   "metrics_service",
		Description: "Analytics and performance metrics collection via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique metric identifier"},
			{Name: "metric_name", Type: "string", Required: true, Description: "Name of the metric"},
			{Name: "metric_type", Type: "string", Required: true, Description: "Type of metric (counter, gauge, histogram)"},
			{Name: "value", Type: "float", Required: true, Description: "Metric value"},
			{Name: "labels", Type: "object", Required: false, Description: "Metric labels and tags"},
			{Name: "entity_type", Type: "string", Required: false, Description: "Associated entity type"},
			{Name: "entity_id", Type: "string", Required: false, Description: "Associated entity ID"},
			{Name: "user_id", Type: "string", Required: false, Description: "Associated user ID"},
			{Name: "timestamp", Type: "datetime", Required: true, Description: "Metric collection timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "record", Description: "Record a metric value", Parameters: []ParameterInfo{
				{Name: "metric_name", Type: "string", Required: true},
				{Name: "value", Type: "float", Required: true},
				{Name: "labels", Type: "object", Required: false},
			}},
			{Name: "query", Description: "Query metrics by name and time range", Parameters: []ParameterInfo{
				{Name: "metric_name", Type: "string", Required: true},
				{Name: "start_time", Type: "datetime", Required: false},
				{Name: "end_time", Type: "datetime", Required: false},
			}},
			{Name: "aggregate", Description: "Aggregate metrics by time period", Parameters: []ParameterInfo{
				{Name: "metric_name", Type: "string", Required: true},
				{Name: "aggregation", Type: "string", Required: true},
				{Name: "period", Type: "string", Required: true},
			}},
		},
		SearchableFields: []string{"metric_name"},
		FilterableFields: []string{"metric_type", "entity_type", "user_id"},
	}
}

// Messages gRPC repository schema
func (r *LLMSchemaRegistry) initializeMessagesSchema() {
	r.schemas[models.MessageType] = &LLMOperationSchema{
		EntityType:  models.MessageType,
		TableName:   "messaging_service",
		Description: "User-to-user messaging system via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique message identifier"},
			{Name: "sender_id", Type: "string", Required: true, Description: "Sender user ID"},
			{Name: "recipient_id", Type: "string", Required: true, Description: "Recipient user ID"},
			{Name: "conversation_id", Type: "string", Required: false, Description: "Conversation thread ID"},
			{Name: "subject", Type: "string", Required: false, Description: "Message subject"},
			{Name: "content", Type: "string", Required: true, Description: "Message content"},
			{Name: "message_type", Type: "string", Required: true, Description: "Message type (text, system, notification)"},
			{Name: "is_read", Type: "boolean", Required: false, Description: "Whether message has been read"},
			{Name: "attachments", Type: "array", Required: false, Description: "Message attachments"},
			{Name: "sent_at", Type: "datetime", Required: false, Description: "Message sent timestamp"},
			{Name: "read_at", Type: "datetime", Required: false, Description: "Message read timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "send", Description: "Send a message", Parameters: []ParameterInfo{
				{Name: "sender_id", Type: "string", Required: true},
				{Name: "recipient_id", Type: "string", Required: true},
				{Name: "content", Type: "string", Required: true},
				{Name: "subject", Type: "string", Required: false},
			}},
			{Name: "find", Description: "Find message by ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
			}},
			{Name: "get_conversation", Description: "Get conversation between users", Parameters: []ParameterInfo{
				{Name: "user1_id", Type: "string", Required: true},
				{Name: "user2_id", Type: "string", Required: true},
			}},
			{Name: "mark_read", Description: "Mark message as read", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "get_inbox", Description: "Get user's inbox", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "page", Type: "integer", Required: false},
				{Name: "page_size", Type: "integer", Required: false},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Message belongs to sender and recipient"},
		},
		SearchableFields: []string{"subject", "content"},
		FilterableFields: []string{"sender_id", "recipient_id", "is_read", "message_type"},
	}
}

// Wishlist gRPC repository schema
func (r *LLMSchemaRegistry) initializeWishlistSchema() {
	r.schemas[models.WishlistType] = &LLMOperationSchema{
		EntityType:  models.WishlistType,
		TableName:   "wishlist_service",
		Description: "User wishlist and favorites management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique wishlist identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "Owner user ID"},
			{Name: "name", Type: "string", Required: true, Description: "Wishlist name"},
			{Name: "description", Type: "string", Required: false, Description: "Wishlist description"},
			{Name: "is_public", Type: "boolean", Required: false, Description: "Whether wishlist is public"},
			{Name: "items", Type: "array", Required: false, Description: "List of wishlist items"},
			{Name: "item_count", Type: "integer", Required: false, Description: "Number of items in wishlist"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Wishlist creation timestamp"},
			{Name: "updated_at", Type: "datetime", Required: false, Description: "Last update timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "create", Description: "Create new wishlist", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "name", Type: "string", Required: true},
				{Name: "description", Type: "string", Required: false},
			}},
			{Name: "find", Description: "Find wishlist by ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
			}},
			{Name: "get_user_wishlists", Description: "Get all wishlists for user", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "add_item", Description: "Add item to wishlist", Parameters: []ParameterInfo{
				{Name: "wishlist_id", Type: "string", Required: true},
				{Name: "product_id", Type: "string", Required: true},
			}},
			{Name: "remove_item", Description: "Remove item from wishlist", Parameters: []ParameterInfo{
				{Name: "wishlist_id", Type: "string", Required: true},
				{Name: "product_id", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Wishlist belongs to a user"},
			{TargetEntity: string(models.ProductType), Type: "has_many", Description: "Wishlist contains multiple products"},
		},
		FilterableFields: []string{"user_id", "is_public"},
	}
}

// Following gRPC repository schema
func (r *LLMSchemaRegistry) initializeFollowingSchema() {
	r.schemas[models.FollowingType] = &LLMOperationSchema{
		EntityType:  models.FollowingType,
		TableName:   "following_service",
		Description: "User following and social connections via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique following relationship identifier"},
			{Name: "follower_id", Type: "string", Required: true, Description: "User ID who is following"},
			{Name: "following_id", Type: "string", Required: true, Description: "User ID being followed"},
			{Name: "status", Type: "string", Required: true, Description: "Relationship status (active, blocked, pending)"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Following relationship start timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "follow", Description: "Follow a user", Parameters: []ParameterInfo{
				{Name: "follower_id", Type: "string", Required: true},
				{Name: "following_id", Type: "string", Required: true},
			}},
			{Name: "unfollow", Description: "Unfollow a user", Parameters: []ParameterInfo{
				{Name: "follower_id", Type: "string", Required: true},
				{Name: "following_id", Type: "string", Required: true},
			}},
			{Name: "get_followers", Description: "Get user's followers", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "get_following", Description: "Get users that user is following", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "check_following", Description: "Check if user follows another user", Parameters: []ParameterInfo{
				{Name: "follower_id", Type: "string", Required: true},
				{Name: "following_id", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Following relationship involves two users"},
		},
		FilterableFields: []string{"follower_id", "following_id", "status"},
	}
}

// Activity gRPC repository schema
func (r *LLMSchemaRegistry) initializeActivitySchema() {
	r.schemas[models.ActivityType] = &LLMOperationSchema{
		EntityType:  models.ActivityType,
		TableName:   "activity_service",
		Description: "User activity tracking and feed management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique activity identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "User performing the activity"},
			{Name: "activity_type", Type: "string", Required: true, Description: "Type of activity (purchase, review, follow, etc.)"},
			{Name: "entity_type", Type: "string", Required: false, Description: "Type of entity involved"},
			{Name: "entity_id", Type: "string", Required: false, Description: "ID of entity involved"},
			{Name: "description", Type: "string", Required: true, Description: "Activity description"},
			{Name: "metadata", Type: "object", Required: false, Description: "Additional activity metadata"},
			{Name: "is_public", Type: "boolean", Required: false, Description: "Whether activity is public"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Activity timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "record", Description: "Record user activity", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "activity_type", Type: "string", Required: true},
				{Name: "description", Type: "string", Required: true},
				{Name: "entity_type", Type: "string", Required: false},
				{Name: "entity_id", Type: "string", Required: false},
			}},
			{Name: "get_user_activity", Description: "Get user's activity feed", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "limit", Type: "integer", Required: false},
			}},
			{Name: "get_feed", Description: "Get activity feed for user's network", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "limit", Type: "integer", Required: false},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Activity belongs to a user"},
		},
		FilterableFields: []string{"user_id", "activity_type", "entity_type", "is_public"},
	}
}

// Media gRPC repository schema
func (r *LLMSchemaRegistry) initializeMediaSchema() {
	r.schemas[models.MediaType] = &LLMOperationSchema{
		EntityType:  models.MediaType,
		TableName:   "media_service",
		Description: "Media file storage and management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique media identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "Owner user ID"},
			{Name: "filename", Type: "string", Required: true, Description: "Original filename"},
			{Name: "content_type", Type: "string", Required: true, Description: "MIME content type"},
			{Name: "file_size", Type: "integer", Required: true, Description: "File size in bytes"},
			{Name: "url", Type: "string", Required: true, Description: "Media access URL"},
			{Name: "thumbnail_url", Type: "string", Required: false, Description: "Thumbnail URL"},
			{Name: "entity_type", Type: "string", Required: false, Description: "Associated entity type"},
			{Name: "entity_id", Type: "string", Required: false, Description: "Associated entity ID"},
			{Name: "alt_text", Type: "string", Required: false, Description: "Alternative text for accessibility"},
			{Name: "is_public", Type: "boolean", Required: false, Description: "Whether media is publicly accessible"},
			{Name: "uploaded_at", Type: "datetime", Required: false, Description: "Upload timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "upload", Description: "Upload media file", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "filename", Type: "string", Required: true},
				{Name: "content_type", Type: "string", Required: true},
			}},
			{Name: "find", Description: "Find media by ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
			}},
			{Name: "get_user_media", Description: "Get user's media files", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
			}},
			{Name: "delete", Description: "Delete media file", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "user_id", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Media belongs to a user"},
		},
		FilterableFields: []string{"user_id", "content_type", "entity_type", "is_public"},
	}
}

// Service gRPC repository schema
func (r *LLMSchemaRegistry) initializeServiceSchema() {
	r.schemas[models.ServiceType] = &LLMOperationSchema{
		EntityType:  models.ServiceType,
		TableName:   "service_listings",
		Description: "Service offerings and professional services via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique service identifier"},
			{Name: "provider_id", Type: "string", Required: true, Description: "Service provider user ID"},
			{Name: "title", Type: "string", Required: true, Description: "Service title"},
			{Name: "description", Type: "string", Required: true, Description: "Service description"},
			{Name: "category", Type: "string", Required: true, Description: "Service category"},
			{Name: "price_type", Type: "string", Required: true, Description: "Pricing type (fixed, hourly, custom)"},
			{Name: "price", Type: "integer", Required: false, Description: "Service price in cents"},
			{Name: "duration", Type: "integer", Required: false, Description: "Service duration in minutes"},
			{Name: "location_type", Type: "string", Required: true, Description: "Location type (remote, on-site, both)"},
			{Name: "skills", Type: "array", Required: false, Description: "Required skills"},
			{Name: "availability", Type: "object", Required: false, Description: "Provider availability"},
			{Name: "is_active", Type: "boolean", Required: true, Description: "Whether service is active"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Service creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "create", Description: "Create service listing", Parameters: []ParameterInfo{
				{Name: "provider_id", Type: "string", Required: true},
				{Name: "title", Type: "string", Required: true},
				{Name: "description", Type: "string", Required: true},
				{Name: "category", Type: "string", Required: true},
			}},
			{Name: "find", Description: "Find service by ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
			}},
			{Name: "search", Description: "Search services", Parameters: []ParameterInfo{
				{Name: "search_term", Type: "string", Required: true},
			}},
			{Name: "filter", Description: "Filter services by criteria", Parameters: []ParameterInfo{
				{Name: "category", Type: "string", Required: false},
				{Name: "location_type", Type: "string", Required: false},
				{Name: "price_type", Type: "string", Required: false},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Service belongs to a provider"},
		},
		SearchableFields: []string{"title", "description", "skills"},
		FilterableFields: []string{"category", "price_type", "location_type", "is_active"},
	}
}

// Shipping gRPC repository schema
func (r *LLMSchemaRegistry) initializeShippingSchema() {
	r.schemas[models.ShippingEntityType] = &LLMOperationSchema{
		EntityType:  models.ShippingEntityType,
		TableName:   "shipping_service",
		Description: "Shipping and logistics management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique shipping identifier"},
			{Name: "order_id", Type: "string", Required: true, Description: "Associated order ID"},
			{Name: "carrier", Type: "string", Required: true, Description: "Shipping carrier"},
			{Name: "service_type", Type: "string", Required: true, Description: "Service type (standard, express, overnight)"},
			{Name: "tracking_number", Type: "string", Required: false, Description: "Carrier tracking number"},
			{Name: "status", Type: "string", Required: true, Description: "Shipping status"},
			{Name: "estimated_delivery", Type: "datetime", Required: false, Description: "Estimated delivery date"},
			{Name: "actual_delivery", Type: "datetime", Required: false, Description: "Actual delivery date"},
			{Name: "shipping_cost", Type: "integer", Required: true, Description: "Shipping cost in cents"},
			{Name: "from_address", Type: "object", Required: true, Description: "Origin address"},
			{Name: "to_address", Type: "object", Required: true, Description: "Destination address"},
		},
		Operations: []OperationDefinition{
			{Name: "create_shipment", Description: "Create shipping label", Parameters: []ParameterInfo{
				{Name: "order_id", Type: "string", Required: true},
				{Name: "carrier", Type: "string", Required: true},
				{Name: "service_type", Type: "string", Required: true},
			}},
			{Name: "track", Description: "Track shipment", Parameters: []ParameterInfo{
				{Name: "tracking_number", Type: "string", Required: true},
			}},
			{Name: "update_status", Description: "Update shipping status", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "status", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.EntityTypeOrder), Type: "belongs_to", Description: "Shipping belongs to an order"},
		},
		FilterableFields: []string{"carrier", "status", "service_type"},
	}
}

// Support gRPC repository schema
func (r *LLMSchemaRegistry) initializeSupportSchema() {
	r.schemas[models.SupportEntityType] = &LLMOperationSchema{
		EntityType:  models.SupportEntityType,
		TableName:   "support_service",
		Description: "Customer support ticket management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique ticket identifier"},
			{Name: "user_id", Type: "string", Required: true, Description: "Customer user ID"},
			{Name: "agent_id", Type: "string", Required: false, Description: "Assigned support agent ID"},
			{Name: "subject", Type: "string", Required: true, Description: "Ticket subject"},
			{Name: "description", Type: "string", Required: true, Description: "Issue description"},
			{Name: "category", Type: "string", Required: true, Description: "Issue category"},
			{Name: "priority", Type: "string", Required: true, Description: "Ticket priority (low, medium, high, urgent)"},
			{Name: "status", Type: "string", Required: true, Description: "Ticket status (open, in_progress, resolved, closed)"},
			{Name: "tags", Type: "array", Required: false, Description: "Ticket tags"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Ticket creation timestamp"},
			{Name: "resolved_at", Type: "datetime", Required: false, Description: "Resolution timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "create_ticket", Description: "Create support ticket", Parameters: []ParameterInfo{
				{Name: "user_id", Type: "string", Required: true},
				{Name: "subject", Type: "string", Required: true},
				{Name: "description", Type: "string", Required: true},
				{Name: "category", Type: "string", Required: true},
			}},
			{Name: "find", Description: "Find ticket by ID", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
			}},
			{Name: "assign_agent", Description: "Assign ticket to agent", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "agent_id", Type: "string", Required: true},
			}},
			{Name: "update_status", Description: "Update ticket status", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "status", Type: "string", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.UserEntityType), Type: "belongs_to", Description: "Ticket belongs to customer and agent"},
		},
		SearchableFields: []string{"subject", "description"},
		FilterableFields: []string{"status", "priority", "category", "agent_id"},
	}
}

// Geocoding gRPC repository schema
func (r *LLMSchemaRegistry) initializeGeocodingSchema() {
	r.schemas[models.GeocodingType] = &LLMOperationSchema{
		EntityType:  models.GeocodingType,
		TableName:   "geocoding_service",
		Description: "Location and geocoding services via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique location identifier"},
			{Name: "address", Type: "string", Required: true, Description: "Full address string"},
			{Name: "latitude", Type: "float", Required: true, Description: "Latitude coordinate"},
			{Name: "longitude", Type: "float", Required: true, Description: "Longitude coordinate"},
			{Name: "street", Type: "string", Required: false, Description: "Street address"},
			{Name: "city", Type: "string", Required: false, Description: "City name"},
			{Name: "state", Type: "string", Required: false, Description: "State or province"},
			{Name: "postal_code", Type: "string", Required: false, Description: "Postal code"},
			{Name: "country", Type: "string", Required: false, Description: "Country name"},
			{Name: "accuracy", Type: "string", Required: false, Description: "Geocoding accuracy level"},
		},
		Operations: []OperationDefinition{
			{Name: "geocode", Description: "Convert address to coordinates", Parameters: []ParameterInfo{
				{Name: "address", Type: "string", Required: true},
			}},
			{Name: "reverse_geocode", Description: "Convert coordinates to address", Parameters: []ParameterInfo{
				{Name: "latitude", Type: "float", Required: true},
				{Name: "longitude", Type: "float", Required: true},
			}},
			{Name: "find_nearby", Description: "Find nearby locations", Parameters: []ParameterInfo{
				{Name: "latitude", Type: "float", Required: true},
				{Name: "longitude", Type: "float", Required: true},
				{Name: "radius", Type: "float", Required: true},
			}},
		},
		SearchableFields: []string{"address", "city", "state"},
		FilterableFields: []string{"country", "state", "city"},
	}
}

// Variant gRPC repository schema
func (r *LLMSchemaRegistry) initializeVariantSchema() {
	r.schemas[models.VariantEntityType] = &LLMOperationSchema{
		EntityType:  models.VariantEntityType,
		TableName:   "variant_service",
		Description: "Product variant management via gRPC",
		Fields: []FieldDefinition{
			{Name: "id", Type: "string", Required: true, Description: "Unique variant identifier"},
			{Name: "product_id", Type: "string", Required: true, Description: "Parent product ID"},
			{Name: "sku", Type: "string", Required: true, Description: "Variant SKU"},
			{Name: "name", Type: "string", Required: true, Description: "Variant name"},
			{Name: "price", Type: "integer", Required: true, Description: "Variant price in cents"},
			{Name: "attributes", Type: "object", Required: false, Description: "Variant attributes (size, color, etc.)"},
			{Name: "stock", Type: "integer", Required: false, Description: "Available stock"},
			{Name: "weight", Type: "integer", Required: false, Description: "Variant weight"},
			{Name: "dimensions", Type: "object", Required: false, Description: "Variant dimensions"},
			{Name: "is_active", Type: "boolean", Required: true, Description: "Whether variant is active"},
			{Name: "created_at", Type: "datetime", Required: false, Description: "Variant creation timestamp"},
		},
		Operations: []OperationDefinition{
			{Name: "create", Description: "Create product variant", Parameters: []ParameterInfo{
				{Name: "product_id", Type: "string", Required: true},
				{Name: "sku", Type: "string", Required: true},
				{Name: "name", Type: "string", Required: true},
				{Name: "price", Type: "integer", Required: true},
			}},
			{Name: "find", Description: "Find variant by ID or SKU", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: false},
				{Name: "sku", Type: "string", Required: false},
			}},
			{Name: "get_product_variants", Description: "Get all variants for product", Parameters: []ParameterInfo{
				{Name: "product_id", Type: "string", Required: true},
			}},
			{Name: "update_stock", Description: "Update variant stock", Parameters: []ParameterInfo{
				{Name: "id", Type: "string", Required: true},
				{Name: "stock", Type: "integer", Required: true},
			}},
		},
		Relationships: []RelationshipDefinition{
			{TargetEntity: string(models.ProductType), Type: "belongs_to", Description: "Variant belongs to a product"},
		},
		SearchableFields: []string{"name", "sku"},
		FilterableFields: []string{"product_id", "is_active"},
	}
}
