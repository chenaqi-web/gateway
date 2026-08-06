package application

import (
	"context"
	"errors"
	"strings"
	"time"

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
	cfg          config.AuthConfig
}

func NewAuthService(rpcClient *rpc.Client, jwtBlackList *cache.JwtBlacklist, cfg *config.Config) *AuthService {
	return &AuthService{rpc: rpcClient, jwtBlackList: jwtBlackList, email: utils.NewEmail(cfg, jwtBlackList.Cache), cfg: cfg.Auth}
}

func (s *AuthService) SendEmailCode(ctx context.Context, req dto.SendEmailCodeRequest) error {
	return s.email.SendCode(ctx, req.Email, req.Purpose)
}

func (s *AuthService) EmailLogin(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {
	if err := s.email.VerifyCode(ctx, req.Email, req.Code); err != nil {
		return nil, "", err
	}
	resp, err := s.rpc.GetAuthClient().EmailLogin(ctx, &authpb.EmailLoginRequest{Email: req.Email})
	if err != nil {
		return nil, "", err
	}
	if resp == nil || resp.GetUser() == nil {
		return nil, "", errors.New("invalid core response")
	}
	user := resp.GetUser()
	claims := utils.JWTClaims{UserID: user.GetId(), Role: user.GetRole(), AuthVersion: user.GetAuthVersion()}
	accessToken, err := utils.CreateAccessToken([]byte(s.cfg.JWTSecret), claims, s.cfg.AccessExpire)
	if err != nil {
		return nil, "", err
	}
	refreshToken, err := utils.CreateRefreshToken([]byte(s.cfg.JWTSecret), claims, s.cfg.RefreshExpire)
	if err != nil {
		return nil, "", err
	}
	return &dto.LoginResponse{
		AccessToken: accessToken, AccessExpiresIn: s.cfg.AccessExpire,
		User: &dto.AuthUser{ID: user.GetId(), Username: user.GetUsername(), Email: user.GetEmail(), Phone: user.GetPhone(), Avatar: user.GetAvatar(), Sex: user.GetSex(), Age: user.GetAge(), Role: user.GetRole(), Status: user.GetStatus(), AuthVersion: user.GetAuthVersion()},
	}, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
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

func (s *AuthService) Config() config.AuthConfig { return s.cfg }
