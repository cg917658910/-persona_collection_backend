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

const relationBaseSelect = `
SELECT
  r.slug,
  COALESCE(NULLIF(r.name, ''), CONCAT(COALESCE(cs.name, r.source_character_slug), ' × ', COALESCE(ct.name, r.target_character_slug))) AS name,
  COALESCE(r.subtitle, '') AS subtitle,
  COALESCE(r.summary, '') AS summary,
  COALESCE(r.one_line_definition, '') AS one_line_definition,
  COALESCE(r.cover_url, cs.cover_url, ct.cover_url, '') AS cover_url,
  COALESCE(r.relation_type_code, '') AS relation_type_code,
  COALESCE(rt.name, '') AS relation_type_name,
  COALESCE(r.work_slug, '') AS work_slug,
  COALESCE(w.title, '') AS work_name,
  COALESCE(r.core_tension, '') AS core_tension,
  COALESCE(r.emotional_tone, '') AS emotional_tone,
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
  COALESCE(r.meta->'relation_keywords', r.tension_tags, '[]'::jsonb)::text AS relation_keywords_json,
  COALESCE(cs.slug, r.source_character_slug) AS source_slug,
  COALESCE(cs.name, r.source_character_slug) AS source_name,
  COALESCE(cs.cover_url, '') AS source_cover_url,
  COALESCE(cs.summary, '') AS source_summary,
  COALESCE(ct.slug, r.target_character_slug) AS target_slug,
  COALESCE(ct.name, r.target_character_slug) AS target_name,
  COALESCE(ct.cover_url, '') AS target_cover_url,
  COALESCE(ct.summary, '') AS target_summary
FROM public.pm_relations r
LEFT JOIN public.pm_relation_types rt ON rt.code = r.relation_type_code
LEFT JOIN public.pm_characters cs ON cs.slug = r.source_character_slug AND cs.is_active = TRUE AND cs.status = 'published'
LEFT JOIN public.pm_characters ct ON ct.slug = r.target_character_slug AND ct.is_active = TRUE AND ct.status = 'published'
LEFT JOIN public.pm_works w ON w.slug = r.work_slug AND w.is_active = TRUE
WHERE r.is_active = TRUE
  AND r.status = 'published'
`

func (r *postgresCatalogRepo) ListRelationships(characterSlug string) ([]dto.RelationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := relationBaseSelect
	args := make([]any, 0, 1)
	if strings.TrimSpace(characterSlug) != "" {
		query += " AND (r.source_character_slug = $1 OR r.target_character_slug = $1)"
		args = append(args, strings.TrimSpace(characterSlug))
	}
	query += " ORDER BY r.sort_order ASC, r.created_at DESC, r.slug ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list relations query: %w", err)
	}
	defer rows.Close()

	return scanRelationRows(rows)
}

func (r *postgresCatalogRepo) GetRelationshipDetail(slug string) (dto.RelationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := relationBaseSelect + `
  AND r.slug = $1
ORDER BY r.sort_order ASC, r.created_at DESC
LIMIT 1
`

	var relation dto.RelationRecord
	if err := scanRelationRow(r.pool.QueryRow(ctx, query, strings.TrimSpace(slug)), &relation); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.RelationRecord{}, errors.New("relationship not found")
		}
		return dto.RelationRecord{}, fmt.Errorf("get relationship detail query: %w", err)
	}

	events, err := r.listRelationEvents(ctx, relation.Slug)
	if err != nil {
		return dto.RelationRecord{}, err
	}
	relation.Events = events

	song, err := r.getPrimaryRelationSong(ctx, relation.Slug)
	if err != nil {
		return dto.RelationRecord{}, err
	}
	relation.PrimarySong = song

	links, err := r.listRelatedRelations(ctx, relation.Slug)
	if err != nil {
		return dto.RelationRecord{}, err
	}
	relation.RelatedRelations = links

	return relation, nil
}

func (r *postgresCatalogRepo) listRelationsByThemeSlug(ctx context.Context, themeSlug string) ([]dto.RelationRecord, error) {
	query := relationBaseSelect + `
  AND EXISTS (
    SELECT 1
    FROM public.pm_relation_themes x
    WHERE x.relation_slug = r.slug
      AND x.theme_slug = $1
  )
ORDER BY r.sort_order ASC, r.created_at DESC, r.slug ASC
`

	rows, err := r.pool.Query(ctx, query, strings.TrimSpace(themeSlug))
	if err != nil {
		return nil, fmt.Errorf("list relations by theme query: %w", err)
	}
	defer rows.Close()

	return scanRelationRows(rows)
}

func (r *postgresCatalogRepo) listRelationEvents(ctx context.Context, slug string) ([]dto.RelationEvent, error) {
	query := `
SELECT
  stage_no,
  COALESCE(stage_code, '') AS stage_code,
  title,
  COALESCE(summary, '') AS summary,
  COALESCE(tension_shift, '') AS tension_shift,
  COALESCE(power_shift, '') AS power_shift,
  COALESCE(fate_impact, '') AS fate_impact,
  COALESCE(source_state, '') AS source_state,
  COALESCE(target_state, '') AS target_state,
  COALESCE(event_quote, '') AS event_quote,
  COALESCE(color_hex, '') AS color_hex
FROM public.pm_relation_events
WHERE relation_slug = $1
ORDER BY sort_order ASC, stage_no ASC
`
	rows, err := r.pool.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("list relation events query: %w", err)
	}
	defer rows.Close()

	events := make([]dto.RelationEvent, 0)
	for rows.Next() {
		var item dto.RelationEvent
		if err := rows.Scan(
			&item.StageNo,
			&item.StageCode,
			&item.Title,
			&item.Summary,
			&item.TensionShift,
			&item.PowerShift,
			&item.FateImpact,
			&item.SourceState,
			&item.TargetState,
			&item.EventQuote,
			&item.ColorHex,
		); err != nil {
			return nil, fmt.Errorf("scan relation event: %w", err)
		}
		events = append(events, item)
	}
	return events, rows.Err()
}

func (r *postgresCatalogRepo) getPrimaryRelationSong(ctx context.Context, relationSlug string) (*dto.RelationSong, error) {
	query := `
SELECT
  slug,
  title,
  COALESCE(subtitle, '') AS subtitle,
  COALESCE(summary, '') AS summary,
  COALESCE(cover_url, '') AS cover_url,
  COALESCE(audio_url, '') AS audio_url,
  COALESCE(song_core_theme, '') AS song_core_theme,
  COALESCE(song_emotional_curve, '') AS song_emotional_curve,
  COALESCE(song_styles, '[]'::jsonb)::text AS song_styles_json,
  COALESCE(vocal_profile, '') AS vocal_profile,
  COALESCE(lyric, '') AS lyric
FROM public.pm_relation_songs
WHERE relation_slug = $1
  AND is_active = TRUE
  AND status = 'published'
ORDER BY CASE WHEN is_primary THEN 0 ELSE 1 END, sort_order ASC, created_at ASC
LIMIT 1
`

	var item dto.RelationSong
	var songStylesJSON string
	err := r.pool.QueryRow(ctx, query, relationSlug).Scan(
		&item.Slug,
		&item.Title,
		&item.Subtitle,
		&item.Summary,
		&item.CoverURL,
		&item.AudioURL,
		&item.SongCoreTheme,
		&item.SongEmotionalCurve,
		&songStylesJSON,
		&item.VocalProfile,
		&item.Lyric,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get primary relation song query: %w", err)
	}
	item.SongStyles = parseJSONStringArray(songStylesJSON)
	return &item, nil
}

func (r *postgresCatalogRepo) listRelatedRelations(ctx context.Context, relationSlug string) ([]dto.RelationLink, error) {
	query := `
SELECT
  rl.link_type_code,
  COALESCE(rl.reason, '') AS reason,
  related.slug,
  COALESCE(NULLIF(related.name, ''), CONCAT(COALESCE(cs.name, related.source_character_slug), ' × ', COALESCE(ct.name, related.target_character_slug))) AS title,
  COALESCE(NULLIF(related.subtitle, ''), COALESCE(w.title, '')) AS subtitle,
  COALESCE(related.cover_url, cs.cover_url, ct.cover_url, '') AS cover_url
FROM public.pm_relation_links rl
JOIN public.pm_relations related ON related.slug = rl.linked_relation_slug
LEFT JOIN public.pm_characters cs ON cs.slug = related.source_character_slug AND cs.is_active = TRUE AND cs.status = 'published'
LEFT JOIN public.pm_characters ct ON ct.slug = related.target_character_slug AND ct.is_active = TRUE AND ct.status = 'published'
LEFT JOIN public.pm_works w ON w.slug = related.work_slug AND w.is_active = TRUE
WHERE rl.relation_slug = $1
  AND related.is_active = TRUE
  AND related.status = 'published'
ORDER BY rl.sort_order ASC, related.sort_order ASC, related.created_at DESC
`

	rows, err := r.pool.Query(ctx, query, relationSlug)
	if err != nil {
		return nil, fmt.Errorf("list related relations query: %w", err)
	}
	defer rows.Close()

	out := make([]dto.RelationLink, 0)
	for rows.Next() {
		var item dto.RelationLink
		if err := rows.Scan(
			&item.LinkTypeCode,
			&item.Reason,
			&item.Slug,
			&item.Title,
			&item.Subtitle,
			&item.CoverURL,
		); err != nil {
			return nil, fmt.Errorf("scan related relation: %w", err)
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func scanRelationRows(rows pgx.Rows) ([]dto.RelationRecord, error) {
	list := make([]dto.RelationRecord, 0)
	for rows.Next() {
		var item dto.RelationRecord
		if err := scanRelationRow(rows, &item); err != nil {
			return nil, fmt.Errorf("scan relation: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanRelationRow(row interface {
	Scan(dest ...any) error
}, relation *dto.RelationRecord) error {
	var phenomenologyJSON string
	var symbolicImagesJSON string
	var themeTagsJSON string
	var relationPaletteJSON string
	var relationKeywordsJSON string

	if err := row.Scan(
		&relation.Slug,
		&relation.Name,
		&relation.Subtitle,
		&relation.Summary,
		&relation.OneLineDefinition,
		&relation.CoverURL,
		&relation.RelationTypeCode,
		&relation.RelationTypeName,
		&relation.WorkSlug,
		&relation.WorkName,
		&relation.CoreTension,
		&relation.EmotionalTone,
		&relation.ConnectionTrigger,
		&relation.SustainingMechanism,
		&relation.RelationConflict,
		&relation.RelationArc,
		&relation.FateImpact,
		&relation.PowerStructure,
		&relation.DependencyPattern,
		&relation.SourcePerspective,
		&relation.SourceDesireInRelation,
		&relation.SourceFearInRelation,
		&relation.SourceUnsaid,
		&relation.TargetPerspective,
		&relation.TargetDesireInRelation,
		&relation.TargetFearInRelation,
		&relation.TargetUnsaid,
		&phenomenologyJSON,
		&symbolicImagesJSON,
		&themeTagsJSON,
		&relationPaletteJSON,
		&relationKeywordsJSON,
		&relation.SourceCharacter.Slug,
		&relation.SourceCharacter.Name,
		&relation.SourceCharacter.CoverURL,
		&relation.SourceCharacter.Summary,
		&relation.TargetCharacter.Slug,
		&relation.TargetCharacter.Name,
		&relation.TargetCharacter.CoverURL,
		&relation.TargetCharacter.Summary,
	); err != nil {
		return err
	}

	relation.Phenomenology = parseRelationPhenomenology(phenomenologyJSON)
	relation.SymbolicImages = parseJSONStringArray(symbolicImagesJSON)
	relation.ThemeTags = parseJSONStringArray(themeTagsJSON)
	relation.RelationPalette = parseRelationPalette(relationPaletteJSON)
	relation.RelationKeywords = parseJSONStringArray(relationKeywordsJSON)
	return nil
}

func parseJSONStringArray(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return []string{}
	}

	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func parseRelationPhenomenology(raw string) dto.RelationPhenomenology {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return dto.RelationPhenomenology{}
	}

	var item dto.RelationPhenomenology
	if err := json.Unmarshal([]byte(raw), &item); err == nil {
		return item
	}

	return dto.RelationPhenomenology{}
}

func parseRelationPalette(raw string) []dto.RelationPaletteItem {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return []dto.RelationPaletteItem{}
	}

	var palette []dto.RelationPaletteItem
	if err := json.Unmarshal([]byte(raw), &palette); err == nil {
		return compactRelationPalette(palette)
	}

	var colors []string
	if err := json.Unmarshal([]byte(raw), &colors); err == nil {
		out := make([]dto.RelationPaletteItem, 0, len(colors))
		for _, color := range colors {
			color = strings.TrimSpace(color)
			if color != "" {
				out = append(out, dto.RelationPaletteItem{Hex: color})
			}
		}
		return out
	}

	return []dto.RelationPaletteItem{}
}

func compactRelationPalette(in []dto.RelationPaletteItem) []dto.RelationPaletteItem {
	out := make([]dto.RelationPaletteItem, 0, len(in))
	for _, item := range in {
		if strings.TrimSpace(item.Hex) == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}
