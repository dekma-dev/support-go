package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"support-go/backend/internal/ticket"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (repository *Repository) Create(value ticket.Ticket) error {
	const query = `
		INSERT INTO tickets (
			id, public_id, title, description, status, priority,
			requester_id, assignee_id, sla_due_at, created_at, updated_at, closed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
	`

	_, err := repository.pool.Exec(context.Background(), query,
		value.ID,
		value.PublicID,
		value.Title,
		value.Description,
		value.Status,
		value.Priority,
		value.RequesterID,
		nullString(value.AssigneeID),
		value.SLADueAt,
		value.CreatedAt,
		value.UpdatedAt,
		value.ClosedAt,
	)

	return err
}

func (repository *Repository) GetByID(id string) (ticket.Ticket, error) {
	const query = `
		SELECT
			id, public_id, title, description, status, priority,
			requester_id, assignee_id, sla_due_at, created_at, updated_at, closed_at
		FROM tickets
		WHERE id = $1
	`

	var result ticket.Ticket
	var assigneeID *string
	err := repository.pool.QueryRow(context.Background(), query, id).Scan(
		&result.ID,
		&result.PublicID,
		&result.Title,
		&result.Description,
		&result.Status,
		&result.Priority,
		&result.RequesterID,
		&assigneeID,
		&result.SLADueAt,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.ClosedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ticket.Ticket{}, ticket.ErrNotFound
	}
	if err != nil {
		return ticket.Ticket{}, err
	}

	if assigneeID != nil {
		result.AssigneeID = *assigneeID
	}

	return result, nil
}

func (repository *Repository) List() []ticket.Ticket {
	const query = `
		SELECT
			id, public_id, title, description, status, priority,
			requester_id, assignee_id, sla_due_at, created_at, updated_at, closed_at
		FROM tickets
		ORDER BY created_at DESC
	`

	rows, err := repository.pool.Query(context.Background(), query)
	if err != nil {
		return []ticket.Ticket{}
	}
	defer rows.Close()

	items := make([]ticket.Ticket, 0)
	for rows.Next() {
		var item ticket.Ticket
		var assigneeID *string
		if scanErr := rows.Scan(
			&item.ID,
			&item.PublicID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Priority,
			&item.RequesterID,
			&assigneeID,
			&item.SLADueAt,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ClosedAt,
		); scanErr != nil {
			return []ticket.Ticket{}
		}

		if assigneeID != nil {
			item.AssigneeID = *assigneeID
		}

		items = append(items, item)
	}

	return items
}

func (repository *Repository) ListWithFilter(ctx context.Context, filter ticket.ListFilter) ([]ticket.Ticket, int, error) {
	where, args := buildWhereClause(filter)
	orderBy := sortClause(filter.Sort)

	// Count query (no limit/offset)
	countQuery := "SELECT COUNT(*) FROM tickets" + where
	var total int
	if err := repository.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Add limit/offset to args for the data query
	dataArgs := append([]any{}, args...)
	dataArgs = append(dataArgs, filter.Limit, filter.Offset)
	limitPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+2)

	dataQuery := `
		SELECT
			id, public_id, title, description, status, priority,
			requester_id, assignee_id, sla_due_at, created_at, updated_at, closed_at
		FROM tickets` + where + orderBy + " LIMIT " + limitPlaceholder + " OFFSET " + offsetPlaceholder

	rows, err := repository.pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]ticket.Ticket, 0)
	for rows.Next() {
		var item ticket.Ticket
		var assigneeID *string
		if scanErr := rows.Scan(
			&item.ID,
			&item.PublicID,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.Priority,
			&item.RequesterID,
			&assigneeID,
			&item.SLADueAt,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ClosedAt,
		); scanErr != nil {
			return nil, 0, scanErr
		}

		if assigneeID != nil {
			item.AssigneeID = *assigneeID
		}

		items = append(items, item)
	}

	return items, total, rows.Err()
}

func buildWhereClause(filter ticket.ListFilter) (string, []any) {
	var clauses []string
	var args []any
	argN := 0

	if len(filter.Statuses) > 0 {
		placeholders := make([]string, 0, len(filter.Statuses))
		for _, s := range filter.Statuses {
			argN++
			placeholders = append(placeholders, fmt.Sprintf("$%d", argN))
			args = append(args, string(s))
		}
		clauses = append(clauses, "status IN ("+strings.Join(placeholders, ",")+")")
	}

	if len(filter.Priorities) > 0 {
		placeholders := make([]string, 0, len(filter.Priorities))
		for _, p := range filter.Priorities {
			argN++
			placeholders = append(placeholders, fmt.Sprintf("$%d", argN))
			args = append(args, string(p))
		}
		clauses = append(clauses, "priority IN ("+strings.Join(placeholders, ",")+")")
	}

	if filter.AssigneeID != "" {
		argN++
		clauses = append(clauses, fmt.Sprintf("assignee_id = $%d", argN))
		args = append(args, filter.AssigneeID)
	}

	if filter.Search != "" {
		argN++
		clauses = append(clauses, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argN, argN))
		args = append(args, "%"+filter.Search+"%")
	}

	if len(clauses) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(clauses, " AND "), args
}

func sortClause(sort ticket.SortOrder) string {
	switch sort {
	case ticket.SortCreatedAtAsc:
		return " ORDER BY created_at ASC"
	case ticket.SortUpdatedAtDesc:
		return " ORDER BY updated_at DESC"
	case ticket.SortUpdatedAtAsc:
		return " ORDER BY updated_at ASC"
	default:
		return " ORDER BY created_at DESC"
	}
}

func (repository *Repository) Update(value ticket.Ticket) error {
	const query = `
		UPDATE tickets
		SET
			public_id = $2,
			title = $3,
			description = $4,
			status = $5,
			priority = $6,
			requester_id = $7,
			assignee_id = $8,
			sla_due_at = $9,
			created_at = $10,
			updated_at = $11,
			closed_at = $12
		WHERE id = $1
	`

	commandTag, err := repository.pool.Exec(context.Background(), query,
		value.ID,
		value.PublicID,
		value.Title,
		value.Description,
		value.Status,
		value.Priority,
		value.RequesterID,
		nullString(value.AssigneeID),
		value.SLADueAt,
		value.CreatedAt,
		value.UpdatedAt,
		value.ClosedAt,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ticket.ErrNotFound
	}

	return nil
}

func nullString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
