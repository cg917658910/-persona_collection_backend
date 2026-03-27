# Vben login contract v1

为兼容 Vben 默认登录流，本版本新增了以下接口：

## 1. Login
`POST /api/v1/auth/login`

request:
```json
{
  "username": "admin",
  "password": "123456"
}
```

response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "accessToken": "jwt-token"
  }
}
```

## 2. User info
`GET /api/v1/user/info`

header:
`Authorization: Bearer <token>`

response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "userId": "admin-local",
    "username": "admin",
    "realName": "管理员",
    "roles": ["admin"]
  }
}
```

## 3. Permission codes
`GET /api/v1/auth/codes`

response:
```json
{
  "code": 0,
  "message": "ok",
  "data": []
}
```

## Notes
- 保留了旧的 `/api/v1/admin/auth/login` 与 `/api/v1/admin/auth/me`
- 新增的是一层兼容 Vben 默认协议的接口
- 当前适合先联调登录，不是最终权限系统
