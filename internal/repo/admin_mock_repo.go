package repo

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"pm-backend/internal/dto"
)

type mockAdminRepo struct {
	characters []dto.AdminCharacter
	songs      []dto.AdminSong
	themes     []dto.AdminTheme
	works      []dto.AdminWork
	creators   []dto.AdminCreator
	dicts      map[string][]dto.AdminDictItem
}

func NewMockAdminRepo() AdminRepo {
	return &mockAdminRepo{
		characters: []dto.AdminCharacter{
			{
				ID: "1", Slug: "lin-daiyu", Name: "林黛玉", Summary: "高敏感、高自尊的真情承受者。", CoverURL: "/assets/images/characters/lin-daiyu.webp",
				Status: "published", Type: "文学人物", CharacterTypeCode: "literary", OneLineDefinition: "她不是脆弱，而是过度敏感、过度清醒、又过度珍视真情的人。",
				CoreIdentity: "真情承受者", CoreFear: "被辜负", CoreConflict: "自尊与依恋并存", EmotionalTone: "清冷哀愁",
				WorkSlugs: []string{"dream-of-the-red-chamber"}, WorkNames: []string{"红楼梦"}, ThemeSlugs: []string{"tragic"}, ThemeNames: []string{"悲剧人格"},
				SongSlugs: []string{"lin-daiyu-theme-v1"}, HasSong: true,
			},
			{
				ID: "2", Slug: "sun-wu-kong", Name: "孙悟空", Summary: "不愿被驯服的自由反抗者。", CoverURL: "/assets/images/characters/sun-wu-kong.webp",
				Status: "published", Type: "文学人物", CharacterTypeCode: "literary", OneLineDefinition: "他不是不守规矩，而是不愿被不合理的规矩困住。",
				CoreIdentity: "自由反抗者", CoreFear: "被驯服", CoreConflict: "野性与责任冲突", EmotionalTone: "热烈桀骜",
				WorkSlugs: []string{"journey-to-the-west"}, WorkNames: []string{"西游记"}, ThemeSlugs: []string{"rebels"}, ThemeNames: []string{"反叛者"},
				SongSlugs: []string{"sun-wu-kong-theme-v1"}, HasSong: true,
			},
		},
		songs: []dto.AdminSong{
			{ID: "1", Slug: "lin-daiyu-theme-v1", Title: "林黛玉之歌", CharacterSlug: "lin-daiyu", CharacterName: "林黛玉", Summary: "清冷哀愁的真情独白。", CoverURL: "/assets/images/songs/lin-daiyu-theme-v1.webp", AudioURL: "/assets/audio/lin-daiyu-theme-v1.mp3", Status: "published", CoreTheme: "真情无所安放", Styles: []string{"国风", "抒情"}, EmotionalCurve: []string{"压抑", "柔痛", "余冷"}},
			{ID: "2", Slug: "sun-wu-kong-theme-v1", Title: "孙悟空之歌", CharacterSlug: "sun-wu-kong", CharacterName: "孙悟空", Summary: "不服与守护交织的火焰。", CoverURL: "/assets/images/songs/sun-wu-kong-theme-v1.webp", AudioURL: "/assets/audio/sun-wu-kong-theme-v1.mp3", Status: "published", CoreTheme: "自由与责任", Styles: []string{"摇滚", "热血"}, EmotionalCurve: []string{"躁动", "爆发", "再燃"}},
		},
		themes: []dto.AdminTheme{
			{ID: "1", Slug: "tragic", Name: "悲剧人格", Code: "tragic", Category: "destiny", Summary: "内在真实与命运难以相容。", CoverURL: "/assets/images/themes/tragic.webp", Status: "published", CharacterSlugs: []string{"lin-daiyu"}},
			{ID: "2", Slug: "rebels", Name: "反叛者", Code: "rebels", Category: "psychology", Summary: "拒绝被不合理秩序驯服。", CoverURL: "/assets/images/themes/rebels.webp", Status: "published", CharacterSlugs: []string{"sun-wu-kong"}},
		},
		works: []dto.AdminWork{
			{ID: "1", Slug: "dream-of-the-red-chamber", Title: "红楼梦", Summary: "家族盛衰与真情消逝。", CoverURL: "/assets/images/works/dream-of-the-red-chamber.webp", Status: "published", WorkTypeCode: "novel", CreatorSlugs: []string{"cao-xueqin"}, CreatorNames: []string{"曹雪芹"}, RegionCode: "china", CulturalRegionCode: "literature", ReleaseYear: 1791},
			{ID: "2", Slug: "journey-to-the-west", Title: "西游记", Summary: "神魔冒险与修行之路。", CoverURL: "/assets/images/works/journey-to-the-west.webp", Status: "published", WorkTypeCode: "novel", CreatorSlugs: []string{"wu-chengen"}, CreatorNames: []string{"吴承恩"}, RegionCode: "china", CulturalRegionCode: "literature", ReleaseYear: 1592},
		},
		creators: []dto.AdminCreator{
			{ID: "1", Slug: "cao-xueqin", Name: "曹雪芹", Summary: "中国古典文学作家。", CoverURL: "/assets/images/creators/cao-xueqin.webp", Status: "published", CreatorTypeCode: "author", WorkSlugs: []string{"dream-of-the-red-chamber"}, WorkNames: []string{"红楼梦"}, RegionCode: "china", CulturalRegionCode: "literature"},
			{ID: "2", Slug: "wu-chengen", Name: "吴承恩", Summary: "中国古典文学作家。", CoverURL: "/assets/images/creators/wu-chengen.webp", Status: "published", CreatorTypeCode: "author", WorkSlugs: []string{"journey-to-the-west"}, WorkNames: []string{"西游记"}, RegionCode: "china", CulturalRegionCode: "literature"},
		},
		dicts: map[string][]dto.AdminDictItem{
			"characterTypes": {
				{ID: "1", Code: "historical", Name: "历史人物", SortOrder: 10, IsActive: true, DictKey: "characterTypes"},
				{ID: "2", Code: "literary", Name: "文学人物", SortOrder: 20, IsActive: true, DictKey: "characterTypes"},
			},
			"workTypes": {
				{ID: "1", Code: "novel", Name: "小说", SortOrder: 10, IsActive: true, DictKey: "workTypes"},
				{ID: "2", Code: "history", Name: "历史语境", SortOrder: 20, IsActive: true, DictKey: "workTypes"},
			},
			"creatorTypes": {
				{ID: "1", Code: "author", Name: "作者", SortOrder: 10, IsActive: true, DictKey: "creatorTypes"},
				{ID: "2", Code: "director", Name: "导演", SortOrder: 20, IsActive: true, DictKey: "creatorTypes"},
			},
			"regions": {
				{ID: "1", Code: "china", Name: "中国", SortOrder: 10, IsActive: true, DictKey: "regions"},
				{ID: "2", Code: "japan", Name: "日本", SortOrder: 20, IsActive: true, DictKey: "regions"},
			},
			"culturalRegions": {
				{ID: "1", Code: "literature", Name: "文学", SortOrder: 10, IsActive: true, DictKey: "culturalRegions"},
				{ID: "2", Code: "history", Name: "历史", SortOrder: 20, IsActive: true, DictKey: "culturalRegions"},
			},
			"motivations": {
				{ID: "1", Code: "freedom", Name: "自由", SortOrder: 10, IsActive: true, DictKey: "motivations"},
				{ID: "2", Code: "control", Name: "控制", SortOrder: 20, IsActive: true, DictKey: "motivations"},
			},
			"themeCategories": {
				{ID: "1", Code: "psychology", Name: "心理结构", SortOrder: 10, IsActive: true, DictKey: "themeCategories"},
				{ID: "2", Code: "destiny", Name: "命运气候", SortOrder: 20, IsActive: true, DictKey: "themeCategories"},
			},
		},
	}
}

func (r *mockAdminRepo) ListAdminCharacters() ([]dto.AdminCharacter, error) { return append([]dto.AdminCharacter{}, r.characters...), nil }
func (r *mockAdminRepo) GetAdminCharacter(ref string) (dto.AdminCharacter, error) {
	for _, item := range r.characters {
		if item.ID == ref || item.Slug == ref {
			return item, nil
		}
	}
	return dto.AdminCharacter{}, errors.New("admin character not found")
}
func (r *mockAdminRepo) CreateAdminCharacter(in dto.AdminCharacter) (dto.AdminCharacter, error) {
	in.ID = fmt.Sprintf("%d", len(r.characters)+1)
	if in.Status == "" { in.Status = "draft" }
	r.characters = append([]dto.AdminCharacter{in}, r.characters...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminCharacter(ref string, in dto.AdminCharacter) (dto.AdminCharacter, error) {
	for idx, item := range r.characters {
		if item.ID == ref || item.Slug == ref {
			in.ID = item.ID
			if in.Status == "" { in.Status = item.Status }
			r.characters[idx] = in
			return in, nil
		}
	}
	return dto.AdminCharacter{}, errors.New("admin character not found")
}
func (r *mockAdminRepo) DeleteAdminCharacter(ref string) error {
	for idx, item := range r.characters {
		if item.ID == ref || item.Slug == ref {
			r.characters = slices.Delete(r.characters, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin character not found")
}

func (r *mockAdminRepo) ListAdminSongs() ([]dto.AdminSong, error) { return append([]dto.AdminSong{}, r.songs...), nil }
func (r *mockAdminRepo) GetAdminSong(ref string) (dto.AdminSong, error) {
	for _, item := range r.songs {
		if item.ID == ref || item.Slug == ref { return item, nil }
	}
	return dto.AdminSong{}, errors.New("admin song not found")
}
func (r *mockAdminRepo) CreateAdminSong(in dto.AdminSong) (dto.AdminSong, error) {
	in.ID = fmt.Sprintf("%d", len(r.songs)+1)
	if in.Status == "" { in.Status = "draft" }
	r.songs = append([]dto.AdminSong{in}, r.songs...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminSong(ref string, in dto.AdminSong) (dto.AdminSong, error) {
	for idx, item := range r.songs {
		if item.ID == ref || item.Slug == ref {
			in.ID = item.ID
			if in.Status == "" { in.Status = item.Status }
			r.songs[idx] = in
			return in, nil
		}
	}
	return dto.AdminSong{}, errors.New("admin song not found")
}
func (r *mockAdminRepo) DeleteAdminSong(ref string) error {
	for idx, item := range r.songs {
		if item.ID == ref || item.Slug == ref {
			r.songs = slices.Delete(r.songs, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin song not found")
}

func (r *mockAdminRepo) ListAdminThemes() ([]dto.AdminTheme, error) { return append([]dto.AdminTheme{}, r.themes...), nil }
func (r *mockAdminRepo) GetAdminTheme(ref string) (dto.AdminTheme, error) {
	for _, item := range r.themes {
		if item.ID == ref || item.Slug == ref { return item, nil }
	}
	return dto.AdminTheme{}, errors.New("admin theme not found")
}
func (r *mockAdminRepo) CreateAdminTheme(in dto.AdminTheme) (dto.AdminTheme, error) {
	in.ID = fmt.Sprintf("%d", len(r.themes)+1)
	if in.Status == "" { in.Status = "draft" }
	r.themes = append([]dto.AdminTheme{in}, r.themes...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminTheme(ref string, in dto.AdminTheme) (dto.AdminTheme, error) {
	for idx, item := range r.themes {
		if item.ID == ref || item.Slug == ref {
			in.ID = item.ID
			if in.Status == "" { in.Status = item.Status }
			r.themes[idx] = in
			return in, nil
		}
	}
	return dto.AdminTheme{}, errors.New("admin theme not found")
}
func (r *mockAdminRepo) DeleteAdminTheme(ref string) error {
	for idx, item := range r.themes {
		if item.ID == ref || item.Slug == ref {
			r.themes = slices.Delete(r.themes, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin theme not found")
}


func (r *mockAdminRepo) ListAdminWorks() ([]dto.AdminWork, error) { return append([]dto.AdminWork{}, r.works...), nil }
func (r *mockAdminRepo) GetAdminWork(ref string) (dto.AdminWork, error) {
	for _, item := range r.works {
		if item.ID == ref || item.Slug == ref { return item, nil }
	}
	return dto.AdminWork{}, errors.New("admin work not found")
}
func (r *mockAdminRepo) CreateAdminWork(in dto.AdminWork) (dto.AdminWork, error) {
	in.ID = fmt.Sprintf("%d", len(r.works)+1)
	if in.Status == "" { in.Status = "draft" }
	r.works = append([]dto.AdminWork{in}, r.works...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminWork(ref string, in dto.AdminWork) (dto.AdminWork, error) {
	for idx, item := range r.works {
		if item.ID == ref || item.Slug == ref {
			in.ID = item.ID
			if in.Status == "" { in.Status = item.Status }
			r.works[idx] = in
			return in, nil
		}
	}
	return dto.AdminWork{}, errors.New("admin work not found")
}
func (r *mockAdminRepo) DeleteAdminWork(ref string) error {
	for idx, item := range r.works {
		if item.ID == ref || item.Slug == ref {
			r.works = slices.Delete(r.works, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin work not found")
}

func (r *mockAdminRepo) ListAdminCreators() ([]dto.AdminCreator, error) { return append([]dto.AdminCreator{}, r.creators...), nil }
func (r *mockAdminRepo) GetAdminCreator(ref string) (dto.AdminCreator, error) {
	for _, item := range r.creators {
		if item.ID == ref || item.Slug == ref { return item, nil }
	}
	return dto.AdminCreator{}, errors.New("admin creator not found")
}
func (r *mockAdminRepo) CreateAdminCreator(in dto.AdminCreator) (dto.AdminCreator, error) {
	in.ID = fmt.Sprintf("%d", len(r.creators)+1)
	if in.Status == "" { in.Status = "draft" }
	r.creators = append([]dto.AdminCreator{in}, r.creators...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminCreator(ref string, in dto.AdminCreator) (dto.AdminCreator, error) {
	for idx, item := range r.creators {
		if item.ID == ref || item.Slug == ref {
			in.ID = item.ID
			if in.Status == "" { in.Status = item.Status }
			r.creators[idx] = in
			return in, nil
		}
	}
	return dto.AdminCreator{}, errors.New("admin creator not found")
}
func (r *mockAdminRepo) DeleteAdminCreator(ref string) error {
	for idx, item := range r.creators {
		if item.ID == ref || item.Slug == ref {
			r.creators = slices.Delete(r.creators, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin creator not found")
}


func (r *mockAdminRepo) ListAdminDictItems(dictKey string) ([]dto.AdminDictItem, error) {
	items, ok := r.dicts[dictKey]
	if !ok {
		return []dto.AdminDictItem{}, nil
	}
	return append([]dto.AdminDictItem{}, items...), nil
}
func (r *mockAdminRepo) GetAdminDictItem(dictKey, ref string) (dto.AdminDictItem, error) {
	for _, item := range r.dicts[dictKey] {
		if item.ID == ref || item.Code == ref {
			return item, nil
		}
	}
	return dto.AdminDictItem{}, errors.New("admin dict item not found")
}
func (r *mockAdminRepo) CreateAdminDictItem(dictKey string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	in.ID = fmt.Sprintf("%d", len(r.dicts[dictKey])+1)
	in.DictKey = dictKey
	in.IsActive = true
	r.dicts[dictKey] = append([]dto.AdminDictItem{in}, r.dicts[dictKey]...)
	return in, nil
}
func (r *mockAdminRepo) UpdateAdminDictItem(dictKey, ref string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	items := r.dicts[dictKey]
	for idx, item := range items {
		if item.ID == ref || item.Code == ref {
			in.ID = item.ID
			in.DictKey = dictKey
			items[idx] = in
			r.dicts[dictKey] = items
			return in, nil
		}
	}
	return dto.AdminDictItem{}, errors.New("admin dict item not found")
}
func (r *mockAdminRepo) DeleteAdminDictItem(dictKey, ref string) error {
	items := r.dicts[dictKey]
	for idx, item := range items {
		if item.ID == ref || item.Code == ref {
			r.dicts[dictKey] = slices.Delete(items, idx, idx+1)
			return nil
		}
	}
	return errors.New("admin dict item not found")
}


func (r *mockAdminRepo) PageAdminCharacters(q dto.PageQuery) (dto.PageResult[dto.AdminCharacter], error) {
	list, _ := r.ListAdminCharacters()
	return pageAdminCharacters(list, q), nil
}
func (r *mockAdminRepo) PageAdminSongs(q dto.PageQuery) (dto.PageResult[dto.AdminSong], error) {
	list, _ := r.ListAdminSongs()
	return pageAdminSongs(list, q), nil
}
func (r *mockAdminRepo) PageAdminThemes(q dto.PageQuery) (dto.PageResult[dto.AdminTheme], error) {
	list, _ := r.ListAdminThemes()
	return pageAdminThemes(list, q), nil
}
func (r *mockAdminRepo) PageAdminWorks(q dto.PageQuery) (dto.PageResult[dto.AdminWork], error) {
	list, _ := r.ListAdminWorks()
	return pageAdminWorks(list, q), nil
}
func (r *mockAdminRepo) PageAdminCreators(q dto.PageQuery) (dto.PageResult[dto.AdminCreator], error) {
	list, _ := r.ListAdminCreators()
	return pageAdminCreators(list, q), nil
}
func (r *mockAdminRepo) PageAdminDictItems(dictKey string, page, pageSize int, keyword string) (dto.PageResult[dto.AdminDictItem], error) {
	list, _ := r.ListAdminDictItems(dictKey)
	return pageAdminDicts(list, page, pageSize, keyword), nil
}

func paginateAdminSlice[T any](items []T, page, pageSize int) dto.PageResult[T] {
	page, pageSize = dto.NormalizePage(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start > total { start = total }
	end := start + pageSize
	if end > total { end = total }
	return dto.PageResult[T]{Items: items[start:end], Total: total, Page: page, PageSize: pageSize}
}
func pageAdminCharacters(items []dto.AdminCharacter, q dto.PageQuery) dto.PageResult[dto.AdminCharacter] {
	filtered := make([]dto.AdminCharacter, 0)
	for _, item := range items {
		if q.Keyword != "" && !strings.Contains(strings.ToLower(item.Name+item.Slug+item.Summary+item.OneLineDefinition), strings.ToLower(q.Keyword)) {
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
	return paginateAdminSlice(filtered, q.Page, q.PageSize)
}
func pageAdminSongs(items []dto.AdminSong, q dto.PageQuery) dto.PageResult[dto.AdminSong] {
	filtered := make([]dto.AdminSong, 0)
	for _, item := range items {
		if q.Keyword != "" && !strings.Contains(strings.ToLower(item.Title+item.Slug+item.Summary+item.CharacterName+item.CoreTheme), strings.ToLower(q.Keyword)) {
			continue
		}
		if q.Status != "" && item.Status != q.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	return paginateAdminSlice(filtered, q.Page, q.PageSize)
}
func pageAdminThemes(items []dto.AdminTheme, q dto.PageQuery) dto.PageResult[dto.AdminTheme] {
	filtered := make([]dto.AdminTheme, 0)
	for _, item := range items {
		if q.Keyword != "" && !strings.Contains(strings.ToLower(item.Name+item.Slug+item.Summary+item.Code+item.Category), strings.ToLower(q.Keyword)) {
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
	return paginateAdminSlice(filtered, q.Page, q.PageSize)
}
func pageAdminWorks(items []dto.AdminWork, q dto.PageQuery) dto.PageResult[dto.AdminWork] {
	filtered := make([]dto.AdminWork, 0)
	for _, item := range items {
		if q.Keyword != "" && !strings.Contains(strings.ToLower(item.Title+item.Slug+item.Summary+item.WorkTypeCode), strings.ToLower(q.Keyword)) {
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
	return paginateAdminSlice(filtered, q.Page, q.PageSize)
}
func pageAdminCreators(items []dto.AdminCreator, q dto.PageQuery) dto.PageResult[dto.AdminCreator] {
	filtered := make([]dto.AdminCreator, 0)
	for _, item := range items {
		if q.Keyword != "" && !strings.Contains(strings.ToLower(item.Name+item.Slug+item.Summary+item.CreatorTypeCode), strings.ToLower(q.Keyword)) {
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
	return paginateAdminSlice(filtered, q.Page, q.PageSize)
}
func pageAdminDicts(items []dto.AdminDictItem, page, pageSize int, keyword string) dto.PageResult[dto.AdminDictItem] {
	filtered := make([]dto.AdminDictItem, 0)
	for _, item := range items {
		if keyword == "" || strings.Contains(strings.ToLower(item.Name+item.Code), strings.ToLower(keyword)) {
			filtered = append(filtered, item)
		}
	}
	return paginateAdminSlice(filtered, page, pageSize)
}
