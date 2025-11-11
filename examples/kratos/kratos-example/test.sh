#!/bin/bash

# Kratos + SA-Token 集成示例测试脚本

BASE_URL="http://localhost:8000"

echo "=========================================="
echo "Kratos + SA-Token 集成示例测试"
echo "=========================================="
echo ""

# 1. 测试公开接口
echo "1. 测试公开接口（无需登录）"
curl -s "$BASE_URL/api/public/info" | jq .
echo -e "\n"

# 2. 测试未登录访问受保护资源
echo "2. 测试未登录访问受保护资源（应该失败）"
curl -s "$BASE_URL/api/user/info" | jq .
echo -e "\n"

# 3. 管理员登录
echo "3. 管理员登录"
ADMIN_TOKEN=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')
echo "Admin Token: $ADMIN_TOKEN"
echo -e "\n"

# 4. 获取管理员用户信息
echo "4. 获取管理员用户信息（需要登录）"
curl -s "$BASE_URL/api/user/info" \
  -H "satoken: $ADMIN_TOKEN" | jq .
echo -e "\n"

# 5. 访问管理员接口
echo "5. 访问管理员接口（需要admin角色）"
curl -s "$BASE_URL/api/admin/dashboard" \
  -H "satoken: $ADMIN_TOKEN" | jq .
echo -e "\n"

# 6. 编辑用户
echo "6. 编辑用户（需要user.edit权限）"
curl -s -X POST "$BASE_URL/api/user/edit" \
  -H "Content-Type: application/json" \
  -H "satoken: $ADMIN_TOKEN" \
  -d '{"user_id":"1002","new_username":"newuser"}' | jq .
echo -e "\n"

# 7. 普通用户登录
echo "7. 普通用户登录"
USER_TOKEN=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}' | jq -r '.token')
echo "User Token: $USER_TOKEN"
echo -e "\n"

# 8. 普通用户获取信息
echo "8. 普通用户获取信息"
curl -s "$BASE_URL/api/user/info" \
  -H "satoken: $USER_TOKEN" | jq .
echo -e "\n"

# 9. 普通用户访问管理员接口（应该失败）
echo "9. 普通用户访问管理员接口（应该失败 - 缺少admin角色）"
curl -s "$BASE_URL/api/admin/dashboard" \
  -H "satoken: $USER_TOKEN" | jq .
echo -e "\n"

# 10. 普通用户编辑用户（应该失败）
echo "10. 普通用户编辑用户（应该失败 - 缺少user.edit权限）"
curl -s -X POST "$BASE_URL/api/user/edit" \
  -H "Content-Type: application/json" \
  -H "satoken: $USER_TOKEN" \
  -d '{"user_id":"1003","new_username":"test"}' | jq .
echo -e "\n"

# 11. Editor 登录
echo "11. Editor 登录"
EDITOR_TOKEN=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"editor","password":"editor123"}' | jq -r '.token')
echo "Editor Token: $EDITOR_TOKEN"
echo -e "\n"

# 12. Editor 编辑用户（应该成功 - 有user.edit权限）
echo "12. Editor 编辑用户（应该成功 - 有user.edit权限）"
curl -s -X POST "$BASE_URL/api/user/edit" \
  -H "Content-Type: application/json" \
  -H "satoken: $EDITOR_TOKEN" \
  -d '{"user_id":"1002","new_username":"edited"}' | jq .
echo -e "\n"

# 13. 管理员登出
echo "13. 管理员登出"
curl -s -X POST "$BASE_URL/api/logout" \
  -H "satoken: $ADMIN_TOKEN" | jq .
echo -e "\n"

# 14. 登出后访问（应该失败）
echo "14. 登出后访问受保护资源（应该失败）"
curl -s "$BASE_URL/api/user/info" \
  -H "satoken: $ADMIN_TOKEN" | jq .
echo -e "\n"

echo "=========================================="
echo "测试完成！"
echo "=========================================="
