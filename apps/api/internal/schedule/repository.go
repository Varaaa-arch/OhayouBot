package schedule

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

func (r *Repository) Create(ctx context.Context, arg db.CreateScheduleParams) (db.Schedule, error) {
	return r.q.CreateSchedule(ctx, arg)
}

func (r *Repository) List(ctx context.Context) ([]db.Schedule, error) {
	return r.q.GetSchedules(ctx)
}
