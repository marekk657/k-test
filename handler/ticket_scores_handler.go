package handler // import "klaus/handler"

import "klaus/repository"

type TicketScoresHandler interface {
	Handle(p Period) ([]TicketScore, error)
}

func NewTicketScoresHandler(repo repository.Repository) TicketScoresHandler {
	return ticketScoresHandler{
		repo: repo,
	}
}

type ticketScoresHandler struct {
	repo repository.Repository
}

func (h ticketScoresHandler) Handle(p Period) ([]TicketScore, error) {
	rows, err := h.repo.QueryTicketWithCategory(p.Start, p.End)
	if err != nil {
		return nil, err
	}

	ticketMap := make(map[uint64]TicketScore, 100)
	for _, r := range rows {
		v, ok := ticketMap[r.TicketID]
		if !ok {
			v = TicketScore{
				TicketID:   r.TicketID,
				categories: make(map[uint64]category, 5),
			}
		}

		cat := v.categories[r.CategoryID]
		cat.Name = r.CategoryName
		cat.Scores = append(cat.Scores, NewScore(r.Rating, r.CategoryWeight))
		v.categories[r.CategoryID] = cat
		ticketMap[r.TicketID] = v
	}

	ticketScores := make([]TicketScore, 0, len(ticketMap))
	for _, t := range ticketMap {
		t.calculateCategoryScores()
		ticketScores = append(ticketScores, t)
	}

	return ticketScores, nil
}
