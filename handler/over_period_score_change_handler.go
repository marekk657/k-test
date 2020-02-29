package handler // import "klaus/handler"

import (
	"klaus/repository"
	"time"
)

type OverPeriodScoreChangeHandler interface {
	Handle(p Period) (OverPeriodChange, error)
}

func NewOverPeriodScoreChangeHandler(repo repository.Repository) OverPeriodScoreChangeHandler {
	return overPeriodScoreChangeHandler{
		repo: repo,
	}
}

type overPeriodScoreChangeHandler struct {
	repo repository.Repository
}

func (h overPeriodScoreChangeHandler) Handle(p Period) (OverPeriodChange, error) {
	selectedPeriod, err := h.repo.QueryRatingsWithWeight(p.Start, p.End)
	if err != nil {
		return OverPeriodChange{}, err
	}

	periodLengthNano := p.End.Sub(p.Start).Nanoseconds()
	previousPeriodStart := p.Start.Add(-(time.Duration(periodLengthNano) * time.Nanosecond))
	previousPeriodEnd := p.End.Add(-(time.Duration(periodLengthNano) * time.Nanosecond))

	previousPeriod, err := h.repo.QueryRatingsWithWeight(previousPeriodStart, previousPeriodEnd)
	if err != nil {
		return OverPeriodChange{}, err
	}

	current := calculateScoreFromRatings(selectedPeriod)
	previous := calculateScoreFromRatings(previousPeriod)
	diff := Score{
		Value: previous.Value.Div(current.Value),
	}
	diff.setFields()
	return OverPeriodChange{
		CurrentPeriod:  current,
		PreviousPeriod: previous,
		Difference:     diff,
	}, nil
}
