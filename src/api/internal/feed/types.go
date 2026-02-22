package feed

import "time"

// Пост в ленте
type FeedPost struct {
	ID        int       `json:"id"`
	HubID     int       `json:"hub_id"`
	HubName   string    `json:"hub_name"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}

// Топ-хаб
type TopHub struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Members     int    `json:"members"`
}
