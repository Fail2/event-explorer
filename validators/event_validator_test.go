package validators

import (
	"net/url"
	"testing"
)

func TestEventsValidator(t *testing.T) {
	tests := []struct {
		name    string
		query   url.Values
		wantErr string
	}{
		{
			name:    "city missing",
			query:   url.Values{},
			wantErr: "city parameter missing",
		},
		{
			name: "city empty",
			query: url.Values{
				"city": {""},
			},
			wantErr: "city value can't be empty",
		},
		{
			name: "countryCode missing",
			query: url.Values{
				"city": {"London"},
			},
			wantErr: "countryCode parameter missing",
		},
		{
			name: "countryCode empty",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {""},
			},
			wantErr: "countryCode value can't be empty",
		},
		{
			name: "valid query",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EventsValidator(tt.query)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("EventsValidator() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("EventsValidator() error = nil, want %q", tt.wantErr)
			}

			if err.Error() != tt.wantErr {
				t.Errorf("EventsValidator() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCachePurgeValidator(t *testing.T) {
	tests := []struct {
		name    string
		query   url.Values
		wantErr string
	}{
		{
			name:    "all parameters omitted",
			query:   url.Values{},
			wantErr: "",
		},
		{
			name: "city missing",
			query: url.Values{
				"countryCode": {"GB"},
				"category":    {"Music"},
			},
			wantErr: "city, countryCode and category must either all be provided or all be omitted",
		},
		{
			name: "countryCode missing",
			query: url.Values{
				"city":     {"London"},
				"category": {"Music"},
			},
			wantErr: "city, countryCode and category must either all be provided or all be omitted",
		},
		{
			name: "category missing",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
			},
			wantErr: "city, countryCode and category must either all be provided or all be omitted",
		},
		{
			name: "city empty",
			query: url.Values{
				"city":        {""},
				"countryCode": {"GB"},
				"category":    {"Music"},
			},
			wantErr: "city value can't be empty",
		},
		{
			name: "countryCode empty",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {""},
				"category":    {"Music"},
			},
			wantErr: "countryCode value can't be empty",
		},
		{
			name: "category empty",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
				"category":    {""},
			},
			wantErr: "category value can't be empty",
		},
		{
			name: "invalid category",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
				"category":    {"Concert"},
			},
			wantErr: "invalid category: Concert. Allowed values are Music and Sports",
		},
		{
			name: "valid Music category",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
				"category":    {"Music"},
			},
			wantErr: "",
		},
		{
			name: "valid Sports category",
			query: url.Values{
				"city":        {"London"},
				"countryCode": {"GB"},
				"category":    {"Sports"},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CachePurgeValidator(tt.query)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("CachePurgeValidator() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("CachePurgeValidator() error = nil, want %q", tt.wantErr)
			}

			if err.Error() != tt.wantErr {
				t.Errorf("CachePurgeValidator() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}
