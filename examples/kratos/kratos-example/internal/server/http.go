package server

import (
	v1 "kratos-example/api/helloworld/v1"
	"kratos-example/internal/conf"
	"kratos-example/internal/service"

	"github.com/click33/sa-token-go/core/manager"
	sakratos "github.com/click33/sa-token-go/integrations/kratos"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UserService, mgr *manager.Manager, logger log.Logger) *http.Server {
	// 创建 sa-token 中间件
	saPlugin := sakratos.NewPlugin(mgr)

	// 配置路由规则
	saPlugin.
		// 跳过公开路由
		Skip(v1.OperationUserLogin).
		// 用户信息需要登录
		For(v1.OperationUserGetUserInfo).RequireLogin().Build().
		// 管理员接口需要admin角色
		For(v1.OperationUserAdminOnly).RequireLogin().RequireRole("admin").Build().
		// 编辑用户需要user.edit权限
		For(v1.OperationUserEditUser).RequireLogin().RequirePermission("user.edit").Build().
		// 登出需要登录
		For(v1.OperationUserLogout).RequireLogin().Build()

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			saPlugin.Server(), // 添加 sa-token 中间件
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUserHTTPServer(srv, user)
	return srv
}
