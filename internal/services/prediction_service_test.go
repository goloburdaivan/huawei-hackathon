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
		t.Errorf("Ожидалась ошибка при недостаточном количестве данных")
	}

	history := []structs.PortInfo{
		{InOctets: 1000},
		{InOctets: 1500},
		{InOctets: 2000},
	}
	predictedStat, err := ps.PredictPortStat(history)
	if err != nil {
		t.Errorf("Не ожидалось ошибки, получили: %v", err)
	}

	if predictedStat.InOctets == 0 {
		t.Errorf("Предсказанное значение InOctets не должно быть нулевым")
	}

	expectedInOctets := uint(2500)
	if predictedStat.InOctets != expectedInOctets {
		t.Errorf("Ожидалось предсказанное InOctets = %d, получили %d", expectedInOctets, predictedStat.InOctets)
	}
}

func TestPredictPortStat_RegressionError(t *testing.T) {
	ps := NewPredictionService()

	history := []structs.PortInfo{
		{InOctets: 1000},
		{InOctets: 1000},
	}
	_, err := ps.PredictPortStat(history)
	if err == nil {
		t.Errorf("Ожидалась ошибка при невозможности вычислить регрессию")
	}
}
