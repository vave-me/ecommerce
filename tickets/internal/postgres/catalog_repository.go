package postgres

import (
	"context"
	"fmt"
	"strings"

	"middleman/internal/postgres"
	"middleman/tickets/internal/domain"

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

func (r CatalogRepository) AddTicket(
	ctx context.Context,
	id, name, description string,
	basePrice int64,
	userSellerID, categoryID, categorySlug, brand string,
	condition domain.TicketCondition,
	model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []domain.Attribute,
	weight, height, width, depth int64,
	status domain.TicketStatus,
	negotiable bool,
	userType domain.UserType,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []domain.Option,
	thumbnail string,
	lat, lng float64, // no 'location' param
) error {

	tagsStr := sliceToString(tags)
	attrsStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	// We'll insert into the 'location' column using ST_SetSRID(ST_MakePoint(:lng, :lat),4326)
	// So in the VALUES list, we do: ST_SetSRID(ST_MakePoint($?, $?), 4326) for lat/lng
	const query = `
      INSERT INTO %s (
        id, name, description,
        base_price, user_seller_id, category_id,category_slug, brand, condition, model,
        tags, attributes,
        manage_stock, stock, sku,
        weight, height, width, depth,
        status, negotiable, user_type, middleman_service,
        shipping_cost, has_variants, options, thumbnail,
        location, lat, lng
      )
      VALUES (
        $1,$2,$3,
        $4,$5,$6,$7,$8,$9,
        $10,$11,
        $12,$13,$14,
        $15,$16,$17,$18,
        $19,$20,$21,$22,
        $23,$24,$25,$26,$27,
        ST_SetSRID(ST_MakePoint($28, $29),4326), $28, $29
      )
    `

	_, err := r.db.ExecContext(ctx, r.table(query),
		id,               //1
		name,             //2
		description,      //3
		basePrice,        //4
		userSellerID,     //5
		categoryID,       //6
		categorySlug,     //7
		brand,            //8
		condition,        //9
		model,            //10
		tagsStr,          //11
		attrsStr,         //12
		manageStock,      //13
		stock,            //14
		sku,              //15
		weight,           //16
		height,           //17
		width,            //18
		depth,            //19
		status,           //20
		negotiable,       //21
		userType,         //22
		middlemanService, //23
		shippingCost,     //24
		hasVariants,      //25
		optsStr,          //26
		thumbnail,        //27
		lng,              // 28pass Lng first (ST_MakePoint expects (lng, lat))
		lat,              //29
	)
	return err
}

func (r CatalogRepository) RebrandTicket(
	ctx context.Context,
	ticketID string,
	newName string,
	newDescription string,
	categoryID string,
	categorySlug string,
	newBrand string,
	newModel string,
	newCondition string,
	newTags []string,
) error {
	tagsStr := sliceToString(newTags)

	const query = `
      UPDATE %s
      SET name = $2,
          description = $3,
          brand = $4,
          category_id = $5,
          category_slug = $6,
          model = $7,
          condition = $8,
          tags = $9,
          updated_at = NOW()
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		ticketID,
		newName,
		newDescription,
		categoryID,
		categorySlug,
		newBrand,
		newModel,
		newCondition,
		tagsStr,
	)
	return err
}

// -----------------------------------------------------------------------------
// 3) UpdateTicket (now storing tags + attributes)
// -----------------------------------------------------------------------------
func (r CatalogRepository) UpdateTicket(
	ctx context.Context, ticketID, name, description string,
	basePrice int64,
	categoryID, categorySlug, brand string,
	condition domain.TicketCondition,
	model string, tags []string, manageStock bool, stock int64, sku string, attributes []domain.Attribute,
	weight, height, width, depth int64, status domain.TicketStatus,
	negotiable bool, userType domain.UserType, middlemanService bool, shippingCost int64, hasVariants bool, options []domain.Option, thumbnail string, lat float64, lng float64) error {
	tagsStr := sliceToString(tags)
	attrsStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	// Update statement now sets tags, attributes as well.
	const query = `
      UPDATE %s
      SET
        name              = $2,
        description       = $3,
        base_price 		  = $4,
        category_id       = $5,
        category_slug     = $6,
        brand             = $7,
        condition         = $8,
        model             = $9,
        tags              = $10,
        attributes        = $11,
        manage_stock      = $12,
        stock             = $13,
        sku               = $14,
        weight            = $15,
        height            = $16,
        width             = $17,
        depth             = $18,
        status            = $19,
        negotiable        = $20,
        user_type       = $21,
        middleman_service = $22,
        shipping_cost     = $23,
        has_variants      = $24,
        options           = $25,
        thumbnail         = $26,
        lat               = $27,
        lng               = $28,
      	location = ST_SetSRID(ST_MakePoint($27, $28),4326),
       updated_at        = NOW()
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		ticketID,
		name,
		description,
		basePrice,
		categoryID,
		categorySlug,
		brand,
		condition,
		model,
		tagsStr,
		attrsStr,
		manageStock,
		stock,
		sku,
		weight,
		height,
		width,
		depth,
		status,
		negotiable,
		userType,
		middlemanService,
		shippingCost,
		hasVariants,
		optsStr,
		thumbnail,
		lat,
		lng,
	)
	return err
}

// -----------------------------------------------------------------------------
// 4) UpdatePrice
// -----------------------------------------------------------------------------
func (r CatalogRepository) UpdatePrice(ctx context.Context, ticketID string, oldPrice, newPrice int64) error {
	const query = `
      UPDATE %s
      SET base_price = $3,
          updated_at = NOW()
      WHERE id = $1
        AND base_price = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), ticketID, oldPrice, newPrice)
	if err != nil {
		return errors.Wrap(err, "updating price")
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "no ticket updated; either ticket not found or base_price mismatch")
	}
	return nil
}

// -----------------------------------------------------------------------------
func (r CatalogRepository) UpdateThumbnail(ctx context.Context, ticketID, thumbnail string) error {
	const query = `
      UPDATE %s
      SET thumbnail = $2,
          updated_at = NOW()
      WHERE id = $1
    `
	res, err := r.db.ExecContext(ctx, r.table(query), ticketID, thumbnail)
	if err != nil {
		return errors.Wrap(err, "updating thumbnail")
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "no ticket updated; either ticket not found or thumbnail mismatch")
	}
	return nil
}

func (r CatalogRepository) GetTicketsWithFilters(
	ctx context.Context,
	name string,
	categoryID string,
	categorySlug string,
	minPrice int64,
	maxPrice int64,
	brand string,
	condition string,
	model string,
	tags []string,
	manageStock bool,
	minStock int64,
	maxStock int64,
	sku string,
	status string,
	negotiable bool,
	userType string,
	middlemanService bool,
	hasVariants bool,
	shippingCost int64,
	minWeight int64,
	maxWeight int64,
	minHeight int64,
	maxHeight int64,
	minWidth int64,
	maxWidth int64,
	minDepth int64,
	maxDepth int64,
	offset int64,
	limit int64,
	lat, lng float64,
	radius int64,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogTicket, int64, error) {

	// 1) Build dynamic WHERE clauses + arguments
	whereClauses := []string{"1=1"}
	args := make([]interface{}, 0, 32)
	argIndex := 1

	// Helper to add a clause: also increments argIndex
	addClause := func(sqlClause string, arg interface{}) {
		whereClauses = append(whereClauses, sqlClause)
		args = append(args, arg)
		argIndex++
	}

	// name partial
	if name != "" {
		addClause(fmt.Sprintf("name ILIKE $%d", argIndex), "%"+name+"%")
	}

	// category exact
	if categoryID != "" {
		addClause(fmt.Sprintf("category_id = $%d", argIndex), categoryID)
	}
	if categorySlug != "" {
		addClause(fmt.Sprintf("category_slug = $%d", argIndex), categorySlug)
	}
	// minPrice / maxPrice
	if minPrice > 0 {
		addClause(fmt.Sprintf("base_price >= $%d", argIndex), minPrice)
	}
	if maxPrice > 0 {
		addClause(fmt.Sprintf("base_price <= $%d", argIndex), maxPrice)
	}

	// brand exact
	if brand != "" {
		addClause(fmt.Sprintf("brand = $%d", argIndex), brand)
	}

	// condition exact
	if condition != "" {
		addClause(fmt.Sprintf("condition = $%d", argIndex), condition)
	}

	// model partial
	if model != "" {
		addClause(fmt.Sprintf("model ILIKE $%d", argIndex), "%"+model+"%")
	}

	// tags naive approach => "tags ILIKE '%someTag%'"
	for _, t := range tags {
		addClause(fmt.Sprintf("tags ILIKE $%d", argIndex), "%"+t+"%")
	}

	// manageStock => boolean
	if manageStock {
		whereClauses = append(whereClauses, "manage_stock = TRUE")
	}

	// minStock / maxStock
	if minStock > 0 {
		addClause(fmt.Sprintf("stock >= $%d", argIndex), minStock)
	}
	if maxStock > 0 {
		addClause(fmt.Sprintf("stock <= $%d", argIndex), maxStock)
	}

	// sku partial
	if sku != "" {
		addClause(fmt.Sprintf("sku ILIKE $%d", argIndex), "%"+sku+"%")
	}

	// status exact
	if status != "" {
		addClause(fmt.Sprintf("status = $%d", argIndex), status)
	}

	// negotiable => bool
	if negotiable {
		whereClauses = append(whereClauses, "negotiable = TRUE")
	}

	// userType => exact
	if userType != "" {
		addClause(fmt.Sprintf("user_type = $%d", argIndex), userType)
	}

	// middleman => bool
	if middlemanService {
		whereClauses = append(whereClauses, "middleman_service = TRUE")
	}

	// hasVariants => bool
	if hasVariants {
		whereClauses = append(whereClauses, "has_variants = TRUE")
	}

	// shippingCost => exact
	if shippingCost > 0 {
		addClause(fmt.Sprintf("shipping_cost = $%d", argIndex), shippingCost)
	}

	// weight range
	if minWeight > 0 {
		addClause(fmt.Sprintf("weight >= $%d", argIndex), minWeight)
	}
	if maxWeight > 0 {
		addClause(fmt.Sprintf("weight <= $%d", argIndex), maxWeight)
	}

	// height range
	if minHeight > 0 {
		addClause(fmt.Sprintf("height >= $%d", argIndex), minHeight)
	}
	if maxHeight > 0 {
		addClause(fmt.Sprintf("height <= $%d", argIndex), maxHeight)
	}

	// width range
	if minWidth > 0 {
		addClause(fmt.Sprintf("width >= $%d", argIndex), minWidth)
	}
	if maxWidth > 0 {
		addClause(fmt.Sprintf("width <= $%d", argIndex), maxWidth)
	}

	// depth range
	if minDepth > 0 {
		addClause(fmt.Sprintf("depth >= $%d", argIndex), minDepth)
	}
	if maxDepth > 0 {
		addClause(fmt.Sprintf("depth <= $%d", argIndex), maxDepth)
	}

	// Optional geo check if radius>0 & lat/lng not zero
	if radius > 0 && lat != 0 && lng != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`
			ST_DWithin(
				location,
				ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography,
				$%d
			)
		`, argIndex, argIndex+1, argIndex+2))
		args = append(args, lng, lat, float64(radius))
		argIndex += 3
	}

	// 2) Figure out final offset/limit from page or from offset/limit
	var finalOffset, finalLimit int64
	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		finalOffset = (page - 1) * pageSize
		finalLimit = pageSize
	} else {
		finalOffset = offset
		finalLimit = limit
		if finalLimit <= 0 {
			finalLimit = 50
		}
	}

	// 3) Validate sort fields
	validSort := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
		"updated_at": true,
	}
	if !validSort[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 4) Build the final WHERE string
	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + joinClauses(whereClauses, " AND ")
	}

	// 5) Query for totalCount (if needed for pagination UI)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", r.tableName, whereSQL)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "GetTicketsWithFilters: counting total items")
	}

	// 6) Build final SELECT for ticket rows
	selectQuery := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			base_price,
			user_seller_id,
			category_id,
			category_slug,
			brand,
			condition,
			model,
			tags,
			attributes,
			manage_stock,
			stock,
			sku,
			weight,
			height,
			width,
			depth,
			status,
			negotiable,
			user_type,
			middleman_service,
			shipping_cost,
			has_variants,
			options,
			thumbnail,
			lat,
			lng
		FROM %s
		%s
		ORDER BY %s %s
		LIMIT $%d
		OFFSET $%d
	`, r.tableName, whereSQL, sortBy, sortOrder, argIndex, argIndex+1)

	// We append finalLimit, finalOffset after the existing args
	args = append(args, finalLimit, finalOffset)

	// 7) Execute the SELECT
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetTicketsWithFilters: query rows error")
	}
	defer rows.Close()

	var results []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		scErr := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if scErr != nil {
			return nil, 0, errors.Wrap(scErr, "GetTicketsWithFilters scanning ticket row")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		results = append(results, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "GetTicketsWithFilters final iteration err")
	}

	return results, totalCount, nil
}

// -----------------------------------------------------------------------------
// 5) AdjustStock
// -----------------------------------------------------------------------------
func (r CatalogRepository) AdjustStock(
	ctx context.Context,
	ticketID string,
	userSellerID string,
	oldStock, newStock int64,
) error {
	const query = `
      UPDATE %s
      SET stock = $3,
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
        AND stock = $4
    `
	res, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID, newStock, oldStock)
	if err != nil {
		return errors.Wrap(err, "adjusting ticket stock")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.Wrap(err, "no ticket updated; stock mismatch or not found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// 6) ToggleNegotiable
// -----------------------------------------------------------------------------
func (r CatalogRepository) ToggleNegotiable(
	ctx context.Context,
	ticketID string,
	userSellerID string,
	currentValue bool,
) error {
	const query = `
      UPDATE %s
      SET negotiable = NOT negotiable,
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
        AND negotiable = $3
    `
	res, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID, currentValue)
	if err != nil {
		return errors.Wrap(err, "toggle negotiable error")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.Wrap(err, "no ticket updated; negotiable mismatch or not found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// 7) RemoveTicket
// -----------------------------------------------------------------------------
func (r CatalogRepository) RemoveTicket(
	ctx context.Context,
	ticketID string,
	userSellerID string,
) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
        AND user_seller_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID)
	return err
}

// -----------------------------------------------------------------------------
// 8) MarkTicketPawned
// -----------------------------------------------------------------------------
func (r CatalogRepository) MarkTicketPawned(
	ctx context.Context,
	ticketID string,
	userSellerID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'pawned',
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID)
	return err
}

// -----------------------------------------------------------------------------
// 9) MarkTicketLeased
// -----------------------------------------------------------------------------
func (r CatalogRepository) MarkTicketLeased(
	ctx context.Context,
	ticketID string,
	userSellerID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'leased',
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID)
	return err
}

// -----------------------------------------------------------------------------
// 10) MarkTicketSold
// -----------------------------------------------------------------------------
func (r CatalogRepository) MarkTicketSold(
	ctx context.Context,
	ticketID string,
	userSellerID string,
	finalPrice int64,
) error {
	const query = `
      UPDATE %s
      SET status = 'sold',
          base_price = $3,
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID, finalPrice)
	return err
}

// -----------------------------------------------------------------------------
// 11) ArchiveTicket
// -----------------------------------------------------------------------------
func (r CatalogRepository) ArchiveTicket(
	ctx context.Context,
	ticketID string,
	userSellerID string,
) error {
	const query = `
      UPDATE %s
      SET status = 'archived',
          updated_at = NOW()
      WHERE id = $1
        AND user_seller_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID, userSellerID)
	return err
}

// -----------------------------------------------------------------------------
// 12) Find
// -----------------------------------------------------------------------------
func (r CatalogRepository) Find(ctx context.Context, ticketID string) (*domain.CatalogTicket, error) {
	const query = `
      SELECT
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        category_slug,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
		lat,
		lng
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), ticketID)
	cp := &domain.CatalogTicket{ID: ticketID}

	var (
		tagsStr  string
		attrsStr string
		optsStr  string
		mmsBool  bool
	)

	err := row.Scan(
		&cp.Name,
		&cp.Description,
		&cp.BasePrice,
		&cp.UserSellerID,
		&cp.CategoryID,
		&cp.CategorySlug,
		&cp.Brand,
		&cp.Condition,
		&cp.Model,
		&tagsStr,
		&attrsStr,
		&cp.ManageStock,
		&cp.Stock,
		&cp.SKU,
		&cp.Weight,
		&cp.Height,
		&cp.Width,
		&cp.Depth,
		&cp.Status,
		&cp.Negotiable,
		&cp.UserType,
		&mmsBool,
		&cp.ShippingCost,
		&cp.HasVariants,
		&optsStr,
		&cp.Lat,
		&cp.Lng,
	)
	if err != nil {
		return nil, errors.Wrap(err, "scanning ticket")
	}

	cp.Tags = stringToSlice(tagsStr)
	cp.Attributes = stringToAttributes(attrsStr)
	cp.MiddlemanService = mmsBool
	cp.Options = stringToOptions(optsStr)

	return cp, nil
}

// -----------------------------------------------------------------------------
// 13) GetCatalog
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogTicket, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_seller_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userSellerID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting tickets")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
		lat,
		lng
      FROM %s
      WHERE user_seller_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userSellerID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying catalog")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning ticket row")
		}
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final ticket rows")
	}

	return tickets, totalCount, nil
}

// -----------------------------------------------------------------------------
// 13) GetCatalog
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetPublicCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogTicket, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_seller_id = $1 AND status = 'available'", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userSellerID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting tickets")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
		lat
      	lng
      FROM %s
      WHERE user_seller_id = $1 AND status = 'available'
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userSellerID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying catalog")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning ticket row")
		}
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final ticket rows")
	}

	return tickets, totalCount, nil
}

func (r CatalogRepository) GetTickets(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogTicket, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
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
		return nil, 0, errors.Wrap(err, "counting tickets")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
		lat,
		lng
      FROM %s
      ORDER BY %s %s
      LIMIT $1
      OFFSET $2
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying all tickets")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning ticket row")
		}
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing ticket rows")
	}

	return tickets, totalCount, nil
}

// -----------------------------------------------------------------------------
// 15) GetTicketsByCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetTicketsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogTicket, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
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
		return nil, 0, errors.Wrap(err, "counting tickets by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_slug,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
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
		return nil, 0, errors.Wrap(err, "querying tickets by category")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{
			CategoryID: categoryID,
		}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning ticket row by category")
		}
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return tickets, totalCount, nil
}

func (r CatalogRepository) GetTicketsByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogTicket, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"name":       true,
		"base_price": true,
		"stock":      true,
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
		return nil, 0, errors.Wrap(err, "counting tickets by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
		lat, 
		lng
      FROM %s
      WHERE category_slug= $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categorySlug, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying tickets by category")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{
			CategorySlug: categorySlug,
		}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning ticket row by category")
		}
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return tickets, totalCount, nil
}

func (r *CatalogRepository) FindByLocation(
	ctx context.Context,
	lat, lng float64,
	radiusMeters float64,
	limit int,
) ([]*domain.CatalogTicket, error) {

	// We SELECT all the columns you also retrieve in GetTickets or Find,
	// plus ST_Distance(...) as dist so you can sort or store it if desired.
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        category_slug,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
        thumbnail,
        ST_Distance(
          location,
          ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
        ) AS dist,
		lat, 
		lng
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
		return nil, errors.Wrap(err, "querying tickets by location")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
			dist     float64 // distance in meters from (lat,lng)
		)
		// We scan all columns in the same order we selected
		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Thumbnail,
			&dist,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning ticket row by location")
		}

		// Convert strings to slices/structs
		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		// Optionally store dist in the domain struct (if you add a Dist field)
		// cp.Dist = dist

		tickets = append(tickets, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing ticket rows by location")
	}

	return tickets, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// For the sake of example, these are minimal stubs.
// In real usage, you might store JSON or parse strings more robustly.

func sliceToString(sl []string) string {
	// e.g. convert ["red","blue"] => "red,blue"
	return fmt.Sprintf("%q", sl)
}

func stringToSlice(s string) []string {
	// e.g. parse "['red','blue']" back to []string
	return []string{} // implement as needed
}

func attributesToString(attrs []domain.Attribute) string {
	// Possibly store as JSON
	return fmt.Sprintf("%q", attrs)
}

func stringToAttributes(s string) []domain.Attribute {
	// parse the string => []Attribute
	return []domain.Attribute{}
}

func optionsToString(opts []domain.Option) string {
	// Possibly store as JSON
	return fmt.Sprintf("%q", opts)
}

func stringToOptions(s string) []domain.Option {
	// parse the string => []Option
	return []domain.Option{}
}

func joinClauses(clauses []string, op string) string {
	if len(clauses) == 0 {
		return ""
	}
	res := clauses[0]
	for i := 1; i < len(clauses); i++ {
		res += " " + op + " " + clauses[i]
	}
	return res
}

// GetTicketBySKU retrieves a single ticket by its SKU
func (r CatalogRepository) GetTicketBySKU(ctx context.Context, sku string) (*domain.CatalogTicket, error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			base_price,
			user_seller_id,
			category_id,
			category_slug,
			brand,
			condition,
			model,
			tags,
			attributes,
			manage_stock,
			stock,
			sku,
			weight,
			height,
			width,
			depth,
			status,
			negotiable,
			user_type,
			middleman_service,
			shipping_cost,
			has_variants,
			options,
			thumbnail,
			lat,
			lng
		FROM %s
		WHERE sku = $1
		LIMIT 1
	`, r.tableName)

	row := r.db.QueryRowContext(ctx, query, sku)

	cp := &domain.CatalogTicket{}
	var (
		tagsStr  string
		attrsStr string
		optsStr  string
		mmsBool  bool
	)

	err := row.Scan(
		&cp.ID,
		&cp.Name,
		&cp.Description,
		&cp.BasePrice,
		&cp.UserSellerID,
		&cp.CategoryID,
		&cp.CategorySlug,
		&cp.Brand,
		&cp.Condition,
		&cp.Model,
		&tagsStr,
		&attrsStr,
		&cp.ManageStock,
		&cp.Stock,
		&cp.SKU,
		&cp.Weight,
		&cp.Height,
		&cp.Width,
		&cp.Depth,
		&cp.Status,
		&cp.Negotiable,
		&cp.UserType,
		&mmsBool,
		&cp.ShippingCost,
		&cp.HasVariants,
		&optsStr,
		&cp.Thumbnail,
		&cp.Lat,
		&cp.Lng,
	)
	if err != nil {
		return nil, errors.Wrap(err, "scanning ticket by SKU")
	}

	cp.Tags = stringToSlice(tagsStr)
	cp.Attributes = stringToAttributes(attrsStr)
	cp.MiddlemanService = mmsBool
	cp.Options = stringToOptions(optsStr)

	return cp, nil
}

// GetTicketsBySKUs retrieves multiple tickets by their SKUs
func (r CatalogRepository) GetTicketsBySKUs(ctx context.Context, skus []string) ([]*domain.CatalogTicket, error) {
	if len(skus) == 0 {
		return []*domain.CatalogTicket{}, nil
	}

	// Build placeholders for IN clause
	placeholders := make([]string, len(skus))
	args := make([]interface{}, len(skus))
	for i, sku := range skus {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = sku
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			base_price,
			user_seller_id,
			category_id,
			category_slug,
			brand,
			condition,
			model,
			tags,
			attributes,
			manage_stock,
			stock,
			sku,
			weight,
			height,
			width,
			depth,
			status,
			negotiable,
			user_type,
			middleman_service,
			shipping_cost,
			has_variants,
			options,
			thumbnail,
			lat,
			lng
		FROM %s
		WHERE sku IN (%s)
	`, r.tableName, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying tickets by SKUs")
	}
	defer rows.Close()

	var tickets []*domain.CatalogTicket
	for rows.Next() {
		cp := &domain.CatalogTicket{}
		var (
			tagsStr  string
			attrsStr string
			optsStr  string
			mmsBool  bool
		)

		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.BasePrice,
			&cp.UserSellerID,
			&cp.CategoryID,
			&cp.CategorySlug,
			&cp.Brand,
			&cp.Condition,
			&cp.Model,
			&tagsStr,
			&attrsStr,
			&cp.ManageStock,
			&cp.Stock,
			&cp.SKU,
			&cp.Weight,
			&cp.Height,
			&cp.Width,
			&cp.Depth,
			&cp.Status,
			&cp.Negotiable,
			&cp.UserType,
			&mmsBool,
			&cp.ShippingCost,
			&cp.HasVariants,
			&optsStr,
			&cp.Thumbnail,
			&cp.Lat,
			&cp.Lng,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning ticket row")
		}

		cp.Tags = stringToSlice(tagsStr)
		cp.Attributes = stringToAttributes(attrsStr)
		cp.MiddlemanService = mmsBool
		cp.Options = stringToOptions(optsStr)

		tickets = append(tickets, cp)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing ticket rows")
	}

	return tickets, nil
}
