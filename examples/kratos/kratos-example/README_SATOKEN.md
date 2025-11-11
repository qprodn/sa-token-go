# Kratos + SA-Token 集成示例

本示例展示如何在 Kratos 框架中集成 SA-Token 进行身份认证和权限控制。

## 功能特性

本示例演示了 SA-Token 的核心功能：

- ✅ **用户登录/登出** - 基于 token 的会话管理
- ✅ **角色验证** - 检查用户是否拥有特定角色
- ✅ **权限验证** - 检查用户是否拥有特定权限
- ✅ **路由保护** - 使用链式 API 配置路由规则
- ✅ **灵活的匹配器** - 支持精确匹配、前缀、后缀、通配符等多种匹配方式

## 目录结构

```
.
├── api/helloworld/v1/          # Proto 定义和生成的代码
│   ├── greeter.proto           # Greeter 服务定义
│   └── user.proto              # User 服务定义（SA-Token 示例）
├── internal/
│   ├── server/
│   │   └── http.go             # HTTP 服务器配置（包含 SA-Token 中间件）
│   └── service/
│       ├── user.go             # 用户服务实现（展示 SA-Token 用法）
│       └── service.go          # 服务提供者（包含 SA-Token Manager）
└── cmd/kratos-example/
    └── main.go                 # 应用入口
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 生成代码（如果修改了 proto）

```bash
make api
```

### 3. 运行服务

```bash
go run cmd/kratos-example/main.go
```

服务将在 `http://0.0.0.0:8000` 启动。

## API 接口

### 公开接口（无需登录）

#### 1. 公开信息
```bash
curl http://localhost:8000/api/public/info
```

#### 2. 用户登录
```bash
# 管理员登录
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 普通用户登录
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}'

# 编辑者登录
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"editor","password":"editor123"}'
```

响应示例：
```json
{
  "token": "abcd1234...",
  "message": "登录成功",
  "user_id": "1001"
}
```

### 需要登录的接口

#### 3. 获取用户信息
```bash
# 使用 Header 传递 token
curl http://localhost:8000/api/user/info \
  -H "satoken: YOUR_TOKEN"

# 或使用 Cookie
curl http://localhost:8000/api/user/info \
  -b "satoken=YOUR_TOKEN"
```

响应示例：
```json
{
  "user_id": "1001",
  "username": "admin",
  "roles": ["admin", "user"],
  "permissions": ["user.view", "user.edit", "user.delete", "admin.dashboard"]
}
```

#### 4. 用户登出
```bash
curl -X POST http://localhost:8000/api/logout \
  -H "satoken: YOUR_TOKEN"
```

### 需要特定角色的接口

#### 5. 管理员接口（需要 admin 角色）
```bash
curl http://localhost:8000/api/admin/dashboard \
  -H "satoken: YOUR_TOKEN"
```

只有 `admin` 用户能访问，其他用户会返回 403 错误。

### 需要特定权限的接口

#### 6. 编辑用户（需要 user.edit 权限）
```bash
curl -X POST http://localhost:8000/api/user/edit \
  -H "Content-Type: application/json" \
  -H "satoken: YOUR_TOKEN" \
  -d '{"user_id":"1002","new_username":"newname"}'
```

只有拥有 `user.edit` 权限的用户能访问（admin 和 editor）。

## 测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | admin123 | admin, user | user.view, user.edit, user.delete, admin.dashboard |
| user | user123 | user | user.view |
| editor | editor123 | editor | user.view, user.edit |

## 核心代码解析

### 1. 创建 SA-Token Manager

在 `internal/service/service.go` 中：

```go
func NewSaTokenManager() *core.Manager {
    storage := memory.NewStorage()
    config := core.NewConfig()
    config.TokenName = "satoken"
    config.Timeout = 2592000 // 30天
    config.IsReadCookie = true
    config.IsWriteHead = true
    
    return core.NewManager(storage, config)
}
```

### 2. 配置 SA-Token 中间件

在 `internal/server/http.go` 中：

```go
// 创建 sa-token 中间件
saPlugin := sakratos.NewPlugin(manager)

// 配置路由规则
saPlugin.
    // 跳过公开路由
    Skip("/api/login", "/api/public/*", "/helloworld/*").
    // 用户信息需要登录
    For("/api/user/info").RequireLogin().Build().
    // 管理员接口需要admin角色
    For("/api/admin/*").RequireLogin().RequireRole("admin").Build().
    // 编辑用户需要user.edit权限
    For("/api/user/edit").RequireLogin().RequirePermission("user.edit").Build().
    // 登出需要登录
    For("/api/logout").RequireLogin().Build()

// 添加到 Kratos 中间件链
var opts = []http.ServerOption{
    http.Middleware(
        recovery.Recovery(),
        saPlugin.Server(), // SA-Token 中间件
    ),
}
```

### 3. 在服务中使用 SA-Token

#### 登录

```go
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
    // 创建 Kratos 上下文适配器
    kratosCtx := sakratos.NewKratosContext(ctx)
    saCtx := core.NewContext(kratosCtx, s.manager)

    // 登录用户
    token, err := saCtx.Login(userID)
    if err != nil {
        return nil, err
    }

    // 设置角色和权限
    s.manager.SetRole(userID, "admin")
    s.manager.SetPermission(userID, "user.edit")

    return &v1.LoginReply{Token: token}, nil
}
```

#### 获取登录用户信息

```go
func (s *UserService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoReply, error) {
    // 从 context 获取 sa-token 上下文（由中间件注入）
    saCtx, ok := ctx.Value("satoken").(*core.Context)
    if !ok {
        return nil, fmt.Errorf("未登录")
    }

    // 获取登录ID
    loginID, err := saCtx.GetLoginID()
    if err != nil {
        return nil, err
    }

    // 获取角色和权限
    roles := s.manager.GetRoleList(loginID)
    permissions := s.manager.GetPermissionList(loginID)

    return &v1.GetUserInfoReply{
        UserId:      loginID,
        Roles:       roles,
        Permissions: permissions,
    }, nil
}
```

## SA-Token 路由规则说明

### 匹配器类型

- `For(pattern)` - 自动选择匹配器（支持通配符 `*`）
- `ForExact(operation)` - 精确匹配
- `ForPrefix(prefix)` - 前缀匹配
- `ForSuffix(suffix)` - 后缀匹配
- `ForPattern(pattern)` - 通配符匹配（支持 `*` 和 `?`）
- `ForRegex(regex)` - 正则表达式匹配
- `ForContains(substring)` - 包含匹配
- `ForFunc(fn)` - 自定义函数匹配

### 检查器类型

- `RequireLogin()` - 需要登录
- `RequireRole(role)` - 需要指定角色
- `RequireRoles(roles...)` - 需要多个角色（AND 逻辑）
- `RequireAnyRole(roles...)` - 需要任一角色（OR 逻辑）
- `RequirePermission(perm)` - 需要指定权限
- `RequirePermissions(perms...)` - 需要多个权限（AND 逻辑）
- `RequireAnyPermission(perms...)` - 需要任一权限（OR 逻辑）
- `CheckNotDisabled()` - 检查账号未被封禁
- `CustomCheck(name, fn)` - 自定义检查逻辑

### 示例：复杂规则配置

```go
saPlugin.
    // 管理员区域：需要 admin 角色且未被封禁
    For("/api/admin/*").
        RequireLogin().
        RequireRole("admin").
        CheckNotDisabled().
        Build().
    
    // 内容管理：需要任一内容相关角色
    For("/api/content/*").
        RequireLogin().
        RequireAnyRole("admin", "editor", "author").
        Build().
    
    // 删除操作：需要删除权限和管理角色
    For("/api/*/delete").
        RequireLogin().
        RequirePermission("delete").
        RequireRole("admin").
        Build().
    
    // 自定义检查：工作时间限制
    For("/api/sensitive/*").
        RequireLogin().
        CustomCheck("working-hours", func(ctx context.Context, manager *core.Manager, loginID string) error {
            hour := time.Now().Hour()
            if hour < 9 || hour > 18 {
                return errors.New("只能在工作时间访问")
            }
            return nil
        }).
        Build()
```

## 进阶特性

### 1. 设置规则优先级

```go
saPlugin.
    For("/api/admin/public").RequireLogin().WithPriority(10).Build(). // 高优先级
    For("/api/admin/*").RequireRole("admin").WithPriority(5).Build()   // 低优先级
```

### 2. 自定义错误处理

```go
plugin := sakratos.NewPlugin(manager)
plugin.SetErrorHandler(func(ctx context.Context, err error) error {
    // 自定义错误响应
    return errors.New("自定义错误: " + err.Error())
})
```

### 3. Token 来源

SA-Token 会自动从以下位置获取 token（按优先级）：

1. HTTP Header: `satoken: YOUR_TOKEN`
2. Cookie: `satoken=YOUR_TOKEN`
3. Query 参数: `?satoken=YOUR_TOKEN`

## 注意事项

1. **本示例使用内存存储**，重启后数据丢失。生产环境建议使用 Redis 存储。
2. **密码硬编码**仅用于演示，生产环境应使用加密存储。
3. **角色和权限**在登录时设置，实际应从数据库加载。

## 扩展阅读

- [SA-Token 官方文档](https://sa-token.cc/)
- [Kratos 官方文档](https://go-kratos.dev/)
- [SA-Token Go 版本](https://github.com/click33/sa-token-go)

## License

MIT
