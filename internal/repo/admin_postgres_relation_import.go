package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
)

func marshalJSONObject(value map[string]any) string {
	if len(value) == 0 {
		return "{}"
	}
	data, err := json.Marshal(value)
	if err != nil || string(data) == "null" || string(data) == "" {
		return "{}"
	}
	return string(data)
}

func (r *postgresAdminRepo) ValidateRelationPackage(pkg dto.RelationImportPackage) (dto.AdminImportResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := newRelationImportResult(pkg)
	errors := make([]string, 0)
	warnings := make([]string, 0)

	existingRelationTypes, err := fetchCodeSet(ctx, r.pool, "public.pm_relation_types")
	if err != nil {
		return result, err
	}
	existingCharacterSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_characters")
	if err != nil {
		return result, err
	}
	existingWorkSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_works")
	if err != nil {
		return result, err
	}
	existingThemeSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_themes")
	if err != nil {
		return result, err
	}
	existingRelationSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_relations")
	if err != nil {
		return result, err
	}

	packageRelationSlugs := map[string]struct{}{}
	for _, raw := range pkg.PmRelations {
		rel := normalizeGeneratedRelation(raw)
		if strings.TrimSpace(rel.Slug) == "" {
			errors = append(errors, "relation requires slug")
			continue
		}
		if _, ok := packageRelationSlugs[rel.Slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate relation slug %s", rel.Slug))
		}
		packageRelationSlugs[rel.Slug] = struct{}{}
	}
	for _, raw := range pkg.PmRelations {
		rel := normalizeGeneratedRelation(raw)
		label := "relation " + rel.Slug
		if strings.TrimSpace(rel.Slug) == "" {
			continue
		}
		if strings.TrimSpace(rel.Name) == "" {
			errors = append(errors, fmt.Sprintf("%s requires name", label))
		}
		if strings.TrimSpace(rel.RelationTypeCode) == "" {
			errors = append(errors, fmt.Sprintf("%s requires relation_type_code", label))
		} else if !hasCode(existingRelationTypes, rel.RelationTypeCode) {
			warnings = append(warnings, fmt.Sprintf("%s relation_type_code %s will be auto-created", label, rel.RelationTypeCode))
		}
		if !hasSlug(nil, existingCharacterSlugs, rel.SourceCharacterSlug) {
			errors = append(errors, fmt.Sprintf("%s references missing source_character_slug %s", label, rel.SourceCharacterSlug))
		}
		if !hasSlug(nil, existingCharacterSlugs, rel.TargetCharacterSlug) {
			errors = append(errors, fmt.Sprintf("%s references missing target_character_slug %s", label, rel.TargetCharacterSlug))
		}
		if rel.SourceCharacterSlug == rel.TargetCharacterSlug && rel.SourceCharacterSlug != "" {
			errors = append(errors, fmt.Sprintf("%s source_character_slug cannot equal target_character_slug", label))
		}
		if rel.WorkSlug != "" && !hasSlug(nil, existingWorkSlugs, rel.WorkSlug) {
			errors = append(errors, fmt.Sprintf("%s references missing work_slug %s", label, rel.WorkSlug))
		}
		for _, themeSlug := range relationThemeSlugs(rel) {
			if !hasSlug(nil, existingThemeSlugs, themeSlug) {
				errors = append(errors, fmt.Sprintf("%s references missing theme_slug %s", label, themeSlug))
			}
		}
		errors = append(errors, validateStatus(label, defaultStatus(rel.Status, "published"))...)
		for _, participant := range rel.Participants {
			if !hasSlug(nil, existingCharacterSlugs, participant.CharacterSlug) {
				errors = append(errors, fmt.Sprintf("%s participant references missing character %s", label, participant.CharacterSlug))
			}
		}
		for _, song := range relationSongs(rel) {
			if strings.TrimSpace(song.Slug) == "" {
				errors = append(errors, fmt.Sprintf("%s has song without slug", label))
			}
			if strings.TrimSpace(song.Title) == "" {
				errors = append(errors, fmt.Sprintf("%s has song %s without title", label, song.Slug))
			}
			errors = append(errors, validateStatus(label+" song "+song.Slug, defaultStatus(song.Status, "published"))...)
		}
		for _, linkedSlug := range relationLinkedSlugs(rel) {
			if linkedSlug == rel.Slug {
				continue
			}
			if !hasSlug(packageRelationSlugs, existingRelationSlugs, linkedSlug) {
				errors = append(errors, fmt.Sprintf("%s references missing linked relation %s", label, linkedSlug))
			}
		}
	}

	result.Warnings = dedupeStrings(warnings)
	result.Errors = dedupeStrings(errors)
	result.Valid = len(result.Errors) == 0
	return result, nil
}

func (r *postgresAdminRepo) ImportRelationPackage(pkg dto.RelationImportPackage) (dto.AdminImportResult, error) {
	result, err := r.ValidateRelationPackage(pkg)
	if err != nil {
		return result, err
	}
	if !result.Valid {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)

	for _, raw := range pkg.PmRelations {
		rel := normalizeGeneratedRelation(raw)
		if err := ensureRelationTypeExists(ctx, tx, rel.RelationTypeCode); err != nil {
			return result, err
		}
		if err := upsertGeneratedRelation(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("upsert relation %s: %w", rel.Slug, err)
		}
		if err := replaceGeneratedRelationParticipants(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("sync relation participants %s: %w", rel.Slug, err)
		}
		if err := replaceGeneratedRelationThemes(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("sync relation themes %s: %w", rel.Slug, err)
		}
		if err := replaceGeneratedRelationEvents(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("sync relation events %s: %w", rel.Slug, err)
		}
		if err := replaceGeneratedRelationSongs(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("sync relation songs %s: %w", rel.Slug, err)
		}
		if err := replaceGeneratedRelationLinks(ctx, tx, rel); err != nil {
			return result, fmt.Errorf("sync relation links %s: %w", rel.Slug, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return result, err
	}
	result.Imported = true
	return result, nil
}

func (r *mockAdminRepo) ValidateRelationPackage(pkg dto.RelationImportPackage) (dto.AdminImportResult, error) {
	result := newRelationImportResult(pkg)
	result.Valid = true
	result.Warnings = []string{"mock mode: relation import validation only checks JSON shape"}
	return result, nil
}

func (r *mockAdminRepo) ImportRelationPackage(pkg dto.RelationImportPackage) (dto.AdminImportResult, error) {
	result, err := r.ValidateRelationPackage(pkg)
	if err != nil {
		return result, err
	}
	result.Imported = result.Valid
	if result.Valid {
		result.Warnings = append(result.Warnings, "mock mode: relation package was not written to a database")
	}
	return result, nil
}

func newRelationImportResult(pkg dto.RelationImportPackage) dto.AdminImportResult {
	participants := 0
	events := 0
	songs := 0
	themes := 0
	links := 0
	for _, raw := range pkg.PmRelations {
		rel := normalizeGeneratedRelation(raw)
		participants += len(rel.Participants)
		events += len(rel.Events)
		songs += len(relationSongs(rel))
		themes += len(relationThemeSlugs(rel))
		links += len(relationLinks(rel))
	}
	return dto.AdminImportResult{
		PackageVersion: "pm-relations-v3",
		Summary: dto.AdminImportSummary{
			Relations:            len(pkg.PmRelations),
			RelationParticipants: participants,
			RelationEvents:       events,
			RelationSongs:        songs,
			RelationThemes:       themes,
			RelationLinks:        links,
		},
	}
}

func normalizeGeneratedRelation(raw dto.GeneratedRelation) dto.GeneratedRelation {
	rel := raw
	if strings.TrimSpace(rel.Name) == "" {
		rel.Name = strings.TrimSpace(rel.Title)
	}
	if len(rel.RelationPalette) == 0 && len(rel.Palette) > 0 {
		rel.RelationPalette = rel.Palette
	}
	if strings.TrimSpace(rel.RelationConflict) == "" && strings.TrimSpace(rel.CoreTension) != "" {
		rel.RelationConflict = rel.CoreTension
	}
	if strings.TrimSpace(rel.CoreTension) == "" && strings.TrimSpace(rel.RelationConflict) != "" {
		rel.CoreTension = rel.RelationConflict
	}
	rel.Status = defaultStatus(rel.Status, "published")
	if rel.IsActive == nil {
		active := rel.Status != "archived"
		rel.IsActive = &active
	}
	if rel.Meta == nil {
		rel.Meta = map[string]any{}
	}
	if strings.TrimSpace(rel.EndingDirection) != "" {
		rel.Meta["ending_direction"] = rel.EndingDirection
	}
	if len(rel.HighlightQuotes) > 0 {
		rel.Meta["ui_highlight_quotes"] = dedupeStrings(rel.HighlightQuotes)
	}
	return rel
}

func defaultStatus(status string, fallback string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return fallback
	}
	return status
}

func relationThemeSlugs(rel dto.GeneratedRelation) []string {
	return uniq(append(append([]string{}, rel.ThemeSlugs...), rel.RelationThemes...))
}

func relationSongs(rel dto.GeneratedRelation) []dto.GeneratedRelationSong {
	if len(rel.Songs) > 0 {
		return rel.Songs
	}
	if rel.Song != nil {
		song := *rel.Song
		if strings.TrimSpace(song.Slug) == "" {
			song.Slug = rel.Slug + "-song"
		}
		song.IsPrimary = true
		if song.Status == "" {
			song.Status = "published"
		}
		return []dto.GeneratedRelationSong{song}
	}
	return nil
}

func relationLinks(rel dto.GeneratedRelation) []dto.GeneratedRelationLink {
	links := make([]dto.GeneratedRelationLink, 0, len(rel.Links)+len(rel.RelatedRelationSlugs)+len(rel.MirrorRelationSlugs)+len(rel.SameWorkRelationSlugs))
	links = append(links, rel.Links...)
	for idx, slug := range uniq(rel.RelatedRelationSlugs) {
		links = append(links, dto.GeneratedRelationLink{LinkedRelationSlug: slug, LinkTypeCode: "related", SortOrder: idx + 1})
	}
	for idx, slug := range uniq(rel.MirrorRelationSlugs) {
		links = append(links, dto.GeneratedRelationLink{LinkedRelationSlug: slug, LinkTypeCode: "mirror", SortOrder: idx + 1})
	}
	for idx, slug := range uniq(rel.SameWorkRelationSlugs) {
		links = append(links, dto.GeneratedRelationLink{LinkedRelationSlug: slug, LinkTypeCode: "same_work", SortOrder: idx + 1})
	}
	return links
}

func relationLinkedSlugs(rel dto.GeneratedRelation) []string {
	out := make([]string, 0)
	for _, link := range relationLinks(rel) {
		out = append(out, link.LinkedRelationSlug)
	}
	return uniq(out)
}

func upsertGeneratedRelation(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	phenomenologyJSON, _ := json.Marshal(rel.Phenomenology)
	symbolicImagesJSON, _ := json.Marshal(uniq(rel.SymbolicImages))
	themeTagsJSON, _ := json.Marshal(uniq(rel.ThemeTags))
	paletteJSON, _ := json.Marshal(rel.RelationPalette)
	tensionTagsJSON, _ := json.Marshal(uniq(rel.TensionTags))
	metaJSON := marshalJSONObject(rel.Meta)
	primarySongSlug := strings.TrimSpace(rel.PrimarySongSlug)
	for _, song := range relationSongs(rel) {
		if song.IsPrimary && strings.TrimSpace(song.Slug) != "" {
			primarySongSlug = song.Slug
			break
		}
	}
	_, err := tx.Exec(ctx, `
INSERT INTO public.pm_relations (
  slug, name, subtitle, relation_type_code, work_slug,
  source_character_slug, target_character_slug,
  summary, one_line_definition, core_dynamic, core_tension, emotional_tone, emotional_temperature,
  connection_trigger, sustaining_mechanism, relation_conflict, relation_arc, fate_impact,
  power_structure, dependency_pattern,
  source_perspective, target_perspective, source_desire_in_relation, source_fear_in_relation, source_unsaid,
  target_desire_in_relation, target_fear_in_relation, target_unsaid,
  phenomenology, symbolic_images, theme_tags, relation_palette, tension_tags,
  cover_url, cover_prompt, song_prompt, primary_song_slug, meta, sort_order, status, is_active
) VALUES (
  $1,$2,NULLIF($3,''),$4,NULLIF($5,''),
  $6,$7,
  NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),NULLIF($12,''),NULLIF($13,''),
  NULLIF($14,''),NULLIF($15,''),NULLIF($16,''),NULLIF($17,''),NULLIF($18,''),
  NULLIF($19,''),NULLIF($20,''),
  NULLIF($21,''),NULLIF($22,''),NULLIF($23,''),NULLIF($24,''),NULLIF($25,''),
  NULLIF($26,''),NULLIF($27,''),NULLIF($28,''),
  COALESCE($29::jsonb, '{}'::jsonb), COALESCE($30::jsonb, '[]'::jsonb), COALESCE($31::jsonb, '[]'::jsonb), COALESCE($32::jsonb, '[]'::jsonb), COALESCE($33::jsonb, '[]'::jsonb),
  NULLIF($34,''), NULLIF($35,''), NULLIF($36,''), NULLIF($37,''), COALESCE($38::jsonb, '{}'::jsonb), COALESCE($39,0), $40, $41
)
ON CONFLICT (slug) DO UPDATE SET
  name = EXCLUDED.name,
  subtitle = EXCLUDED.subtitle,
  relation_type_code = EXCLUDED.relation_type_code,
  work_slug = EXCLUDED.work_slug,
  source_character_slug = EXCLUDED.source_character_slug,
  target_character_slug = EXCLUDED.target_character_slug,
  summary = EXCLUDED.summary,
  one_line_definition = EXCLUDED.one_line_definition,
  core_dynamic = EXCLUDED.core_dynamic,
  core_tension = EXCLUDED.core_tension,
  emotional_tone = EXCLUDED.emotional_tone,
  emotional_temperature = EXCLUDED.emotional_temperature,
  connection_trigger = EXCLUDED.connection_trigger,
  sustaining_mechanism = EXCLUDED.sustaining_mechanism,
  relation_conflict = EXCLUDED.relation_conflict,
  relation_arc = EXCLUDED.relation_arc,
  fate_impact = EXCLUDED.fate_impact,
  power_structure = EXCLUDED.power_structure,
  dependency_pattern = EXCLUDED.dependency_pattern,
  source_perspective = EXCLUDED.source_perspective,
  target_perspective = EXCLUDED.target_perspective,
  source_desire_in_relation = EXCLUDED.source_desire_in_relation,
  source_fear_in_relation = EXCLUDED.source_fear_in_relation,
  source_unsaid = EXCLUDED.source_unsaid,
  target_desire_in_relation = EXCLUDED.target_desire_in_relation,
  target_fear_in_relation = EXCLUDED.target_fear_in_relation,
  target_unsaid = EXCLUDED.target_unsaid,
  phenomenology = EXCLUDED.phenomenology,
  symbolic_images = EXCLUDED.symbolic_images,
  theme_tags = EXCLUDED.theme_tags,
  relation_palette = EXCLUDED.relation_palette,
  tension_tags = EXCLUDED.tension_tags,
  cover_url = EXCLUDED.cover_url,
  cover_prompt = EXCLUDED.cover_prompt,
  song_prompt = EXCLUDED.song_prompt,
  primary_song_slug = EXCLUDED.primary_song_slug,
  meta = EXCLUDED.meta,
  sort_order = EXCLUDED.sort_order,
  status = EXCLUDED.status,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, rel.Slug, rel.Name, rel.Subtitle, rel.RelationTypeCode, rel.WorkSlug,
		rel.SourceCharacterSlug, rel.TargetCharacterSlug,
		rel.Summary, rel.OneLineDefinition, rel.CoreDynamic, rel.CoreTension, rel.EmotionalTone, rel.EmotionalTemperature,
		rel.ConnectionTrigger, rel.SustainingMechanism, rel.RelationConflict, rel.RelationArc, rel.FateImpact,
		rel.PowerStructure, rel.DependencyPattern,
		rel.SourcePerspective, rel.TargetPerspective, rel.SourceDesireInRelation, rel.SourceFearInRelation, rel.SourceUnsaid,
		rel.TargetDesireInRelation, rel.TargetFearInRelation, rel.TargetUnsaid,
		string(phenomenologyJSON), string(symbolicImagesJSON), string(themeTagsJSON), string(paletteJSON), string(tensionTagsJSON),
		rel.CoverURL, rel.CoverPrompt, rel.SongPrompt, primarySongSlug, metaJSON, rel.SortOrder, normalizeAdminStatus(rel.Status, "published"), *rel.IsActive)
	return err
}

func replaceGeneratedRelationParticipants(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_participants WHERE relation_slug = $1`, rel.Slug); err != nil {
		return err
	}
	for idx, item := range rel.Participants {
		metaJSON := marshalJSONObject(item.Meta)
		roleCode := strings.TrimSpace(item.RoleCode)
		if roleCode == "" {
			roleCode = "participant"
		}
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_participants (
  relation_slug, character_slug, role_code, role_name, perspective_summary, desire_in_relation, fear_in_relation, unsaid, sort_order, meta
) VALUES (
  $1,$2,$3,NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),$9,COALESCE($10::jsonb, '{}'::jsonb)
)
`, rel.Slug, item.CharacterSlug, roleCode, item.RoleName, item.PerspectiveSummary, item.DesireInRelation, item.FearInRelation, item.Unsaid, sortOrder, metaJSON); err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedRelationThemes(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_themes WHERE relation_slug = $1`, rel.Slug); err != nil {
		return err
	}
	for idx, themeSlug := range relationThemeSlugs(rel) {
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_themes (relation_slug, theme_slug, is_primary, sort_order)
VALUES ($1,$2,$3,$4)
`, rel.Slug, themeSlug, idx == 0, idx+1); err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedRelationEvents(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_events WHERE relation_slug = $1`, rel.Slug); err != nil {
		return err
	}
	for idx, item := range rel.Events {
		metaJSON := marshalJSONObject(item.Meta)
		stageNo := item.StageNo
		if stageNo <= 0 {
			stageNo = item.Stage
		}
		if stageNo <= 0 {
			stageNo = idx + 1
		}
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_events (
  relation_slug, stage_no, stage_code, title, summary, tension_shift, power_shift, fate_impact, source_state, target_state, event_quote, color_hex, sort_order, meta
) VALUES (
  $1,$2,NULLIF($3,''),$4,NULLIF($5,''),NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),NULLIF($12,''),$13,COALESCE($14::jsonb, '{}'::jsonb)
)
`, rel.Slug, stageNo, item.StageCode, item.Title, item.Summary, item.TensionShift, item.PowerShift, item.FateImpact, item.SourceState, item.TargetState, item.EventQuote, item.ColorHex, sortOrder, metaJSON); err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedRelationSongs(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_songs WHERE relation_slug = $1`, rel.Slug); err != nil {
		return err
	}
	for idx, item := range relationSongs(rel) {
		stylesJSON, _ := json.Marshal(uniq(item.SongStyles))
		status := normalizeAdminStatus(defaultStatus(item.Status, "published"), "published")
		isActive := status != "archived"
		if item.IsActive != nil {
			isActive = *item.IsActive
		}
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_songs (
  slug, relation_slug, title, subtitle, summary, cover_url, audio_url, duration_sec,
  song_core_theme, song_emotional_curve, song_styles, tempo_bpm, vocal_profile, lyric, prompt,
  is_primary, sort_order, status, is_active
) VALUES (
  $1,$2,$3,NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),NULLIF($7,''),NULLIF($8,0),
  NULLIF($9,''),NULLIF($10,''),COALESCE($11::jsonb, '[]'::jsonb),NULLIF($12,0),NULLIF($13,''),NULLIF($14,''),NULLIF($15,''),
  $16,$17,$18,$19
)
ON CONFLICT (slug) DO UPDATE SET
  relation_slug = EXCLUDED.relation_slug,
  title = EXCLUDED.title,
  subtitle = EXCLUDED.subtitle,
  summary = EXCLUDED.summary,
  cover_url = EXCLUDED.cover_url,
  audio_url = EXCLUDED.audio_url,
  duration_sec = EXCLUDED.duration_sec,
  song_core_theme = EXCLUDED.song_core_theme,
  song_emotional_curve = EXCLUDED.song_emotional_curve,
  song_styles = EXCLUDED.song_styles,
  tempo_bpm = EXCLUDED.tempo_bpm,
  vocal_profile = EXCLUDED.vocal_profile,
  lyric = EXCLUDED.lyric,
  prompt = EXCLUDED.prompt,
  is_primary = EXCLUDED.is_primary,
  sort_order = EXCLUDED.sort_order,
  status = EXCLUDED.status,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, item.Slug, rel.Slug, item.Title, item.Subtitle, item.Summary, item.CoverURL, item.AudioURL, item.DurationSec,
			item.SongCoreTheme, item.SongEmotionalCurve, string(stylesJSON), item.TempoBPM, item.VocalProfile, item.Lyric, item.Prompt,
			item.IsPrimary, sortOrder, status, isActive); err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedRelationLinks(ctx context.Context, tx pgx.Tx, rel dto.GeneratedRelation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_relation_links WHERE relation_slug = $1`, rel.Slug); err != nil {
		return err
	}
	for idx, item := range relationLinks(rel) {
		sortOrder := item.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO public.pm_relation_links (relation_slug, linked_relation_slug, link_type_code, sort_order, reason)
VALUES ($1,$2,$3,$4,NULLIF($5,''))
`, rel.Slug, item.LinkedRelationSlug, item.LinkTypeCode, sortOrder, item.Reason); err != nil {
			return err
		}
	}
	return nil
}
