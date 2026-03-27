package dto

type AdminResourceType string

const (
	AdminResourceTypeImage AdminResourceType = "image"
	AdminResourceTypeAudio AdminResourceType = "audio"
)

type AdminResourceItem struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	URL          string            `json:"url"`
	Type         AdminResourceType `json:"type"`
	MimeType     string            `json:"mimeType,omitempty"`
	Size         int64             `json:"size,omitempty"`
	LinkedModule string            `json:"linkedModule,omitempty"`
	LinkedCount  int               `json:"linkedCount,omitempty"`
	CreatedAt    string            `json:"createdAt,omitempty"`
}

type AdminCreateResourceRequest struct {
	Name         string            `json:"name"`
	URL          string            `json:"url"`
	Type         AdminResourceType `json:"type"`
	MimeType     string            `json:"mimeType,omitempty"`
	Size         int64             `json:"size,omitempty"`
	LinkedModule string            `json:"linkedModule,omitempty"`
	LinkedCount  int               `json:"linkedCount,omitempty"`
}
