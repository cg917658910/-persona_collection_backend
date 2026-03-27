package repo

import "pm-backend/internal/dto"

type AdminCharacterPageRepo interface {
	PageAdminCharacters(q dto.PageQuery) (dto.PageResult[dto.AdminCharacter], error)
}
type AdminSongPageRepo interface {
	PageAdminSongs(q dto.PageQuery) (dto.PageResult[dto.AdminSong], error)
}
type AdminThemePageRepo interface {
	PageAdminThemes(q dto.PageQuery) (dto.PageResult[dto.AdminTheme], error)
}
type AdminWorkPageRepo interface {
	PageAdminWorks(q dto.PageQuery) (dto.PageResult[dto.AdminWork], error)
}
type AdminCreatorPageRepo interface {
	PageAdminCreators(q dto.PageQuery) (dto.PageResult[dto.AdminCreator], error)
}
type AdminDictPageRepo interface {
	PageAdminDictItems(dictKey string, page, pageSize int, keyword string) (dto.PageResult[dto.AdminDictItem], error)
}
