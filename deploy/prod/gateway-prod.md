# Gateway 生产部署

本文说明如何使用 Docker 独立部署 Gateway。Gateway、Core Server、Agent Server 和 Redis 必须加入同一个 Docker 网络（示例：`chenaqi-net`）。

## 1. 准备生产配置

```bash
docker network create chenaqi-net 2>/dev/null || true
mkdir -p /home/docker/gateway/uploads
```

## 2. 构建镜像

```bash
cd /home/newweb/gateway
docker build \
  -f deploy/prod/Dockerfile \
  -t renai-gateway:prod \
  .
```

Dockerfile 默认使用 DaoCloud 镜像源、`goproxy.cn` 和 `conf/config.prod.yaml`，无需额外参数。

## 3. 启动 Gateway

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

## 4. 检查运行状态

```bash
docker ps --filter name=gateway
docker logs -f gateway
curl http://127.0.0.1:8079/health
```

如果需要公网访问，还要在云安全组和服务器防火墙放行 TCP `8079`，或通过 Nginx 反向代理并启用 HTTPS。

## 5. 更新部署

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
