package prometheus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 请求计数器:统计累计请求量，计算QPS
	// 示例：http_requests_total{method="GET", path="/user/info", status="200"} 15678
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// 请求耗时直方图
	// 记录每次请求耗时，自动归入对应的桶
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// 注册自定义指标
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	// 注册系统级指标（自动采集）
	prometheus.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: nil},
		),
	))
	prometheus.MustRegister(collectors.NewProcessCollector(
		collectors.ProcessCollectorOpts{},
	))
}

// GinMiddleware 返回 Gin 中间件，用于采集 Prometheus 指标
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// 获取路径
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		method := c.Request.Method
		statusStr := fmt.Sprintf("%d", status)

		// 记录指标
		httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// Handler 返回 /metrics 端点
// http://localhost:8080/metrics
func Handler() http.Handler {
	return promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// 可选：启用压缩
			EnableOpenMetrics: true,

			// 可选：自定义注册表
			// Registry: prometheus.DefaultRegisterer,
		},
	)
}
