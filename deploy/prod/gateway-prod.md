# Gateway 生产部署

本文说明如何使用 Docker 独立部署 Gateway。Gateway、Core Server、Agent Server 和 Redis 必须加入同一个 Docker 网络（示例：`chenaqi-net`）。

## 1. 准备生产配置

复制生产配置到服务器，并修改其中的密钥、数据库、Redis、LLM 和邮箱配置：

```bash
mkdir -p /home/newweb/gateway-config
cp conf/config.prod.yaml /home/newweb/gateway-config/config.yaml
chmod 600 /home/newweb/gateway-config/config.yaml
vi /home/newweb/gateway-config/config.yaml
```

生产配置中的服务地址应使用 Docker 网络中的服务名：

```yaml
server:
  addr: "0.0.0.0:8079"
rpc:
  core_server_addr: "core-server:8081"
http:
  agent_server_addr: "http://agent-server:8080"
Redis:
  Host: "redis"
```

请务必替换所有 `CHANGE_ME_*`、`YOUR_EMAIL@example.com` 和 `YOUR_DOMAIN_OR_PUBLIC_IP` 占位符。不要把真实密钥提交到代码仓库。

## 2. 创建网络和持久化目录

```bash
docker network create chenaqi-net 2>/dev/null || true
mkdir -p /home/docker/gateway/uploads
```

## 3. 构建镜像

```bash
cd /home/newweb/gateway
docker build \
  -f deploy/prod/Dockerfile \
  -t renai-gateway:prod \
  .
```

Dockerfile 默认使用 DaoCloud 镜像源和 `goproxy.cn`。如需更换镜像源，可传入：

```bash
docker build \
  --build-arg GO_IMAGE=golang:1.26-alpine \
  --build-arg ALPINE_IMAGE=alpine:3.21 \
  -f deploy/prod/Dockerfile \
  -t renai-gateway:prod \
  .
```

## 4. 启动 Gateway

```bash
docker rm -f gateway 2>/dev/null || true
docker run -d \
  --name gateway \
  --network chenaqi-net \
  -p 8079:8079 \
  -v /home/newweb/gateway-config/config.yaml:/app/conf/config.yaml:ro \
  -v /home/docker/gateway/uploads:/app/static/upload \
  --restart unless-stopped \
  renai-gateway:prod
```

挂载配置文件后，容器使用服务器上的生产配置；挂载上传目录可避免更新容器时丢失用户文件。

## 5. 检查运行状态

```bash
docker ps --filter name=gateway？
docker logs -f gateway
curl http://127.0.0.1:8079/health
```

如果需要公网访问，还要在云安全组和服务器防火墙放行 TCP `8079`，或通过 Nginx 反向代理并启用 HTTPS。

## 6. 更新部署

```bash
cd /home/newweb/gateway
docker build -f deploy/prod/Dockerfile -t renai-gateway:prod .
docker rm -f gateway
docker run -d \
  --name gateway \
  --network chenaqi-net \
  -p 8079:8079 \
  -v /home/newweb/gateway-config/config.yaml:/app/conf/config.yaml:ro \
  -v /home/docker/gateway/uploads:/app/static/upload \
  --restart unless-stopped \
  renai-gateway:prod
```
