package service

import (
	"strings"

	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type AdminPageService struct {
	repo repo.AdminRepo
}

func NewAdminPageService(r repo.AdminRepo) *AdminPageService {
	return &AdminPageService{repo: r}
}

func paginateSlice[T any](items []T, page, pageSize int) dto.PageResult[T] {
	page, pageSize = dto.NormalizePage(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return dto.PageResult[T]{Items: items[start:end], Total: total, Page: page, PageSize: pageSize}
}

func containsKeyword(s string, keyword string) bool {
	if keyword == "" {
		return true
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(keyword))
}

func (s *AdminPageService) Characters(q dto.PageQuery) (dto.PageResult[dto.AdminCharacter], error) {
	if r, ok := s.repo.(repo.AdminCharacterPageRepo); ok {
		return r.PageAdminCharacters(q)
	}
	list, err := s.repo.ListAdminCharacters()
	if err != nil {
		return dto.PageResult[dto.AdminCharacter]{}, err
	}
	filtered := make([]dto.AdminCharacter, 0)
	for _, item := range list {
		if !containsKeyword(item.Name+item.Slug+item.Summary+item.OneLineDefinition, q.Keyword) {
			continue
		}
		if q.CharacterTypeCode != "" && item.CharacterTypeCode != q.CharacterTypeCode {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateSlice(filtered, q.Page, q.PageSize), nil
}
func (s *AdminPageService) Songs(q dto.PageQuery) (dto.PageResult[dto.AdminSong], error) {
	if r, ok := s.repo.(repo.AdminSongPageRepo); ok {
		return r.PageAdminSongs(q)
	}
	list, err := s.repo.ListAdminSongs()
	if err != nil {
		return dto.PageResult[dto.AdminSong]{}, err
	}
	filtered := make([]dto.AdminSong, 0)
	for _, item := range list {
		if !containsKeyword(item.Title+item.Slug+item.Summary+item.CharacterName+item.CoreTheme, q.Keyword) {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateSlice(filtered, q.Page, q.PageSize), nil
}
func (s *AdminPageService) Themes(q dto.PageQuery) (dto.PageResult[dto.AdminTheme], error) {
	if r, ok := s.repo.(repo.AdminThemePageRepo); ok {
		return r.PageAdminThemes(q)
	}
	list, err := s.repo.ListAdminThemes()
	if err != nil {
		return dto.PageResult[dto.AdminTheme]{}, err
	}
	filtered := make([]dto.AdminTheme, 0)
	for _, item := range list {
		if !containsKeyword(item.Name+item.Slug+item.Summary+item.Code+item.Category, q.Keyword) {
			continue
		}
		if q.Category != "" && item.Category != q.Category {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateSlice(filtered, q.Page, q.PageSize), nil
}
func (s *AdminPageService) Works(q dto.PageQuery) (dto.PageResult[dto.AdminWork], error) {
	if r, ok := s.repo.(repo.AdminWorkPageRepo); ok {
		return r.PageAdminWorks(q)
	}
	list, err := s.repo.ListAdminWorks()
	if err != nil {
		return dto.PageResult[dto.AdminWork]{}, err
	}
	filtered := make([]dto.AdminWork, 0)
	for _, item := range list {
		if !containsKeyword(item.Title+item.Slug+item.Summary+item.WorkTypeCode, q.Keyword) {
			continue
		}
		if q.WorkTypeCode != "" && item.WorkTypeCode != q.WorkTypeCode {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateSlice(filtered, q.Page, q.PageSize), nil
}
func (s *AdminPageService) Creators(q dto.PageQuery) (dto.PageResult[dto.AdminCreator], error) {
	if r, ok := s.repo.(repo.AdminCreatorPageRepo); ok {
		return r.PageAdminCreators(q)
	}
	list, err := s.repo.ListAdminCreators()
	if err != nil {
		return dto.PageResult[dto.AdminCreator]{}, err
	}
	filtered := make([]dto.AdminCreator, 0)
	for _, item := range list {
		if !containsKeyword(item.Name+item.Slug+item.Summary+item.CreatorTypeCode, q.Keyword) {
			continue
		}
		if q.CreatorTypeCode != "" && item.CreatorTypeCode != q.CreatorTypeCode {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateSlice(filtered, q.Page, q.PageSize), nil
}
func (s *AdminPageService) Dicts(dictKey string, page, pageSize int, keyword string) (dto.PageResult[dto.AdminDictItem], error) {
	if r, ok := s.repo.(repo.AdminDictPageRepo); ok {
		return r.PageAdminDictItems(dictKey, page, pageSize, keyword)
	}
	list, err := s.repo.ListAdminDictItems(dictKey)
	if err != nil { return dto.PageResult[dto.AdminDictItem]{}, err }
	filtered := make([]dto.AdminDictItem, 0)
	for _, item := range list {
		if containsKeyword(item.Name+item.Code, keyword) {
			filtered = append(filtered, item)
		}
	}
	return paginateSlice(filtered, page, pageSize), nil
}
