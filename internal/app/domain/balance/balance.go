package balance

//go:generate easyjson -all balance.go
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
	UserID    string  `json:"-"`
}

func New(current float64, withdrawn int, userID string) Balance {
	return Balance{
		Current:   current,
		Withdrawn: withdrawn,
		UserID:    userID,
	}
}
