package friends

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

type FriendRequest struct {
	ID        int       `json:"id"`
	FromID    string    `json:"from_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Список друзей
type Friend struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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

func (s *Service) DeclineRequest(ctx context.Context, fromID, toID string) error {
	_, err := s.DB.Exec(ctx, `
		DELETE FROM friends_requests
		WHERE from_user_id=$1 AND to_user_id=$2
	`, fromID, toID)
	return err
}

func (s *Service) IncomingRequests(ctx context.Context, userID string) ([]FriendRequest, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT fr.id, u.id, u.username, u.email, fr.created_at
		FROM friends_requests fr
		JOIN users u ON u.id = fr.from_user_id
		WHERE fr.to_user_id = $1
		ORDER BY fr.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []FriendRequest

	for rows.Next() {
		var r FriendRequest
		if err := rows.Scan(
			&r.ID,
			&r.FromID,
			&r.Username,
			&r.Email,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}
