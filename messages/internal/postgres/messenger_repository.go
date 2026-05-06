package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/messages/internal/domain"
)

type MessengerRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MessengerRepository = (*MessengerRepository)(nil)

func NewMessengerRepository(tableName string, db postgres.DB) MessengerRepository {
	return MessengerRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MessengerRepository) SendMessage(ctx context.Context, messageID, conversationID, senderID, recipientID, itemID, body string, isRead bool) error {
	const query = `INSERT INTO %s (id,conversation_id, sender_id, recipient_id, item_id, body, is_read) VALUES ($1, $2, $3, $4, $5, $6,$7)`

	_, err := r.db.ExecContext(ctx, r.table(query), messageID, conversationID, senderID, recipientID, itemID, body, isRead)

	return err
}
func (r MessengerRepository) All(ctx context.Context, conversationID string) (messages []*domain.MiddlemanMessage, err error) {
	const query = `SELECT id, conversation_id, sender_id,recipient_id,item_id, body,is_read FROM %s WHERE conversation_id = $1`

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), conversationID)
	if err != nil {
		return nil, errors.Wrap(err, "querying messages")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing item rows")
		}
	}(rows)

	for rows.Next() {
		message := &domain.MiddlemanMessage{
			ConversationID: conversationID,
		}
		err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.RecipientID, &message.ItemID, &message.Body, &message.IsRead)
		if err != nil {
			return nil, errors.Wrap(err, "scanning message")
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing message rows")
	}

	return messages, nil
}

func (r MessengerRepository) Find(ctx context.Context, messageID string) (*domain.MiddlemanMessage, error) {
	const query = `SELECT conversation_id, sender_id, recipient_id, body,is_read FROM %s WHERE id = $1 LIMIT 1`

	item := &domain.MiddlemanMessage{
		ID: messageID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), messageID).Scan(&item.ConversationID, &item.SenderID, &item.RecipientID, &item.ItemID, &item.Body, &item.IsRead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("message with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning message")
	}

	return item, nil
}
func (r MessengerRepository) Delete(ctx context.Context, messageID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), messageID)

	return err
}
func (r MessengerRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
