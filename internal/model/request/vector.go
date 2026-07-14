package request

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
