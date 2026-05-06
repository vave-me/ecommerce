package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"middleman/activity/internal/domain"
	"middleman/internal/postgres"
)

type MiddlemanInteractionRepository struct {
	tableName    string // e.g. "interactions"
	counterTable string // e.g. "item_interactions_counts"
	db           postgres.DB
}

// Compile-time check: ensure MiddlemanInteractionRepository implements domain.MiddlemanInteractionRepository
var _ domain.MiddlemanInteractionRepository = (*MiddlemanInteractionRepository)(nil)

// NewMiddlemanInteractionRepository constructs a new repository for the given tableName (e.g. interactions)
// and counterTable (e.g. item_interactions_counts), plus a DB connection.
func NewMiddlemanInteractionRepository(tableName, counterTable string, db postgres.DB) MiddlemanInteractionRepository {
	return MiddlemanInteractionRepository{
		tableName:    tableName,
		counterTable: counterTable,
		db:           db,
	}
}

// Add inserts a new interaction record into the 'interactions' table.
func (r MiddlemanInteractionRepository) Add(
	ctx context.Context,
	interactionID, activityID, itemID, itemType, actionType string,
) error {
	const query = `INSERT INTO %s (id, activity_id, item_id, item_type, action_type)
                   VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, r.table(query), interactionID, activityID, itemID, itemType, actionType)
	if err != nil {
		return errors.Wrap(err, "inserting new interaction")
	}
	return nil
}

// Update changes the actionType of an existing interaction in the 'interactions' table.
func (r MiddlemanInteractionRepository) Update(
	ctx context.Context,
	interactionID, actionType string,
) error {
	const query = `UPDATE %s
                      SET action_type = $2
                    WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), interactionID, actionType)
	if err != nil {
		return errors.Wrap(err, "updating interaction actionType")
	}
	return nil
}

// Remove deletes an interaction record from the 'interactions' table by its ID.
func (r MiddlemanInteractionRepository) Remove(
	ctx context.Context,
	interactionID string,
) error {
	const query = `DELETE FROM %s
                    WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), interactionID)
	if err != nil {
		return errors.Wrap(err, "deleting interaction")
	}
	return nil
}

// Find retrieves a single interaction by its ID from the 'interactions' table.
func (r MiddlemanInteractionRepository) Find(
	ctx context.Context,
	interactionID string,
) (*domain.MiddlemanInteraction, error) {
	const query = `SELECT activity_id, item_id, item_type, action_type
                     FROM %s
                    WHERE id = $1
                    LIMIT 1`

	interaction := &domain.MiddlemanInteraction{ID: interactionID}

	err := r.db.QueryRowContext(ctx, r.table(query), interactionID).
		Scan(&interaction.ActivityID, &interaction.ItemID, &interaction.ItemType, &interaction.ActionType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("interaction with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning interaction")
	}
	return interaction, nil
}

// All returns all interactions for a given activityID from the 'interactions' table.
func (r MiddlemanInteractionRepository) All(
	ctx context.Context,
	activityID string,
) ([]*domain.MiddlemanInteraction, error) {
	const query = `SELECT id, activity_id, item_id, item_type, action_type
                     FROM %s
                    WHERE activity_id = $1`

	rows, err := r.db.QueryContext(ctx, r.table(query), activityID)
	if err != nil {
		return nil, errors.Wrap(err, "querying interactions")
	}
	defer rows.Close()

	var interactions []*domain.MiddlemanInteraction
	for rows.Next() {
		interaction := &domain.MiddlemanInteraction{}
		if err := rows.Scan(&interaction.ID, &interaction.ActivityID, &interaction.ItemID, &interaction.ItemType, &interaction.ActionType); err != nil {
			return nil, errors.Wrap(err, "scanning interaction row")
		}
		interactions = append(interactions, interaction)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing interactions rows")
	}
	return interactions, nil
}

// --------------------------------------------------------------------
// The next methods use the aggregator table (counterTable), e.g. "item_interactions_counts",
// which might have columns like:
//    item_id text
//    item_type text
//    total_count bigint
//    like_count bigint
//    dislike_count bigint
// We’ll query "like_count" or "dislike_count" and sort descending for the top results.
// --------------------------------------------------------------------

// GetMostLiked retrieves the top 'limit' items by descending like_count from the aggregator table.
func (r MiddlemanInteractionRepository) GetMostLiked(
	ctx context.Context,
	itemType string,
	limit int64,
) ([]*domain.MostReactionResult, error) {
	// Example aggregator query (assuming aggregator columns: like_count, item_id, item_type)
	// SELECT item_id, item_type, 'like' AS action, like_count AS cnt
	//   FROM <counter_table>
	//  WHERE item_type = $1
	//  ORDER BY like_count DESC
	//  LIMIT $2

	const query = `
        SELECT item_id, item_type, 'like' AS action, like_count AS cnt
          FROM %s
         WHERE item_type = $1
         ORDER BY like_count DESC
         LIMIT $2
    `
	rows, err := r.db.QueryContext(ctx, r.countTable(query), itemType, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying most liked items from aggregator")
	}
	defer rows.Close()

	var results []*domain.MostReactionResult
	for rows.Next() {
		var (
			itemID, itemTypeCol, action string
			count                       int64
		)
		if scanErr := rows.Scan(&itemID, &itemTypeCol, &action, &count); scanErr != nil {
			return nil, errors.Wrap(scanErr, "scanning most liked result")
		}
		results = append(results, &domain.MostReactionResult{
			ItemID:   itemID,
			ItemType: itemTypeCol,
			Action:   action,
			Count:    count,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing most liked rows")
	}
	return results, nil
}

// GetMostDisliked retrieves the top 'limit' items by descending dislike_count from the aggregator table.
func (r MiddlemanInteractionRepository) GetMostDisliked(
	ctx context.Context,
	itemType string,
	limit int64,
) ([]*domain.MostReactionResult, error) {
	// Example aggregator query (assuming aggregator columns: dislike_count, item_id, item_type)
	// SELECT item_id, item_type, 'dislike' AS action, dislike_count AS cnt
	//   FROM <counter_table>
	//  WHERE item_type = $1
	//  ORDER BY dislike_count DESC
	//  LIMIT $2
	const query = `
        SELECT item_id, item_type, 'dislike' AS action, dislike_count AS cnt
          FROM %s
         WHERE item_type = $1
         ORDER BY dislike_count DESC
         LIMIT $2
    `
	rows, err := r.db.QueryContext(ctx, r.countTable(query), itemType, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying most disliked items from aggregator")
	}
	defer rows.Close()

	var results []*domain.MostReactionResult
	for rows.Next() {
		var (
			itemID, itemTypeCol, action string
			count                       int64
		)
		if scanErr := rows.Scan(&itemID, &itemTypeCol, &action, &count); scanErr != nil {
			return nil, errors.Wrap(scanErr, "scanning most disliked result")
		}
		results = append(results, &domain.MostReactionResult{
			ItemID:   itemID,
			ItemType: itemTypeCol,
			Action:   action,
			Count:    count,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing most disliked rows")
	}
	return results, nil
}

// table formats the query string to include the actual "interactions" table name.
func (r MiddlemanInteractionRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// countTable formats the query string to include the aggregator table name,
// e.g. "item_interactions_counts".
func (r MiddlemanInteractionRepository) countTable(query string) string {
	return fmt.Sprintf(query, r.counterTable)
}
