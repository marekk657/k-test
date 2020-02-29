package handler // import "klaus/handler"

import "klaus/repository"

type OverallQualityScoreHandler interface {
	Handle(p Period) (Score, error)
}

func NewOverallQualityScoreHandler(repo repository.Repository) OverallQualityScoreHandler {
	return overallQualityScoreHandler{
		repo: repo,
	}
}

type overallQualityScoreHandler struct {
	repo repository.Repository
}

func (h overallQualityScoreHandler) Handle(p Period) (Score, error) {
	ratings, err := h.repo.QueryRatingsWithWeight(p.Start, p.End)
	if err != nil {
		return Score{}, err
	}

	return calculateScoreFromRatings(ratings), nil
}
