package repository_test

import (
	"klaus/repository"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestQueryTicketWithCategory(t *testing.T) {
	db := createDbConnection(t)
	defer db.Close()

	log.SetOutput(os.Stdout)

	repo := repository.NewRepository(db)

	start := time.Date(2019, 03, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, 03, 31, 0, 0, 0, 0, time.UTC)

	// Act
	res, err := repo.QueryTicketWithCategory(start, end)
	if err != nil {
		t.Fatal("error occurred:", err)
	}

	if len(res) == 0 {
		t.Fatal("0 records returned")
	}
}

func TestQueryRatingsWithWeight(t *testing.T) {
	db := createDbConnection(t)
	defer db.Close()

	repo := repository.NewRepository(db)

	start := time.Date(2019, 03, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, 03, 31, 0, 0, 0, 0, time.UTC)

	// Act
	res, err := repo.QueryRatingsWithWeight(start, end)
	if err != nil {
		t.Fatal("error occurred:", err)
	}

	if len(res) == 0 {
		t.Fatal("0 records returned")
	}
}

func TestCountCategoryWeights(t *testing.T) {
	db := createDbConnection(t)
	defer db.Close()

	log.SetOutput(os.Stdout)

	repo := repository.NewRepository(db)

	start := time.Date(2019, 03, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, 03, 31, 0, 0, 0, 0, time.UTC)

	// Act
	res, err := repo.CountCategoryWeights(start, end)
	if err != nil {
		t.Fatal("error occurred:", err)
	}

	if len(res) == 0 {
		t.Fatal("0 records returned")
	}
}

func createDbConnection(t *testing.T) *sqlx.DB {
	db, err := repository.CreateDBConnection()
	if err != nil {
		t.Fatal("failed to connect with db:", err)
	}
	return db
}
