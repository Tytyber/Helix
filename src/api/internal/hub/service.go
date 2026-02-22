package hub

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}

func (s *Service) GetAllHubs(ctx context.Context) ([]Hub, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, name, description, avatar, created_at FROM hubs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hubs := []Hub{}
	for rows.Next() {
		var h Hub
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Avatar, &h.CreatedAt); err != nil {
			return nil, err
		}
		hubs = append(hubs, h)
	}
	return hubs, nil
}

func (s *Service) GetHub(ctx context.Context, hubID int) (*Hub, error) {
	var h Hub
	err := s.DB.QueryRow(ctx, "SELECT id, name, description, avatar, created_at FROM hubs WHERE id=$1", hubID).
		Scan(&h.ID, &h.Name, &h.Description, &h.Avatar, &h.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *Service) JoinHub(ctx context.Context, hubID int, userID string) error {
	_, err := s.DB.Exec(ctx,
		"INSERT INTO hub_members (hub_id, user_id, joined_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		hubID, userID, time.Now())
	return err
}

func (s *Service) LeaveHub(ctx context.Context, hubID int, userID string) error {
	_, err := s.DB.Exec(ctx, "DELETE FROM hub_members WHERE hub_id=$1 AND user_id=$2", hubID, userID)
	return err
}

func (s *Service) CreatePost(ctx context.Context, hubID int, authorID, content string) (*HubPost, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}
	var postID int
	err := s.DB.QueryRow(ctx,
		"INSERT INTO hub_posts (hub_id, author_id, content, likes, created_at) VALUES ($1,$2,$3,0,$4) RETURNING id",
		hubID, authorID, content, time.Now()).Scan(&postID)
	if err != nil {
		return nil, err
	}
	return &HubPost{
		ID:        postID,
		HubID:     hubID,
		AuthorID:  authorID,
		Content:   content,
		Likes:     0,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) GetHubPosts(ctx context.Context, hubID int) ([]HubPost, error) {
	rows, err := s.DB.Query(ctx, `
        SELECT id, hub_id, author_id, content, likes, created_at
        FROM hub_posts
        WHERE hub_id=$1
        ORDER BY created_at DESC
    `, hubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []HubPost{}
	for rows.Next() {
		var p HubPost
		if err := rows.Scan(&p.ID, &p.HubID, &p.AuthorID, &p.Content, &p.Likes, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
