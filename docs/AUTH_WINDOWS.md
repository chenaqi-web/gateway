# Windows 下的 Gateway 身份认证

Windows 开发与集成的权威指南位于：

[`core-server/docs/auth-windows.md`](../../core-server/docs/auth-windows.md)

该文档涵盖以下内容：

* 独立 MySQL 环境
* Docker Redis
* 环境变量
* 代码生成器
* 测试
* 服务启动顺序
* HTTP 请求示例
* Cookie 行为
* Refresh Token 轮换机制
* CORS 配置
* 手动创建初始管理员的流程

Gateway 相关配置位于 `gateway/conf/config.yaml`：

* `server.addr` 是 Gateway 的 HTTP 监听地址。应使用配置文件中的当前值，不要强制使用端口 `9000`。
* `rpc.core_server_addr` 是 Core 服务的 gRPC 地址。
* `auth.refresh_expire` 和 `auth.cookie_secure` 分别用于控制 Refresh Cookie 的有效期和 `Secure` 属性。
* `storage.base_url` 与 Gateway 的监听地址无关，不应为了身份认证功能修改该配置。

身份认证处理器只会通过 gRPC 调用 Core 服务。

它们不会访问 Gateway 中现有的以下内容：

* GORM
* MySQL
* Redis
* AI Chat Repository
* Wire Provider

请按照权威指南中的说明，在 Gateway 项目根目录下运行相关命令。
