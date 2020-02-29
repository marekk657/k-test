package handler // import "klaus/handler"

import (
	"klaus/repository"
)

type CategoryScoresHandler interface {
	Handle(p Period) ([]CategoryScore, error)
}

func NewCategoryScoresHandler(repo repository.Repository) CategoryScoresHandler {
	return categoryScoreHandler{
		repo: repo,
	}
}

type categoryScoreHandler struct {
	repo repository.Repository
}

func (h categoryScoreHandler) Handle(p Period) ([]CategoryScore, error) {
	categories, err := h.repo.CountCategoryWeights(p.Start, p.End)
	if err != nil {
		return nil, err
	}

	periodLengthDays := p.End.Sub(p.Start).Hours() / 24
	shouldAggregateByWeeks := periodLengthDays >= 30

	categoriesMap := h.categoriesToMap(categories, shouldAggregateByWeeks)
	categoryScores := h.createFinalCategoryScores(categoriesMap, shouldAggregateByWeeks)

	return categoryScores, nil
}

func (h categoryScoreHandler) categoriesToMap(categories []repository.CountedCategoryWeight, shouldAggregateByWeeks bool) map[uint64]CategoryScore {
	categoriesMap := make(map[uint64]CategoryScore, 10)
	for _, c := range categories {
		v, ok := categoriesMap[c.ID]
		if !ok {
			v = CategoryScore{
				Name:         c.Name,
				RatingsCount: c.Count,
				CategoryID:   c.ID,
				scores:       make(Scores, 0, 100),
			}

			if shouldAggregateByWeeks {
				v.ScoredDates = make([]ScoredDate, 53) // initialize array of 53 elements, so every week would inserted into specific index
			} else {
				v.ScoredDates = make([]ScoredDate, 0, 100)
			}
		}

		v.RatingsCount += c.Count
		v.scores = append(v.scores, NewScore(c.Rating, c.Weight))

		if !shouldAggregateByWeeks {
			v.ScoredDates = append(v.ScoredDates, ScoredDate{
				Date:  c.CreatedAt,
				Score: NewScore(c.Rating, c.Weight),
			})
		} else {
			_, week := c.CreatedAt.ISOWeek()
			scoredDate := v.ScoredDates[week]
			if scoredDate.Date.IsZero() {
				scoredDate.Date = c.CreatedAt
			}

			scoredDate.scores = append(scoredDate.scores, NewScore(c.Rating, c.Weight))
			v.ScoredDates[week] = scoredDate
		}

		categoriesMap[c.ID] = v
	}

	return categoriesMap
}

func (h categoryScoreHandler) createFinalCategoryScores(categoriesMap map[uint64]CategoryScore, shouldAggregateByWeeks bool) []CategoryScore {
	categoryScores := make([]CategoryScore, 0, len(categoriesMap))
	for _, v := range categoriesMap {
		v.TotalScore = v.scores.Avg()
		v.scores = nil

		if !shouldAggregateByWeeks {

			categoryScores = append(categoryScores, v)
			// daily aggregations already have correct Score
			continue
		}

		// recalculate weekly aggregations Score
		scoredDates := make([]ScoredDate, 0, 53)
		for _, sd := range v.ScoredDates {
			if sd.Date.IsZero() {
				// skip empty weeks
				continue
			}

			scoredDates = append(scoredDates, ScoredDate{
				Date:  sd.Date,
				Score: sd.scores.Avg(),
			})
			sd.scores = nil
		}
		v.ScoredDates = scoredDates
		categoryScores = append(categoryScores, v)
	}
	return categoryScores
}
