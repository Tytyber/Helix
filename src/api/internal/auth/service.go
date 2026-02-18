package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB        *pgxpool.Pool
	JWTSecret string
}

func NewService(db *pgxpool.Pool, secret string) *Service {
	return &Service{
		DB:        db,
		JWTSecret: secret,
	}
}

// Register user
func (s *Service) Register(username, email, password string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	var userID string
	err = s.DB.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		username, email, hash,
	).Scan(&userID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := GenerateAccessToken(userID, s.JWTSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, _ := GenerateRefreshToken()
	s.SaveRefreshToken(userID, refreshToken, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

// Login user
func (s *Service) Login(email, password string) (string, string, error) {
	var userID string
	var hash string

	err := s.DB.QueryRow(context.Background(),
		`SELECT id, password_hash FROM users WHERE email=$1`,
		email,
	).Scan(&userID, &hash)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := GenerateAccessToken(userID, s.JWTSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, _ := GenerateRefreshToken()
	s.SaveRefreshToken(userID, refreshToken, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

// Refresh token methods
func (s *Service) SaveRefreshToken(userID, token string, duration time.Duration) error {
	hash := sha256.Sum256([]byte(token))
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hex.EncodeToString(hash[:]), time.Now().Add(duration),
	)
	return err
}

func (s *Service) ValidateRefreshToken(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	var userID string
	var expiresAt time.Time

	err := s.DB.QueryRow(context.Background(),
		`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash=$1`,
		hex.EncodeToString(hash[:]),
	).Scan(&userID, &expiresAt)

	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("refresh token expired")
	}

	return userID, nil
}

func (s *Service) DeleteRefreshToken(token string) error {
	hash := sha256.Sum256([]byte(token))
	_, err := s.DB.Exec(context.Background(),
		`DELETE FROM refresh_tokens WHERE token_hash=$1`,
		hex.EncodeToString(hash[:]),
	)
	return err
}

// AddReputation позволяет поставить +1 или -1 другому пользователю
func (s *Service) AddReputation(userID, fromUserID string, score int) error {
	if score != 1 && score != -1 {
		return errors.New("score must be +1 or -1")
	}

	// Вставляем или обновляем оценку
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO user_reputation (user_id, from_user_id, score) 
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, from_user_id) 
		 DO UPDATE SET score = EXCLUDED.score`,
		userID, fromUserID, score,
	)
	if err != nil {
		return err
	}

	// Обновляем суммарную репутацию пользователя
	_, err = s.DB.Exec(context.Background(),
		`UPDATE users
		 SET reputation = (SELECT COALESCE(SUM(score),0) FROM user_reputation WHERE user_id=$1)
		 WHERE id=$1`,
		userID,
	)

	return err
}
