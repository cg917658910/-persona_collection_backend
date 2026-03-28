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
)

const adminRelationSummarySelect = `
SELECT
  r.slug AS id,
  r.slug,
  COALESCE(r.name, '') AS name,
  COALESCE(r.subtitle, '') AS subtitle,
  COALESCE(r.summary, '') AS summary,
  COALESCE(r.one_line_definition, '') AS one_line_definition,
  COALESCE(r.cover_url, '') AS cover_url,
  r.status,
  COALESCE(r.sort_order, 0) AS sort_order,
  COALESCE(r.relation_type_code, '') AS relation_type_code,
  COALESCE(rt.name, r.relation_type_code) AS relation_type_name,
  COALESCE(r.work_slug, '') AS work_slug,
  COALESCE(w.title, '') AS work_name,
  COALESCE(r.source_character_slug, '') AS source_character_slug,
  COALESCE(cs.name, r.source_character_slug) AS source_character_name,
  COALESCE(r.target_character_slug, '') AS target_character_slug,
  COALESCE(ct.name, r.target_character_slug) AS target_character_name
FROM public.pm_relations r
LEFT JOIN public.pm_relation_types rt ON rt.code = r.relation_type_code
LEFT JOIN public.pm_works w ON w.slug = r.work_slug
LEFT JOIN public.pm_characters cs ON cs.slug = r.source_character_slug
LEFT JOIN public.pm_characters ct ON ct.slug = r.target_character_slug
WHERE 1 = 1
`

const adminRelationDetailSelect = `
SELECT
  r.slug AS id,
  r.slug,
  COALESCE(r.name, '') AS name,
  COALESCE(r.subtitle, '') AS subtitle,
  COALESCE(r.summary, '') AS summary,
  COALESCE(r.one_line_definition, '') AS one_line_definition,
  COALESCE(r.cover_url, '') AS cover_url,
  r.status,
  COALESCE(r.sort_order, 0) AS sort_order,
  COALESCE(r.relation_type_code, '') AS relation_type_code,
  COALESCE(rt.name, r.relation_type_code) AS relation_type_name,
  COALESCE(r.work_slug, '') AS work_slug,
  COALESCE(w.title, '') AS work_name,
  COALESCE(r.source_character_slug, '') AS source_character_slug,
  COALESCE(cs.name, r.source_character_slug) AS source_character_name,
  COALESCE(r.target_character_slug, '') AS target_character_slug,
  COALESCE(ct.name, r.target_character_slug) AS target_character_name,
  COALESCE(r.core_dynamic, '') AS core_dynamic,
  COALESCE(r.core_tension, '') AS core_tension,
  COALESCE(r.emotional_tone, '') AS emotional_tone,
  COALESCE(r.emotional_temperature, '') AS emotional_temperature,
  COALESCE(r.connection_trigger, '') AS connection_trigger,
  COALESCE(r.sustaining_mechanism, '') AS sustaining_mechanism,
  COALESCE(r.relation_conflict, '') AS relation_conflict,
  COALESCE(r.relation_arc, '') AS relation_arc,
  COALESCE(r.fate_impact, '') AS fate_impact,
  COALESCE(r.power_structure, '') AS power_structure,
  COALESCE(r.dependency_pattern, '') AS dependency_pattern,
  COALESCE(r.source_perspective, '') AS source_perspective,
  COALESCE(r.source_desire_in_relation, '') AS source_desire_in_relation,
  COALESCE(r.source_fear_in_relation, '') AS source_fear_in_relation,
  COALESCE(r.source_unsaid, '') AS source_unsaid,
  COALESCE(r.target_perspective, '') AS target_perspective,
  COALESCE(r.target_desire_in_relation, '') AS target_desire_in_relation,
  COALESCE(r.target_fear_in_relation, '') AS target_fear_in_relation,
  COALESCE(r.target_unsaid, '') AS target_unsaid,
  COALESCE(r.phenomenology, '{}'::jsonb)::text AS phenomenology_json,
  COALESCE(r.symbolic_images, '[]'::jsonb)::text AS symbolic_images_json,
  COALESCE(r.theme_tags, '[]'::jsonb)::text AS theme_tags_json,
  COALESCE(r.relation_palette, '[]'::jsonb)::text AS relation_palette_json,
  COALESCE(r.meta->'relation_keywords', '[]'::jsonb)::text AS relation_keywords_json,
  COALESCE(r.tension_tags, '[]'::jsonb)::text AS tension_tags_json,
  COALESCE(r.cover_prompt, '') AS cover_prompt,
  COALESCE(r.song_prompt, '') AS song_prompt,
  COALESCE(r.primary_song_slug, '') AS primary_song_slug
FROM public.pm_relations r
LEFT JOIN public.pm_relation_types rt ON rt.code = r.relation_type_code
LEFT JOIN public.pm_works w ON w.slug = r.work_slug
LEFT JOIN public.pm_characters cs ON cs.slug = r.source_character_slug
LEFT JOIN public.pm_characters ct ON ct.slug = r.target_character_slug
WHERE 1 = 1
`

func (r *postgresAdminRepo) ListAdminRelations() ([]dto.AdminRelation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, adminRelationSummarySelect+`
ORDER BY r.sort_order ASC, r.updated_at DESC, r.slug ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin relations query: %w", err)
	}
	defer rows.Close()

	list := make([]dto.AdminRelation, 0)
	for rows.Next() {
		item, err := scanAdminRelationSummary(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresAdminRepo) GetAdminRelation(ref string) (dto.AdminRelation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.pool.QueryRow(ctx, adminRelationDetailSelect+`
  AND r.slug = $1
LIMIT 1
`, strings.TrimSpace(ref))

	item, err := scanAdminRelationDetail(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AdminRelation{}, errors.New("admin relation not found")
		}
		return dto.AdminRelation{}, fmt.Errorf("get admin relation query: %w", err)
	}
	if err := r.populateAdminRelationChildren(ctx, &item); err != nil {
		return dto.AdminRelation{}, err
	}
	return item, nil
}

func (r *postgresAdminRepo) CreateAdminRelation(in dto.AdminRelation) (dto.AdminRelation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	defer tx.Rollback(ctx)

	in = normalizeAdminRelationInput(in)
	if err := ensureRelationRefs(ctx, tx, in); err != nil {
		return dto.AdminRelation{}, err
	}

	phenomenologyJSON, _ := json.Marshal(in.Phenomenology)
	symbolicImagesJSON, _ := json.Marshal(uniq(in.SymbolicImages))
	themeTagsJSON, _ := json.Marshal(uniq(in.ThemeTags))
	relationPaletteJSON, _ := json.Marshal(in.RelationPalette)
	tensionTagsJSON, _ := json.Marshal(uniq(in.TensionTags))
	metaJSON, _ := json.Marshal(buildAdminRelationMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")
	primarySongSlug := pickPrimaryRelationSongSlug(in)

	_, err = tx.Exec(ctx, `
INSERT INTO public.pm_relations (
  slug, name, subtitle, relation_type_code, work_slug,
  source_character_slug, target_character_slug,
  summary, one_line_definition, core_dynamic, core_tension, emotional_tone, emotional_temperature,
  connection_trigger, sustaining_mechanism, relation_conflict, relation_arc, fate_impact,
  power_structure, dependency_pattern,
  source_perspective, target_perspective,
  source_desire_in_relation, source_fear_in_relation, source_unsaid,
  target_desire_in_relation, target_fear_in_relation, target_unsaid,
  phenomenology, symbolic_images, theme_tags, relation_palette, tension_tags,
  cover_url, cover_prompt, song_prompt, primary_song_slug, meta, sort_order, status, is_active
) VALUES (
  $1,$2,NULLIF($3,''),$4,NULLIF($5,''),
  $6,$7,
  NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),NULLIF($12,''),NULLIF($13,''),
  NULLIF($14,''),NULLIF($15,''),NULLIF($16,''),NULLIF($17,''),NULLIF($18,''),
  NULLIF($19,''),NULLIF($20,''),
  NULLIF($21,''),NULLIF($22,''),
  NULLIF($23,''),NULLIF($24,''),NULLIF($25,''),
  NULLIF($26,''),NULLIF($27,''),NULLIF($28,''),
  COALESCE($29::jsonb, '{}'::jsonb), COALESCE($30::jsonb, '[]'::jsonb), COALESCE($31::jsonb, '[]'::jsonb), COALESCE($32::jsonb, '[]'::jsonb), COALESCE($33::jsonb, '[]'::jsonb),
  NULLIF($34,''), NULLIF($35,''), NULLIF($36,''), NULLIF($37,''), COALESCE($38::jsonb, '{}'::jsonb), COALESCE($39,0), $40, $41
)
`, in.Slug, in.Name, in.Subtitle, in.RelationTypeCode, in.WorkSlug,
		in.SourceCharacterSlug, in.TargetCharacterSlug,
		in.Summary, in.OneLineDefinition, in.CoreDynamic, in.CoreTension, in.EmotionalTone, in.EmotionalTemperature,
		in.ConnectionTrigger, in.SustainingMechanism, in.RelationConflict, in.RelationArc, in.FateImpact,
		in.PowerStructure, in.DependencyPattern,
		in.SourcePerspective, in.TargetPerspective,
		in.SourceDesireInRelation, in.SourceFearInRelation, in.SourceUnsaid,
		in.TargetDesireInRelation, in.TargetFearInRelation, in.TargetUnsaid,
		string(phenomenologyJSON), string(symbolicImagesJSON), string(themeTagsJSON), string(relationPaletteJSON), string(tensionTagsJSON),
		in.CoverURL, in.CoverPrompt, in.SongPrompt, primarySongSlug, string(metaJSON), in.SortOrder, status, status != "archived",
	)
	if err != nil {
		return dto.AdminRelation{}, fmt.Errorf("insert relation: %w", err)
	}

	if err := replaceRelationChildren(ctx, tx, in.Slug, in); err != nil {
		return dto.AdminRelation{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return dto.AdminRelation{}, err
	}
	return r.GetAdminRelation(in.Slug)
}

func (r *postgresAdminRepo) UpdateAdminRelation(ref string, in dto.AdminRelation) (dto.AdminRelation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	defer tx.Rollback(ctx)

	ref = strings.TrimSpace(ref)
	in = normalizeAdminRelationInput(in)
	if ref == "" {
		ref = in.Slug
	}
	if strings.TrimSpace(in.Slug) != "" && in.Slug != ref {
		return dto.AdminRelation{}, errors.New("relation slug cannot be changed once created")
	}
	if in.Slug == "" {
		in.Slug = ref
	}
	if err := ensureRelationRefs(ctx, tx, in); err != nil {
		return dto.AdminRelation{}, err
	}

	phenomenologyJSON, _ := json.Marshal(in.Phenomenology)
	symbolicImagesJSON, _ := json.Marshal(uniq(in.SymbolicImages))
	themeTagsJSON, _ := json.Marshal(uniq(in.ThemeTags))
	relationPaletteJSON, _ := json.Marshal(in.RelationPalette)
	tensionTagsJSON, _ := json.Marshal(uniq(in.TensionTags))
	metaJSON, _ := json.Marshal(buildAdminRelationMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")
	primarySongSlug := pickPrimaryRelationSongSlug(in)

	tag, err := tx.Exec(ctx, `
UPDATE public.pm_relations
SET slug = $2,
    name = $3,
    subtitle = NULLIF($4,''),
    relation_type_code = $5,
    work_slug = NULLIF($6,''),
    source_character_slug = $7,
    target_character_slug = $8,
    summary = NULLIF($9,''),
    one_line_definition = NULLIF($10,''),
    core_dynamic = NULLIF($11,''),
    core_tension = NULLIF($12,''),
    emotional_tone = NULLIF($13,''),
    emotional_temperature = NULLIF($14,''),
    connection_trigger = NULLIF($15,''),
    sustaining_mechanism = NULLIF($16,''),
    relation_conflict = NULLIF($17,''),
    relation_arc = NULLIF($18,''),
    fate_impact = NULLIF($19,''),
    power_structure = NULLIF($20,''),
    dependency_pattern = NULLIF($21,''),
    source_perspective = NULLIF($22,''),
    target_perspective = NULLIF($23,''),
    source_desire_in_relation = NULLIF($24,''),
    source_fear_in_relation = NULLIF($25,''),
    source_unsaid = NULLIF($26,''),
    target_desire_in_relation = NULLIF($27,''),
    target_fear_in_relation = NULLIF($28,''),
    target_unsaid = NULLIF($29,''),
    phenomenology = COALESCE($30::jsonb, '{}'::jsonb),
    symbolic_images = COALESCE($31::jsonb, '[]'::jsonb),
    theme_tags = COALESCE($32::jsonb, '[]'::jsonb),
    relation_palette = COALESCE($33::jsonb, '[]'::jsonb),
    tension_tags = COALESCE($34::jsonb, '[]'::jsonb),
    cover_url = NULLIF($35,''),
    cover_prompt = NULLIF($36,''),
    song_prompt = NULLIF($37,''),
    primary_song_slug = NULLIF($38,''),
    meta = COALESCE($39::jsonb, '{}'::jsonb),
    sort_order = COALESCE($40,0),
    status = $41,
    is_active = $42,
    updated_at = NOW()
WHERE slug = $1
`, ref, in.Slug, in.Name, in.Subtitle, in.RelationTypeCode, in.WorkSlug,
		in.SourceCharacterSlug, in.TargetCharacterSlug,
		in.Summary, in.OneLineDefinition, in.CoreDynamic, in.CoreTension, in.EmotionalTone, in.EmotionalTemperature,
		in.ConnectionTrigger, in.SustainingMechanism, in.RelationConflict, in.RelationArc, in.FateImpact,
		in.PowerStructure, in.DependencyPattern,
		in.SourcePerspective, in.TargetPerspective,
		in.SourceDesireInRelation, in.SourceFearInRelation, in.SourceUnsaid,
		in.TargetDesireInRelation, in.TargetFearInRelation, in.TargetUnsaid,
		string(phenomenologyJSON), string(symbolicImagesJSON), string(themeTagsJSON), string(relationPaletteJSON), string(tensionTagsJSON),
		in.CoverURL, in.CoverPrompt, in.SongPrompt, primarySongSlug, string(metaJSON), in.SortOrder, status, status != "archived",
	)
	if err != nil {
		return dto.AdminRelation{}, fmt.Errorf("update relation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return dto.AdminRelation{}, errors.New("admin relation not found")
	}

	if err := replaceRelationChildren(ctx, tx, in.Slug, in); err != nil {
		return dto.AdminRelation{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return dto.AdminRelation{}, err
	}
	return r.GetAdminRelation(in.Slug)
}

func (r *postgresAdminRepo) DeleteAdminRelation(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
UPDATE public.pm_relations
SET is_active = FALSE, status = 'archived', updated_at = NOW()
WHERE slug = $1
`, strings.TrimSpace(ref))
	if err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin relation not found")
	}
	return nil
}

func (r *postgresAdminRepo) PageAdminRelations(q dto.PageQuery) (dto.PageResult[dto.AdminRelation], error) {
	page, pageSize := dto.NormalizePage(q.Page, q.PageSize)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pattern := "%" + q.Keyword + "%"
	var total int
	err := r.pool.QueryRow(ctx, `
SELECT COUNT(1)
FROM public.pm_relations r
LEFT JOIN public.pm_works w ON w.slug = r.work_slug
LEFT JOIN public.pm_characters cs ON cs.slug = r.source_character_slug
LEFT JOIN public.pm_characters ct ON ct.slug = r.target_character_slug
WHERE ($1 = '' OR r.name ILIKE $2 OR r.slug ILIKE $2 OR COALESCE(r.summary,'') ILIKE $2 OR COALESCE(r.one_line_definition,'') ILIKE $2 OR COALESCE(cs.name,'') ILIKE $2 OR COALESCE(ct.name,'') ILIKE $2 OR COALESCE(w.title,'') ILIKE $2)
  AND ($3 = '' OR r.relation_type_code = $3)
  AND ($4 = '' OR r.status = $4)
`, q.Keyword, pattern, q.RelationTypeCode, q.Status).Scan(&total)
	if err != nil {
		return dto.PageResult[dto.AdminRelation]{}, fmt.Errorf("count admin relations: %w", err)
	}

	limit, offset := pageClause(page, pageSize)
	rows, err := r.pool.Query(ctx, adminRelationSummarySelect+`
  AND ($1 = '' OR r.name ILIKE $2 OR r.slug ILIKE $2 OR COALESCE(r.summary,'') ILIKE $2 OR COALESCE(r.one_line_definition,'') ILIKE $2 OR COALESCE(cs.name,'') ILIKE $2 OR COALESCE(ct.name,'') ILIKE $2 OR COALESCE(w.title,'') ILIKE $2)
  AND ($3 = '' OR r.relation_type_code = $3)
  AND ($4 = '' OR r.status = $4)
ORDER BY r.sort_order ASC, r.updated_at DESC, r.slug ASC
LIMIT $5 OFFSET $6
`, q.Keyword, pattern, q.RelationTypeCode, q.Status, limit, offset)
	if err != nil {
		return dto.PageResult[dto.AdminRelation]{}, fmt.Errorf("page admin relations: %w", err)
	}
	defer rows.Close()

	items := make([]dto.AdminRelation, 0)
	for rows.Next() {
		item, err := scanAdminRelationSummary(rows)
		if err != nil {
			return dto.PageResult[dto.AdminRelation]{}, err
		}
		items = append(items, item)
	}
	return dto.PageResult[dto.AdminRelation]{Items: items, Total: total, Page: page, PageSize: pageSize}, rows.Err()
}

func scanAdminRelationSummary(row interface{ Scan(dest ...any) error }) (dto.AdminRelation, error) {
	var item dto.AdminRelation
	err := row.Scan(
		&item.ID,
		&item.Slug,
		&item.Name,
		&item.Subtitle,
		&item.Summary,
		&item.OneLineDefinition,
		&item.CoverURL,
		&item.Status,
		&item.SortOrder,
		&item.RelationTypeCode,
		&item.RelationTypeName,
		&item.WorkSlug,
		&item.WorkName,
		&item.SourceCharacterSlug,
		&item.SourceCharacterName,
		&item.TargetCharacterSlug,
		&item.TargetCharacterName,
	)
	return item, err
}

func scanAdminRelationDetail(row interface{ Scan(dest ...any) error }) (dto.AdminRelation, error) {
	var item dto.AdminRelation
	var phenomenologyJSON, symbolicImagesJSON, themeTagsJSON, relationPaletteJSON, relationKeywordsJSON, tensionTagsJSON string
	err := row.Scan(
		&item.ID,
		&item.Slug,
		&item.Name,
		&item.Subtitle,
		&item.Summary,
		&item.OneLineDefinition,
		&item.CoverURL,
		&item.Status,
		&item.SortOrder,
		&item.RelationTypeCode,
		&item.RelationTypeName,
		&item.WorkSlug,
		&item.WorkName,
		&item.SourceCharacterSlug,
		&item.SourceCharacterName,
		&item.TargetCharacterSlug,
		&item.TargetCharacterName,
		&item.CoreDynamic,
		&item.CoreTension,
		&item.EmotionalTone,
		&item.EmotionalTemperature,
		&item.ConnectionTrigger,
		&item.SustainingMechanism,
		&item.RelationConflict,
		&item.RelationArc,
		&item.FateImpact,
		&item.PowerStructure,
		&item.DependencyPattern,
		&item.SourcePerspective,
		&item.SourceDesireInRelation,
		&item.SourceFearInRelation,
		&item.SourceUnsaid,
		&item.TargetPerspective,
		&item.TargetDesireInRelation,
		&item.TargetFearInRelation,
		&item.TargetUnsaid,
		&phenomenologyJSON,
		&symbolicImagesJSON,
		&themeTagsJSON,
		&relationPaletteJSON,
		&relationKeywordsJSON,
		&tensionTagsJSON,
		&item.CoverPrompt,
		&item.SongPrompt,
		&item.PrimarySongSlug,
	)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	item.Phenomenology = adminPhenomenologyFromRelation(parseRelationPhenomenology(phenomenologyJSON))
	item.SymbolicImages = parseJSONStringArray(symbolicImagesJSON)
	item.ThemeTags = parseJSONStringArray(themeTagsJSON)
	item.RelationPalette = adminPaletteFromRelation(parseRelationPalette(relationPaletteJSON))
	item.RelationKeywords = parseJSONStringArray(relationKeywordsJSON)
	item.TensionTags = parseJSONStringArray(tensionTagsJSON)
	return item, nil
}

func (r *postgresAdminRepo) populateAdminRelationChildren(ctx context.Context, item *dto.AdminRelation) error {
	themeSlugs, err := listAdminRelationThemes(ctx, r.pool, item.Slug)
	if err != nil {
		return err
	}
	events, err := listAdminRelationEvents(ctx, r.pool, item.Slug)
	if err != nil {
		return err
	}
	songs, err := listAdminRelationSongs(ctx, r.pool, item.Slug)
	if err != nil {
		return err
	}
	links, err := listAdminRelationLinks(ctx, r.pool, item.Slug)
	if err != nil {
		return err
	}
	item.ThemeSlugs = themeSlugs
	item.Events = events
	item.Songs = songs
	item.Links = links
	return nil
}

func listAdminRelationThemes(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, relationSlug string) ([]string, error) {
	rows, err := q.Query(ctx, `
SELECT theme_slug
FROM public.pm_relation_themes
WHERE relation_slug = $1
ORDER BY is_primary DESC, sort_order ASC, theme_slug ASC
`, relationSlug)
	if err != nil {
		return nil, fmt.Errorf("list relation themes: %w", err)
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		out = append(out, slug)
	}
	return out, rows.Err()
}

func listAdminRelationEvents(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, relationSlug string) ([]dto.AdminRelationEvent, error) {
	rows, err := q.Query(ctx, `
SELECT stage_no, COALESCE(stage_code,''), title, COALESCE(summary,''), COALESCE(tension_shift,''), COALESCE(power_shift,''), COALESCE(fate_impact,''), COALESCE(source_state,''), COALESCE(target_state,''), COALESCE(event_quote,''), COALESCE(color_hex,''), COALESCE(sort_order,0)
FROM public.pm_relation_events
WHERE relation_slug = $1
ORDER BY sort_order ASC, stage_no ASC
`, relationSlug)
	if err != nil {
		return nil, fmt.Errorf("list relation events: %w", err)
	}
	defer rows.Close()
	out := make([]dto.AdminRelationEvent, 0)
	for rows.Next() {
		var item dto.AdminRelationEvent
		if err := rows.Scan(&item.StageNo, &item.StageCode, &item.Title, &item.Summary, &item.TensionShift, &item.PowerShift, &item.FateImpact, &item.SourceState, &item.TargetState, &item.EventQuote, &item.ColorHex, &item.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func listAdminRelationSongs(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, relationSlug string) ([]dto.AdminRelationSong, error) {
	rows, err := q.Query(ctx, `
SELECT slug, title, COALESCE(subtitle,''), COALESCE(summary,''), COALESCE(cover_url,''), COALESCE(audio_url,''), COALESCE(song_core_theme,''), COALESCE(song_emotional_curve,''), COALESCE(song_styles, '[]'::jsonb)::text, COALESCE(tempo_bpm,0), COALESCE(vocal_profile,''), COALESCE(lyric,''), COALESCE(prompt,''), is_primary, COALESCE(sort_order,0), status
FROM public.pm_relation_songs
WHERE relation_slug = $1
ORDER BY CASE WHEN is_primary THEN 0 ELSE 1 END, sort_order ASC, slug ASC
`, relationSlug)
	if err != nil {
		return nil, fmt.Errorf("list relation songs: %w", err)
	}
	defer rows.Close()
	out := make([]dto.AdminRelationSong, 0)
	for rows.Next() {
		var item dto.AdminRelationSong
		var stylesJSON string
		if err := rows.Scan(&item.Slug, &item.Title, &item.Subtitle, &item.Summary, &item.CoverURL, &item.AudioURL, &item.SongCoreTheme, &item.SongEmotionalCurve, &stylesJSON, &item.TempoBPM, &item.VocalProfile, &item.Lyric, &item.Prompt, &item.IsPrimary, &item.SortOrder, &item.Status); err != nil {
			return nil, err
		}
		item.SongStyles = parseJSONStringArray(stylesJSON)
		out = append(out, item)
	}
	return out, rows.Err()
}

func listAdminRelationLinks(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, relationSlug string) ([]dto.AdminRelationLink, error) {
	rows, err := q.Query(ctx, `
SELECT linked_relation_slug, link_type_code, COALESCE(reason,''), COALESCE(sort_order,0)
FROM public.pm_relation_links
WHERE relation_slug = $1
ORDER BY sort_order ASC, linked_relation_slug ASC
`, relationSlug)
	if err != nil {
		return nil, fmt.Errorf("list relation links: %w", err)
	}
	defer rows.Close()
	out := make([]dto.AdminRelationLink, 0)
	for rows.Next() {
		var item dto.AdminRelationLink
		if err := rows.Scan(&item.LinkedRelationSlug, &item.LinkTypeCode, &item.Reason, &item.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func replaceRelationChildren(ctx context.Context, tx pgx.Tx, relationSlug string, in dto.AdminRelation) error {
	if err := replaceRelationThemes(ctx, tx, relationSlug, in.ThemeSlugs); err != nil {
		return err
	}
	if err := replaceRelationEvents(ctx, tx, relationSlug, in.Events); err != nil {
		return err
	}
	if err := replaceRelationSongs(ctx, tx, relationSlug, in.Songs); err != nil {
		return err
	}
	if err := replaceRelationLinks(ctx, tx, relationSlug, in.Links); err != nil {
		return err
	}
	return nil
}

func replaceRelationThemes(ctx context.Context, tx pgx.Tx, relationSlug string, themeSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_themes WHERE relation_slug = $1`, relationSlug); err != nil {
		return fmt.Errorf("clear relation themes: %w", err)
	}
	for idx, themeSlug := range uniq(themeSlugs) {
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_themes (relation_slug, theme_slug, is_primary, sort_order)
VALUES ($1, $2, $3, $4)
`, relationSlug, themeSlug, idx == 0, idx); err != nil {
			return fmt.Errorf("insert relation theme %s: %w", themeSlug, err)
		}
	}
	return nil
}

func replaceRelationEvents(ctx context.Context, tx pgx.Tx, relationSlug string, items []dto.AdminRelationEvent) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_events WHERE relation_slug = $1`, relationSlug); err != nil {
		return fmt.Errorf("clear relation events: %w", err)
	}
	for idx, item := range items {
		stageNo := item.StageNo
		if stageNo <= 0 {
			stageNo = idx + 1
		}
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_events (
  relation_slug, stage_no, stage_code, title, summary, tension_shift, power_shift, fate_impact, source_state, target_state, event_quote, color_hex, sort_order
) VALUES (
  $1,$2,NULLIF($3,''),$4,NULLIF($5,''),NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),NULLIF($12,''),$13
)
`, relationSlug, stageNo, item.StageCode, item.Title, item.Summary, item.TensionShift, item.PowerShift, item.FateImpact, item.SourceState, item.TargetState, item.EventQuote, item.ColorHex, sortOrder); err != nil {
			return fmt.Errorf("insert relation event %d: %w", stageNo, err)
		}
	}
	return nil
}

func replaceRelationSongs(ctx context.Context, tx pgx.Tx, relationSlug string, items []dto.AdminRelationSong) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_songs WHERE relation_slug = $1`, relationSlug); err != nil {
		return fmt.Errorf("clear relation songs: %w", err)
	}
	for idx, item := range items {
		stylesJSON, _ := json.Marshal(uniq(item.SongStyles))
		status := normalizeAdminStatus(item.Status, "draft")
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_songs (
  slug, relation_slug, title, subtitle, summary, cover_url, audio_url, song_core_theme, song_emotional_curve, song_styles, tempo_bpm, vocal_profile, lyric, prompt, is_primary, sort_order, status, is_active
) VALUES (
  $1,$2,$3,NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),COALESCE($10::jsonb, '[]'::jsonb),NULLIF($11,0),NULLIF($12,''),NULLIF($13,''),NULLIF($14,''),$15,$16,$17,$18
)
`, item.Slug, relationSlug, item.Title, item.Subtitle, item.Summary, item.CoverURL, item.AudioURL, item.SongCoreTheme, item.SongEmotionalCurve, string(stylesJSON), item.TempoBPM, item.VocalProfile, item.Lyric, item.Prompt, item.IsPrimary, sortOrder, status, status != "archived"); err != nil {
			return fmt.Errorf("insert relation song %s: %w", item.Slug, err)
		}
	}
	return nil
}

func replaceRelationLinks(ctx context.Context, tx pgx.Tx, relationSlug string, items []dto.AdminRelationLink) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_links WHERE relation_slug = $1`, relationSlug); err != nil {
		return fmt.Errorf("clear relation links: %w", err)
	}
	for idx, item := range items {
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_links (relation_slug, linked_relation_slug, link_type_code, sort_order, reason)
VALUES ($1,$2,$3,$4,NULLIF($5,''))
`, relationSlug, item.LinkedRelationSlug, item.LinkTypeCode, sortOrder, item.Reason); err != nil {
			return fmt.Errorf("insert relation link %s: %w", item.LinkedRelationSlug, err)
		}
	}
	return nil
}

func ensureRelationRefs(ctx context.Context, tx pgx.Tx, in dto.AdminRelation) error {
	if err := ensureRelationTypeExists(ctx, tx, in.RelationTypeCode); err != nil {
		return err
	}
	if strings.TrimSpace(in.SourceCharacterSlug) == "" {
		return errors.New("source character is required")
	}
	if strings.TrimSpace(in.TargetCharacterSlug) == "" {
		return errors.New("target character is required")
	}
	if in.SourceCharacterSlug == in.TargetCharacterSlug {
		return errors.New("source and target characters must be different")
	}
	if err := ensureSlugExists(ctx, tx, "public.pm_characters", in.SourceCharacterSlug, "source character"); err != nil {
		return err
	}
	if err := ensureSlugExists(ctx, tx, "public.pm_characters", in.TargetCharacterSlug, "target character"); err != nil {
		return err
	}
	if strings.TrimSpace(in.WorkSlug) != "" {
		if err := ensureSlugExists(ctx, tx, "public.pm_works", in.WorkSlug, "work"); err != nil {
			return err
		}
	}
	for _, themeSlug := range uniq(in.ThemeSlugs) {
		if err := ensureSlugExists(ctx, tx, "public.pm_themes", themeSlug, "theme"); err != nil {
			return err
		}
	}
	for _, link := range in.Links {
		if strings.TrimSpace(link.LinkedRelationSlug) == "" || link.LinkedRelationSlug == in.Slug {
			continue
		}
		if err := ensureOptionalRelationExists(ctx, tx, link.LinkedRelationSlug); err != nil {
			return err
		}
	}
	return nil
}

func ensureRelationTypeExists(ctx context.Context, tx pgx.Tx, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.New("relation_type_code is required")
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_types (code, name, sort_order, is_active)
VALUES ($1, $2, 0, TRUE)
ON CONFLICT (code) DO NOTHING
`, code, humanizeCode(code)); err != nil {
		return fmt.Errorf("ensure relation type %s: %w", code, err)
	}
	return nil
}

func ensureSlugExists(ctx context.Context, tx pgx.Tx, table, slug, label string) error {
	var exists bool
	err := tx.QueryRow(ctx, fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE slug = $1)`, table), strings.TrimSpace(slug)).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check %s ref: %w", label, err)
	}
	if !exists {
		return fmt.Errorf("%s %s not found", label, slug)
	}
	return nil
}

func ensureOptionalRelationExists(ctx context.Context, tx pgx.Tx, slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil
	}
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM public.pm_relations WHERE slug = $1)`, slug).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check linked relation ref: %w", err)
	}
	if !exists {
		return fmt.Errorf("linked relation %s not found", slug)
	}
	return nil
}

func normalizeAdminRelationInput(in dto.AdminRelation) dto.AdminRelation {
	in.Slug = strings.TrimSpace(in.Slug)
	in.Name = strings.TrimSpace(in.Name)
	in.RelationTypeCode = strings.TrimSpace(in.RelationTypeCode)
	in.WorkSlug = strings.TrimSpace(in.WorkSlug)
	in.SourceCharacterSlug = strings.TrimSpace(in.SourceCharacterSlug)
	in.TargetCharacterSlug = strings.TrimSpace(in.TargetCharacterSlug)
	in.CoreTension = strings.TrimSpace(in.CoreTension)
	in.Status = normalizeAdminStatus(in.Status, "draft")
	if in.SourceCharacterSlug == in.TargetCharacterSlug && in.SourceCharacterSlug != "" {
		in.TargetCharacterSlug = ""
	}
	return in
}

func buildAdminRelationMeta(in dto.AdminRelation) map[string]any {
	return map[string]any{
		"relation_keywords": uniq(in.RelationKeywords),
	}
}

func pickPrimaryRelationSongSlug(in dto.AdminRelation) string {
	if strings.TrimSpace(in.PrimarySongSlug) != "" {
		return strings.TrimSpace(in.PrimarySongSlug)
	}
	for _, item := range in.Songs {
		if item.IsPrimary && strings.TrimSpace(item.Slug) != "" {
			return strings.TrimSpace(item.Slug)
		}
	}
	if len(in.Songs) > 0 {
		return strings.TrimSpace(in.Songs[0].Slug)
	}
	return ""
}

func cascadeRelationSlug(ctx context.Context, tx pgx.Tx, oldSlug, newSlug string) error {
	for _, statement := range []string{
		`UPDATE public.pm_relation_events SET relation_slug = $2 WHERE relation_slug = $1`,
		`UPDATE public.pm_relation_songs SET relation_slug = $2 WHERE relation_slug = $1`,
		`UPDATE public.pm_relation_themes SET relation_slug = $2 WHERE relation_slug = $1`,
		`UPDATE public.pm_relation_links SET relation_slug = $2 WHERE relation_slug = $1`,
		`UPDATE public.pm_relation_links SET linked_relation_slug = $2 WHERE linked_relation_slug = $1`,
		`UPDATE public.pm_relation_participants SET relation_slug = $2 WHERE relation_slug = $1`,
		`UPDATE public.pm_relation_works SET relation_slug = $2 WHERE relation_slug = $1`,
	} {
		if _, err := tx.Exec(ctx, statement, oldSlug, newSlug); err != nil {
			return fmt.Errorf("cascade relation slug %s -> %s: %w", oldSlug, newSlug, err)
		}
	}
	return nil
}

func adminPhenomenologyFromRelation(in dto.RelationPhenomenology) dto.AdminRelationPhenomenology {
	return dto.AdminRelationPhenomenology{
		Body:     in.Body,
		Time:     in.Time,
		Space:    in.Space,
		Gaze:     in.Gaze,
		Language: in.Language,
	}
}

func adminPaletteFromRelation(in []dto.RelationPaletteItem) []dto.AdminRelationPaletteItem {
	out := make([]dto.AdminRelationPaletteItem, 0, len(in))
	for _, item := range in {
		out = append(out, dto.AdminRelationPaletteItem{Name: item.Name, Hex: item.Hex})
	}
	return out
}

func humanizeCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	parts := strings.Fields(strings.NewReplacer("_", " ", "-", " ").Replace(code))
	for i := range parts {
		if len(parts[i]) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, " ")
}
