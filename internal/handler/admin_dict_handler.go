package handler

import (
	"net/http"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminDictHandler struct {
	service *service.AdminDictService
}

func NewAdminDictHandler(s *service.AdminDictService) *AdminDictHandler {
	return &AdminDictHandler{service: s}
}

func (h *AdminDictHandler) List(c *gin.Context) {
	dictKey := c.Param("dictKey")
	items, err := h.service.List(dictKey)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{
		"items":    items,
		"total":    len(items),
		"page":     1,
		"pageSize": len(items),
	}, nil)
}

func (h *AdminDictHandler) Get(c *gin.Context) {
	dictKey := c.Param("dictKey")
	ref := c.Param("ref")
	item, err := h.service.Get(dictKey, ref)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, item, nil)
}

func (h *AdminDictHandler) Create(c *gin.Context) {
	dictKey := c.Param("dictKey")
	var in dto.AdminDictItem
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Create(dictKey, in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, item, nil)
}

func (h *AdminDictHandler) Update(c *gin.Context) {
	dictKey := c.Param("dictKey")
	ref := c.Param("ref")
	var in dto.AdminDictItem
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Update(dictKey, ref, in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, item, nil)
}

func (h *AdminDictHandler) Delete(c *gin.Context) {
	dictKey := c.Param("dictKey")
	ref := c.Param("ref")
	if err := h.service.Delete(dictKey, ref); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"success": true}, nil)
}
