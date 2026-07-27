package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/gateway/internal/client/rpc/core-rpc/authpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeAuthRPC struct {
	sendEmailCode func(context.Context, *authpb.SendEmailCodeRequest) (*authpb.SendEmailCodeResponse, error)
	register      func(context.Context, *authpb.RegisterRequest) (*authpb.RegisterResponse, error)
	login         func(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error)
	emailLogin    func(context.Context, *authpb.EmailLoginRequest) (*authpb.LoginResponse, error)
	refresh       func(context.Context, *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error)
	logout        func(context.Context, *authpb.LogoutRequest) (*authpb.LogoutResponse, error)
	reset         func(context.Context, *authpb.ResetPasswordByEmailRequest) (*authpb.ResetPasswordByEmailResponse, error)
}

func (f *fakeAuthRPC) SendEmailCode(ctx context.Context, request *authpb.SendEmailCodeRequest, _ ...grpc.CallOption) (*authpb.SendEmailCodeResponse, error) {
	if f.sendEmailCode != nil {
		return f.sendEmailCode(ctx, request)
	}
	return &authpb.SendEmailCodeResponse{Success: true}, nil
}

func (f *fakeAuthRPC) Register(ctx context.Context, request *authpb.RegisterRequest, _ ...grpc.CallOption) (*authpb.RegisterResponse, error) {
	if f.register != nil {
		return f.register(ctx, request)
	}
	return &authpb.RegisterResponse{Success: true, User: testAuthUser()}, nil
}

func (f *fakeAuthRPC) Login(ctx context.Context, request *authpb.LoginRequest, _ ...grpc.CallOption) (*authpb.LoginResponse, error) {
	if f.login != nil {
		return f.login(ctx, request)
	}
	return testLoginResponse(), nil
}

func (f *fakeAuthRPC) EmailLogin(ctx context.Context, request *authpb.EmailLoginRequest, _ ...grpc.CallOption) (*authpb.LoginResponse, error) {
	if f.emailLogin != nil {
		return f.emailLogin(ctx, request)
	}
	return testLoginResponse(), nil
}

func (f *fakeAuthRPC) RefreshToken(ctx context.Context, request *authpb.RefreshTokenRequest, _ ...grpc.CallOption) (*authpb.RefreshTokenResponse, error) {
	if f.refresh != nil {
		return f.refresh(ctx, request)
	}
	return &authpb.RefreshTokenResponse{Tokens: &authpb.TokenPair{
		AccessToken: "new-access", RefreshToken: "new-refresh", AccessExpiresIn: 1200, RefreshExpiresIn: 604800,
	}}, nil
}

func (f *fakeAuthRPC) Logout(ctx context.Context, request *authpb.LogoutRequest, _ ...grpc.CallOption) (*authpb.LogoutResponse, error) {
	if f.logout != nil {
		return f.logout(ctx, request)
	}
	return &authpb.LogoutResponse{Success: true}, nil
}

func (f *fakeAuthRPC) ResetPasswordByEmail(ctx context.Context, request *authpb.ResetPasswordByEmailRequest, _ ...grpc.CallOption) (*authpb.ResetPasswordByEmailResponse, error) {
	if f.reset != nil {
		return f.reset(ctx, request)
	}
	return &authpb.ResetPasswordByEmailResponse{Success: true}, nil
}

func TestAuthControllerSendEmailCodeMapsPurpose(t *testing.T) {
	var received *authpb.SendEmailCodeRequest
	controller := newTestAuthController(t, &fakeAuthRPC{sendEmailCode: func(_ context.Context, request *authpb.SendEmailCodeRequest) (*authpb.SendEmailCodeResponse, error) {
		received = request
		return &authpb.SendEmailCodeResponse{Success: true}, nil
	}})
	recorder := performAuthRequest(controller.SendEmailCode, "/api/v1/auth/send-email-code", `{"email":"user@qq.com","purpose":"reset_password"}`, nil)
	if recorder.Code != http.StatusOK || received.GetPurpose() != authpb.EmailCodePurpose_EMAIL_CODE_PURPOSE_RESET_PASSWORD {
		t.Fatalf("response=%d %s request=%+v", recorder.Code, recorder.Body.String(), received)
	}
}

func TestAuthControllerRegisterReturnsSafeUserSummary(t *testing.T) {
	var received *authpb.RegisterRequest
	controller := newTestAuthController(t, &fakeAuthRPC{register: func(_ context.Context, request *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
		received = request
		return &authpb.RegisterResponse{Success: true, User: testAuthUser()}, nil
	}})
	body := `{"username":"user","email":"user@qq.com","password":"abc12345","confirm_password":"abc12345","email_code":"123456","role":"admin","status":"disabled","auth_version":99}`
	recorder := performAuthRequest(controller.Register, "/api/v1/auth/register", body, nil)
	if recorder.Code != http.StatusOK || received.GetUsername() != "user" || received.GetEmailCode() != "123456" {
		t.Fatalf("response=%d %s request=%+v", recorder.Code, recorder.Body.String(), received)
	}
	assertResponseOmitsSecrets(t, recorder.Body.String())
}

func TestAuthControllerLoginAndEmailLoginSetRefreshCookieWithoutReturningIt(t *testing.T) {
	tests := []struct {
		name    string
		handler func(*AuthController) gin.HandlerFunc
		path    string
		body    string
	}{
		{name: "username", handler: func(controller *AuthController) gin.HandlerFunc { return controller.Login }, path: "/api/v1/auth/login", body: `{"username":"user","password":"abc12345"}`},
		{name: "email", handler: func(controller *AuthController) gin.HandlerFunc { return controller.EmailLogin }, path: "/api/v1/auth/email-login", body: `{"email":"user@qq.com","password":"abc12345"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := newTestAuthController(t, &fakeAuthRPC{})
			recorder := performAuthRequest(test.handler(controller), test.path, test.body, nil)
			if recorder.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
			assertResponseOmitsSecrets(t, recorder.Body.String())
			var payload struct {
				Data struct {
					AccessToken     string `json:"access_token"`
					AccessExpiresIn int64  `json:"access_expires_in"`
				} `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if payload.Data.AccessToken != "access" || payload.Data.AccessExpiresIn != 1200 {
				t.Fatalf("login data = %+v", payload.Data)
			}
			assertRefreshCookie(t, recorder, "refresh")
		})
	}
}

func TestAuthControllerRefreshReadsCookieRotatesItAndHasNoRequestBody(t *testing.T) {
	var received string
	controller := newTestAuthController(t, &fakeAuthRPC{refresh: func(_ context.Context, request *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
		received = request.GetRefreshToken()
		return &authpb.RefreshTokenResponse{Tokens: &authpb.TokenPair{
			AccessToken: "new-access", RefreshToken: "new-refresh", AccessExpiresIn: 1200, RefreshExpiresIn: 604800,
		}}, nil
	}})
	recorder := performAuthRequest(controller.Refresh, "/api/v1/auth/refresh", "", &http.Cookie{Name: refreshCookieName, Value: "old-refresh"})
	if recorder.Code != http.StatusOK || received != "old-refresh" {
		t.Fatalf("status=%d body=%s refresh=%q", recorder.Code, recorder.Body.String(), received)
	}
	assertResponseOmitsSecrets(t, recorder.Body.String())
	assertRefreshCookie(t, recorder, "new-refresh")
}

func TestAuthControllerRefreshFailureClearsCookie(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid", err: status.Error(codes.Unauthenticated, "invalid refresh"), wantStatus: http.StatusUnauthorized},
		{name: "disabled", err: status.Error(codes.PermissionDenied, "user disabled"), wantStatus: http.StatusUnauthorized},
		{name: "core unavailable", err: status.Error(codes.Unavailable, "core unavailable"), wantStatus: http.StatusServiceUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := newTestAuthController(t, &fakeAuthRPC{refresh: func(context.Context, *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
				return nil, test.err
			}})
			recorder := performAuthRequest(controller.Refresh, "/api/v1/auth/refresh", "", &http.Cookie{Name: refreshCookieName, Value: "bad-refresh"})
			if recorder.Code != test.wantStatus {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
			assertClearedRefreshCookie(t, recorder)
		})
	}
}

func TestAuthControllerLogoutIsIdempotentAndAlwaysClearsCookie(t *testing.T) {
	controller := newTestAuthController(t, &fakeAuthRPC{logout: func(context.Context, *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
		return nil, status.Error(codes.Unauthenticated, "already gone")
	}})
	for _, cookie := range []*http.Cookie{nil, {Name: refreshCookieName, Value: "refresh"}} {
		recorder := performAuthRequest(controller.Logout, "/api/v1/auth/logout", "", cookie)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
		}
		assertClearedRefreshCookie(t, recorder)
	}
}

func TestAuthControllerLogoutCoreUnavailableStillClearsCookie(t *testing.T) {
	controller := newTestAuthController(t, &fakeAuthRPC{logout: func(context.Context, *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
		return nil, status.Error(codes.Unavailable, "core unavailable")
	}})
	recorder := performAuthRequest(controller.Logout, "/api/v1/auth/logout", "", &http.Cookie{Name: refreshCookieName, Value: "refresh"})
	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	assertClearedRefreshCookie(t, recorder)
}

func TestAuthControllerResetPasswordClearsCookieWithoutReturningTokens(t *testing.T) {
	controller := newTestAuthController(t, &fakeAuthRPC{})
	body := `{"email":"user@qq.com","email_code":"123456","new_password":"newabc123","confirm_password":"newabc123"}`
	recorder := performAuthRequest(controller.ResetPasswordByEmail, "/api/v1/auth/reset-password-by-email", body, &http.Cookie{Name: refreshCookieName, Value: "refresh"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	assertResponseOmitsSecrets(t, recorder.Body.String())
	assertClearedRefreshCookie(t, recorder)
}

func TestAuthControllerMapsLoginErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "credentials", err: status.Error(codes.Unauthenticated, "hidden"), wantStatus: http.StatusUnauthorized, wantCode: `"code":10002`},
		{name: "active session", err: status.Error(codes.AlreadyExists, "active"), wantStatus: http.StatusConflict, wantCode: `"code":10007`},
		{name: "disabled", err: status.Error(codes.PermissionDenied, "disabled"), wantStatus: http.StatusForbidden, wantCode: `"code":10012`},
		{name: "unavailable", err: status.Error(codes.Unavailable, "redis detail"), wantStatus: http.StatusServiceUnavailable, wantCode: `"code":503`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := newTestAuthController(t, &fakeAuthRPC{login: func(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
				return nil, test.err
			}})
			recorder := performAuthRequest(controller.Login, "/api/v1/auth/login", `{"username":"user","password":"abc12345"}`, nil)
			if recorder.Code != test.wantStatus || !strings.Contains(recorder.Body.String(), test.wantCode) || strings.Contains(recorder.Body.String(), "redis detail") {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestAuthControllerMapsCommonRPCErrors(t *testing.T) {
	tests := []struct {
		name       string
		operation  authOperation
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "invalid argument", operation: authOperationRegister, err: status.Error(codes.InvalidArgument, "password detail"), wantStatus: http.StatusBadRequest, wantCode: `"code":10001`},
		{name: "username exists", operation: authOperationRegister, err: status.Error(codes.AlreadyExists, "username already exists"), wantStatus: http.StatusConflict, wantCode: `"code":10003`},
		{name: "email exists", operation: authOperationRegister, err: status.Error(codes.AlreadyExists, "email already exists"), wantStatus: http.StatusConflict, wantCode: `"code":10004`},
		{name: "email code", operation: authOperationGeneric, err: status.Error(codes.FailedPrecondition, "code detail"), wantStatus: http.StatusBadRequest, wantCode: `"code":10005`},
		{name: "mail rate", operation: authOperationGeneric, err: status.Error(codes.ResourceExhausted, "redis detail"), wantStatus: http.StatusTooManyRequests, wantCode: `"code":10010`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			engine.GET("/error", func(c *gin.Context) { writeAuthRPCError(c, test.operation, test.err) })
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/error", nil))
			if recorder.Code != test.wantStatus || !strings.Contains(recorder.Body.String(), test.wantCode) || strings.Contains(recorder.Body.String(), "detail") {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func newTestAuthController(t *testing.T, client authRPC) *AuthController {
	t.Helper()
	controller, err := newAuthController(client, time.Second, newTestRefreshCookieManager(t, false))
	if err != nil {
		t.Fatalf("newAuthController() error = %v", err)
	}
	return controller
}

func performAuthRequest(handler gin.HandlerFunc, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.POST(path, handler)
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func testLoginResponse() *authpb.LoginResponse {
	return &authpb.LoginResponse{User: testAuthUser(), Tokens: &authpb.TokenPair{
		AccessToken: "access", RefreshToken: "refresh", AccessExpiresIn: 1200, RefreshExpiresIn: 604800,
	}}
}

func testAuthUser() *authpb.UserInfo {
	return &authpb.UserInfo{
		Id: 7, Username: "user", Email: "user@qq.com", Phone: "secret-phone",
		Avatar: "avatar.png", Sex: "unknown", Age: 20, Role: "user", Status: "active", AuthVersion: 3,
	}
}

func assertResponseOmitsSecrets(t *testing.T, body string) {
	t.Helper()
	for _, forbidden := range []string{"refresh_token", "refresh", "password", "secret-phone", "auth_version"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("response leaked %q: %s", forbidden, body)
		}
	}
}

func assertRefreshCookie(t *testing.T, recorder *httptest.ResponseRecorder, value string) {
	t.Helper()
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies=%d headers=%v", len(cookies), recorder.Header())
	}
	cookie := cookies[0]
	if cookie.Name != refreshCookieName || cookie.Value != value || cookie.Path != refreshCookiePath ||
		!cookie.HttpOnly || cookie.Secure || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != 604800 || cookie.Domain != "" {
		t.Fatalf("refresh cookie = %+v", cookie)
	}
}

func assertClearedRefreshCookie(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != refreshCookieName || cookies[0].Value != "" || cookies[0].MaxAge != -1 ||
		cookies[0].Path != refreshCookiePath || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("cleared cookie = %+v", cookies)
	}
}
