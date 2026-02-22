package feed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}

// Получаем последние посты из всех хабов, где пользователь состоит
func (s *Service) GetUserFeed(ctx context.Context, userID string, limit int) ([]FeedPost, error) {
	rows, err := s.DB.Query(ctx, `
        SELECT hp.id, hp.hub_id, hp.author_id, hp.content, hp.likes, hp.created_at, h.name AS hub_name
        FROM hub_posts hp
        INNER JOIN hub_members hm ON hm.hub_id = hp.hub_id
        INNER JOIN hubs h ON h.id = hp.hub_id
        WHERE hm.user_id=$1
        ORDER BY hp.created_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []FeedPost{}
	for rows.Next() {
		var p FeedPost
		if err := rows.Scan(&p.ID, &p.HubID, &p.AuthorID, &p.Content, &p.Likes, &p.CreatedAt, &p.HubName); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// Можно добавить метод топ-хабов по кол-ву участников
func (s *Service) GetTopHubs(ctx context.Context, limit int) ([]TopHub, error) {
	rows, err := s.DB.Query(ctx, `
        SELECT h.id, h.name, h.description, h.avatar, COUNT(hm.user_id) AS members
        FROM hubs h
        LEFT JOIN hub_members hm ON hm.hub_id = h.id
        GROUP BY h.id
        ORDER BY members DESC
        LIMIT $1
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hubs := []TopHub{}
	for rows.Next() {
		var h TopHub
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Avatar, &h.Members); err != nil {
			return nil, err
		}
		hubs = append(hubs, h)
	}
	return hubs, nil
}
