package kratos

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/transport"
	"net/http"

	"github.com/click33/sa-token-go/core/adapter"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type KratosContext struct {
	ctx    context.Context
	values map[string]interface{}
	mu     sync.RWMutex
}

// NewKratosContext creates a Kratos context adapter | 创建Kratos上下文适配器
// This constructor accepts any request/response objects that implement the KratosRequest/KratosResponse interfaces
func NewKratosContext(ctx context.Context) adapter.RequestContext {
	return &KratosContext{
		ctx: ctx,
	}
}

// GetHeader gets request header | 获取请求头
func (k *KratosContext) GetHeader(key string) string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		return tr.RequestHeader().Get(key)
	}
	return ""
}

// GetQuery gets query parameter | 获取查询参数
func (k *KratosContext) GetQuery(key string) string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok {
			request := htr.Request()
			return request.URL.Query().Get(key)
		}
	}
	return ""
}

// GetCookie gets cookie | 获取Cookie
func (k *KratosContext) GetCookie(key string) string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok {
			request := htr.Request()
			cookie, err := request.Cookie(key)
			if err != nil {
				return ""
			}
			return cookie.Value
		}
	}
	return ""
}

// SetHeader sets response header | 设置响应头
func (k *KratosContext) SetHeader(key, value string) {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		tr.ReplyHeader().Set(key, value)
	}
}

// SetCookie sets cookie | 设置Cookie
func (k *KratosContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	}
	khttp.SetCookie(k.ctx, cookie)
}

// GetClientIP gets client IP address | 获取客户端IP地址
func (k *KratosContext) GetClientIP() string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		// 尝试从X-Forwarded-For获取
		if xff := tr.RequestHeader().Get("X-Forwarded-For"); xff != "" {
			// X-Forwarded-For可能包含多个IP，取第一个
			if idx := indexOf(xff, ","); idx > 0 {
				return trimSpace(xff[:idx])
			}
			return trimSpace(xff)
		}

		// 尝试从X-Real-IP获取
		if xri := tr.RequestHeader().Get("X-Real-IP"); xri != "" {
			return trimSpace(xri)
		}

		// 如果是HTTP transport，尝试从Request获取
		if htr, ok := tr.(*khttp.Transport); ok {
			request := htr.Request()
			if request.RemoteAddr != "" {
				// RemoteAddr格式: "IP:Port"，需要去掉端口
				if idx := lastIndexOf(request.RemoteAddr, ":"); idx > 0 {
					return request.RemoteAddr[:idx]
				}
				return request.RemoteAddr
			}
		}
	}
	return ""
}

// GetMethod gets request method | 获取请求方法
func (k *KratosContext) GetMethod() string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok {
			request := htr.Request()
			return request.Method
		}
	}
	return ""
}

// GetPath gets request path | 获取请求路径
func (k *KratosContext) GetPath() string {
	if tr, ok := transport.FromServerContext(k.ctx); ok {
		return tr.Operation()
	}
	return ""
}

// Set sets context value | 设置上下文值
func (k *KratosContext) Set(key string, value interface{}) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.values == nil {
		k.values = make(map[string]interface{})
	}

	k.values[key] = value
}

// Get gets context value | 获取上下文值
func (k *KratosContext) Get(key string) (interface{}, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	value, exists := k.values[key]
	return value, exists
}
