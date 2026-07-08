package model

import "time"

type Reservation struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ServiceID int64     `json:"service_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReservationServiceProviderSummary struct {
	ID           int64  `json:"id"`
	BusinessName string `json:"business_name"`
}

type ReservationServiceSummary struct {
	ID       int64                             `json:"id"`
	Title    string                            `json:"title"`
	Provider ReservationServiceProviderSummary `json:"provider"`
}

type ReservationView struct {
	ID        int64                     `json:"id"`
	Service   ReservationServiceSummary `json:"service"`
	StartTime time.Time                 `json:"start_time"`
	EndTime   time.Time                 `json:"end_time"`
	Status    string                    `json:"status"`
	Note      string                    `json:"note,omitempty"`
	CreatedAt time.Time                 `json:"created_at"`
}
