package service

import (
	"errors"
	"time"

	"pm-backend/internal/config"
	"pm-backend/internal/dto"

	"github.com/golang-jwt/jwt/v5"
)

type AdminAuthService struct {
	cfg config.Config
}

func NewAdminAuthService(cfg config.Config) *AdminAuthService {
	return &AdminAuthService{cfg: cfg}
}

func (s *AdminAuthService) Login(username, password string) (dto.AdminLoginResponse, error) {
	if username != s.cfg.AdminLoginUsername || password != s.cfg.AdminLoginPassword {
		return dto.AdminLoginResponse{}, errors.New("username or password is invalid")
	}

	now := time.Now()
	expireAt := now.Add(time.Duration(s.cfg.AdminTokenExpireHours) * time.Hour)

	claims := jwt.MapClaims{
		"sub":   "admin-local",
		"typ":   "admin",
		"uname": s.cfg.AdminLoginUsername,
		"roles": []string{"admin"},
		"iat":   now.Unix(),
		"exp":   expireAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.AdminJWTSecret))
	if err != nil {
		return dto.AdminLoginResponse{}, err
	}

	return dto.AdminLoginResponse{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: int64(time.Duration(s.cfg.AdminTokenExpireHours) * time.Hour / time.Second),
		User: dto.AdminUser{
			ID:       "admin-local",
			Username: s.cfg.AdminLoginUsername,
			Nickname: "管理员",
			Roles:    []string{"admin"},
		},
	}, nil
}

func (s *AdminAuthService) LoginForVben(username, password string) (dto.VbenLoginResponse, error) {
	resp, err := s.Login(username, password)
	if err != nil {
		return dto.VbenLoginResponse{}, err
	}
	return dto.VbenLoginResponse{AccessToken: resp.Token}, nil
}

func (s *AdminAuthService) ParseToken(tokenString string) (dto.AdminUser, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AdminJWTSecret), nil
	})
	if err != nil {
		return dto.AdminUser{}, err
	}
	if !token.Valid {
		return dto.AdminUser{}, errors.New("token is invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return dto.AdminUser{}, errors.New("token claims is invalid")
	}
	username, _ := claims["uname"].(string)
	roles := []string{"admin"}
	if rawRoles, ok := claims["roles"].([]interface{}); ok {
		roles = make([]string, 0, len(rawRoles))
		for _, v := range rawRoles {
			if s, ok := v.(string); ok {
				roles = append(roles, s)
			}
		}
	}
	return dto.AdminUser{
		ID:       "admin-local",
		Username: username,
		Nickname: "管理员",
		Roles:    roles,
	}, nil
}

func (s *AdminAuthService) ParseTokenForVben(tokenString string) (dto.VbenUserInfo, error) {
	user, err := s.ParseToken(tokenString)
	if err != nil {
		return dto.VbenUserInfo{}, err
	}
	return dto.VbenUserInfo{
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.Nickname,
		Roles:    user.Roles,
	}, nil
}

func (s *AdminAuthService) PermissionCodes() []string {
	return []string{}
}
