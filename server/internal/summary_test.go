package internal

import "testing"

func TestGetSummary(t *testing.T) {
	puid := "fmt/123"
	mime := "application/pdf"
	version := "1.7"
	score := 0.8

	tests := []struct {
		name     string
		sets     []FeatureSet
		results  []ToolResult
		validate func(t *testing.T, s Summary)
	}{
		{
			name: "full valid summary",
			sets: []FeatureSet{
				{
					Score: 0.8,
					Features: map[string]MergeFeatureValue{
						"format:valid":    {Value: true},
						"format:puid":     {Value: puid},
						"format:mimeType": {Value: mime},
						"format:version":  {Value: version},
					},
				},
			},
			results: []ToolResult{{Score: &score}},
			validate: func(t *testing.T, s Summary) {
				t.Helper()
				if !s.Valid {
					t.Fatal("expected valid=true")
				}
				if s.Invalid {
					t.Fatal("expected invalid=false")
				}
				if s.FormatUncertain {
					t.Fatal("expected formatUncertain=false")
				}
				if s.Error {
					t.Fatal("expected error=false")
				}
				if s.PUID == nil || *s.PUID != puid {
					t.Fatalf("unexpected puid: %#v", s.PUID)
				}
				if s.MimeType == nil || *s.MimeType != mime {
					t.Fatalf("unexpected mime type: %#v", s.MimeType)
				}
				if s.FormatVersion == nil || *s.FormatVersion != version {
					t.Fatalf("unexpected format version: %#v", s.FormatVersion)
				}
			},
		},
		{
			name:    "empty sets are uncertain",
			sets:    nil,
			results: nil,
			validate: func(t *testing.T, s Summary) {
				t.Helper()
				if !s.FormatUncertain {
					t.Fatal("expected formatUncertain=true")
				}
			},
		},
		{
			name: "error in tool results sets error flag",
			sets: []FeatureSet{
				{
					Score: 0.8,
					Features: map[string]MergeFeatureValue{
						"format:valid": {Value: false},
					},
				},
			},
			results: []ToolResult{
				{Error: &[]string{"boom"}[0]},
			},
			validate: func(t *testing.T, s Summary) {
				t.Helper()
				if !s.Error {
					t.Fatal("expected error=true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := GetSummary(tt.sets, tt.results)
			tt.validate(t, summary)
		})
	}
}
