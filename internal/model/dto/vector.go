package dto

// VectorSearchRequest 向量检索请求，对齐 agent-server VectorSearchQuery。
type VectorSearchRequest struct {
	Content         string   `json:"content"`
	TopK            *int     `json:"top_k,omitempty"`
	KnowledgeType   *string  `json:"knowledge_type,omitempty"`
	KnowledgeDomain *string  `json:"knowledge_domain,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Source          *string  `json:"source,omitempty"`
	DocID           *string  `json:"doc_id,omitempty"`
	Status          *int     `json:"status,omitempty"`
	Version         *int     `json:"version,omitempty"`
	CreatedAtMin    *int     `json:"created_at_min,omitempty"`
	CreatedAtMax    *int     `json:"created_at_max,omitempty"`
	UpdatedAtMin    *int     `json:"updated_at_min,omitempty"`
	UpdatedAtMax    *int     `json:"updated_at_max,omitempty"`
}

// VectorSearchItem 单条向量检索结果。
type VectorSearchItem struct {
	ChunkID         string   `json:"chunk_id"`
	DocID           string   `json:"doc_id"`
	ChunkNo         int      `json:"chunk_no"`
	Score           *float64 `json:"score,omitempty"`
	Distance        *float64 `json:"distance,omitempty"`
	Content         string   `json:"content"`
	Title           string   `json:"title,omitempty"`
	KnowledgeType   string   `json:"knowledge_type,omitempty"`
	KnowledgeDomain string   `json:"knowledge_domain,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Source          string   `json:"source,omitempty"`
	Status          int      `json:"status,omitempty"`
	Version         int      `json:"version,omitempty"`
	CreatedAt       *int     `json:"created_at,omitempty"`
	UpdatedAt       *int     `json:"updated_at,omitempty"`
	Ext             string   `json:"ext,omitempty"`
}

// VectorSearchResponse 向量检索响应。
type VectorSearchResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []VectorSearchItem `json:"data"`
	Count   int                `json:"count"`
}
