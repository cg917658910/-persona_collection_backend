package dto

import "encoding/json"

type AdminRelationPaletteItem struct {
	Name string `json:"name,omitempty"`
	Hex  string `json:"hex"`
}

type AdminRelationPhenomenology struct {
	Body     string `json:"body,omitempty"`
	Time     string `json:"time,omitempty"`
	Space    string `json:"space,omitempty"`
	Gaze     string `json:"gaze,omitempty"`
	Language string `json:"language,omitempty"`
}

type AdminRelationEvent struct {
	StageNo      int    `json:"stageNo"`
	StageCode    string `json:"stageCode,omitempty"`
	Title        string `json:"title"`
	Summary      string `json:"summary,omitempty"`
	TensionShift string `json:"tensionShift,omitempty"`
	PowerShift   string `json:"powerShift,omitempty"`
	FateImpact   string `json:"fateImpact,omitempty"`
	SourceState  string `json:"sourceState,omitempty"`
	TargetState  string `json:"targetState,omitempty"`
	EventQuote   string `json:"eventQuote,omitempty"`
	ColorHex     string `json:"colorHex,omitempty"`
	SortOrder    int    `json:"sortOrder"`
}

type AdminRelationSong struct {
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Subtitle           string   `json:"subtitle,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	CoverURL           string   `json:"coverUrl,omitempty"`
	AudioURL           string   `json:"audioUrl,omitempty"`
	SongCoreTheme      string   `json:"songCoreTheme,omitempty"`
	SongEmotionalCurve string   `json:"songEmotionalCurve,omitempty"`
	SongStyles         []string `json:"songStyles,omitempty"`
	TempoBPM           int      `json:"tempoBpm,omitempty"`
	VocalProfile       string   `json:"vocalProfile,omitempty"`
	Lyric              string   `json:"lyric,omitempty"`
	Prompt             string   `json:"prompt,omitempty"`
	IsPrimary          bool     `json:"isPrimary"`
	SortOrder          int      `json:"sortOrder"`
	Status             string   `json:"status"`
}

type AdminRelationLink struct {
	LinkedRelationSlug string `json:"linkedRelationSlug"`
	LinkTypeCode       string `json:"linkTypeCode"`
	Reason             string `json:"reason,omitempty"`
	SortOrder          int    `json:"sortOrder"`
}

type AdminRelation struct {
	ID                     string                     `json:"id,omitempty"`
	Slug                   string                     `json:"slug"`
	Name                   string                     `json:"name"`
	Subtitle               string                     `json:"subtitle,omitempty"`
	Summary                string                     `json:"summary,omitempty"`
	OneLineDefinition      string                     `json:"oneLineDefinition,omitempty"`
	CoverURL               string                     `json:"coverUrl,omitempty"`
	Status                 string                     `json:"status"`
	SortOrder              int                        `json:"sortOrder"`
	RelationTypeCode       string                     `json:"relationTypeCode"`
	RelationTypeName       string                     `json:"relationTypeName,omitempty"`
	WorkSlug               string                     `json:"workSlug,omitempty"`
	WorkName               string                     `json:"workName,omitempty"`
	SourceCharacterSlug    string                     `json:"sourceCharacterSlug"`
	SourceCharacterName    string                     `json:"sourceCharacterName,omitempty"`
	TargetCharacterSlug    string                     `json:"targetCharacterSlug"`
	TargetCharacterName    string                     `json:"targetCharacterName,omitempty"`
	CoreDynamic            string                     `json:"coreDynamic,omitempty"`
	CoreTension            string                     `json:"coreTension,omitempty"`
	EmotionalTone          string                     `json:"emotionalTone,omitempty"`
	EmotionalTemperature   string                     `json:"emotionalTemperature,omitempty"`
	ConnectionTrigger      string                     `json:"connectionTrigger,omitempty"`
	SustainingMechanism    string                     `json:"sustainingMechanism,omitempty"`
	RelationConflict       string                     `json:"relationConflict,omitempty"`
	RelationArc            string                     `json:"relationArc,omitempty"`
	FateImpact             string                     `json:"fateImpact,omitempty"`
	PowerStructure         string                     `json:"powerStructure,omitempty"`
	DependencyPattern      string                     `json:"dependencyPattern,omitempty"`
	SourcePerspective      string                     `json:"sourcePerspective,omitempty"`
	SourceDesireInRelation string                     `json:"sourceDesireInRelation,omitempty"`
	SourceFearInRelation   string                     `json:"sourceFearInRelation,omitempty"`
	SourceUnsaid           string                     `json:"sourceUnsaid,omitempty"`
	TargetPerspective      string                     `json:"targetPerspective,omitempty"`
	TargetDesireInRelation string                     `json:"targetDesireInRelation,omitempty"`
	TargetFearInRelation   string                     `json:"targetFearInRelation,omitempty"`
	TargetUnsaid           string                     `json:"targetUnsaid,omitempty"`
	Phenomenology          AdminRelationPhenomenology `json:"phenomenology,omitempty"`
	SymbolicImages         []string                   `json:"symbolicImages,omitempty"`
	ThemeTags              []string                   `json:"themeTags,omitempty"`
	RelationPalette        []AdminRelationPaletteItem `json:"relationPalette,omitempty"`
	RelationKeywords       []string                   `json:"relationKeywords,omitempty"`
	TensionTags            []string                   `json:"tensionTags,omitempty"`
	CoverPrompt            string                     `json:"coverPrompt,omitempty"`
	SongPrompt             string                     `json:"songPrompt,omitempty"`
	PrimarySongSlug        string                     `json:"primarySongSlug,omitempty"`
	ThemeSlugs             []string                   `json:"themeSlugs,omitempty"`
	Events                 []AdminRelationEvent       `json:"events,omitempty"`
	Songs                  []AdminRelationSong        `json:"songs,omitempty"`
	Links                  []AdminRelationLink        `json:"links,omitempty"`
}

type RelationImportPackage struct {
	PmRelations []GeneratedRelation `json:"pm_relations"`
}

func (p *RelationImportPackage) UnmarshalJSON(data []byte) error {
	data = json.RawMessage(data)
	if len(data) == 0 || string(data) == "null" {
		p.PmRelations = nil
		return nil
	}
	var arr []GeneratedRelation
	if err := json.Unmarshal(data, &arr); err == nil {
		p.PmRelations = arr
		return nil
	}
	type wrapper struct {
		PmRelations []GeneratedRelation `json:"pm_relations"`
	}
	var w wrapper
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	p.PmRelations = w.PmRelations
	return nil
}

type AdminRelationImportRequest struct {
	Package RelationImportPackage `json:"package"`
}

type GeneratedRelation struct {
	Slug                   string                         `json:"slug"`
	Name                   string                         `json:"name"`
	Title                  string                         `json:"title"`
	Subtitle               string                         `json:"subtitle"`
	RelationTypeCode       string                         `json:"relation_type_code"`
	WorkSlug               string                         `json:"work_slug"`
	SourceCharacterSlug    string                         `json:"source_character_slug"`
	TargetCharacterSlug    string                         `json:"target_character_slug"`
	Summary                string                         `json:"summary"`
	OneLineDefinition      string                         `json:"one_line_definition"`
	CoreDynamic            string                         `json:"core_dynamic"`
	EmotionalTone          string                         `json:"emotional_tone"`
	EmotionalTemperature   string                         `json:"emotional_temperature"`
	ConnectionTrigger      string                         `json:"connection_trigger"`
	SustainingMechanism    string                         `json:"sustaining_mechanism"`
	RelationConflict       string                         `json:"relation_conflict"`
	CoreTension            string                         `json:"core_tension"`
	RelationArc            string                         `json:"relation_arc"`
	FateImpact             string                         `json:"fate_impact"`
	EndingDirection        string                         `json:"ending_direction"`
	PowerStructure         string                         `json:"power_structure"`
	DependencyPattern      string                         `json:"dependency_pattern"`
	SourcePerspective      string                         `json:"source_perspective"`
	TargetPerspective      string                         `json:"target_perspective"`
	SourceDesireInRelation string                         `json:"source_desire_in_relation"`
	SourceFearInRelation   string                         `json:"source_fear_in_relation"`
	SourceUnsaid           string                         `json:"source_unsaid"`
	TargetDesireInRelation string                         `json:"target_desire_in_relation"`
	TargetFearInRelation   string                         `json:"target_fear_in_relation"`
	TargetUnsaid           string                         `json:"target_unsaid"`
	Phenomenology          AdminRelationPhenomenology     `json:"phenomenology"`
	SymbolicImages         []string                       `json:"symbolic_images"`
	ThemeTags              []string                       `json:"theme_tags"`
	RelationPalette        []AdminRelationPaletteItem     `json:"relation_palette"`
	Palette                []AdminRelationPaletteItem     `json:"palette"`
	TensionTags            []string                       `json:"tension_tags"`
	CoverURL               string                         `json:"cover_url"`
	CoverPrompt            string                         `json:"cover_prompt"`
	SongPrompt             string                         `json:"song_prompt"`
	PrimarySongSlug        string                         `json:"primary_song_slug"`
	RelatedRelationSlugs   []string                       `json:"related_relation_slugs"`
	MirrorRelationSlugs    []string                       `json:"mirror_relation_slugs"`
	SameWorkRelationSlugs  []string                       `json:"same_work_relation_slugs"`
	HighlightQuotes        []string                       `json:"ui_highlight_quotes"`
	Meta                   map[string]any                 `json:"meta"`
	SortOrder              int                            `json:"sort_order"`
	Status                 string                         `json:"status"`
	IsActive               *bool                          `json:"is_active"`
	ThemeSlugs             []string                       `json:"theme_slugs"`
	RelationThemes         []string                       `json:"relation_themes"`
	Events                 []GeneratedRelationEvent       `json:"events"`
	Song                   *GeneratedRelationSong         `json:"song"`
	Songs                  []GeneratedRelationSong        `json:"songs"`
	Participants           []GeneratedRelationParticipant `json:"participants"`
	Links                  []GeneratedRelationLink        `json:"links"`
}

type GeneratedRelationParticipant struct {
	CharacterSlug      string         `json:"character_slug"`
	RoleCode           string         `json:"role_code"`
	RoleName           string         `json:"role_name"`
	PerspectiveSummary string         `json:"perspective_summary"`
	DesireInRelation   string         `json:"desire_in_relation"`
	FearInRelation     string         `json:"fear_in_relation"`
	Unsaid             string         `json:"unsaid"`
	SortOrder          int            `json:"sort_order"`
	Meta               map[string]any `json:"meta"`
}

type GeneratedRelationEvent struct {
	Stage        int            `json:"stage"`
	StageNo      int            `json:"stage_no"`
	StageCode    string         `json:"stage_code"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	TensionShift string         `json:"tension_shift"`
	PowerShift   string         `json:"power_shift"`
	FateImpact   string         `json:"fate_impact"`
	SourceState  string         `json:"source_state"`
	TargetState  string         `json:"target_state"`
	EventQuote   string         `json:"event_quote"`
	ColorHex     string         `json:"color_hex"`
	SortOrder    int            `json:"sort_order"`
	Meta         map[string]any `json:"meta"`
}

type GeneratedRelationSong struct {
	Slug               string         `json:"slug"`
	Title              string         `json:"title"`
	Subtitle           string         `json:"subtitle"`
	Summary            string         `json:"summary"`
	CoverURL           string         `json:"cover_url"`
	AudioURL           string         `json:"audio_url"`
	DurationSec        int            `json:"duration_sec"`
	SongCoreTheme      string         `json:"song_core_theme"`
	SongEmotionalCurve string         `json:"song_emotional_curve"`
	SongStyles         []string       `json:"song_styles"`
	TempoBPM           int            `json:"tempo_bpm"`
	VocalProfile       string         `json:"vocal_profile"`
	Lyric              string         `json:"lyric"`
	Prompt             string         `json:"prompt"`
	IsPrimary          bool           `json:"is_primary"`
	SortOrder          int            `json:"sort_order"`
	Status             string         `json:"status"`
	IsActive           *bool          `json:"is_active"`
	Meta               map[string]any `json:"meta"`
}

type GeneratedRelationLink struct {
	LinkedRelationSlug string `json:"linked_relation_slug"`
	LinkTypeCode       string `json:"link_type_code"`
	Reason             string `json:"reason"`
	SortOrder          int    `json:"sort_order"`
}
