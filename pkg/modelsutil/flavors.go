package modelsutil

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
	"strings"
)

var ErrEmptyFlavor = errors.New("empty flavor")

func NormalizeFlavor(s string) string {
	return strings.TrimSpace(s)
}

func MarshalFlavors(flavors []string) (datatypes.JSON, error) {
	if flavors == nil {
		flavors = []string{}
	}
	b, err := json.Marshal(flavors)
	if err != nil {
		return nil, err
	}
	raw := datatypes.JSON(b)
	return raw, nil
}

func UnmarshalFlavors(raw datatypes.JSON) ([]string, error) {
	if len(raw) == 0 {
		return []string{}, nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []string{}
	}
	return out, nil
}

// AddFlavorUnique - добавляет вкус, если такого еще нет (сравнение)
// Возвращает новый слайс и флаг chanded
func AddFlavorUnique(flavors []string, value string) ([]string, bool, error) {
	val := NormalizeFlavor(value)
	if val == "" {
		return flavors, false, ErrEmptyFlavor
	}
	for _, f := range flavors {
		if strings.EqualFold(NormalizeFlavor(f), val) {
			return flavors, false, nil // уже есть, идемпотентно
		}
	}
	return append(flavors, val), true, nil
}

func RemoveFlavor(flavors []string, value string) ([]string, bool, error) {
	val := NormalizeFlavor(value)
	if val == "" {
		return flavors, false, ErrEmptyFlavor
	}
	out := make([]string, 0, len(flavors))
	removed := false
	for _, f := range flavors {
		if strings.EqualFold(NormalizeFlavor(f), val) {
			removed = true
			continue
		}
		out = append(out, f)
	}
	return out, removed, nil
}
