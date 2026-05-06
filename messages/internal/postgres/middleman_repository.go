package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/messages/internal/domain"
)

type MiddlemanRepository struct {
	tableName string
	db        postgres.DB
}

// Compile-time check: ensures MiddlemanRepository implements domain.MiddlemanRepository
var _ domain.MiddlemanRepository = (*MiddlemanRepository)(nil)

func NewMiddlemanRepository(tableName string, db postgres.DB) MiddlemanRepository {
	return MiddlemanRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanRepository) Add(
	ctx context.Context,
	conversationID, senderID, recipientID, itemID string,
) error {
	const query = `
        INSERT INTO %s (id, sender_id, recipient_id, item_id, active)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(
		ctx,
		r.table(query),
		conversationID,
		senderID,
		recipientID,
		itemID,
		true, // new conversations start active
	)
	return err
}

func (r MiddlemanRepository) Find(
	ctx context.Context,
	conversationID string,
) (*domain.MiddlemanConversation, error) {
	const query = `
        SELECT id, sender_id, recipient_id, item_id, active
        FROM %s
        WHERE id = $1
        LIMIT 1
    `
	convo := &domain.MiddlemanConversation{}
	err := r.db.QueryRowContext(ctx, r.table(query), conversationID).
		Scan(&convo.ID, &convo.SenderID, &convo.RecipientID, &convo.ItemID, &convo.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msgf("conversation '%s' not found", conversationID)
		}
		return nil, errors.Wrap(err, "scanning conversation")
	}
	return convo, nil
}

func (r MiddlemanRepository) FindByRecipientAndItem(
	ctx context.Context,
	senderID, recipientID, itemID string,
) (*domain.MiddlemanConversation, error) {
	const query = `
        SELECT id, sender_id, recipient_id, item_id, active
        FROM %s
        WHERE sender_id = $1 AND recipient_id = $2
          AND item_id = $3 LIMIT 1`

	convo := &domain.MiddlemanConversation{}
	err := r.db.QueryRowContext(ctx, r.table(query), senderID, recipientID, itemID).
		Scan(&convo.ID, &convo.SenderID, &convo.RecipientID, &convo.ItemID, &convo.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msgf(
				"no conversation for sernder= %s, recipient=%s, item=%s",
				senderID, recipientID, itemID,
			)
		}
		return nil, errors.Wrap(err, "scanning conversation by recipient/item")
	}
	return convo, nil
}

func (r MiddlemanRepository) All(
	ctx context.Context,
	userID string,
) ([]*domain.MiddlemanConversation, error) {
	const query = `
        SELECT id, sender_id, recipient_id, item_id, active
        FROM %s
        WHERE (sender_id = $1) OR (recipient_id = $1)
    `
	rows, err := r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying conversations")
	}
	defer rows.Close()

	var conversations []*domain.MiddlemanConversation
	for rows.Next() {
		c := new(domain.MiddlemanConversation)
		if scanErr := rows.Scan(
			&c.ID,
			&c.SenderID,
			&c.RecipientID,
			&c.ItemID,
			&c.Active,
		); scanErr != nil {
			return nil, errors.Wrap(scanErr, "scanning conversation row")
		}
		conversations = append(conversations, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing conversation rows")
	}

	return conversations, nil
}

// helper to format the query
func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
