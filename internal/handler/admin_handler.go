package handler

import (
	"net/http"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h *AdminHandler) ListCharacters(c *gin.Context) {
	list, err := h.service.ListCharacters()
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"list": list}, nil)
}
func (h *AdminHandler) GetCharacter(c *gin.Context) {
	item, err := h.service.GetCharacter(c.Param("ref"))
	if err != nil { response.Error(c, http.StatusNotFound, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) CreateCharacter(c *gin.Context) {
	var in dto.AdminCharacter
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.CreateCharacter(in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) UpdateCharacter(c *gin.Context) {
	var in dto.AdminCharacter
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.UpdateCharacter(c.Param("ref"), in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) DeleteCharacter(c *gin.Context) {
	if err := h.service.DeleteCharacter(c.Param("ref")); err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"success": true}, nil)
}

func (h *AdminHandler) ListSongs(c *gin.Context) {
	list, err := h.service.ListSongs()
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"list": list}, nil)
}
func (h *AdminHandler) GetSong(c *gin.Context) {
	item, err := h.service.GetSong(c.Param("ref"))
	if err != nil { response.Error(c, http.StatusNotFound, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) CreateSong(c *gin.Context) {
	var in dto.AdminSong
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.CreateSong(in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) UpdateSong(c *gin.Context) {
	var in dto.AdminSong
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.UpdateSong(c.Param("ref"), in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) DeleteSong(c *gin.Context) {
	if err := h.service.DeleteSong(c.Param("ref")); err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"success": true}, nil)
}

func (h *AdminHandler) ListThemes(c *gin.Context) {
	list, err := h.service.ListThemes()
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"list": list}, nil)
}
func (h *AdminHandler) GetTheme(c *gin.Context) {
	item, err := h.service.GetTheme(c.Param("ref"))
	if err != nil { response.Error(c, http.StatusNotFound, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) CreateTheme(c *gin.Context) {
	var in dto.AdminTheme
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.CreateTheme(in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) UpdateTheme(c *gin.Context) {
	var in dto.AdminTheme
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.UpdateTheme(c.Param("ref"), in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminHandler) DeleteTheme(c *gin.Context) {
	if err := h.service.DeleteTheme(c.Param("ref")); err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"success": true}, nil)
}
