package dto

type SendEmailCodeRequest struct {
	Email   string `json:"email" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthUserSummary struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	AccessToken     string           `json:"access_token"`
	AccessExpiresIn int64            `json:"access_expires_in"`
	User            *AuthUserSummary `json:"user"`
}
