package domain

type SymptomProfile struct {
	ID         int64
	CauseID    int64
	SymptomID  int64
	Likelihood float64
}
