package withdrawal

import "time"

//go:generate easyjson -all withdrawal.go
type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

//easyjson:json
type Withdrawals []Withdrawal

func New(order string, sum float64, processedAt time.Time) Withdrawal {
	return Withdrawal{
		Order:       order,
		Sum:         sum,
		ProcessedAt: processedAt,
	}
}
