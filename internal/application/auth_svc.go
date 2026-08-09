package application

import (
	"context"
	"errors"
	"strings"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/dto"
	"gateway/internal/utils"
)

type AuthService struct {
	rpc          *rpc.Client
	jwtBlackList *cache.JwtBlacklist
	email        *utils.Email
	cfg          *config.Config
}

func NewAuthService(rpcClient *rpc.Client, jwtBlackList *cache.JwtBlacklist, cfg *config.Config) *AuthService {
	return &AuthService{
		rpc:          rpcClient,
		jwtBlackList: jwtBlackList,
		email:        utils.NewEmail(cfg, jwtBlackList.Cache),
		cfg:          cfg,
	}
}

func (s *AuthService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5e9)
	defer cancel()
	return s.email.SendCode(ctx, req.Email, req.Purpose)
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "register"); err != nil {
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
	accessToken, err := utils.CreateAccessToken([]byte(s.cfg.Auth.JWTSecret), claims, s.cfg.Auth.AccessExpire)
	if err != nil {
		return nil, "", err
	}
	refreshToken, err := utils.CreateRefreshToken([]byte(s.cfg.Auth.JWTSecret), claims, s.cfg.Auth.RefreshExpire)
	if err != nil {
		return nil, "", err
	}

	return &dto.LoginResponse{
		AccessToken:     accessToken,
		AccessExpiresIn: s.cfg.Auth.AccessExpire,
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
	// 直接从 authorization 提取 token
	const prefix = "Bearer "
	accessToken := strings.TrimPrefix(authorization, prefix)

	if accessToken == authorization { // 没有 Bearer 前缀
		return errors.New("invalid authorization header")
	}

	if err := s.jwtBlackList.BlacklistToken(ctx, accessToken, s.cfg.Auth.AccessExpire); err != nil {
		return err
	}
	return s.jwtBlackList.BlacklistToken(ctx, refreshToken, s.cfg.Auth.RefreshExpire)
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "forgot_password"); err != nil {
		return err
	}

	_, err := s.rpc.GetAuthClient().ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Email:           req.Email,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	})
	return err
}
