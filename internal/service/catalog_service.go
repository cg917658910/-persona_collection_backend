package service

import (
	"sort"
	"strings"

	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type CatalogService struct {
	repo   repo.CatalogRepo
	assets *AssetURLBuilder
}

func NewCatalogService(r repo.CatalogRepo, assets *AssetURLBuilder) *CatalogService {
	return &CatalogService{repo: r, assets: assets}
}

func (s *CatalogService) Home() (dto.HomeResponseData, error) {
	home, err := s.repo.Home()
	if err != nil {
		return dto.HomeResponseData{}, err
	}
	chars, err := s.repo.ListCharacters()
	if err != nil {
		return dto.HomeResponseData{}, err
	}
	works, _ := s.repo.ListWorks()
	creators, _ := s.repo.ListCreators()
	themes, _ := s.repo.ListThemes()

	out := dto.HomeResponseData{
		LatestCharacters:  make([]dto.HomeCharacterCard, 0, len(home.LatestCharacters)),
		FeaturedSongs:     make([]dto.HomeSongCard, 0, len(home.FeaturedSongs)),
		RecommendedWorks:  make([]dto.HomeWorkCard, 0, len(home.RecommendedWorks)),
		RecommendedThemes: make([]dto.HomeThemeCard, 0, len(home.Themes)),
		CategoryCounts:    s.buildCategoryCounts(chars, works, creators, themes),
	}
	if home.FeaturedCharacter.Slug != "" {
		featuredDetail, _ := s.repo.GetCharacterDetail(home.FeaturedCharacter.Slug)
		var songRef *dto.HomeSongRef
		if len(featuredDetail.RelatedSongs) > 0 {
			song := s.normalizeSong(featuredDetail.RelatedSongs[0])
			songRef = &dto.HomeSongRef{
				ID: song.Slug, Slug: song.Slug, Title: song.Title, CoverURL: song.CoverURL,
				AudioURL: song.AudioURL, CharacterSlug: song.CharacterSlug, CharacterName: featuredDetail.Name,
			}
		}
		workTitle := ""
		if len(featuredDetail.RelatedWorks) > 0 {
			workTitle = featuredDetail.RelatedWorks[0].Title
		}
		out.FeaturedCharacter = &dto.HomeFeaturedCharacter{
			ID: home.FeaturedCharacter.Slug, Slug: home.FeaturedCharacter.Slug, Name: home.FeaturedCharacter.Name,
			CoverURL: s.assets.Normalize(home.FeaturedCharacter.CoverURL),
			Summary:  home.FeaturedCharacter.Summary, OneLineDefinition: home.FeaturedCharacter.OneLineDefinition,
			CharacterTypeCode: home.FeaturedCharacter.CharacterTypeCode, WorkTitle: workTitle,
			Tags: s.characterTags(featuredDetail), Song: songRef,
		}
	}
	for _, item := range home.LatestCharacters {
		detail, _ := s.repo.GetCharacterDetail(item.Slug)
		card := dto.HomeCharacterCard{
			ID: item.Slug, Slug: item.Slug, Name: item.Name,
			CoverURL: s.assets.Normalize(item.CoverURL),
			Summary:  item.Summary, OneLineDefinition: item.OneLineDefinition,
			CharacterTypeCode: item.CharacterTypeCode,
			WorkTitle:         s.firstWorkTitle(detail.RelatedWorks),
			Tags:              s.characterTags(detail),
			HasSong:           len(detail.RelatedSongs) > 0,
		}
		out.LatestCharacters = append(out.LatestCharacters, card)
	}
	for _, item := range home.FeaturedSongs {
		song := s.normalizeSong(item)
		out.FeaturedSongs = append(out.FeaturedSongs, dto.HomeSongCard{
			ID: song.Slug, Slug: song.Slug, Title: song.Title, CoverURL: song.CoverURL,
			AudioURL: song.AudioURL, CharacterSlug: song.CharacterSlug, CharacterName: s.characterName(chars, song.CharacterSlug),
			Summary: song.Summary, SongCoreTheme: song.SongCoreTheme,
		})
	}
	for _, item := range home.RecommendedWorks {
		work := s.normalizeWork(item)
		out.RecommendedWorks = append(out.RecommendedWorks, dto.HomeWorkCard{
			ID: work.Slug, Slug: work.Slug, Title: work.Title, CoverURL: work.CoverURL,
			Summary: work.Summary, WorkTypeCode: work.TypeCode,
		})
	}
	for _, item := range home.Themes {
		themeDetail, _ := s.repo.GetThemeDetail(item.Slug)
		theme := s.normalizeTheme(item)
		itemCount := len(themeDetail.Characters)
		if normalizeThemeSubjectType(theme.SubjectType) == "relation" {
			itemCount = len(themeDetail.Relationships)
		}
		out.RecommendedThemes = append(out.RecommendedThemes, dto.HomeThemeCard{
			ID: theme.Slug, Slug: theme.Slug, Name: theme.Name, CoverURL: theme.CoverURL,
			Summary: theme.Summary, CharacterCount: itemCount,
		})
	}
	return out, nil
}

func (s *CatalogService) RandomCharacter(theme string, exclude []string) (dto.HomeFeaturedCharacter, error) {
	item, err := s.repo.RandomCharacter(theme, exclude)
	if err != nil {
		return dto.HomeFeaturedCharacter{}, err
	}
	detail, _ := s.repo.GetCharacterDetail(item.Slug)
	var songRef *dto.HomeSongRef
	if len(detail.RelatedSongs) > 0 {
		song := s.normalizeSong(detail.RelatedSongs[0])
		songRef = &dto.HomeSongRef{
			ID: song.Slug, Slug: song.Slug, Title: song.Title, CoverURL: song.CoverURL,
			AudioURL: song.AudioURL, CharacterSlug: song.CharacterSlug, CharacterName: detail.Name,
		}
	}
	return dto.HomeFeaturedCharacter{
		ID: item.Slug, Slug: item.Slug, Name: item.Name, CoverURL: s.assets.Normalize(item.CoverURL),
		Summary: item.Summary, OneLineDefinition: item.OneLineDefinition, CharacterTypeCode: item.CharacterTypeCode,
		WorkTitle: s.firstWorkTitle(detail.RelatedWorks), Tags: s.characterTags(detail), Song: songRef,
	}, nil
}

func (s *CatalogService) ListCharacters() (dto.CharacterListResponse, error) {
	list, err := s.repo.ListCharacters()
	if err != nil {
		return dto.CharacterListResponse{}, err
	}
	items := make([]dto.CharacterListItemResponse, 0, len(list))
	for _, item := range list {
		items = append(items, s.buildCharacterListItem(item))
	}
	return dto.CharacterListResponse{Items: items, Total: len(items), Page: 1, PageSize: len(items)}, nil
}

func (s *CatalogService) GetCharacterDetail(slug string) (dto.CharacterDetailResponse, error) {
	item, err := s.repo.GetCharacterDetail(slug)
	if err != nil {
		return dto.CharacterDetailResponse{}, err
	}
	item = s.normalizeCharacterDetail(item)

	var workRef *dto.CharacterDetailRef
	if len(item.RelatedWorks) > 0 {
		w := item.RelatedWorks[0]
		workRef = &dto.CharacterDetailRef{Slug: w.Slug, Title: w.Title, Name: w.Title, CoverURL: w.CoverURL, Summary: w.Summary}
	}
	var creatorRef *dto.CharacterDetailRef
	if len(item.RelatedCreator) > 0 {
		c := item.RelatedCreator[0]
		creatorRef = &dto.CharacterDetailRef{Slug: c.Slug, Title: c.Name, Name: c.Name, CoverURL: c.CoverURL, Summary: c.Summary}
	}
	songs := make([]dto.CharacterDetailSong, 0, len(item.RelatedSongs))
	for _, rawSong := range item.RelatedSongs {
		songs = append(songs, s.toCharacterDetailSong(rawSong))
	}
	var songRef *dto.CharacterDetailSong
	if len(songs) > 0 {
		songRef = &songs[0]
	}
	similar := s.buildSimilarCharacters(item)
	return dto.CharacterDetailResponse{
		ID:   firstNonEmpty(item.ID, item.Slug),
		Slug: item.Slug, Name: item.Name, CoverURL: item.CoverURL, Summary: item.Summary,
		OneLineDefinition: item.OneLineDefinition, CharacterTypeCode: item.CharacterTypeCode,
		Work: workRef, Creator: creatorRef, Song: songRef, Songs: songs,
		CoreIdentity:         item.CoreIdentity,
		PublicImage:          item.PublicImage,
		HiddenSelf:           item.HiddenSelf,
		PrimaryMotivation:    item.PrimaryMotivation,
		CoreFear:             item.CoreFear,
		PsychologicalWound:   item.PsychologicalWound,
		CoreConflict:         item.CoreConflict,
		EmotionalTone:        item.EmotionalTone,
		Origin:               item.Origin,
		FateArc:              item.FateArc,
		EndingState:          item.EndingState,
		SurfaceTraits:        s.sliceOrDefault(item.SurfaceTraits, []string{}),
		DeepTraits:           s.sliceOrDefault(item.DeepTraits, []string{}),
		DominantEmotions:     s.sliceOrDefault(item.DominantEmotions, []string{}),
		SuppressedEmotions:   s.sliceOrDefault(item.SuppressedEmotions, []string{}),
		ValuesTags:           s.sliceOrDefault(item.ValuesTags, []string{}),
		DisplayTags:          s.characterTags(item),
		BottomLines:          s.sliceOrDefault(item.BottomLines, []string{}),
		Timeline:             s.timelineItemsV2(item.Timeline),
		RelationshipProfile:  s.localizedRelationshipProfile(item.RelationshipProfile),
		RelationshipPatterns: s.localizedRelationshipPatterns(item.RelationshipProfile),
		Colors:               s.colorItems(item.Colors),
		SymbolicImages:       s.sliceOrDefault(item.SymbolicImages, []string{}),
		Elements:             s.sliceOrDefault(item.Elements, []string{}),
		SoundscapeKeywords:   s.sliceOrDefault(item.SoundscapeKeywords, []string{}),
		KeyRelationships:     s.buildKeyRelationships(item.Slug, 4),
		SimilarCharacters:    similar,
	}, nil
}

func (s *CatalogService) ListWorks() ([]dto.WorkListItemResponse, error) {
	list, err := s.repo.ListWorks()
	if err != nil {
		return nil, err
	}
	out := make([]dto.WorkListItemResponse, 0, len(list))
	for _, item := range list {
		work := s.normalizeWork(item)
		out = append(out, dto.WorkListItemResponse{
			ID: work.Slug, Slug: work.Slug, Title: work.Title, CoverURL: work.CoverURL,
			Summary: work.Summary, WorkTypeCode: work.TypeCode,
			CreatorName:    strings.Join(work.CreatorNames, " / "),
			CharacterCount: len(work.CharacterSlugs),
		})
	}
	return out, nil
}

func (s *CatalogService) GetWorkDetail(slug string) (dto.WorkDetailResponse, error) {
	item, err := s.repo.GetWorkDetail(slug)
	if err != nil {
		return dto.WorkDetailResponse{}, err
	}
	work := s.normalizeWork(item)
	allChars, _ := s.repo.ListCharacters()
	characters := make([]dto.CharacterListItemResponse, 0)
	for _, c := range allChars {
		if contains(c.WorkSlugs, slug) {
			characters = append(characters, s.buildCharacterListItem(c))
		}
	}
	var creator *dto.CharacterDetailRef
	if len(work.CreatorSlugs) > 0 {
		creator = &dto.CharacterDetailRef{Slug: work.CreatorSlugs[0], Name: firstOrEmpty(work.CreatorNames), Title: firstOrEmpty(work.CreatorNames)}
	}
	return dto.WorkDetailResponse{
		ID: work.Slug, Slug: work.Slug, Title: work.Title, CoverURL: work.CoverURL,
		Summary: work.Summary, WorkTypeCode: work.TypeCode, Creator: creator,
		CharacterCount: len(characters), Characters: characters,
	}, nil
}

func (s *CatalogService) ListCreators() ([]dto.CreatorListItemResponse, error) {
	list, err := s.repo.ListCreators()
	if err != nil {
		return nil, err
	}
	out := make([]dto.CreatorListItemResponse, 0, len(list))
	for _, item := range list {
		creator := s.normalizeCreator(item)
		out = append(out, dto.CreatorListItemResponse{
			ID: creator.Slug, Slug: creator.Slug, Name: creator.Name, CoverURL: creator.CoverURL,
			Summary: creator.Summary, CreatorTypeCode: creator.CreatorTypeCode, EraText: creator.EraText,
			WorkCount: len(creator.WorkSlugs),
		})
	}
	return out, nil
}

func (s *CatalogService) GetCreatorDetail(slug string) (dto.CreatorDetailResponse, error) {
	item, err := s.repo.GetCreatorDetail(slug)
	if err != nil {
		return dto.CreatorDetailResponse{}, err
	}
	creator := s.normalizeCreator(item)
	allWorks, _ := s.repo.ListWorks()
	works := make([]dto.WorkListItemResponse, 0)
	for _, w := range allWorks {
		if contains(creator.WorkSlugs, w.Slug) {
			w = s.normalizeWork(w)
			works = append(works, dto.WorkListItemResponse{
				ID: w.Slug, Slug: w.Slug, Title: w.Title, CoverURL: w.CoverURL, Summary: w.Summary,
				WorkTypeCode: w.TypeCode, CreatorName: strings.Join(w.CreatorNames, " / "), CharacterCount: len(w.CharacterSlugs),
			})
		}
	}
	return dto.CreatorDetailResponse{
		ID: creator.Slug, Slug: creator.Slug, Name: creator.Name, CoverURL: creator.CoverURL, Summary: creator.Summary,
		CreatorTypeCode: creator.CreatorTypeCode, EraText: creator.EraText, Works: works,
	}, nil
}

func (s *CatalogService) ListThemes() ([]dto.ThemeListItemResponse, error) {
	list, err := s.repo.ListThemes()
	if err != nil {
		return nil, err
	}
	out := make([]dto.ThemeListItemResponse, 0, len(list))
	for _, item := range list {
		detail, _ := s.repo.GetThemeDetail(item.Slug)
		theme := s.normalizeTheme(item)
		itemCount := len(detail.Characters)
		if normalizeThemeSubjectType(theme.SubjectType) == "relation" {
			itemCount = len(detail.Relationships)
		}
		out = append(out, dto.ThemeListItemResponse{
			ID: theme.Slug, Slug: theme.Slug, Name: theme.Name, CoverURL: theme.CoverURL,
			Summary: theme.Summary, Category: theme.Category, SubjectType: theme.SubjectType,
			ItemCount: itemCount, CharacterCount: itemCount,
		})
	}
	return out, nil
}

func (s *CatalogService) GetThemeDetail(slug string) (dto.ThemeDetailResponse, error) {
	item, err := s.repo.GetThemeDetail(slug)
	if err != nil {
		return dto.ThemeDetailResponse{}, err
	}
	item.Theme = s.normalizeTheme(item.Theme)
	item.Characters = s.normalizeCharacters(item.Characters)
	item.Relationships = s.normalizeRelationRecords(item.Relationships)
	allCharacters, _ := s.repo.ListCharacters()
	characterBySlug := make(map[string]dto.Character, len(allCharacters))
	for _, current := range allCharacters {
		characterBySlug[current.Slug] = current
	}
	chars := make([]dto.CharacterListItemResponse, 0, len(item.Characters))
	for _, c := range item.Characters {
		if full, ok := characterBySlug[c.Slug]; ok {
			chars = append(chars, s.buildCharacterListItem(full))
			continue
		}
		chars = append(chars, s.buildCharacterListItem(c))
	}
	relationships := make([]dto.RelationshipListItemResponse, 0, len(item.Relationships))
	for _, relation := range item.Relationships {
		relationships = append(relationships, s.buildRelationshipListItem(relation, ""))
	}
	return dto.ThemeDetailResponse{
		ID: item.Theme.Slug, Slug: item.Theme.Slug, Name: item.Theme.Name, CoverURL: item.Theme.CoverURL,
		Summary: item.Theme.Summary, Category: item.Theme.Category, SubjectType: item.Theme.SubjectType,
		Characters: chars, Relationships: relationships,
	}, nil
}

func (s *CatalogService) ListSongs() ([]dto.SongListItemResponse, error) {
	list, err := s.repo.ListSongs()
	if err != nil {
		return nil, err
	}
	out := make([]dto.SongListItemResponse, 0, len(list))
	for _, item := range list {
		song := s.normalizeSong(item)
		out = append(out, dto.SongListItemResponse{
			ID: song.Slug, Slug: song.Slug, Title: song.Title, CharacterSlug: song.CharacterSlug,
			CoverURL: song.CoverURL, AudioURL: song.AudioURL, Styles: song.Styles,
		})
	}
	return out, nil
}

func (s *CatalogService) SearchCatalog(keyword string, limit int) (dto.SearchResponseData, error) {
	out, err := s.repo.SearchCatalog(keyword, limit)
	if err != nil {
		return dto.SearchResponseData{}, err
	}
	allCharacters, listErr := s.repo.ListCharacters()
	characterBySlug := make(map[string]dto.Character, len(allCharacters))
	if listErr == nil {
		for _, item := range allCharacters {
			characterBySlug[item.Slug] = item
		}
	}
	for i := range out.Characters {
		if item, ok := characterBySlug[out.Characters[i].Slug]; ok {
			out.Characters[i] = s.buildCharacterListItem(item)
			continue
		}
		out.Characters[i].CoverURL = s.assets.Normalize(out.Characters[i].CoverURL)
		if out.Characters[i].HasSong && strings.TrimSpace(out.Characters[i].ThemeSongTitle) == "" {
			out.Characters[i].ThemeSongTitle = "人物之歌"
		}
	}
	for i := range out.Works {
		out.Works[i].CoverURL = s.assets.Normalize(out.Works[i].CoverURL)
	}
	for i := range out.Creators {
		out.Creators[i].CoverURL = s.assets.Normalize(out.Creators[i].CoverURL)
	}
	for i := range out.Themes {
		out.Themes[i].CoverURL = s.assets.Normalize(out.Themes[i].CoverURL)
	}
	for i := range out.Songs {
		out.Songs[i].CoverURL = s.assets.Normalize(out.Songs[i].CoverURL)
		out.Songs[i].AudioURL = s.assets.Normalize(out.Songs[i].AudioURL)
	}
	return out, nil
}

func (s *CatalogService) buildCategoryCounts(chars []dto.Character, works []dto.Work, creators []dto.Creator, themes []dto.Theme) dto.HomeCategoryCounts {
	out := dto.HomeCategoryCounts{Characters: len(chars), Creators: len(creators), Works: len(works), Themes: len(themes)}
	for _, item := range chars {
		switch strings.TrimSpace(item.CharacterTypeCode) {
		case "historical":
			out.Historical++
		case "literary":
			out.Literary++
		case "film_tv":
			out.FilmTV++
		case "anime":
			out.Anime++
		}
	}
	return out
}

func (s *CatalogService) characterTags(item dto.CharacterDetail) []string {
	return s.characterCardTags(firstThemeName(item.RelatedThemes), item.SurfaceTraits)
}

func (s *CatalogService) characterCardTags(primaryTheme string, surfaceTraits []string) []string {
	out := make([]string, 0, 3)
	if strings.TrimSpace(primaryTheme) != "" {
		out = append(out, primaryTheme)
	}
	for _, trait := range surfaceTraits {
		if strings.TrimSpace(trait) != "" {
			out = append(out, trait)
		}
		if len(out) >= 3 {
			break
		}
	}
	return s.sliceCompact(out)
}

func (s *CatalogService) buildCharacterListItem(item dto.Character) dto.CharacterListItemResponse {
	themeSongTitle := strings.TrimSpace(item.PrimarySongTitle)
	if themeSongTitle == "" && len(item.SongSlugs) > 0 {
		themeSongTitle = "人物之歌"
	}
	return dto.CharacterListItemResponse{
		ID:                item.Slug,
		Slug:              item.Slug,
		Name:              item.Name,
		CoverURL:          s.assets.Normalize(item.CoverURL),
		Summary:           item.Summary,
		OneLineDefinition: item.OneLineDefinition,
		CharacterTypeCode: item.CharacterTypeCode,
		WorkTitle:         item.PrimaryWorkTitle,
		Tags:              s.characterCardTags(item.PrimaryThemeName, item.SurfaceTraits),
		HasSong:           len(item.SongSlugs) > 0,
		ThemeSongTitle:    themeSongTitle,
	}
}

func (s *CatalogService) toCharacterDetailSong(song dto.Song) dto.CharacterDetailSong {
	song = s.normalizeSong(song)
	return dto.CharacterDetailSong{
		Slug:           song.Slug,
		Title:          song.Title,
		CoverURL:       song.CoverURL,
		AudioURL:       song.AudioURL,
		SongCoreTheme:  song.SongCoreTheme,
		EmotionalCurve: s.sliceOrDefault(song.EmotionalCurve, []string{}),
		SongStyles:     s.sliceOrDefault(song.Styles, []string{}),
		VocalProfile:   song.VocalProfile,
		Lyrics:         []string{},
	}
}

func (s *CatalogService) firstWorkTitle(works []dto.Work) string {
	if len(works) == 0 {
		return ""
	}
	return works[0].Title
}

func (s *CatalogService) characterName(chars []dto.Character, slug string) string {
	for _, item := range chars {
		if item.Slug == slug {
			return item.Name
		}
	}
	return slug
}

func (s *CatalogService) timelineItems(in []string) []dto.CharacterTimelineItem {
	out := make([]dto.CharacterTimelineItem, 0, len(in))
	for i, item := range in {
		out = append(out, dto.CharacterTimelineItem{
			Year:    string(rune('0' + (i+1)%10)),
			Event:   item,
			Emotion: item,
		})
	}
	return out
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func normalizeThemeSubjectType(subjectType string) string {
	switch strings.ToLower(strings.TrimSpace(subjectType)) {
	case "relation":
		return "relation"
	default:
		return "character"
	}
}

func firstOrEmpty(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[0]
}

func firstThemeName(themes []dto.Theme) string {
	if len(themes) == 0 {
		return ""
	}
	return themes[0].Name
}

func (s *CatalogService) sliceCompact(in []string) []string {
	out := make([]string, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, item := range in {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func (s *CatalogService) sliceOrDefault(in []string, def []string) []string {
	if len(in) == 0 {
		return def
	}
	return in
}

func (s *CatalogService) timelineItemsV2(in []dto.CharacterTimelineItem) []dto.CharacterTimelineItem {
	if len(in) == 0 {
		return []dto.CharacterTimelineItem{}
	}
	return in
}

func (s *CatalogService) localizedRelationshipProfile(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for _, item := range s.localizedRelationshipPatterns(in) {
		out[item.Title] = item.Value
	}
	return out
}

func (s *CatalogService) localizedRelationshipPatterns(in map[string]string) []dto.CharacterRelationshipPattern {
	if len(in) == 0 {
		return []dto.CharacterRelationshipPattern{}
	}
	ordered := []struct {
		key   string
		title string
	}{
		{key: "toward_love", title: "对爱情"},
		{key: "love", title: "对爱情"},
		{key: "toward_authority", title: "对权威"},
		{key: "authority", title: "对权威"},
		{key: "toward_friends", title: "对朋友"},
		{key: "friends", title: "对朋友"},
		{key: "toward_enemies", title: "对敌人"},
		{key: "enemies", title: "对敌人"},
	}
	out := make([]dto.CharacterRelationshipPattern, 0, len(in))
	used := make(map[string]struct{}, len(in))
	for _, item := range ordered {
		if value := strings.TrimSpace(in[item.key]); value != "" {
			out = append(out, dto.CharacterRelationshipPattern{Title: item.title, Value: value})
			used[item.key] = struct{}{}
		}
	}
	extraKeys := make([]string, 0, len(in))
	for key, value := range in {
		if _, ok := used[key]; ok || strings.TrimSpace(value) == "" {
			continue
		}
		extraKeys = append(extraKeys, key)
	}
	sort.Strings(extraKeys)
	for _, key := range extraKeys {
		out = append(out, dto.CharacterRelationshipPattern{
			Title: s.relationshipLabel(key),
			Value: strings.TrimSpace(in[key]),
		})
	}
	return out
}

func (s *CatalogService) relationshipLabel(key string) string {
	switch strings.TrimSpace(key) {
	case "toward_love", "love":
		return "对爱情"
	case "toward_authority", "authority":
		return "对权威"
	case "toward_friends", "friends":
		return "对朋友"
	case "toward_enemies", "enemies":
		return "对敌人"
	default:
		return strings.TrimSpace(key)
	}
}

func (s *CatalogService) mapOrDefault(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	return in
}

func (s *CatalogService) colorItems(in []dto.CharacterColorItem) []dto.CharacterColorItem {
	if len(in) == 0 {
		return []dto.CharacterColorItem{}
	}
	return in
}

func (s *CatalogService) normalizeCharacterDetail(in dto.CharacterDetail) dto.CharacterDetail {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.RelatedWorks = s.normalizeWorks(in.RelatedWorks)
	in.RelatedThemes = s.normalizeThemes(in.RelatedThemes)
	in.RelatedSongs = s.normalizeSongs(in.RelatedSongs)
	in.RelatedCreator = s.normalizeCreators(in.RelatedCreator)
	return in
}

func (s *CatalogService) buildSimilarCharacters(item dto.CharacterDetail) []dto.CharacterListItemResponse {
	all, err := s.repo.ListCharacters()
	if err != nil {
		return []dto.CharacterListItemResponse{}
	}

	primaryThemeSlug := ""
	if len(item.RelatedThemes) > 0 {
		primaryThemeSlug = strings.TrimSpace(item.RelatedThemes[0].Slug)
	}
	if primaryThemeSlug == "" && len(item.ThemeSlugs) > 0 {
		primaryThemeSlug = strings.TrimSpace(item.ThemeSlugs[0])
	}
	if primaryThemeSlug == "" {
		return []dto.CharacterListItemResponse{}
	}

	out := make([]dto.CharacterListItemResponse, 0, 4)
	seen := map[string]struct{}{item.Slug: {}}
	for _, c := range all {
		if _, ok := seen[c.Slug]; ok {
			continue
		}

		matched := false
		for _, ts := range c.ThemeSlugs {
			if strings.TrimSpace(ts) == primaryThemeSlug {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		seen[c.Slug] = struct{}{}
		out = append(out, s.buildCharacterListItem(c))
		if len(out) >= 4 {
			break
		}
	}
	return out
}

func firstNonEmpty(items ...string) string {
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			return item
		}
	}
	return ""
}

func (s *CatalogService) normalizeCharacters(in []dto.Character) []dto.Character {
	out := make([]dto.Character, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeCharacter(v))
	}
	return out
}

func (s *CatalogService) normalizeCharacter(in dto.Character) dto.Character {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	return in
}

func (s *CatalogService) normalizeWorks(in []dto.Work) []dto.Work {
	out := make([]dto.Work, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeWork(v))
	}
	return out
}

func (s *CatalogService) normalizeWork(in dto.Work) dto.Work {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	return in
}

func (s *CatalogService) normalizeCreators(in []dto.Creator) []dto.Creator {
	out := make([]dto.Creator, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeCreator(v))
	}
	return out
}

func (s *CatalogService) normalizeCreator(in dto.Creator) dto.Creator {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	return in
}

func (s *CatalogService) normalizeThemes(in []dto.Theme) []dto.Theme {
	out := make([]dto.Theme, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeTheme(v))
	}
	return out
}

func (s *CatalogService) normalizeTheme(in dto.Theme) dto.Theme {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.SubjectType = normalizeThemeSubjectType(in.SubjectType)
	return in
}

func (s *CatalogService) normalizeSongs(in []dto.Song) []dto.Song {
	out := make([]dto.Song, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeSong(v))
	}
	return out
}

func (s *CatalogService) normalizeRelationRecords(in []dto.RelationRecord) []dto.RelationRecord {
	out := make([]dto.RelationRecord, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeRelationRecord(v))
	}
	return out
}

func (s *CatalogService) normalizeSong(in dto.Song) dto.Song {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.AudioURL = s.assets.Normalize(in.AudioURL)
	return in
}
