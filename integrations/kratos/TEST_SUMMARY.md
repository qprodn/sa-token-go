# Kratos Integration Unit Tests

## 测试文件

为 `/integrations/kratos` 模块创建了以下单元测试文件：

### 1. checker_test.go (9.8KB)
测试所有检查器（Checker）的功能：

- **LoginChecker**: 登录检查器测试
  - 有登录ID的情况
  - 无登录ID的情况

- **PermissionChecker**: 权限检查器测试
  - 拥有权限
  - 缺少权限

- **PermissionsAndChecker**: 多权限AND检查器测试
  - 拥有所有权限
  - 缺少部分权限
  - 空权限列表

- **PermissionsOrChecker**: 多权限OR检查器测试
  - 拥有任一权限
  - 无任何权限
  - 空权限列表

- **RoleChecker**: 角色检查器测试
  - 拥有角色
  - 缺少角色

- **RolesAndChecker**: 多角色AND检查器测试
  - 拥有所有角色
  - 缺少部分角色
  - 空角色列表

- **RolesOrChecker**: 多角色OR检查器测试
  - 拥有任一角色
  - 无任何角色
  - 空角色列表

- **DisableChecker**: 账号封禁检查器测试
  - 未封禁状态
  - 已封禁状态

- **CustomChecker**: 自定义检查器测试
  - 检查通过
  - 检查失败

- **AndChecker**: AND组合检查器测试
  - 所有检查通过
  - 部分检查失败

- **OrChecker**: OR组合检查器测试
  - 任一检查通过
  - 所有检查失败
  - 空检查器列表

- **构造函数测试**: 所有检查器的构造函数

### 2. matcher_test.go (7.2KB)
测试所有操作匹配器（OperationMatcher）的功能：

- **ExactMatcher**: 精确匹配测试
  - 完全匹配
  - 不匹配
  - 部分匹配

- **PrefixMatcher**: 前缀匹配测试
  - 精确匹配
  - 前缀匹配
  - 不匹配

- **SuffixMatcher**: 后缀匹配测试
  - 精确匹配
  - 后缀匹配
  - 不匹配

- **WildcardMatcher**: 通配符匹配测试
  - 星号通配符 (*)
  - 问号通配符 (?)
  - 多个通配符
  - 通配符在末尾

- **RegexMatcher**: 正则表达式匹配测试
  - 数字模式
  - 单词模式

- **ContainsMatcher**: 包含匹配测试
  - 包含子串
  - 不包含子串

- **FuncMatcher**: 函数匹配测试
  - 自定义函数逻辑

- **AndMatcher**: AND组合匹配器测试
  - 所有匹配
  - 部分匹配

- **OrMatcher**: OR组合匹配器测试
  - 任一匹配
  - 全不匹配

- **NotMatcher**: NOT匹配器测试
  - 反转匹配结果

- **辅助函数测试**: And(), Or(), Not(), newPatternMatcher()

### 3. util_test.go (2.3KB)
测试工具函数：

- **indexOf**: 查找子串索引
  - 在开头找到
  - 在中间找到
  - 在末尾找到
  - 未找到
  - 空子串

- **lastIndexOf**: 查找子串最后出现的索引
  - 在末尾找到
  - 在中间找到（最后一次）
  - 在开头找到

- **trimSpace**: 去除空白字符
  - 无空格
  - 前导空格
  - 尾随空格
  - 两端空格
  - 制表符
  - 混合空白字符

- **contains**: 检查是否包含子串
  - 包含
  - 不包含
  - 空子串
  - 相同字符串

## 测试结果

```bash
PASS
coverage: 33.6% of statements
ok      github.com/click33/sa-token-go/integrations/kratos    0.005s
```

所有测试用例均通过，包括：
- 12个检查器测试函数
- 13个匹配器测试函数
- 4个工具函数测试函数

总计测试用例数：**约130+个子测试**

## 测试覆盖的功能模块

1. **认证检查（Checker）**
   - 登录验证
   - 权限验证（单个、多个AND、多个OR）
   - 角色验证（单个、多个AND、多个OR）
   - 账号封禁检查
   - 自定义检查
   - 组合检查器（AND、OR）

2. **路由匹配（Matcher）**
   - 精确匹配
   - 前缀/后缀匹配
   - 通配符匹配
   - 正则表达式匹配
   - 包含匹配
   - 函数匹配
   - 组合匹配器（AND、OR、NOT）

3. **工具函数（Util）**
   - 字符串查找
   - 空白字符处理
   - 子串包含检查

## 依赖的核心模块

测试使用以下依赖：
- `github.com/click33/sa-token-go/core/config` - 配置管理
- `github.com/click33/sa-token-go/core/manager` - 核心管理器
- `github.com/click33/sa-token-go/storage/memory` - 内存存储

## 注意事项

1. 测试使用内存存储，无需外部依赖
2. 所有测试用例均采用表驱动测试（Table-Driven Tests）模式
3. 测试覆盖了正常流程和异常流程
4. 使用子测试（subtests）提高测试的可读性和可维护性
