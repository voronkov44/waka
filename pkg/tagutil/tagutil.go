package tagutil

import (
	"regexp"
	"strings"
)

type VisualTag struct {
	Label     string
	BgColor   string
	TextColor string
	Outlined  bool
}

var hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func NormalizeVisualTag(in VisualTag) (VisualTag, bool) {
	out := VisualTag{
		Label:     strings.TrimSpace(in.Label),
		BgColor:   strings.TrimSpace(in.BgColor),
		TextColor: strings.TrimSpace(in.TextColor),
		Outlined:  in.Outlined,
	}

	if out.Label == "" {
		return VisualTag{}, false
	}
	if !IsValidHexColor(out.BgColor) || !IsValidHexColor(out.TextColor) {
		return VisualTag{}, false
	}

	return out, true
}

func IsValidHexColor(s string) bool {
	return hexColorRe.MatchString(strings.TrimSpace(s))
}
