package friends

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}

// Отправить запрос в друзья
func (s *Service) SendRequest(ctx context.Context, fromID, toID string) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO friends_requests (from_user_id, to_user_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, fromID, toID, time.Now())
	return err
}

// Принять запрос
func (s *Service) AcceptRequest(ctx context.Context, fromID, toID string) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO friends (user1_id, user2_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, fromID, toID, time.Now())
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(ctx, `
		DELETE FROM friends_requests
		WHERE from_user_id=$1 AND to_user_id=$2
	`, fromID, toID)
	return err
}

// Список друзей
type Friend struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Service) ListFriends(ctx context.Context, userID string) ([]Friend, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT u.id, u.username, u.email
		FROM friends f
		JOIN users u ON (u.id=f.user1_id OR u.id=f.user2_id) AND u.id<>$1
		WHERE f.user1_id=$1 OR f.user2_id=$1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := []Friend{}
	for rows.Next() {
		var f Friend
		if err := rows.Scan(&f.ID, &f.Username, &f.Email); err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}
