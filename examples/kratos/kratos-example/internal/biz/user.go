package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// User 用户业务实体
type User struct {
	ID          string
	Username    string
	Password    string
	Roles       []string
	Permissions []string
}

// UserRepo 用户仓储接口
type UserRepo interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

// UserUsecase 用户业务逻辑
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase 创建用户业务逻辑
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 简单密码验证（实际应该使用加密）
	if user.Password != password {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// GetUserInfo 获取用户信息
func (uc *UserUsecase) GetUserInfo(ctx context.Context, userID string) (*User, error) {
	return uc.repo.FindByID(ctx, userID)
}

// EditUser 编辑用户
func (uc *UserUsecase) EditUser(ctx context.Context, userID, newUsername string) error {
	// 这里只是简单示例，实际应该调用repo更新
	uc.log.Infof("编辑用户: %s -> %s", userID, newUsername)
	return nil
}
