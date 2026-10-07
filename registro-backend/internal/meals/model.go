package meals

import (
	"time"
)

type SpecialDiet struct {
	ID                    string     `json:"id" db:"id"`
	SchoolID              string     `json:"school_id" db:"school_id"`
	StudentID             string     `json:"student_id" db:"student_id"`
	StudentName           string     `json:"student_name,omitempty"`
	DietCategory          string     `json:"diet_category" db:"diet_category"` // "sanitaria", "etico_religiosa"
	SpecificDiet          string     `json:"specific_diet" db:"specific_diet"` // "celiachia", "intolleranza_lattosio", "favismo", "no_maiale", "vegetariana", "vegana"
	MedicalCertificateURL string     `json:"medical_certificate_url" db:"medical_certificate_url"`
	CertificateExpiryDate *time.Time `json:"certificate_expiry_date,omitempty" db:"certificate_expiry_date"`
	AllergensList         string     `json:"allergens_list" db:"allergens_list"`
	IsApproved            bool       `json:"is_approved" db:"is_approved"`
	Notes                 string     `json:"notes" db:"notes"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
}

type MealRollCallRequest struct {
	ClassID    string   `json:"class_id" binding:"required"`
	Date       string   `json:"date" binding:"required"` // YYYY-MM-DD
	PresentIDs []string `json:"present_ids"`
}

type CateringOrderReport struct {
	ID               string         `json:"id"`
	SchoolID         string         `json:"school_id"`
	Date             time.Time      `json:"date"`
	TransmissionTime time.Time      `json:"transmission_time"`
	TotalStandard    int            `json:"total_standard"`
	TotalDiets       int            `json:"total_diets"`
	TotalMeals       int            `json:"total_meals"`
	DietsBreakdown   map[string]int `json:"diets_breakdown"`
	CateringProvider string         `json:"catering_provider"`
	Status           string         `json:"status"` // "trasmesso", "in_attesa"
}

type MealWallet struct {
	ID                  string    `json:"id" db:"id"`
	SchoolID            string    `json:"school_id" db:"school_id"`
	StudentID           string    `json:"student_id" db:"student_id"`
	Balance             float64   `json:"balance" db:"balance"`
	MealPrice           float64   `json:"meal_price" db:"meal_price"`
	LowBalanceThreshold float64   `json:"low_balance_threshold" db:"low_balance_threshold"`
	IsBlocked           bool      `json:"is_blocked" db:"is_blocked"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type WalletTransaction struct {
	ID           string    `json:"id" db:"id"`
	WalletID     string    `json:"wallet_id" db:"wallet_id"`
	MovementType string    `json:"movement_type" db:"movement_type"` // "addebito_pasto", "ricarica_pagopa"
	Amount       float64   `json:"amount" db:"amount"`
	BalanceAfter float64   `json:"balance_after" db:"balance_after"`
	PagoPAIUV    string    `json:"pagopa_iuv,omitempty" db:"pagopa_iuv"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type SpecialDietRequest struct {
	StudentID             string `json:"student_id" binding:"required"`
	DietCategory          string `json:"diet_category" binding:"required"` // "sanitaria", "etico_religiosa"
	SpecificDiet          string `json:"specific_diet" binding:"required"`
	MedicalCertificateURL string `json:"medical_certificate_url"`
	CertificateExpiryDate string `json:"certificate_expiry_date"` // YYYY-MM-DD
	AllergensList         string `json:"allergens_list"`
	Notes                 string `json:"notes"`
}
