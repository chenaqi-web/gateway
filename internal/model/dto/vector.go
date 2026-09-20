package dto

type CreateVectorCollectionRequest struct {
	Name string `json:"name"`
}

type DelVectorCollectionRequest struct {
	Name string `json:"name"`
}

type ListDocumentsRequest struct {
	Name     string `json:"name"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type ListDocumentsResponse struct {
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Items    []VectorDocument `json:"items"`
}

type DocsSearchRequest struct {
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

type Document struct {
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

type DocsSearchResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data []Document `json:"data"`
}

type Collection struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type VectorDocument struct {
	ChunkID       string `json:"chunk_id"`
	DocID         string `json:"doc_id"`
	ChunkNo       int    `json:"chunk_no"`
	Title         string `json:"title"`
	Source        string `json:"source"`
	KnowledgeType string `json:"knowledge_type"`
	Status        int    `json:"status"`
	CreatedAt     *int   `json:"created_at"`
	UpdatedAt     *int   `json:"updated_at"`
	Content       string `json:"content"`
}
