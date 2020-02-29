package handler_test

import (
	"klaus/handler"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewScore(t *testing.T) {
	expectedScore := handler.Score{
		Value:             decimal.NewFromFloat(44.80),
		ValueAsPercentage: decimal.NewFromFloat(0.4480),
		ValueAsLabeled:    "44.80%",
	}

	// Act
	score := handler.NewScore(4, decimal.NewFromFloat(0.7))

	// Assert
	if !expectedScore.ValueAsPercentage.Equal(score.ValueAsPercentage) {
		t.Fatalf("wrong value %% returned! got=%v; want=%v", score.ValueAsPercentage, expectedScore.ValueAsPercentage)
	}

	if expectedScore.ValueAsLabeled != score.ValueAsLabeled {
		t.Fatalf("wrong label returned! got=%s; want=%s", score.ValueAsLabeled, expectedScore.ValueAsLabeled)
	}

	if !expectedScore.Value.Equal(score.Value) {
		t.Fatalf("wrong value returned! got=%v; want=%v", score.Value, expectedScore.Value)
	}
}

func TestScores_Avg(t *testing.T) {
	scores := handler.Scores{
		handler.NewScore(3, decimal.NewFromFloat(1)),
		handler.NewScore(5, decimal.NewFromFloat(.75)),
	}

	expectedScore := handler.Score{
		ValueAsPercentage: decimal.NewFromFloat(0.6171),
		ValueAsLabeled:    "61.71%",
	}

	// Act
	score := scores.Avg()

	// Assert
	if !expectedScore.ValueAsPercentage.Equal(score.ValueAsPercentage) {
		t.Fatalf("wrong value %% returned! got=%v; want=%v", score.ValueAsPercentage, expectedScore.ValueAsPercentage)
	}

	if expectedScore.ValueAsLabeled != score.ValueAsLabeled {
		t.Fatalf("wrong label returned! got=%s; want=%s", score.ValueAsLabeled, expectedScore.ValueAsLabeled)
	}
}

func TestScores_Avg_ZeroWeights(t *testing.T) {
	scores := handler.Scores{
		handler.NewScore(3, decimal.Zero),
	}

	expectedScore := handler.Score{
		ValueAsPercentage: decimal.Zero,
		ValueAsLabeled:    "N/A",
	}

	// Act
	score := scores.Avg()

	// Assert
	if !expectedScore.ValueAsPercentage.Equal(score.ValueAsPercentage) {
		t.Fatalf("wrong value %% returned! got=%v; want=%v", score.ValueAsPercentage, expectedScore.ValueAsPercentage)
	}

	if expectedScore.ValueAsLabeled != score.ValueAsLabeled {
		t.Fatalf("wrong label returned! got=%s; want=%s", score.ValueAsLabeled, expectedScore.ValueAsLabeled)
	}
}
