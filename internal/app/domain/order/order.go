package order

import "time"

const (
	NewStatus        = "NEW"        //заказ загружен в систему, но не попал в обработку;
	ProcessingStatus = "PROCESSING" //вознаграждение за заказ рассчитывается;
	InvalidStatus    = "INVALID"    //система расчёта вознаграждений отказала в расчёте;
	ProcessedStatus  = "PROCESSED"  //данные по заказу проверены и информация о расчёте успешно получена.
)

//go:generate easyjson -all order.go
type Order struct {
	ID         string    `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UserID     string    `json:"-"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

//easyjson:json
type Orders []Order

func New(id string, number string, status string, userID string,
	accrual float64, uploadedAt time.Time) Order {
	return Order{
		ID:         id,
		Number:     number,
		Status:     status,
		UserID:     userID,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
	}
}
