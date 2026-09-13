package domain

type Symptom struct {
	ID          int64
	Code        string
	Name        string
	Description string
	CategoryID  int64
}

type Category struct {
	ID   int64
	Name string
}

type SymptomProfile struct {
	ID         int64
	CauseID    int64
	SymptomID  int64
	Likelihood float64
}

type SymptomEvidence struct {
	CauseID     int64
	SymptomID   int64
	SymptomCode string
	Likelihood  *float64
}
