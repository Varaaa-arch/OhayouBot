package message

import (
	"context"

	"github.com/Varaaa-arch/OhayouBot/generated/db"
	"github.com/google/uuid"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

func (r *Repository) Create(ctx context.Context, arg db.CreateMessageLogParams) (db.MessageLog, error) {
	return r.q.CreateMessageLog(ctx, arg)
}

func (r *Repository) ListByContact(ctx context.Context, contactID uuid.NullUUID) ([]db.MessageLog, error) {
	return r.q.GetMessageLogsByContact(ctx, contactID)
}
