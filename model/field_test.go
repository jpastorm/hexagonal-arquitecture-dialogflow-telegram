package model

import "testing"

func TestFields_ValidateNames(t *testing.T) {
	tests := []struct {
		name          string
		fields        Fields
		allowedFields []string
		wantErr       bool
	}{
		{
			name:          "currency_id not allowed",
			fields:        Fields{{Name: "currency_id", Value: 5}, {Name: "campaign_id", Value: 7}},
			allowedFields: []string{"campaign_id"},
			wantErr:       true,
		},
		{
			name:          "currency_id and campaign_id not allowed",
			fields:        Fields{{Name: "currency_id", Value: 5}, {Name: "campaign_id", Value: 7}},
			allowedFields: []string{"id", "is_current"},
			wantErr:       true,
		},
		{
			name:          "allowed fields",
			fields:        Fields{{Name: "currency_id", Value: 5}, {Name: "campaign_id", Value: 7}, {Name: "is_current", Value: true}},
			allowedFields: []string{"campaign_id", "id", "currency_id", "is_current"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.ValidateNames(tt.allowedFields); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNames() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
