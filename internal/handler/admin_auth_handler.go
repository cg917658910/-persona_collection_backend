package handler

import (
	"net/http"
	"strings"

	"pm-backend/internal/dto"
	"pm-backend/internal/service"
	"pm-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminAuthHandler struct {
	service *service.AdminAuthService
}

func NewAdminAuthHandler(s *service.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{service: s}
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.service.Login(strings.TrimSpace(req.Username), req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, resp, nil)
}

func (h *AdminAuthHandler) Me(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		response.Error(c, http.StatusUnauthorized, "authorization header is required")
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		response.Error(c, http.StatusUnauthorized, "authorization header format is invalid")
		return
	}
	user, err := h.service.ParseToken(strings.TrimSpace(parts[1]))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, user, nil)
}

func (h *AdminAuthHandler) VbenLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.service.LoginForVben(strings.TrimSpace(req.Username), req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, resp, nil)
}

func (h *AdminAuthHandler) VbenUserInfo(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		response.Error(c, http.StatusUnauthorized, "authorization header is required")
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		response.Error(c, http.StatusUnauthorized, "authorization header format is invalid")
		return
	}
	user, err := h.service.ParseTokenForVben(strings.TrimSpace(parts[1]))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, user, nil)
}

func (h *AdminAuthHandler) VbenPermissionCodes(c *gin.Context) {
	response.OK(c, h.service.PermissionCodes(), nil)
}
