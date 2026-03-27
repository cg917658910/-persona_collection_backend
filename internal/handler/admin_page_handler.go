package handler

import (
	"net/http"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminPageHandler struct {
	service *service.AdminPageService
}

func NewAdminPageHandler(s *service.AdminPageService) *AdminPageHandler {
	return &AdminPageHandler{service: s}
}

func (h *AdminPageHandler) Characters(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Characters(q)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
func (h *AdminPageHandler) Songs(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Songs(q)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
func (h *AdminPageHandler) Themes(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Themes(q)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
func (h *AdminPageHandler) Works(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Works(q)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
func (h *AdminPageHandler) Creators(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Creators(q)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
func (h *AdminPageHandler) Dicts(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	res, err := h.service.Dicts(c.Param("dictKey"), q.Page, q.PageSize, q.Keyword)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, res, nil)
}
