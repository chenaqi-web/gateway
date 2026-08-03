package dto

type SendEmailCodeRequest struct {
	Email   string `json:"email" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthUser struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Sex         string `json:"sex"`
	Age         uint32 `json:"age"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	AuthVersion uint64 `json:"auth_version"`
}

type LoginResponse struct {
	AccessToken     string    `json:"access_token"`
	AccessExpiresIn int       `json:"access_expires_in"`
	User            *AuthUser `json:"user"`
}
