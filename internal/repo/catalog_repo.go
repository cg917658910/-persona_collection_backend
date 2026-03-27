package repo

import "pm-backend/internal/dto"

type CatalogRepo interface {
	Home() (dto.HomePayload, error)
	RandomCharacter(theme string, exclude []string) (dto.Character, error)

	ListCharacters() ([]dto.Character, error)
	GetCharacterDetail(slug string) (dto.CharacterDetail, error)

	ListWorks() ([]dto.Work, error)
	GetWorkDetail(slug string) (dto.Work, error)

	ListCreators() ([]dto.Creator, error)
	GetCreatorDetail(slug string) (dto.Creator, error)

	ListThemes() ([]dto.Theme, error)
	GetThemeDetail(slug string) (dto.ThemeDetail, error)

	ListSongs() ([]dto.Song, error)
	SearchCatalog(keyword string, limit int) (dto.SearchResponseData, error)
}
