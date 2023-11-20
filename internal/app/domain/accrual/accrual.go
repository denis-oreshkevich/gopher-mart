package accrual

const (
	RegisteredStatus = "REGISTERED" //заказ зарегистрирован, но вознаграждение не рассчитано;
	InvalidStatus    = "INVALID"    //заказ не принят к расчёту, и вознаграждение не будет начислено;
	ProcessingStatus = "PROCESSING" //расчёт начисления в процессе;
	ProcessedStatus  = "PROCESSED"  //расчёт начисления окончен;
)

//go:generate easyjson -all accrual.go
type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func New(order string, status string, accrual float64) Accrual {
	return Accrual{
		Order:   order,
		Status:  status,
		Accrual: accrual,
	}
}
