package model

import (
	"testing"
)

func TestValidateNil(t *testing.T) {
	tests := []struct {
		Model   interface{}
		WantErr error
	}{
		{Model: User{}, WantErr: nil},
		{Model: &User{}, WantErr: nil},
		{Model: nil, WantErr: ErrNilPointer},
	}

	for _, tt := range tests {
		gotErr := ValidateStructNil(tt.Model)
		if gotErr != tt.WantErr {
			t.Fatalf("want: %v, got: %v", tt.WantErr, gotErr)
		}
	}
}
