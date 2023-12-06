package accrual

type Status string

const (
	StatusRegistered Status = "REGISTERED" //заказ зарегистрирован, но вознаграждение не рассчитано;
	StatusInvalid    Status = "INVALID"    //заказ не принят к расчёту, и вознаграждение не будет начислено;
	StatusProcessing Status = "PROCESSING" //расчёт начисления в процессе;
	StatusProcessed  Status = "PROCESSED"  //расчёт начисления окончен;
)

//go:generate easyjson -all accrual.go
type Accrual struct {
	Order   string  `json:"order"`
	Status  Status  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func New(order string, status Status, accrual float64) Accrual {
	return Accrual{
		Order:   order,
		Status:  status,
		Accrual: accrual,
	}
}
