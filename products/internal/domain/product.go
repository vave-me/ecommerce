package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const ProductAggregate = "products.Product"

var (
	ErrProductNameIsBlank     = errors.Wrap(errors.ErrBadRequest, "the product name cannot be blank")
	ErrProductPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the product price cannot be negative")
	ErrNotAPriceIncrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Product struct {
	es.Aggregate
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        ProductCondition
	Model            string
	Tags             []string
	ManageStock      bool
	Stock            int64
	SKU              string
	Attributes       []Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           ProductStatus
	Negotiable       bool
	UserType         UserType // e.g. "private", "business"
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	URLReference     string
	MerchantName     string
	TypeOfOffer      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Product)(nil)

func NewProduct(id string) *Product {
	return &Product{
		Aggregate: es.NewAggregate(id, ProductAggregate),
	}
}

// Key implements registry.Registerable
func (Product) Key() string { return ProductAggregate }
func (p *Product) InitProduct(
	id, name, description string,
	basePrice int64,
	userSellerID, categoryID, categorySlug, brand string,
	condition ProductCondition, model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []Attribute,
	weight, height, width, depth int64,
	status ProductStatus,
	negotiable bool,
	userType UserType,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []Option,
	thumbnail string,
	lat, lng float64,
) (ddd.Event, error) {

	if name == "" {
		return nil, ErrProductNameIsBlank
	}
	if basePrice < 0 {
		return nil, ErrProductPriceIsNegative
	}

	p.AddEvent(ProductAddedEvent, &ProductAdded{
		Name:             name,
		Description:      description,
		BasePrice:        basePrice,
		UserSellerID:     userSellerID,
		CategoryID:       categoryID,
		CategorySlug:     categorySlug,
		Brand:            brand,
		Condition:        condition,
		Model:            model,
		Tags:             tags,
		ManageStock:      manageStock,
		Stock:            stock,
		SKU:              sku,
		Attributes:       attributes,
		Weight:           weight,
		Height:           height,
		Width:            width,
		Depth:            depth,
		Status:           status,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		ShippingCost:     shippingCost,
		HasVariants:      hasVariants,
		Options:          options,
		Thumbnail:        thumbnail,
		Lat:              lat,
		Lng:              lng,
	})
	return ddd.NewEvent(ProductAddedEvent, p), nil
}
func (p *Product) Update(
	name, description string,
	basePrice int64,
	userSellerID, categoryID, categorySlug, brand string,
	condition ProductCondition, model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []Attribute,
	weight, height, width, depth int64,
	status ProductStatus,
	thumbnail string,
	negotiable bool,
	userType UserType,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []Option,
	lat, lng float64,
) (ddd.Event, error) {

	if name == "" {
		name = p.Name
	}

	if description == "" {
		description = p.Description
	}
	if basePrice < 0 {
		basePrice = p.BasePrice
	}
	if userSellerID == "" {
		userSellerID = p.UserSellerID
	}
	if categoryID == "" {
		categoryID = p.CategoryID
	}
	if categorySlug == "" {
		categorySlug = p.CategorySlug
	}
	if brand == "" {
		brand = p.Brand
	}
	if condition == "" {
		condition = p.Condition
	}
	if model == "" {
		model = p.Model
	}
	if tags == nil {
		tags = p.Tags
	}
	if manageStock {
		manageStock = p.ManageStock
	}
	if attributes == nil {
		attributes = p.Attributes
	}

	if weight == 0 {
		weight = p.Weight
	}
	if height == 0 {
		height = p.Height
	}
	if width == 0 {
		width = p.Width
	}
	if depth == 0 {
		depth = p.Depth
	}
	if status == "" {
		status = p.Status
	}

	if stock == 0 {
		stock = p.Stock
	}
	if shippingCost == 0 {
		shippingCost = p.ShippingCost
	}
	if hasVariants {
		hasVariants = p.HasVariants
	}
	if negotiable {
		negotiable = p.Negotiable
	}
	if userType == "" {
		userType = p.UserType
	}
	if thumbnail == "" {
		thumbnail = p.Thumbnail
	}
	if attributes == nil {
		attributes = p.Attributes
	}
	if options == nil {
		options = p.Options
	}
	if lat == 0 {
		lat = p.Lat
	}
	if lng == 0 {
		lng = p.Lng
	}

	p.AddEvent(ProductUpdatedEvent, &ProductUpdated{
		Name:             name,
		Description:      description,
		BasePrice:        basePrice,
		UserSellerID:     userSellerID,
		CategoryID:       categoryID,
		CategorySlug:     categorySlug,
		Brand:            brand,
		Condition:        condition,
		Model:            model,
		Tags:             tags,
		ManageStock:      manageStock,
		Stock:            stock,
		SKU:              sku,
		Attributes:       attributes,
		Weight:           weight,
		Height:           height,
		Width:            width,
		Depth:            depth,
		Status:           status,
		Thumbnail:        thumbnail,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		ShippingCost:     shippingCost,
		HasVariants:      hasVariants,
		Options:          options,
		Lat:              lat,
		Lng:              lng,
	})
	return ddd.NewEvent(ProductUpdatedEvent, p), nil
}
func (p *Product) Rebrand(name, description string) (ddd.Event, error) {
	p.AddEvent(ProductRebrandedEvent, &ProductRebranded{
		Name:        name,
		Description: description,
		Brand:       p.Brand, // or pass new brand in param if you want
		Model:       p.Model,
	})
	return ddd.NewEvent(ProductRebrandedEvent, p), nil
}
func (p *Product) AddThumbnail(thumbnail string) (ddd.Event, error) {
	p.AddEvent(ProductThumbnailAddedEvent, &ProductThumbnailAdded{
		Thumbnail: thumbnail,
	})
	return ddd.NewEvent(ProductThumbnailAddedEvent, p), nil
}
func (p *Product) UpdateThumbnail(thumbnail string) (ddd.Event, error) {
	p.AddEvent(ProductThumbnailUpdatedEvent, &ProductThumbnailUpdated{
		Thumbnail: thumbnail,
	})
	return ddd.NewEvent(ProductThumbnailUpdatedEvent, p), nil
}

func (p *Product) IncreasePrice(newPrice int64) (ddd.Event, error) {
	if newPrice < p.BasePrice {
		return nil, ErrNotAPriceIncrease
	}
	p.AddEvent(ProductPriceIncreasedEvent, &ProductPriceIncreased{
		ProductID: p.ID(),
		OldPrice:  p.BasePrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(ProductPriceIncreasedEvent, p), nil
}

func (p *Product) DecreasePrice(newPrice int64) (ddd.Event, error) {
	if newPrice > p.BasePrice {
		return nil, ErrNotAPriceDecrease
	}
	p.AddEvent(ProductPriceDecreasedEvent, &ProductPriceDecreased{
		ProductID: p.ID(),
		OldPrice:  p.BasePrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(ProductPriceDecreasedEvent, p), nil
}

func (p *Product) Remove(id, userSellerID string) (ddd.Event, error) {
	p.AddEvent(ProductRemovedEvent, &ProductRemoved{
		ProductID:    id,
		UserSellerID: userSellerID,
	})
	return ddd.NewEvent(ProductRemovedEvent, p), nil
}

func (p *Product) AdjustStock(newStock int64) (ddd.Event, error) {
	oldStock := p.Stock
	p.AddEvent(ProductStockAdjustedEvent, &ProductStockAdjusted{
		ProductID: p.ID(),
		OldStock:  oldStock,
		NewStock:  newStock,
	})
	return ddd.NewEvent(ProductStockAdjustedEvent, p), nil
}

func (p *Product) Archive(userSellerID string) (ddd.Event, error) {
	p.AddEvent(ProductArchivedEvent, &ProductArchived{
		ProductID:    p.ID(),
		UserSellerID: userSellerID,
	})
	return ddd.NewEvent(ProductArchivedEvent, p), nil
}

func (p *Product) MarkSold(userSellerID string, finalPrice int64) (ddd.Event, error) {
	p.AddEvent(ProductSoldEvent, &ProductSold{
		ProductID:    p.ID(),
		UserSellerID: userSellerID,
		FinalPrice:   finalPrice,
	})
	return ddd.NewEvent(ProductSoldEvent, p), nil
}

func (p *Product) MarkLeased(userSellerID string, monthlyPrice int64, leaseTermMonths int64) (ddd.Event, error) {
	p.AddEvent(ProductLeasedEvent, &ProductLeased{
		ProductID:    p.ID(),
		UserSellerID: userSellerID,
	})
	return ddd.NewEvent(ProductLeasedEvent, p), nil
}

func (p *Product) MarkPawned(userSellerID string, lockedPrice, redemptionFee int64) (ddd.Event, error) {
	p.AddEvent(ProductPawnedEvent, &ProductPawned{
		ProductID:    p.ID(),
		UserSellerID: userSellerID,
	})
	return ddd.NewEvent(ProductPawnedEvent, p), nil
}

func (p *Product) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *ProductAdded:
		p.Name = e.Name
		p.Description = e.Description
		p.BasePrice = e.BasePrice
		p.UserSellerID = e.UserSellerID
		p.CategoryID = e.CategoryID
		p.CategorySlug = e.CategorySlug
		p.Brand = e.Brand
		p.Condition = e.Condition
		p.Model = e.Model
		p.Tags = e.Tags
		p.ManageStock = e.ManageStock
		p.Stock = e.Stock
		p.SKU = e.SKU
		p.Attributes = e.Attributes
		p.Weight = e.Weight
		p.Height = e.Height
		p.Width = e.Width
		p.Depth = e.Depth
		p.Status = e.Status
		p.Thumbnail = e.Thumbnail
		p.Negotiable = e.Negotiable
		p.UserType = e.UserType
		p.MiddlemanService = e.MiddlemanService
		p.ShippingCost = e.ShippingCost
		p.HasVariants = e.HasVariants
		p.Options = e.Options
		p.Lat = e.Lat
		p.Lng = e.Lng

	case *ProductUpdated:
		p.Name = e.Name
		p.Description = e.Description
		p.BasePrice = e.BasePrice
		p.UserSellerID = e.UserSellerID
		p.CategoryID = e.CategoryID
		p.CategorySlug = e.CategorySlug
		p.Brand = e.Brand
		p.Condition = e.Condition
		p.Model = e.Model
		p.Tags = e.Tags
		p.ManageStock = e.ManageStock
		p.Stock = e.Stock
		p.SKU = e.SKU
		p.Attributes = e.Attributes
		p.Weight = e.Weight
		p.Height = e.Height
		p.Width = e.Width
		p.Depth = e.Depth
		p.Status = e.Status
		p.Thumbnail = e.Thumbnail
		p.Negotiable = e.Negotiable
		p.UserType = e.UserType
		p.MiddlemanService = e.MiddlemanService
		p.ShippingCost = e.ShippingCost
		p.HasVariants = e.HasVariants
		p.Options = e.Options
		p.Lat = e.Lat
		p.Lng = e.Lng

	case *ProductThumbnailAdded:
		p.Thumbnail = e.Thumbnail
	case *ProductThumbnailUpdated:
		p.Thumbnail = e.Thumbnail
	case *ProductRebranded:
		p.Name = e.Name
		p.Description = e.Description
		if e.Brand != "" {
			p.Brand = e.Brand
		}
		if e.Model != "" {
			p.Model = e.Model
		}

	case *ProductPriceIncreased:
		p.BasePrice = e.NewPrice

	case *ProductPriceDecreased:
		p.BasePrice = e.NewPrice

	case *ProductStockAdjusted:
		p.Stock = e.NewStock

	case *ProductArchived:
		// p.Status = ProductStatusArchived or
		// p.IsActive = false, etc.

	case *ProductSold:
		// p.Status = ProductStatusSold or set final price, etc.

	case *ProductLeased:
		// p.Status = ProductStatusLeased or whatever logic

	case *ProductPawned:
		// p.Status = ProductStatusPawned or set locked price, etc.

	case *ProductRemoved:
		// Possibly set an internal "removed" flag or do nothing if the
		// aggregator is ephemeral

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			p, event.EventName(), e)
	}
	return nil
}

func (p *Product) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ProductV1:
		p.Name = ss.Name
		p.Description = ss.Description
		p.BasePrice = ss.BasePrice
		p.UserSellerID = ss.UserSellerID
		p.CategoryID = ss.CategoryID
		p.CategorySlug = ss.CategorySlug
		p.Brand = ss.Brand
		p.Condition = ss.Condition
		p.Model = ss.Model
		p.Tags = ss.Tags
		p.ManageStock = ss.ManageStock
		p.Stock = ss.Stock
		p.SKU = ss.SKU
		p.Weight = ss.Weight
		p.Height = ss.Height
		p.Width = ss.Width
		p.Depth = ss.Depth
		p.Status = ss.Status
		p.Thumbnail = ss.Thumbnail
		p.Negotiable = ss.Negotiable
		p.UserType = ss.UserType
		p.MiddlemanService = ss.MiddlemanService
		p.ShippingCost = ss.ShippingCost
		p.HasVariants = ss.HasVariants
		p.Options = ss.Options
		p.Lat = ss.Lat
		p.Lng = ss.Lng

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", p, snapshot)
	}
	return nil
}

func (p Product) ToSnapshot() es.Snapshot {
	return ProductV1{
		Name:             p.Name,
		Description:      p.Description,
		BasePrice:        p.BasePrice,
		UserSellerID:     p.UserSellerID,
		CategoryID:       p.CategoryID,
		CategorySlug:     p.CategorySlug,
		Brand:            p.Brand,
		Condition:        p.Condition,
		Model:            p.Model,
		Tags:             p.Tags,
		ManageStock:      p.ManageStock,
		Stock:            p.Stock,
		SKU:              p.SKU,
		Attributes:       p.Attributes,
		Weight:           p.Weight,
		Height:           p.Height,
		Width:            p.Width,
		Depth:            p.Depth,
		Status:           p.Status,
		Thumbnail:        p.Thumbnail,
		Negotiable:       p.Negotiable,
		UserType:         p.UserType,
		MiddlemanService: p.MiddlemanService,
		ShippingCost:     p.ShippingCost,
		HasVariants:      p.HasVariants,
		Options:          p.Options,
	}
}
