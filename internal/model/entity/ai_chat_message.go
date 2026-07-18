package entity

import "time"

type AiChatMessage struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time  `gorm:"index:idx_ai_chat_message_session_created;comment:创建时间"`
	SessionID string     `gorm:"size:64;not null;index:idx_ai_chat_message_session_created,priority:1;comment:会话ID"`
	UserID    string     `gorm:"size:64;not null;comment:用户ID"`
	Role      AiChatRole `gorm:"size:16;not null;comment:消息角色 user/assistant"`
	Content   string     `gorm:"type:text;not null;comment:消息内容"`
	AiModel   string     `gorm:"size:64;not null;default:'';comment:使用的模型"`
}

func (AiChatMessage) TableName() string {
	return "ai_chat_message"
}

type AiChatRole string

const (
	AiChatRoleUser      AiChatRole = "user"
	AiChatRoleAssistant AiChatRole = "assistant"
)

func (r AiChatRole) String() string {
	return string(r)
}
