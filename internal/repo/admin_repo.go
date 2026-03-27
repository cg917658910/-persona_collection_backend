package repo

import "pm-backend/internal/dto"

type AdminRepo interface {
	ListAdminCharacters() ([]dto.AdminCharacter, error)
	GetAdminCharacter(ref string) (dto.AdminCharacter, error)
	CreateAdminCharacter(in dto.AdminCharacter) (dto.AdminCharacter, error)
	UpdateAdminCharacter(ref string, in dto.AdminCharacter) (dto.AdminCharacter, error)
	DeleteAdminCharacter(ref string) error

	ListAdminSongs() ([]dto.AdminSong, error)
	GetAdminSong(ref string) (dto.AdminSong, error)
	CreateAdminSong(in dto.AdminSong) (dto.AdminSong, error)
	UpdateAdminSong(ref string, in dto.AdminSong) (dto.AdminSong, error)
	DeleteAdminSong(ref string) error

	ListAdminThemes() ([]dto.AdminTheme, error)
	GetAdminTheme(ref string) (dto.AdminTheme, error)
	CreateAdminTheme(in dto.AdminTheme) (dto.AdminTheme, error)
	UpdateAdminTheme(ref string, in dto.AdminTheme) (dto.AdminTheme, error)
	DeleteAdminTheme(ref string) error

	ListAdminWorks() ([]dto.AdminWork, error)
	GetAdminWork(ref string) (dto.AdminWork, error)
	CreateAdminWork(in dto.AdminWork) (dto.AdminWork, error)
	UpdateAdminWork(ref string, in dto.AdminWork) (dto.AdminWork, error)
	DeleteAdminWork(ref string) error

	ListAdminCreators() ([]dto.AdminCreator, error)
	GetAdminCreator(ref string) (dto.AdminCreator, error)
	CreateAdminCreator(in dto.AdminCreator) (dto.AdminCreator, error)
	UpdateAdminCreator(ref string, in dto.AdminCreator) (dto.AdminCreator, error)
	DeleteAdminCreator(ref string) error

	ListAdminDictItems(dictKey string) ([]dto.AdminDictItem, error)
	GetAdminDictItem(dictKey, ref string) (dto.AdminDictItem, error)
	CreateAdminDictItem(dictKey string, in dto.AdminDictItem) (dto.AdminDictItem, error)
	UpdateAdminDictItem(dictKey, ref string, in dto.AdminDictItem) (dto.AdminDictItem, error)
	DeleteAdminDictItem(dictKey, ref string) error

	ValidateGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error)
	ImportGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error)
}
