package repo

import (
	"errors"
	"math/rand/v2"
	"strings"

	"pm-backend/internal/dto"
)

type mockCatalogRepo struct {
	characters    []dto.Character
	relationships []dto.RelationRecord
	works         []dto.Work
	creators      []dto.Creator
	themes        []dto.Theme
	songs         []dto.Song
}

func NewMockCatalogRepo() CatalogRepo {
	works := []dto.Work{
		{
			Slug:     "dream-of-the-red-chamber",
			Title:    "Dream of the Red Chamber",
			Summary:  "A novel about family decline and fragile feeling.",
			CoverURL: "/assets/images/works/dream-of-the-red-chamber.webp",
			TypeCode: "novel",
		},
		{
			Slug:     "journey-to-the-west",
			Title:    "Journey to the West",
			Summary:  "A mythic pilgrimage shaped by rebellion and discipline.",
			CoverURL: "/assets/images/works/journey-to-the-west.webp",
			TypeCode: "novel",
		},
	}

	creators := []dto.Creator{
		{
			Slug:     "cao-xueqin",
			Name:     "Cao Xueqin",
			Summary:  "A classic Chinese novelist.",
			CoverURL: "/assets/images/creators/cao-xueqin.webp",
		},
		{
			Slug:     "wu-chengen",
			Name:     "Wu Chengen",
			Summary:  "A classic Chinese novelist.",
			CoverURL: "/assets/images/creators/wu-chengen.webp",
		},
	}

	themes := []dto.Theme{
		{
			Slug:     "tragic",
			Name:     "Tragic Temperament",
			Summary:  "Characters who cannot live lightly with truth and feeling.",
			CoverURL: "/assets/images/themes/tragic.webp",
		},
		{
			Slug:     "rebels",
			Name:     "Rebels",
			Summary:  "Characters who reject unreasonable order.",
			CoverURL: "/assets/images/themes/rebels.webp",
		},
	}

	songs := []dto.Song{
		{
			Slug:          "lin-daiyu-theme-v1",
			Title:         "Lin Daiyu Theme",
			CharacterSlug: "lin-daiyu",
			CoverURL:      "/assets/images/songs/lin-daiyu-theme-v1.webp",
			AudioURL:      "/assets/audio/lin-daiyu-theme-v1.mp3",
			Styles:        []string{"guofeng", "lyrical"},
		},
		{
			Slug:          "sun-wu-kong-theme-v1",
			Title:         "Sun Wukong Theme",
			CharacterSlug: "sun-wu-kong",
			CoverURL:      "/assets/images/songs/sun-wu-kong-theme-v1.webp",
			AudioURL:      "/assets/audio/sun-wu-kong-theme-v1.mp3",
			Styles:        []string{"rock", "energetic"},
		},
	}

	characters := []dto.Character{
		{
			Slug:              "lin-daiyu",
			Name:              "Lin Daiyu",
			CharacterTypeCode: "literary",
			Summary:           "A character of sensitivity, dignity, and emotional lucidity.",
			OneLineDefinition: "She is not weak; she is too sensitive and too honest to live numbly.",
			CoverURL:          "/assets/images/characters/lin-daiyu.webp",
			ThemeSlugs:        []string{"tragic"},
			WorkSlugs:         []string{"dream-of-the-red-chamber"},
			SongSlugs:         []string{"lin-daiyu-theme-v1"},
			SurfaceTraits:     []string{"sensitive", "lucid"},
			PrimaryWorkTitle:  "Dream of the Red Chamber",
			PrimaryThemeName:  "Tragic Temperament",
			PrimarySongTitle:  "Lin Daiyu Theme",
		},
		{
			Slug:              "sun-wu-kong",
			Name:              "Sun Wukong",
			CharacterTypeCode: "literary",
			Summary:           "A figure of freedom who refuses submission to absurd order.",
			OneLineDefinition: "He does not resist rules in general; he resists rules that should never rule him.",
			CoverURL:          "/assets/images/characters/sun-wu-kong.webp",
			ThemeSlugs:        []string{"rebels"},
			WorkSlugs:         []string{"journey-to-the-west"},
			SongSlugs:         []string{"sun-wu-kong-theme-v1"},
			SurfaceTraits:     []string{"defiant", "free"},
			PrimaryWorkTitle:  "Journey to the West",
			PrimaryThemeName:  "Rebels",
			PrimarySongTitle:  "Sun Wukong Theme",
		},
	}

	relationships := []dto.RelationRecord{
		{
			Slug:                "lin-daiyu--sun-wu-kong--mirror",
			Name:                "Lin Daiyu × Sun Wukong",
			Summary:             "Two people who defend their inner truth in opposite but equally stubborn ways.",
			OneLineDefinition:   "Sensitivity and defiance become two different methods of refusing a false life.",
			CoverURL:            "/assets/images/characters/lin-daiyu.webp",
			RelationTypeCode:    "mirror",
			RelationTypeName:    "Mirror",
			EmotionalTone:       "Tense but illuminating",
			ConnectionTrigger:   "They recognize in each other a refusal to submit to a false order.",
			SustainingMechanism: "Contrast keeps the relationship vivid: grief answers revolt, revolt answers grief.",
			PowerStructure:      "Neither dominates; their energy comes from contrast rather than control.",
			DependencyPattern:   "They do not need each other directly, but they reveal alternate forms of resistance.",
			RelationConflict:    "Their refusal points in different directions: one toward grief, one toward revolt.",
			RelationArc:         "The relation begins as contrast and deepens into mutual illumination.",
			FateImpact:          "Each becomes a lens for the other's unfinished possibility.",
			SymbolicImages:      []string{"lantern", "staff"},
			ThemeTags:           []string{"tragic", "rebels"},
			SourceCharacter: dto.RelationshipCharacterRef{
				Slug:     "lin-daiyu",
				Name:     "Lin Daiyu",
				CoverURL: "/assets/images/characters/lin-daiyu.webp",
				Summary:  "A character of sensitivity, dignity, and emotional lucidity.",
			},
			TargetCharacter: dto.RelationshipCharacterRef{
				Slug:     "sun-wu-kong",
				Name:     "Sun Wukong",
				CoverURL: "/assets/images/characters/sun-wu-kong.webp",
				Summary:  "A figure of freedom who refuses submission to absurd order.",
			},
			Phenomenology: dto.RelationPhenomenology{
				Body:     "The body tightens before either one yields.",
				Time:     "Time feels suspended between recognition and divergence.",
				Space:    "Space becomes charged whenever they occupy the same moral field.",
				Gaze:     "Their gaze is less affectionate than clarifying.",
				Language: "Language cuts rather than soothes.",
			},
			RelationPalette: []dto.RelationPaletteItem{
				{Name: "Lantern Gold", Hex: "#D6B36A"},
				{Name: "Storm Ink", Hex: "#1A1D23"},
			},
			RelationKeywords: []string{"mirror", "defiance", "lucidity"},
			Events: []dto.RelationEvent{
				{StageNo: 1, StageCode: "recognition", Title: "Recognition", Summary: "They first see in each other an unwillingness to live falsely."},
				{StageNo: 2, StageCode: "contrast", Title: "Contrast", Summary: "Their responses diverge, making the distance between grief and revolt visible."},
			},
			PrimarySong: &dto.RelationSong{
				Slug:               "lin-daiyu--sun-wu-kong--mirror-theme",
				Title:              "Mirror of Defiance",
				Summary:            "A relationship song about two incompatible but illuminating refusals.",
				CoverURL:           "/assets/images/characters/lin-daiyu.webp",
				AudioURL:           "/assets/audio/lin-daiyu-theme-v1.mp3",
				SongCoreTheme:      "mirror resistance",
				SongEmotionalCurve: "restraint -> flare -> recognition",
				SongStyles:         []string{"cinematic", "lyrical"},
				VocalProfile:       "dual vocal tension",
			},
		},
	}

	return &mockCatalogRepo{
		characters:    characters,
		relationships: relationships,
		works:         works,
		creators:      creators,
		themes:        themes,
		songs:         songs,
	}
}

func (m *mockCatalogRepo) Home() (dto.HomePayload, error) {
	featured := dto.Character{}
	if len(m.characters) > 0 {
		featured = m.characters[0]
	}

	return dto.HomePayload{
		FeaturedCharacter: featured,
		LatestCharacters:  m.characters,
		FeaturedSongs:     m.songs,
		RecommendedWorks:  m.works,
		Themes:            m.themes,
	}, nil
}

func (m *mockCatalogRepo) RandomCharacter(theme string, exclude []string) (dto.Character, error) {
	pool := make([]dto.Character, 0, len(m.characters))
	excluded := make(map[string]struct{}, len(exclude))
	for _, slug := range exclude {
		excluded[slug] = struct{}{}
	}

	for _, item := range m.characters {
		if _, skip := excluded[item.Slug]; skip {
			continue
		}
		if theme != "" && !containsString(item.ThemeSlugs, theme) {
			continue
		}
		pool = append(pool, item)
	}

	if len(pool) == 0 {
		return dto.Character{}, errors.New("no characters")
	}
	return pool[rand.IntN(len(pool))], nil
}

func (m *mockCatalogRepo) ListCharacters() ([]dto.Character, error) { return m.characters, nil }

func (m *mockCatalogRepo) GetCharacterDetail(slug string) (dto.CharacterDetail, error) {
	for _, character := range m.characters {
		if character.Slug != slug {
			continue
		}

		themes := make([]dto.Theme, 0)
		for _, themeSlug := range character.ThemeSlugs {
			for _, theme := range m.themes {
				if theme.Slug == themeSlug {
					themes = append(themes, theme)
				}
			}
		}

		works := make([]dto.Work, 0)
		for _, workSlug := range character.WorkSlugs {
			for _, work := range m.works {
				if work.Slug == workSlug {
					works = append(works, work)
				}
			}
		}

		creators := make([]dto.Creator, 0)
		for _, work := range works {
			switch work.Slug {
			case "dream-of-the-red-chamber":
				creators = append(creators, m.creators[0])
			case "journey-to-the-west":
				creators = append(creators, m.creators[1])
			}
		}

		return dto.CharacterDetail{
			Character:          character,
			CoreIdentity:       character.Summary,
			PrimaryMotivation:  "To protect inner truth from compromise.",
			CoreFear:           "A life emptied of meaning.",
			PsychologicalWound: "The self is wounded whenever truth is denied or humiliated.",
			CoreConflict:       "The desire to stay true collides with the world’s demand for adaptation.",
			EmotionalTone:      "Tense, lucid, and quietly intense.",
			SurfaceTraits:      character.SurfaceTraits,
			ValuesTags:         []string{"truth", "dignity"},
			SymbolicImages:     []string{"rain", "lantern"},
			Timeline:           []dto.CharacterTimelineItem{},
			RelatedWorks:       works,
			RelatedThemes:      themes,
			RelatedSongs:       filterSongsByCharacter(m.songs, slug),
			RelatedCreator:     creators,
		}, nil
	}

	return dto.CharacterDetail{}, errors.New("character not found")
}

func (m *mockCatalogRepo) ListRelationships(characterSlug string) ([]dto.RelationRecord, error) {
	if strings.TrimSpace(characterSlug) == "" {
		return m.relationships, nil
	}

	list := make([]dto.RelationRecord, 0, len(m.relationships))
	for _, item := range m.relationships {
		if item.SourceCharacter.Slug == characterSlug || item.TargetCharacter.Slug == characterSlug {
			list = append(list, item)
		}
	}
	return list, nil
}

func (m *mockCatalogRepo) GetRelationshipDetail(slug string) (dto.RelationRecord, error) {
	for _, item := range m.relationships {
		if item.Slug == slug {
			return item, nil
		}
	}
	return dto.RelationRecord{}, errors.New("relationship not found")
}

func (m *mockCatalogRepo) ListWorks() ([]dto.Work, error) { return m.works, nil }

func (m *mockCatalogRepo) GetWorkDetail(slug string) (dto.Work, error) {
	for _, item := range m.works {
		if item.Slug == slug {
			return item, nil
		}
	}
	return dto.Work{}, errors.New("work not found")
}

func (m *mockCatalogRepo) ListCreators() ([]dto.Creator, error) { return m.creators, nil }

func (m *mockCatalogRepo) GetCreatorDetail(slug string) (dto.Creator, error) {
	for _, item := range m.creators {
		if item.Slug == slug {
			return item, nil
		}
	}
	return dto.Creator{}, errors.New("creator not found")
}

func (m *mockCatalogRepo) ListThemes() ([]dto.Theme, error) { return m.themes, nil }

func (m *mockCatalogRepo) GetThemeDetail(slug string) (dto.ThemeDetail, error) {
	for _, theme := range m.themes {
		if theme.Slug != slug {
			continue
		}

		characters := make([]dto.Character, 0)
		for _, character := range m.characters {
			if containsString(character.ThemeSlugs, slug) {
				characters = append(characters, character)
			}
		}
		return dto.ThemeDetail{Theme: theme, Characters: characters}, nil
	}
	return dto.ThemeDetail{}, errors.New("theme not found")
}

func (m *mockCatalogRepo) ListSongs() ([]dto.Song, error) { return m.songs, nil }

func (m *mockCatalogRepo) SearchCatalog(keyword string, limit int) (dto.SearchResponseData, error) {
	if limit <= 0 {
		limit = 8
	}
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	match := func(parts ...string) bool {
		if keyword == "" {
			return true
		}
		return strings.Contains(strings.ToLower(strings.Join(parts, " ")), keyword)
	}

	out := dto.SearchResponseData{
		Characters: make([]dto.CharacterListItemResponse, 0, limit),
		Works:      make([]dto.WorkListItemResponse, 0, limit),
		Creators:   make([]dto.CreatorListItemResponse, 0, limit),
		Themes:     make([]dto.ThemeListItemResponse, 0, limit),
		Songs:      make([]dto.SongListItemResponse, 0, limit),
	}

	for _, item := range m.characters {
		if len(out.Characters) >= limit || !match(item.Name, item.Summary, item.OneLineDefinition) {
			continue
		}
		themeSongTitle := ""
		if len(item.SongSlugs) > 0 {
			themeSongTitle = item.PrimarySongTitle
		}
		out.Characters = append(out.Characters, dto.CharacterListItemResponse{
			ID:                item.Slug,
			Slug:              item.Slug,
			Name:              item.Name,
			CoverURL:          item.CoverURL,
			Summary:           item.Summary,
			OneLineDefinition: item.OneLineDefinition,
			CharacterTypeCode: item.CharacterTypeCode,
			WorkTitle:         item.PrimaryWorkTitle,
			Tags:              []string{item.PrimaryThemeName},
			HasSong:           len(item.SongSlugs) > 0,
			ThemeSongTitle:    themeSongTitle,
		})
	}

	for _, item := range m.works {
		if len(out.Works) >= limit || !match(item.Title, item.Summary) {
			continue
		}
		out.Works = append(out.Works, dto.WorkListItemResponse{
			ID:             item.Slug,
			Slug:           item.Slug,
			Title:          item.Title,
			CoverURL:       item.CoverURL,
			Summary:        item.Summary,
			WorkTypeCode:   item.TypeCode,
			CharacterCount: len(item.CharacterSlugs),
		})
	}

	for _, item := range m.creators {
		if len(out.Creators) >= limit || !match(item.Name, item.Summary, item.EraText) {
			continue
		}
		out.Creators = append(out.Creators, dto.CreatorListItemResponse{
			ID:              item.Slug,
			Slug:            item.Slug,
			Name:            item.Name,
			CoverURL:        item.CoverURL,
			Summary:         item.Summary,
			CreatorTypeCode: item.CreatorTypeCode,
			EraText:         item.EraText,
			WorkCount:       len(item.WorkSlugs),
		})
	}

	for _, item := range m.themes {
		if len(out.Themes) >= limit || !match(item.Name, item.Summary, item.Category) {
			continue
		}
		out.Themes = append(out.Themes, dto.ThemeListItemResponse{
			ID:       item.Slug,
			Slug:     item.Slug,
			Name:     item.Name,
			CoverURL: item.CoverURL,
			Summary:  item.Summary,
			Category: item.Category,
		})
	}

	for _, item := range m.songs {
		if len(out.Songs) >= limit || !match(item.Title, item.Summary, item.SongCoreTheme) {
			continue
		}
		out.Songs = append(out.Songs, dto.SongListItemResponse{
			ID:            item.Slug,
			Slug:          item.Slug,
			Title:         item.Title,
			CharacterSlug: item.CharacterSlug,
			CoverURL:      item.CoverURL,
			AudioURL:      item.AudioURL,
			Styles:        item.Styles,
		})
	}

	return out, nil
}

func filterSongsByCharacter(in []dto.Song, slug string) []dto.Song {
	out := make([]dto.Song, 0)
	for _, item := range in {
		if item.CharacterSlug == slug {
			out = append(out, item)
		}
	}
	return out
}

func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
