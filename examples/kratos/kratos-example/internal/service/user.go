package service

import (
	"context"
	"fmt"
	"github.com/click33/sa-token-go/stputil"

	v1 "kratos-example/api/helloworld/v1"

	"github.com/click33/sa-token-go/core/manager"
	sakratos "github.com/click33/sa-token-go/integrations/kratos"
)

// UserService 用户服务
type UserService struct {
	v1.UnimplementedUserServer
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{}
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	// 简化处理：直接硬编码用户验证
	var userID string
	var roles []string
	var permissions []string

	switch req.Username {
	case "admin":
		if req.Password != "admin123" {
			return nil, fmt.Errorf("密码错误")
		}
		userID = "1001"
		roles = []string{"admin", "user"}
		permissions = []string{"user.view", "user.edit", "user.delete", "admin.dashboard"}
	case "user":
		if req.Password != "user123" {
			return nil, fmt.Errorf("密码错误")
		}
		userID = "1002"
		roles = []string{"user"}
		permissions = []string{"user.view"}
	case "editor":
		if req.Password != "editor123" {
			return nil, fmt.Errorf("密码错误")
		}
		userID = "1003"
		roles = []string{"editor"}
		permissions = []string{"user.view", "user.edit"}
	default:
		return nil, fmt.Errorf("用户不存在")
	}

	// 登录用户
	token, err := stputil.Login(userID)
	if err != nil {
		return nil, err
	}

	// 设置角色和权限
	if err := stputil.SetRoles(userID, roles); err != nil {
		return nil, err
	}
	if err := stputil.SetPermissions(userID, permissions); err != nil {
		return nil, err
	}

	return &v1.LoginReply{
		Token:   token,
		Message: "登录成功",
		UserId:  userID,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutReply, error) {
	// 从 Kratos context 获取 token
	kratosCtx := sakratos.NewKratosContext(ctx)
	token := kratosCtx.GetHeader(stputil.GetConfig().TokenName)
	if token == "" {
		token = kratosCtx.GetCookie(s.manager.GetConfig().TokenName)
	}

	if token == "" {
		return nil, fmt.Errorf("未登录")
	}

	// 登出
	if err := s.manager.Logout(token); err != nil {
		return nil, err
	}

	return &v1.LogoutReply{
		Message: "登出成功",
	}, nil
}

// GetUserInfo 获取用户信息（需要登录）
func (s *UserService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoReply, error) {
	// 从 Kratos context 获取 token
	kratosCtx := sakratos.NewKratosContext(ctx)
	token := kratosCtx.GetHeader(s.manager.GetConfig().TokenName)
	if token == "" {
		token = kratosCtx.GetCookie(s.manager.GetConfig().TokenName)
	}

	// 获取登录ID
	loginID, err := s.manager.GetLoginID(token)
	if err != nil {
		return nil, fmt.Errorf("未登录")
	}

	// 获取角色和权限
	roles, _ := s.manager.GetRoles(loginID)
	permissions, _ := s.manager.GetPermissions(loginID)

	// 硬编码用户名映射
	usernameMap := map[string]string{
		"1001": "admin",
		"1002": "user",
		"1003": "editor",
	}

	return &v1.GetUserInfoReply{
		UserId:      loginID,
		Username:    usernameMap[loginID],
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// AdminOnly 管理员专用接口（需要admin角色）
func (s *UserService) AdminOnly(ctx context.Context, req *v1.AdminRequest) (*v1.AdminReply, error) {
	// 从 Kratos context 获取 token
	kratosCtx := sakratos.NewKratosContext(ctx)
	token := kratosCtx.GetHeader(s.manager.GetConfig().TokenName)
	if token == "" {
		token = kratosCtx.GetCookie(s.manager.GetConfig().TokenName)
	}

	// 获取登录ID
	loginID, err := s.manager.GetLoginID(token)
	if err != nil {
		return nil, fmt.Errorf("未登录")
	}

	return &v1.AdminReply{
		Message:   "欢迎管理员",
		AdminInfo: fmt.Sprintf("管理员ID: %s, 系统运行正常", loginID),
	}, nil
}

// EditUser 编辑用户（需要user.edit权限）
func (s *UserService) EditUser(ctx context.Context, req *v1.EditUserRequest) (*v1.EditUserReply, error) {
	return &v1.EditUserReply{
		Message: fmt.Sprintf("用户 %s 已更新为 %s", req.UserId, req.NewUsername),
	}, nil
}

// PublicInfo 公开信息（无需登录）
func (s *UserService) PublicInfo(ctx context.Context, req *v1.PublicInfoRequest) (*v1.PublicInfoReply, error) {
	return &v1.PublicInfoReply{
		Message: "这是公开信息，无需登录",
		Version: "1.0.0",
	}, nil
}
