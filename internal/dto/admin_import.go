package dto

type AdminImportRequest struct {
	Package GeneratedPackage `json:"package"`
}

type AdminImportSummary struct {
	Themes               int `json:"themes"`
	Creators             int `json:"creators"`
	Works                int `json:"works"`
	Characters           int `json:"characters"`
	Songs                int `json:"songs"`
	Relations            int `json:"relations,omitempty"`
	RelationParticipants int `json:"relationParticipants,omitempty"`
	RelationEvents       int `json:"relationEvents,omitempty"`
	RelationSongs        int `json:"relationSongs,omitempty"`
	RelationThemes       int `json:"relationThemes,omitempty"`
	RelationLinks        int `json:"relationLinks,omitempty"`
}

type AdminImportResult struct {
	Valid          bool               `json:"valid"`
	Imported       bool               `json:"imported"`
	PackageVersion string             `json:"packageVersion"`
	Summary        AdminImportSummary `json:"summary"`
	Warnings       []string           `json:"warnings,omitempty"`
	Errors         []string           `json:"errors,omitempty"`
}

type GeneratedPackage struct {
	PackageVersion string               `json:"package_version"`
	PmThemes       []GeneratedTheme     `json:"pm_themes"`
	PmCreators     []GeneratedCreator   `json:"pm_creators"`
	PmWorks        []GeneratedWork      `json:"pm_works"`
	PmCharacters   []GeneratedCharacter `json:"pm_characters"`
	PmSongs        []GeneratedSong      `json:"pm_songs"`
}

type GeneratedTheme struct {
	Code        string `json:"code"`
	Slug        string `json:"slug"`
	NameZH      string `json:"name_zh"`
	Summary     string `json:"summary"`
	Category    string `json:"category"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

type GeneratedCreator struct {
	Name               string         `json:"name"`
	Slug               string         `json:"slug"`
	Aliases            []string       `json:"aliases"`
	CreatorTypeCode    string         `json:"creator_type_code"`
	RegionCode         string         `json:"region_code"`
	CulturalRegionCode string         `json:"cultural_region_code"`
	EraText            string         `json:"era_text"`
	Summary            string         `json:"summary"`
	Introduction       string         `json:"introduction"`
	CoverURL           string         `json:"cover_url"`
	MajorWorks         []string       `json:"major_works"`
	IdentityTags       []string       `json:"identity_tags"`
	Meta               map[string]any `json:"meta"`
	SortOrder          int            `json:"sort_order"`
	IsActive           bool           `json:"is_active"`
}

type GeneratedWorkCreatorRole struct {
	CreatorSlug string `json:"creator_slug"`
	RoleCode    string `json:"role_code"`
	IsPrimary   bool   `json:"is_primary"`
	SortOrder   int    `json:"sort_order"`
	Note        string `json:"note"`
}

type GeneratedWork struct {
	Title              string                     `json:"title"`
	Slug               string                     `json:"slug"`
	Subtitle           string                     `json:"subtitle"`
	OriginalTitle      string                     `json:"original_title"`
	Aliases            []string                   `json:"aliases"`
	WorkTypeCode       string                     `json:"work_type_code"`
	RegionCode         string                     `json:"region_code"`
	CulturalRegionCode string                     `json:"cultural_region_code"`
	EraText            string                     `json:"era_text"`
	ReleaseYear        *int                       `json:"release_year"`
	Summary            string                     `json:"summary"`
	Introduction       string                     `json:"introduction"`
	CoverURL           string                     `json:"cover_url"`
	Themes             []string                   `json:"themes"`
	VersionNotes       []string                   `json:"version_notes"`
	Meta               map[string]any             `json:"meta"`
	SortOrder          int                        `json:"sort_order"`
	IsActive           bool                       `json:"is_active"`
	CreatorRoles       []GeneratedWorkCreatorRole `json:"creator_roles"`
}

type GeneratedColor struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type GeneratedCharacter struct {
	Name                 string            `json:"name"`
	Slug                 string            `json:"slug"`
	Aliases              []string          `json:"aliases"`
	CharacterTypeCode    string            `json:"character_type_code"`
	Subtype              string            `json:"subtype"`
	Gender               string            `json:"gender"`
	RegionCode           string            `json:"region_code"`
	CulturalRegionCode   string            `json:"cultural_region_code"`
	EraText              string            `json:"era_text"`
	DynastyPeriodText    string            `json:"dynasty_period_text"`
	Summary              string            `json:"summary"`
	CoverURL             string            `json:"cover_url"`
	CoverPrompt          string            `json:"cover_prompt"`
	OneLineDefinition    string            `json:"one_line_definition"`
	CoreIdentity         string            `json:"core_identity"`
	PublicImage          string            `json:"public_image"`
	HiddenSelf           string            `json:"hidden_self"`
	CoreFear             string            `json:"core_fear"`
	PsychologicalWound   string            `json:"psychological_wound"`
	CoreConflict         string            `json:"core_conflict"`
	EmotionalTone        string            `json:"emotional_tone"`
	EmotionalTemperature string            `json:"emotional_temperature"`
	Origin               string            `json:"origin"`
	FateArc              string            `json:"fate_arc"`
	EndingState          string            `json:"ending_state"`
	MbtiGuess            []string          `json:"mbti_guess"`
	MbtiConfidence       string            `json:"mbti_confidence"`
	CognitiveStyle       []string          `json:"cognitive_style"`
	SurfaceTraits        []string          `json:"surface_traits"`
	DeepTraits           []string          `json:"deep_traits"`
	BehaviorPatterns     []string          `json:"behavior_patterns"`
	StressResponse       []string          `json:"stress_response"`
	DominantEmotions     []string          `json:"dominant_emotions"`
	SuppressedEmotions   []string          `json:"suppressed_emotions"`
	ValuesTags           []string          `json:"values_tags"`
	BottomLines          []string          `json:"bottom_lines"`
	Taboos               []string          `json:"taboos"`
	SymbolicImages       []string          `json:"symbolic_images"`
	Colors               []GeneratedColor  `json:"colors"`
	Elements             []string          `json:"elements"`
	SoundscapeKeywords   []string          `json:"soundscape_keywords"`
	RelationshipProfile  map[string]string `json:"relationship_profile"`
	Psychology           map[string]any    `json:"psychology"`
	Timeline             []AdminTimeline   `json:"timeline"`
	Meta                 map[string]any    `json:"meta"`
	SortOrder            int               `json:"sort_order"`
	Status               string            `json:"status"`
	IsActive             bool              `json:"is_active"`
	MotivationCodes      []string          `json:"motivation_codes"`
	PrimaryMotivation    string            `json:"primary_motivation_code"`
	ThemeCodes           []string          `json:"theme_codes"`
	PrimaryTheme         string            `json:"primary_theme_code"`
	WorkSlugs            []string          `json:"work_slugs"`
	PrimaryWork          string            `json:"primary_work_slug"`
}

type GeneratedSong struct {
	CharacterSlug      string         `json:"character_slug"`
	Title              string         `json:"title"`
	Slug               string         `json:"slug"`
	Subtitle           string         `json:"subtitle"`
	Summary            string         `json:"summary"`
	CoverURL           string         `json:"cover_url"`
	AudioURL           string         `json:"audio_url"`
	SongCoreTheme      string         `json:"song_core_theme"`
	SongSummary        string         `json:"song_summary"`
	SongEmotionalCurve []string       `json:"song_emotional_curve"`
	SongStyles         []string       `json:"song_styles"`
	TempoBPM           *int           `json:"tempo_bpm"`
	VocalProfile       string         `json:"vocal_profile"`
	LyricKeywords      []string       `json:"lyric_keywords"`
	ForbiddenCliches   []string       `json:"forbidden_cliches"`
	SymbolImages       []string       `json:"symbol_images"`
	EndingFeeling      string         `json:"ending_feeling"`
	Prompt             string         `json:"prompt"`
	Lyrics             string         `json:"lyrics"`
	VersionNo          int            `json:"version_no"`
	Status             string         `json:"status"`
	Meta               map[string]any `json:"meta"`
	SortOrder          int            `json:"sort_order"`
	IsActive           bool           `json:"is_active"`
}
