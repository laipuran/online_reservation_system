package model

import "time"

type ServiceProvider struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	BusinessName string    `json:"business_name"`
	Description  string    `json:"description,omitempty"`
	Address      string    `json:"address,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
