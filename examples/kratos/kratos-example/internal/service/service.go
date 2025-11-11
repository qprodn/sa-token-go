package service

import (
	"github.com/click33/sa-token-go/core/config"
	"github.com/click33/sa-token-go/core/manager"
	"github.com/click33/sa-token-go/storage/memory"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService, NewUserService, NewSaTokenManager)

// NewSaTokenManager 创建 sa-token 管理器
func NewSaTokenManager() *manager.Manager {
	storage := memory.NewStorage()
	cfg := &config.Config{
		TokenName:     "satoken",
		Timeout:       2592000, // 30天
		IsReadCookie:  true,
		IsReadHeader:  true,
		ActiveTimeout: -1,
		IsConcurrent:  true,
		IsShare:       false,
		MaxLoginCount: -1,
		TokenStyle:    config.TokenStyleUUID,
	}

	return manager.NewManager(storage, cfg)
}
