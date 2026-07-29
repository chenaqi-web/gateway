# 在 Docker 中运行 Prometheus 和 Grafana（最小化说明）

本文档仅记录使用 Docker 运行 Prometheus 和 Grafana 的最小示例（当前不挂载本地配置文件）。如需自定义 Prometheus 配置或持久化数据，可在以后通过挂载或制作镜像来实现。

## 1. 使用 Docker 直接运行（最小示例，不挂载）

启动 Prometheus（使用镜像内置的默认配置）：

```powershell
docker run -d --name prometheus -p 9090:9090 prom/prometheus:latest
```

启动 Grafana（默认管理员密码设置为 admin）：

```powershell
docker run -d --name grafana -p 3000:3000 -e GF_SECURITY_ADMIN_PASSWORD=admin grafana/grafana:latest
```

说明：以上示例不将主机文件挂载到容器中，Prometheus 和 Grafana 会使用镜像内的默认配置与存储（非持久化）。在生产或需要自定义 scrape 配置时，请改为挂载自定义 `prometheus.yml` 和数据卷。

## 2. 使用 docker-compose（示例，无挂载）

仓库内已包含一个最小的 `deploy/docker-compose-prometheus.yml`，可用于同时启动 Prometheus 与 Grafana（不挂载主机文件）：

```yaml
version: '3.7'
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    restart: unless-stopped
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    restart: unless-stopped

# 启动命令：
# docker-compose -f deploy/docker-compose-prometheus.yml up -d
```

## 3. 在本项目中抓取指标

本项目已在 `internal/facade/router.go` 注册了 Prometheus 的中间件，并提供 `/metrics` 路径。

默认网关服务地址为 `http://<host>:<port>/metrics`，例如开发时可能是 `http://localhost:8080/metrics`。

如果你运行的是上面的最小化 Docker 示例，Prometheus 将使用镜像内的默认配置，可能不会自动包含本项目的 `/metrics` 目标。要将本项目加入抓取列表，请在未来自定义 `prometheus.yml`，例如：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gateway'
    static_configs:
      - targets: ['host.docker.internal:8080']
```

注意：如果 Prometheus 容器在 Docker 中运行且服务在宿主机上，Windows/macOS 下可使用 `host.docker.internal:8080` 作为目标地址。



---

如果你希望我为项目生成一个默认的 `deploy/prometheus.yml`（Prometheus 抓取配置）并把它集成到 `docker-compose-prometheus.yml`（通过挂载或构建镜像），我可以继续为你创建。当前文档保持为只记录 Docker 启动示例且不挂载主机文件。
