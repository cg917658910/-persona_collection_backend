# pm-backend scaffold v7

Go + Gin backend scaffold for 拾荒者·人物集.

This version adds **Vben compatible login endpoints**.

## Run

```bash
cp .env.example .env
go mod tidy
go run ./cmd/server
```

## Vben compatible auth APIs

### Login
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

### User info
`GET /api/v1/user/info`

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

### Permission codes
`GET /api/v1/auth/codes`

response:
```json
{
  "code": 0,
  "message": "ok",
  "data": []
}
```

## Original admin auth APIs kept
- `POST /api/v1/admin/auth/login`
- `GET /api/v1/admin/auth/me`

## Notes
- This is for fast Vben login integration.
- Current admin CRUD is still not globally protected by JWT middleware.
- Next step can add auth middleware to `/api/v1/admin/*`.
