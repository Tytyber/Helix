package group

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

// Создать группу
func (s *Service) CreateGroup(name, ownerID string) (string, error) {
	var groupID string
	err := s.DB.QueryRow(context.Background(),
		`INSERT INTO groups (name, owner_id) VALUES ($1, $2) RETURNING id`,
		name, ownerID,
	).Scan(&groupID)
	return groupID, err
}

// Добавить участника
func (s *Service) JoinGroup(groupID, userID string) error {
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		groupID, userID,
	)
	return err
}

// Выйти из группы
func (s *Service) LeaveGroup(groupID, userID string) error {
	_, err := s.DB.Exec(context.Background(),
		`DELETE FROM group_members WHERE group_id=$1 AND user_id=$2`,
		groupID, userID,
	)
	return err
}

// Отправить сообщение в группу
func (s *Service) SendMessage(groupID, senderID, content string) error {
	_, err := s.DB.Exec(context.Background(),
		`INSERT INTO group_messages (group_id, sender_id, content) VALUES ($1, $2, $3)`,
		groupID, senderID, content,
	)
	return err
}

// Получить сообщения группы
func (s *Service) GetMessages(groupID string) ([]Message, error) {
	rows, err := s.DB.Query(context.Background(),
		`SELECT id, sender_id, content, created_at FROM group_messages WHERE group_id=$1 ORDER BY created_at ASC`,
		groupID,
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

type Message struct {
	ID        string
	SenderID  string
	Content   string
	CreatedAt string
}
