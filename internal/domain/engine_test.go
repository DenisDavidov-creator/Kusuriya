package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func f(v float64) *float64 { return &v }

// Expected probabilities below are computed by hand, not with the engine's
// formula, so the test cannot agree with a bug by construction:
//
//	score(cause) = prior * Π likelihood * Π modifier   (unmatched evidence = 0.05)
//	P(cause)     = score(cause) / Σ score
func TestEngineCause(t *testing.T) {
	tests := []struct {
		name     string
		causes   []Cause
		factors  []FactorEvidence
		symptoms []SymptomEvidence
		wantP    []SortedDiagnose
		wantEv   map[int64][]Contribution
	}{
		{
			name:   "no evidence: ranking follows priors, evidence is empty",
			causes: []Cause{{ID: 1, Prior: 0.1}, {ID: 2, Prior: 0.3}},
			// A = 0.1, B = 0.3, sum = 0.4
			wantP:  []SortedDiagnose{{CauseID: 2, P: 0.75}, {CauseID: 1, P: 0.25}},
			wantEv: map[int64][]Contribution{},
		},
		{
			name:   "single cause always has P = 1",
			causes: []Cause{{ID: 1, Prior: 0.01}},
			symptoms: []SymptomEvidence{
				{CauseID: 1, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.8)},
			},
			wantP: []SortedDiagnose{{CauseID: 1, P: 1}},
			wantEv: map[int64][]Contribution{
				1: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.8), Matched: true}},
			},
		},
		{
			name:   "equal priors: matched symptom wins over unmatched",
			causes: []Cause{{ID: 1, Prior: 0.5}, {ID: 2, Prior: 0.5}},
			symptoms: []SymptomEvidence{
				{CauseID: 1, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.8)},
				{CauseID: 2, SymptomID: 3, SymptomCode: "s1", Likelihood: nil},
			},
			// A = 0.5*0.8 = 0.4, B = 0.5*0.05 = 0.025, sum = 0.425
			wantP: []SortedDiagnose{{CauseID: 1, P: 0.4 / 0.425}, {CauseID: 2, P: 0.025 / 0.425}},
			wantEv: map[int64][]Contribution{
				1: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.8), Matched: true}},
				2: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: nil, Matched: false}},
			},
		},
		{
			name:   "strong evidence overturns a weaker prior",
			causes: []Cause{{ID: 1, Prior: 0.1}, {ID: 2, Prior: 0.3}},
			symptoms: []SymptomEvidence{
				{CauseID: 1, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.9)},
				{CauseID: 2, SymptomID: 3, SymptomCode: "s1", Likelihood: nil},
			},
			// A = 0.1*0.9 = 0.09, B = 0.3*0.05 = 0.015, sum = 0.105
			wantP: []SortedDiagnose{{CauseID: 1, P: 0.09 / 0.105}, {CauseID: 2, P: 0.015 / 0.105}},
			wantEv: map[int64][]Contribution{
				1: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.9), Matched: true}},
				2: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: nil, Matched: false}},
			},
		},
		{
			name:   "factors and symptoms multiply together; evidence lists factors first",
			causes: []Cause{{ID: 1, Prior: 0.2}, {ID: 2, Prior: 0.2}},
			factors: []FactorEvidence{
				{CauseID: 1, FactorID: 8, FactorCode: "f1", Modifier: f(0.5)},
				{CauseID: 2, FactorID: 8, FactorCode: "f1", Modifier: f(0.5)},
			},
			symptoms: []SymptomEvidence{
				{CauseID: 1, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.6)},
				{CauseID: 2, SymptomID: 3, SymptomCode: "s1", Likelihood: f(0.3)},
			},
			// A = 0.2*0.5*0.6 = 0.06, B = 0.2*0.5*0.3 = 0.03, sum = 0.09
			wantP: []SortedDiagnose{{CauseID: 1, P: 0.06 / 0.09}, {CauseID: 2, P: 0.03 / 0.09}},
			wantEv: map[int64][]Contribution{
				1: {
					{Kind: "factor", KindID: 8, Code: "f1", Weight: f(0.5), Matched: true},
					{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.6), Matched: true},
				},
				2: {
					{Kind: "factor", KindID: 8, Code: "f1", Weight: f(0.5), Matched: true},
					{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.3), Matched: true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, gotEv := EngineCause(tt.causes, tt.factors, tt.symptoms)

			require.Len(t, gotP, len(tt.wantP))
			sum := 0.0
			for i := range tt.wantP {
				assert.Equal(t, tt.wantP[i].CauseID, gotP[i].CauseID, "rank %d", i)
				assert.InDelta(t, tt.wantP[i].P, gotP[i].P, 1e-9, "P for cause %d", gotP[i].CauseID)
				sum += gotP[i].P
			}
			assert.InDelta(t, 1.0, sum, 1e-9, "probabilities must sum to 1")

			assert.Equal(t, tt.wantEv, gotEv)
		})
	}
}

func TestCombineCausesWithScores(t *testing.T) {
	scores := []SortedDiagnose{{CauseID: 2, P: 0.7}, {CauseID: 1, P: 0.3}, {CauseID: 9, P: 0.0}}
	causes := []CauseWithRemedy{
		{ID: 1, Code: "c1"},
		{ID: 2, Code: "c2"},
		// cause 9 is deliberately absent: the service only loads the top-N causes
	}
	ev := map[int64][]Contribution{
		2: {{Kind: "symptom", KindID: 3, Code: "s1", Weight: f(0.8), Matched: true}},
	}

	got := CombineCausesWithScores(scores, causes, ev)

	want := []DiagnosisResponse{
		{P: 0.7, CauseWithRemedy: CauseWithRemedy{ID: 2, Code: "c2"}, Evidence: ev[2]},
		{P: 0.3, CauseWithRemedy: CauseWithRemedy{ID: 1, Code: "c1"}, Evidence: nil},
	}
	assert.Equal(t, want, got, "keeps score order, drops causes that were not loaded, attaches evidence")
}
