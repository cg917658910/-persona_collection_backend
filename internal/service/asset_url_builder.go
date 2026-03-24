package service

import "strings"

type AssetURLBuilder struct {
	publicBaseURL string
	staticPrefix  string
}

func NewAssetURLBuilder(publicBaseURL, staticPrefix string) *AssetURLBuilder {
	return &AssetURLBuilder{
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
		staticPrefix:  "/" + strings.Trim(strings.TrimSpace(staticPrefix), "/"),
	}
}

func (b *AssetURLBuilder) Normalize(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}

	clean := raw
	if strings.HasPrefix(clean, "/assets/") {
		return b.publicBaseURL + b.staticPrefix + clean
	}
	if strings.HasPrefix(clean, "assets/") {
		return b.publicBaseURL + b.staticPrefix + "/" + clean
	}
	if strings.HasPrefix(clean, b.staticPrefix+"/") {
		return b.publicBaseURL + clean
	}
	if strings.HasPrefix(clean, "/") {
		return b.publicBaseURL + clean
	}
	return b.publicBaseURL + "/" + clean
}
