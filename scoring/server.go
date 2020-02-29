package scoring // import "klaus"

import (
	context "context"
	"errors"
	"klaus/handler"

	"github.com/golang/protobuf/ptypes"
)

type server struct {
	categoryScores               handler.CategoryScoresHandler
	ticketScoresHandler          handler.TicketScoresHandler
	overallQualityScoreHandler   handler.OverallQualityScoreHandler
	overPeriodScoreChangeHandler handler.OverPeriodScoreChangeHandler
}

func NewServer(
	categoryScores handler.CategoryScoresHandler,
	ticketScoresHandler handler.TicketScoresHandler,
	overallQualityScoreHandler handler.OverallQualityScoreHandler,
	overPeriodScoreChangeHandler handler.OverPeriodScoreChangeHandler) ScoringServer {
	return &server{
		categoryScores:               categoryScores,
		ticketScoresHandler:          ticketScoresHandler,
		overallQualityScoreHandler:   overallQualityScoreHandler,
		overPeriodScoreChangeHandler: overPeriodScoreChangeHandler,
	}
}

func (s *server) GetCategoryScoresInPeriod(ctx context.Context, req *Request) (*CategoryScoresResponse, error) {
	p, err := s.toPeriod(req)
	if err != nil {
		return nil, err
	}

	scores, err := s.categoryScores.Handle(p)
	if err != nil {
		return nil, err
	}

	return s.toCategoryScoresResponse(scores), nil
}

func (s *server) GetTicketScoresInPeriod(ctx context.Context, req *Request) (*TicketScoresResponse, error) {
	p, err := s.toPeriod(req)
	if err != nil {
		return nil, err
	}

	tickets, err := s.ticketScoresHandler.Handle(p)
	if err != nil {
		return nil, err
	}

	return s.toTicketScoresResponse(tickets), nil
}

func (s *server) GetOverallQualityScore(ctx context.Context, req *Request) (*OverallQualityScoreResponse, error) {
	p, err := s.toPeriod(req)
	if err != nil {
		return nil, err
	}

	score, err := s.overallQualityScoreHandler.Handle(p)
	if err != nil {
		return nil, err
	}

	scoreValue, _ := score.ValueAsPercentage.Float64()
	return &OverallQualityScoreResponse{
		Score: &Score{
			Score:        scoreValue,
			ScoreLabeled: score.ValueAsLabeled,
		},
	}, nil
}

func (s *server) GetOverPeriodScoreChange(ctx context.Context, req *Request) (*OverPeriodScoreChangeResponse, error) {
	p, err := s.toPeriod(req)
	if err != nil {
		return nil, err
	}

	overPeriodChange, err := s.overPeriodScoreChangeHandler.Handle(p)
	if err != nil {
		return nil, err
	}

	return s.toOverPeriodScoreChangeResponse(overPeriodChange), nil
}

func (s *server) toPeriod(req *Request) (handler.Period, error) {
	start, err := ptypes.Timestamp(req.GetPeriodStart())
	if err != nil {
		return handler.Period{}, errors.New("invalid period start")
	}

	end, err := ptypes.Timestamp(req.GetPeriodEnd())
	if err != nil {
		return handler.Period{}, errors.New("invalid period end")
	}

	return handler.Period{
		Start: start,
		End:   end,
	}, nil
}

func (s *server) toCategoryScoresResponse(scores []handler.CategoryScore) *CategoryScoresResponse {
	resp := CategoryScoresResponse{
		CategoryScores: make([]*CategoryScoresResponse_CategoryScore, 0, len(scores)),
	}

	for _, cs := range scores {
		totalScore, _ := cs.TotalScore.ValueAsPercentage.Float64()
		responseCategoryScore := CategoryScoresResponse_CategoryScore{
			Category:     cs.Name,
			RatingsCount: cs.RatingsCount,
			TotalScore: &Score{
				Score:        totalScore,
				ScoreLabeled: cs.TotalScore.ValueAsLabeled,
			},
			ScoredDates: make([]*CategoryScoresResponse_ScoredDate, 0, len(cs.ScoredDates)),
		}

		for _, sd := range cs.ScoredDates {
			tp, _ := ptypes.TimestampProto(sd.Date)
			dateScore, _ := sd.Score.ValueAsPercentage.Float64()
			responseScoredDate := CategoryScoresResponse_ScoredDate{
				Date: tp,
				Score: &Score{
					Score:        dateScore,
					ScoreLabeled: sd.Score.ValueAsLabeled,
				},
			}
			responseCategoryScore.ScoredDates = append(responseCategoryScore.ScoredDates, &responseScoredDate)
		}

		resp.CategoryScores = append(resp.CategoryScores, &responseCategoryScore)
	}

	return &resp
}

func (s *server) toTicketScoresResponse(tickets []handler.TicketScore) *TicketScoresResponse {
	resp := TicketScoresResponse{
		Tickets: make([]*TicketScoresResponse_TicketScore, 0, len(tickets)),
	}

	for _, t := range tickets {
		responseTicketScore := TicketScoresResponse_TicketScore{
			Id:         t.TicketID,
			Categories: make([]*TicketScoresResponse_Category, 0, len(t.Categories)),
		}

		for _, c := range t.Categories {
			score, _ := c.TotalScore.ValueAsPercentage.Float64()
			responseCategory := TicketScoresResponse_Category{
				CategoryName: c.Name,
				Score: &Score{
					Score:        score,
					ScoreLabeled: c.TotalScore.ValueAsLabeled,
				},
			}

			responseTicketScore.Categories = append(responseTicketScore.Categories, &responseCategory)
		}

		resp.Tickets = append(resp.Tickets, &responseTicketScore)
	}

	return &resp
}

func (s *server) toOverPeriodScoreChangeResponse(overPeriodChange handler.OverPeriodChange) *OverPeriodScoreChangeResponse {
	currentValue, _ := overPeriodChange.CurrentPeriod.ValueAsPercentage.Float64()
	previousValue, _ := overPeriodChange.PreviousPeriod.ValueAsPercentage.Float64()
	differenceValue, _ := overPeriodChange.Difference.ValueAsPercentage.Float64()
	return &OverPeriodScoreChangeResponse{
		CurrentPeriod: &Score{
			Score:        currentValue,
			ScoreLabeled: overPeriodChange.CurrentPeriod.ValueAsLabeled,
		},
		PreviousPeriod: &Score{
			Score:        previousValue,
			ScoreLabeled: overPeriodChange.PreviousPeriod.ValueAsLabeled,
		},
		Difference: &Score{
			Score:        differenceValue,
			ScoreLabeled: overPeriodChange.Difference.ValueAsLabeled,
		},
	}
}
