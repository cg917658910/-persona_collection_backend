package repo

import (
	"errors"
	"math/rand/v2"

	"pm-backend/internal/dto"
)

type mockCatalogRepo struct {
	characters []dto.Character
	works      []dto.Work
	creators   []dto.Creator
	themes     []dto.Theme
	songs      []dto.Song
}

func NewMockCatalogRepo() CatalogRepo {
	works := []dto.Work{
		{Slug: "dream-of-the-red-chamber", Title: "红楼梦", Summary: "家族盛衰与真情消逝。", CoverURL: "/assets/images/works/dream-of-the-red-chamber.webp", TypeCode: "novel"},
		{Slug: "journey-to-the-west", Title: "西游记", Summary: "神魔冒险与修行之路。", CoverURL: "/assets/images/works/journey-to-the-west.webp", TypeCode: "novel"},
	}
	creators := []dto.Creator{
		{Slug: "cao-xueqin", Name: "曹雪芹", Summary: "中国古典文学作家。", CoverURL: "/assets/images/creators/cao-xueqin.webp"},
		{Slug: "wu-chengen", Name: "吴承恩", Summary: "中国古典文学作家。", CoverURL: "/assets/images/creators/wu-chengen.webp"},
	}
	themes := []dto.Theme{
		{Slug: "tragic", Name: "悲剧人格", Summary: "内在真实与命运难以相容。", CoverURL: "/assets/images/themes/tragic.webp"},
		{Slug: "rebels", Name: "反叛者", Summary: "拒绝被不合理秩序驯服。", CoverURL: "/assets/images/themes/rebels.webp"},
	}
	songs := []dto.Song{
		{Slug: "lin-daiyu-theme-v1", Title: "林黛玉之歌", CharacterSlug: "lin-daiyu", CoverURL: "/assets/images/songs/lin-daiyu-theme-v1.webp", AudioURL: "/assets/audio/lin-daiyu-theme-v1.mp3", Styles: []string{"国风", "抒情"}},
		{Slug: "sun-wu-kong-theme-v1", Title: "孙悟空之歌", CharacterSlug: "sun-wu-kong", CoverURL: "/assets/images/songs/sun-wu-kong-theme-v1.webp", AudioURL: "/assets/audio/sun-wu-kong-theme-v1.mp3", Styles: []string{"摇滚", "热血"}},
	}
	characters := []dto.Character{
		{
			Slug: "lin-daiyu", Name: "林黛玉", CharacterTypeCode: "literary",
			Summary: "高敏感、高自尊的真情承受者。", OneLineDefinition: "她不是脆弱，而是过度敏感、过度清醒、又过度珍视真情的人。",
			CoverURL: "/assets/images/characters/lin-daiyu.webp", ThemeSlugs: []string{"tragic"}, WorkSlugs: []string{"dream-of-the-red-chamber"}, SongSlugs: []string{"lin-daiyu-theme-v1"},
		},
		{
			Slug: "sun-wu-kong", Name: "孙悟空", CharacterTypeCode: "literary",
			Summary: "不愿被驯服的自由反抗者。", OneLineDefinition: "他不是不守规矩，而是不愿被不合理的规矩困住。",
			CoverURL: "/assets/images/characters/sun-wu-kong.webp", ThemeSlugs: []string{"rebels"}, WorkSlugs: []string{"journey-to-the-west"}, SongSlugs: []string{"sun-wu-kong-theme-v1"},
		},
	}

	return &mockCatalogRepo{
		characters: characters,
		works:      works,
		creators:   creators,
		themes:     themes,
		songs:      songs,
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
		if theme != "" {
			matched := false
			for _, themeSlug := range item.ThemeSlugs {
				if themeSlug == theme {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
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
	for _, c := range m.characters {
		if c.Slug == slug {
			return dto.CharacterDetail{
				Character:      c,
				CoreIdentity:   c.Summary,
				CoreFear:       "待补充",
				CoreConflict:   "待补充",
				EmotionalTone:  "待补充",
				SurfaceTraits:  []string{"待补充"},
				Timeline:       []string{"待补充"},
				RelatedWorks:   m.works,
				RelatedThemes:  m.themes,
				RelatedSongs:   filterSongsByCharacter(m.songs, slug),
				RelatedCreator: m.creators,
			}, nil
		}
	}
	return dto.CharacterDetail{}, errors.New("character not found")
}

func (m *mockCatalogRepo) ListWorks() ([]dto.Work, error) { return m.works, nil }

func (m *mockCatalogRepo) GetWorkDetail(slug string) (dto.Work, error) {
	for _, v := range m.works {
		if v.Slug == slug {
			return v, nil
		}
	}
	return dto.Work{}, errors.New("work not found")
}

func (m *mockCatalogRepo) ListCreators() ([]dto.Creator, error) { return m.creators, nil }

func (m *mockCatalogRepo) GetCreatorDetail(slug string) (dto.Creator, error) {
	for _, v := range m.creators {
		if v.Slug == slug {
			return v, nil
		}
	}
	return dto.Creator{}, errors.New("creator not found")
}

func (m *mockCatalogRepo) ListThemes() ([]dto.Theme, error) { return m.themes, nil }

func (m *mockCatalogRepo) GetThemeDetail(slug string) (dto.ThemeDetail, error) {
	for _, t := range m.themes {
		if t.Slug == slug {
			chars := make([]dto.Character, 0)
			for _, c := range m.characters {
				for _, ts := range c.ThemeSlugs {
					if ts == slug {
						chars = append(chars, c)
						break
					}
				}
			}
			return dto.ThemeDetail{Theme: t, Characters: chars}, nil
		}
	}
	return dto.ThemeDetail{}, errors.New("theme not found")
}

func (m *mockCatalogRepo) ListSongs() ([]dto.Song, error) { return m.songs, nil }

func filterSongsByCharacter(in []dto.Song, slug string) []dto.Song {
	out := make([]dto.Song, 0)
	for _, s := range in {
		if s.CharacterSlug == slug {
			out = append(out, s)
		}
	}
	return out
}
