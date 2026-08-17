package postgres

import (
	"database/sql"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, err
}

// TODO:
// Implement it
func (p *PostgresRepo) Save(roomID int, message []byte) error {
	return nil
}

// TODO:
// Same thing goes to this
func (p *PostgresRepo) GetRecent(roomID int, limit int) ([][]byte, error) {
	return nil, nil
}
