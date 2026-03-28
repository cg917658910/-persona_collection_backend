package service

import (
	"pm-backend/internal/dto"
	"pm-backend/internal/repo"
)

type AdminImportService struct {
	repo repo.AdminRepo
}

func NewAdminImportService(r repo.AdminRepo) *AdminImportService {
	return &AdminImportService{repo: r}
}

func (s *AdminImportService) Validate(in dto.AdminImportRequest) (dto.AdminImportResult, error) {
	return s.repo.ValidateGeneratedPackage(in.Package)
}

func (s *AdminImportService) Run(in dto.AdminImportRequest) (dto.AdminImportResult, error) {
	return s.repo.ImportGeneratedPackage(in.Package)
}

func (s *AdminImportService) ValidateRelations(in dto.AdminRelationImportRequest) (dto.AdminImportResult, error) {
	return s.repo.ValidateRelationPackage(in.Package)
}

func (s *AdminImportService) RunRelations(in dto.AdminRelationImportRequest) (dto.AdminImportResult, error) {
	return s.repo.ImportRelationPackage(in.Package)
}
