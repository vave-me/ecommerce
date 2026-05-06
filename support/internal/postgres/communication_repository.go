package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	
	"github.com/lib/pq"
	"middleman/support/internal/domain"
)

type CommunicationRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.CommunicationRepository = (*CommunicationRepository)(nil)

func NewCommunicationRepository(tableName string, db *sql.DB) CommunicationRepository {
	return CommunicationRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CommunicationRepository) Add(ctx context.Context, comm *domain.Communication) error {
	const query = `
		INSERT INTO %s (id, ticket_id, author_id, author_type, content, is_public, mentioned_users, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	metadataJSON, err := json.Marshal(comm.Metadata)
	if err != nil {
		return err
	}
	
	_, err = r.db.ExecContext(ctx, r.table(query),
		comm.ID,
		comm.TicketID,
		comm.AuthorID,
		comm.AuthorType,
		comm.Content,
		comm.IsPublic,
		pq.Array(comm.MentionedUsers),
		metadataJSON,
		comm.CreatedAt,
	)
	
	if err != nil {
		return err
	}
	
	// Handle attachments
	if len(comm.Attachments) > 0 {
		err = r.addAttachments(ctx, comm.ID, "communication", comm.Attachments)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (r CommunicationRepository) GetByTicketID(ctx context.Context, ticketID string, includeInternal bool, limit, offset int) ([]*domain.Communication, error) {
	query := `
		SELECT c.id, c.ticket_id, c.author_id, c.author_type, c.content, c.is_public, c.mentioned_users, c.metadata, c.created_at
		FROM %s c
		WHERE c.ticket_id = $1
	`
	
	args := []interface{}{ticketID}
	argCount := 1
	
	if !includeInternal {
		query += fmt.Sprintf(" AND c.is_public = $%d", argCount+1)
		args = append(args, true)
		argCount++
	}
	
	query += " ORDER BY c.created_at ASC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit)
		argCount++
	}
	
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount+1)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var communications []*domain.Communication
	
	for rows.Next() {
		comm := &domain.Communication{}
		var metadataJSON []byte
		var mentionedUsers pq.StringArray
		
		err := rows.Scan(
			&comm.ID,
			&comm.TicketID,
			&comm.AuthorID,
			&comm.AuthorType,
			&comm.Content,
			&comm.IsPublic,
			&mentionedUsers,
			&metadataJSON,
			&comm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		comm.MentionedUsers = []string(mentionedUsers)
		
		if len(metadataJSON) > 0 {
			err = json.Unmarshal(metadataJSON, &comm.Metadata)
			if err != nil {
				return nil, err
			}
		}
		
		// Load attachments
		attachments, err := r.getAttachments(ctx, comm.ID)
		if err != nil {
			return nil, err
		}
		comm.Attachments = attachments
		
		communications = append(communications, comm)
	}
	
	return communications, rows.Err()
}

func (r CommunicationRepository) Count(ctx context.Context, ticketID string, includeInternal bool) (int, error) {
	query := `SELECT COUNT(*) FROM %s WHERE ticket_id = $1`
	args := []interface{}{ticketID}
	
	if !includeInternal {
		query += " AND is_public = $2"
		args = append(args, true)
	}
	
	var count int
	err := r.db.QueryRowContext(ctx, r.table(query), args...).Scan(&count)
	
	return count, err
}

func (r CommunicationRepository) addAttachments(ctx context.Context, entityID, entityType string, attachments []domain.Attachment) error {
	const query = `
		INSERT INTO attachments (id, entity_id, entity_type, filename, content_type, size_bytes, url, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	for _, attachment := range attachments {
		_, err := r.db.ExecContext(ctx, query,
			attachment.ID,
			entityID,
			entityType,
			attachment.Filename,
			attachment.ContentType,
			attachment.SizeBytes,
			attachment.URL,
			attachment.UploadedAt,
		)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (r CommunicationRepository) getAttachments(ctx context.Context, entityID string) ([]domain.Attachment, error) {
	const query = `
		SELECT id, filename, content_type, size_bytes, url, uploaded_at
		FROM attachments
		WHERE entity_id = $1 AND entity_type = 'communication'
		ORDER BY uploaded_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var attachments []domain.Attachment
	
	for rows.Next() {
		var attachment domain.Attachment
		err := rows.Scan(
			&attachment.ID,
			&attachment.Filename,
			&attachment.ContentType,
			&attachment.SizeBytes,
			&attachment.URL,
			&attachment.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	
	return attachments, rows.Err()
}

func (r CommunicationRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}