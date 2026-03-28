package handler

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"path"
	"strings"

	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	service         *service.CatalogService
	shareBaseURL    string
	frontendBaseURL string
}

func NewCatalogHandler(s *service.CatalogService, shareBaseURL string, frontendBaseURL string) *CatalogHandler {
	return &CatalogHandler{
		service:         s,
		shareBaseURL:    strings.TrimRight(shareBaseURL, "/"),
		frontendBaseURL: strings.TrimRight(frontendBaseURL, "/"),
	}
}

func (h *CatalogHandler) Home(c *gin.Context) {
	data, err := h.service.Home()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data, nil)
}

func (h *CatalogHandler) RandomCharacter(c *gin.Context) {
	exclude := c.QueryArray("exclude")
	if len(exclude) == 0 {
		if raw := strings.TrimSpace(c.Query("exclude")); raw != "" {
			exclude = strings.Split(raw, ",")
		}
	}
	data, err := h.service.RandomCharacter(strings.TrimSpace(c.Query("theme")), exclude)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, data, nil)
}

func (h *CatalogHandler) ListCharacters(c *gin.Context) {
	data, err := h.service.ListCharacters()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data, nil)
}

func (h *CatalogHandler) GetCharacterDetail(c *gin.Context) {
	data, err := h.service.GetCharacterDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListRelationships(c *gin.Context) {
	list, err := h.service.ListRelationships(strings.TrimSpace(c.Query("character")))
	h.handleList(c, gin.H{"items": list}, err)
}

func (h *CatalogHandler) GetRelationshipDetail(c *gin.Context) {
	data, err := h.service.GetRelationshipDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) CharacterSharePage(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	character, err := h.service.GetCharacterDetail(slug)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(`<!doctype html><html><head><meta charset="utf-8"><title>Character Not Found</title></head><body>Character not found.</body></html>`))
		return
	}

	title := html.EscapeString(character.Name)
	description := html.EscapeString(strings.TrimSpace(character.OneLineDefinition))
	if description == "" {
		description = html.EscapeString(strings.TrimSpace(character.Summary))
	}
	imageURL := html.EscapeString(h.shareImageURL(strings.TrimSpace(character.CoverURL)))
	shareURL := html.EscapeString(strings.TrimRight(h.shareBaseURL, "/") + "/share/character/" + slug)
	redirectURL := html.EscapeString(strings.TrimRight(h.frontendBaseURL, "/") + "/#/character/" + slug)
	clientRedirectURL := strings.TrimRight(h.frontendBaseURL, "/") + "/#/character/" + slug

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(h.shareHTMLPage(title, description, imageURL, shareURL, redirectURL, clientRedirectURL, "正在打开人物详情...", "如果没有自动跳转，请点这里")))
}

func (h *CatalogHandler) RelationSharePage(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	relation, err := h.service.GetRelationshipDetail(slug)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(`<!doctype html><html><head><meta charset="utf-8"><title>Relation Not Found</title></head><body>Relation not found.</body></html>`))
		return
	}

	title := html.EscapeString(strings.TrimSpace(relation.Name))
	if title == "" {
		title = html.EscapeString(strings.TrimSpace(relation.SourceCharacterName + " × " + relation.TargetCharacterName))
	}
	description := html.EscapeString(strings.TrimSpace(relation.OneLineDefinition))
	if description == "" {
		description = html.EscapeString(strings.TrimSpace(relation.Summary))
	}
	imageURL := html.EscapeString(h.shareImageURL(strings.TrimSpace(relation.CoverURL)))
	shareURL := html.EscapeString(strings.TrimRight(h.shareBaseURL, "/") + "/share/relation/" + slug)
	redirectURL := html.EscapeString(strings.TrimRight(h.frontendBaseURL, "/") + "/#/relationship/" + slug)
	clientRedirectURL := strings.TrimRight(h.frontendBaseURL, "/") + "/#/relationship/" + slug

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(h.shareHTMLPage(title, description, imageURL, shareURL, redirectURL, clientRedirectURL, "正在打开关系详情...", "如果没有自动跳转，请点这里")))
}

func (h *CatalogHandler) shareHTMLPage(title, description, imageURL, shareURL, redirectURL, clientRedirectURL, loadingText, linkText string) string {
	loadingText = html.EscapeString(loadingText)
	linkText = html.EscapeString(linkText)

	return fmt.Sprintf(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
  <title>%s</title>
  <meta name="description" content="%s">
  <meta property="og:type" content="website">
  <meta property="og:title" content="%s">
  <meta property="og:description" content="%s">
  <meta property="og:image" content="%s">
  <meta property="og:image:secure_url" content="%s">
  <meta property="og:image:type" content="image/png">
  <meta property="og:image:width" content="1200">
  <meta property="og:image:height" content="630">
  <meta property="og:url" content="%s">
  <meta name="image" content="%s">
  <link rel="image_src" href="%s">
  <meta name="twitter:card" content="summary_large_image">
  <meta name="twitter:title" content="%s">
  <meta name="twitter:description" content="%s">
  <meta name="twitter:image" content="%s">
  <meta name="twitter:image:src" content="%s">
  <meta http-equiv="refresh" content="0;url=%s">
  <script>window.location.replace(%q)</script>
</head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#0f1115;color:#fff;padding:24px;">
  <p>%s</p>
  <p><a href="%s" style="color:#d6b36a;">%s</a></p>
</body>
</html>`,
		title,
		description,
		title,
		description,
		imageURL,
		imageURL,
		shareURL,
		imageURL,
		imageURL,
		title,
		description,
		imageURL,
		imageURL,
		redirectURL,
		clientRedirectURL,
		loadingText,
		redirectURL,
		linkText,
	)
}

func (h *CatalogHandler) shareImageURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return h.defaultShareImageURL()
	}
	lower := strings.ToLower(raw)
	if strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		return raw
	}
	if strings.HasSuffix(lower, ".webp") {
		return h.defaultShareImageURL()
	}
	if parsed, err := url.Parse(raw); err == nil {
		ext := strings.ToLower(path.Ext(parsed.Path))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			return raw
		}
	}
	return h.defaultShareImageURL()
}

func (h *CatalogHandler) defaultShareImageURL() string {
	return strings.TrimRight(h.shareBaseURL, "/") + "/static/assets/images/share/default-share.png"
}

func (h *CatalogHandler) ListWorks(c *gin.Context) {
	list, err := h.service.ListWorks()
	h.handleList(c, gin.H{"items": list}, err)
}

func (h *CatalogHandler) GetWorkDetail(c *gin.Context) {
	data, err := h.service.GetWorkDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListCreators(c *gin.Context) {
	list, err := h.service.ListCreators()
	h.handleList(c, gin.H{"items": list}, err)
}

func (h *CatalogHandler) GetCreatorDetail(c *gin.Context) {
	data, err := h.service.GetCreatorDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListThemes(c *gin.Context) {
	list, err := h.service.ListThemes()
	h.handleList(c, gin.H{"items": list}, err)
}

func (h *CatalogHandler) GetThemeDetail(c *gin.Context) {
	data, err := h.service.GetThemeDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListSongs(c *gin.Context) {
	list, err := h.service.ListSongs()
	h.handleList(c, gin.H{"items": list}, err)
}

func (h *CatalogHandler) Search(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		keyword = strings.TrimSpace(c.Query("keyword"))
	}

	limit := 8
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		var parsed int
		if _, err := fmt.Sscanf(raw, "%d", &parsed); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	data, err := h.service.SearchCatalog(keyword, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data, nil)
}

func (h *CatalogHandler) handleList(c *gin.Context, list interface{}, err error) {
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, list, nil)
}

func (h *CatalogHandler) handleOne(c *gin.Context, data interface{}, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		response.Error(c, status, err.Error())
		return
	}
	response.OK(c, data, nil)
}
