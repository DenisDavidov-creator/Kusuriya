package domain

type Factor struct {
	ID          int64
	Code        string
	Name        string
	Description string
}

type FactorProfile struct {
	ID       int64
	CauseID  int64
	FactorID int64
	Modifier float64
}

type FactorEvidence struct {
	CauseID    int64
	FactorID   int64
	FactorCode string
	Modifier   *float64
}
