package repository // import "klaus/repository"

import (
	"time"

	"github.com/shopspring/decimal"
)

type TicketCategoryAggregate struct {
	TicketID       uint64          `db:"id"`
	CategoryName   string          `db:"name"`
	Rating         uint64          `db:"rating"`
	CategoryID     uint64          `db:"category_id"`
	CategoryWeight decimal.Decimal `db:"weight"`
}

type RatingWeightAggregate struct {
	Rating uint64          `db:"rating"`
	Weight decimal.Decimal `db:"weight"`
}

type CountedCategoryWeight struct {
	ID        uint64          `db:"id"`
	Name      string          `db:"name"`
	Count     uint64          `db:"total"`
	Rating    uint64          `db:"rating"`
	Weight    decimal.Decimal `db:"weight"`
	CreatedAt time.Time       `db:"created_at"`
}
