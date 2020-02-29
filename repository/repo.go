package repository // import "klaus/repository"

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDBConnection() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", "file:database.db?mode=ro")
}

type Repository interface {
	QueryTicketWithCategory(start, end time.Time) ([]TicketCategoryAggregate, error)
	QueryRatingsWithWeight(start, end time.Time) ([]RatingWeightAggregate, error)
	CountCategoryWeights(start, end time.Time) ([]CountedCategoryWeight, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) QueryTicketWithCategory(start, end time.Time) ([]TicketCategoryAggregate, error) {
	sql := `SELECT r.ticket_id AS id, r.rating, rc.name, rc.weight, rc.id AS category_id
	FROM ratings AS r
	JOIN rating_categories AS rc ON r.rating_category_id = rc.id
	WHERE r.created_at > $1 AND r.created_at < $2
	ORDER BY r.ticket_id ASC`

	var result []TicketCategoryAggregate
	if err := r.db.Select(&result, sql, start, end); err != nil {
		return nil, err
	}
	return result, nil
}

func (r repository) QueryRatingsWithWeight(start, end time.Time) ([]RatingWeightAggregate, error) {
	sql := `SELECT r.rating, rc.weight
    FROM ratings as r
	JOIN rating_categories AS rc on r.rating_category_id = rc.id
	WHERE r.created_at > $1 AND r.created_at < $2`

	var result []RatingWeightAggregate
	if err := r.db.Select(&result, sql, start, end); err != nil {
		return nil, err
	}
	return result, nil
}

func (r repository) CountCategoryWeights(start, end time.Time) ([]CountedCategoryWeight, error) {
	sql := `SELECT rc.id, r.rating, rc.name, rc.weight, r.created_at, COUNT(r.rating) AS total
    FROM ratings AS r
	JOIN rating_categories AS rc ON r.rating_category_id = rc.id
	WHERE r.created_at > $1 AND r.created_at < $2
	GROUP BY rc.name, rc.weight, r.rating
	ORDER BY r.created_at ASC`

	var result []CountedCategoryWeight
	if err := r.db.Select(&result, sql, start, end); err != nil {
		return nil, err
	}
	return result, nil
}
