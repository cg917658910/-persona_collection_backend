package handler

import (
	"net/http"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminWorkCreatorHandler struct {
	service *service.AdminWorkCreatorService
}

func NewAdminWorkCreatorHandler(s *service.AdminWorkCreatorService) *AdminWorkCreatorHandler {
	return &AdminWorkCreatorHandler{service: s}
}

func (h *AdminWorkCreatorHandler) ListWorks(c *gin.Context) {
	list, err := h.service.ListWorks()
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"items": list, "total": len(list), "page": 1, "pageSize": len(list)}, nil)
}
func (h *AdminWorkCreatorHandler) GetWork(c *gin.Context) {
	item, err := h.service.GetWork(c.Param("ref"))
	if err != nil { response.Error(c, http.StatusNotFound, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) CreateWork(c *gin.Context) {
	var in dto.AdminWork
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.CreateWork(in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) UpdateWork(c *gin.Context) {
	var in dto.AdminWork
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.UpdateWork(c.Param("ref"), in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) DeleteWork(c *gin.Context) {
	if err := h.service.DeleteWork(c.Param("ref")); err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"success": true}, nil)
}

func (h *AdminWorkCreatorHandler) ListCreators(c *gin.Context) {
	list, err := h.service.ListCreators()
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"items": list, "total": len(list), "page": 1, "pageSize": len(list)}, nil)
}
func (h *AdminWorkCreatorHandler) GetCreator(c *gin.Context) {
	item, err := h.service.GetCreator(c.Param("ref"))
	if err != nil { response.Error(c, http.StatusNotFound, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) CreateCreator(c *gin.Context) {
	var in dto.AdminCreator
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.CreateCreator(in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) UpdateCreator(c *gin.Context) {
	var in dto.AdminCreator
	if err := c.ShouldBindJSON(&in); err != nil { response.Error(c, http.StatusBadRequest, err.Error()); return }
	item, err := h.service.UpdateCreator(c.Param("ref"), in)
	if err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, item, nil)
}
func (h *AdminWorkCreatorHandler) DeleteCreator(c *gin.Context) {
	if err := h.service.DeleteCreator(c.Param("ref")); err != nil { response.Error(c, http.StatusInternalServerError, err.Error()); return }
	response.OK(c, gin.H{"success": true}, nil)
}
