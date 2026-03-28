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
	SurfaceTraits     []string `json:"surfaceTraits,omitempty"`
	PrimaryWorkTitle  string   `json:"primaryWorkTitle,omitempty"`
	PrimaryThemeName  string   `json:"primaryThemeName,omitempty"`
	PrimarySongTitle  string   `json:"primarySongTitle,omitempty"`
}

type CharacterTimelineItem struct {
	Year    string `json:"year"`
	Event   string `json:"event"`
	Emotion string `json:"emotion"`
}

type CharacterColorItem struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type CharacterDetail struct {
	Character
	ID                  string                  `json:"id,omitempty"`
	CoreIdentity        string                  `json:"coreIdentity"`
	PublicImage         string                  `json:"publicImage,omitempty"`
	HiddenSelf          string                  `json:"hiddenSelf,omitempty"`
	PrimaryMotivation   string                  `json:"primaryMotivation,omitempty"`
	CoreFear            string                  `json:"coreFear"`
	PsychologicalWound  string                  `json:"psychologicalWound,omitempty"`
	CoreConflict        string                  `json:"coreConflict"`
	EmotionalTone       string                  `json:"emotionalTone"`
	Origin              string                  `json:"origin,omitempty"`
	FateArc             string                  `json:"fateArc,omitempty"`
	EndingState         string                  `json:"endingState,omitempty"`
	SurfaceTraits       []string                `json:"surfaceTraits"`
	DeepTraits          []string                `json:"deepTraits,omitempty"`
	DominantEmotions    []string                `json:"dominantEmotions,omitempty"`
	SuppressedEmotions  []string                `json:"suppressedEmotions,omitempty"`
	ValuesTags          []string                `json:"valuesTags,omitempty"`
	BottomLines         []string                `json:"bottomLines,omitempty"`
	SymbolicImages      []string                `json:"symbolicImages,omitempty"`
	Colors              []CharacterColorItem    `json:"colors,omitempty"`
	Elements            []string                `json:"elements,omitempty"`
	SoundscapeKeywords  []string                `json:"soundscapeKeywords,omitempty"`
	RelationshipProfile map[string]string       `json:"relationshipProfile,omitempty"`
	Timeline            []CharacterTimelineItem `json:"timeline"`
	RelatedWorks        []Work                  `json:"relatedWorks"`
	RelatedThemes       []Theme                 `json:"relatedThemes"`
	RelatedSongs        []Song                  `json:"relatedSongs"`
	RelatedCreator      []Creator               `json:"relatedCreators"`
}

type Work struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	CoverURL       string   `json:"coverUrl"`
	TypeCode       string   `json:"workTypeCode"`
	CreatorSlugs   []string `json:"creatorSlugs,omitempty"`
	CreatorNames   []string `json:"creatorNames,omitempty"`
	CharacterSlugs []string `json:"characterSlugs,omitempty"`
}

type Creator struct {
	Slug            string   `json:"slug"`
	Name            string   `json:"name"`
	Summary         string   `json:"summary"`
	CoverURL        string   `json:"coverUrl"`
	CreatorTypeCode string   `json:"creatorTypeCode,omitempty"`
	EraText         string   `json:"eraText,omitempty"`
	WorkSlugs       []string `json:"workSlugs,omitempty"`
}

type Theme struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	CoverURL    string `json:"coverUrl"`
	Category    string `json:"category,omitempty"`
	SubjectType string `json:"subjectType,omitempty"`
}

type ThemeDetail struct {
	Theme
	Characters    []Character      `json:"characters"`
	Relationships []RelationRecord `json:"relationships"`
}

type Song struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	CharacterSlug  string   `json:"characterSlug"`
	CoverURL       string   `json:"coverUrl"`
	AudioURL       string   `json:"audioUrl"`
	Summary        string   `json:"summary,omitempty"`
	SongCoreTheme  string   `json:"songCoreTheme,omitempty"`
	Styles         []string `json:"styles"`
	EmotionalCurve []string `json:"emotionalCurve,omitempty"`
	VocalProfile   string   `json:"vocalProfile,omitempty"`
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

type HomeCategoryCounts struct {
	Characters int `json:"characters"`
	Creators   int `json:"creators"`
	Historical int `json:"historical"`
	Literary   int `json:"literary"`
	FilmTV     int `json:"film_tv"`
	Anime      int `json:"anime"`
	Works      int `json:"works"`
	Themes     int `json:"themes"`
}

type HomeSongRef struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	CoverURL      string `json:"coverUrl"`
	AudioURL      string `json:"audioUrl"`
	CharacterSlug string `json:"characterSlug"`
	CharacterName string `json:"characterName,omitempty"`
}

type HomeFeaturedCharacter struct {
	ID                string       `json:"id"`
	Slug              string       `json:"slug"`
	Name              string       `json:"name"`
	CoverURL          string       `json:"coverUrl"`
	Summary           string       `json:"summary"`
	OneLineDefinition string       `json:"oneLineDefinition"`
	CharacterTypeCode string       `json:"characterTypeCode"`
	WorkTitle         string       `json:"workTitle,omitempty"`
	Tags              []string     `json:"tags"`
	Song              *HomeSongRef `json:"song,omitempty"`
}

type HomeCharacterCard struct {
	ID                string   `json:"id"`
	Slug              string   `json:"slug"`
	Name              string   `json:"name"`
	CoverURL          string   `json:"coverUrl"`
	Summary           string   `json:"summary"`
	OneLineDefinition string   `json:"oneLineDefinition"`
	CharacterTypeCode string   `json:"characterTypeCode"`
	WorkTitle         string   `json:"workTitle,omitempty"`
	Tags              []string `json:"tags"`
	HasSong           bool     `json:"hasSong"`
}

type HomeSongCard struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	CoverURL      string `json:"coverUrl"`
	AudioURL      string `json:"audioUrl"`
	CharacterSlug string `json:"characterSlug"`
	CharacterName string `json:"characterName"`
	Summary       string `json:"summary,omitempty"`
	SongCoreTheme string `json:"songCoreTheme,omitempty"`
}

type HomeWorkCard struct {
	ID           string `json:"id"`
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	CoverURL     string `json:"coverUrl"`
	Summary      string `json:"summary"`
	WorkTypeCode string `json:"workTypeCode"`
	CreatorName  string `json:"creatorName,omitempty"`
}

type HomeThemeCard struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	CoverURL       string `json:"coverUrl"`
	Summary        string `json:"summary"`
	CharacterCount int    `json:"characterCount"`
}

type HomeResponseData struct {
	FeaturedCharacter *HomeFeaturedCharacter `json:"featuredCharacter"`
	LatestCharacters  []HomeCharacterCard    `json:"latestCharacters"`
	FeaturedSongs     []HomeSongCard         `json:"featuredSongs"`
	RecommendedWorks  []HomeWorkCard         `json:"recommendedWorks"`
	RecommendedThemes []HomeThemeCard        `json:"recommendedThemes"`
	CategoryCounts    HomeCategoryCounts     `json:"categoryCounts"`
}

type CharacterListItemResponse struct {
	ID                string   `json:"id"`
	Slug              string   `json:"slug"`
	Name              string   `json:"name"`
	CoverURL          string   `json:"coverUrl"`
	Summary           string   `json:"summary"`
	OneLineDefinition string   `json:"oneLineDefinition"`
	CharacterTypeCode string   `json:"characterTypeCode"`
	WorkTitle         string   `json:"workTitle,omitempty"`
	Tags              []string `json:"tags"`
	HasSong           bool     `json:"hasSong"`
	ThemeSongTitle    string   `json:"themeSongTitle,omitempty"`
}

type CharacterListResponse struct {
	Items    []CharacterListItemResponse `json:"items"`
	Total    int                         `json:"total,omitempty"`
	Page     int                         `json:"page,omitempty"`
	PageSize int                         `json:"pageSize,omitempty"`
}

type CharacterDetailRef struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	CoverURL string `json:"coverUrl"`
	Summary  string `json:"summary,omitempty"`
}

type CharacterDetailSong struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	CoverURL       string   `json:"coverUrl"`
	AudioURL       string   `json:"audioUrl"`
	SongCoreTheme  string   `json:"songCoreTheme,omitempty"`
	EmotionalCurve []string `json:"emotionalCurve"`
	SongStyles     []string `json:"songStyles"`
	VocalProfile   string   `json:"vocalProfile,omitempty"`
	Lyrics         []string `json:"lyrics"`
}

type CharacterRelationshipPattern struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type RelationshipCharacterRef struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	CoverURL string `json:"coverUrl"`
	Summary  string `json:"summary,omitempty"`
}

type RelationPaletteItem struct {
	Name string `json:"name,omitempty"`
	Hex  string `json:"hex"`
}

type RelationPhenomenology struct {
	Body     string `json:"body,omitempty"`
	Time     string `json:"time,omitempty"`
	Space    string `json:"space,omitempty"`
	Gaze     string `json:"gaze,omitempty"`
	Language string `json:"language,omitempty"`
}

type RelationEvent struct {
	StageNo      int    `json:"stage_no"`
	StageCode    string `json:"stage_code,omitempty"`
	Title        string `json:"title"`
	Summary      string `json:"summary,omitempty"`
	TensionShift string `json:"tension_shift,omitempty"`
	PowerShift   string `json:"power_shift,omitempty"`
	FateImpact   string `json:"fate_impact,omitempty"`
	SourceState  string `json:"source_state,omitempty"`
	TargetState  string `json:"target_state,omitempty"`
	EventQuote   string `json:"event_quote,omitempty"`
	ColorHex     string `json:"color_hex,omitempty"`
}

type RelationSong struct {
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Subtitle           string   `json:"subtitle,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	CoverURL           string   `json:"cover_url,omitempty"`
	AudioURL           string   `json:"audio_url,omitempty"`
	SongCoreTheme      string   `json:"song_core_theme,omitempty"`
	SongEmotionalCurve string   `json:"song_emotional_curve,omitempty"`
	SongStyles         []string `json:"song_styles"`
	VocalProfile       string   `json:"vocal_profile,omitempty"`
	Lyric              string   `json:"lyric,omitempty"`
}

type RelationLink struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	CoverURL     string `json:"cover_url,omitempty"`
	LinkTypeCode string `json:"link_type_code"`
	Reason       string `json:"reason,omitempty"`
}

type RelationRecord struct {
	Slug                   string
	Name                   string
	Subtitle               string
	Summary                string
	OneLineDefinition      string
	CoverURL               string
	RelationTypeCode       string
	RelationTypeName       string
	WorkSlug               string
	WorkName               string
	CoreTension            string
	EmotionalTone          string
	ConnectionTrigger      string
	SustainingMechanism    string
	RelationConflict       string
	RelationArc            string
	FateImpact             string
	PowerStructure         string
	DependencyPattern      string
	SourcePerspective      string
	SourceDesireInRelation string
	SourceFearInRelation   string
	SourceUnsaid           string
	TargetPerspective      string
	TargetDesireInRelation string
	TargetFearInRelation   string
	TargetUnsaid           string
	Phenomenology          RelationPhenomenology
	SymbolicImages         []string
	ThemeTags              []string
	RelationPalette        []RelationPaletteItem
	RelationKeywords       []string
	SourceCharacter        RelationshipCharacterRef
	TargetCharacter        RelationshipCharacterRef
	Events                 []RelationEvent
	PrimarySong            *RelationSong
	RelatedRelations       []RelationLink
}

type RelationshipListItemResponse struct {
	ID                string                    `json:"id"`
	Slug              string                    `json:"slug"`
	Name              string                    `json:"name"`
	Summary           string                    `json:"summary,omitempty"`
	OneLineDefinition string                    `json:"oneLineDefinition,omitempty"`
	CoverURL          string                    `json:"coverUrl"`
	RelationType      string                    `json:"relationType"`
	RelationLabel     string                    `json:"relationLabel,omitempty"`
	WorkTitle         string                    `json:"workTitle,omitempty"`
	Intensity         float64                   `json:"intensity,omitempty"`
	Tags              []string                  `json:"tags"`
	SourceCharacter   RelationshipCharacterRef  `json:"sourceCharacter"`
	TargetCharacter   RelationshipCharacterRef  `json:"targetCharacter"`
	Counterpart       *RelationshipCharacterRef `json:"counterpart,omitempty"`
}

type RelationshipDetailResponse struct {
	Slug                    string                `json:"slug"`
	Name                    string                `json:"name"`
	Subtitle                string                `json:"subtitle,omitempty"`
	RelationTypeCode        string                `json:"relation_type_code,omitempty"`
	RelationTypeName        string                `json:"relation_type_name,omitempty"`
	SourceCharacterSlug     string                `json:"source_character_slug"`
	SourceCharacterName     string                `json:"source_character_name"`
	SourceCharacterCoverURL string                `json:"source_character_cover_url,omitempty"`
	TargetCharacterSlug     string                `json:"target_character_slug"`
	TargetCharacterName     string                `json:"target_character_name"`
	TargetCharacterCoverURL string                `json:"target_character_cover_url,omitempty"`
	WorkSlug                string                `json:"work_slug,omitempty"`
	WorkName                string                `json:"work_name,omitempty"`
	CoreTension             string                `json:"core_tension,omitempty"`
	EmotionalTone           string                `json:"emotional_tone,omitempty"`
	OneLineDefinition       string                `json:"one_line_definition,omitempty"`
	Summary                 string                `json:"summary,omitempty"`
	CoverURL                string                `json:"cover_url,omitempty"`
	ConnectionTrigger       string                `json:"connection_trigger,omitempty"`
	SustainingMechanism     string                `json:"sustaining_mechanism,omitempty"`
	RelationConflict        string                `json:"relation_conflict,omitempty"`
	RelationArc             string                `json:"relation_arc,omitempty"`
	FateImpact              string                `json:"fate_impact,omitempty"`
	PowerStructure          string                `json:"power_structure,omitempty"`
	DependencyPattern       string                `json:"dependency_pattern,omitempty"`
	SourcePerspective       string                `json:"source_perspective,omitempty"`
	SourceDesireInRelation  string                `json:"source_desire_in_relation,omitempty"`
	SourceFearInRelation    string                `json:"source_fear_in_relation,omitempty"`
	SourceUnsaid            string                `json:"source_unsaid,omitempty"`
	TargetPerspective       string                `json:"target_perspective,omitempty"`
	TargetDesireInRelation  string                `json:"target_desire_in_relation,omitempty"`
	TargetFearInRelation    string                `json:"target_fear_in_relation,omitempty"`
	TargetUnsaid            string                `json:"target_unsaid,omitempty"`
	Phenomenology           RelationPhenomenology `json:"phenomenology"`
	RelationPalette         []RelationPaletteItem `json:"relation_palette"`
	SymbolicImages          []string              `json:"symbolic_images"`
	RelationKeywords        []string              `json:"relation_keywords"`
	RelationEvents          []RelationEvent       `json:"relation_events"`
	PrimarySong             *RelationSong         `json:"primary_song"`
	RelatedRelations        []RelationLink        `json:"related_relations"`
}

type CharacterDetailResponse struct {
	ID                   string                         `json:"id"`
	Slug                 string                         `json:"slug"`
	Name                 string                         `json:"name"`
	CoverURL             string                         `json:"coverUrl"`
	Summary              string                         `json:"summary"`
	OneLineDefinition    string                         `json:"oneLineDefinition"`
	CharacterTypeCode    string                         `json:"characterTypeCode"`
	Work                 *CharacterDetailRef            `json:"work"`
	Creator              *CharacterDetailRef            `json:"creator"`
	Song                 *CharacterDetailSong           `json:"song"`
	Songs                []CharacterDetailSong          `json:"songs"`
	CoreIdentity         string                         `json:"coreIdentity,omitempty"`
	PublicImage          string                         `json:"publicImage,omitempty"`
	HiddenSelf           string                         `json:"hiddenSelf,omitempty"`
	PrimaryMotivation    string                         `json:"primaryMotivation,omitempty"`
	CoreFear             string                         `json:"coreFear,omitempty"`
	PsychologicalWound   string                         `json:"psychologicalWound,omitempty"`
	CoreConflict         string                         `json:"coreConflict,omitempty"`
	EmotionalTone        string                         `json:"emotionalTone,omitempty"`
	Origin               string                         `json:"origin,omitempty"`
	FateArc              string                         `json:"fateArc,omitempty"`
	EndingState          string                         `json:"endingState,omitempty"`
	SurfaceTraits        []string                       `json:"surfaceTraits"`
	DeepTraits           []string                       `json:"deepTraits"`
	DominantEmotions     []string                       `json:"dominantEmotions"`
	SuppressedEmotions   []string                       `json:"suppressedEmotions"`
	ValuesTags           []string                       `json:"valuesTags"`
	DisplayTags          []string                       `json:"displayTags"`
	BottomLines          []string                       `json:"bottomLines"`
	Timeline             []CharacterTimelineItem        `json:"timeline"`
	RelationshipProfile  map[string]string              `json:"relationshipProfile"`
	RelationshipPatterns []CharacterRelationshipPattern `json:"relationshipPatterns"`
	Colors               []CharacterColorItem           `json:"colors"`
	SymbolicImages       []string                       `json:"symbolicImages"`
	Elements             []string                       `json:"elements"`
	SoundscapeKeywords   []string                       `json:"soundscapeKeywords"`
	KeyRelationships     []RelationshipListItemResponse `json:"keyRelationships"`
	SimilarCharacters    []CharacterListItemResponse    `json:"similarCharacters"`
}

type WorkListItemResponse struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	CoverURL       string `json:"coverUrl"`
	Summary        string `json:"summary"`
	WorkTypeCode   string `json:"workTypeCode"`
	CreatorName    string `json:"creatorName,omitempty"`
	CharacterCount int    `json:"characterCount"`
}

type WorkDetailResponse struct {
	ID             string                      `json:"id"`
	Slug           string                      `json:"slug"`
	Title          string                      `json:"title"`
	CoverURL       string                      `json:"coverUrl"`
	Summary        string                      `json:"summary"`
	WorkTypeCode   string                      `json:"workTypeCode"`
	Creator        *CharacterDetailRef         `json:"creator"`
	CharacterCount int                         `json:"characterCount"`
	Characters     []CharacterListItemResponse `json:"characters"`
}

type CreatorListItemResponse struct {
	ID              string `json:"id"`
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	CoverURL        string `json:"coverUrl"`
	Summary         string `json:"summary"`
	CreatorTypeCode string `json:"creatorTypeCode,omitempty"`
	EraText         string `json:"eraText,omitempty"`
	WorkCount       int    `json:"workCount"`
}

type CreatorDetailResponse struct {
	ID              string                 `json:"id"`
	Slug            string                 `json:"slug"`
	Name            string                 `json:"name"`
	CoverURL        string                 `json:"coverUrl"`
	Summary         string                 `json:"summary"`
	CreatorTypeCode string                 `json:"creatorTypeCode,omitempty"`
	EraText         string                 `json:"eraText,omitempty"`
	Works           []WorkListItemResponse `json:"works"`
}

type ThemeListItemResponse struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	CoverURL       string `json:"coverUrl"`
	Summary        string `json:"summary"`
	Category       string `json:"category,omitempty"`
	SubjectType    string `json:"subjectType,omitempty"`
	ItemCount      int    `json:"itemCount"`
	CharacterCount int    `json:"characterCount"`
}

type ThemeDetailResponse struct {
	ID            string                         `json:"id"`
	Slug          string                         `json:"slug"`
	Name          string                         `json:"name"`
	CoverURL      string                         `json:"coverUrl"`
	Summary       string                         `json:"summary"`
	Category      string                         `json:"category,omitempty"`
	SubjectType   string                         `json:"subjectType,omitempty"`
	Characters    []CharacterListItemResponse    `json:"characters"`
	Relationships []RelationshipListItemResponse `json:"relationships"`
}

type SongListItemResponse struct {
	ID            string   `json:"id"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	CharacterSlug string   `json:"characterSlug"`
	CoverURL      string   `json:"coverUrl"`
	AudioURL      string   `json:"audioUrl"`
	Styles        []string `json:"styles"`
}

type SearchResponseData struct {
	Characters []CharacterListItemResponse `json:"characters"`
	Works      []WorkListItemResponse      `json:"works"`
	Creators   []CreatorListItemResponse   `json:"creators"`
	Themes     []ThemeListItemResponse     `json:"themes"`
	Songs      []SongListItemResponse      `json:"songs"`
}
