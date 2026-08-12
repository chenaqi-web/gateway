package repo

import (
	"context"
	"errors"
	"time"

	"gateway/internal/model/entity"

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

func (r *AiChatRepo) ListSessionsByUser(ctx context.Context, userID string) ([]*entity.AiChatSession, error) {
	var list []*entity.AiChatSession
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC, id DESC").
		Find(&list).Error
	return list, err
}

func (r *AiChatRepo) UpdateSessionTitleByUser(ctx context.Context, userID, sessionID, title string) error {
	result := r.DB.WithContext(ctx).Model(&entity.AiChatSession{}).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Update("title", title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAiChatSessionNotFound
	}
	return nil
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

func (r *AiChatRepo) DeleteSessionByUser(ctx context.Context, userID, sessionID string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND session_id = ?", userID, sessionID).Delete(&entity.AiChatSession{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrAiChatSessionNotFound
		}
		return tx.Where("session_id = ?", sessionID).Delete(&entity.AiChatMessage{}).Error
	})
}
