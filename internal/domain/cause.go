package domain

type Cause struct {
	ID          int64
	Code        string
	Name        string
	Type        string
	Prior       float64
	Description string
	RemedyID    *int64
}

type CauseWithRemedy struct {
	ID                int64
	Code              string
	Name              string
	Type              string
	Prior             float64
	Description       string
	RemedyID          *int64
	RemedyName        *string
	RemedyDescription *string
	RemedyNotes       *string
}

type DiagnosisRequest struct {
	SymptomCodes []string
	FactorCodes  []string
}

type DiagnosisResponse struct {
	P float64
	CauseWithRemedy
	Evidence []Contribution
}

type SortedDiagnose struct {
	P       float64
	CauseID int64
}

type Contribution struct {
	Kind    string
	KindID  int64
	Code    string
	Weight  *float64
	Matched bool
}
