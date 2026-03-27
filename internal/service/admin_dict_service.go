package service

import (
	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type AdminDictService struct {
	repo repo.AdminRepo
}

func NewAdminDictService(r repo.AdminRepo) *AdminDictService {
	return &AdminDictService{repo: r}
}

func (s *AdminDictService) List(dictKey string) ([]dto.AdminDictItem, error) {
	return s.repo.ListAdminDictItems(dictKey)
}
func (s *AdminDictService) Get(dictKey, ref string) (dto.AdminDictItem, error) {
	return s.repo.GetAdminDictItem(dictKey, ref)
}
func (s *AdminDictService) Create(dictKey string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	in.DictKey = dictKey
	return s.repo.CreateAdminDictItem(dictKey, in)
}
func (s *AdminDictService) Update(dictKey, ref string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	in.DictKey = dictKey
	return s.repo.UpdateAdminDictItem(dictKey, ref, in)
}
func (s *AdminDictService) Delete(dictKey, ref string) error {
	return s.repo.DeleteAdminDictItem(dictKey, ref)
}
