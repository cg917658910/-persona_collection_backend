package service

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"pm-backend/internal/config"
	"pm-backend/internal/dto"
)

type AdminResourceService struct {
	cfg    config.Config
	assets *AssetURLBuilder
}

func NewAdminResourceService(cfg config.Config, assets *AssetURLBuilder) *AdminResourceService {
	return &AdminResourceService{cfg: cfg, assets: assets}
}

func (s *AdminResourceService) ListResources(resourceType dto.AdminResourceType, keyword string) ([]dto.AdminResourceItem, error) {
	baseDir, urlPrefix, err := s.typePaths(resourceType)
	if err != nil {
		return nil, err
	}
	entries := make([]dto.AdminResourceItem, 0)
	if _, err := os.Stat(baseDir); err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(s.cfg.StaticLocalDir, path)
		if err != nil {
			return err
		}
		url := s.assets.Normalize(filepath.ToSlash(rel))
		item := dto.AdminResourceItem{
			ID:        filepath.ToSlash(rel),
			Name:      info.Name(),
			URL:       url,
			Type:      resourceType,
			MimeType:  mime.TypeByExtension(filepath.Ext(info.Name())),
			Size:      info.Size(),
			CreatedAt: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		if keyword == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(keyword)) || strings.Contains(strings.ToLower(item.URL), strings.ToLower(keyword)) || strings.Contains(strings.ToLower(item.MimeType), strings.ToLower(keyword)) {
			entries = append(entries, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt > entries[j].CreatedAt
	})
	_ = urlPrefix
	return entries, nil
}

func (s *AdminResourceService) CreateResource(in dto.AdminCreateResourceRequest) (dto.AdminResourceItem, error) {
	if strings.TrimSpace(in.URL) == "" {
		return dto.AdminResourceItem{}, errors.New("url is required")
	}
	item := dto.AdminResourceItem{
		ID:           s.assets.ToStorage(in.URL),
		Name:         chooseName(in.Name, filepath.Base(in.URL)),
		URL:          s.assets.Normalize(in.URL),
		Type:         in.Type,
		MimeType:     in.MimeType,
		Size:         in.Size,
		LinkedModule: in.LinkedModule,
		LinkedCount:  in.LinkedCount,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}
	return item, nil
}

func (s *AdminResourceService) UploadResource(resourceType dto.AdminResourceType, fileHeader *multipart.FileHeader) (dto.AdminResourceItem, error) {
	baseDir, _, err := s.typePaths(resourceType)
	if err != nil {
		return dto.AdminResourceItem{}, err
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return dto.AdminResourceItem{}, err
	}
	safeName := sanitizeFileName(fileHeader.Filename)
	targetPath := filepath.Join(baseDir, fmt.Sprintf("%d_%s", time.Now().Unix(), safeName))

	src, err := fileHeader.Open()
	if err != nil {
		return dto.AdminResourceItem{}, err
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return dto.AdminResourceItem{}, err
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		return dto.AdminResourceItem{}, err
	}

	rel, err := filepath.Rel(s.cfg.StaticLocalDir, targetPath)
	if err != nil {
		return dto.AdminResourceItem{}, err
	}
	storagePath := filepath.ToSlash(rel)
	item := dto.AdminResourceItem{
		ID:        storagePath,
		Name:      filepath.Base(targetPath),
		URL:       s.assets.Normalize(storagePath),
		Type:      resourceType,
		MimeType:  mime.TypeByExtension(filepath.Ext(fileHeader.Filename)),
		Size:      written,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	return item, nil
}

func (s *AdminResourceService) DeleteResource(resourceType dto.AdminResourceType, ref string) error {
	if strings.TrimSpace(ref) == "" {
		return errors.New("resource ref is required")
	}
	storagePath := s.assets.ToStorage(ref)
	if storagePath == "" {
		storagePath = ref
	}
	storagePath = strings.TrimPrefix(storagePath, "/")
	fullPath := filepath.Join(s.cfg.StaticLocalDir, filepath.FromSlash(storagePath))
	baseDir, _, err := s.typePaths(resourceType)
	if err != nil {
		return err
	}
	cleanBase, _ := filepath.Abs(baseDir)
	cleanTarget, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(cleanTarget, cleanBase) {
		return errors.New("resource path is invalid")
	}
	if _, err := os.Stat(cleanTarget); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(cleanTarget)
}

func (s *AdminResourceService) typePaths(resourceType dto.AdminResourceType) (string, string, error) {
	switch resourceType {
	case dto.AdminResourceTypeImage:
		return filepath.Join(s.cfg.StaticLocalDir, "assets", "images"), "/assets/images", nil
	case dto.AdminResourceTypeAudio:
		return filepath.Join(s.cfg.StaticLocalDir, "assets", "audio"), "/assets/audio", nil
	default:
		return "", "", errors.New("resource type is invalid")
	}
}

var invalidFileChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = invalidFileChars.ReplaceAllString(name, "_")
	if name == "" {
		return "file"
	}
	return name
}

func chooseName(name, fallback string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	return fallback
}
