package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"pm-backend/internal/dto"

	"github.com/jackc/pgx/v5"
)

const generatedPackageVersion = "pm-character-gen-v2"

func (r *postgresAdminRepo) ValidateGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := newImportResult(pkg)
	errors := make([]string, 0)
	warnings := make([]string, 0)

	if strings.TrimSpace(pkg.PackageVersion) != generatedPackageVersion {
		errors = append(errors, fmt.Sprintf("unsupported package_version: %s", pkg.PackageVersion))
	}

	charTypeCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_character_types")
	if err != nil {
		return result, err
	}
	creatorTypeCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_creator_types")
	if err != nil {
		return result, err
	}
	workTypeCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_work_types")
	if err != nil {
		return result, err
	}
	motivationCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_motivation_dict")
	if err != nil {
		return result, err
	}
	regionCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_regions")
	if err != nil {
		return result, err
	}
	culturalCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_cultural_regions")
	if err != nil {
		return result, err
	}
	existingThemeCodes, err := fetchCodeSet(ctx, r.pool, "public.pm_themes")
	if err != nil {
		return result, err
	}
	existingCreatorSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_creators")
	if err != nil {
		return result, err
	}
	existingWorkSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_works")
	if err != nil {
		return result, err
	}
	existingCharacterSlugs, err := fetchSlugSet(ctx, r.pool, "public.pm_characters")
	if err != nil {
		return result, err
	}

	themeCodesInPackage := make(map[string]struct{}, len(pkg.PmThemes))
	creatorSlugsInPackage := make(map[string]struct{}, len(pkg.PmCreators))
	workSlugsInPackage := make(map[string]struct{}, len(pkg.PmWorks))
	characterSlugsInPackage := make(map[string]struct{}, len(pkg.PmCharacters))

	errors = append(errors, validateUniqueThemes(pkg.PmThemes, themeCodesInPackage)...)
	errors = append(errors, validateUniqueCreators(pkg.PmCreators, creatorSlugsInPackage)...)
	errors = append(errors, validateUniqueWorks(pkg.PmWorks, workSlugsInPackage)...)
	errors = append(errors, validateUniqueCharacters(pkg.PmCharacters, characterSlugsInPackage)...)
	errors = append(errors, validateUniqueSongs(pkg.PmSongs)...)

	for _, item := range pkg.PmCreators {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Slug) == "" {
			errors = append(errors, "creator requires name and slug")
		}
		if !hasCode(creatorTypeCodes, item.CreatorTypeCode) {
			errors = append(errors, fmt.Sprintf("creator %s references missing creator_type_code %s", item.Slug, item.CreatorTypeCode))
		}
		warnings = append(warnings, validateOptionalAutocreateCode("creator "+item.Slug+" region_code", item.RegionCode, regionCodes, "pm_regions")...)
		warnings = append(warnings, validateOptionalAutocreateCode("creator "+item.Slug+" cultural_region_code", item.CulturalRegionCode, culturalCodes, "pm_cultural_regions")...)
	}

	for _, item := range pkg.PmWorks {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Slug) == "" {
			errors = append(errors, "work requires title and slug")
		}
		warnings = append(warnings, validateAutocreateRequiredCode("work "+item.Slug+" work_type_code", item.WorkTypeCode, workTypeCodes, "pm_work_types")...)
		warnings = append(warnings, validateOptionalAutocreateCode("work "+item.Slug+" region_code", item.RegionCode, regionCodes, "pm_regions")...)
		warnings = append(warnings, validateOptionalAutocreateCode("work "+item.Slug+" cultural_region_code", item.CulturalRegionCode, culturalCodes, "pm_cultural_regions")...)
		for _, role := range item.CreatorRoles {
			if strings.TrimSpace(role.CreatorSlug) == "" {
				errors = append(errors, fmt.Sprintf("work %s has empty creator_slug in creator_roles", item.Slug))
				continue
			}
			if !hasSlug(creatorSlugsInPackage, existingCreatorSlugs, role.CreatorSlug) {
				errors = append(errors, fmt.Sprintf("work %s references missing creator %s", item.Slug, role.CreatorSlug))
			}
		}
	}

	for _, item := range pkg.PmCharacters {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Slug) == "" {
			errors = append(errors, "character requires name and slug")
		}
		if !hasCode(charTypeCodes, item.CharacterTypeCode) {
			errors = append(errors, fmt.Sprintf("character %s references missing character_type_code %s", item.Slug, item.CharacterTypeCode))
		}
		warnings = append(warnings, validateOptionalAutocreateCode("character "+item.Slug+" region_code", item.RegionCode, regionCodes, "pm_regions")...)
		warnings = append(warnings, validateOptionalAutocreateCode("character "+item.Slug+" cultural_region_code", item.CulturalRegionCode, culturalCodes, "pm_cultural_regions")...)
		errors = append(errors, validateStatus("character "+item.Slug, item.Status)...)
		for _, code := range item.MotivationCodes {
			if !hasCode(motivationCodes, code) {
				errors = append(errors, fmt.Sprintf("character %s references missing motivation_code %s", item.Slug, code))
			}
		}
		if item.PrimaryMotivation != "" && !slices.Contains(item.MotivationCodes, item.PrimaryMotivation) {
			warnings = append(warnings, fmt.Sprintf("character %s primary_motivation_code is not in motivation_codes", item.Slug))
		}
		for _, code := range item.ThemeCodes {
			if !hasThemeCode(themeCodesInPackage, existingThemeCodes, code) {
				errors = append(errors, fmt.Sprintf("character %s references missing theme_code %s", item.Slug, code))
			}
		}
		if item.PrimaryTheme != "" && !slices.Contains(item.ThemeCodes, item.PrimaryTheme) {
			warnings = append(warnings, fmt.Sprintf("character %s primary_theme_code is not in theme_codes", item.Slug))
		}
		for _, slug := range item.WorkSlugs {
			if !hasSlug(workSlugsInPackage, existingWorkSlugs, slug) {
				errors = append(errors, fmt.Sprintf("character %s references missing work_slug %s", item.Slug, slug))
			}
		}
		if item.PrimaryWork != "" && !slices.Contains(item.WorkSlugs, item.PrimaryWork) {
			warnings = append(warnings, fmt.Sprintf("character %s primary_work_slug is not in work_slugs", item.Slug))
		}
	}

	for _, item := range pkg.PmSongs {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Slug) == "" {
			errors = append(errors, "song requires title and slug")
		}
		if !hasSlug(characterSlugsInPackage, existingCharacterSlugs, item.CharacterSlug) {
			errors = append(errors, fmt.Sprintf("song %s references missing character %s", item.Slug, item.CharacterSlug))
		}
		errors = append(errors, validateStatus("song "+item.Slug, item.Status)...)
	}

	result.Warnings = warnings
	result.Errors = dedupeStrings(errors)
	result.Valid = len(result.Errors) == 0
	return result, nil
}

func (r *postgresAdminRepo) ImportGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error) {
	result, err := r.ValidateGeneratedPackage(pkg)
	if err != nil {
		return result, err
	}
	if !result.Valid {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)

	autoWarnings, err := ensureGeneratedPackageDictionaries(ctx, tx, pkg)
	if err != nil {
		return result, fmt.Errorf("ensure import dictionaries: %w", err)
	}
	result.Warnings = dedupeStrings(append(result.Warnings, autoWarnings...))

	for _, item := range pkg.PmThemes {
		if err := upsertGeneratedTheme(ctx, tx, item); err != nil {
			return result, fmt.Errorf("upsert theme %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmCreators {
		if err := upsertGeneratedCreator(ctx, tx, item); err != nil {
			return result, fmt.Errorf("upsert creator %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmWorks {
		if err := upsertGeneratedWork(ctx, tx, item); err != nil {
			return result, fmt.Errorf("upsert work %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmCharacters {
		if err := upsertGeneratedCharacter(ctx, tx, item); err != nil {
			return result, fmt.Errorf("upsert character %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmWorks {
		if err := replaceGeneratedWorkCreators(ctx, tx, item); err != nil {
			return result, fmt.Errorf("sync work creators %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmCharacters {
		if err := replaceGeneratedCharacterWorks(ctx, tx, item); err != nil {
			return result, fmt.Errorf("sync character works %s: %w", item.Slug, err)
		}
		if err := replaceGeneratedCharacterMotivations(ctx, tx, item); err != nil {
			return result, fmt.Errorf("sync character motivations %s: %w", item.Slug, err)
		}
		if err := replaceGeneratedCharacterThemes(ctx, tx, item); err != nil {
			return result, fmt.Errorf("sync character themes %s: %w", item.Slug, err)
		}
	}
	for _, item := range pkg.PmSongs {
		if err := upsertGeneratedSong(ctx, tx, item); err != nil {
			return result, fmt.Errorf("upsert song %s: %w", item.Slug, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return result, err
	}
	result.Imported = true
	return result, nil
}

func (r *mockAdminRepo) ValidateGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error) {
	result := newImportResult(pkg)
	if strings.TrimSpace(pkg.PackageVersion) != generatedPackageVersion {
		result.Errors = []string{fmt.Sprintf("unsupported package_version: %s", pkg.PackageVersion)}
	} else {
		result.Warnings = []string{"mock mode: validation only checks package_version"}
	}
	result.Valid = len(result.Errors) == 0
	return result, nil
}

func (r *mockAdminRepo) ImportGeneratedPackage(pkg dto.GeneratedPackage) (dto.AdminImportResult, error) {
	result, err := r.ValidateGeneratedPackage(pkg)
	if err != nil {
		return result, err
	}
	if result.Valid {
		result.Warnings = append(result.Warnings, "mock mode: package was not written to a database")
	}
	return result, nil
}

func newImportResult(pkg dto.GeneratedPackage) dto.AdminImportResult {
	return dto.AdminImportResult{
		PackageVersion: pkg.PackageVersion,
		Summary: dto.AdminImportSummary{
			Themes:     len(pkg.PmThemes),
			Creators:   len(pkg.PmCreators),
			Works:      len(pkg.PmWorks),
			Characters: len(pkg.PmCharacters),
			Songs:      len(pkg.PmSongs),
		},
	}
}

func validateUniqueThemes(items []dto.GeneratedTheme, seen map[string]struct{}) []string {
	errors := make([]string, 0)
	slugSeen := map[string]struct{}{}
	for _, item := range items {
		code := strings.TrimSpace(item.Code)
		slug := strings.TrimSpace(item.Slug)
		if code == "" {
			errors = append(errors, "theme requires code")
			continue
		}
		if _, ok := seen[code]; ok {
			errors = append(errors, fmt.Sprintf("duplicate theme code %s", code))
		}
		seen[code] = struct{}{}
		if slug == "" {
			errors = append(errors, fmt.Sprintf("theme %s requires slug", code))
			continue
		}
		if _, ok := slugSeen[slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate theme slug %s", slug))
		}
		slugSeen[slug] = struct{}{}
	}
	return errors
}

func validateUniqueCreators(items []dto.GeneratedCreator, seen map[string]struct{}) []string {
	errors := make([]string, 0)
	for _, item := range items {
		slug := strings.TrimSpace(item.Slug)
		if slug == "" {
			errors = append(errors, "creator requires slug")
			continue
		}
		if _, ok := seen[slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate creator slug %s", slug))
		}
		seen[slug] = struct{}{}
	}
	return errors
}

func validateUniqueWorks(items []dto.GeneratedWork, seen map[string]struct{}) []string {
	errors := make([]string, 0)
	for _, item := range items {
		slug := strings.TrimSpace(item.Slug)
		if slug == "" {
			errors = append(errors, "work requires slug")
			continue
		}
		if _, ok := seen[slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate work slug %s", slug))
		}
		seen[slug] = struct{}{}
	}
	return errors
}

func validateUniqueCharacters(items []dto.GeneratedCharacter, seen map[string]struct{}) []string {
	errors := make([]string, 0)
	for _, item := range items {
		slug := strings.TrimSpace(item.Slug)
		if slug == "" {
			errors = append(errors, "character requires slug")
			continue
		}
		if _, ok := seen[slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate character slug %s", slug))
		}
		seen[slug] = struct{}{}
	}
	return errors
}

func validateUniqueSongs(items []dto.GeneratedSong) []string {
	errors := make([]string, 0)
	seen := map[string]struct{}{}
	for _, item := range items {
		slug := strings.TrimSpace(item.Slug)
		if slug == "" {
			errors = append(errors, "song requires slug")
			continue
		}
		if _, ok := seen[slug]; ok {
			errors = append(errors, fmt.Sprintf("duplicate song slug %s", slug))
		}
		seen[slug] = struct{}{}
	}
	return errors
}

func validateOptionalCode(label, code string, existing map[string]struct{}) []string {
	if strings.TrimSpace(code) == "" {
		return nil
	}
	if _, ok := existing[code]; ok {
		return nil
	}
	return []string{fmt.Sprintf("%s references missing code %s", label, code)}
}

func validateAutocreateRequiredCode(label, code string, existing map[string]struct{}, table string) []string {
	if strings.TrimSpace(code) == "" {
		return []string{fmt.Sprintf("%s is required", label)}
	}
	if _, ok := existing[code]; ok {
		return nil
	}
	return []string{fmt.Sprintf("%s references missing code %s and will be auto-created in %s", label, code, table)}
}

func validateOptionalAutocreateCode(label, code string, existing map[string]struct{}, table string) []string {
	if strings.TrimSpace(code) == "" {
		return nil
	}
	if _, ok := existing[code]; ok {
		return nil
	}
	return []string{fmt.Sprintf("%s references missing code %s and will be auto-created in %s", label, code, table)}
}

func validateStatus(label, status string) []string {
	status = normalizeAdminStatus(status, "draft")
	if status != "draft" && status != "published" && status != "archived" {
		return []string{fmt.Sprintf("%s has invalid status %s", label, status)}
	}
	return nil
}

func hasCode(existing map[string]struct{}, code string) bool {
	_, ok := existing[code]
	return ok
}

func hasSlug(packageSet, dbSet map[string]struct{}, slug string) bool {
	if _, ok := packageSet[slug]; ok {
		return true
	}
	_, ok := dbSet[slug]
	return ok
}

func hasThemeCode(packageSet, dbSet map[string]struct{}, code string) bool {
	if _, ok := packageSet[code]; ok {
		return true
	}
	_, ok := dbSet[code]
	return ok
}

func fetchCodeSet(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, table string) (map[string]struct{}, error) {
	rows, err := q.Query(ctx, fmt.Sprintf(`SELECT code FROM %s`, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]struct{}{}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		result[code] = struct{}{}
	}
	return result, rows.Err()
}

func fetchSlugSet(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, table string) (map[string]struct{}, error) {
	rows, err := q.Query(ctx, fmt.Sprintf(`SELECT slug FROM %s`, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]struct{}{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		result[slug] = struct{}{}
	}
	return result, rows.Err()
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func ensureGeneratedPackageDictionaries(ctx context.Context, tx pgx.Tx, pkg dto.GeneratedPackage) ([]string, error) {
	workTypeCodes, err := fetchCodeSet(ctx, tx, "public.pm_work_types")
	if err != nil {
		return nil, err
	}
	regionCodes, err := fetchCodeSet(ctx, tx, "public.pm_regions")
	if err != nil {
		return nil, err
	}
	culturalCodes, err := fetchCodeSet(ctx, tx, "public.pm_cultural_regions")
	if err != nil {
		return nil, err
	}

	warnings := make([]string, 0)
	seenWorkTypes := map[string]struct{}{}
	for _, item := range pkg.PmWorks {
		code := strings.TrimSpace(item.WorkTypeCode)
		if code == "" {
			continue
		}
		if _, ok := workTypeCodes[code]; ok {
			continue
		}
		if _, ok := seenWorkTypes[code]; ok {
			continue
		}
		if err := upsertAutoWorkType(ctx, tx, code); err != nil {
			return nil, err
		}
		seenWorkTypes[code] = struct{}{}
		workTypeCodes[code] = struct{}{}
		warnings = append(warnings, fmt.Sprintf("auto-created pm_work_types code %s", code))
	}

	regionRefs := collectGeneratedRegionRefs(pkg)
	seenRegions := map[string]struct{}{}
	for _, code := range regionRefs {
		if _, ok := regionCodes[code]; ok {
			continue
		}
		if _, ok := seenRegions[code]; ok {
			continue
		}
		if err := upsertAutoRegion(ctx, tx, code); err != nil {
			return nil, err
		}
		seenRegions[code] = struct{}{}
		regionCodes[code] = struct{}{}
		warnings = append(warnings, fmt.Sprintf("auto-created pm_regions code %s", code))
	}

	culturalRefs := collectGeneratedCulturalRegionRefs(pkg)
	seenCultural := map[string]struct{}{}
	for _, code := range culturalRefs {
		if _, ok := culturalCodes[code]; ok {
			continue
		}
		if _, ok := seenCultural[code]; ok {
			continue
		}
		if err := upsertAutoCulturalRegion(ctx, tx, code); err != nil {
			return nil, err
		}
		seenCultural[code] = struct{}{}
		culturalCodes[code] = struct{}{}
		warnings = append(warnings, fmt.Sprintf("auto-created pm_cultural_regions code %s", code))
	}

	return warnings, nil
}

func collectGeneratedRegionRefs(pkg dto.GeneratedPackage) []string {
	values := make([]string, 0, len(pkg.PmCreators)+len(pkg.PmWorks)+len(pkg.PmCharacters))
	for _, item := range pkg.PmCreators {
		if code := strings.TrimSpace(item.RegionCode); code != "" {
			values = append(values, code)
		}
	}
	for _, item := range pkg.PmWorks {
		if code := strings.TrimSpace(item.RegionCode); code != "" {
			values = append(values, code)
		}
	}
	for _, item := range pkg.PmCharacters {
		if code := strings.TrimSpace(item.RegionCode); code != "" {
			values = append(values, code)
		}
	}
	return uniq(values)
}

func collectGeneratedCulturalRegionRefs(pkg dto.GeneratedPackage) []string {
	values := make([]string, 0, len(pkg.PmCreators)+len(pkg.PmWorks)+len(pkg.PmCharacters))
	for _, item := range pkg.PmCreators {
		if code := strings.TrimSpace(item.CulturalRegionCode); code != "" {
			values = append(values, code)
		}
	}
	for _, item := range pkg.PmWorks {
		if code := strings.TrimSpace(item.CulturalRegionCode); code != "" {
			values = append(values, code)
		}
	}
	for _, item := range pkg.PmCharacters {
		if code := strings.TrimSpace(item.CulturalRegionCode); code != "" {
			values = append(values, code)
		}
	}
	return uniq(values)
}

func generatedDictName(code string) string {
	text := strings.TrimSpace(code)
	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, "-", " ")
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return code
	}
	return text
}

func upsertAutoWorkType(ctx context.Context, tx pgx.Tx, code string) error {
	name := generatedDictName(code)
	_, err := tx.Exec(ctx, `
INSERT INTO public.pm_work_types (code, name_zh, description, sort_order, is_active)
VALUES ($1,$2,$3,9990,TRUE)
ON CONFLICT (code) DO UPDATE SET
  name_zh = COALESCE(NULLIF(public.pm_work_types.name_zh, ''), EXCLUDED.name_zh),
  description = COALESCE(public.pm_work_types.description, EXCLUDED.description),
  updated_at = NOW()
`, code, name, "auto-created by admin import")
	return err
}

func upsertAutoRegion(ctx context.Context, tx pgx.Tx, code string) error {
	name := generatedDictName(code)
	regionType := "region"
	if strings.Contains(strings.ToLower(code), "asia") || strings.Contains(strings.ToLower(code), "europe") || strings.Contains(strings.ToLower(code), "america") || strings.Contains(strings.ToLower(code), "africa") || strings.Contains(strings.ToLower(code), "oceania") {
		regionType = "continent"
	}
	_, err := tx.Exec(ctx, `
INSERT INTO public.pm_regions (code, name_zh, region_type, parent_id, description, cover_url, sort_order, is_active)
VALUES ($1,$2,$3,NULL,$4,NULL,9990,TRUE)
ON CONFLICT (code) DO UPDATE SET
  name_zh = COALESCE(NULLIF(public.pm_regions.name_zh, ''), EXCLUDED.name_zh),
  region_type = COALESCE(NULLIF(public.pm_regions.region_type, ''), EXCLUDED.region_type),
  description = COALESCE(public.pm_regions.description, EXCLUDED.description),
  updated_at = NOW()
`, code, name, regionType, "auto-created by admin import")
	return err
}

func upsertAutoCulturalRegion(ctx context.Context, tx pgx.Tx, code string) error {
	name := generatedDictName(code)
	_, err := tx.Exec(ctx, `
INSERT INTO public.pm_cultural_regions (code, name_zh, parent_id, description, cover_url, sort_order, is_active)
VALUES ($1,$2,NULL,$3,NULL,9990,TRUE)
ON CONFLICT (code) DO UPDATE SET
  name_zh = COALESCE(NULLIF(public.pm_cultural_regions.name_zh, ''), EXCLUDED.name_zh),
  description = COALESCE(public.pm_cultural_regions.description, EXCLUDED.description),
  updated_at = NOW()
`, code, name, "auto-created by admin import")
	return err
}

func upsertGeneratedTheme(ctx context.Context, tx pgx.Tx, item dto.GeneratedTheme) error {
	_, err := tx.Exec(ctx, `
INSERT INTO public.pm_themes (
  code, slug, name_zh, summary, category, description, cover_url, sort_order, is_active
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9
)
ON CONFLICT (code) DO UPDATE SET
  slug = EXCLUDED.slug,
  name_zh = EXCLUDED.name_zh,
  summary = EXCLUDED.summary,
  category = EXCLUDED.category,
  description = EXCLUDED.description,
  cover_url = EXCLUDED.cover_url,
  sort_order = EXCLUDED.sort_order,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, item.Code, item.Slug, item.NameZH, item.Summary, item.Category, item.Description, item.CoverURL, item.SortOrder, item.IsActive)
	return err
}

func upsertGeneratedCreator(ctx context.Context, tx pgx.Tx, item dto.GeneratedCreator) error {
	creatorTypeID, err := lookupIDByCode(ctx, tx, "public.pm_creator_types", item.CreatorTypeCode)
	if err != nil {
		return err
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", item.RegionCode)
	if err != nil {
		return err
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", item.CulturalRegionCode)
	if err != nil {
		return err
	}
	metaJSON, _ := json.Marshal(item.Meta)
	_, err = tx.Exec(ctx, `
INSERT INTO public.pm_creators (
  creator_type_id, name, slug, aliases, region_id, cultural_region_id, era_text,
  summary, introduction, cover_url, major_works, identity_tags, meta, sort_order, is_active
) VALUES (
  $1,$2,$3,COALESCE($4, ARRAY[]::text[]),$5,$6,$7,$8,$9,$10,COALESCE($11, ARRAY[]::text[]),COALESCE($12, ARRAY[]::text[]),COALESCE($13::jsonb,'{}'::jsonb),$14,$15
)
ON CONFLICT (slug) DO UPDATE SET
  creator_type_id = EXCLUDED.creator_type_id,
  name = EXCLUDED.name,
  aliases = EXCLUDED.aliases,
  region_id = EXCLUDED.region_id,
  cultural_region_id = EXCLUDED.cultural_region_id,
  era_text = EXCLUDED.era_text,
  summary = EXCLUDED.summary,
  introduction = EXCLUDED.introduction,
  cover_url = EXCLUDED.cover_url,
  major_works = EXCLUDED.major_works,
  identity_tags = EXCLUDED.identity_tags,
  meta = EXCLUDED.meta,
  sort_order = EXCLUDED.sort_order,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, creatorTypeID, item.Name, item.Slug, item.Aliases, regionID, culturalID, item.EraText, item.Summary, item.Introduction, item.CoverURL, item.MajorWorks, item.IdentityTags, string(metaJSON), item.SortOrder, item.IsActive)
	return err
}

func upsertGeneratedWork(ctx context.Context, tx pgx.Tx, item dto.GeneratedWork) error {
	workTypeID, err := lookupIDByCode(ctx, tx, "public.pm_work_types", item.WorkTypeCode)
	if err != nil {
		return err
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", item.RegionCode)
	if err != nil {
		return err
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", item.CulturalRegionCode)
	if err != nil {
		return err
	}
	metaJSON, _ := json.Marshal(item.Meta)
	_, err = tx.Exec(ctx, `
INSERT INTO public.pm_works (
  work_type_id, title, slug, subtitle, original_title, aliases, region_id, cultural_region_id,
  era_text, release_year, summary, introduction, cover_url, themes, version_notes, meta, sort_order, is_active
) VALUES (
  $1,$2,$3,NULLIF($4,''),NULLIF($5,''),COALESCE($6, ARRAY[]::text[]),$7,$8,$9,$10,$11,$12,$13,COALESCE($14, ARRAY[]::text[]),COALESCE($15, ARRAY[]::text[]),COALESCE($16::jsonb,'{}'::jsonb),$17,$18
)
ON CONFLICT (slug) DO UPDATE SET
  work_type_id = EXCLUDED.work_type_id,
  title = EXCLUDED.title,
  subtitle = EXCLUDED.subtitle,
  original_title = EXCLUDED.original_title,
  aliases = EXCLUDED.aliases,
  region_id = EXCLUDED.region_id,
  cultural_region_id = EXCLUDED.cultural_region_id,
  era_text = EXCLUDED.era_text,
  release_year = EXCLUDED.release_year,
  summary = EXCLUDED.summary,
  introduction = EXCLUDED.introduction,
  cover_url = EXCLUDED.cover_url,
  themes = EXCLUDED.themes,
  version_notes = EXCLUDED.version_notes,
  meta = EXCLUDED.meta,
  sort_order = EXCLUDED.sort_order,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, workTypeID, item.Title, item.Slug, item.Subtitle, item.OriginalTitle, item.Aliases, regionID, culturalID, item.EraText, item.ReleaseYear, item.Summary, item.Introduction, item.CoverURL, item.Themes, item.VersionNotes, string(metaJSON), item.SortOrder, item.IsActive)
	return err
}

func upsertGeneratedCharacter(ctx context.Context, tx pgx.Tx, item dto.GeneratedCharacter) error {
	charTypeID, err := lookupIDByCode(ctx, tx, "public.pm_character_types", item.CharacterTypeCode)
	if err != nil {
		return err
	}
	regionID, err := optionalLookupIDByCode(ctx, tx, "public.pm_regions", item.RegionCode)
	if err != nil {
		return err
	}
	culturalID, err := optionalLookupIDByCode(ctx, tx, "public.pm_cultural_regions", item.CulturalRegionCode)
	if err != nil {
		return err
	}
	meta := item.Meta
	if meta == nil {
		meta = map[string]any{}
	}
	if strings.TrimSpace(item.CoverPrompt) != "" {
		if _, ok := meta["cover_prompt"]; !ok {
			meta["cover_prompt"] = item.CoverPrompt
		}
	}
	colorsJSON, _ := json.Marshal(item.Colors)
	relationshipJSON, _ := json.Marshal(item.RelationshipProfile)
	psychologyJSON, _ := json.Marshal(item.Psychology)
	timelineJSON, _ := json.Marshal(item.Timeline)
	metaJSON, _ := json.Marshal(meta)
	status := normalizeAdminStatus(item.Status, "published")
	isActive := item.IsActive && status == "published"

	_, err = tx.Exec(ctx, `
INSERT INTO public.pm_characters (
  character_type_id, name, slug, aliases, subtype, gender, region_id, cultural_region_id,
  era_text, dynasty_period_text, summary, cover_url, one_line_definition, core_identity,
  public_image, hidden_self, core_fear, psychological_wound, core_conflict,
  emotional_tone, emotional_temperature, origin, fate_arc, ending_state,
  mbti_guess, mbti_confidence, cognitive_style, surface_traits, deep_traits,
  behavior_patterns, stress_response, dominant_emotions, suppressed_emotions,
  values_tags, bottom_lines, taboos, symbolic_images, colors, elements,
  soundscape_keywords, relationship_profile, psychology, timeline, meta, sort_order, status, is_active
) VALUES (
  $1,$2,$3,COALESCE($4, ARRAY[]::text[]),NULLIF($5,''),NULLIF($6,''),$7,$8,$9,NULLIF($10,''),$11,$12,$13,$14,
  $15,$16,$17,$18,$19,$20,NULLIF($21,''),$22,$23,$24,COALESCE($25, ARRAY[]::text[]),NULLIF($26,''),COALESCE($27, ARRAY[]::text[]),COALESCE($28, ARRAY[]::text[]),COALESCE($29, ARRAY[]::text[]),
  COALESCE($30, ARRAY[]::text[]),COALESCE($31, ARRAY[]::text[]),COALESCE($32, ARRAY[]::text[]),COALESCE($33, ARRAY[]::text[]),COALESCE($34, ARRAY[]::text[]),COALESCE($35, ARRAY[]::text[]),COALESCE($36, ARRAY[]::text[]),COALESCE($37, ARRAY[]::text[]),COALESCE($38::jsonb,'[]'::jsonb),COALESCE($39, ARRAY[]::text[]),
  COALESCE($40, ARRAY[]::text[]),COALESCE($41::jsonb,'{}'::jsonb),COALESCE($42::jsonb,'{}'::jsonb),COALESCE($43::jsonb,'[]'::jsonb),COALESCE($44::jsonb,'{}'::jsonb),$45,$46,$47
)
ON CONFLICT (slug) DO UPDATE SET
  character_type_id = EXCLUDED.character_type_id,
  name = EXCLUDED.name,
  aliases = EXCLUDED.aliases,
  subtype = EXCLUDED.subtype,
  gender = EXCLUDED.gender,
  region_id = EXCLUDED.region_id,
  cultural_region_id = EXCLUDED.cultural_region_id,
  era_text = EXCLUDED.era_text,
  dynasty_period_text = EXCLUDED.dynasty_period_text,
  summary = EXCLUDED.summary,
  cover_url = EXCLUDED.cover_url,
  one_line_definition = EXCLUDED.one_line_definition,
  core_identity = EXCLUDED.core_identity,
  public_image = EXCLUDED.public_image,
  hidden_self = EXCLUDED.hidden_self,
  core_fear = EXCLUDED.core_fear,
  psychological_wound = EXCLUDED.psychological_wound,
  core_conflict = EXCLUDED.core_conflict,
  emotional_tone = EXCLUDED.emotional_tone,
  emotional_temperature = EXCLUDED.emotional_temperature,
  origin = EXCLUDED.origin,
  fate_arc = EXCLUDED.fate_arc,
  ending_state = EXCLUDED.ending_state,
  mbti_guess = EXCLUDED.mbti_guess,
  mbti_confidence = EXCLUDED.mbti_confidence,
  cognitive_style = EXCLUDED.cognitive_style,
  surface_traits = EXCLUDED.surface_traits,
  deep_traits = EXCLUDED.deep_traits,
  behavior_patterns = EXCLUDED.behavior_patterns,
  stress_response = EXCLUDED.stress_response,
  dominant_emotions = EXCLUDED.dominant_emotions,
  suppressed_emotions = EXCLUDED.suppressed_emotions,
  values_tags = EXCLUDED.values_tags,
  bottom_lines = EXCLUDED.bottom_lines,
  taboos = EXCLUDED.taboos,
  symbolic_images = EXCLUDED.symbolic_images,
  colors = EXCLUDED.colors,
  elements = EXCLUDED.elements,
  soundscape_keywords = EXCLUDED.soundscape_keywords,
  relationship_profile = EXCLUDED.relationship_profile,
  psychology = EXCLUDED.psychology,
  timeline = EXCLUDED.timeline,
  meta = EXCLUDED.meta,
  sort_order = EXCLUDED.sort_order,
  status = EXCLUDED.status,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, charTypeID, item.Name, item.Slug, item.Aliases, item.Subtype, item.Gender, regionID, culturalID, item.EraText, item.DynastyPeriodText, item.Summary, item.CoverURL, item.OneLineDefinition, item.CoreIdentity, item.PublicImage, item.HiddenSelf, item.CoreFear, item.PsychologicalWound, item.CoreConflict, item.EmotionalTone, item.EmotionalTemperature, item.Origin, item.FateArc, item.EndingState, item.MbtiGuess, item.MbtiConfidence, item.CognitiveStyle, item.SurfaceTraits, item.DeepTraits, item.BehaviorPatterns, item.StressResponse, item.DominantEmotions, item.SuppressedEmotions, item.ValuesTags, item.BottomLines, item.Taboos, item.SymbolicImages, string(colorsJSON), item.Elements, item.SoundscapeKeywords, string(relationshipJSON), string(psychologyJSON), string(timelineJSON), string(metaJSON), item.SortOrder, status, isActive)
	return err
}

func upsertGeneratedSong(ctx context.Context, tx pgx.Tx, item dto.GeneratedSong) error {
	characterID, err := resolveIDByRef(ctx, tx, "public.pm_characters", item.CharacterSlug)
	if err != nil {
		return err
	}
	metaJSON, _ := json.Marshal(item.Meta)
	status := normalizeAdminStatus(item.Status, "published")
	isActive := item.IsActive && status == "published"
	versionNo := item.VersionNo
	if versionNo <= 0 {
		versionNo = 1
	}
	_, err = tx.Exec(ctx, `
INSERT INTO public.pm_songs (
  character_id, title, slug, subtitle, summary, cover_url, audio_url,
  song_core_theme, song_summary, song_emotional_curve, song_styles, tempo_bpm,
  vocal_profile, lyric_keywords, forbidden_cliches, symbol_images, ending_feeling,
  prompt, lyrics, version_no, status, meta, sort_order, is_active
) VALUES (
  $1,$2,$3,NULLIF($4,''),$5,$6,$7,$8,$9,COALESCE($10, ARRAY[]::text[]),COALESCE($11, ARRAY[]::text[]),$12,
  $13,COALESCE($14, ARRAY[]::text[]),COALESCE($15, ARRAY[]::text[]),COALESCE($16, ARRAY[]::text[]),NULLIF($17,''),$18,$19,$20,$21,COALESCE($22::jsonb,'{}'::jsonb),$23,$24
)
ON CONFLICT (slug) DO UPDATE SET
  character_id = EXCLUDED.character_id,
  title = EXCLUDED.title,
  subtitle = EXCLUDED.subtitle,
  summary = EXCLUDED.summary,
  cover_url = EXCLUDED.cover_url,
  audio_url = EXCLUDED.audio_url,
  song_core_theme = EXCLUDED.song_core_theme,
  song_summary = EXCLUDED.song_summary,
  song_emotional_curve = EXCLUDED.song_emotional_curve,
  song_styles = EXCLUDED.song_styles,
  tempo_bpm = EXCLUDED.tempo_bpm,
  vocal_profile = EXCLUDED.vocal_profile,
  lyric_keywords = EXCLUDED.lyric_keywords,
  forbidden_cliches = EXCLUDED.forbidden_cliches,
  symbol_images = EXCLUDED.symbol_images,
  ending_feeling = EXCLUDED.ending_feeling,
  prompt = EXCLUDED.prompt,
  lyrics = EXCLUDED.lyrics,
  version_no = EXCLUDED.version_no,
  status = EXCLUDED.status,
  meta = EXCLUDED.meta,
  sort_order = EXCLUDED.sort_order,
  is_active = EXCLUDED.is_active,
  updated_at = NOW()
`, characterID, item.Title, item.Slug, item.Subtitle, item.Summary, item.CoverURL, item.AudioURL, item.SongCoreTheme, item.SongSummary, item.SongEmotionalCurve, item.SongStyles, item.TempoBPM, item.VocalProfile, item.LyricKeywords, item.ForbiddenCliches, item.SymbolImages, item.EndingFeeling, item.Prompt, item.Lyrics, versionNo, status, string(metaJSON), item.SortOrder, isActive)
	return err
}

func replaceGeneratedWorkCreators(ctx context.Context, tx pgx.Tx, item dto.GeneratedWork) error {
	workID, err := resolveIDByRef(ctx, tx, "public.pm_works", item.Slug)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_work_creators WHERE work_id=$1`, workID); err != nil {
		return err
	}
	for _, role := range item.CreatorRoles {
		creatorID, err := resolveIDByRef(ctx, tx, "public.pm_creators", role.CreatorSlug)
		if err != nil {
			return err
		}
		roleCode := strings.TrimSpace(role.RoleCode)
		if roleCode == "" {
			roleCode = "author"
		}
		sortOrder := role.SortOrder
		if sortOrder <= 0 {
			sortOrder = 10
		}
		_, err = tx.Exec(ctx, `
INSERT INTO public.pm_work_creators (work_id, creator_id, role_code, is_primary, sort_order, note)
VALUES ($1,$2,$3,$4,$5,NULLIF($6,''))
`, workID, creatorID, roleCode, role.IsPrimary, sortOrder, role.Note)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedCharacterWorks(ctx context.Context, tx pgx.Tx, item dto.GeneratedCharacter) error {
	characterID, err := resolveIDByRef(ctx, tx, "public.pm_characters", item.Slug)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_works WHERE character_id=$1`, characterID); err != nil {
		return err
	}
	for index, slug := range uniq(item.WorkSlugs) {
		workID, err := resolveIDByRef(ctx, tx, "public.pm_works", slug)
		if err != nil {
			return err
		}
		sortOrder := (index + 1) * 10
		_, err = tx.Exec(ctx, `
INSERT INTO public.pm_character_works (character_id, work_id, relation_type, is_primary, sort_order, note)
VALUES ($1,$2,'belongs_to',$3,$4,NULL)
`, characterID, workID, slug == item.PrimaryWork, sortOrder)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedCharacterMotivations(ctx context.Context, tx pgx.Tx, item dto.GeneratedCharacter) error {
	characterID, err := resolveIDByRef(ctx, tx, "public.pm_characters", item.Slug)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_motivations WHERE character_id=$1`, characterID); err != nil {
		return err
	}
	for index, code := range uniq(item.MotivationCodes) {
		motivationID, err := lookupIDByCode(ctx, tx, "public.pm_motivation_dict", code)
		if err != nil {
			return err
		}
		weight := 0.9 - float64(index)*0.1
		if code == item.PrimaryMotivation {
			weight = 1.0
		}
		if weight < 0.5 {
			weight = 0.5
		}
		_, err = tx.Exec(ctx, `
INSERT INTO public.pm_character_motivations (character_id, motivation_id, weight, is_primary, note)
VALUES ($1,$2,$3,$4,NULL)
`, characterID, motivationID, weight, code == item.PrimaryMotivation)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceGeneratedCharacterThemes(ctx context.Context, tx pgx.Tx, item dto.GeneratedCharacter) error {
	characterID, err := resolveIDByRef(ctx, tx, "public.pm_characters", item.Slug)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM public.pm_character_themes WHERE character_id=$1`, characterID); err != nil {
		return err
	}
	for index, code := range uniq(item.ThemeCodes) {
		themeID, err := lookupIDByCodeAny(ctx, tx, "public.pm_themes", code)
		if err != nil {
			return err
		}
		weight := 0.9 - float64(index)*0.1
		if code == item.PrimaryTheme {
			weight = 1.0
		}
		if weight < 0.5 {
			weight = 0.5
		}
		_, err = tx.Exec(ctx, `
INSERT INTO public.pm_character_themes (character_id, theme_id, weight, is_primary, note)
VALUES ($1,$2,$3,$4,NULL)
`, characterID, themeID, weight, code == item.PrimaryTheme)
		if err != nil {
			return err
		}
	}
	return nil
}

func lookupIDByCodeAny(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, table, code string) (string, error) {
	var id string
	err := q.QueryRow(ctx, fmt.Sprintf(`SELECT id::text FROM %s WHERE code=$1 LIMIT 1`, table), code).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
