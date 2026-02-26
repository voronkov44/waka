package models

import (
	"errors"
	"testing"
)

func TestNormalizeStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "active", input: "active", want: StatusActive},
		{name: "active case and spaces", input: "  AcTiVe ", want: StatusActive},
		{name: "hidden", input: "hidden", want: StatusHidden},
		{name: "archive", input: " ARCHIVE ", want: StatusArchive},
		{name: "empty defaults hidden", input: "   ", want: StatusHidden},
		{name: "reject unknown", input: "draft", wantErr: ErrInvalidArgument},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeStatus(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("normalizeStatus() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if got != tc.want {
				t.Fatalf("normalizeStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}
