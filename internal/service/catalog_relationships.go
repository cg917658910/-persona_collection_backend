package service

import (
	"strings"

	"pm-backend/internal/dto"
)

func (s *CatalogService) ListRelationships(characterSlug string) ([]dto.RelationshipListItemResponse, error) {
	list, err := s.repo.ListRelationships(characterSlug)
	if err != nil {
		return nil, err
	}

	out := make([]dto.RelationshipListItemResponse, 0, len(list))
	for _, item := range list {
		out = append(out, s.buildRelationshipListItem(item, characterSlug))
	}
	return out, nil
}

func (s *CatalogService) GetRelationshipDetail(slug string) (dto.RelationshipDetailResponse, error) {
	item, err := s.repo.GetRelationshipDetail(slug)
	if err != nil {
		return dto.RelationshipDetailResponse{}, err
	}

	item = s.normalizeRelationRecord(item)
	return dto.RelationshipDetailResponse{
		Slug:                    item.Slug,
		Name:                    s.relationshipName(item),
		Subtitle:                item.Subtitle,
		RelationTypeCode:        item.RelationTypeCode,
		RelationTypeName:        s.relationshipDisplayLabel(item),
		SourceCharacterSlug:     item.SourceCharacter.Slug,
		SourceCharacterName:     item.SourceCharacter.Name,
		SourceCharacterCoverURL: item.SourceCharacter.CoverURL,
		TargetCharacterSlug:     item.TargetCharacter.Slug,
		TargetCharacterName:     item.TargetCharacter.Name,
		TargetCharacterCoverURL: item.TargetCharacter.CoverURL,
		WorkSlug:                item.WorkSlug,
		WorkName:                item.WorkName,
		CoreTension:             item.CoreTension,
		EmotionalTone:           item.EmotionalTone,
		OneLineDefinition:       item.OneLineDefinition,
		Summary:                 item.Summary,
		CoverURL:                item.CoverURL,
		ConnectionTrigger:       item.ConnectionTrigger,
		SustainingMechanism:     item.SustainingMechanism,
		RelationConflict:        item.RelationConflict,
		RelationArc:             item.RelationArc,
		FateImpact:              item.FateImpact,
		PowerStructure:          item.PowerStructure,
		DependencyPattern:       item.DependencyPattern,
		SourcePerspective:       item.SourcePerspective,
		SourceDesireInRelation:  item.SourceDesireInRelation,
		SourceFearInRelation:    item.SourceFearInRelation,
		SourceUnsaid:            item.SourceUnsaid,
		TargetPerspective:       item.TargetPerspective,
		TargetDesireInRelation:  item.TargetDesireInRelation,
		TargetFearInRelation:    item.TargetFearInRelation,
		TargetUnsaid:            item.TargetUnsaid,
		Phenomenology:           item.Phenomenology,
		RelationPalette:         s.sliceOrDefaultPalette(item.RelationPalette),
		SymbolicImages:          s.sliceOrDefault(item.SymbolicImages, []string{}),
		RelationKeywords:        s.sliceOrDefault(item.RelationKeywords, []string{}),
		RelationEvents:          s.sliceOrDefaultRelationEvents(item.Events),
		PrimarySong:             item.PrimarySong,
		RelatedRelations:        s.sliceOrDefaultRelationLinks(item.RelatedRelations),
	}, nil
}

func (s *CatalogService) buildKeyRelationships(characterSlug string, limit int) []dto.RelationshipListItemResponse {
	list, err := s.repo.ListRelationships(characterSlug)
	if err != nil {
		return []dto.RelationshipListItemResponse{}
	}

	out := make([]dto.RelationshipListItemResponse, 0, smallerInt(limit, len(list)))
	for _, item := range list {
		out = append(out, s.buildRelationshipListItem(item, characterSlug))
		if len(out) >= limit {
			break
		}
	}
	return out
}

func (s *CatalogService) buildRelationshipListItem(item dto.RelationRecord, anchorCharacterSlug string) dto.RelationshipListItemResponse {
	item = s.normalizeRelationRecord(item)
	label := s.relationshipDisplayLabel(item)
	return dto.RelationshipListItemResponse{
		ID:                item.Slug,
		Slug:              item.Slug,
		Name:              s.relationshipName(item),
		Summary:           item.Summary,
		OneLineDefinition: item.OneLineDefinition,
		CoverURL:          item.CoverURL,
		RelationType:      item.RelationTypeCode,
		RelationLabel:     label,
		WorkTitle:         item.WorkName,
		Tags:              s.relationshipTags(label, item.ThemeTags),
		SourceCharacter:   item.SourceCharacter,
		TargetCharacter:   item.TargetCharacter,
		Counterpart:       s.relationshipCounterpart(item, anchorCharacterSlug),
	}
}

func (s *CatalogService) relationshipCounterpart(item dto.RelationRecord, anchorCharacterSlug string) *dto.RelationshipCharacterRef {
	if strings.TrimSpace(anchorCharacterSlug) == "" {
		return nil
	}
	switch strings.TrimSpace(anchorCharacterSlug) {
	case item.SourceCharacter.Slug:
		return &item.TargetCharacter
	case item.TargetCharacter.Slug:
		return &item.SourceCharacter
	default:
		return nil
	}
}

func (s *CatalogService) relationshipName(item dto.RelationRecord) string {
	if value := strings.TrimSpace(item.Name); value != "" {
		return value
	}
	return strings.TrimSpace(item.SourceCharacter.Name + " × " + item.TargetCharacter.Name)
}

func (s *CatalogService) relationshipDisplayLabel(item dto.RelationRecord) string {
	if value := strings.TrimSpace(item.RelationTypeName); value != "" {
		return value
	}

	switch strings.TrimSpace(item.RelationTypeCode) {
	case "lover":
		return "恋人"
	case "spouse":
		return "伴侣"
	case "rival":
		return "对手"
	case "mirror":
		return "镜像"
	case "friend":
		return "朋友"
	case "mentor_student":
		return "师徒"
	case "parent_child":
		return "亲子"
	case "siblings":
		return "手足"
	case "ruler_subject":
		return "统治-臣属"
	case "oppressor_oppressed":
		return "压迫-被压迫"
	case "co_conspirator":
		return "共谋者"
	case "institution_subject":
		return "制度-主体"
	default:
		return strings.TrimSpace(item.RelationTypeCode)
	}
}

func (s *CatalogService) relationshipTags(label string, themeTags []string) []string {
	tags := make([]string, 0, len(themeTags)+1)
	if value := strings.TrimSpace(label); value != "" {
		tags = append(tags, value)
	}
	tags = append(tags, themeTags...)
	return s.sliceCompact(tags)
}

func (s *CatalogService) normalizeRelationRecord(in dto.RelationRecord) dto.RelationRecord {
	in.CoverURL = s.assets.Normalize(in.CoverURL)
	in.SourceCharacter.CoverURL = s.assets.Normalize(in.SourceCharacter.CoverURL)
	in.TargetCharacter.CoverURL = s.assets.Normalize(in.TargetCharacter.CoverURL)
	if in.PrimarySong != nil {
		in.PrimarySong.CoverURL = s.assets.Normalize(in.PrimarySong.CoverURL)
		in.PrimarySong.AudioURL = s.assets.Normalize(in.PrimarySong.AudioURL)
	}
	for i := range in.RelatedRelations {
		in.RelatedRelations[i].CoverURL = s.assets.Normalize(in.RelatedRelations[i].CoverURL)
	}
	return in
}

func (s *CatalogService) sliceOrDefaultRelationEvents(in []dto.RelationEvent) []dto.RelationEvent {
	if len(in) == 0 {
		return []dto.RelationEvent{}
	}
	return in
}

func (s *CatalogService) sliceOrDefaultRelationLinks(in []dto.RelationLink) []dto.RelationLink {
	if len(in) == 0 {
		return []dto.RelationLink{}
	}
	return in
}

func (s *CatalogService) sliceOrDefaultPalette(in []dto.RelationPaletteItem) []dto.RelationPaletteItem {
	if len(in) == 0 {
		return []dto.RelationPaletteItem{}
	}
	return in
}

func smallerInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
