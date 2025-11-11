package data

import (
	"context"
	"errors"

	"kratos-example/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// userRepo 用户仓储实现
type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo 创建用户仓储
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// FindByUsername 根据用户名查找用户
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	// 模拟数据库查询，实际应该从数据库获取
	users := r.getMockUsers()

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("用户不存在")
}

// FindByID 根据用户ID查找用户
func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	// 模拟数据库查询
	users := r.getMockUsers()

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors.New("用户不存在")
}

// getMockUsers 获取模拟用户数据
func (r *userRepo) getMockUsers() []*biz.User {
	return []*biz.User{
		{
			ID:       "1001",
			Username: "admin",
			Password: "admin123",
			Roles:    []string{"admin", "user"},
			Permissions: []string{
				"user.view",
				"user.edit",
				"user.delete",
				"admin.dashboard",
			},
		},
		{
			ID:       "1002",
			Username: "user",
			Password: "user123",
			Roles:    []string{"user"},
			Permissions: []string{
				"user.view",
			},
		},
		{
			ID:       "1003",
			Username: "editor",
			Password: "editor123",
			Roles:    []string{"editor"},
			Permissions: []string{
				"user.view",
				"user.edit",
			},
		},
	}
}
