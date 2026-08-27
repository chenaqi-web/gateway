# Gateway 生产部署

本文说明如何使用 Docker 独立部署 Gateway。Gateway、Core Server、Agent Server 和 Redis 必须加入同一个 Docker 网络（示例：`chenaqi-net`）。

## 1. 准备运行环境

```bash
docker network inspect chenaqi-net >/dev/null 2>&1 || docker network create --driver bridge chenaqi-net
mkdir -p /home/docker/gateway/uploads
```

首次使用上传目录挂载前，如已有头像、封面或正文图片，先将旧上传文件复制到 `/home/docker/gateway/uploads`。该目录会覆盖容器内的 `/app/static/upload`；未迁移的历史文件将无法通过 `/static/upload/...` 访问。

```bash
cp -a /home/newweb/gateway/static/upload/. /home/docker/gateway/uploads/
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
docker run -d \
  --name gateway \
  --network chenaqi-net \
  -p 8079:8079 \
  -v /home/docker/gateway/uploads:/app/static/upload \
  --restart unless-stopped \
  renai-gateway:prod
```

镜像构建时已将 `conf/config.prod.yaml` 写入容器的 `/app/conf/config.yaml`。挂载上传目录可避免更新容器时丢失用户文件；修改生产配置后需重新构建镜像并重新创建容器。上传后的访问地址格式为 `storage.base_url/static/upload/...`。

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
  -v /home/docker/gateway/uploads:/app/static/upload \
  --restart unless-stopped \
  renai-gateway:prod
```
