package handler

import (
	"net/http"
	"strconv"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminResourceHandler struct {
	service *service.AdminResourceService
}

func NewAdminResourceHandler(s *service.AdminResourceService) *AdminResourceHandler {
	return &AdminResourceHandler{service: s}
}

func (h *AdminResourceHandler) ListResources(c *gin.Context) {
	resourceType := dto.AdminResourceType(c.Param("type"))
	keyword := c.Query("keyword")
	list, err := h.service.ListResources(resourceType, keyword)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	page, pageSize := dto.NormalizePage(toIntDefault(c.Query("page"), 1), toIntDefault(c.Query("pageSize"), 10))
	total := len(list)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	response.OK(c, gin.H{
		"items":    list[start:end],
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, nil)
}

func (h *AdminResourceHandler) CreateResource(c *gin.Context) {
	resourceType := dto.AdminResourceType(c.Param("type"))
	var req dto.AdminCreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.Type == "" {
		req.Type = resourceType
	}
	item, err := h.service.CreateResource(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, item, nil)
}

func (h *AdminResourceHandler) UploadResource(c *gin.Context) {
	resourceType := dto.AdminResourceType(c.Param("type"))
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.UploadResource(resourceType, fileHeader)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, item, nil)
}

func (h *AdminResourceHandler) DeleteResource(c *gin.Context) {
	resourceType := dto.AdminResourceType(c.Param("type"))
	ref := c.Param("ref")
	if err := h.service.DeleteResource(resourceType, ref); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"success": true}, nil)
}


func toIntDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}
