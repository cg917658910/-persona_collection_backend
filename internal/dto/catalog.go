package dto

type Character struct {
	Slug              string   `json:"slug"`
	Name              string   `json:"name"`
	CharacterTypeCode string   `json:"characterTypeCode"`
	Summary           string   `json:"summary"`
	OneLineDefinition string   `json:"oneLineDefinition"`
	CoverURL          string   `json:"coverUrl"`
	ThemeSlugs        []string `json:"themeSlugs"`
	WorkSlugs         []string `json:"workSlugs"`
	SongSlugs         []string `json:"songSlugs"`
}

type CharacterDetail struct {
	Character
	CoreIdentity   string   `json:"coreIdentity"`
	CoreFear       string   `json:"coreFear"`
	CoreConflict   string   `json:"coreConflict"`
	EmotionalTone  string   `json:"emotionalTone"`
	SurfaceTraits  []string `json:"surfaceTraits"`
	Timeline       []string `json:"timeline"`
	RelatedWorks   []Work   `json:"relatedWorks"`
	RelatedThemes  []Theme  `json:"relatedThemes"`
	RelatedSongs   []Song   `json:"relatedSongs"`
	RelatedCreator []Creator `json:"relatedCreators"`
}

type Work struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	CoverURL  string `json:"coverUrl"`
	TypeCode  string `json:"workTypeCode"`
}

type Creator struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	CoverURL string `json:"coverUrl"`
}

type Theme struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	CoverURL string `json:"coverUrl"`
}

type ThemeDetail struct {
	Theme
	Characters []Character `json:"characters"`
}

type Song struct {
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	CharacterSlug string   `json:"characterSlug"`
	CoverURL      string   `json:"coverUrl"`
	AudioURL      string   `json:"audioUrl"`
	Styles        []string `json:"styles"`
}

type HomePayload struct {
	FeaturedCharacter Character   `json:"featuredCharacter"`
	LatestCharacters  []Character `json:"latestCharacters"`
	FeaturedSongs     []Song      `json:"featuredSongs"`
	RecommendedWorks  []Work      `json:"recommendedWorks"`
	Themes            []Theme     `json:"themes"`
}

type ListMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
