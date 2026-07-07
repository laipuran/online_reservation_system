package model

import "time"

type Service struct {
	ID              int64     `json:"id"`
	ProviderID      int64     `json:"provider_id"`
	CategoryID      int64     `json:"category_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	Price           float64   `json:"price"`
	DurationMinutes int       `json:"duration_minutes"`
	ImageURL        string    `json:"image_url,omitempty"`
	Status          string    `json:"status"`
	AvgRating       float64   `json:"avg_rating"`
	ReviewCount     int       `json:"review_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ServiceProviderSummary struct {
	ID           int64  `json:"id"`
	BusinessName string `json:"business_name"`
}

type CategorySummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ServiceView struct {
	ID              int64                  `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description,omitempty"`
	Provider        ServiceProviderSummary `json:"provider"`
	Category        CategorySummary        `json:"category"`
	Price           float64                `json:"price"`
	DurationMinutes int                    `json:"duration_minutes"`
	ImageURL        string                 `json:"image_url,omitempty"`
	Status          string                 `json:"status"`
	AvgRating       float64                `json:"avg_rating"`
	ReviewCount     int                    `json:"review_count"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type ServiceFilter struct {
	Keyword    string
	CategoryID *int64
	ProviderID *int64
	MinPrice   *float64
	MaxPrice   *float64
	Status     string
	SortBy     string
	SortOrder  string
	Page       int
	PageSize   int
}
