package main

import (
	"context"
	"klaus/scoring"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := scoring.NewScoringClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start, _ := ptypes.TimestampProto(time.Date(2019, 03, 01, 0, 0, 0, 0, time.UTC))
	end, _ := ptypes.TimestampProto(time.Date(2019, 03, 31, 0, 0, 0, 0, time.UTC))

	r, err := c.GetCategoryScoresInPeriod(ctx, &scoring.Request{
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		log.Fatalf("couldnt send [GetCategoryScoresInPeriod] message: %v", err)
	}
	log.Printf("[GetCategoryScoresInPeriod] response: %v", r.GetCategoryScores())

	r2, err := c.GetOverPeriodScoreChange(ctx, &scoring.Request{
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		log.Fatalf("couldnt send [GetOverPeriodScoreChange] message: %v", err)
	}
	log.Printf("[GetOverPeriodScoreChange] response: CurrentPeriod %v", r2.GetCurrentPeriod())
	log.Printf("[GetOverPeriodScoreChange] response: PreviousPeriod %v", r2.GetPreviousPeriod())
	log.Printf("[GetOverPeriodScoreChange] response: PeriodDiff %v", r2.GetDifference())

	r3, err := c.GetOverallQualityScore(ctx, &scoring.Request{
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		log.Fatalf("couldnt send [GetOverallQualityScore] message: %v", err)
	}
	log.Printf("[GetOverPeriodScoreChange] response: GetScore %v", r3.GetScore())

	// r4, err := c.GetTicketScoresInPeriod(ctx, &scoring.Request{
	// 	PeriodStart: start,
	// 	PeriodEnd:   end,
	// })
	// if err != nil {
	// 	log.Fatalf("couldnt send [GetTicketScoresInPeriod] message: %v", err)
	// }
	// log.Printf("[GetTicketScoresInPeriod] response: GetScore %v", r4.GetTickets())
}
