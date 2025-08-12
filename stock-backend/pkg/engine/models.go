package engine

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	Code       uuid.UUID  `json:"code"`
	Ticker     *string    `json:"ticker"`
	Company    *string    `json:"company"`
	Brokerage  *string    `json:"brokerage"`
	Action     *string    `json:"action"`
	RatingFrom *string    `json:"rating_from"`
	RatingTo   *string    `json:"rating_to"`
	TargetFrom *float64   `json:"target_from"`
	TargetTo   *float64   `json:"target_to"`
	RecordTime *time.Time `json:"record_time"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type PaginatedStocksResponse struct {
	Stocks      []Stock `json:"stocks"`
	CurrentPage int     `json:"current_page"`
	NextPage    *int    `json:"next_page"`
	Total       int     `json:"total"`
	PerPage     int     `json:"per_page"`
}

type Recommendation struct {
	Stocks   []Stock `json:"stocks"`
	Message  *string `json:"message"`
	Priority *string `json:"priority"`
}
