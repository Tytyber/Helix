package calls

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

type Call struct {
	ID        int       `json:"id"`
	CallerID  string    `json:"caller_id"`
	CalleeID  string    `json:"callee_id"`
	Status    string    `json:"status"` // pending, accepted, rejected, ended
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) CreateCall(ctx context.Context, callerID, calleeID string) (*Call, error) {
	var callID int
	err := s.DB.QueryRow(ctx, `
		INSERT INTO calls (caller_id, callee_id, status, created_at)
		VALUES ($1, $2, 'pending', $3) RETURNING id
	`, callerID, calleeID, time.Now()).Scan(&callID)
	if err != nil {
		return nil, err
	}
	return &Call{
		ID:        callID,
		CallerID:  callerID,
		CalleeID:  calleeID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) UpdateCallStatus(ctx context.Context, callID int, status string) error {
	_, err := s.DB.Exec(ctx, "UPDATE calls SET status=$1 WHERE id=$2", status, callID)
	return err
}
