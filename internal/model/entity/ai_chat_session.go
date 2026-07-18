package entity

import "time"

type AiChatSession struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	UserID    string    `gorm:"size:64;not null;index:idx_ai_chat_session_user_id;comment:用户ID"`
	SessionID string    `gorm:"size:64;not null;uniqueIndex:uk_ai_chat_session_id;comment:会话ID"`
	Title     string    `gorm:"size:256;not null;default:'';comment:会话标题"`
}

func (AiChatSession) TableName() string {
	return "ai_chat_session"
}
