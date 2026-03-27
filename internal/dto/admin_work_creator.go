package dto

type AdminWork struct {
	ID                 string   `json:"id,omitempty"`
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Summary            string   `json:"summary"`
	CoverURL           string   `json:"coverUrl"`
	Status             string   `json:"status"`
	WorkTypeCode       string   `json:"workTypeCode,omitempty"`
	CreatorSlugs       []string `json:"creatorSlugs,omitempty"`
	CreatorNames       []string `json:"creatorNames,omitempty"`
	RegionCode         string   `json:"regionCode,omitempty"`
	CulturalRegionCode string   `json:"culturalRegionCode,omitempty"`
	ReleaseYear        int      `json:"releaseYear,omitempty"`
	SortOrder          int      `json:"sortOrder"`
	Recommended        bool     `json:"recommended"`
	RecommendSort      int      `json:"recommendSort"`
}

type AdminCreator struct {
	ID                 string   `json:"id,omitempty"`
	Slug               string   `json:"slug"`
	Name               string   `json:"name"`
	Summary            string   `json:"summary"`
	CoverURL           string   `json:"coverUrl"`
	Status             string   `json:"status"`
	CreatorTypeCode    string   `json:"creatorTypeCode,omitempty"`
	WorkSlugs          []string `json:"workSlugs,omitempty"`
	WorkNames          []string `json:"workNames,omitempty"`
	RegionCode         string   `json:"regionCode,omitempty"`
	CulturalRegionCode string   `json:"culturalRegionCode,omitempty"`
	SortOrder          int      `json:"sortOrder"`
}
