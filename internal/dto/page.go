package dto

type PageQuery struct {
	Page              int    `form:"page"`
	PageSize          int    `form:"pageSize"`
	Keyword           string `form:"keyword"`
	CharacterTypeCode string `form:"characterTypeCode"`
	Status            string `form:"status"`
	Category          string `form:"category"`
	WorkTypeCode      string `form:"workTypeCode"`
	CreatorTypeCode   string `form:"creatorTypeCode"`
}

type PageResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func NormalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
