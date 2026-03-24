package service

import (
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

func (s *CatalogService) Home() (dto.HomePayload, error) {
	data, err := s.repo.Home()
	if err != nil {
		return dto.HomePayload{}, err
	}
	data.FeaturedCharacter = s.normalizeCharacter(data.FeaturedCharacter)
	data.LatestCharacters = s.normalizeCharacters(data.LatestCharacters)
	data.FeaturedSongs = s.normalizeSongs(data.FeaturedSongs)
	data.RecommendedWorks = s.normalizeWorks(data.RecommendedWorks)
	data.Themes = s.normalizeThemes(data.Themes)
	return data, nil
}

func (s *CatalogService) RandomCharacter(theme string, exclude []string) (dto.Character, error) {
	item, err := s.repo.RandomCharacter(theme, exclude)
	if err != nil {
		return dto.Character{}, err
	}
	return s.normalizeCharacter(item), nil
}

func (s *CatalogService) ListCharacters() ([]dto.Character, error) {
	list, err := s.repo.ListCharacters()
	if err != nil {
		return nil, err
	}
	return s.normalizeCharacters(list), nil
}

func (s *CatalogService) GetCharacterDetail(slug string) (dto.CharacterDetail, error) {
	item, err := s.repo.GetCharacterDetail(slug)
	if err != nil {
		return dto.CharacterDetail{}, err
	}
	item.Character = s.normalizeCharacter(item.Character)
	item.RelatedWorks = s.normalizeWorks(item.RelatedWorks)
	item.RelatedThemes = s.normalizeThemes(item.RelatedThemes)
	item.RelatedSongs = s.normalizeSongs(item.RelatedSongs)
	item.RelatedCreator = s.normalizeCreators(item.RelatedCreator)
	return item, nil
}

func (s *CatalogService) ListWorks() ([]dto.Work, error) {
	list, err := s.repo.ListWorks()
	if err != nil {
		return nil, err
	}
	return s.normalizeWorks(list), nil
}

func (s *CatalogService) GetWorkDetail(slug string) (dto.Work, error) {
	item, err := s.repo.GetWorkDetail(slug)
	if err != nil {
		return dto.Work{}, err
	}
	return s.normalizeWork(item), nil
}

func (s *CatalogService) ListCreators() ([]dto.Creator, error) {
	list, err := s.repo.ListCreators()
	if err != nil {
		return nil, err
	}
	return s.normalizeCreators(list), nil
}

func (s *CatalogService) GetCreatorDetail(slug string) (dto.Creator, error) {
	item, err := s.repo.GetCreatorDetail(slug)
	if err != nil {
		return dto.Creator{}, err
	}
	return s.normalizeCreator(item), nil
}

func (s *CatalogService) ListThemes() ([]dto.Theme, error) {
	list, err := s.repo.ListThemes()
	if err != nil {
		return nil, err
	}
	return s.normalizeThemes(list), nil
}

func (s *CatalogService) GetThemeDetail(slug string) (dto.ThemeDetail, error) {
	item, err := s.repo.GetThemeDetail(slug)
	if err != nil {
		return dto.ThemeDetail{}, err
	}
	item.Theme = s.normalizeTheme(item.Theme)
	item.Characters = s.normalizeCharacters(item.Characters)
	return item, nil
}

func (s *CatalogService) ListSongs() ([]dto.Song, error) {
	list, err := s.repo.ListSongs()
	if err != nil {
		return nil, err
	}
	return s.normalizeSongs(list), nil
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
	return in
}

func (s *CatalogService) normalizeSongs(in []dto.Song) []dto.Song {
	out := make([]dto.Song, 0, len(in))
	for _, v := range in {
		out = append(out, s.normalizeSong(v))
	}
	return out
}

func (s *CatalogService) normalizeSong(in dto.Song) dto.Song {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.AudioURL = s.assets.Normalize(in.AudioURL)
	return in
}
