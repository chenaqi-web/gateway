package utils

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"gateway/internal/config"

	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
)

const (
	emailCodeTTL       = 5 * time.Minute
	emailCodeKeyPrefix = "auth:"
)

type Email struct {
	cfg   config.EmailConfig
	cache *redis.Client
}

func NewEmail(cfg *config.Config, cache *redis.Client) *Email {
	return &Email{cfg: cfg.Email, cache: cache}
}

func (s *Email) SendCode(ctx context.Context, email, purpose string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if purpose != "register" && purpose != "login" {
		return fmt.Errorf("unsupported verification purpose")
	}
	if strings.TrimSpace(s.cfg.Host) == "" || strings.TrimSpace(s.cfg.From) == "" || strings.TrimSpace(s.cfg.Secret) == "" {
		return fmt.Errorf("email service is not configured")
	}

	codeKey := emailCodeKeyPrefix + purpose + ":" + email

	// 检查是否已有有效的验证码（可选）
	exists, err := s.cache.Exists(ctx, codeKey).Result()
	if err != nil {
		return fmt.Errorf("check existing code: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("verification code already sent, please check your email")
	}

	code, err := GenerateCode()
	if err != nil {
		return fmt.Errorf("generate email code: %w", err)
	}

	// 存储验证码，5分钟过期
	if err := s.cache.Set(ctx, codeKey, code, emailCodeTTL).Err(); err != nil {
		return fmt.Errorf("store email code: %w", err)
	}

	// 发送邮件
	if err := s.send(email, code); err != nil {
		// 发送失败，删除已存储的验证码
		_ = s.cache.Del(ctx, codeKey).Err()
		return fmt.Errorf("send email code: %w", err)
	}

	return nil
}

func (s *Email) send(to, code string) error {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	auth := smtp.PlainAuth("", s.cfg.From, s.cfg.Secret, s.cfg.Host)

	// 配置email对象
	e := email.NewEmail()

	// 格式化发件人地址为 "昵称 <邮箱>"
	e.From = fmt.Sprintf("%s <%s>", s.cfg.Nickname, s.cfg.From)

	// 设置收件人、主题和邮件内容
	e.To = []string{to}
	e.Subject = "Renai 邮箱验证码"
	e.HTML = []byte(`您好，<br/><br/>
您正在进行 Renai 的邮箱验证。验证码为：<br/><br/>
<strong style="font-size:24px;letter-spacing:6px;color:#2563eb">` + code + `</strong><br/><br/>
验证码 5 分钟内有效，请勿向任何人泄露。若非本人操作，请忽略此邮件。<br/><br/>
Renai 团队`)

	// 定义错误变量
	var err error

	if s.cfg.IsSSL {
		// 使用带 TLS 的邮件发送
		err = e.SendWithTLS(addr, auth, &tls.Config{ServerName: s.cfg.Host})
	} else {
		// 使用普通的邮件发送
		err = e.Send(addr, auth)
	}

	return err
}

func GenerateCode() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", number.Int64()), nil
}

func (s *Email) VerifyCode(ctx context.Context, email, code, purpose string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	code = strings.TrimSpace(code)
	if email == "" || code == "" {
		return fmt.Errorf("email and code are required")
	}

	key := emailCodeKeyPrefix + purpose + ":" + email
	storedCode, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("email code is invalid or expired")
		}
		return fmt.Errorf("get email code: %w", err)
	}
	if storedCode != code {
		return fmt.Errorf("email code is invalid or expired")
	}
	return s.cache.Del(ctx, key).Err()
}
