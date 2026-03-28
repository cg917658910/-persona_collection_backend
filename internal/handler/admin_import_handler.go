package handler

import (
	"net/http"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminImportHandler struct {
	service *service.AdminImportService
}

func NewAdminImportHandler(s *service.AdminImportService) *AdminImportHandler {
	return &AdminImportHandler{service: s}
}

func (h *AdminImportHandler) Validate(c *gin.Context) {
	var in dto.AdminImportRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.Validate(in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result, nil)
}

func (h *AdminImportHandler) Run(c *gin.Context) {
	var in dto.AdminImportRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.Run(in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result, nil)
}

func (h *AdminImportHandler) ValidateRelations(c *gin.Context) {
	var in dto.AdminRelationImportRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.ValidateRelations(in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result, nil)
}

func (h *AdminImportHandler) RunRelations(c *gin.Context) {
	var in dto.AdminRelationImportRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.RunRelations(in)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result, nil)
}
