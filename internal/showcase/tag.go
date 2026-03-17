package showcase

import "rest_waka/pkg/tagutil"

func normalizeTag(in ItemTag) (ItemTag, error) {
	visual, ok := tagutil.NormalizeVisualTag(tagutil.VisualTag{
		Label:     in.Label,
		BgColor:   in.BgColor,
		TextColor: in.TextColor,
		Outlined:  in.Outlined,
	})
	if !ok {
		return ItemTag{}, ErrInvalidArgument
	}

	return ItemTag{
		Label:     visual.Label,
		BgColor:   visual.BgColor,
		TextColor: visual.TextColor,
		Outlined:  visual.Outlined,
	}, nil
}
