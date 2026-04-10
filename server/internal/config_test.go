package internal

import "testing"

func TestFeatureConditionIsFulfilled(t *testing.T) {
	tests := []struct {
		name     string
		cond     FeatureCondition
		value    interface{}
		expected bool
	}{
		{
			name:     "exact string match",
			cond:     FeatureCondition{Value: "abc"},
			value:    "abc",
			expected: true,
		},
		{
			name:     "exact string mismatch",
			cond:     FeatureCondition{Value: "abc"},
			value:    "def",
			expected: false,
		},
		{
			name:     "exact bool match",
			cond:     FeatureCondition{Value: true},
			value:    true,
			expected: true,
		},
		{
			name:     "no condition configured",
			cond:     FeatureCondition{},
			value:    "abc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cond.IsFulfilled(tt.value); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestTriggerIsTriggered(t *testing.T) {
	trigger := Trigger{
		Conditions: []FeatureCondition{
			{Feature: "format:type", Value: "pdf"},
			{Feature: "format:valid", Value: true},
		},
	}

	results := map[string]ToolResult{
		"tool-a": {
			Features: map[string]ToolFeatureValue{
				"format:type":  {Value: "pdf"},
				"format:valid": {Value: true},
			},
		},
	}

	isTriggered, matches := trigger.IsTriggered(results)
	if !isTriggered {
		t.Fatal("expected trigger to fire")
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestToolConfigIsTriggered(t *testing.T) {
	tc := ToolConfig{
		Triggers: []Trigger{
			{
				Conditions: []FeatureCondition{
					{Feature: "format:type", Value: "pdf"},
				},
			},
		},
	}

	results := map[string]ToolResult{
		"tool-a": {
			Features: map[string]ToolFeatureValue{
				"format:type": {Value: "pdf"},
			},
		},
	}

	isTriggered, _ := tc.IsTriggered(results)
	if !isTriggered {
		t.Fatal("expected tool config to be triggered")
	}
}

func TestConditionalWeightIsFulfilled(t *testing.T) {
	tests := []struct {
		name     string
		weight   ConditionalWeight
		result   ToolResult
		expected bool
	}{
		{
			name: "all conditions fulfilled",
			weight: ConditionalWeight{
				Value: 0.9,
				Conditions: []FeatureCondition{
					{Feature: "format:type", Value: "pdf"},
					{Feature: "format:valid", Value: true},
				},
			},
			result: ToolResult{
				Features: map[string]ToolFeatureValue{
					"format:type":  {Value: "pdf"},
					"format:valid": {Value: true},
				},
			},
			expected: true,
		},
		{
			name: "missing feature",
			weight: ConditionalWeight{
				Value: 0.9,
				Conditions: []FeatureCondition{
					{Feature: "format:type", Value: "pdf"},
				},
			},
			result:   ToolResult{Features: map[string]ToolFeatureValue{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.weight.IsFulfilled(tt.result); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestWeightGetWeight(t *testing.T) {
	score := 0.7

	tests := []struct {
		name     string
		weight   Weight
		result   ToolResult
		expected float64
	}{
		{
			name: "conditional weight wins",
			weight: Weight{
				Default: 0.1,
				ConditionalWeights: []ConditionalWeight{
					{
						Value: 0.9,
						Conditions: []FeatureCondition{
							{Feature: "format:type", Value: "pdf"},
						},
					},
				},
				ProvidedByTool: true,
			},
			result: ToolResult{
				Score: &score,
				Features: map[string]ToolFeatureValue{
					"format:type": {Value: "pdf"},
				},
			},
			expected: 0.9,
		},
		{
			name: "tool provided weight wins when no conditional match",
			weight: Weight{
				Default:        0.1,
				ProvidedByTool: true,
			},
			result:   ToolResult{Score: &score},
			expected: 0.7,
		},
		{
			name: "default weight used last",
			weight: Weight{
				Default: 0.1,
			},
			result:   ToolResult{},
			expected: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.weight.GetWeight(tt.result); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMergeConditionIsFulfilled(t *testing.T) {
	extraction := `^v(\d+)$`

	tests := []struct {
		name        string
		cond        MergeCondition
		featureKey  string
		fs1         map[string]MergeFeatureValue
		fs2         map[string]ToolFeatureValue
		expectedOK  bool
		expectedHit bool
	}{
		{
			name:       "exact equality strong link",
			cond:       MergeCondition{},
			featureKey: "format:type",
			fs1: map[string]MergeFeatureValue{
				"format:type": {Value: "pdf"},
			},
			fs2: map[string]ToolFeatureValue{
				"format:type": {Value: "pdf"},
			},
			expectedOK:  true,
			expectedHit: true,
		},
		{
			name:        "regex extraction strong link",
			cond:        MergeCondition{ValueRegEx: &extraction},
			featureKey:  "format:version",
			fs1:         map[string]MergeFeatureValue{"format:version": {Value: "v12"}},
			fs2:         map[string]ToolFeatureValue{"format:version": {Value: "v12"}},
			expectedOK:  true,
			expectedHit: true,
		},
		{
			name:       "missing feature on one side still mergeable",
			cond:       MergeCondition{},
			featureKey: "format:type",
			fs1:        map[string]MergeFeatureValue{"format:type": {Value: "pdf"}},
			fs2:        map[string]ToolFeatureValue{},
			expectedOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, strongLink := tt.cond.IsFulfilled(tt.featureKey, tt.fs1, tt.fs2)
			if ok != tt.expectedOK {
				t.Fatalf("expected ok=%v, got %v", tt.expectedOK, ok)
			}
			if strongLink != tt.expectedHit {
				t.Fatalf("expected strongLink=%v, got %v", tt.expectedHit, strongLink)
			}
		})
	}
}

func TestFeatureSetConfigGetFeatureConfig(t *testing.T) {
	cfg := FeatureSetConfig{
		Features: []FeatureConfig{
			{Key: "format:type", MergeOrder: 10},
			{Key: "format:version", MergeOrder: 20},
		},
	}

	tests := []struct {
		name     string
		key      string
		expected bool
		order    uint
	}{
		{
			name:     "found",
			key:      "format:version",
			expected: true,
			order:    20,
		},
		{
			name:     "not found",
			key:      "missing",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := cfg.GetFeatureConfig(tt.key)
			if ok != tt.expected {
				t.Fatalf("expected ok=%v, got %v", tt.expected, ok)
			}
			if ok && got.MergeOrder != tt.order {
				t.Fatalf("expected merge order %d, got %d", tt.order, got.MergeOrder)
			}
		})
	}
}
