package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
)

func (r *postgresAdminRepo) ListAdminWorks() ([]dto.AdminWork, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.pool.Query(ctx, `
SELECT
  w.id::text,
  w.slug,
  w.title,
  COALESCE(w.summary,''),
  COALESCE(w.cover_url,''),
  COALESCE(w.sort_order,0),
  COALESCE(NULLIF(w.meta->>'status',''), CASE WHEN w.is_active THEN 'published' ELSE 'archived' END) AS status,
  COALESCE(wt.code,'') AS work_type_code,
  COALESCE(rg.code,'') AS region_code,
  COALESCE(cr.code,'') AS cultural_region_code,
  COALESCE(w.release_year,0),
  COALESCE(w.meta::text, '{}'::text),
  COALESCE((
    SELECT array_agg(c.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE
  ), ARRAY[]::text[]),
  COALESCE((
    SELECT array_agg(c.name ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE
  ), ARRAY[]::text[])
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
LEFT JOIN public.pm_regions rg ON rg.id = w.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = w.cultural_region_id
ORDER BY w.sort_order ASC, w.updated_at DESC, w.title ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin works query: %w", err)
	}
	defer rows.Close()
	list := make([]dto.AdminWork, 0)
	for rows.Next() {
		var item dto.AdminWork
		var metaJSON string
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.WorkTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.ReleaseYear, &metaJSON, &item.CreatorSlugs, &item.CreatorNames); err != nil {
			return nil, fmt.Errorf("scan admin work: %w", err)
		}
		applyAdminWorkMeta(&item, metaJSON)
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresAdminRepo) GetAdminWork(ref string) (dto.AdminWork, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var item dto.AdminWork
	var metaJSON string
	err := r.pool.QueryRow(ctx, `
SELECT
  w.id::text,
  w.slug,
  w.title,
  COALESCE(w.summary,''),
  COALESCE(w.cover_url,''),
  COALESCE(w.sort_order,0),
  COALESCE(NULLIF(w.meta->>'status',''), CASE WHEN w.is_active THEN 'published' ELSE 'archived' END) AS status,
  COALESCE(wt.code,'') AS work_type_code,
  COALESCE(rg.code,'') AS region_code,
  COALESCE(cr.code,'') AS cultural_region_code,
  COALESCE(w.release_year,0),
  COALESCE(w.meta::text, '{}'::text),
  COALESCE((
    SELECT array_agg(c.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE
  ), ARRAY[]::text[]),
  COALESCE((
    SELECT array_agg(c.name ORDER BY wc.is_primary DESC, wc.sort_order ASC, c.name ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_creators c ON c.id = wc.creator_id
    WHERE wc.work_id = w.id AND c.is_active = TRUE
  ), ARRAY[]::text[])
FROM public.pm_works w
LEFT JOIN public.pm_work_types wt ON wt.id = w.work_type_id
LEFT JOIN public.pm_regions rg ON rg.id = w.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = w.cultural_region_id
WHERE w.id::text=$1 OR w.slug=$1
LIMIT 1
`, ref).Scan(&item.ID, &item.Slug, &item.Title, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.WorkTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.ReleaseYear, &metaJSON, &item.CreatorSlugs, &item.CreatorNames)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AdminWork{}, errors.New("admin work not found")
		}
		return dto.AdminWork{}, err
	}
	applyAdminWorkMeta(&item, metaJSON)
	return item, nil
}

func (r *postgresAdminRepo) CreateAdminWork(in dto.AdminWork) (dto.AdminWork, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminWork{}, err
	}
	defer tx.Rollback(ctx)

	workTypeID, err := optionalLookupIDByCode(ctx, tx, "public.pm_work_types", in.WorkTypeCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup work type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup cultural region: %w", err)
	}
	metaJSON, _ := json.Marshal(buildAdminWorkMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	var id string
	err = tx.QueryRow(ctx, `
INSERT INTO public.pm_works (title, slug, summary, cover_url, work_type_id, region_id, cultural_region_id, release_year, meta, sort_order, is_active)
VALUES ($1,$2,$3,$4,$5,$6,$7,NULLIF($8,0),COALESCE($9::jsonb,'{}'::jsonb),COALESCE($10,0),$11)
RETURNING id::text
`, in.Title, in.Slug, in.Summary, in.CoverURL, workTypeID, regionID, culturalID, in.ReleaseYear, string(metaJSON), in.SortOrder, isPublishedStatus(status)).Scan(&id)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("insert work: %w", err)
	}

	if err := replaceWorkCreators(ctx, tx, id, in.CreatorSlugs); err != nil {
		return dto.AdminWork{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminWork{}, err
	}
	return r.GetAdminWork(id)
}

func (r *postgresAdminRepo) UpdateAdminWork(ref string, in dto.AdminWork) (dto.AdminWork, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminWork{}, err
	}
	defer tx.Rollback(ctx)

	workID, err := resolveIDByRef(ctx, tx, "public.pm_works", ref)
	if err != nil {
		return dto.AdminWork{}, err
	}
	workTypeID, err := optionalLookupIDByCode(ctx, tx, "public.pm_work_types", in.WorkTypeCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup work type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("lookup cultural region: %w", err)
	}
	metaJSON, _ := json.Marshal(buildAdminWorkMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	_, err = tx.Exec(ctx, `
UPDATE public.pm_works
SET title=$2, slug=$3, summary=$4, cover_url=$5, work_type_id=$6, region_id=$7, cultural_region_id=$8, release_year=NULLIF($9,0),
    meta=COALESCE(pm_works.meta, '{}'::jsonb) || COALESCE($10::jsonb, '{}'::jsonb), sort_order=COALESCE($11, pm_works.sort_order), is_active=$12, updated_at=NOW()
WHERE id=$1
`, workID, in.Title, in.Slug, in.Summary, in.CoverURL, workTypeID, regionID, culturalID, in.ReleaseYear, string(metaJSON), in.SortOrder, isPublishedStatus(status))
	if err != nil {
		return dto.AdminWork{}, fmt.Errorf("update work: %w", err)
	}

	if err := replaceWorkCreators(ctx, tx, workID, in.CreatorSlugs); err != nil {
		return dto.AdminWork{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminWork{}, err
	}
	return r.GetAdminWork(workID)
}

func (r *postgresAdminRepo) DeleteAdminWork(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := r.pool.Exec(ctx, `UPDATE public.pm_works SET is_active=FALSE, meta=COALESCE(meta, '{}'::jsonb) || '{"status":"archived"}'::jsonb, updated_at=NOW() WHERE id::text=$1 OR slug=$1`, ref)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin work not found")
	}
	return nil
}

func (r *postgresAdminRepo) ListAdminCreators() ([]dto.AdminCreator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.pool.Query(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.summary,''),
  COALESCE(c.cover_url,''),
  COALESCE(c.sort_order,0),
  COALESCE(NULLIF(c.meta->>'status',''), CASE WHEN c.is_active THEN 'published' ELSE 'archived' END) AS status,
  COALESCE(ct.code,'') AS creator_type_code,
  COALESCE(rg.code,'') AS region_code,
  COALESCE(cr.code,'') AS cultural_region_code,
  COALESCE((
    SELECT array_agg(w.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE
  ), ARRAY[]::text[]),
  COALESCE((
    SELECT array_agg(w.title ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE
  ), ARRAY[]::text[])
FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
ORDER BY c.sort_order ASC, c.updated_at DESC, c.name ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list admin creators query: %w", err)
	}
	defer rows.Close()
	list := make([]dto.AdminCreator, 0)
	for rows.Next() {
		var item dto.AdminCreator
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.CreatorTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.WorkSlugs, &item.WorkNames); err != nil {
			return nil, fmt.Errorf("scan admin creator: %w", err)
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *postgresAdminRepo) GetAdminCreator(ref string) (dto.AdminCreator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var item dto.AdminCreator
	err := r.pool.QueryRow(ctx, `
SELECT
  c.id::text,
  c.slug,
  c.name,
  COALESCE(c.summary,''),
  COALESCE(c.cover_url,''),
  COALESCE(c.sort_order,0),
  COALESCE(NULLIF(c.meta->>'status',''), CASE WHEN c.is_active THEN 'published' ELSE 'archived' END) AS status,
  COALESCE(ct.code,'') AS creator_type_code,
  COALESCE(rg.code,'') AS region_code,
  COALESCE(cr.code,'') AS cultural_region_code,
  COALESCE((
    SELECT array_agg(w.slug ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE
  ), ARRAY[]::text[]),
  COALESCE((
    SELECT array_agg(w.title ORDER BY wc.is_primary DESC, wc.sort_order ASC, w.title ASC)
    FROM public.pm_work_creators wc
    JOIN public.pm_works w ON w.id = wc.work_id
    WHERE wc.creator_id = c.id AND w.is_active = TRUE
  ), ARRAY[]::text[])
FROM public.pm_creators c
LEFT JOIN public.pm_creator_types ct ON ct.id = c.creator_type_id
LEFT JOIN public.pm_regions rg ON rg.id = c.region_id
LEFT JOIN public.pm_cultural_regions cr ON cr.id = c.cultural_region_id
WHERE c.id::text=$1 OR c.slug=$1
LIMIT 1
`, ref).Scan(&item.ID, &item.Slug, &item.Name, &item.Summary, &item.CoverURL, &item.SortOrder, &item.Status, &item.CreatorTypeCode, &item.RegionCode, &item.CulturalRegionCode, &item.WorkSlugs, &item.WorkNames)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AdminCreator{}, errors.New("admin creator not found")
		}
		return dto.AdminCreator{}, err
	}
	return item, nil
}

func (r *postgresAdminRepo) CreateAdminCreator(in dto.AdminCreator) (dto.AdminCreator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminCreator{}, err
	}
	defer tx.Rollback(ctx)

	creatorTypeID, err := optionalLookupIDByCode(ctx, tx, "public.pm_creator_types", in.CreatorTypeCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup creator type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup cultural region: %w", err)
	}
	metaJSON, _ := json.Marshal(buildAdminCreatorMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	var id string
	err = tx.QueryRow(ctx, `
INSERT INTO public.pm_creators (name, slug, summary, cover_url, creator_type_id, region_id, cultural_region_id, meta, sort_order, is_active)
VALUES ($1,$2,$3,$4,$5,$6,$7,COALESCE($8::jsonb,'{}'::jsonb),COALESCE($9,0),$10)
RETURNING id::text
`, in.Name, in.Slug, in.Summary, in.CoverURL, creatorTypeID, regionID, culturalID, string(metaJSON), in.SortOrder, isPublishedStatus(status)).Scan(&id)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("insert creator: %w", err)
	}

	if err := replaceCreatorWorks(ctx, tx, id, in.WorkSlugs); err != nil {
		return dto.AdminCreator{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminCreator{}, err
	}
	return r.GetAdminCreator(id)
}

func (r *postgresAdminRepo) UpdateAdminCreator(ref string, in dto.AdminCreator) (dto.AdminCreator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return dto.AdminCreator{}, err
	}
	defer tx.Rollback(ctx)

	creatorID, err := resolveIDByRef(ctx, tx, "public.pm_creators", ref)
	if err != nil {
		return dto.AdminCreator{}, err
	}
	creatorTypeID, err := optionalLookupIDByCode(ctx, tx, "public.pm_creator_types", in.CreatorTypeCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup creator type: %w", err)
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", in.RegionCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup region: %w", err)
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", in.CulturalRegionCode)
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("lookup cultural region: %w", err)
	}
	metaJSON, _ := json.Marshal(buildAdminCreatorMeta(in))
	status := normalizeAdminStatus(in.Status, "draft")

	_, err = tx.Exec(ctx, `
UPDATE public.pm_creators
SET name=$2, slug=$3, summary=$4, cover_url=$5, creator_type_id=$6, region_id=$7, cultural_region_id=$8,
    meta=COALESCE(pm_creators.meta, '{}'::jsonb) || COALESCE($9::jsonb, '{}'::jsonb),
    sort_order=COALESCE($10, pm_creators.sort_order), is_active=$11, updated_at=NOW()
WHERE id=$1
`, creatorID, in.Name, in.Slug, in.Summary, in.CoverURL, creatorTypeID, regionID, culturalID, string(metaJSON), in.SortOrder, isPublishedStatus(status))
	if err != nil {
		return dto.AdminCreator{}, fmt.Errorf("update creator: %w", err)
	}

	if err := replaceCreatorWorks(ctx, tx, creatorID, in.WorkSlugs); err != nil {
		return dto.AdminCreator{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AdminCreator{}, err
	}
	return r.GetAdminCreator(creatorID)
}

func (r *postgresAdminRepo) DeleteAdminCreator(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := r.pool.Exec(ctx, `UPDATE public.pm_creators SET is_active=FALSE, meta=COALESCE(meta, '{}'::jsonb) || '{"status":"archived"}'::jsonb, updated_at=NOW() WHERE id::text=$1 OR slug=$1`, ref)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("admin creator not found")
	}
	return nil
}

func replaceWorkCreators(ctx context.Context, tx pgx.Tx, workID string, creatorSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_work_creators WHERE work_id=$1`, workID); err != nil {
		return err
	}
	for idx, slug := range uniq(creatorSlugs) {
		creatorID, err := resolveIDByRef(ctx, tx, "public.pm_creators", slug)
		if err != nil {
			return fmt.Errorf("lookup creator %s: %w", slug, err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.pm_work_creators (work_id, creator_id, role_code, is_primary, sort_order) VALUES ($1,$2,'author',$3,$4)`, workID, creatorID, idx == 0, (idx+1)*10)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceCreatorWorks(ctx context.Context, tx pgx.Tx, creatorID string, workSlugs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_work_creators WHERE creator_id=$1`, creatorID); err != nil {
		return err
	}
	for idx, slug := range uniq(workSlugs) {
		workID, err := resolveIDByRef(ctx, tx, "public.pm_works", slug)
		if err != nil {
			return fmt.Errorf("lookup work %s: %w", slug, err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.pm_work_creators (work_id, creator_id, role_code, is_primary, sort_order) VALUES ($1,$2,'author',$3,$4)`, workID, creatorID, idx == 0, (idx+1)*10)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildAdminWorkMeta(in dto.AdminWork) map[string]any {
	return map[string]any{
		"status":         normalizeAdminStatus(in.Status, "draft"),
		"is_recommended": in.Recommended,
		"recommend_sort": in.RecommendSort,
	}
}

func applyAdminWorkMeta(item *dto.AdminWork, metaJSON string) {
	var meta map[string]any
	if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
		return
	}
	item.Recommended = adminMetaBool(meta, "is_recommended")
	item.RecommendSort = adminMetaInt(meta, "recommend_sort")
	if v, ok := meta["status"].(string); ok && strings.TrimSpace(v) != "" {
		item.Status = normalizeAdminStatus(v, item.Status)
	}
}

func buildAdminCreatorMeta(in dto.AdminCreator) map[string]any {
	return map[string]any{
		"status": normalizeAdminStatus(in.Status, "draft"),
	}
}
