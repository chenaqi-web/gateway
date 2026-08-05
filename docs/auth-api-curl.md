# Gateway Auth cURL

服务地址：`http://127.0.0.1:8079`

## 发送邮箱验证码

```bash
curl --request POST 'http://127.0.0.1:8079/api/v1/auth/send-email-code' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "email": "user@example.com",
    "purpose": "register"
  }'
```

`purpose` 可填写 `register` 或 `reset_password`。

## 登录

```bash
curl --include \
  --request POST 'http://127.0.0.1:8079/api/v1/auth/login' \
  --header 'Content-Type: application/json' \
  --cookie-jar 'auth-cookie.txt' \
  --data-raw '{
    "username": "admin",
    "password": "123456"
  }'
```

从响应体复制 `data.access_token`，后续替换命令中的 `ACCESS_TOKEN`。

## 验证鉴权

```bash
curl --include \
  --request GET 'http://127.0.0.1:8079/api/v1/types/list' \
  --header 'Authorization: Bearer ACCESS_TOKEN' \
  --cookie 'auth-cookie.txt'
```

不传 Token 或修改 Token 后请求，应返回 `401`。

## 验证自动刷新

等待 access token 过期后，继续执行鉴权请求：

```bash
curl --include \
  --request GET 'http://127.0.0.1:8079/api/v1/types/list' \
  --header 'Authorization: Bearer EXPIRED_ACCESS_TOKEN' \
  --cookie 'auth-cookie.txt'
```

refresh token 有效时，响应头会返回：

```text
Authorization: Bearer NEW_ACCESS_TOKEN
```

## 注销

```bash
curl --include \
  --request POST 'http://127.0.0.1:8079/api/v1/auth/logout' \
  --header 'Authorization: Bearer ACCESS_TOKEN' \
  --cookie 'auth-cookie.txt'
```

## 验证黑名单

注销后再次使用原 access token：

```bash
curl --include \
  --request GET 'http://127.0.0.1:8079/api/v1/types/list' \
  --header 'Authorization: Bearer ACCESS_TOKEN' \
  --cookie 'auth-cookie.txt'
```

应返回 `401`。
