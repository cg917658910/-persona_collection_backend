package service

import (
	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type AdminWorkCreatorService struct {
	repo   repo.AdminRepo
	assets *AssetURLBuilder
}

func NewAdminWorkCreatorService(r repo.AdminRepo, assets *AssetURLBuilder) *AdminWorkCreatorService {
	return &AdminWorkCreatorService{repo: r, assets: assets}
}

func (s *AdminWorkCreatorService) ListWorks() ([]dto.AdminWork, error) {
	list, err := s.repo.ListAdminWorks()
	if err != nil { return nil, err }
	for i := range list { list[i].CoverURL = s.assets.Normalize(list[i].CoverURL) }
	return list, nil
}
func (s *AdminWorkCreatorService) GetWork(ref string) (dto.AdminWork, error) {
	item, err := s.repo.GetAdminWork(ref)
	if err != nil { return dto.AdminWork{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) CreateWork(in dto.AdminWork) (dto.AdminWork, error) {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	item, err := s.repo.CreateAdminWork(in)
	if err != nil { return dto.AdminWork{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) UpdateWork(ref string, in dto.AdminWork) (dto.AdminWork, error) {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	item, err := s.repo.UpdateAdminWork(ref, in)
	if err != nil { return dto.AdminWork{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) DeleteWork(ref string) error { return s.repo.DeleteAdminWork(ref) }

func (s *AdminWorkCreatorService) ListCreators() ([]dto.AdminCreator, error) {
	list, err := s.repo.ListAdminCreators()
	if err != nil { return nil, err }
	for i := range list { list[i].CoverURL = s.assets.Normalize(list[i].CoverURL) }
	return list, nil
}
func (s *AdminWorkCreatorService) GetCreator(ref string) (dto.AdminCreator, error) {
	item, err := s.repo.GetAdminCreator(ref)
	if err != nil { return dto.AdminCreator{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) CreateCreator(in dto.AdminCreator) (dto.AdminCreator, error) {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	item, err := s.repo.CreateAdminCreator(in)
	if err != nil { return dto.AdminCreator{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) UpdateCreator(ref string, in dto.AdminCreator) (dto.AdminCreator, error) {
	in.CoverURL = s.assets.ToStorage(in.CoverURL)
	item, err := s.repo.UpdateAdminCreator(ref, in)
	if err != nil { return dto.AdminCreator{}, err }
	item.CoverURL = s.assets.Normalize(item.CoverURL)
	return item, nil
}
func (s *AdminWorkCreatorService) DeleteCreator(ref string) error { return s.repo.DeleteAdminCreator(ref) }
