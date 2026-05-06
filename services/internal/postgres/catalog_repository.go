package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"middleman/internal/postgres"
	"middleman/services/internal/domain"

	"github.com/stackus/errors"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) AddService(
	ctx context.Context,
	id, name, description, serviceType string,
	basePrice int64, pricing []string, availability string,
	providerName, userID, categoryID, categorySlug string,
	descriptionShort, descriptionLong string,
	qualifications []string,
	contact, faq string,
	tags []string,
	status domain.ServiceStatus,
	userType domain.UserType,
	shippingCost int64,
	negotiable, hasVariants, middlemanService bool,
	attributes []domain.Attribute,
	options []domain.Option,
	thumbnail string,
	lat, lng float64,
) error {
	tagsStr := sliceToString(tags)
	pricingStr := sliceToString(pricing)
	qualificationsStr := sliceToString(qualifications)
	attrsStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	// We'll insert into the 'location' column using ST_SetSRID(ST_MakePoint(:lng, :lat),4326)
	const query = `
      INSERT INTO %s (
        id, name, description, service_type,
        base_price, pricing, availability,
        provider_name, user_id, category_id, category_slug,
        description_short, description_long,
        qualifications, contact, faq,
        tags, status, negotiable,
        user_type, middleman_service, shipping_cost,
        has_variants, attributes, options, thumbnail,
        location, lat, lng
      )
      VALUES (
        $1, $2, $3, $4,
        $5, $6, $7,
        $8, $9, $10, $11,
        $12, $13,
        $14, $15, $16,
        $17, $18, $19,
        $20, $21, $22,
        $23, $24, $25, $26,
        ST_SetSRID(ST_MakePoint($28, $27), 4326), $27, $28
      )
    `

	_, err := r.db.ExecContext(ctx, r.table(query),
		id,                // $1
		name,              // $2
		description,       // $3
		serviceType,       // $4
		basePrice,         // $5
		pricingStr,        // $6
		availability,      // $7
		providerName,      // $8
		userID,            // $9
		categoryID,        // $10
		categorySlug,      // $11
		descriptionShort,  // $12
		descriptionLong,   // $13
		qualificationsStr, // $14
		contact,           // $15
		faq,               // $16
		tagsStr,           // $17
		status,            // $18
		negotiable,        // $19
		userType,          // $20
		middlemanService,  // $21
		shippingCost,      // $22
		hasVariants,       // $23
		attrsStr,          // $24
		optsStr,           // $25
		thumbnail,         // $26
		lat,               // $27
		lng,               // $28
	)
	return err
}

func (r CatalogRepository) RebrandService(
	ctx context.Context,
	serviceID, name, description string,
	tags, qualifications []string,
	faq string,
) error {
	tagsStr := sliceToString(tags)
	qualificationsStr := sliceToString(qualifications)

	const query = `
      UPDATE %s
      SET name = $2,
          description = $3,
          tags = $4,
          qualifications = $5,
          faq = $6,
          updated_at = NOW()
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		serviceID,
		name,
		description,
		tagsStr,
		qualificationsStr,
		faq,
	)
	return err
}

func (r CatalogRepository) UpdateService(
	ctx context.Context,
	id, name, description, serviceType string,
	basePrice int64, pricing []string, availability string,
	providerName, userID, categoryID, categorySlug string,
	descriptionShort, descriptionLong string,
	qualifications []string,
	contact, faq string,
	tags []string,
	status domain.ServiceStatus,
	userType domain.UserType,
	shippingCost int64,
	negotiable, hasVariants, middlemanService bool,
	attributes []domain.Attribute,
	options []domain.Option,
	thumbnail string,
	lat, lng float64,
) error {
	tagsStr := sliceToString(tags)
	pricingStr := sliceToString(pricing)
	qualificationsStr := sliceToString(qualifications)
	attrsStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	// Update statement with all fields from the domain struct
	const query = `
      UPDATE %s
      SET
        name                = $2,
        description         = $3,
        service_type        = $4,
        base_price          = $5,
        pricing             = $6,
        availability        = $7,
        provider_name       = $8,
        category_id         = $10,
        category_slug       = $11,
        description_short   = $12,
        description_long    = $13,
        qualifications      = $14,
        contact             = $15,
        faq                 = $16,
        tags                = $17,
        status              = $18,
        negotiable          = $19,
        user_type       = $20,
        middleman_service   = $21,
        shipping_cost       = $22,
        has_variants        = $23,
        attributes          = $24,
        options             = $25,
        thumbnail           = $26,
        location            = ST_SetSRID(ST_MakePoint($28, $27), 4326),
        lat                 = $27,
        lng                 = $28,
        updated_at          = NOW()
      WHERE id = $1
        AND user_id = $9
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,                // $1
		name,              // $2
		description,       // $3
		serviceType,       // $4
		basePrice,         // $5
		pricingStr,        // $6
		availability,      // $7
		providerName,      // $8
		userID,            // $9
		categoryID,        // $10
		categorySlug,      // $11
		descriptionShort,  // $12
		descriptionLong,   // $13
		qualificationsStr, // $14
		contact,           // $15
		faq,               // $16
		tagsStr,           // $17
		status,            // $18
		negotiable,        // $19
		userType,          // $20
		middlemanService,  // $21
		shippingCost,      // $22
		hasVariants,       // $23
		attrsStr,          // $24
		optsStr,           // $25
		thumbnail,         // $26
		lat,               // $27
		lng,               // $28
	)
	return err
}

func (r CatalogRepository) UpdatePrice(ctx context.Context, serviceID string, oldPrice, newPrice int64) error {
	const query = `
      UPDATE %s
      SET base_price = $3,
          updated_at = NOW()
      WHERE id = $1
        AND base_price = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), serviceID, oldPrice, newPrice)
	if err != nil {
		return errors.Wrap(err, "updating price")
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "no service updated; either service not found or base_price mismatch")
	}
	return nil
}

func (r CatalogRepository) AdjustStock(
	ctx context.Context,
	serviceID string,
	userID string,
	oldStock, newStock int64,
) error {
	// Note: This seems like it doesn't apply to services but keeping for interface compatibility
	const query = `
      UPDATE %s
      SET updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID)
	if err != nil {
		return errors.Wrap(err, "adjusting service")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.Wrap(err, "no service updated; service not found")
	}
	return nil
}

func (r CatalogRepository) ToggleNegotiable(
	ctx context.Context,
	serviceID string,
	userID string,
	currentValue bool,
) error {
	const query = `
      UPDATE %s
      SET negotiable = NOT negotiable,
          updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
        AND negotiable = $3
    `
	res, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID, currentValue)
	if err != nil {
		return errors.Wrap(err, "toggle negotiable error")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.Wrap(err, "no service updated; negotiable mismatch or not found")
	}
	return nil
}

func (r CatalogRepository) RemoveService(
	ctx context.Context,
	serviceID string,
	userID string,
) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID)
	return err
}

func (r CatalogRepository) MarkServicePawned(
	ctx context.Context,
	serviceID string,
	userID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'pawned',
          updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID)
	return err
}

func (r CatalogRepository) MarkServiceLeased(
	ctx context.Context,
	serviceID string,
	userID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'leased',
          updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID)
	return err
}

func (r CatalogRepository) MarkServiceSold(
	ctx context.Context,
	serviceID string,
	userID string,
	finalPrice int64,
) error {
	const query = `
      UPDATE %s
      SET status = 'sold',
          base_price = $3,
          updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID, finalPrice)
	return err
}

func (r CatalogRepository) ArchiveService(
	ctx context.Context,
	serviceID string,
	userID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'archived',
          updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), serviceID, userID)
	return err
}

func (r CatalogRepository) Find(ctx context.Context, serviceID string) (*domain.CatalogService, error) {
	const query = `
      SELECT
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), serviceID)
	cp := &domain.CatalogService{ID: serviceID}

	var (
		tagsStr           string
		pricingStr        string
		qualificationsStr string
		attrsStr          string
		optsStr           string
	)

	err := row.Scan(
		&cp.Name,
		&cp.Description,
		&cp.ServiceType,
		&cp.BasePrice,
		&pricingStr,
		&cp.Availability,
		&cp.ProviderName,
		&cp.UserID,
		&cp.CategoryID,
		&cp.CategorySlug,
		&cp.DescriptionShort,
		&cp.DescriptionLong,
		&qualificationsStr,
		&cp.Contact,
		&cp.Faq,
		&tagsStr,
		&cp.Status,
		&cp.UserType,
		&cp.ShippingCost,
		&cp.Negotiable,
		&cp.HasVariants,
		&cp.MiddlemanService,
		&attrsStr,
		&optsStr,
		&cp.Thumbnail,
		&cp.Lat,
		&cp.Lng,
	)
	if err != nil {
		return nil, errors.Wrap(err, "scanning service")
	}

	cp.Tags = stringToSlice(tagsStr)
	cp.Pricing = stringToSlice(pricingStr)
	cp.Qualifications = stringToSlice(qualificationsStr)
	cp.Attributes = stringToAttributes(attrsStr)
	cp.Options = stringToOptions(optsStr)

	return cp, nil
}

func (r CatalogRepository) GetCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        category_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE user_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying catalog")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{UserID: userID}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final service rows")
	}

	return services, totalCount, nil
}
func (r CatalogRepository) GetPublicCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        category_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE user_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying catalog")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{UserID: userID}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final service rows")
	}

	return services, totalCount, nil
}

func (r CatalogRepository) GetServices(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      ORDER BY %s %s
      LIMIT $1
      OFFSET $2
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying all services")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.UserID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing service rows")
	}

	return services, totalCount, nil
}

func (r CatalogRepository) GetServicesByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_id = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, categoryID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE category_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categoryID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying services by category")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{CategoryID: categoryID}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.UserID,
			&cp.CategorySlug,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row by category")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return services, totalCount, nil
}
func (r CatalogRepository) GetServicesByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_slug = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, categorySlug).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_id,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE category_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categorySlug, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying services by category")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{CategorySlug: categorySlug}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.UserID,
			&cp.CategoryID,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row by category")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return services, totalCount, nil
}

func (r CatalogRepository) GetServicesWithFilterSlug(ctx context.Context, categorySlug, categoryID string, serviceType string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogService, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_slug = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, categorySlug).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting services by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_id,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng
      FROM %s
      WHERE category_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categorySlug, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying services by category")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{CategorySlug: categorySlug}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.UserID,
			&cp.CategoryID,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning service row by category")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return services, totalCount, nil
}

func (r CatalogRepository) FindByLocation(
	ctx context.Context,
	lat, lng float64,
	radiusMeters float64,
	limit int,
) ([]*domain.CatalogService, error) {
	// We SELECT all columns, plus distance for sorting
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        service_type,
        base_price,
        pricing,
        availability,
        provider_name,
        user_id,
        category_id,
        category_slug,
        description_short,
        description_long,
        qualifications,
        contact,
        faq,
        tags,
        status,
        user_type,
        shipping_cost,
        negotiable,
        has_variants,
        middleman_service,
        attributes,
        options,
        thumbnail,
        lat,
        lng,
        ST_Distance(
          location,
          ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
        ) AS dist
      FROM %s
      WHERE ST_DWithin(
        location,
        ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
        $3
      )
      ORDER BY dist
      LIMIT $4
    `, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, lng, lat, radiusMeters, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying services by location")
	}
	defer rows.Close()

	var services []*domain.CatalogService
	for rows.Next() {
		cp := &domain.CatalogService{}
		var (
			tagsStr           string
			pricingStr        string
			qualificationsStr string
			attrsStr          string
			optsStr           string
			dist              float64 // distance in meters from (lat,lng)
		)

		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.ServiceType,
			&cp.BasePrice,
			&pricingStr,
			&cp.Availability,
			&cp.ProviderName,
			&cp.UserID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.DescriptionShort,
			&cp.DescriptionLong,
			&qualificationsStr,
			&cp.Contact,
			&cp.Faq,
			&tagsStr,
			&cp.Status,
			&cp.UserType,
			&cp.ShippingCost,
			&cp.Negotiable,
			&cp.HasVariants,
			&cp.MiddlemanService,
			&attrsStr,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
			&dist,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning service row by location")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Pricing = stringToSlice(pricingStr)
		cp.Qualifications = stringToSlice(qualificationsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.Options = stringToOptions(optsStr)

		// Optionally store dist in the domain struct (if you add a Dist field)
		// cp.Dist = dist

		services = append(services, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing service rows by location")
	}

	return services, nil
}

// GetServicesWithFilter now matches the target interface signature
func (r CatalogRepository) GetServicesWithFilter(
	ctx context.Context,
	categoryID string, categorySlug string, serviceType string,
	userID string,
	status domain.ServiceStatus, // Assuming domain.ServiceStatus is the type for ServiceStatus
	searchText string,
	minPrice int64, maxPrice int64,
	lat float64, lng float64, radius float64, // radius is now float64
	availableFrom time.Time, availableTo time.Time,
	hasVariants bool, negotiable bool, middlemanService bool, // Direct names from interface
	userType domain.UserType, // Assuming domain.UserType is the type for UserType
	tags []string, qualifications []string,
	page int64, pageSize int64,
	sortBy string, sortOrder string,
) ([]*domain.CatalogService, int64, error) {

	whereClauses := []string{"1=1"}
	argsForWhere := make([]interface{}, 0, 32)
	argIndex := 1

	addClause := func(sqlClause string, arg interface{}) {
		whereClauses = append(whereClauses, sqlClause)
		argsForWhere = append(argsForWhere, arg)
		argIndex++
	}
	addSimpleClause := func(sqlClause string) {
		whereClauses = append(whereClauses, sqlClause)
	}

	if categoryID != "" {
		addClause(fmt.Sprintf("s.category_id = $%d", argIndex), categoryID)
	}
	if categorySlug != "" {
		addClause(fmt.Sprintf("s.category_slug = $%d", argIndex), categorySlug)
	}
	if serviceType != "" {
		addClause(fmt.Sprintf("s.service_type = $%d", argIndex), serviceType)
	}
	if userID != "" {
		addClause(fmt.Sprintf("s.user_id = $%d", argIndex), userID)
	}
	if status != "" { // Ensure empty string is a valid "no filter" for status
		addClause(fmt.Sprintf("s.status = $%d", argIndex), string(status))
	}
	if searchText != "" {
		searchPattern := "%" + strings.ToLower(searchText) + "%"
		clause := fmt.Sprintf("(LOWER(s.name) ILIKE $%d OR LOWER(s.description) ILIKE $%d)", argIndex, argIndex+1)
		whereClauses = append(whereClauses, clause)
		argsForWhere = append(argsForWhere, searchPattern, searchPattern)
		argIndex += 2
	}
	if minPrice > 0 {
		addClause(fmt.Sprintf("s.base_price >= $%d", argIndex), minPrice)
	}
	if maxPrice > 0 {
		addClause(fmt.Sprintf("s.base_price <= $%d", argIndex), maxPrice)
	}

	// Geo-spatial filter (only adds WHERE condition, does not select distance)
	if radius > 0 && lat != 0 && lng != 0 {
		clause := fmt.Sprintf(`
			ST_DWithin(
				s.location,
				ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography,
				$%d
			)
		`, argIndex, argIndex+1, argIndex+2) // lng, lat, radius order for ST_MakePoint
		whereClauses = append(whereClauses, clause)
		argsForWhere = append(argsForWhere, lng, lat, radius) // radius is float64
		argIndex += 3
	}

	if !availableFrom.IsZero() {
		addClause(fmt.Sprintf("s.availability >= $%d", argIndex), availableFrom)
	}
	if !availableTo.IsZero() {
		addClause(fmt.Sprintf("s.availability <= $%d", argIndex), availableTo)
	}

	if hasVariants {
		addSimpleClause("s.has_variants = TRUE")
	}
	if negotiable { // Direct parameter name from interface
		addSimpleClause("s.negotiable = TRUE")
	}
	if middlemanService { // Direct parameter name from interface
		addSimpleClause("s.middleman_service = TRUE")
	}
	if userType != "" { // Ensure empty string is a valid "no filter"
		addClause(fmt.Sprintf("s.user_type = $%d", argIndex), string(userType))
	}

	// Tags & Qualifications: "Naive approach" ILIKE '%value%' per product example
	for _, t := range tags {
		if t != "" {
			addClause(fmt.Sprintf("s.tags ILIKE $%d", argIndex), "%"+t+"%")
		}
	}
	for _, q := range qualifications {
		if q != "" {
			addClause(fmt.Sprintf("s.qualifications ILIKE $%d", argIndex), "%"+q+"%")
		}
	}

	// Pagination: Uses page and pageSize only, as offsetInput/limitInput are not in interface
	var finalOffset, finalLimit int64
	if pageSize <= 0 {
		pageSize = 20 // Default page size if not specified or invalid
	}
	if page < 1 {
		page = 1
	}
	finalOffset = (page - 1) * pageSize
	finalLimit = pageSize

	// Sorting: "distance" is not a sortable field as it's not selected.
	validSort := map[string]bool{
		"name":       true,
		"base_price": true,
		"updated_at": true,
		"created_at": true, // Ensure these columns exist in your 'services' table
	}
	if !validSort[sortBy] {
		sortBy = "name"
	}
	if strings.ToLower(sortOrder) != "desc" {
		sortOrder = "asc"
	}

	whereSQL := "WHERE " + joinClauses(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s s %s", r.tableName, whereSQL)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, argsForWhere...).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "GetServicesWithFilter: counting total items")
	}

	if totalCount == 0 {
		return []*domain.CatalogService{}, 0, nil
	}

	argsForSelect := make([]interface{}, len(argsForWhere))
	copy(argsForSelect, argsForWhere)
	argsForSelect = append(argsForSelect, finalLimit, finalOffset)

	// Select columns based on domain.CatalogService struct fields
	// pricing, qualifications, tags, attributes, options are selected as strings.
	// lat, lng are selected. Distance is NOT selected.
	// created_at, updated_at are NOT selected as they are not in provided domain.CatalogService struct.
	selectColumns := `
		s.id, s.name, s.description, s.service_type, s.base_price,
		s.pricing AS pricing_str, s.availability, s.provider_name, s.user_id,
		s.category_id, s.category_slug, s.description_short,
		s.description_long, s.qualifications AS qualifications_str, s.contact, s.faq,
		s.tags AS tags_str, s.status, s.user_type, s.shipping_cost,
		s.negotiable, s.has_variants, s.middleman_service,
		s.attributes AS attributes_str, s.options AS options_str, s.thumbnail,
		s.lat, s.lng
	`

	selectQuery := fmt.Sprintf(`
		SELECT %s
		FROM %s s
		%s
		ORDER BY s.%s %s
		LIMIT $%d OFFSET $%d
	`, selectColumns, r.tableName, whereSQL, sortBy, sortOrder, argIndex, argIndex+1)

	rows, err := r.db.QueryContext(ctx, selectQuery, argsForSelect...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetServicesWithFilter: querying rows")
	}
	defer rows.Close()

	var results []*domain.CatalogService
	for rows.Next() {
		cs := domain.CatalogService{} // Distance field will remain 0.0
		var (
			pricingStr        sql.NullString
			qualificationsStr sql.NullString
			tagsStr           sql.NullString
			attributesStr     sql.NullString
			optionsStr        sql.NullString
			availabilityStr   sql.NullString
			providerNameStr   sql.NullString
			categorySlugStr   sql.NullString
			descShortStr      sql.NullString
			descLongStr       sql.NullString
			contactStr        sql.NullString
			faqStr            sql.NullString
			thumbnailStr      sql.NullString
			serviceTypeStr    sql.NullString // If service_type can be NULL
		)

		// Adjust scan destinations to match the `selectColumns` and `domain.CatalogService` struct
		err := rows.Scan(
			&cs.ID, &cs.Name, &cs.Description, &serviceTypeStr, &cs.BasePrice,
			&pricingStr, &availabilityStr, &providerNameStr, &cs.UserID,
			&cs.CategoryID, &categorySlugStr, &descShortStr,
			&descLongStr, &qualificationsStr, &contactStr, &faqStr,
			&tagsStr, &cs.Status, &cs.UserType, &cs.ShippingCost, // UserType and Status are enums, direct scan might work if driver handles it or they are string aliases.
			&cs.Negotiable, &cs.HasVariants, &cs.MiddlemanService,
			&attributesStr, &optionsStr, &thumbnailStr,
			&cs.Lat, &cs.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "GetServicesWithFilter: scanning service row")
		}

		// Populate fields from sql.NullString and deserialize
		if serviceTypeStr.Valid {
			cs.ServiceType = serviceTypeStr.String
		}
		if pricingStr.Valid {
			cs.Pricing = stringToSlice(pricingStr.String)
		}
		if availabilityStr.Valid {
			cs.Availability = availabilityStr.String
		} // Or parse if cs.Availability is time.Time
		if providerNameStr.Valid {
			cs.ProviderName = providerNameStr.String
		}
		if categorySlugStr.Valid {
			cs.CategorySlug = categorySlugStr.String
		}
		if descShortStr.Valid {
			cs.DescriptionShort = descShortStr.String
		}
		if descLongStr.Valid {
			cs.DescriptionLong = descLongStr.String
		}
		if qualificationsStr.Valid {
			cs.Qualifications = stringToSlice(qualificationsStr.String)
		}
		if contactStr.Valid {
			cs.Contact = contactStr.String
		}
		if faqStr.Valid {
			cs.Faq = faqStr.String
		}
		if tagsStr.Valid {
			cs.Tags = stringToSlice(tagsStr.String)
		}
		if attributesStr.Valid {
			cs.Attributes = stringToAttributes(attributesStr.String)
		}
		if optionsStr.Valid {
			cs.Options = stringToOptions(optionsStr.String)
		}
		if thumbnailStr.Valid {
			cs.Thumbnail = thumbnailStr.String
		}
		// cs.Distance will remain 0.0 as it's not selected

		results = append(results, &cs)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "GetServicesWithFilter: final iteration error")
	}

	return results, totalCount, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// Helper functions for converting between data types

func sliceToString(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	// Use JSON for more reliable serialization
	jsonBytes, err := json.Marshal(sl)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func stringToSlice(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		// Fallback to simple parsing if JSON unmarshal fails
		for _, item := range strings.Split(s, ",") {
			trimmed := strings.Trim(item, "\" []'")
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}

func attributesToString(attrs []domain.Attribute) string {
	if len(attrs) == 0 {
		return ""
	}

	// Use JSON for more reliable serialization
	jsonBytes, err := json.Marshal(attrs)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func stringToAttributes(s string) []domain.Attribute {
	if s == "" {
		return []domain.Attribute{}
	}

	var attrs []domain.Attribute
	err := json.Unmarshal([]byte(s), &attrs)
	if err != nil {
		return []domain.Attribute{}
	}
	return attrs
}

func optionsToString(options []domain.Option) string {
	if len(options) == 0 {
		return ""
	}

	// Use JSON for more reliable serialization
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func stringToOptions(s string) []domain.Option {
	if s == "" {
		return []domain.Option{}
	}

	var options []domain.Option
	err := json.Unmarshal([]byte(s), &options)
	if err != nil {
		return []domain.Option{}
	}
	return options
}
func joinClauses(clauses []string, separator string) string {
	return strings.Join(clauses, separator)
}
