package services

import (
	"Hackathon/internal/core/structs"
	"testing"
)

func TestPredictPortStat(t *testing.T) {
	ps := NewPredictionService()

	historyInsufficient := []structs.PortInfo{
		{InOctets: 1000},
	}
	_, err := ps.PredictPortStat(historyInsufficient)
	if err == nil {
		t.Errorf("Expected an error for insufficient data")
	}

	history := []structs.PortInfo{
		{InOctets: 1000},
		{InOctets: 1500},
		{InOctets: 2000},
	}
	predictedStat, err := ps.PredictPortStat(history)
	if err != nil {
		t.Errorf("Did not expect an error, got: %v", err)
	}

	if predictedStat.InOctets == 0 {
		t.Errorf("Predicted InOctets value should not be zero")
	}

	expectedInOctets := uint(2500)
	if predictedStat.InOctets != expectedInOctets {
		t.Errorf("Expected predicted InOctets = %d, got %d", expectedInOctets, predictedStat.InOctets)
	}
}
