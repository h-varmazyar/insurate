package core

import (
	"context"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
	scoreRepo "github.com/h-varmazyar/insurate/internal/core/repository/score"
	"time"
)

type ScoreCalculateParams struct {
	*scoreRepo.Score
}

func (params *ScoreCalculateParams) CalculateScore(_ context.Context) float64 {
	scoreValue := float64(0)
	scoreCount := float64(0)
	if s := params.calculateDrivingOffencesComplexity(); s != nil {
		scoreValue += *s
		scoreCount++
	}
	if s := params.calculateNegativeScoreComplexity(); s != nil {
		scoreValue += *s
		scoreCount++
	}
	if s := params.calculateAgeComplexity(); s != nil {
		scoreValue += *s
		scoreCount++
	}
	if s := params.calculateGenderComplexity(); s != nil {
		scoreValue += *s
		scoreCount++
	}
	return scoreValue / scoreCount
}

func (params *ScoreCalculateParams) calculateDrivingOffencesComplexity() *float64 {
	if params.DrivingOffences == nil {
		return nil
	}
	sumScore := float64(0)
	for _, bill := range params.DrivingOffences {
		switch bill.Code {
		case "2001", "2013", "2015", "2009", "2003", "2002", "2007", "2004", "2005", "2010", "2012", "2039", "2048", "2077", "2018":
			sumScore += 1
		default:
			sumScore += 0.5
		}
	}
	score := sumScore / float64(len(params.DrivingOffences))

	return &score
}

func (params *ScoreCalculateParams) calculateNegativeScoreComplexity() *float64 {
	if params.NegativeScore < 31 {
		params.NegativeScore /= 2
	} else if params.NegativeScore > 55 {
		params.NegativeScore = 55
	}
	resp := float64(params.NegativeScore / 55)
	return &resp
}

//reference: https://p3o.ir/women10
func (params *ScoreCalculateParams) calculateGenderComplexity() *float64 {
	if params.Person.Gender == personRepo.Unknown {
		return nil
	}
	resp := float64(0)
	if params.Person.Gender == personRepo.Men {
		resp = 0.9
	} else {
		resp = 0.1
	}
	return &resp
}

func (params *ScoreCalculateParams) calculateAgeComplexity() *float64 {
	if params.Person.BirthDate == time.Unix(0, 0) {
		return nil
	}
	ageDuration := time.Now().Sub(params.Person.BirthDate)
	ageYear := ageDuration / (time.Hour * 24 * 365)
	resp := float64(0)
	if ageYear < 18 {
		resp = 1
	} else if ageYear < 23 { //18-22
		resp = 0.9
	} else if ageYear < 26 { //23-25
		resp = 1
	} else if ageYear < 31 { //26-30
		resp = 0.9
	} else if ageYear < 41 { //31-40
		resp = 0.1
	} else if ageYear < 61 { //40-60
		resp = 0.5
	} else { //more than 60
		resp = 0.8
	}
	return &resp
}
