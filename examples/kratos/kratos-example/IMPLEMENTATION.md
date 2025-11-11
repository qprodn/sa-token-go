# Kratos 集成 SA-Token 示例实现说明

## 📋 实现总结

本示例展示了如何在 Kratos 微服务框架中集成 SA-Token 进行身份认证和权限控制。

## 🎯 核心要点

### 1. SA-Token Manager 初始化

在 `internal/service/service.go` 中创建全局 Manager：

```go
func NewSaTokenManager() *manager.Manager {
    storage := memory.NewStorage()  // 使用内存存储
    cfg := &config.Config{
        TokenName:    "satoken",
        Timeout:      2592000,  // 30天过期
        IsReadCookie: true,
        IsReadHeader: true,
    }
    return manager.NewManager(storage, cfg)
}
```

### 2. 中间件配置

在 `internal/server/http.go` 中配置 SA-Token 中间件：

```go
// 创建插件
saPlugin := sakratos.NewPlugin(mgr)

// 配置路由规则 - 核心特性展示
saPlugin.
    Skip("/api/login", "/api/public/*", "/helloworld/*").  // 跳过公开路由
    For("/api/user/info").RequireLogin().Build().           // 需要登录
    For("/api/admin/*").RequireLogin().RequireRole("admin").Build(). // 需要角色
    For("/api/user/edit").RequireLogin().RequirePermission("user.edit").Build() // 需要权限

// 添加到中间件链
http.Middleware(
    recovery.Recovery(),
    saPlugin.Server(),  // SA-Token 认证中间件
)
```

### 3. 服务实现

#### 登录接口
```go
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
    // 验证用户（此处硬编码示例）
    userID := "1001"
    
    // 调用 SA-Token 登录
    token, err := s.manager.Login(userID)
    
    // 设置角色和权限
    s.manager.SetRoles(userID, []string{"admin"})
    s.manager.SetPermissions(userID, []string{"user.edit"})
    
    return &v1.LoginReply{Token: token}, nil
}
```

#### 获取用户信息
```go
func (s *UserService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoReply, error) {
    // 从请求中提取 token
    kratosCtx := sakratos.NewKratosContext(ctx)
    token := kratosCtx.GetHeader("satoken")
    
    // 获取登录用户ID
    loginID, err := s.manager.GetLoginID(token)
    
    // 获取角色和权限
    roles, _ := s.manager.GetRoles(loginID)
    permissions, _ := s.manager.GetPermissions(loginID)
    
    return &v1.GetUserInfoReply{UserId: loginID, Roles: roles, Permissions: permissions}, nil
}
```

## 🔑 SA-Token 核心功能演示

| 功能 | 代码位置 | 说明 |
|------|----------|------|
| **登录/登出** | `user.go::Login/Logout` | token 生成和销毁 |
| **角色验证** | `http.go::RequireRole("admin")` | 检查用户是否拥有 admin 角色 |
| **权限验证** | `http.go::RequirePermission("user.edit")` | 检查用户是否拥有编辑权限 |
| **路由保护** | `http.go::For().Build()` | 链式配置路由规则 |
| **灵活匹配** | `http.go::Skip()` | 支持通配符路径匹配 |

## 📦 测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | admin123 | admin, user | user.view, user.edit, user.delete, admin.dashboard |
| user | user123 | user | user.view |
| editor | editor123 | editor | user.view, user.edit |

## 🚀 快速测试

1. 启动服务：
```bash
go run cmd/kratos-example/main.go cmd/kratos-example/wire_gen.go
```

2. 运行测试脚本：
```bash
./test.sh
```

3. 手动测试：
```bash
# 登录
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 使用返回的 token 访问受保护资源
curl http://localhost:8000/api/user/info \
  -H "satoken: YOUR_TOKEN"
```

## 🔧 关键技术点

### 1. Kratos Context 适配器
`sakratos.NewKratosContext(ctx)` 将 Kratos 的 context 转换为 SA-Token 可识别的接口，支持：
- 从 Header 读取 token
- 从 Cookie 读取 token  
- 从 Query 参数读取 token

### 2. 链式 API 设计
```go
saPlugin.
    For("/api/admin/*").
    RequireLogin().
    RequireRole("admin").
    CheckNotDisabled().
    WithPriority(10).
    Build()
```

### 3. 多种匹配器
- `For(pattern)` - 自动选择匹配器
- `ForExact(op)` - 精确匹配
- `ForPrefix(pre)` - 前缀匹配
- `ForPattern(pat)` - 通配符（`*` 和 `?`）
- `ForRegex(regex)` - 正则表达式

## ⚠️ 注意事项

1. **简化实现**：本示例为演示目的，使用了硬编码的用户数据和内存存储
2. **生产环境**：实际应用应使用 Redis 存储和数据库查询
3. **密码安全**：生产环境应使用加密存储密码
4. **错误处理**：示例中简化了错误处理，实际应更完善

## 📚 扩展阅读

详细文档请参考 `README_SATOKEN.md`

---

✨ 这个示例完整展示了 SA-Token 在 Kratos 中的集成方式，重点突出了中间件配置和服务层使用。
