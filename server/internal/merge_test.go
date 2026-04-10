package internal

import "testing"

func TestFeatureSetIsEqual(t *testing.T) {
	s1 := FeatureSet{SupportingTools: []string{"tool-b", "tool-a"}}
	s2 := FeatureSet{SupportingTools: []string{"tool-a", "tool-b"}}

	if !s1.IsEqual(s2) {
		t.Fatal("expected sets to be equal")
	}
}

func TestFilterDuplicateSets(t *testing.T) {
	sets := []FeatureSet{
		{
			SupportingTools: []string{"tool-a", "tool-b"},
			Score:           0.2,
		},
		{
			SupportingTools: []string{"tool-b", "tool-a"},
			Score:           0.9,
		},
		{
			SupportingTools: []string{"tool-c"},
			Score:           0.5,
		},
	}

	filtered := filterDuplicateSets(sets)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 sets after filtering, got %d", len(filtered))
	}

	if filtered[0].Score != 0.9 && filtered[1].Score != 0.9 {
		t.Fatal("expected higher score to remain for duplicate set")
	}
}

func TestNormalizeSetScore(t *testing.T) {
	sets := []FeatureSet{
		{Score: 2},
		{Score: 1},
	}

	normalized := normalizeSetScore(sets)
	if len(normalized) != 2 {
		t.Fatalf("expected 2 normalized sets, got %d", len(normalized))
	}

	if normalized[0].Score+normalized[1].Score < 0.999 || normalized[0].Score+normalized[1].Score > 1.001 {
		t.Fatalf("expected normalized scores to sum to 1, got %v", normalized[0].Score+normalized[1].Score)
	}
}

func TestSetFileIdentity(t *testing.T) {
	sets := []FeatureSet{
		{Score: 0.2},
		{Score: 0.3},
		{Score: 0.5},
	}

	adjusted := setFileIdentity(sets, 1)

	expected := []float64{0, 1, 0}
	for i, set := range adjusted {
		if set.Score != expected[i] {
			t.Fatalf("expected score %v at index %d, got %v", expected[i], i, set.Score)
		}
	}
}

func TestApplyFileIdentityRules_NoMatch(t *testing.T) {
	serverConfig = ServerConfig{
		FileIdentityRules: []FileIdentityRule{
			{
				Conditions: []FeatureCondition{
					{Feature: "format:type", Value: "pdf"},
				},
			},
		},
	}

	sets := []FeatureSet{
		{Score: 0.8, Features: map[string]MergeFeatureValue{"format:type": {Value: "xml"}}},
	}

	result := applyFileIdentityRules(sets)
	if len(result) != 1 || result[0].Score != 0.8 {
		t.Fatal("expected set to remain unchanged")
	}
}
