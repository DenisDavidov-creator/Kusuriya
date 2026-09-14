package handler

import (
	"Kusuriya/internal/domain"
	"Kusuriya/internal/httpx"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// checkValidate asserts the (problems, err) pair every Validate returns:
// no problems must mean err == nil; any problem must mean ErrValidation.
func checkValidate(t *testing.T, wantProblems []string, problems []string, err error) {
	t.Helper()
	if len(wantProblems) == 0 {
		assert.Nil(t, problems)
		assert.NoError(t, err)
		return
	}
	assert.Equal(t, wantProblems, problems)
	assert.ErrorIs(t, err, httpx.ErrValidation)
}

func TestCreateCauseRequest_Validate(t *testing.T) {
	valid := CreateCauseRequest{Code: "c1", Name: "n", Type: "poison", Prior: 0.5}

	tests := []struct {
		name   string
		mutate func(r *CreateCauseRequest)
		want   []string
	}{
		{"valid poison", func(r *CreateCauseRequest) {}, nil},
		{"valid disease", func(r *CreateCauseRequest) { r.Type = "disease" }, nil},
		{"prior = 1 is allowed", func(r *CreateCauseRequest) { r.Prior = 1 }, nil},
		{"prior = 0 is out of range", func(r *CreateCauseRequest) { r.Prior = 0 }, []string{"prior out of range (0,1]"}},
		{"prior > 1 is out of range", func(r *CreateCauseRequest) { r.Prior = 1.01 }, []string{"prior out of range (0,1]"}},
		{"negative prior", func(r *CreateCauseRequest) { r.Prior = -0.1 }, []string{"prior out of range (0,1]"}},
		{"empty code", func(r *CreateCauseRequest) { r.Code = "" }, []string{"code is empty"}},
		{"empty name", func(r *CreateCauseRequest) { r.Name = "" }, []string{"name is empty"}},
		{"unknown type", func(r *CreateCauseRequest) { r.Type = "banana" }, []string{"type is incorrect"}},
		{"type is case sensitive", func(r *CreateCauseRequest) { r.Type = "Poison" }, []string{"type is incorrect"}},
		{
			"all problems are reported together, in field order",
			func(r *CreateCauseRequest) { *r = CreateCauseRequest{} },
			[]string{"code is empty", "name is empty", "type is incorrect", "prior out of range (0,1]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := valid
			tt.mutate(&r)
			problems, err := r.Validate()
			checkValidate(t, tt.want, problems, err)
		})
	}
}

func TestCreateCauseRequest_ToDomain(t *testing.T) {
	remedy := int64(7)
	r := CreateCauseRequest{Code: "c1", Name: "n", Type: "poison", Prior: 0.5, Description: "d", RemedyID: &remedy}

	assert.Equal(t, domain.Cause{Code: "c1", Name: "n", Type: "poison", Prior: 0.5, Description: "d", RemedyID: &remedy}, r.ToDomain())
}

func TestUpdateCauseRequest_ToDomain(t *testing.T) {
	// the id in the body is ignored; the URL id wins
	r := UpdateCauseRequest{ID: 999, Code: "c1", Name: "n", Type: "poison", Prior: 0.5}

	assert.Equal(t, int64(42), r.ToDomain(42).ID)
}

func TestSimpleRequests_Validate(t *testing.T) {
	t.Run("remedy", func(t *testing.T) {
		p, err := (&CreateRemedyRequest{Name: "r"}).Validate()
		checkValidate(t, nil, p, err)
		p, err = (&CreateRemedyRequest{}).Validate()
		checkValidate(t, []string{"name is empty"}, p, err)
	})
	t.Run("factor", func(t *testing.T) {
		p, err := CreateFactorRequest{Code: "f", Name: "n"}.Validate()
		checkValidate(t, nil, p, err)
		p, err = CreateFactorRequest{}.Validate()
		checkValidate(t, []string{"code is empty", "name is empty"}, p, err)
	})
	t.Run("category", func(t *testing.T) {
		p, err := CreateCategoryRequest{Name: "c"}.Validate()
		checkValidate(t, nil, p, err)
		p, err = CreateCategoryRequest{}.Validate()
		checkValidate(t, []string{"name is empty"}, p, err)
	})
	t.Run("symptom", func(t *testing.T) {
		p, err := CreateSymptomRequest{Code: "s", Name: "n", CategoryID: 1}.Validate()
		checkValidate(t, nil, p, err)
		p, err = CreateSymptomRequest{Code: "s", Name: "n"}.Validate()
		checkValidate(t, []string{"category is empty"}, p, err)
	})
}

func TestCreateSymptomsProfileRequests_Validate(t *testing.T) {
	tests := []struct {
		name     string
		profiles []CreateSymptomProfileRequest
		want     []string
	}{
		{"empty list passes", nil, nil},
		{"valid", []CreateSymptomProfileRequest{{SymptomID: 1, Likelihood: 0.5}, {SymptomID: 2, Likelihood: 1}}, nil},
		{"likelihood = 0", []CreateSymptomProfileRequest{{SymptomID: 1, Likelihood: 0}}, []string{"1: likelihood is empty"}},
		{"likelihood > 1", []CreateSymptomProfileRequest{{SymptomID: 1, Likelihood: 1.5}}, []string{"1: likelihood is empty"}},
		{"missing symptom_id", []CreateSymptomProfileRequest{{Likelihood: 0.5}}, []string{"0: symptom_id is empty"}},
		{
			"problems in several items are all reported",
			[]CreateSymptomProfileRequest{{SymptomID: 1, Likelihood: 0}, {SymptomID: 2, Likelihood: 2}},
			[]string{"1: likelihood is empty", "2: likelihood is empty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := CreateSymptomsProfileRequests{Profiles: tt.profiles}.Validate()
			checkValidate(t, tt.want, p, err)
		})
	}
}

func TestCreateSymptomsProfileRequests_ToDomains(t *testing.T) {
	rs := CreateSymptomsProfileRequests{Profiles: []CreateSymptomProfileRequest{
		{SymptomID: 10, Likelihood: 0.8},
		{SymptomID: 11, Likelihood: 0.2},
	}}

	got := rs.ToDomains(5)

	assert.Equal(t, []domain.SymptomProfile{
		{CauseID: 5, SymptomID: 10, Likelihood: 0.8},
		{CauseID: 5, SymptomID: 11, Likelihood: 0.2},
	}, got, "cause id is stamped on every profile")
	assert.Nil(t, CreateSymptomsProfileRequests{}.ToDomains(5), "no profiles -> nil")
}

func TestCreateFactorsProfileRequests_Validate(t *testing.T) {
	p, err := CreateFactorsProfileRequests{Profiles: []CreateFactorProfileRequest{{FactorID: 1, Modifier: 0.5}}}.Validate()
	checkValidate(t, nil, p, err)

	p, err = CreateFactorsProfileRequests{Profiles: []CreateFactorProfileRequest{{Modifier: 3}}}.Validate()
	checkValidate(t, []string{"0: factor_id is empty", "0: modifier is empty"}, p, err)
}

func TestResponseJSONFieldNames(t *testing.T) {
	// JSON keys are the API contract; pin them so a tag typo is caught here, not by a client.
	b, _ := json.Marshal(FactorProfileResponse{ID: 1, CauseID: 2, FactorID: 3, Modifier: 0.5})
	assert.JSONEq(t, `{"id":1,"cause_id":2,"factor_id":3,"modifier":0.5}`, string(b))

	b, _ = json.Marshal(SymptomProfileResponse{ID: 1, CauseID: 2, SymptomID: 3, Likelihood: 0.5})
	assert.JSONEq(t, `{"id":1,"cause_id":2,"symptom_id":3,"likelihood":0.5}`, string(b))
}
