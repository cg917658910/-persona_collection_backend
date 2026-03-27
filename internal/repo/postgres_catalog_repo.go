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
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.surface_traits, ARRAY[]::text[]) AS surface_traits
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
		&item.SurfaceTraits,
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
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.surface_traits, ARRAY[]::text[]) AS surface_traits
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
		&item.SurfaceTraits,
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
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.surface_traits, ARRAY[]::text[]) AS surface_traits
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
			&item.SurfaceTraits,
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
  COALESCE(s.audio_url, '') AS audio_url,
  COALESCE(s.summary, '') AS summary,
  COALESCE(s.song_core_theme, '') AS song_core_theme,
  COALESCE(s.song_styles, ARRAY[]::text[]) AS song_styles,
  COALESCE(s.song_emotional_curve, ARRAY[]::text[]) AS song_emotional_curve,
  COALESCE(s.vocal_profile, '') AS vocal_profile
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
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL, &item.Summary, &item.SongCoreTheme, &item.Styles, &item.EmotionalCurve, &item.VocalProfile); err != nil {
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
  COALESCE(t.cover_url, '') AS cover_url,
  COALESCE(t.category, '') AS category
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
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.Category); err != nil {
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
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.surface_traits, ARRAY[]::text[]) AS surface_traits
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
			&item.SurfaceTraits,
		); err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}
		if works, err := r.listWorksByCharacterSlug(ctx, item.Slug); err == nil {
			for i, work := range works {
				item.WorkSlugs = append(item.WorkSlugs, work.Slug)
				if i == 0 {
					item.PrimaryWorkTitle = work.Title
				}
			}
		}
		if themes, err := r.listThemesByCharacterSlug(ctx, item.Slug); err == nil {
			for i, theme := range themes {
				item.ThemeSlugs = append(item.ThemeSlugs, theme.Slug)
				if i == 0 {
					item.PrimaryThemeName = theme.Name
				}
			}
		}
		if songs, err := r.listSongsByCharacterSlug(ctx, item.Slug); err == nil {
			for i, song := range songs {
				item.SongSlugs = append(item.SongSlugs, song.Slug)
				if i == 0 {
					item.PrimarySongTitle = song.Title
				}
			}
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) GetCharacterDetail(slug string) (dto.CharacterDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d dto.CharacterDetail
	var rawColorsJSON string
	var profileJSON, timelineJSON string
	err := r.pool.QueryRow(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.core_identity, '') AS core_identity,
  COALESCE(c.public_image, '') AS public_image,
  COALESCE(c.hidden_self, '') AS hidden_self,
  COALESCE((SELECT m.description FROM public.pm_character_motivations cm JOIN public.pm_motivation_dict m ON m.id = cm.motivation_id WHERE cm.character_id = c.id ORDER BY cm.is_primary DESC, cm.weight DESC, cm.created_at ASC LIMIT 1),
           c.psychology->>'primary_motivation',
           c.psychology->>'primaryMotivation',
           c.psychology->>'desire',
           c.psychology->>'pursuit',
           '') AS primary_motivation,
  COALESCE(c.core_fear, '') AS core_fear,
  COALESCE(c.psychological_wound, '') AS psychological_wound,
  COALESCE(c.core_conflict, '') AS core_conflict,
  COALESCE(c.emotional_tone, '') AS emotional_tone,
  COALESCE(c.origin, '') AS origin,
  COALESCE(c.fate_arc, '') AS fate_arc,
  COALESCE(c.ending_state, '') AS ending_state,
  COALESCE(c.surface_traits, ARRAY[]::text[]) AS surface_traits,
  COALESCE(c.deep_traits, ARRAY[]::text[]) AS deep_traits,
  COALESCE(c.dominant_emotions, ARRAY[]::text[]) AS dominant_emotions,
  COALESCE(c.suppressed_emotions, ARRAY[]::text[]) AS suppressed_emotions,
  COALESCE(c.values_tags, ARRAY[]::text[]) AS values_tags,
  COALESCE(c.bottom_lines, ARRAY[]::text[]) AS bottom_lines,
  COALESCE(c.symbolic_images, ARRAY[]::text[]) AS symbolic_images,
  COALESCE(c.colors, '[]'::jsonb)::text AS colors_json,
  COALESCE(c.elements, ARRAY[]::text[]) AS elements,
  COALESCE(c.soundscape_keywords, ARRAY[]::text[]) AS soundscape_keywords,
  COALESCE(c.relationship_profile::text, '{}') AS relationship_profile_json,
  COALESCE(c.timeline::text, '[]') AS timeline_json
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.slug = $1
  AND c.is_active = TRUE
  AND c.status = 'published'
`, slug).Scan(
		&d.ID,
		&d.Slug,
		&d.Name,
		&d.CharacterTypeCode,
		&d.Summary,
		&d.OneLineDefinition,
		&d.CoverURL,
		&d.CoreIdentity,
		&d.PublicImage,
		&d.HiddenSelf,
		&d.PrimaryMotivation,
		&d.CoreFear,
		&d.PsychologicalWound,
		&d.CoreConflict,
		&d.EmotionalTone,
		&d.Origin,
		&d.FateArc,
		&d.EndingState,
		&d.SurfaceTraits,
		&d.DeepTraits,
		&d.DominantEmotions,
		&d.SuppressedEmotions,
		&d.ValuesTags,
		&d.BottomLines,
		&d.SymbolicImages,
		&rawColorsJSON,
		&d.Elements,
		&d.SoundscapeKeywords,
		&profileJSON,
		&timelineJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.CharacterDetail{}, errors.New("character not found")
		}
		return dto.CharacterDetail{}, fmt.Errorf("get character detail query: %w", err)
	}

	d.Colors = parseCharacterColorsJSON(rawColorsJSON)
	_ = json.Unmarshal([]byte(profileJSON), &d.RelationshipProfile)
	if d.RelationshipProfile == nil {
		d.RelationshipProfile = map[string]string{}
	}
	var rawTimeline []map[string]any
	if json.Unmarshal([]byte(timelineJSON), &rawTimeline) == nil {
		d.Timeline = make([]dto.CharacterTimelineItem, 0, len(rawTimeline))
		for idx, item := range rawTimeline {
			year := strings.TrimSpace(fmt.Sprint(item["stage"]))
			if year == "" || year == "<nil>" {
				year = strings.TrimSpace(fmt.Sprint(item["year"]))
			}
			if year == "" || year == "<nil>" {
				year = fmt.Sprintf("%02d", idx+1)
			}
			event := strings.TrimSpace(fmt.Sprint(item["title"]))
			if event == "" || event == "<nil>" {
				event = strings.TrimSpace(fmt.Sprint(item["event"]))
			}
			emotion := strings.TrimSpace(fmt.Sprint(item["summary"]))
			if emotion == "" || emotion == "<nil>" {
				emotion = strings.TrimSpace(fmt.Sprint(item["emotion"]))
			}
			d.Timeline = append(d.Timeline, dto.CharacterTimelineItem{Year: year, Event: event, Emotion: emotion})
		}
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
		if creators, err := r.listCreatorsByWorkSlug(ctx, item.Slug); err == nil {
			for _, creator := range creators {
				item.CreatorSlugs = append(item.CreatorSlugs, creator.Slug)
				item.CreatorNames = append(item.CreatorNames, creator.Name)
			}
		}
		if characters, err := r.listCharactersByWorkSlug(ctx, item.Slug); err == nil {
			for _, character := range characters {
				item.CharacterSlugs = append(item.CharacterSlugs, character.Slug)
			}
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
	if creators, err := r.listCreatorsByWorkSlug(ctx, slug); err == nil {
		for _, creator := range creators {
			item.CreatorSlugs = append(item.CreatorSlugs, creator.Slug)
			item.CreatorNames = append(item.CreatorNames, creator.Name)
		}
	}
	if characters, err := r.listCharactersByWorkSlug(ctx, slug); err == nil {
		for _, character := range characters {
			item.CharacterSlugs = append(item.CharacterSlugs, character.Slug)
		}
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
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(ct.code, '') AS creator_type_code,
  COALESCE(c.era_text, '') AS era_text
FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
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
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.CreatorTypeCode, &item.EraText); err != nil {
			return nil, fmt.Errorf("scan creator: %w", err)
		}
		if works, err := r.listWorksByCreatorSlug(ctx, item.Slug); err == nil {
			for _, work := range works {
				item.WorkSlugs = append(item.WorkSlugs, work.Slug)
			}
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
	if works, err := r.listWorksByCreatorSlug(ctx, slug); err == nil {
		for _, work := range works {
			item.WorkSlugs = append(item.WorkSlugs, work.Slug)
		}
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
  COALESCE(t.cover_url, '') AS cover_url,
  COALESCE(t.category, '') AS category
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
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.Category); err != nil {
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
  COALESCE(t.cover_url, '') AS cover_url,
  COALESCE(t.category, '') AS category
FROM public.pm_themes t
WHERE t.slug = $1
  AND t.is_active = TRUE
`, slug).Scan(&d.Slug, &d.Name, &d.Summary, &d.CoverURL, &d.Category)
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
  COALESCE(s.audio_url, '') AS audio_url,
  COALESCE(s.summary, '') AS summary,
  COALESCE(s.song_core_theme, '') AS song_core_theme,
  COALESCE(s.song_styles, ARRAY[]::text[]) AS song_styles,
  COALESCE(s.song_emotional_curve, ARRAY[]::text[]) AS song_emotional_curve,
  COALESCE(s.vocal_profile, '') AS vocal_profile
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
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL, &item.Summary, &item.SongCoreTheme, &item.Styles, &item.EmotionalCurve, &item.VocalProfile); err != nil {
			return nil, fmt.Errorf("scan song: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) SearchCatalog(keyword string, limit int) (dto.SearchResponseData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 8
	}
	keyword = strings.TrimSpace(keyword)

	out := dto.SearchResponseData{
		Characters: make([]dto.CharacterListItemResponse, 0, limit),
		Works:      make([]dto.WorkListItemResponse, 0, limit),
		Creators:   make([]dto.CreatorListItemResponse, 0, limit),
		Themes:     make([]dto.ThemeListItemResponse, 0, limit),
		Songs:      make([]dto.SongListItemResponse, 0, limit),
	}

	charRows, err := r.pool.Query(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  ct.code AS character_type_code
FROM public.pm_characters c
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE c.is_active = TRUE
  AND c.status = 'published'
  AND (
    $1 = ''
    OR c.name ILIKE '%' || $1 || '%'
    OR COALESCE(c.summary, '') ILIKE '%' || $1 || '%'
    OR COALESCE(c.one_line_definition, '') ILIKE '%' || $1 || '%'
  )
ORDER BY c.sort_order ASC, c.created_at DESC, c.name ASC
LIMIT $2
`, keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, fmt.Errorf("search characters query: %w", err)
	}
	defer charRows.Close()

	for charRows.Next() {
		var item dto.CharacterListItemResponse
		if err := charRows.Scan(&item.ID, &item.Slug, &item.Name, &item.CoverURL, &item.Summary, &item.OneLineDefinition, &item.CharacterTypeCode); err != nil {
			return dto.SearchResponseData{}, fmt.Errorf("scan search character: %w", err)
		}
		if works, err := r.listWorksByCharacterSlug(ctx, item.Slug); err == nil && len(works) > 0 {
			item.WorkTitle = works[0].Title
		}
		if themes, err := r.listThemesByCharacterSlug(ctx, item.Slug); err == nil && len(themes) > 0 {
			item.Tags = []string{themes[0].Name}
		}
		if songs, err := r.listSongsByCharacterSlug(ctx, item.Slug); err == nil && len(songs) > 0 {
			item.HasSong = true
			item.ThemeSongTitle = "人物之歌"
		}
		out.Characters = append(out.Characters, item)
	}
	if err := charRows.Err(); err != nil {
		return dto.SearchResponseData{}, err
	}

	workRows, err := r.pool.Query(ctx, `
SELECT
  w.id::text,
  w.slug,
  w.title,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(w.summary, '') AS summary,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE w.is_active = TRUE
  AND (
    $1 = ''
    OR w.title ILIKE '%' || $1 || '%'
    OR COALESCE(w.summary, '') ILIKE '%' || $1 || '%'
    OR COALESCE(w.original_title, '') ILIKE '%' || $1 || '%'
  )
ORDER BY w.sort_order ASC, w.created_at DESC, w.title ASC
LIMIT $2
`, keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, fmt.Errorf("search works query: %w", err)
	}
	defer workRows.Close()

	for workRows.Next() {
		var item dto.WorkListItemResponse
		if err := workRows.Scan(&item.ID, &item.Slug, &item.Title, &item.CoverURL, &item.Summary, &item.WorkTypeCode); err != nil {
			return dto.SearchResponseData{}, fmt.Errorf("scan search work: %w", err)
		}
		if creators, err := r.listCreatorsByWorkSlug(ctx, item.Slug); err == nil && len(creators) > 0 {
			item.CreatorName = creators[0].Name
		}
		if chars, err := r.listCharactersByWorkSlug(ctx, item.Slug); err == nil {
			item.CharacterCount = len(chars)
		}
		out.Works = append(out.Works, item)
	}
	if err := workRows.Err(); err != nil {
		return dto.SearchResponseData{}, err
	}

	creatorRows, err := r.pool.Query(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.cover_url, '') AS cover_url,
  COALESCE(c.summary, '') AS summary,
  COALESCE(ct.code, '') AS creator_type_code,
  COALESCE(c.era_text, '') AS era_text
FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
WHERE c.is_active = TRUE
  AND (
    $1 = ''
    OR c.name ILIKE '%' || $1 || '%'
    OR COALESCE(c.summary, '') ILIKE '%' || $1 || '%'
    OR COALESCE(c.era_text, '') ILIKE '%' || $1 || '%'
  )
ORDER BY c.sort_order ASC, c.created_at DESC, c.name ASC
LIMIT $2
`, keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, fmt.Errorf("search creators query: %w", err)
	}
	defer creatorRows.Close()

	for creatorRows.Next() {
		var item dto.CreatorListItemResponse
		if err := creatorRows.Scan(&item.ID, &item.Slug, &item.Name, &item.CoverURL, &item.Summary, &item.CreatorTypeCode, &item.EraText); err != nil {
			return dto.SearchResponseData{}, fmt.Errorf("scan search creator: %w", err)
		}
		if works, err := r.listWorksByCreatorSlug(ctx, item.Slug); err == nil {
			item.WorkCount = len(works)
		}
		out.Creators = append(out.Creators, item)
	}
	if err := creatorRows.Err(); err != nil {
		return dto.SearchResponseData{}, err
	}

	themeRows, err := r.pool.Query(ctx, `
SELECT
  t.id::text,
  t.slug,
  t.name_zh AS name,
  COALESCE(t.cover_url, '') AS cover_url,
  COALESCE(t.summary, '') AS summary,
  COALESCE(t.category, '') AS category,
  COUNT(ct.character_id)::INT AS character_count
FROM public.pm_themes t
LEFT JOIN public.pm_character_themes ct ON ct.theme_id = t.id
LEFT JOIN public.pm_characters c ON c.id = ct.character_id AND c.is_active = TRUE AND c.status = 'published'
WHERE t.is_active = TRUE
  AND (
    $1 = ''
    OR t.name_zh ILIKE '%' || $1 || '%'
    OR COALESCE(t.summary, '') ILIKE '%' || $1 || '%'
    OR COALESCE(t.category, '') ILIKE '%' || $1 || '%'
  )
GROUP BY t.id
ORDER BY t.sort_order ASC, t.created_at DESC, t.name_zh ASC
LIMIT $2
`, keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, fmt.Errorf("search themes query: %w", err)
	}
	defer themeRows.Close()

	for themeRows.Next() {
		var item dto.ThemeListItemResponse
		if err := themeRows.Scan(&item.ID, &item.Slug, &item.Name, &item.CoverURL, &item.Summary, &item.Category, &item.CharacterCount); err != nil {
			return dto.SearchResponseData{}, fmt.Errorf("scan search theme: %w", err)
		}
		out.Themes = append(out.Themes, item)
	}
	if err := themeRows.Err(); err != nil {
		return dto.SearchResponseData{}, err
	}

	songRows, err := r.pool.Query(ctx, `
SELECT
  s.id::text,
  s.slug,
  s.title,
  c.slug AS character_slug,
  COALESCE(s.cover_url, '') AS cover_url,
  COALESCE(s.audio_url, '') AS audio_url,
  COALESCE(s.song_styles, ARRAY[]::text[]) AS styles
FROM public.pm_songs s
JOIN public.pm_characters c ON c.id = s.character_id
WHERE s.is_active = TRUE
  AND s.status = 'published'
  AND (
    $1 = ''
    OR s.title ILIKE '%' || $1 || '%'
    OR COALESCE(s.summary, '') ILIKE '%' || $1 || '%'
    OR COALESCE(s.song_core_theme, '') ILIKE '%' || $1 || '%'
  )
ORDER BY s.sort_order ASC, s.created_at DESC, s.title ASC
LIMIT $2
`, keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, fmt.Errorf("search songs query: %w", err)
	}
	defer songRows.Close()

	for songRows.Next() {
		var item dto.SongListItemResponse
		if err := songRows.Scan(&item.ID, &item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL, &item.Styles); err != nil {
			return dto.SearchResponseData{}, fmt.Errorf("scan search song: %w", err)
		}
		out.Songs = append(out.Songs, item)
	}
	if err := songRows.Err(); err != nil {
		return dto.SearchResponseData{}, err
	}

	return out, nil
}

func (r *postgresCatalogRepo) listSongsByCharacterSlug(ctx context.Context, slug string) ([]dto.Song, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  s.slug,
  s.title,
  c.slug AS character_slug,
  COALESCE(s.cover_url, '') AS cover_url,
  COALESCE(s.audio_url, '') AS audio_url,
  COALESCE(s.summary, '') AS summary,
  COALESCE(s.song_core_theme, '') AS song_core_theme,
  COALESCE(s.song_styles, ARRAY[]::text[]) AS song_styles,
  COALESCE(s.song_emotional_curve, ARRAY[]::text[]) AS song_emotional_curve,
  COALESCE(s.vocal_profile, '') AS vocal_profile
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
		if err := rows.Scan(&item.Slug, &item.Title, &item.CharacterSlug, &item.CoverURL, &item.AudioURL, &item.Summary, &item.SongCoreTheme, &item.Styles, &item.EmotionalCurve, &item.VocalProfile); err != nil {
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
  COALESCE(t.cover_url, '') AS cover_url,
  COALESCE(t.category, '') AS category
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
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.Category); err != nil {
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
  COALESCE(cr.cover_url, '') AS cover_url,
  COALESCE(ct.code, '') AS creator_type_code,
  COALESCE(cr.era_text, '') AS era_text
FROM public.pm_character_works cw
JOIN public.pm_characters c ON c.id = cw.character_id
JOIN public.pm_work_creators wc ON wc.work_id = cw.work_id
JOIN public.pm_creators cr ON cr.id = wc.creator_id
LEFT JOIN public.pm_creator_types ct ON ct.id = cr.creator_type_id
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
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.CreatorTypeCode, &item.EraText); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listCreatorsByWorkSlug(ctx context.Context, slug string) ([]dto.Creator, error) {
	rows, err := r.pool.Query(ctx, `
SELECT DISTINCT
  cr.slug,
  cr.name,
  COALESCE(cr.summary, '') AS summary,
  COALESCE(cr.cover_url, '') AS cover_url,
  COALESCE(ct.code, '') AS creator_type_code,
  COALESCE(cr.era_text, '') AS era_text
FROM public.pm_work_creators wc
JOIN public.pm_works w ON w.id = wc.work_id
JOIN public.pm_creators cr ON cr.id = wc.creator_id
LEFT JOIN public.pm_creator_types ct ON ct.id = cr.creator_type_id
WHERE w.slug = $1
  AND cr.is_active = TRUE
ORDER BY wc.sort_order ASC, cr.name ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Creator, 0)
	for rows.Next() {
		var item dto.Creator
		if err := rows.Scan(&item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.CreatorTypeCode, &item.EraText); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listCharactersByWorkSlug(ctx context.Context, slug string) ([]dto.Character, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  c.slug,
  c.name,
  ct.code AS character_type_code,
  COALESCE(c.summary, '') AS summary,
  COALESCE(c.one_line_definition, '') AS one_line_definition,
  COALESCE(c.cover_url, '') AS cover_url
FROM public.pm_character_works cw
JOIN public.pm_works w ON w.id = cw.work_id
JOIN public.pm_characters c ON c.id = cw.character_id
JOIN public.pm_character_types ct ON ct.id = c.character_type_id
WHERE w.slug = $1
  AND c.is_active = TRUE
  AND c.status = 'published'
ORDER BY cw.is_primary DESC, cw.sort_order ASC, c.sort_order ASC, c.name ASC
`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.Character, 0)
	for rows.Next() {
		var item dto.Character
		if err := rows.Scan(&item.Slug, &item.Name, &item.CharacterTypeCode, &item.Summary, &item.OneLineDefinition, &item.CoverURL); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresCatalogRepo) listWorksByCreatorSlug(ctx context.Context, slug string) ([]dto.Work, error) {
	rows, err := r.pool.Query(ctx, `
SELECT
  w.slug,
  w.title,
  COALESCE(w.summary, '') AS summary,
  COALESCE(w.cover_url, '') AS cover_url,
  COALESCE(wt.code, '') AS work_type_code
FROM public.pm_work_creators wc
JOIN public.pm_creators cr ON cr.id = wc.creator_id
JOIN public.pm_works w ON w.id = wc.work_id
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
WHERE cr.slug = $1
  AND w.is_active = TRUE
ORDER BY wc.sort_order ASC, w.sort_order ASC, w.title ASC
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

func parseCharacterColorsJSON(raw string) []dto.CharacterColorItem {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return nil
	}

	var objectItems []dto.CharacterColorItem
	if err := json.Unmarshal([]byte(raw), &objectItems); err == nil {
		out := make([]dto.CharacterColorItem, 0, len(objectItems))
		for _, item := range objectItems {
			item.Name = strings.TrimSpace(item.Name)
			item.Hex = strings.TrimSpace(item.Hex)
			if item.Hex == "" {
				continue
			}
			if item.Name == "" {
				item.Name = item.Hex
			}
			out = append(out, item)
		}
		if len(out) > 0 {
			return out
		}
	}

	var stringItems []string
	if err := json.Unmarshal([]byte(raw), &stringItems); err == nil {
		return parseCharacterColors(stringItems)
	}

	return parseCharacterColors([]string{raw})
}

func parseCharacterColors(values []string) []dto.CharacterColorItem {
	out := make([]dto.CharacterColorItem, 0, len(values))
	for _, raw := range values {
		raw = strings.TrimSpace(raw)
		raw = strings.Trim(raw, "\"")
		if raw == "" {
			continue
		}

		if strings.HasPrefix(raw, "{") && strings.Contains(raw, "\"hex\"") {
			var item dto.CharacterColorItem
			if json.Unmarshal([]byte(raw), &item) == nil && strings.TrimSpace(item.Hex) != "" {
				if strings.TrimSpace(item.Name) == "" {
					item.Name = item.Hex
				}
				out = append(out, item)
				continue
			}
		}

		hex := extractHexColor(raw)
		name := strings.TrimSpace(strings.NewReplacer(hex, "", ":", " ", "|", " ", "｜", " ", ",", " ").Replace(raw))
		name = strings.Join(strings.Fields(name), " ")
		if hex == "" {
			hex = "#6C7A89"
		}
		if name == "" {
			name = raw
		}
		out = append(out, dto.CharacterColorItem{Name: name, Hex: hex})
	}
	return out
}

func extractHexColor(raw string) string {
	for i := 0; i < len(raw); i++ {
		if raw[i] != '#' {
			continue
		}
		j := i + 1
		for j < len(raw) {
			ch := raw[j]
			if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F') {
				j++
				continue
			}
			break
		}
		switch j - i {
		case 4, 7, 9:
			return raw[i:j]
		}
	}
	return ""
}
