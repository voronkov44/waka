package modelsutil

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"gorm.io/datatypes"
)

var ErrEmptyFlavor = errors.New("empty flavor")

func NormalizeFlavor(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeFlavors - строгая нормализация для write-path:
// - trim
// - пустые значения -> ошибка
// - dedupe case-insensitive
// - сортировка по алфавиту case-insensitive
func NormalizeFlavors(flavors []string) ([]string, error) {
	if flavors == nil {
		return []string{}, nil
	}

	seen := make(map[string]struct{}, len(flavors))
	out := make([]string, 0, len(flavors))

	for _, f := range flavors {
		v := NormalizeFlavor(f)
		if v == "" {
			return nil, ErrEmptyFlavor
		}

		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i]) < strings.ToLower(out[j])
	})

	return out, nil
}

// CleanupFlavors - мягкая очистка для backfill:
// - trim
// - пустые значения пропускаем
// - dedupe case-insensitive
// - сортировка по алфавиту case-insensitive
func CleanupFlavors(flavors []string) []string {
	if flavors == nil {
		return []string{}
	}

	seen := make(map[string]struct{}, len(flavors))
	out := make([]string, 0, len(flavors))

	for _, f := range flavors {
		v := NormalizeFlavor(f)
		if v == "" {
			continue
		}

		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i]) < strings.ToLower(out[j])
	})

	return out
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

// AddFlavorUnique - добавляет вкус, если такого еще нет
func AddFlavorUnique(flavors []string, value string) ([]string, bool, error) {
	val := NormalizeFlavor(value)
	if val == "" {
		return flavors, false, ErrEmptyFlavor
	}

	for _, f := range flavors {
		if strings.EqualFold(NormalizeFlavor(f), val) {
			return flavors, false, nil
		}
	}

	next := append(append([]string{}, flavors...), val)
	next, err := NormalizeFlavors(next)
	if err != nil {
		return flavors, false, err
	}

	return next, true, nil
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

	if !removed {
		return flavors, false, nil
	}

	out, err := NormalizeFlavors(out)
	if err != nil {
		return flavors, false, err
	}

	return out, true, nil
}
