package models

import (
	"encoding/json"
	"regexp"
	"strings"

	"gorm.io/datatypes"
)

var hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func normalizeTag(in *ModelTag) (*ModelTag, error) {
	if in == nil {
		return nil, nil
	}

	tag := &ModelTag{
		Key:       strings.ToLower(strings.TrimSpace(in.Key)),
		Label:     strings.TrimSpace(in.Label),
		BgColor:   strings.TrimSpace(in.BgColor),
		TextColor: strings.TrimSpace(in.TextColor),
	}

	if tag.Key == "" {
		tag.Key = "custom"
	}
	if tag.Label == "" {
		return nil, ErrInvalidArgument
	}
	if !isValidHexColor(tag.BgColor) || !isValidHexColor(tag.TextColor) {
		return nil, ErrInvalidArgument
	}

	return tag, nil
}

func isValidHexColor(s string) bool {
	return hexColorRe.MatchString(strings.TrimSpace(s))
}

func marshalTag(tag *ModelTag) (datatypes.JSON, error) {
	if tag == nil {
		return nil, nil
	}
	raw, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func unmarshalTag(raw datatypes.JSON) (*ModelTag, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var tag ModelTag
	if err := json.Unmarshal(raw, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}
