package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"support-go/backend/internal/ticket"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (repository *CommentRepository) Create(comment ticket.Comment) error {
	const query = `
		INSERT INTO ticket_comments (
			id, ticket_id, author_id, body, is_internal, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := repository.pool.Exec(context.Background(), query,
		comment.ID,
		comment.TicketID,
		comment.AuthorID,
		comment.Body,
		comment.IsInternal,
		comment.CreatedAt,
	)

	return err
}

func (repository *CommentRepository) ListByTicketID(ticketID string) ([]ticket.Comment, error) {
	const query = `
		SELECT id, ticket_id, author_id, body, is_internal, created_at
		FROM ticket_comments
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := repository.pool.Query(context.Background(), query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ticket.Comment, 0)
	for rows.Next() {
		var item ticket.Comment
		if scanErr := rows.Scan(
			&item.ID,
			&item.TicketID,
			&item.AuthorID,
			&item.Body,
			&item.IsInternal,
			&item.CreatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
