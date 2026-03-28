package repo

import (
	"context"
	"fmt"
	"time"

	"pm-backend/internal/dto"
)

func pageClause(page, pageSize int) (int, int) {
	page, pageSize = dto.NormalizePage(page, pageSize)
	return pageSize, (page - 1) * pageSize
}

func (r *postgresAdminRepo) PageAdminCharacters(q dto.PageQuery) (dto.PageResult[dto.AdminCharacter], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + q.Keyword + "%"

	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1)
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE ($1 = '' OR c.name ILIKE $2 OR c.slug ILIKE $2 OR COALESCE(c.summary,'') ILIKE $2 OR COALESCE(c.one_line_definition,'') ILIKE $2)
  AND ($3 = '' OR ct.code = $3)
  AND ($4 = '' OR c.status = $4)
`, q.Keyword, pattern, q.CharacterTypeCode, q.Status).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminCharacter]{}, fmt.Errorf("count admin characters: %w", err)
	}

	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, `
SELECT
  c.id::text, c.slug, c.name, COALESCE(c.summary,''), COALESCE(c.cover_url,''), c.status,
  COALESCE(ct.name_zh, ct.code), ct.code, COALESCE(c.one_line_definition,''), COALESCE(c.core_identity,''), COALESCE(c.core_fear,''), COALESCE(c.core_conflict,''),
  COALESCE(c.emotional_tone,''), COALESCE(c.emotional_temperature,''), COALESCE(c.gender,''), COALESCE(rg.code,''), COALESCE(cr.code,''),
  COALESCE(c.dominant_emotions, ARRAY[]::text[]), COALESCE(c.suppressed_emotions, ARRAY[]::text[]), COALESCE(c.values_tags, ARRAY[]::text[]), COALESCE(c.symbolic_images, ARRAY[]::text[]), COALESCE(c.elements, ARRAY[]::text[]),
  COALESCE(c.sort_order, 0),
  COALESCE(c.relationship_profile::text, '{}'::text), COALESCE(c.timeline::text, '[]'), COALESCE(c.meta::text, '{}'::text),
  COALESCE((SELECT array_agg(w.slug ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC) FROM public.pm_character_works x JOIN public.pm_works w ON w.id = x.work_id WHERE x.character_id = c.id), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(w.title ORDER BY x.is_primary DESC, x.sort_order ASC, w.title ASC) FROM public.pm_character_works x JOIN public.pm_works w ON w.id = x.work_id WHERE x.character_id = c.id), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(t.slug ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC) FROM public.pm_character_themes x JOIN public.pm_themes t ON t.id = x.theme_id WHERE x.character_id = c.id), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(t.name_zh ORDER BY x.is_primary DESC, x.weight DESC, t.name_zh ASC) FROM public.pm_character_themes x JOIN public.pm_themes t ON t.id = x.theme_id WHERE x.character_id = c.id), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(s.slug ORDER BY s.sort_order ASC, s.title ASC) FROM public.pm_songs s WHERE s.character_id = c.id AND s.is_active = TRUE), ARRAY[]::text[]),
  EXISTS(SELECT 1 FROM public.pm_songs s WHERE s.character_id = c.id AND s.is_active = TRUE)
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
WHERE ($1 = '' OR c.name ILIKE $2 OR c.slug ILIKE $2 OR COALESCE(c.summary,'') ILIKE $2 OR COALESCE(c.one_line_definition,'') ILIKE $2)
  AND ($3 = '' OR ct.code = $3)
  AND ($4 = '' OR c.status = $4)
ORDER BY c.sort_order ASC, c.updated_at DESC, c.name ASC
LIMIT $5 OFFSET $6
`, q.Keyword, pattern, q.CharacterTypeCode, q.Status, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminCharacter]{}, fmt.Errorf("page admin characters: %w", err)
	}
	defer rows.Close()
	items := make([]dto.AdminCharacter, 0)
	for rows.Next() {
		item, err := scanAdminCharacter(rows)
		if err != nil {
			return dto.PageResult[dto.AdminCharacter]{}, err
		}
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminCharacter]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func (r *postgresAdminRepo) PageAdminSongs(q dto.PageQuery) (dto.PageResult[dto.AdminSong], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + q.Keyword + "%"
	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1)
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE ($1 = '' OR s.title ILIKE $2 OR s.slug ILIKE $2 OR COALESCE(s.summary,'') ILIKE $2 OR c.name ILIKE $2 OR COALESCE(s.song_core_theme,'') ILIKE $2)
  AND ($3 = '' OR s.status = $3)
`, q.Keyword, pattern, q.Status).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminSong]{}, err
	}
	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, `
SELECT s.id::text, s.slug, s.title, COALESCE(s.summary,''), COALESCE(s.cover_url,''), COALESCE(s.audio_url,''), s.status,
       c.slug, c.name, COALESCE(s.song_core_theme,''), COALESCE(s.song_styles, ARRAY[]::text[]), COALESCE(s.song_emotional_curve, ARRAY[]::text[]),
       COALESCE(s.prompt,''), COALESCE(s.lyrics,''), COALESCE(s.sort_order,0), COALESCE(s.meta::text, '{}'::text)
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE ($1 = '' OR s.title ILIKE $2 OR s.slug ILIKE $2 OR COALESCE(s.summary,'') ILIKE $2 OR c.name ILIKE $2 OR COALESCE(s.song_core_theme,'') ILIKE $2)
  AND ($3 = '' OR s.status = $3)
ORDER BY s.sort_order ASC, s.updated_at DESC, s.title ASC
LIMIT $4 OFFSET $5
`, q.Keyword, pattern, q.Status, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminSong]{}, err
	}
	defer rows.Close()
	items := make([]dto.AdminSong, 0)
	for rows.Next() {
		var item dto.AdminSong
		var metaJSON string
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.AudioURL, &item.Status, &item.CharacterSlug, &item.CharacterName, &item.CoreTheme, &item.Styles, &item.EmotionalCurve, &item.Prompt, &item.Lyrics, &item.SortOrder, &metaJSON); err != nil {
			return dto.PageResult[dto.AdminSong]{}, err
		}
		applyAdminSongMeta(&item, metaJSON)
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminSong]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func (r *postgresAdminRepo) PageAdminThemes(q dto.PageQuery) (dto.PageResult[dto.AdminTheme], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + q.Keyword + "%"
	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1) FROM public.pm_themes t
WHERE ($1 = '' OR t.slug ILIKE $2 OR t.code ILIKE $2 OR COALESCE(t.name_zh,'') ILIKE $2 OR COALESCE(t.summary,'') ILIKE $2 OR COALESCE(t.category,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(t.category,'') = $3)
  AND ($4 = '' OR CASE WHEN t.is_active THEN 'published' ELSE 'archived' END = $4)
  AND ($5 = '' OR COALESCE(t.subject_type, 'character') = $5)
`, q.Keyword, pattern, q.Category, q.Status, q.SubjectType).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminTheme]{}, err
	}
	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, `
SELECT t.id::text, t.slug, t.name_zh, t.code, t.category, COALESCE(t.summary,''), COALESCE(t.cover_url,''), COALESCE(t.sort_order,0),
       COALESCE(t.subject_type, 'character'),
       CASE WHEN t.is_active THEN 'published' ELSE 'archived' END,
       COALESCE((SELECT array_agg(c.slug ORDER BY x.is_primary DESC, x.weight DESC, c.name ASC)
                 FROM public.pm_character_themes x JOIN public.pm_characters c ON c.id = x.character_id
                 WHERE x.theme_id = t.id AND c.is_active = TRUE), ARRAY[]::text[]),
       COALESCE((SELECT array_agg(r.slug ORDER BY x.is_primary DESC, x.sort_order ASC, r.slug ASC)
                 FROM public.pm_relation_themes x JOIN public.pm_relations r ON r.slug = x.relation_slug
                 WHERE x.theme_slug = t.slug AND r.is_active = TRUE), ARRAY[]::text[])
FROM public.pm_themes t
WHERE ($1 = '' OR t.slug ILIKE $2 OR t.code ILIKE $2 OR COALESCE(t.name_zh,'') ILIKE $2 OR COALESCE(t.summary,'') ILIKE $2 OR COALESCE(t.category,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(t.category,'') = $3)
  AND ($4 = '' OR CASE WHEN t.is_active THEN 'published' ELSE 'archived' END = $4)
  AND ($5 = '' OR COALESCE(t.subject_type, 'character') = $5)
ORDER BY t.created_at DESC, t.updated_at DESC, t.name_zh ASC
 LIMIT $6 OFFSET $7
`, q.Keyword, pattern, q.Category, q.Status, q.SubjectType, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminTheme]{}, err
	}
	defer rows.Close()
	items := make([]dto.AdminTheme, 0)
	for rows.Next() {
		var item dto.AdminTheme
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &item.Code, &item.Category, &item.Summary, &item.CoverURL, &item.SortOrder, &item.SubjectType, &item.Status, &item.CharacterSlugs, &item.RelationSlugs); err != nil {
			return dto.PageResult[dto.AdminTheme]{}, err
		}
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminTheme]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func (r *postgresAdminRepo) PageAdminWorks(q dto.PageQuery) (dto.PageResult[dto.AdminWork], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + q.Keyword + "%"
	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1) FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE ($1 = '' OR w.title ILIKE $2 OR w.slug ILIKE $2 OR COALESCE(w.summary,'') ILIKE $2 OR COALESCE(wt.code,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(wt.code,'') = $3)
  AND ($4 = '' OR CASE WHEN w.is_active = FALSE THEN 'archived' ELSE COALESCE(w.meta->>'status', 'published') END = $4)
`, q.Keyword, pattern, q.WorkTypeCode, q.Status).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminWork]{}, err
	}
	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, `
SELECT
  w.id::text, w.slug, w.title, COALESCE(w.summary,''), COALESCE(w.cover_url,''),
  COALESCE(w.sort_order,0),
  CASE WHEN w.is_active = FALSE THEN 'archived' ELSE COALESCE(w.meta->>'status', 'published') END,
  COALESCE(wt.code,''), COALESCE(rg.code,''), COALESCE(cr.code,''), COALESCE(w.release_year,0), COALESCE(w.meta::text, '{}'::text),
  COALESCE((SELECT array_agg(c.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(c.name ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE), ARRAY[]::text[])
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
LEFT JOIN public.pm_regions rg ON rg.id = w.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = w.cultural_region_id
WHERE ($1 = '' OR w.title ILIKE $2 OR w.slug ILIKE $2 OR COALESCE(w.summary,'') ILIKE $2 OR COALESCE(wt.code,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(wt.code,'') = $3)
  AND ($4 = '' OR CASE WHEN w.is_active = FALSE THEN 'archived' ELSE COALESCE(w.meta->>'status', 'published') END = $4)
ORDER BY w.sort_order ASC, w.updated_at DESC, w.title ASC
LIMIT $5 OFFSET $6
`, q.Keyword, pattern, q.WorkTypeCode, q.Status, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminWork]{}, err
	}
	defer rows.Close()
	items := make([]dto.AdminWork, 0)
	for rows.Next() {
		var item dto.AdminWork
		var metaJSON string
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.WorkTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.ReleaseYear, &metaJSON, &item.CreatorSlugs, &item.CreatorNames); err != nil {
			return dto.PageResult[dto.AdminWork]{}, err
		}
		applyAdminWorkMeta(&item, metaJSON)
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminWork]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func (r *postgresAdminRepo) PageAdminCreators(q dto.PageQuery) (dto.PageResult[dto.AdminCreator], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + q.Keyword + "%"
	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1) FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
WHERE ($1 = '' OR c.name ILIKE $2 OR c.slug ILIKE $2 OR COALESCE(c.summary,'') ILIKE $2 OR COALESCE(ct.code,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(ct.code,'') = $3)
  AND ($4 = '' OR CASE WHEN c.is_active = FALSE THEN 'archived' ELSE COALESCE(c.meta->>'status', 'published') END = $4)
`, q.Keyword, pattern, q.CreatorTypeCode, q.Status).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminCreator]{}, err
	}
	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, `
SELECT
  c.id::text, c.slug, c.name, COALESCE(c.summary,''), COALESCE(c.cover_url,''),
  COALESCE(c.sort_order,0),
  CASE WHEN c.is_active = FALSE THEN 'archived' ELSE COALESCE(c.meta->>'status', 'published') END,
  COALESCE(ct.code,''), COALESCE(rg.code,''), COALESCE(cr.code,''),
  COALESCE((SELECT array_agg(w.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE), ARRAY[]::text[]),
  COALESCE((SELECT array_agg(w.title ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE), ARRAY[]::text[])
FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
WHERE ($1 = '' OR c.name ILIKE $2 OR c.slug ILIKE $2 OR COALESCE(c.summary,'') ILIKE $2 OR COALESCE(ct.code,'') ILIKE $2)
  AND ($3 = '' OR COALESCE(ct.code,'') = $3)
  AND ($4 = '' OR CASE WHEN c.is_active = FALSE THEN 'archived' ELSE COALESCE(c.meta->>'status', 'published') END = $4)
ORDER BY c.sort_order ASC, c.updated_at DESC, c.name ASC
LIMIT $5 OFFSET $6
`, q.Keyword, pattern, q.CreatorTypeCode, q.Status, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminCreator]{}, err
	}
	defer rows.Close()
	items := make([]dto.AdminCreator, 0)
	for rows.Next() {
		var item dto.AdminCreator
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.CreatorTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.WorkSlugs, &item.WorkNames); err != nil {
			return dto.PageResult[dto.AdminCreator]{}, err
		}
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminCreator]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func (r *postgresAdminRepo) PageAdminDictItems(dictKey string, page, pageSize int, keyword string) (dto.PageResult[dto.AdminDictItem], error) {
	page, pageSize = dto.NormalizePage(page, pageSize)
	meta, err := dictMeta(dictKey)
	if err != nil {
		return dto.PageResult[dto.AdminDictItem]{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pattern := "%" + keyword + "%"
	var total int
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE ($1 = '' OR code ILIKE $2 OR COALESCE(%s,'') ILIKE $2)`, meta.Table, meta.NameColumn), keyword, pattern).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminDictItem]{}, err
	}
	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
SELECT id::text, code, COALESCE(%s,''), COALESCE(sort_order,0), COALESCE(is_active,TRUE)
FROM %s
WHERE ($1 = '' OR code ILIKE $2 OR COALESCE(%s,'') ILIKE $2)
ORDER BY sort_order ASC, code ASC
LIMIT $3 OFFSET $4
`, meta.NameColumn, meta.Table, meta.NameColumn), keyword, pattern, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminDictItem]{}, err
	}
	defer rows.Close()
	items := make([]dto.AdminDictItem, 0)
	for rows.Next() {
		var item dto.AdminDictItem
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.SortOrder, &item.IsActive); err != nil {
			return dto.PageResult[dto.AdminDictItem]{}, err
		}
		item.DictKey = dictKey
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminDictItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}
