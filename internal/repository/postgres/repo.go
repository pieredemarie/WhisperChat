package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"whisperchat/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(dsn string) (*PostgresRepo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PostgresRepo{
		DB: db,
	}, nil
}

func (p *PostgresRepo) Save(msg *domain.Message) error {
	query := `
		INSERT INTO messages (room_id, display_name, content, created_at)
		VALUES ($1,$2,$3,$4) 
	`

	_, err := p.DB.Exec(
		query,
		msg.RoomID,
		msg.DisplayName,
		msg.Content,
		msg.CreatedAt,
	)

	return err
}

func (p *PostgresRepo) GetRecent(roomID string, limit int) ([]*domain.Message, error) {
	query := `
		SELECT room_id, display_name, content, created_at
		FROM messages 
		WHERE room_id = $1 
		ORDER BY created_at ASC  
		LIMIT $2
	`

	rows, err := p.DB.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message

	for rows.Next() {
		msg := &domain.Message{}
		err := rows.Scan(&msg.RoomID, &msg.DisplayName, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (p *PostgresRepo) Close() error {
	return p.DB.Close()
}
