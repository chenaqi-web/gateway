package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/dto"
	"gateway/internal/utils"
)

var (
	ErrInvalidOrExpiredVerificationCode = errors.New("invalid or expired verification code")
	ErrPasswordsDoNotMatch              = errors.New("passwords do not match")
)

type AuthService struct {
	rpc          *rpc.Client
	jwtBlackList *cache.JwtBlacklist
	email        *utils.Email
	cfg          config.AuthConfig
}

func NewAuthService(rpcClient *rpc.Client, jwtBlackList *cache.JwtBlacklist, cfg *config.Config) *AuthService {
	return &AuthService{
		rpc:          rpcClient,
		jwtBlackList: jwtBlackList,
		email:        utils.NewEmail(cfg, jwtBlackList.Cache),
		cfg:          cfg.Auth,
	}
}

func (s *AuthService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5e9)
	defer cancel()
	return s.email.SendCode(ctx, req.Email, req.Purpose)
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "register"); err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return fmt.Errorf("%w: %v", ErrInvalidOrExpiredVerificationCode, err)
		}
		return err
	}
	_, err := s.rpc.GetAuthClient().Register(ctx, &authpb.RegisterRequest{
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
	})
	return err
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {
	resp, err := s.rpc.GetAuthClient().Login(ctx, &authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, "", err
	}
	if resp == nil || resp.GetUser() == nil {
		return nil, "", errors.New("invalid core response")
	}
	return s.createLoginResult(resp.GetUser())
}

func (s *AuthService) EmailLogin(ctx context.Context, req dto.EmailLoginRequest) (*dto.LoginResponse, string, error) {
	// 1. 校验验证码是否正确
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "login"); err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return nil, "", fmt.Errorf("%w: %v", ErrInvalidOrExpiredVerificationCode, err)
		}
		return nil, "", err
	}

	// 2.调用rpc服务
	resp, err := s.rpc.GetAuthClient().EmailLogin(ctx, &authpb.EmailLoginRequest{Email: req.Email})
	if err != nil {
		return nil, "", err
	}
	if resp == nil || resp.GetUser() == nil {
		return nil, "", errors.New("invalid core response")
	}

	// 3.准备token
	return s.createLoginResult(resp.GetUser())
}

func (s *AuthService) createLoginResult(user *authpb.UserInfo) (*dto.LoginResponse, string, error) {
	claims := utils.JWTClaims{
		UserID:      user.GetId(),
		Role:        user.GetRole(),
		AuthVersion: user.GetAuthVersion(),
	}
	accessToken, err := utils.CreateAccessToken([]byte(s.cfg.JWTSecret), claims, s.cfg.AccessExpire)
	if err != nil {
		return nil, "", err
	}
	refreshToken, err := utils.CreateRefreshToken([]byte(s.cfg.JWTSecret), claims, s.cfg.RefreshExpire)
	if err != nil {
		return nil, "", err
	}

	return &dto.LoginResponse{
		AccessToken:     accessToken,
		AccessExpiresIn: s.cfg.AccessExpire,
		User: &dto.AuthUser{
			ID:          user.GetId(),
			Username:    user.GetUsername(),
			Email:       user.GetEmail(),
			Phone:       user.GetPhone(),
			Avatar:      user.GetAvatar(),
			Sex:         user.GetSex(),
			Age:         user.GetAge(),
			Role:        user.GetRole(),
			Status:      user.GetStatus(),
			AuthVersion: user.GetAuthVersion(),
		},
	}, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, authorization, refreshToken string) error {
	accessToken := accessTokenFromAuthorization(authorization)
	if strings.TrimSpace(accessToken) != "" {
		if err := s.jwtBlackList.BlacklistToken(ctx, accessToken, s.cfg.AccessExpire); err != nil {
			return err
		}
	}
	if strings.TrimSpace(refreshToken) != "" {
		return s.jwtBlackList.BlacklistToken(ctx, refreshToken, s.cfg.RefreshExpire)
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordsDoNotMatch
	}
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "reset_password"); err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return fmt.Errorf("%w: %v", ErrInvalidOrExpiredVerificationCode, err)
		}
		return err
	}

	_, err := s.rpc.GetAuthClient().ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Email:           req.Email,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	})
	return err
}

func (s *AuthService) RefreshCookieConfig() config.AuthConfig { return s.cfg }

func accessTokenFromAuthorization(authorization string) string {
	fields := strings.Fields(authorization)
	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
		return fields[1]
	}
	return ""
}
