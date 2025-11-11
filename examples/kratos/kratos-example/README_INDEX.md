# 📚 Kratos + SA-Token 集成示例 - 文档索引

欢迎使用 Kratos + SA-Token 集成示例！本项目提供了完整的代码和文档。

## 📖 文档列表

### 1. [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - 项目总结 ⭐ **推荐首读**
- ✅ 完成内容清单
- 📊 核心技术展示
- 🎯 设计亮点
- 📁 文件结构说明
- 💡 学习要点

### 2. [IMPLEMENTATION.md](IMPLEMENTATION.md) - 实现说明 ⭐ **核心文档**
- 🔑 SA-Token Manager 初始化
- 🔧 中间件配置详解
- 📝 服务实现代码
- 🔑 核心功能演示表
- 📦 测试账号说明

### 3. [README_SATOKEN.md](README_SATOKEN.md) - 详细使用文档
- 🚀 快速开始指南
- 🔌 API 接口说明
- 📝 完整的 curl 示例
- 🎓 进阶特性教程
- 💡 注意事项

## 🚀 快速开始

### 1️⃣ 构建项目
```bash
go mod tidy
go build -o bin/server cmd/kratos-example/*.go
```

### 2️⃣ 运行服务
```bash
./bin/server
# 或
go run cmd/kratos-example/main.go cmd/kratos-example/wire_gen.go
```

### 3️⃣ 测试（3选1）

#### 方式1：自动化测试脚本
```bash
./test.sh
```

#### 方式2：手动测试
```bash
# 登录
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 使用返回的 token 访问
curl http://localhost:8000/api/user/info \
  -H "satoken: YOUR_TOKEN"
```

#### 方式3：查看详细文档
参考 [README_SATOKEN.md](README_SATOKEN.md) 中的完整测试用例

## 🎯 学习路线

```
1. PROJECT_SUMMARY.md (5分钟)
   └─ 了解项目整体结构和完成内容

2. IMPLEMENTATION.md (10分钟)  
   └─ 学习核心实现代码和关键技术点

3. 运行并测试 (10分钟)
   └─ 实际体验功能

4. README_SATOKEN.md (15分钟)
   └─ 深入了解所有特性和高级用法
```

## 📦 测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | admin123 | admin, user | user.view, user.edit, user.delete, admin.dashboard |
| user | user123 | user | user.view |
| editor | editor123 | editor | user.view, user.edit |

## 🎓 核心功能速查

| 功能 | 接口 | 要求 |
|------|------|------|
| 登录 | `POST /api/login` | 无 |
| 登出 | `POST /api/logout` | 需要登录 |
| 用户信息 | `GET /api/user/info` | 需要登录 |
| 管理面板 | `GET /api/admin/dashboard` | 需要 admin 角色 |
| 编辑用户 | `POST /api/user/edit` | 需要 user.edit 权限 |
| 公开信息 | `GET /api/public/info` | 无 |

## 🔑 核心代码位置

| 功能 | 文件 | 说明 |
|------|------|------|
| Manager 初始化 | `internal/service/service.go` | SA-Token 管理器创建 |
| 中间件配置 | `internal/server/http.go` | 路由规则和中间件 |
| 服务实现 | `internal/service/user.go` | 业务逻辑和 API 使用 |
| Proto 定义 | `api/helloworld/v1/user.proto` | 接口定义 |

## 💡 提示

- 🔥 想快速上手？直接看 [IMPLEMENTATION.md](IMPLEMENTATION.md)
- 📚 想了解细节？阅读 [README_SATOKEN.md](README_SATOKEN.md)
- 🎯 想全面了解？从 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) 开始

## 🔗 相关链接

- [SA-Token 官方文档](https://sa-token.cc/)
- [SA-Token Go GitHub](https://github.com/click33/sa-token-go)
- [Kratos 官方文档](https://go-kratos.dev/)

---

**Happy Coding! 🎉**
