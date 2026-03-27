package dto

type AdminDictItem struct {
	ID        string `json:"id,omitempty"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder,omitempty"`
	IsActive  bool   `json:"isActive"`
	DictKey   string `json:"dictKey,omitempty"`
}
