package dto

type AiChatStatusResponse struct {
	AssistantEnabled bool `json:"assistantEnabled"`
}

type UpdateAiChatStatusRequest struct {
	AssistantEnabled bool `json:"assistantEnabled"`
}

type AiSettingsResponse struct {
	AssistantEnabled bool   `json:"assistantEnabled"`
	Provider         string `json:"provider"`
	ModelID          string `json:"modelId"`
	APIKey           string `json:"apiKey"`
	APIKeyConfigured bool   `json:"apiKeyConfigured"`
}

type UpdateAiSettingsRequest struct {
	Provider string `json:"provider" binding:"required"`
	ModelID  string `json:"modelId" binding:"required"`
	APIKey   string `json:"apiKey"`
}
