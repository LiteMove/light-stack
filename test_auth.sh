#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080/api/v1"

echo "=== LightStack-Go 认证模块测试 ==="
echo

echo "1. 测试健康检查..."
curl -s "${BASE_URL}/../health" | echo "健康检查: $(cat)"
echo

echo "2. 测试Ping..."
curl -s "${BASE_URL}/ping" | echo "Ping: $(cat)"
echo

echo "3. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123",
    "nickname": "测试用户"
  }')
echo "注册响应: $REGISTER_RESPONSE"
echo

echo "4. 测试管理员登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')
echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"
echo

if [ -n "$TOKEN" ]; then
  echo "5. 测试获取用户信息..."
  curl -s "${BASE_URL}/user/profile" \
    -H "Authorization: Bearer $TOKEN" | echo "用户信息: $(cat)"
  echo

  echo "6. 测试更新用户信息..."
  curl -s -X PUT "${BASE_URL}/user/profile" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"nickname": "超级管理员"}' | echo "更新响应: $(cat)"
  echo

  echo "7. 测试管理员权限..."
  curl -s "${BASE_URL}/admin/users" \
    -H "Authorization: Bearer $TOKEN" | echo "管理员接口: $(cat)"
  echo

  echo "8. 测试刷新Token..."
  curl -s -X POST "${BASE_URL}/auth/refresh" \
    -H "Authorization: Bearer $TOKEN" | echo "刷新Token: $(cat)"
  echo

  echo "9. 测试登出..."
  curl -s -X POST "${BASE_URL}/auth/logout" \
    -H "Authorization: Bearer $TOKEN" | echo "登出响应: $(cat)"
  echo
else
  echo "登录失败，无法获取Token"
fi

echo "=== 测试完成 ==="