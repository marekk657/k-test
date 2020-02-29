package main

import (
	"klaus/handler"
	"klaus/repository"
	"klaus/scoring"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := repository.CreateDBConnection()
	if err != nil {
		log.Fatal("failed to connect with db:", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	scoring.RegisterScoringServer(s, createServer(db))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func createServer(db *sqlx.DB) scoring.ScoringServer {
	repo := repository.NewRepository(db)

	categoryScoresHandler := handler.NewCategoryScoresHandler(repo)
	ticketScoresHandler := handler.NewTicketScoresHandler(repo)
	overallQualityScoreHandler := handler.NewOverallQualityScoreHandler(repo)
	overPeriodScoreChangeHandler := handler.NewOverPeriodScoreChangeHandler(repo)

	return scoring.NewServer(categoryScoresHandler,
		ticketScoresHandler,
		overallQualityScoreHandler,
		overPeriodScoreChangeHandler)
}
