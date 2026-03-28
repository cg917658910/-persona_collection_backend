package service

import (
	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type AdminService struct {
	repo   repo.AdminRepo
	assets *AssetURLBuilder
}

func NewAdminService(r repo.AdminRepo, assets *AssetURLBuilder) *AdminService {
	return &AdminService{repo: r, assets: assets}
}

func (s *AdminService) ListCharacters() ([]dto.AdminCharacter, error) {
	list, err := s.repo.ListAdminCharacters()
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = s.normalizeCharacter(list[i])
	}
	return list, nil
}
func (s *AdminService) GetCharacter(ref string) (dto.AdminCharacter, error) {
	item, err := s.repo.GetAdminCharacter(ref)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	return s.normalizeCharacter(item), nil
}
func (s *AdminService) CreateCharacter(in dto.AdminCharacter) (dto.AdminCharacter, error) {
	in = s.toStorageCharacter(in)
	item, err := s.repo.CreateAdminCharacter(in)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	return s.normalizeCharacter(item), nil
}
func (s *AdminService) UpdateCharacter(ref string, in dto.AdminCharacter) (dto.AdminCharacter, error) {
	in = s.toStorageCharacter(in)
	item, err := s.repo.UpdateAdminCharacter(ref, in)
	if err != nil {
		return dto.AdminCharacter{}, err
	}
	return s.normalizeCharacter(item), nil
}
func (s *AdminService) DeleteCharacter(ref string) error { return s.repo.DeleteAdminCharacter(ref) }

func (s *AdminService) ListSongs() ([]dto.AdminSong, error) {
	list, err := s.repo.ListAdminSongs()
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = s.normalizeSong(list[i])
	}
	return list, nil
}
func (s *AdminService) GetSong(ref string) (dto.AdminSong, error) {
	item, err := s.repo.GetAdminSong(ref)
	if err != nil {
		return dto.AdminSong{}, err
	}
	return s.normalizeSong(item), nil
}
func (s *AdminService) CreateSong(in dto.AdminSong) (dto.AdminSong, error) {
	in = s.toStorageSong(in)
	item, err := s.repo.CreateAdminSong(in)
	if err != nil {
		return dto.AdminSong{}, err
	}
	return s.normalizeSong(item), nil
}
func (s *AdminService) UpdateSong(ref string, in dto.AdminSong) (dto.AdminSong, error) {
	in = s.toStorageSong(in)
	item, err := s.repo.UpdateAdminSong(ref, in)
	if err != nil {
		return dto.AdminSong{}, err
	}
	return s.normalizeSong(item), nil
}
func (s *AdminService) DeleteSong(ref string) error { return s.repo.DeleteAdminSong(ref) }

func (s *AdminService) ListRelations() ([]dto.AdminRelation, error) {
	list, err := s.repo.ListAdminRelations()
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = s.normalizeRelation(list[i])
	}
	return list, nil
}
func (s *AdminService) GetRelation(ref string) (dto.AdminRelation, error) {
	item, err := s.repo.GetAdminRelation(ref)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	return s.normalizeRelation(item), nil
}
func (s *AdminService) CreateRelation(in dto.AdminRelation) (dto.AdminRelation, error) {
	in = s.toStorageRelation(in)
	item, err := s.repo.CreateAdminRelation(in)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	return s.normalizeRelation(item), nil
}
func (s *AdminService) UpdateRelation(ref string, in dto.AdminRelation) (dto.AdminRelation, error) {
	in = s.toStorageRelation(in)
	item, err := s.repo.UpdateAdminRelation(ref, in)
	if err != nil {
		return dto.AdminRelation{}, err
	}
	return s.normalizeRelation(item), nil
}
func (s *AdminService) DeleteRelation(ref string) error { return s.repo.DeleteAdminRelation(ref) }

func (s *AdminService) ListThemes() ([]dto.AdminTheme, error) {
	list, err := s.repo.ListAdminThemes()
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = s.normalizeTheme(list[i])
	}
	return list, nil
}
func (s *AdminService) GetTheme(ref string) (dto.AdminTheme, error) {
	item, err := s.repo.GetAdminTheme(ref)
	if err != nil {
		return dto.AdminTheme{}, err
	}
	return s.normalizeTheme(item), nil
}
func (s *AdminService) CreateTheme(in dto.AdminTheme) (dto.AdminTheme, error) {
	in = s.toStorageTheme(in)
	item, err := s.repo.CreateAdminTheme(in)
	if err != nil {
		return dto.AdminTheme{}, err
	}
	return s.normalizeTheme(item), nil
}
func (s *AdminService) UpdateTheme(ref string, in dto.AdminTheme) (dto.AdminTheme, error) {
	in = s.toStorageTheme(in)
	item, err := s.repo.UpdateAdminTheme(ref, in)
	if err != nil {
		return dto.AdminTheme{}, err
	}
	return s.normalizeTheme(item), nil
}
func (s *AdminService) DeleteTheme(ref string) error { return s.repo.DeleteAdminTheme(ref) }

func (s *AdminService) normalizeCharacter(in dto.AdminCharacter) dto.AdminCharacter {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	return in
}
func (s *AdminService) toStorageCharacter(in dto.AdminCharacter) dto.AdminCharacter {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	return in
}
func (s *AdminService) normalizeSong(in dto.AdminSong) dto.AdminSong {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.AudioURL = s.assets.Normalize(in.AudioURL)
	return in
}
func (s *AdminService) toStorageSong(in dto.AdminSong) dto.AdminSong {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	in.AudioURL = s.assets.ToStorage(in.AudioURL)
	return in
}
func (s *AdminService) normalizeRelation(in dto.AdminRelation) dto.AdminRelation {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	for i := range in.Songs {
		in.Songs[i].CoverURL = s.assets.Normalize(in.Songs[i].CoverURL)
		in.Songs[i].AudioURL = s.assets.Normalize(in.Songs[i].AudioURL)
	}
	return in
}
func (s *AdminService) toStorageRelation(in dto.AdminRelation) dto.AdminRelation {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	for i := range in.Songs {
		in.Songs[i].CoverURL = s.assets.ToStorage(in.Songs[i].CoverURL)
		in.Songs[i].AudioURL = s.assets.ToStorage(in.Songs[i].AudioURL)
	}
	return in
}
func (s *AdminService) normalizeTheme(in dto.AdminTheme) dto.AdminTheme {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	return in
}
func (s *AdminService) toStorageTheme(in dto.AdminTheme) dto.AdminTheme {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	return in
}
