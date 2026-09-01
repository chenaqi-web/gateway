package application

import (
	"context"
	"errors"
	"gateway/internal/infras/clog"
	"strings"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/dto"
	"gateway/internal/utils"

	"go.uber.org/zap"
)

type AuthService struct {
	cfg *config.Config
	rpc *rpc.Client
	log *clog.Log

	blackList *cache.Blacklist
	email     *utils.Email
}

func NewAuthService(
	rpcClient *rpc.Client,
	blackList *cache.Blacklist,
	cfg *config.Config,
	log *clog.Log,
) *AuthService {
	return &AuthService{
		cfg:       cfg,
		rpc:       rpcClient,
		log:       log,
		blackList: blackList,
		email:     utils.NewEmail(cfg, blackList.Cache),
	}
}

func (s *AuthService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) error {
	err := s.email.SendCode(ctx, req.Email, req.Purpose)
	if err != nil {
		s.log.Error("AuthService/SendEmailCode error",
			zap.String("Purpose:", req.Purpose),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "register"); err != nil {
		s.log.Error("AuthService/Register error", zap.Error(err))
		return err
	}

	if _, err := s.rpc.GetAuthClient().Register(ctx, &authpb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		s.log.Error("AuthService/Register error", zap.Error(err))
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {
	resp, err := s.rpc.GetAuthClient().Login(ctx, &authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		s.log.Error("AuthService/Login error", zap.Error(err))
		return nil, "", err
	}
	return s.createLoginResult(resp)
}

func (s *AuthService) EmailLogin(ctx context.Context, req dto.EmailLoginRequest) (*dto.LoginResponse, string, error) {
	// 1. 校验验证码是否正确
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "login"); err != nil {
		s.log.Error("AuthService/EmailLogin error", zap.Error(err))
		return nil, "", err
	}

	// 2.调用rpc服务
	resp, err := s.rpc.GetAuthClient().EmailLogin(ctx, &authpb.EmailLoginRequest{Email: req.Email})
	if err != nil {
		s.log.Error("AuthService/EmailLogin error", zap.Error(err))
		return nil, "", err
	}

	// 3.准备token
	return s.createLoginResult(resp)
}

func (s *AuthService) Logout(ctx context.Context, authorization, refreshToken string) error {
	// 直接从 authorization 提取 token
	const prefix = "Bearer"
	accessToken := strings.TrimPrefix(authorization, prefix)

	if accessToken == authorization { // 没有 Bearer 前缀
		return errors.New("invalid authorization header")
	}

	if err := s.blackList.AddToken(ctx, accessToken, s.cfg.Auth.AccessExpire); err != nil {
		s.log.Error("AuthService/Logout error", zap.Error(err))
		return err
	}

	if err := s.blackList.AddToken(ctx, refreshToken, s.cfg.Auth.RefreshExpire); err != nil {
		s.log.Error("AuthService/Logout error", zap.Error(err))
		return err
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code, "forgot_password"); err != nil {
		s.log.Error("AuthService/ForgotPassword error", zap.Error(err))
		return err
	}

	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	if _, err := s.rpc.GetAuthClient().ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Email:           req.Email,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}); err != nil {
		s.log.Error("AuthService/ForgotPassword error", zap.Error(err))
		return err
	}
	return nil
}

// =====================================================================================================================

func (s *AuthService) createLoginResult(user *authpb.LoginResponse) (*dto.LoginResponse, string, error) {
	claims := utils.JWTClaims{
		UserID: user.GetId(),
		Role:   user.GetRole(),
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
		User: &dto.User{
			ID:       user.GetId(),
			Username: user.GetUsername(),
			Email:    user.GetEmail(),
			Phone:    user.GetPhone(),
			Avatar:   user.GetAvatar(),
			Sex:      user.GetSex(),
			Age:      user.GetAge(),
			Role:     user.GetRole(),
			Status:   user.GetStatus(),
		},
	}, refreshToken, nil
}
