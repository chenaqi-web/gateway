package dto

type UploadResponse struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Provider string `json:"provider"`
}

type DeleteUploadRequest struct {
	Key string `json:"key" binding:"required"`
}
