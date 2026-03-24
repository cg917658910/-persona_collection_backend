# pm-backend scaffold v4

Go + Gin backend scaffold for 拾荒者·人物集, with PostgreSQL repo skeleton and static URL normalization.

## Run

```bash
cp .env.example .env
go mod tidy
go run ./cmd/server
```

## Modes

### Mock mode
```env
USE_MOCK=true
```

### PostgreSQL mode
```env
USE_MOCK=false
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pmdb?sslmode=disable
```

## Static assets and media URLs

Environment:

```env
PUBLIC_BASE_URL=http://localhost:8080
STATIC_MOUNT_PREFIX=/static
STATIC_LOCAL_DIR=./public
```

Behavior:
- if DB/mock value is already absolute (`http://` / `https://`) -> keep as-is
- if value starts with `/assets/...` -> convert to:
  - `http://localhost:8080/static/assets/...`

## PostgreSQL implemented now
- Home
- Discover random
- ListCharacters
- GetCharacterDetail
- ListWorks
- GetWorkDetail
- ListCreators
- GetCreatorDetail
- ListThemes
- GetThemeDetail
- ListSongs

## Recommended next step

1. Add admin CRUD routes under `/api/v1/admin`
2. Add authentication for admin
3. Add resource upload and media binding
