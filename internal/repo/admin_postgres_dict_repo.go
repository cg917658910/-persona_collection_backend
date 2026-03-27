package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
)

type dictTableMeta struct {
	Table      string
	NameColumn string
}

func dictMeta(dictKey string) (dictTableMeta, error) {
	switch dictKey {
	case "characterTypes":
		return dictTableMeta{Table: "public.pm_character_types", NameColumn: "name_zh"}, nil
	case "workTypes":
		return dictTableMeta{Table: "public.pm_work_types", NameColumn: "name_zh"}, nil
	case "creatorTypes":
		return dictTableMeta{Table: "public.pm_creator_types", NameColumn: "name_zh"}, nil
	case "regions":
		return dictTableMeta{Table: "public.pm_regions", NameColumn: "name_zh"}, nil
	case "culturalRegions":
		return dictTableMeta{Table: "public.pm_cultural_regions", NameColumn: "name_zh"}, nil
	case "motivations":
		return dictTableMeta{Table: "public.pm_motivation_dict", NameColumn: "name_zh"}, nil
	case "themeCategories":
		return dictTableMeta{Table: "public.pm_theme_categories", NameColumn: "name_zh"}, nil
	default:
		return dictTableMeta{}, errors.New("dictKey is invalid")
	}
}

func (r *postgresAdminRepo) ListAdminDictItems(dictKey string) ([]dto.AdminDictItem, error) {
	meta, err := dictMeta(dictKey)
	if err != nil { return nil, err }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
SELECT id::text, code, COALESCE(%s,''), COALESCE(sort_order,0), COALESCE(is_active,TRUE)
FROM %s
ORDER BY sort_order ASC, code ASC
`, meta.NameColumn, meta.Table))
	if err != nil { return nil, fmt.Errorf("list dict query: %w", err) }
	defer rows.Close()

	list := make([]dto.AdminDictItem, 0)
	for rows.Next() {
		var item dto.AdminDictItem
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.SortOrder, &item.IsActive); err != nil {
			return nil, fmt.Errorf("scan dict item: %w", err)
		}
		item.DictKey = dictKey
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresAdminRepo) GetAdminDictItem(dictKey, ref string) (dto.AdminDictItem, error) {
	meta, err := dictMeta(dictKey)
	if err != nil { return dto.AdminDictItem{}, err }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item dto.AdminDictItem
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`
SELECT id::text, code, COALESCE(%s,''), COALESCE(sort_order,0), COALESCE(is_active,TRUE)
FROM %s
WHERE id::text=$1 OR code=$1
LIMIT 1
`, meta.NameColumn, meta.Table), ref).Scan(&item.ID, &item.Code, &item.Name, &item.SortOrder, &item.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return dto.AdminDictItem{}, errors.New("admin dict item not found") }
		return dto.AdminDictItem{}, err
	}
	item.DictKey = dictKey
	return item, nil
}

func (r *postgresAdminRepo) CreateAdminDictItem(dictKey string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	meta, err := dictMeta(dictKey)
	if err != nil { return dto.AdminDictItem{}, err }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`
INSERT INTO %s (code, %s, sort_order, is_active)
VALUES ($1,$2,$3,COALESCE($4,TRUE))
RETURNING id::text
`, meta.Table, meta.NameColumn), in.Code, in.Name, in.SortOrder, in.IsActive).Scan(&id)
	if err != nil { return dto.AdminDictItem{}, fmt.Errorf("insert dict item: %w", err) }
	return r.GetAdminDictItem(dictKey, id)
}

func (r *postgresAdminRepo) UpdateAdminDictItem(dictKey, ref string, in dto.AdminDictItem) (dto.AdminDictItem, error) {
	meta, err := dictMeta(dictKey)
	if err != nil { return dto.AdminDictItem{}, err }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := resolveIDByRef(ctx, r.pool, meta.Table, ref)
	if err != nil { return dto.AdminDictItem{}, err }
	_, err = r.pool.Exec(ctx, fmt.Sprintf(`
UPDATE %s
SET code=$2, %s=$3, sort_order=$4, is_active=$5, updated_at=NOW()
WHERE id=$1
`, meta.Table, meta.NameColumn), id, in.Code, in.Name, in.SortOrder, in.IsActive)
	if err != nil { return dto.AdminDictItem{}, fmt.Errorf("update dict item: %w", err) }
	return r.GetAdminDictItem(dictKey, id)
}

func (r *postgresAdminRepo) DeleteAdminDictItem(dictKey, ref string) error {
	meta, err := dictMeta(dictKey)
	if err != nil { return err }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET is_active=FALSE, updated_at=NOW() WHERE id::text=$1 OR code=$1`, meta.Table), ref)
	return err
}
