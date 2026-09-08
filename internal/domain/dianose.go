package domain

type DiagnosisRequest struct {
	SymptomCodes []string
	FactorCodes  []string
}

type DiagnosisResponse struct {
	P float64
	CauseWithRemedy
}

type SortedDiagnose struct {
	P       float64
	CauseID int64
}
