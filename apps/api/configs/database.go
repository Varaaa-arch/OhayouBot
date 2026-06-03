package configs

import (
	"database/sql"

	"github.com/Varaaa-arch/OhayouBot/generated/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(dsn string) (*db.Queries, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db.New(conn), nil
}
