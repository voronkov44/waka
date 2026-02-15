package modelsutil_test

import (
	"errors"
	"reflect"
	"testing"

	"gorm.io/datatypes"

	"rest_waka/pkg/modelsutil"
)

func TestMarshalFlavors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []string
		wantRaw string
	}{
		{
			name:    "nil slice marshals to empty array",
			input:   nil,
			wantRaw: "[]",
		},
		{
			name:    "non-empty slice marshals as-is",
			input:   []string{"mint", "cola"},
			wantRaw: "[\"mint\",\"cola\"]",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			raw, err := modelsutil.MarshalFlavors(tc.input)
			if err != nil {
				t.Fatalf("MarshalFlavors() error = %v", err)
			}

			if string(raw) != tc.wantRaw {
				t.Fatalf("MarshalFlavors() raw = %q, want %q", string(raw), tc.wantRaw)
			}
		})
	}
}

func TestUnmarshalFlavors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     datatypes.JSON
		want      []string
		wantError bool
	}{
		{
			name:  "empty json value returns empty slice",
			input: nil,
			want:  []string{},
		},
		{
			name:  "json null returns empty slice",
			input: datatypes.JSON([]byte("null")),
			want:  []string{},
		},
		{
			name:  "valid json array",
			input: datatypes.JSON([]byte("[\"mint\",\"cola\"]")),
			want:  []string{"mint", "cola"},
		},
		{
			name:      "invalid json returns error",
			input:     datatypes.JSON([]byte("{")),
			wantError: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := modelsutil.UnmarshalFlavors(tc.input)
			if tc.wantError {
				if err == nil {
					t.Fatal("UnmarshalFlavors() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalFlavors() error = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("UnmarshalFlavors() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestAddFlavorUnique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		base        []string
		value       string
		want        []string
		wantChanged bool
		wantErr     error
	}{
		{
			name:        "adds trimmed flavor",
			base:        []string{"mint"},
			value:       "  cola ",
			want:        []string{"mint", "cola"},
			wantChanged: true,
		},
		{
			name:        "duplicate exact is idempotent",
			base:        []string{"mint"},
			value:       "mint",
			want:        []string{"mint"},
			wantChanged: false,
		},
		{
			name:        "duplicate case-insensitive is idempotent",
			base:        []string{"Mint"},
			value:       "mInT",
			want:        []string{"Mint"},
			wantChanged: false,
		},
		{
			name:        "duplicate with spaces in existing value is idempotent",
			base:        []string{"  Mint  "},
			value:       "mint",
			want:        []string{"  Mint  "},
			wantChanged: false,
		},
		{
			name:        "empty flavor is rejected",
			base:        []string{"mint"},
			value:       "   ",
			want:        []string{"mint"},
			wantChanged: false,
			wantErr:     modelsutil.ErrEmptyFlavor,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, changed, err := modelsutil.AddFlavorUnique(tc.base, tc.value)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("AddFlavorUnique() err = %v, want %v", err, tc.wantErr)
			}
			if changed != tc.wantChanged {
				t.Fatalf("AddFlavorUnique() changed = %v, want %v", changed, tc.wantChanged)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("AddFlavorUnique() flavors = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestRemoveFlavor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		base        []string
		value       string
		want        []string
		wantChanged bool
		wantErr     error
	}{
		{
			name:        "removes existing flavor case-insensitive",
			base:        []string{"Mint", "Cola"},
			value:       "mint",
			want:        []string{"Cola"},
			wantChanged: true,
		},
		{
			name:        "removes existing flavor with spaces in stored value",
			base:        []string{"  Mint  ", "Cola"},
			value:       "mint",
			want:        []string{"Cola"},
			wantChanged: true,
		},
		{
			name:        "removes all duplicates",
			base:        []string{"mint", "MINT", "cola"},
			value:       "mint",
			want:        []string{"cola"},
			wantChanged: true,
		},
		{
			name:        "remove missing flavor is idempotent",
			base:        []string{"mint"},
			value:       "cola",
			want:        []string{"mint"},
			wantChanged: false,
		},
		{
			name:        "empty flavor is rejected",
			base:        []string{"mint"},
			value:       "",
			want:        []string{"mint"},
			wantChanged: false,
			wantErr:     modelsutil.ErrEmptyFlavor,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, changed, err := modelsutil.RemoveFlavor(tc.base, tc.value)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("RemoveFlavor() err = %v, want %v", err, tc.wantErr)
			}
			if changed != tc.wantChanged {
				t.Fatalf("RemoveFlavor() changed = %v, want %v", changed, tc.wantChanged)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("RemoveFlavor() flavors = %#v, want %#v", got, tc.want)
			}
		})
	}
}
