package middleware

import (
	"net/http"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// 每秒允许的请求数
	Rate rate.Limit
	// 令牌桶大小
	Burst int
	// 根据客户端 IP 生成 key 的函数
	KeyGenerator func(*gin.Context) string
	// 限流失败时的处理函数
	OnLimitReached func(*gin.Context)
}

// RateLimiter 限流器
type RateLimiter struct {
	limiter *rate.Limiter
	config  RateLimiterConfig
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.Rate == 0 {
		config.Rate = rate.Limit(10) // 默认每秒 10 个请求
	}
	if config.Burst == 0 {
		config.Burst = 20 // 默认令牌桶大小
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultKeyGenerator
	}
	if config.OnLimitReached == nil {
		config.OnLimitReached = defaultOnLimitReached
	}

	return &RateLimiter{
		limiter: rate.NewLimiter(config.Rate, config.Burst),
		config:  config,
	}
}

// defaultKeyGenerator 默认的 key 生成器（使用客户端 IP）
func defaultKeyGenerator(c *gin.Context) string {
	return c.ClientIP()
}

// defaultOnLimitReached 默认的限流处理函数
func defaultOnLimitReached(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"success":     false,
		"error":       "RATE_LIMIT_EXCEEDED",
		"message":     "请求过于频繁，请稍后再试",
		"retry_after": "1m",
	})
	c.Abort()
}

// Middleware 返回限流中间件函数
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	// 创建一个 map 来存储每个客户端的限流器
	clients := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		key := rl.config.KeyGenerator(c)

		// 为每个客户端创建独立的限流器
		if _, exists := clients[key]; !exists {
			clients[key] = rate.NewLimiter(rl.config.Rate, rl.config.Burst)
		}

		limiter := clients[key]

		if !limiter.Allow() {
			rl.config.OnLimitReached(c)
			return
		}

		c.Next()
	}
}

// SimpleRateLimitMiddleware 简单的全局限流中间件
func SimpleRateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// PerIPRateLimitMiddleware 每个独立 IP 的限流中间件
func PerIPRateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	config := RateLimiterConfig{
		Rate:  r,
		Burst: b,
		KeyGenerator: func(c *gin.Context) string {
			return c.ClientIP()
		},
		OnLimitReached: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":     false,
				"error":       "RATE_LIMIT_EXCEEDED",
				"message":     "请求过于频繁，请稍后再试",
				"retry_after": "1m",
			})
			c.Abort()
		},
	}

	limiter := NewRateLimiter(config)
	return limiter.Middleware()
}
