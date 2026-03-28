package dto

type PageQuery struct {
	Page              int    `form:"page"`
	PageSize          int    `form:"pageSize"`
	Keyword           string `form:"keyword"`
	SubjectType       string `form:"subjectType"`
	CharacterTypeCode string `form:"characterTypeCode"`
	Status            string `form:"status"`
	Category          string `form:"category"`
	WorkTypeCode      string `form:"workTypeCode"`
	CreatorTypeCode   string `form:"creatorTypeCode"`
	RelationTypeCode  string `form:"relationTypeCode"`
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
