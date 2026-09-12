package domain

import (
	"math"
	"sort"
)

func EngineCause(causeDoms []Cause, factors []FactorEvidence, symptoms []SymptomEvidence) []SortedDiagnose {
	var causeScores = make(map[int64]float64, len(causeDoms))
	for _, cause := range causeDoms {
		causeScores[cause.ID] = math.Log(cause.Prior)
	}
	for _, factor := range factors {
		var factorModifier float64
		if factor.Modifier != nil {
			factorModifier = math.Log(*factor.Modifier)
		} else {
			factorModifier = math.Log(0.05)
		}
		score, ok := causeScores[factor.CauseID]
		if ok {
			causeScores[factor.CauseID] = score + factorModifier
		} else {
			causeScores[factor.CauseID] = math.Log(0.05) + factorModifier

		}
	}
	for _, symptom := range symptoms {
		var symptomLikelihood float64
		if symptom.Likelihood != nil {
			symptomLikelihood = math.Log(*symptom.Likelihood)
		} else {
			symptomLikelihood = math.Log(0.05)
		}
		score, ok := causeScores[symptom.CauseID]
		if ok {
			causeScores[symptom.CauseID] = score + symptomLikelihood
		} else {
			causeScores[symptom.CauseID] = math.Log(0.05) + symptomLikelihood
		}
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

	return sortedScores
}

func CombineCausesWithScores(scores []SortedDiagnose, causes []CauseWithRemedy) []DiagnosisResponse {
	var causesMap = make(map[int64]CauseWithRemedy, len(causes))
	for _, value := range causes {
		causesMap[value.ID] = value
	}

	var finalScores = make([]DiagnosisResponse, 0, len(causes))
	for _, value := range scores {
		cause, ok := causesMap[value.CauseID]
		if ok {
			score := DiagnosisResponse{
				P:               value.P,
				CauseWithRemedy: cause,
			}
			finalScores = append(finalScores, score)
		}
	}
	return finalScores
}
