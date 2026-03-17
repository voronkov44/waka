package models

import (
	"encoding/json"
	"strings"

	"gorm.io/datatypes"
	"rest_waka/pkg/tagutil"
)

func normalizeTag(in *ModelTag) (*ModelTag, error) {
	if in == nil {
		return nil, nil
	}

	key := strings.ToLower(strings.TrimSpace(in.Key))
	if key == "" {
		key = "custom"
	}

	visual, ok := tagutil.NormalizeVisualTag(tagutil.VisualTag{
		Label:     in.Label,
		BgColor:   in.BgColor,
		TextColor: in.TextColor,
		Outlined:  in.Outlined,
	})
	if !ok {
		return nil, ErrInvalidArgument
	}

	return &ModelTag{
		Key:       key,
		Label:     visual.Label,
		BgColor:   visual.BgColor,
		TextColor: visual.TextColor,
		Outlined:  visual.Outlined,
	}, nil
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
