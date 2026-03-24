package handler

import (
	"errors"
	"net/http"
	"strings"

	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	service *service.CatalogService
}

func NewCatalogHandler(s *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: s}
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

	normalizedExclude := make([]string, 0, len(exclude))
	for _, item := range exclude {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		normalizedExclude = append(normalizedExclude, item)
	}

	data, err := h.service.RandomCharacter(strings.TrimSpace(c.Query("theme")), normalizedExclude)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, data, nil)
}

func (h *CatalogHandler) ListCharacters(c *gin.Context) {
	list, err := h.service.ListCharacters()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": list}, gin.H{"page": 1, "pageSize": len(list), "total": len(list)})
}

func (h *CatalogHandler) GetCharacterDetail(c *gin.Context) {
	data, err := h.service.GetCharacterDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListWorks(c *gin.Context) {
	list, err := h.service.ListWorks()
	h.handleList(c, list, err)
}

func (h *CatalogHandler) GetWorkDetail(c *gin.Context) {
	data, err := h.service.GetWorkDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListCreators(c *gin.Context) {
	list, err := h.service.ListCreators()
	h.handleList(c, list, err)
}

func (h *CatalogHandler) GetCreatorDetail(c *gin.Context) {
	data, err := h.service.GetCreatorDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListThemes(c *gin.Context) {
	list, err := h.service.ListThemes()
	h.handleList(c, list, err)
}

func (h *CatalogHandler) GetThemeDetail(c *gin.Context) {
	data, err := h.service.GetThemeDetail(c.Param("slug"))
	h.handleOne(c, data, err)
}

func (h *CatalogHandler) ListSongs(c *gin.Context) {
	list, err := h.service.ListSongs()
	h.handleList(c, list, err)
}

func (h *CatalogHandler) handleList(c *gin.Context, list interface{}, err error) {
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": list}, nil)
}

func (h *CatalogHandler) handleOne(c *gin.Context, data interface{}, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errors.New("not found")) {
			status = http.StatusNotFound
		}
		if err.Error() == "character not found" || err.Error() == "work not found" || err.Error() == "creator not found" || err.Error() == "theme not found" {
			status = http.StatusNotFound
		}
		response.Error(c, status, err.Error())
		return
	}
	response.OK(c, data, nil)
}
