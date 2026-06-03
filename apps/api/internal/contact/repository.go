package contact

import (
	"context"

	"github.com/Varaaa-arch/OhayouBot/generated/db"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) *Repository {
	return &Repository{q: q}
}

func (r *Repository) Create(ctx context.Context, arg db.CreateContactParams) (db.Contact, error) {
	return r.q.CreateContact(ctx, arg)
}

func (r *Repository) List(ctx context.Context) ([]db.Contact, error) {
	return r.q.GetContacts(ctx)
}
