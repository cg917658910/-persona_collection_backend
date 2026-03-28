package dto

type AdminCharacter struct {
	ID                   string            `json:"id,omitempty"`
	Slug                 string            `json:"slug"`
	Name                 string            `json:"name"`
	Summary              string            `json:"summary"`
	CoverURL             string            `json:"coverUrl"`
	Status               string            `json:"status"`
	Type                 string            `json:"type"`
	CharacterTypeCode    string            `json:"characterTypeCode"`
	Gender               string            `json:"gender,omitempty"`
	RegionCode           string            `json:"regionCode,omitempty"`
	CulturalRegionCode   string            `json:"culturalRegionCode,omitempty"`
	OneLineDefinition    string            `json:"oneLineDefinition"`
	CoreIdentity         string            `json:"coreIdentity"`
	MotivationNote       string            `json:"motivationNote,omitempty"`
	CoreFear             string            `json:"coreFear"`
	CoreConflict         string            `json:"coreConflict"`
	EmotionalTone        string            `json:"emotionalTone"`
	EmotionalTemperature string            `json:"emotionalTemperature,omitempty"`
	PrimaryMotivation    string            `json:"primaryMotivation,omitempty"`
	WorkSlugs            []string          `json:"workSlugs,omitempty"`
	WorkNames            []string          `json:"workNames,omitempty"`
	ThemeSlugs           []string          `json:"themeSlugs,omitempty"`
	ThemeNames           []string          `json:"themeNames,omitempty"`
	SongSlugs            []string          `json:"songSlugs,omitempty"`
	HasSong              bool              `json:"hasSong"`
	DominantEmotions     []string          `json:"dominantEmotions,omitempty"`
	SuppressedEmotions   []string          `json:"suppressedEmotions,omitempty"`
	ValuesTags           []string          `json:"valuesTags,omitempty"`
	SymbolicImages       []string          `json:"symbolicImages,omitempty"`
	Elements             []string          `json:"elements,omitempty"`
	RelationshipProfile  map[string]string `json:"relationshipProfile,omitempty"`
	Timeline             []AdminTimeline   `json:"timeline,omitempty"`
	SortOrder            int               `json:"sortOrder"`
	HomeToday            bool              `json:"homeToday"`
	FeaturedHome         bool              `json:"featuredHome"`
	HomeSort             int               `json:"homeSort"`
	DiscoverWeight       float64           `json:"discoverWeight"`
}

type AdminTimeline struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type AdminSong struct {
	ID             string   `json:"id,omitempty"`
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	CoverURL       string   `json:"coverUrl"`
	AudioURL       string   `json:"audioUrl"`
	Status         string   `json:"status"`
	CharacterSlug  string   `json:"characterSlug"`
	CharacterName  string   `json:"characterName,omitempty"`
	CoreTheme      string   `json:"coreTheme"`
	Styles         []string `json:"styles,omitempty"`
	EmotionalCurve []string `json:"emotionalCurve,omitempty"`
	Prompt         string   `json:"prompt,omitempty"`
	Lyrics         string   `json:"lyrics,omitempty"`
	SortOrder      int      `json:"sortOrder"`
	FeaturedHome   bool     `json:"featuredHome"`
	HomeSort       int      `json:"homeSort"`
}

type AdminTheme struct {
	ID             string   `json:"id,omitempty"`
	Slug           string   `json:"slug"`
	Name           string   `json:"name"`
	Code           string   `json:"code"`
	SubjectType    string   `json:"subjectType,omitempty"`
	Category       string   `json:"category"`
	Summary        string   `json:"summary"`
	CoverURL       string   `json:"coverUrl"`
	Status         string   `json:"status"`
	CharacterSlugs []string `json:"characterSlugs,omitempty"`
	RelationSlugs  []string `json:"relationSlugs,omitempty"`
	SortOrder      int      `json:"sortOrder"`
}
