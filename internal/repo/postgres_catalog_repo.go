package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresCatalogRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresCatalogRepo(pool *pgxpool.Pool) CatalogRepo {
	return &postgresCatalogRepo{pool: pool}
}

func (r *postgresCatalogRepo) Home() (dto.HomePayload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	featured, err := r.queryHomeFeaturedCharacter(ctx)
	if err != nil {
		return dto.HomePayload{}, err
	}
	latestChars, err := r.queryHomeLatestCharacters(ctx, 8)
	if err != nil {
		return dto.HomePayload{}, err
	}
	featuredSongs, err := r.queryHomeFeaturedSongs(ctx, 8)
	if err != nil {
		return dto.HomePayload{}, err
	}
	recommendedWorks, err := r.queryHomeRecommendedWorks(ctx, 8)
	if err != nil {
		return dto.HomePayload{}, err
	}
	homeThemes, err := r.queryHomeThemes(ctx, 8)
	if err != nil {
		return dto.HomePayload{}, err
	}

	return dto.HomePayload{
		FeaturedCharacter: featured,
		LatestCharacters:  latestChars,
		FeaturedSongs:     featuredSongs,
		RecommendedWorks:  recommendedWorks,
		Themes:            homeThemes,
	}, nil
}

func (r *postgresCatalogRepo) RandomCharacter(theme string, exclude []string) (dto.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item dto.Character
	err := r.pool.QueryRow(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.is_active = TRUE
  AND c.status = 'published'
  AND ($1 = '' OR EXISTS (
    SELECT 1
    FROM public.pm_character_themes x
    JOIN public.pm_themes t ON t.id = x.theme_id
    WHERE x.character_id = c.id
      AND t.slug = $1
      AND t.is_active = TRUE
  ))
  AND (COALESCE(array_length($2::text[], 1), 0) = 0 OR NOT c.slug = ANY($2))
ORDER BY
  random() * CASE
    WHEN COALESCE(c.meta->>'discover_weight', '') ~ '^[0-9]+(\.[0-9]+)?$'
      THEN GREATEST((c.meta->>'discover_weight')::DOUBLE PRECISION, 0.1)
    ELSE 1
  END DESC,
  c.sort_order ASC,
  c.created_at DESC
LIMIT 1
`, theme, exclude).Scan(
		&item.Slug,
		&item.Name,
		&item.CharacterTypeCode,
		&item.Summary,
		&item.OneLineDefinition,
		&item.CoverURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Character{}, errors.New("no characters")
		}
		return dto.Character{}, fmt.Errorf("random character query: %w", err)
	}
	return item, nil
}

func (r *postgresCatalogRepo) queryHomeFeaturedCharacter(ctx context.Context) (dto.Character, error) {
	var item dto.Character
	err := r.pool.QueryRow(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.is_active = TRUE
  AND c.status = 'published'
ORDER BY
  CASE
    WHEN COALESCE(c.meta->>'home_today', 'false') = 'true' THEN 0
    WHEN COALESCE(c.meta->>'is_featured_home', 'false') = 'true' THEN 1
    ELSE 2
  END ASC,
  CASE
    WHEN COALESCE(c.meta->>'home_sort', '') ~ '^-?[0-9]+$' THEN (c.meta->>'home_sort')::INTEGER
    ELSE c.sort_order
  END ASC,
  c.created_at DESC,
  c.name ASC
LIMIT 1
`).Scan(
		&item.Slug,
		&item.Name,
		&item.CharacterTypeCode,
		&item.Summary,
		&item.OneLineDefinition,
		&item.CoverURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Character{}, nil
		}
		return dto.Character{}, fmt.Errorf("home featured character query: %w", err)
	}
	return item, nil
}

func (r *postgresCatalogRepo) queryHomeLatestCharacters(ctx context.Context, limit int) ([]dto.Character, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.is_active = TRUE
  AND c.status = 'published'
ORDER BY c.created_at DESC, c.sort_order ASC, c.name ASC
LIMIT $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("home latest characters query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Character, 0, limit)
	for rows.Next() {
		var item dto.Character
		if err := rows.Scan(
			&item.Slug,
			&item.Name,
			&item.CharacterTypeCode,
			&item.Summary,
			&item.OneLineDefinition,
			&item.CoverURL,
		); err != nil {
			return nil, fmt.Errorf("scan home latest character: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) queryHomeFeaturedSongs(ctx context.Context, limit int) ([]dto.Song, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  s.slug,
  s.title,
  c.slug AS character_slug,
  COALESCE(s.cover_url, '') AS cover_url,
  COALESCE(s.audio_url, '') AS audio_url
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE s.is_active = TRUE
  AND s.status = 'published'
ORDER BY
  CASE
    WHEN COALESCE(s.meta->>'is_featured_home', 'false') = 'true' THEN 0
    ELSE 1
  END ASC,
  CASE
    WHEN COALESCE(s.meta->>'home_sort', '') ~ '^-?[0-9]+$' THEN (s.meta->>'home_sort')::INTEGER
    ELSE s.sort_order
  END ASC,
  s.created_at DESC,
  s.title ASC
LIMIT $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("home featured songs query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Song, 0, limit)
	for rows.Next() {
		var item dto.Song
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL); err != nil {
			return nil, fmt.Errorf("scan home featured song: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) queryHomeRecommendedWorks(ctx context.Context, limit int) ([]dto.Work, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  w.slug,
  w.title,
  COALESCE(w.summary, '') AS summary,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE w.is_active = TRUE
ORDER BY
  CASE
    WHEN COALESCE(w.meta->>'is_recommended', 'false') = 'true' THEN 0
    ELSE 1
  END ASC,
  CASE
    WHEN COALESCE(w.meta->>'recommend_sort', '') ~ '^-?[0-9]+$' THEN (w.meta->>'recommend_sort')::INTEGER
    ELSE w.sort_order
  END ASC,
  w.created_at DESC,
  w.title ASC
LIMIT $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("home recommended works query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Work, 0, limit)
	for rows.Next() {
		var item dto.Work
		if err := rows.Scan(&item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.TypeCode); err != nil {
			return nil, fmt.Errorf("scan home recommended work: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) queryHomeThemes(ctx context.Context, limit int) ([]dto.Theme, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  t.slug,
  t.name_zh AS name,
  COALESCE(t.summary, '') AS summary,
  COALESCE(t.cover_url, '') AS cover_url
FROM public.pm_themes t
WHERE t.is_active = TRUE
ORDER BY t.sort_order ASC, t.created_at DESC, t.name_zh ASC
LIMIT $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("home themes query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Theme, 0, limit)
	for rows.Next() {
		var item dto.Theme
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL); err != nil {
			return nil, fmt.Errorf("scan home theme: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) ListCharacters() ([]dto.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.is_active = TRUE
  AND c.status = 'published'
ORDER BY c.sort_order ASC, c.name ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list characters query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Character, 0)
	for rows.Next() {
		var item dto.Character
		if err := rows.Scan(
			&item.Slug,
			&item.Name,
			&item.CharacterTypeCode,
			&item.Summary,
			&item.OneLineDefinition,
			&item.CoverURL,
		); err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) GetCharacterDetail(slug string) (dto.CharacterDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d dto.CharacterDetail
	err := r.pool.QueryRow(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.core_identity, '') AS core_identity,
  COALESCE(c.core_fear, '') AS core_fear,
  COALESCE(c.core_conflict, '') AS core_conflict,
  COALESCE(c.emotional_tone, '') AS emotional_tone
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.slug = $1
  AND c.is_active = TRUE
  AND c.status = 'published'
`, slug).Scan(
		&d.Slug,
		&d.Name,
		&d.CharacterTypeCode,
		&d.Summary,
		&d.OneLineDefinition,
		&d.CoverURL,
		&d.CoreIdentity,
		&d.CoreFear,
		&d.CoreConflict,
		&d.EmotionalTone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.CharacterDetail{}, errors.New("character not found")
		}
		return dto.CharacterDetail{}, fmt.Errorf("get character detail query: %w", err)
	}

	d.RelatedSongs, _ = r.listSongsByCharacterSlug(ctx, slug)
	d.RelatedThemes, _ = r.listThemesByCharacterSlug(ctx, slug)
	d.RelatedWorks, _ = r.listWorksByCharacterSlug(ctx, slug)
	d.RelatedCreator, _ = r.listCreatorsByCharacterSlug(ctx, slug)

	return d, nil
}

func (r *postgresCatalogRepo) ListWorks() ([]dto.Work, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  w.slug,
  w.title,
  COALESCE(w.summary, '') AS summary,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE w.is_active = TRUE
ORDER BY w.sort_order ASC, w.title ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list works query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Work, 0)
	for rows.Next() {
		var item dto.Work
		if err := rows.Scan(&item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.TypeCode); err != nil {
			return nil, fmt.Errorf("scan work: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) GetWorkDetail(slug string) (dto.Work, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item dto.Work
	err := r.pool.QueryRow(ctx, `
SELECT
  w.slug,
  w.title,
  COALESCE(w.summary, '') AS summary,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE w.slug = $1
  AND w.is_active = TRUE
`, slug).Scan(&item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.TypeCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Work{}, errors.New("work not found")
		}
		return dto.Work{}, fmt.Errorf("get work detail query: %w", err)
	}
	return item, nil
}

func (r *postgresCatalogRepo) ListCreators() ([]dto.Creator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  c.slug,
  c.name,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_creators c
WHERE c.is_active = TRUE
ORDER BY c.sort_order ASC, c.name ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list creators query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Creator, 0)
	for rows.Next() {
		var item dto.Creator
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL); err != nil {
			return nil, fmt.Errorf("scan creator: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) GetCreatorDetail(slug string) (dto.Creator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item dto.Creator
	err := r.pool.QueryRow(ctx, `
SELECT
  c.slug,
  c.name,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_creators c
WHERE c.slug = $1
  AND c.is_active = TRUE
`, slug).Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Creator{}, errors.New("creator not found")
		}
		return dto.Creator{}, fmt.Errorf("get creator detail query: %w", err)
	}
	return item, nil
}

func (r *postgresCatalogRepo) ListThemes() ([]dto.Theme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  t.slug,
  t.name_zh AS name,
  COALESCE(t.summary, '') AS summary,
  COALESCE(t.cover_url, '') AS cover_url
FROM public.pm_themes t
WHERE t.is_active = TRUE
ORDER BY t.sort_order ASC, t.name_zh ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list themes query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Theme, 0)
	for rows.Next() {
		var item dto.Theme
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL); err != nil {
			return nil, fmt.Errorf("scan theme: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) GetThemeDetail(slug string) (dto.ThemeDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d dto.ThemeDetail
	err := r.pool.QueryRow(ctx, `
SELECT
  t.slug,
  t.name_zh AS name,
  COALESCE(t.summary, '') AS summary,
  COALESCE(t.cover_url, '') AS cover_url
FROM public.pm_themes t
WHERE t.slug = $1
  AND t.is_active = TRUE
`, slug).Scan(&d.Slug, &d.Name, &d.Summary, &d.CoverURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ThemeDetail{}, errors.New("theme not found")
		}
		return dto.ThemeDetail{}, fmt.Errorf("get theme detail query: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_character_themes x
JOIN public.pm_themes t ON t.id = x.theme_id
JOIN public.pm_characters c ON c.id = x.character_id
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE t.slug = $1
  AND c.is_active = TRUE
  AND c.status = 'published'
ORDER BY x.is_primary DESC, x.weight DESC, c.sort_order ASC, c.name ASC
`, slug)
	if err != nil {
		return dto.ThemeDetail{}, fmt.Errorf("theme characters query: %w", err)
	}
	defer rows.Close()

	chars := make([]dto.Character, 0)
	for rows.Next() {
		var item dto.Character
		if err := rows.Scan(
			&item.Slug,
			&item.Name,
			&item.CharacterTypeCode,
			&item.Summary,
			&item.OneLineDefinition,
			&item.CoverURL,
		); err != nil {
			return dto.ThemeDetail{}, fmt.Errorf("scan theme character: %w", err)
		}
		chars = append(chars, item)
	}
	d.Characters = chars
	return d, rows.Err()
}

func (r *postgresCatalogRepo) ListSongs() ([]dto.Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  s.slug,
  s.title,
  c.slug AS character_slug,
  COALESCE(s.cover_url, '') AS cover_url,
  COALESCE(s.audio_url, '') AS audio_url
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE s.is_active = TRUE
  AND s.status = 'published'
ORDER BY s.sort_order ASC, s.title ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list songs query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.Song, 0)
	for rows.Next() {
		var item dto.Song
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL); err != nil {
			return nil, fmt.Errorf("scan song: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listSongsByCharacterSlug(ctx context.Context, slug string) ([]dto.Song, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  s.slug,
  s.title,
  c.slug AS character_slug,
  COALESCE(s.cover_url, '') AS cover_url,
  COALESCE(s.audio_url, '') AS audio_url
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE c.slug = $1
  AND s.is_active = TRUE
  AND s.status = 'published'
ORDER BY s.sort_order ASC, s.title ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Song, 0)
	for rows.Next() {
		var item dto.Song
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listThemesByCharacterSlug(ctx context.Context, slug string) ([]dto.Theme, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  t.slug,
  t.name_zh AS name,
  COALESCE(t.summary, '') AS summary,
  COALESCE(t.cover_url, '') AS cover_url
FROM public.pm_character_themes x
JOIN public.pm_themes t ON t.id = x.theme_id
JOIN public.pm_characters c ON c.id = x.character_id
WHERE c.slug = $1
  AND t.is_active = TRUE
ORDER BY x.is_primary DESC, x.weight DESC, t.sort_order ASC, t.name_zh ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Theme, 0)
	for rows.Next() {
		var item dto.Theme
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listWorksByCharacterSlug(ctx context.Context, slug string) ([]dto.Work, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  w.slug,
  w.title,
  COALESCE(w.summary, '') AS summary,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_character_works x
JOIN public.pm_works w ON w.id = x.work_id
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
JOIN public.pm_characters c ON c.id = x.character_id
WHERE c.slug = $1
  AND w.is_active = TRUE
ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Work, 0)
	for rows.Next() {
		var item dto.Work
		if err := rows.Scan(&item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.TypeCode); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listCreatorsByCharacterSlug(ctx context.Context, slug string) ([]dto.Creator, error) {
	rows, err := r.pool.Query(ctx, `
SELECT DISTINCT
  cr.slug,
  cr.name,
  COALESCE(cr.summary, '') AS summary,
  COALESCE(cr.cover_url, '') AS cover_url
FROM public.pm_character_works cw
JOIN public.pm_characters c ON c.id = cw.character_id
JOIN public.pm_work_creators wc ON wc.work_id = cw.work_id
JOIN public.pm_creators cr ON cr.id = wc.creator_id
WHERE c.slug = $1
  AND cr.is_active = TRUE
ORDER BY cr.name ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Creator, 0)
	for rows.Next() {
		var item dto.Creator
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func limitCharacters(in []dto.Character, n int) []dto.Character {
	if len(in) <= n {
		return in
	}
	return in[:n]
}

func limitWorks(in []dto.Work, n int) []dto.Work {
	if len(in) <= n {
		return in
	}
	return in[:n]
}

func limitThemes(in []dto.Theme, n int) []dto.Theme {
	if len(in) <= n {
		return in
	}
	return in[:n]
}

func limitSongs(in []dto.Song, n int) []dto.Song {
	if len(in) <= n {
		return in
	}
	return in[:n]
}
