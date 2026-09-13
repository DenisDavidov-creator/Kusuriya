package domain

import (
	"math"
	"sort"
)

func EngineCause(causeDoms []Cause, factors []FactorEvidence, symptoms []SymptomEvidence) ([]SortedDiagnose, map[int64][]Contribution) {
	var causeScores = make(map[int64]float64, len(causeDoms))
	var evidence = make(map[int64][]Contribution, len(causeDoms))
	for _, cause := range causeDoms {
		causeScores[cause.ID] = math.Log(cause.Prior)
	}
	for _, factor := range factors {
		var factorModifier float64
		contr := Contribution{
			Kind:   "factor",
			KindID: factor.FactorID,
			Code:   factor.FactorCode,
		}
		if factor.Modifier != nil {
			factorModifier = math.Log(*factor.Modifier)
			contr.Weight = factor.Modifier
			contr.Matched = true
		} else {
			factorModifier = math.Log(0.05)
			contr.Weight = factor.Modifier
			contr.Matched = false
		}
		score := causeScores[factor.CauseID]
		causeScores[factor.CauseID] = score + factorModifier

		evidence[factor.CauseID] = append(evidence[factor.CauseID], contr)
	}
	for _, symptom := range symptoms {
		var symptomLikelihood float64
		contr := Contribution{
			Kind:   "symptom",
			KindID: symptom.SymptomID,
			Code:   symptom.SymptomCode,
		}
		if symptom.Likelihood != nil {
			symptomLikelihood = math.Log(*symptom.Likelihood)
			contr.Weight = symptom.Likelihood
			contr.Matched = true
		} else {
			symptomLikelihood = math.Log(0.05)
			contr.Weight = symptom.Likelihood
			contr.Matched = false
		}
		score := causeScores[symptom.CauseID]
		causeScores[symptom.CauseID] = score + symptomLikelihood

		evidence[symptom.CauseID] = append(evidence[symptom.CauseID], contr)
	}

	SummOfE := 0.0
	for ID, k := range causeScores {
		pow := math.Pow(math.E, k)
		causeScores[ID] = pow
		SummOfE += pow
	}
	var sortedScores = make([]SortedDiagnose, 0, len(causeDoms))
	for ID, k := range causeScores {
		score := SortedDiagnose{
			P:       k / SummOfE,
			CauseID: ID,
		}
		sortedScores = append(sortedScores, score)
	}

	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].P > sortedScores[j].P
	})

	return sortedScores, evidence
}

func CombineCausesWithScores(scores []SortedDiagnose, causes []CauseWithRemedy, contribution map[int64][]Contribution) []DiagnosisResponse {
	var causesMap = make(map[int64]CauseWithRemedy, len(causes))
	for _, value := range causes {
		causesMap[value.ID] = value
	}

	var finalScores = make([]DiagnosisResponse, 0, len(causes))
	for _, value := range scores {
		cause, ok := causesMap[value.CauseID]
		evidence := contribution[value.CauseID]
		if ok {
			score := DiagnosisResponse{
				P:               value.P,
				CauseWithRemedy: cause,
				Evidence:        evidence,
			}
			finalScores = append(finalScores, score)
		}
	}
	return finalScores
}
