package model

import "time"

type Review struct {
	ID            int64     `json:"id"`
	ReservationID int64     `json:"reservation_id"`
	UserID        int64     `json:"user_id"`
	ServiceID     int64     `json:"service_id"`
	Rating        int16     `json:"rating"`
	Comment       string    `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
