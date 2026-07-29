package dto

type SendEmailCodeRequest struct {
	Email   string `json:"email" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	EmailCode       string `json:"email_code" binding:"required"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar"`
	Sex             string `json:"sex"`
	Age             uint32 `json:"age"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordByEmailRequest struct {
	Email           string `json:"email" binding:"required"`
	EmailCode       string `json:"email_code" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
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

type RefreshResponse struct {
	AccessToken     string `json:"access_token"`
	AccessExpiresIn int64  `json:"access_expires_in"`
}

type RegisterResponse struct {
	User *AuthUserSummary `json:"user"`
}
