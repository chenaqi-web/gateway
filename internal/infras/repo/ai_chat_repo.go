package repo

import (
	"context"
	"errors"
	"time"

	"backend/gateway/internal/model/entity"

	"gorm.io/gorm"
)

var ErrAiChatSessionNotFound = errors.New("ai chat session not found")

type AiChatRepo struct {
	*DBClient
}

func NewAiChatRepo(client *DBClient) *AiChatRepo {
	return &AiChatRepo{DBClient: client}
}

func (r *AiChatRepo) CreateSession(ctx context.Context, session *entity.AiChatSession) error {
	return r.DB.WithContext(ctx).Create(session).Error
}

func (r *AiChatRepo) GetSessionByUser(ctx context.Context, userID, sessionID string) (*entity.AiChatSession, error) {
	var session entity.AiChatSession
	err := r.DB.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAiChatSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AiChatRepo) ListSessionsByUser(ctx context.Context, userID string, page, pageSize int) ([]*entity.AiChatSession, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []*entity.AiChatSession
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC, id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error
	return list, err
}

func (r *AiChatRepo) TouchSession(ctx context.Context, sessionID string) error {
	return r.DB.WithContext(ctx).Model(&entity.AiChatSession{}).
		Where("session_id = ?", sessionID).
		Update("updated_at", time.Now()).Error
}

func (r *AiChatRepo) CreateMessages(ctx context.Context, messages []*entity.AiChatMessage) error {
	if len(messages) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&messages).Error
}

func (r *AiChatRepo) ListMessagesBySession(ctx context.Context, sessionID string) ([]*entity.AiChatMessage, error) {
	var list []*entity.AiChatMessage
	err := r.DB.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC, id ASC").
		Find(&list).Error
	return list, err
}
