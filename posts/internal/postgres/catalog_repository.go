package postgres

import (
	"context"
	"fmt"

	"middleman/internal/postgres"
	"middleman/posts/internal/domain"

	"github.com/stackus/errors"
)

// CatalogRepository implements domain.CatalogRepository
// for a PostgreSQL database.
type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

// Ensure CatalogRepository satisfies the interface
var _ domain.CatalogRepository = (*CatalogRepository)(nil)

// NewCatalogRepository returns a new instance
func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

// -----------------------------------------------------------------------------
// 1) AddPost
// -----------------------------------------------------------------------------
func (r CatalogRepository) AddPost(
	ctx context.Context,
	id, name, description string, typeOfPost domain.TypeOfPost, userID string, userType domain.UserType,
	category_id, category_slug string,
	tags []string,
	status domain.PostStatus,
	thumbnail string,
	lat, lng float64,
) error {

	tagsStr := sliceToString(tags)

	// Minimal columns: id, user_id, name, description, tags, status, thumbnail, location
	// Example schema must match your DB.
	// Corrected SQL query with the right number of placeholders:
	const query = `
      INSERT INTO %s (
        id, name, description, type_of_post, user_id, user_type, category_id, category_slug, tags, status, thumbnail, location, lat, lng
      )
      VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
        ST_SetSRID(ST_MakePoint($12, $13),4326), $12, $13
      )
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,            // $1
		name,          // $2
		description,   // $3
		typeOfPost,    // $4
		userID,        // $5
		userType,      // $6
		category_id,   // $7
		category_slug, // $8
		tagsStr,       // $9
		status,        // $10
		thumbnail,     // $11
		lng,           // $12 (for ST_MakePoint)
		lat,           // $13 (for ST_MakePoint)
	)
	return err // Return original error, potentially wrapped by caller if needed
}

// -----------------------------------------------------------------------------
// 2) UpdatePost
// -----------------------------------------------------------------------------
func (r CatalogRepository) UpdatePost(
	ctx context.Context,
	postID, name, description string, typeOfPost domain.TypeOfPost, userID string,
	tags []string,
	status domain.PostStatus,
	thumbnail string,
) error {

	tagsStr := sliceToString(tags)

	const query = `
      UPDATE %s
      SET
        name      = $2,
        description    = $3,
        type_of_post     = $4,
        tags       = $5,
        status      = $6,
        thumbnail  = $7,
        updated_at = NOW()
      WHERE id = $1
        AND user_id = $8
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		postID,
		name,
		description,
		typeOfPost,
		tagsStr,
		status.String(),
		thumbnail,
		userID,
	)
	return err
}

func (r CatalogRepository) UpdatePostThumbnail(
	ctx context.Context,
	postID,
	thumbnail string,
) error {

	const query = `
      UPDATE %s
      SET
        thumbnail  = $2,
        updated_at = NOW()
      WHERE id = $1 `
	_, err := r.db.ExecContext(ctx, r.table(query),
		postID,

		thumbnail,
	)
	return err
}

// -----------------------------------------------------------------------------
// 3) ArchivePost
// -----------------------------------------------------------------------------
func (r CatalogRepository) ArchivePost(ctx context.Context, postID, userID string) error {
	const query = `
      UPDATE %s
      SET
        status = 'ARCHIVED',
        updated_at = NOW()
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), postID, userID)
	return err
}

// -----------------------------------------------------------------------------
// 4) RemovePost
// -----------------------------------------------------------------------------
func (r CatalogRepository) RemovePost(ctx context.Context, postID, userID string) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
        AND user_id = $2
    `
	_, err := r.db.ExecContext(ctx, r.table(query), postID, userID)
	return err
}

// -----------------------------------------------------------------------------
// 5) Find
// -----------------------------------------------------------------------------
func (r CatalogRepository) Find(ctx context.Context, postID string) (*domain.CatalogPost, error) {
	const query = `
      SELECT   
        name,
        description,
        type_of_post,
        user_id,
        user_type,
        category_id,
        category_slug,
        tags,
        status,
        thumbnail,
        lng,
     	lat
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), postID)

	cp := &domain.CatalogPost{ID: postID}
	var tagsStr string

	err := row.Scan(
		&cp.Name,
		&cp.Description,
		&cp.TypeOfPost,
		&cp.UserID,
		&cp.UserType,
		&cp.CategoryID,
		&cp.CategorySlug,
		&tagsStr,
		&cp.Status,
		&cp.Thumbnail,
		&cp.Lng,
		&cp.Lat,
	)
	if err != nil {
		return nil, errors.Wrap(err, "scanning post")
	}
	cp.Tags = stringToSlice(tagsStr)
	return cp, nil
}

// -----------------------------------------------------------------------------
// 6) GetPosts
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetPosts(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogPost, int64, error) {

	offset := (page - 1) * pageSize

	validSortFields := map[string]bool{
		"updated_at": true,
		"created_at": true,
		"name":       true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts")
	}

	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        type_of_post,
        user_id,
        user_type,
        category_id,
        category_slug,
        tags,
        status,
        thumbnail,
      	lng,
       	lat
      FROM %s
      ORDER BY %s %s
      LIMIT $1
      OFFSET $2
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost

	for rows.Next() {
		cp := &domain.CatalogPost{}
		var tagsStr string

		if err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.UserID,
			&cp.UserType,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	return posts, totalCount, nil
}

func (r CatalogRepository) GetPostsWithFilters(
	ctx context.Context,
	name, description string,
	tags []string,
	offset, limit int64,
	lat, lng float64,
	radius, page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogPost, int64, error) {

	//-----------------------------------------------------------------------
	// 1) Build a dynamic WHERE clause + args for filtering
	//-----------------------------------------------------------------------
	whereClauses := "1=1"
	args := []interface{}{}

	// Filter by name (case-insensitive)
	if name != "" {
		whereClauses += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+name+"%")
	}
	// Filter by description
	if description != "" {
		whereClauses += fmt.Sprintf(" AND description ILIKE $%d", len(args)+1)
		args = append(args, "%"+description+"%")
	}
	// Filter by tags
	if len(tags) > 0 {
		// Example: "tags && $X" if 'tags' is a text[] column
		whereClauses += fmt.Sprintf(" AND tags && $%d", len(args)+1)
		args = append(args, sliceToString(tags)) // or handle text[] properly
	}
	// location + radius
	if radius > 0 {
		// Add ST_DWithin for bounding
		whereClauses += fmt.Sprintf(`
         AND ST_DWithin(
           location,
           ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography,
           $%d
         )
       `, len(args)+1, len(args)+2, len(args)+3)
		args = append(args, lng, lat, float64(radius))
	}

	//-----------------------------------------------------------------------
	// 2) Count the total rows that match (no LIMIT/OFFSET)
	//-----------------------------------------------------------------------
	countQuery := fmt.Sprintf(`
       SELECT COUNT(*)
       FROM %s
       WHERE %s
    `, r.tableName, whereClauses)

	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts with filters")
	}

	//-----------------------------------------------------------------------
	// 3) Build the main query: fetch the actual rows
	//-----------------------------------------------------------------------
	// We'll select ST_X(...), ST_Y(...), ST_Distance(...) AS dist
	// and order by dist (if you want to prioritize nearest first).
	// If you want a different sorting, adjust accordingly.
	mainQuery := fmt.Sprintf(`
       SELECT
          id,
          user_id,
          name,
          description,
          tags,
          status,
          thumbnail,
          lng,
          lat,
          ST_Distance(
            location,
            ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography
          ) AS dist
       FROM %s
       WHERE %s
       ORDER BY dist
       LIMIT $%d
       OFFSET $%d
    `,
		// The next 2 placeholders are for (lng, lat) in ST_Distance:
		len(args)+1, len(args)+2,
		r.tableName,
		whereClauses,
		// Then we append limit, offset at the end:
		len(args)+3, len(args)+4,
	)

	// We'll re-use 'args' but must append (lng, lat) again for the distance function,
	// plus limit + offset. So let's copy them carefully.
	mainArgs := make([]interface{}, len(args))
	copy(mainArgs, args)

	// We need to add the same (lng, lat) used in ST_DDistance.
	// If radius > 0, we already appended them once for ST_DWithin,
	// so let's do it again for ST_DDistance. If radius = 0, we do it anyway, for a consistent approach.
	mainArgs = append(mainArgs, lng, lat)

	// Finally add limit, offset
	mainArgs = append(mainArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, mainArgs...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts with filters")
	}
	defer rows.Close()

	//-----------------------------------------------------------------------
	// 4) Parse rows into CatalogPost
	//-----------------------------------------------------------------------
	var posts []*domain.CatalogPost
	for rows.Next() {
		cp := &domain.CatalogPost{}
		var tagsStr string
		var dist float64 // if you want to store or log distance

		if err := rows.Scan(
			&cp.ID,
			&cp.UserID,
			&cp.Name,
			&cp.Description,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
			&dist, // optional
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row with filters")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}

	return posts, totalCount, nil
}

func (r CatalogRepository) GetUserPosts(
	ctx context.Context,
	userID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogPost, int64, error) {

	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"updated_at": true,
		"created_at": true,
		"name":       true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts")
	}

	query := fmt.Sprintf(`
     SELECT
       id,
       name,
       description,
       type_of_post,
       user_id,
       user_type,
       category_id,
       category_slug,
       tags,
       status,
       thumbnail,
    	lng,
       lat
     FROM %s
     WHERE user_id = $1
     ORDER BY %s %s
     LIMIT $2
     OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost

	for rows.Next() {
		cp := &domain.CatalogPost{}
		var tagsStr string

		if err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.UserID,
			&cp.UserType,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	return posts, totalCount, nil
}

func (r CatalogRepository) GetPublicCatalog(
	ctx context.Context,
	userID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogPost, int64, error) {

	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"updated_at": true,
		"created_at": true,
		"name":       true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts")
	}

	query := fmt.Sprintf(`
     SELECT
       id,
       name,
       description,
       type_of_post,
       user_id,
       user_type,
       category_id,
       category_slug,
       tags,
       status,
       thumbnail,
      	lng,
       lat
     FROM %s
     WHERE user_id = $1
     ORDER BY %s %s
     LIMIT $2
     OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost

	for rows.Next() {
		cp := &domain.CatalogPost{}
		var tagsStr string

		if err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.UserID,
			&cp.UserType,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	return posts, totalCount, nil
}

func (r CatalogRepository) GetPostsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogPost, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"updated_at": true,
		"created_at": true,
		"name":       true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_id = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, categoryID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        type_of_post,
        user_id,
        user_type,
        category_id,
        category_slug,
        tags,
        status,
        thumbnail,
        lng,
        lat
      FROM %s
      WHERE category_id = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categoryID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts by category")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost
	for rows.Next() {
		cp := &domain.CatalogPost{
			CategoryID: categoryID,
		}
		var tagsStr string

		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.UserID,
			&cp.UserType,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row by category")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return posts, totalCount, nil
}

func (r CatalogRepository) GetPostsByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.CatalogPost, int64, error) {
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"updated_at": true,
		"created_at": true,
		"name":       true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// 1) Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_slug = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, categorySlug).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting posts by category")
	}

	// 2) Query
	query := fmt.Sprintf(`
      SELECT
        id,
        name,
        description,
        type_of_post,
        user_id,
        user_type,
        category_id,
        category_slug,
        tags,
        status,
        thumbnail,
        lng,
        lat
      FROM %s
      WHERE category_slug = $1
      ORDER BY %s %s
      LIMIT $2
      OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, categorySlug, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying posts by category")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost
	for rows.Next() {
		cp := &domain.CatalogPost{
			CategorySlug: categorySlug,
		}
		var tagsStr string

		err := rows.Scan(
			&cp.ID,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.UserID,
			&cp.UserType,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning post row by category")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading final rows by category")
	}
	return posts, totalCount, nil
}

// -----------------------------------------------------------------------------
// 9) FindByLocation
// -----------------------------------------------------------------------------
func (r CatalogRepository) FindByLocation(
	ctx context.Context,
	lat, lng float64,
	radiusMeters float64,
	limit int,
) ([]*domain.CatalogPost, error) {
	const queryTpl = `
      SELECT
        id,
        user_id,
        user_type,
        name,
        description,
        type_of_post,
        category_id,
        category_slug,
        tags,
        status,
        thumbnail,
         lng,      -- <--- cast to geometry
        lat,      -- <--- cast to geometry
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
    `
	query := r.table(queryTpl)
	rows, err := r.db.QueryContext(ctx, query, lng, lat, radiusMeters, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying posts by location")
	}
	defer rows.Close()

	var posts []*domain.CatalogPost
	for rows.Next() {
		cp := &domain.CatalogPost{}
		var (
			tagsStr string
			dist    float64 // optional distance
		)
		err := rows.Scan(
			&cp.ID,
			&cp.UserID,
			&cp.UserType,
			&cp.Name,
			&cp.Description,
			&cp.TypeOfPost,
			&cp.CategoryID,
			&cp.CategorySlug,
			&tagsStr,
			&cp.Status,
			&cp.Thumbnail,
			&cp.Lng,
			&cp.Lat,
			&dist, // ignoring or storing if needed
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning post row by location")
		}
		cp.Tags = stringToSlice(tagsStr)
		posts = append(posts, cp)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing location-based rows")
	}

	return posts, nil
}

// table is a helper to format table names in queries
func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// sliceToString converts []string into a stored string format
func sliceToString(sl []string) string {
	// For example, store as JSON or comma-separated
	return fmt.Sprintf("%q", sl)
}

// stringToSlice parses a stored string back to a []string
func stringToSlice(_ string) []string {
	// Implement a real parser if you store as JSON or CSV
	return []string{}
}
