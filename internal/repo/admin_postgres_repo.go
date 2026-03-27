package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresAdminRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresAdminRepo(pool *pgxpool.Pool) AdminRepo {
	return &postgresAdminRepo{pool: pool}
}

func (r *postgresAdminRepo) ListAdminCharacters() ([]dto.AdminCharacter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.cover_url, '') AS cover_url,
  c.status,
  COALESCE(ct.name_zh, ct.code) AS type_name,
  ct.code AS character_type_code,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.core_identity, '') AS core_identity,
  COALESCE(c.core_fear, '') AS core_fear,
  COALESCE(c.core_conflict, '') AS core_conflict,
  COALESCE(c.emotional_tone, '') AS emotional_tone,
  COALESCE(c.emotional_temperature, '') AS emotional_temperature,
  COALESCE(c.gender, '') AS gender,
  COALESCE(rg.code, '') AS region_code,
  COALESCE(cr.code, '') AS cultural_region_code,
  COALESCE(c.dominant_emotions, ARRAY[]::text[]),
  COALESCE(c.suppressed_emotions, ARRAY[]::text[]),
  COALESCE(c.values_tags, ARRAY[]::text[]),
  COALESCE(c.symbolic_images, ARRAY[]::text[]),
  COALESCE(c.elements, ARRAY[]::text[]),
  COALESCE(c.sort_order, 0),
  COALESCE(c.relationship_profile::text, '{}'::text),
  COALESCE(c.timeline::text, '[]'),
  COALESCE(c.meta::text, '{}'::text),
  COALESCE((
    SELECT array_agg(w.slug ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC)
    FROM public.pm_character_works x
    JOIN public.pm_works w ON w.id = x.work_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS work_slugs,
  COALESCE((
    SELECT array_agg(w.title ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC)
    FROM public.pm_character_works x
    JOIN public.pm_works w ON w.id = x.work_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS work_names,
  COALESCE((
    SELECT array_agg(t.slug ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC)
    FROM public.pm_character_themes x
    JOIN public.pm_themes t ON t.id = x.theme_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS theme_slugs,
  COALESCE((
    SELECT array_agg(t.name_zh ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC)
    FROM public.pm_character_themes x
    JOIN public.pm_themes t ON t.id = x.theme_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS theme_names,
  COALESCE((
    SELECT array_agg(s.slug ORDER BY s.sort_order ASC, s.title ASC)
    FROM public.pm_songs s
    WHERE s.character_id = c.id AND s.is_active = TRUE
  ), ARRAY[]::text[]) AS song_slugs,
  EXISTS(SELECT 1 FROM public.pm_songs s WHERE s.character_id = c.id AND s.is_active = TRUE) AS has_song
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
ORDER BY c.sort_order ASC, c.updated_at DESC, c.name ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin characters query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.AdminCharacter, 0)
	for rows.Next() {
		item, err := scanAdminCharacter(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresAdminRepo) GetAdminCharacter(ref string) (dto.AdminCharacter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.cover_url, '') AS cover_url,
  c.status,
  COALESCE(ct.name_zh, ct.code) AS type_name,
  ct.code AS character_type_code,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.core_identity, '') AS core_identity,
  COALESCE(c.core_fear, '') AS core_fear,
  COALESCE(c.core_conflict, '') AS core_conflict,
  COALESCE(c.emotional_tone, '') AS emotional_tone,
  COALESCE(c.emotional_temperature, '') AS emotional_temperature,
  COALESCE(c.gender, '') AS gender,
  COALESCE(rg.code, '') AS region_code,
  COALESCE(cr.code, '') AS cultural_region_code,
  COALESCE(c.dominant_emotions, ARRAY[]::text[]),
  COALESCE(c.suppressed_emotions, ARRAY[]::text[]),
  COALESCE(c.values_tags, ARRAY[]::text[]),
  COALESCE(c.symbolic_images, ARRAY[]::text[]),
  COALESCE(c.elements, ARRAY[]::text[]),
  COALESCE(c.sort_order, 0),
  COALESCE(c.relationship_profile::text, '{}'::text),
  COALESCE(c.timeline::text, '[]'),
  COALESCE(c.meta::text, '{}'::text),
  COALESCE((
    SELECT array_agg(w.slug ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC)
    FROM public.pm_character_works x
    JOIN public.pm_works w ON w.id = x.work_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS work_slugs,
  COALESCE((
    SELECT array_agg(w.title ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC)
    FROM public.pm_character_works x
    JOIN public.pm_works w ON w.id = x.work_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS work_names,
  COALESCE((
    SELECT array_agg(t.slug ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC)
    FROM public.pm_character_themes x
    JOIN public.pm_themes t ON t.id = x.theme_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS theme_slugs,
  COALESCE((
    SELECT array_agg(t.name_zh ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC)
    FROM public.pm_character_themes x
    JOIN public.pm_themes t ON t.id = x.theme_id
    WHERE x.character_id = c.id
  ), ARRAY[]::text[]) AS theme_names,
  COALESCE((
    SELECT array_agg(s.slug ORDER BY s.sort_order ASC, s.title ASC)
    FROM public.pm_songs s
    WHERE s.character_id = c.id AND s.is_active = TRUE
  ), ARRAY[]::text[]) AS song_slugs,
  EXISTS(SELECT 1 FROM public.pm_songs s WHERE s.character_id = c.id AND s.is_active = TRUE) AS has_song
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
WHERE c.id::text = $1 OR c.slug = $1
LIMIT 1
`, ref)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("get admin character query: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return dto.AdminCharacter{}, errors.New("admin character not found")
	}
	item, err := scanAdminCharacter(rows)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	return item, nil
}

func (r *postgresAdminRepo) CreateAdminCharacter(in dto.AdminCharacter) (dto.AdminCharacter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	defer tx.Rollback(ctx)

	charTypeID, err := lookupIDByCode(ctx, tx, "public.pm_character_types", in.CharacterTypeCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup character type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup cultural region: %w", err)
	}

	relJSON, _ := json.Marshal(in.RelationshipProfile)
	timelineJSON, _ := json.Marshal(in.Timeline)
	metaJSON, _ := json.Marshal(buildAdminCharacterMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	var id string
	err = tx.QueryRow(ctx, `
INSERT INTO public.pm_characters (
  character_type_id, name, slug, aliases, gender, region_id, cultural_region_id, summary, cover_url,
  one_line_definition, core_identity, core_fear, core_conflict, emotional_tone, emotional_temperature,
  dominant_emotions, suppressed_emotions, values_tags, symbolic_images, elements,
  relationship_profile, timeline, meta, sort_order, status, is_active
) VALUES (
  $1,$2,$3,COALESCE($4, ARRAY[]::text[]),NULLIF($5,''),$6,$7,$8,$9,
  $10,$11,$12,$13,$14,NULLIF($15,''),COALESCE($16, ARRAY[]::text[]),COALESCE($17, ARRAY[]::text[]),COALESCE($18, ARRAY[]::text[]),COALESCE($19, ARRAY[]::text[]),COALESCE($20, ARRAY[]::text[]),
  COALESCE($21::jsonb,'{}'::jsonb),COALESCE($22::jsonb,'[]'::jsonb),COALESCE($23::jsonb,'{}'::jsonb),COALESCE($24,0),$25,$26
) RETURNING id::text
`, charTypeID, in.Name, in.Slug, []string{}, emptyToNil(in.Gender), regionID, culturalID, in.Summary, in.CoverURL,
		in.OneLineDefinition, in.CoreIdentity, in.CoreFear, in.CoreConflict, in.EmotionalTone, emptyToNil(in.EmotionalTemperature),
		in.DominantEmotions, in.SuppressedEmotions, in.ValuesTags, in.SymbolicImages, in.Elements,
		string(relJSON), string(timelineJSON), string(metaJSON), in.SortOrder, status, isPublishedStatus(status),
	).Scan(&id)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("insert character: %w", err)
	}

	if err := replaceCharacterWorks(ctx, tx, id, in.WorkSlugs); err != nil {
		return dto.AdminCharacter{}, err
	}
	if err := replaceCharacterThemes(ctx, tx, id, in.ThemeSlugs); err != nil {
		return dto.AdminCharacter{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminCharacter{}, err
	}
	return r.GetAdminCharacter(id)
}

func (r *postgresAdminRepo) UpdateAdminCharacter(ref string, in dto.AdminCharacter) (dto.AdminCharacter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	defer tx.Rollback(ctx)

	charID, err := resolveIDByRef(ctx, tx, "public.pm_characters", ref)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	charTypeID, err := lookupIDByCode(ctx, tx, "public.pm_character_types", in.CharacterTypeCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup character type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("lookup cultural region: %w", err)
	}

	relJSON, _ := json.Marshal(in.RelationshipProfile)
	timelineJSON, _ := json.Marshal(in.Timeline)
	metaJSON, _ := json.Marshal(buildAdminCharacterMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	_, err = tx.Exec(ctx, `
UPDATE public.pm_characters
SET character_type_id=$2,
    name=$3,
    slug=$4,
    gender=NULLIF($5,''),
    region_id=$6,
    cultural_region_id=$7,
    summary=$8,
    cover_url=$9,
    one_line_definition=$10,
    core_identity=$11,
    core_fear=$12,
    core_conflict=$13,
    emotional_tone=$14,
    emotional_temperature=NULLIF($15,''),
    dominant_emotions=COALESCE($16, ARRAY[]::text[]),
    suppressed_emotions=COALESCE($17, ARRAY[]::text[]),
    values_tags=COALESCE($18, ARRAY[]::text[]),
    symbolic_images=COALESCE($19, ARRAY[]::text[]),
    elements=COALESCE($20, ARRAY[]::text[]),
    sort_order=COALESCE($21, pm_characters.sort_order),
    relationship_profile=COALESCE($22::jsonb,'{}'::jsonb),
    timeline=COALESCE($23::jsonb,'[]'::jsonb),
    meta=COALESCE(pm_characters.meta, '{}'::jsonb) || COALESCE($24::jsonb,'{}'::jsonb),
    status=$25,
    is_active=$26,
    updated_at=NOW()
WHERE id=$1
`, charID, charTypeID, in.Name, in.Slug, emptyToNil(in.Gender), regionID, culturalID, in.Summary, in.CoverURL,
		in.OneLineDefinition, in.CoreIdentity, in.CoreFear, in.CoreConflict, in.EmotionalTone, emptyToNil(in.EmotionalTemperature),
		in.DominantEmotions, in.SuppressedEmotions, in.ValuesTags, in.SymbolicImages, in.Elements,
		in.SortOrder, string(relJSON), string(timelineJSON), string(metaJSON), status, isPublishedStatus(status))
	if err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("update character: %w", err)
	}

	if err := replaceCharacterWorks(ctx, tx, charID, in.WorkSlugs); err != nil {
		return dto.AdminCharacter{}, err
	}
	if err := replaceCharacterThemes(ctx, tx, charID, in.ThemeSlugs); err != nil {
		return dto.AdminCharacter{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminCharacter{}, err
	}
	return r.GetAdminCharacter(charID)
}

func (r *postgresAdminRepo) DeleteAdminCharacter(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := r.pool.Exec(ctx, `UPDATE public.pm_characters SET is_active=FALSE, status='archived', updated_at=NOW() WHERE id::text=$1 OR slug=$1`, ref)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin character not found")
	}
	return nil
}

func (r *postgresAdminRepo) ListAdminSongs() ([]dto.AdminSong, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.pool.Query(ctx, `
SELECT s.id::text, s.slug, s.title, COALESCE(s.summary,''), COALESCE(s.cover_url,''), COALESCE(s.audio_url,''), s.status,
       c.slug, c.name, COALESCE(s.song_core_theme,''), COALESCE(s.song_styles, ARRAY[]::text[]), COALESCE(s.song_emotional_curve, ARRAY[]::text[]),
       COALESCE(s.prompt,''), COALESCE(s.lyrics,''), COALESCE(s.sort_order,0), COALESCE(s.meta::text, '{}'::text)
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
ORDER BY s.sort_order ASC, s.updated_at DESC, s.title ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin songs query: %w", err)
	}
	defer rows.Close()
	list := make([]dto.AdminSong, 0)
	for rows.Next() {
		var item dto.AdminSong
		var metaJSON string
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.AudioURL, &item.Status, &item.CharacterSlug, &item.CharacterName, &item.CoreTheme, &item.Styles, &item.EmotionalCurve, &item.Prompt, &item.Lyrics, &item.SortOrder, &metaJSON); err != nil {
			return nil, fmt.Errorf("scan admin song: %w", err)
		}
		applyAdminSongMeta(&item, metaJSON)
		list = append(list, item)
	}
	return list, rows.Err()
}
func (r *postgresAdminRepo) GetAdminSong(ref string) (dto.AdminSong, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var item dto.AdminSong
	var metaJSON string
	err := r.pool.QueryRow(ctx, `
SELECT s.id::text, s.slug, s.title, COALESCE(s.summary,''), COALESCE(s.cover_url,''), COALESCE(s.audio_url,''), s.status,
       c.slug, c.name, COALESCE(s.song_core_theme,''), COALESCE(s.song_styles, ARRAY[]::text[]), COALESCE(s.song_emotional_curve, ARRAY[]::text[]),
       COALESCE(s.prompt,''), COALESCE(s.lyrics,''), COALESCE(s.sort_order,0), COALESCE(s.meta::text, '{}'::text)
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE s.id::text=$1 OR s.slug=$1
LIMIT 1
`, ref).Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.AudioURL, &item.Status, &item.CharacterSlug, &item.CharacterName, &item.CoreTheme, &item.Styles, &item.EmotionalCurve, &item.Prompt, &item.Lyrics, &item.SortOrder, &metaJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AdminSong{}, errors.New("admin song not found")
		}
		return dto.AdminSong{}, err
	}
	applyAdminSongMeta(&item, metaJSON)
	return item, nil
}
func (r *postgresAdminRepo) CreateAdminSong(in dto.AdminSong) (dto.AdminSong, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	charID, err := lookupIDBySlug(ctx, r.pool, "public.pm_characters", in.CharacterSlug)
	if err != nil {
		return dto.AdminSong{}, fmt.Errorf("lookup character: %w", err)
	}
	var id string
	metaJSON, _ := json.Marshal(buildAdminSongMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")
	err = r.pool.QueryRow(ctx, `
INSERT INTO public.pm_songs (
  character_id, title, slug, summary, cover_url, audio_url, song_core_theme, song_styles, song_emotional_curve, prompt, lyrics, meta, sort_order, status, is_active
) VALUES ($1,$2,$3,$4,$5,$6,$7,COALESCE($8, ARRAY[]::text[]),COALESCE($9, ARRAY[]::text[]),$10,$11,COALESCE($12::jsonb,'{}'::jsonb),COALESCE($13,0),$14,$15)
RETURNING id::text
`, charID, in.Title, in.Slug, in.Summary, in.CoverURL, in.AudioURL, in.CoreTheme, in.Styles, in.EmotionalCurve, in.Prompt, in.Lyrics, string(metaJSON), in.SortOrder, status, isPublishedStatus(status)).Scan(&id)
	if err != nil {
		return dto.AdminSong{}, fmt.Errorf("insert song: %w", err)
	}
	return r.GetAdminSong(id)
}
func (r *postgresAdminRepo) UpdateAdminSong(ref string, in dto.AdminSong) (dto.AdminSong, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	songID, err := resolveIDByRef(ctx, r.pool, "public.pm_songs", ref)
	if err != nil {
		return dto.AdminSong{}, err
	}
	charID, err := lookupIDBySlug(ctx, r.pool, "public.pm_characters", in.CharacterSlug)
	if err != nil {
		return dto.AdminSong{}, fmt.Errorf("lookup character: %w", err)
	}
	metaJSON, _ := json.Marshal(buildAdminSongMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")
	_, err = r.pool.Exec(ctx, `
UPDATE public.pm_songs
SET character_id=$2, title=$3, slug=$4, summary=$5, cover_url=$6, audio_url=$7, song_core_theme=$8,
    song_styles=COALESCE($9, ARRAY[]::text[]), song_emotional_curve=COALESCE($10, ARRAY[]::text[]),
    prompt=$11, lyrics=$12, meta=COALESCE(pm_songs.meta, '{}'::jsonb) || COALESCE($13::jsonb, '{}'::jsonb), sort_order=COALESCE($14, pm_songs.sort_order),
    status=$15, is_active=$16, updated_at=NOW()
WHERE id=$1
`, songID, charID, in.Title, in.Slug, in.Summary, in.CoverURL, in.AudioURL, in.CoreTheme, in.Styles, in.EmotionalCurve, in.Prompt, in.Lyrics, string(metaJSON), in.SortOrder, status, isPublishedStatus(status))
	if err != nil {
		return dto.AdminSong{}, fmt.Errorf("update song: %w", err)
	}
	return r.GetAdminSong(songID)
}
func (r *postgresAdminRepo) DeleteAdminSong(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := r.pool.Exec(ctx, `UPDATE public.pm_songs SET is_active=FALSE, status='archived', updated_at=NOW() WHERE id::text=$1 OR slug=$1`, ref)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin song not found")
	}
	return nil
}

func (r *postgresAdminRepo) ListAdminThemes() ([]dto.AdminTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.pool.Query(ctx, `
SELECT t.id::text, t.slug, t.name_zh, t.code, t.category, COALESCE(t.summary,''), COALESCE(t.cover_url,''), COALESCE(t.sort_order,0),
       CASE WHEN t.is_active THEN 'published' ELSE 'archived' END AS status,
       COALESCE((SELECT array_agg(c.slug ORDER BY x.is_primary DESC, x.weight DESC, c.name ASC)
                 FROM public.pm_character_themes x
                 JOIN public.pm_characters c ON c.id = x.character_id
                 WHERE x.theme_id = t.id AND c.is_active = TRUE), ARRAY[]::text[]) AS character_slugs
FROM public.pm_themes t
ORDER BY t.sort_order ASC, t.updated_at DESC, t.name_zh ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin themes query: %w", err)
	}
	defer rows.Close()
	list := make([]dto.AdminTheme, 0)
	for rows.Next() {
		var item dto.AdminTheme
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &item.Code, &item.Category, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.CharacterSlugs); err != nil {
			return nil, fmt.Errorf("scan admin theme: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
func (r *postgresAdminRepo) GetAdminTheme(ref string) (dto.AdminTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var item dto.AdminTheme
	err := r.pool.QueryRow(ctx, `
SELECT t.id::text, t.slug, t.name_zh, t.code, t.category, COALESCE(t.summary,''), COALESCE(t.cover_url,''), COALESCE(t.sort_order,0),
       CASE WHEN t.is_active THEN 'published' ELSE 'archived' END AS status,
       COALESCE((SELECT array_agg(c.slug ORDER BY x.is_primary DESC, x.weight DESC, c.name ASC)
                 FROM public.pm_character_themes x
                 JOIN public.pm_characters c ON c.id = x.character_id
                 WHERE x.theme_id = t.id AND c.is_active = TRUE), ARRAY[]::text[]) AS character_slugs
FROM public.pm_themes t
WHERE t.id::text=$1 OR t.slug=$1
LIMIT 1
`, ref).Scan(&item.ID, &item.Slug, &item.Name, &item.Code, &item.Category, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.CharacterSlugs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AdminTheme{}, errors.New("admin theme not found")
		}
		return dto.AdminTheme{}, err
	}
	return item, nil
}
func (r *postgresAdminRepo) CreateAdminTheme(in dto.AdminTheme) (dto.AdminTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminTheme{}, err
	}
	defer tx.Rollback(ctx)

	var id string
	isActive := isPublishedStatus(in.Status)
	err = tx.QueryRow(ctx, `
INSERT INTO public.pm_themes (code, slug, name_zh, category, summary, cover_url, sort_order, is_active)
VALUES ($1,$2,$3,COALESCE(NULLIF($4,''),'general'),$5,$6,COALESCE($7,0),$8)
RETURNING id::text
`, in.Code, in.Slug, in.Name, emptyToNil(in.Category), in.Summary, in.CoverURL, in.SortOrder, isActive).Scan(&id)
	if err != nil {
		return dto.AdminTheme{}, fmt.Errorf("insert theme: %w", err)
	}

	if err := replaceThemeCharacters(ctx, tx, id, in.CharacterSlugs); err != nil {
		return dto.AdminTheme{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminTheme{}, err
	}
	return r.GetAdminTheme(id)
}
func (r *postgresAdminRepo) UpdateAdminTheme(ref string, in dto.AdminTheme) (dto.AdminTheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminTheme{}, err
	}
	defer tx.Rollback(ctx)

	themeID, err := resolveIDByRef(ctx, tx, "public.pm_themes", ref)
	if err != nil {
		return dto.AdminTheme{}, err
	}

	_, err = tx.Exec(ctx, `
UPDATE public.pm_themes
SET code=$2, slug=$3, name_zh=$4, category=COALESCE(NULLIF($5,''),category), summary=$6, cover_url=$7, sort_order=COALESCE($8, pm_themes.sort_order), is_active=$9, updated_at=NOW()
WHERE id=$1
`, themeID, in.Code, in.Slug, in.Name, emptyToNil(in.Category), in.Summary, in.CoverURL, in.SortOrder, isPublishedStatus(in.Status))
	if err != nil {
		return dto.AdminTheme{}, fmt.Errorf("update theme: %w", err)
	}

	if err := replaceThemeCharacters(ctx, tx, themeID, in.CharacterSlugs); err != nil {
		return dto.AdminTheme{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminTheme{}, err
	}
	return r.GetAdminTheme(themeID)
}
func (r *postgresAdminRepo) DeleteAdminTheme(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := r.pool.Exec(ctx, `UPDATE public.pm_themes SET is_active=FALSE, updated_at=NOW() WHERE id::text=$1 OR slug=$1`, ref)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin theme not found")
	}
	return nil
}

func scanAdminCharacter(rows pgx.Rows) (dto.AdminCharacter, error) {
	var item dto.AdminCharacter
	var relJSON, timelineJSON, metaJSON string
	if err := rows.Scan(
		&item.ID, &item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.Status,
		&item.Type, &item.CharacterTypeCode, &item.OneLineDefinition, &item.CoreIdentity, &item.CoreFear, &item.CoreConflict,
		&item.EmotionalTone, &item.EmotionalTemperature, &item.Gender, &item.RegionCode, &item.CulturalRegionCode,
		&item.DominantEmotions, &item.SuppressedEmotions, &item.ValuesTags, &item.SymbolicImages, &item.Elements, &item.SortOrder,
		&relJSON, &timelineJSON, &metaJSON, &item.WorkSlugs, &item.WorkNames, &item.ThemeSlugs, &item.ThemeNames, &item.SongSlugs, &item.HasSong,
	); err != nil {
		return dto.AdminCharacter{}, fmt.Errorf("scan admin character: %w", err)
	}
	_ = json.Unmarshal([]byte(relJSON), &item.RelationshipProfile)
	_ = json.Unmarshal([]byte(timelineJSON), &item.Timeline)
	var meta map[string]any
	if err := json.Unmarshal([]byte(metaJSON), &meta); err == nil {
		if v, ok := meta["primary_motivation"].(string); ok {
			item.PrimaryMotivation = v
		}
		item.HomeToday = adminMetaBool(meta, "home_today")
		item.FeaturedHome = adminMetaBool(meta, "is_featured_home")
		item.HomeSort = adminMetaInt(meta, "home_sort")
		item.DiscoverWeight = adminMetaFloat(meta, "discover_weight")
	}
	return item, nil
}

func replaceCharacterWorks(ctx context.Context, tx pgx.Tx, characterID string, workSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_works WHERE character_id=$1`, characterID); err != nil {
		return err
	}
	for idx, slug := range uniq(workSlugs) {
		if strings.TrimSpace(slug) == "" {
			continue
		}
		workID, err := resolveIDByRef(ctx, tx, "public.pm_works", slug)
		if err != nil {
			return fmt.Errorf("lookup work %s: %w", slug, err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.pm_character_works (character_id, work_id, relation_type, is_primary, sort_order) VALUES ($1,$2,'belongs_to',$3,$4)`, characterID, workID, idx == 0, (idx+1)*10)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceCharacterThemes(ctx context.Context, tx pgx.Tx, characterID string, themeSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_themes WHERE character_id=$1`, characterID); err != nil {
		return err
	}
	for idx, slug := range uniq(themeSlugs) {
		if strings.TrimSpace(slug) == "" {
			continue
		}
		themeID, err := resolveIDByRef(ctx, tx, "public.pm_themes", slug)
		if err != nil {
			return fmt.Errorf("lookup theme %s: %w", slug, err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.pm_character_themes (character_id, theme_id, weight, is_primary) VALUES ($1,$2,$3,$4)`, characterID, themeID, 1.0, idx == 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceThemeCharacters(ctx context.Context, tx pgx.Tx, themeID string, characterSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_themes WHERE theme_id=$1`, themeID); err != nil {
		return err
	}
	for idx, slug := range uniq(characterSlugs) {
		if strings.TrimSpace(slug) == "" {
			continue
		}
		charID, err := resolveIDByRef(ctx, tx, "public.pm_characters", slug)
		if err != nil {
			return fmt.Errorf("lookup character %s: %w", slug, err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.pm_character_themes (character_id, theme_id, weight, is_primary) VALUES ($1,$2,$3,$4)`, charID, themeID, 1.0, idx == 0)
		if err != nil {
			return err
		}
	}
	return nil
}
func lookupIDByCode(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, table, code string) (string, error) {
	var id string
	err := q.QueryRow(ctx, fmt.Sprintf(`SELECT id::text FROM %s WHERE code=$1 AND is_active=TRUE LIMIT 1`, table), code).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("code not found: %s", code)
		}
		return "", err
	}
	return id, nil
}

func optionalLookupIDByCode(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, table, code string) (any, error) {
	if strings.TrimSpace(code) == "" {
		return nil, nil
	}
	id, err := lookupIDByCode(ctx, q, table, code)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func resolveIDByRef(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, table, ref string) (string, error) {
	var id string
	err := q.QueryRow(ctx, fmt.Sprintf(`SELECT id::text FROM %s WHERE (id::text=$1 OR slug=$1) LIMIT 1`, table), ref).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("ref not found: %s", ref)
		}
		return "", err
	}
	return id, nil
}

func lookupIDBySlug(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, table, slug string) (string, error) {
	return resolveIDByRef(ctx, q, table, slug)
}

func emptyToNil(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func normalizeAdminStatus(status string, fallback string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "draft", "published", "archived":
		return strings.ToLower(strings.TrimSpace(status))
	default:
		return fallback
	}
}

func isPublishedStatus(status string) bool {
	return normalizeAdminStatus(status, "draft") == "published"
}

func uniq(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func buildAdminCharacterMeta(in dto.AdminCharacter) map[string]any {
	return map[string]any{
		"primary_motivation": in.PrimaryMotivation,
		"home_today":         in.HomeToday,
		"is_featured_home":   in.FeaturedHome,
		"home_sort":          in.HomeSort,
		"discover_weight":    in.DiscoverWeight,
	}
}

func buildAdminSongMeta(in dto.AdminSong) map[string]any {
	return map[string]any{
		"is_featured_home": in.FeaturedHome,
		"home_sort":        in.HomeSort,
	}
}

func applyAdminSongMeta(item *dto.AdminSong, metaJSON string) {
	var meta map[string]any
	if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
		return
	}
	item.FeaturedHome = adminMetaBool(meta, "is_featured_home")
	item.HomeSort = adminMetaInt(meta, "home_sort")
}

func adminMetaBool(meta map[string]any, key string) bool {
	switch v := meta[key].(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true")
	default:
		return false
	}
}

func adminMetaInt(meta map[string]any, key string) int {
	switch v := meta[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		var out int
		if _, err := fmt.Sscanf(strings.TrimSpace(v), "%d", &out); err == nil {
			return out
		}
	}
	return 0
}

func adminMetaFloat(meta map[string]any, key string) float64 {
	switch v := meta[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		var out float64
		if _, err := fmt.Sscanf(strings.TrimSpace(v), "%f", &out); err == nil {
			return out
		}
	}
	return 0
}
