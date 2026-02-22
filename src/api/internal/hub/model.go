package hub

import "time"

type Hub struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

type HubPost struct {
	ID        int       `json:"id"`
	HubID     int       `json:"hub_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
