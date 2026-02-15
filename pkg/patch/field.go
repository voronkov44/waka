package patch

import (
	"bytes"
	"encoding/json"
)

type Field[T any] struct {
	Set   bool
	Null  bool
	Value T
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Set = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		f.Null = true
		var zero T
		f.Value = zero
		return nil
	}

	f.Null = false
	return json.Unmarshal(data, &f.Value)
}
