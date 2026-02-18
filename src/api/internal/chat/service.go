package chat

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}

// Создать чат 1-на-1
func (s *Service) CreateChat(user1ID, user2ID string) (string, error) {
	if user1ID == user2ID {
		return "", errors.New("cannot create chat with yourself")
	}

	var chatID string
	err := s.DB.QueryRow(context.Background(),
		`INSERT INTO chats (user1_id, user2_id) 
		 VALUES ($1, $2)
		 ON CONFLICT (user1_id, user2_id) DO NOTHING
		 RETURNING id`,
		user1ID, user2ID,
	).Scan(&chatID)

	// Если чат уже существует, достаем id
	if err != nil {
		err = s.DB.QueryRow(context.Background(),
			`SELECT id FROM chats WHERE (user1_id=$1 AND user2_id=$2) OR (user1_id=$2 AND user2_id=$1)`,
			user1ID, user2ID,
		).Scan(&chatID)
		if err != nil {
			return "", err
		}
	}

	return chatID, nil
}

// Получить все сообщения чата
func (s *Service) GetMessages(chatID string) ([]Message, error) {
	rows, err := s.DB.Query(context.Background(),
		`SELECT id, sender_id, content, created_at FROM messages WHERE chat_id=$1 ORDER BY created_at ASC`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// Отправить сообщение
func (s *Service) SendMessage(chatID, senderID, content string) error {
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO messages (chat_id, sender_id, content) VALUES ($1, $2, $3)`,
		chatID, senderID, content,
	)
	return err
}

type Message struct {
	ID        string
	SenderID  string
	Content   string
	CreatedAt string
}
