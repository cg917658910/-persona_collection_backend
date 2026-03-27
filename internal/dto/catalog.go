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
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	CoverURL string `json:"coverUrl"`
	Category string `json:"category,omitempty"`
}

type ThemeDetail struct {
	Theme
	Characters []Character `json:"characters"`
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
	CharacterCount int    `json:"characterCount"`
}

type ThemeDetailResponse struct {
	ID         string                      `json:"id"`
	Slug       string                      `json:"slug"`
	Name       string                      `json:"name"`
	CoverURL   string                      `json:"coverUrl"`
	Summary    string                      `json:"summary"`
	Category   string                      `json:"category,omitempty"`
	Characters []CharacterListItemResponse `json:"characters"`
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
