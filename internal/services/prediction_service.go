package services

import (
	"Hackathon/internal/core/structs"
	"errors"
)

type PredictionService struct {
}

func NewPredictionService() *PredictionService {
	return &PredictionService{}
}

func (p *PredictionService) PredictPortStat(history []structs.PortInfo) (structs.PortInfo, error) {

	if len(history) < 2 {
		return structs.PortInfo{}, errors.New("Insufficient data for forecasting")
	}

	n := float64(len(history))
	var sumX, sumY, sumXY, sumX2 float64

	for i, stat := range history {
		x := float64(i)
		y := float64(stat.InOctets)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return structs.PortInfo{}, errors.New("Error calculating regression")
	}
	a := (n*sumXY - sumX*sumY) / denominator
	b := (sumY - a*sumX) / n

	nextX := n
	predictedInOctets := a*nextX + b

	predictedStat := history[len(history)-1]
	predictedStat.InOctets = uint(predictedInOctets)

	return predictedStat, nil
}
