package handler // import "klaus/handler"

import (
	"fmt"
	"klaus/repository"
	"time"

	"github.com/shopspring/decimal"
)

var (
	// Rating steps: 0 - 5
	ratingStepValue = decimal.NewFromFloat(100 / 6)
)

type Period struct {
	Start time.Time
	End   time.Time
}

type ScoredDate struct {
	Date  time.Time
	Score Score

	scores Scores
}

type CategoryScore struct {
	CategoryID   uint64
	Name         string
	RatingsCount uint64
	ScoredDates  []ScoredDate
	TotalScore   Score

	scores Scores
}

type OverPeriodChange struct {
	CurrentPeriod  Score
	PreviousPeriod Score
	Difference     Score
}

type TicketScore struct {
	TicketID   uint64
	Categories []CategoryScore

	categories map[uint64]category
}

type category struct {
	Name   string
	Scores Scores
}

func (ts *TicketScore) calculateCategoryScores() {
	categories := make([]CategoryScore, 0, len(ts.categories))
	for _, cat := range ts.categories {
		categoryScore := CategoryScore{
			Name:         cat.Name,
			RatingsCount: uint64(len(cat.Scores)),
			TotalScore:   cat.Scores.Avg(),
		}
		categories = append(categories, categoryScore)
	}
	ts.categories = nil
	ts.Categories = categories
}

type Score struct {
	weight decimal.Decimal

	Value             decimal.Decimal
	ValueAsPercentage decimal.Decimal
	ValueAsLabeled    string
}

type Scores []Score

func (ss Scores) Avg() Score {
	var valueSum decimal.Decimal
	var weightsSum decimal.Decimal
	for _, s := range ss {
		valueSum = valueSum.Add(s.Value)
		weightsSum = weightsSum.Add(s.weight)
	}

	if weightsSum.IsZero() {
		var score Score
		score.setFields()
		return score
	}

	weightedValue := valueSum.Div(weightsSum)

	s := Score{
		Value: weightedValue,
	}
	s.setFields()
	return s
}

func NewScore(rating uint64, weight decimal.Decimal) Score {
	score := decimal.NewFromInt(int64(rating)).Mul(weight).Mul(ratingStepValue)
	s := Score{
		Value:  score,
		weight: weight,
	}
	s.setFields()
	return s
}

func (s *Score) setFields() {
	s.ValueAsLabeled = fmt.Sprintf("%s%%", s.Value.StringFixed(2))
	if s.Value.IsZero() {
		s.ValueAsLabeled = "N/A"
	}

	s.ValueAsPercentage = s.Value.DivRound(decimal.NewFromInt(100), 4)
}

func calculateScoreFromRatings(ratings []repository.RatingWeightAggregate) Score {
	scores := make(Scores, 0, len(ratings))
	for _, s := range ratings {
		scores = append(scores, NewScore(s.Rating, s.Weight))
	}

	return scores.Avg()
}
