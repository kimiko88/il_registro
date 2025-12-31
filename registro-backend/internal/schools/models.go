package schools

import "time"

// School represents a school entity (aligned with 001_initial_schema.sql + code column)
type School struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Code          string    `json:"code" db:"code"` // Codice meccanografico - frontend expects this
	Address       string    `json:"address" db:"address"`
	City          string    `json:"city" db:"city"`
	Province      string    `json:"province" db:"province"`
	Cap           string    `json:"cap" db:"zip_code"`
	Phone         string    `json:"phone" db:"phone"`
	Email         string    `json:"email" db:"email"`
	CodiceFiscale string    `json:"codice_fiscale" db:"codice_fiscale"`
	Iban          string    `json:"iban" db:"iban"`
	Website       string    `json:"website" db:"website"`
	LogoURL       string    `json:"logo_url" db:"logo_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSchoolRequest represents the request to create a school (frontend fields)
type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
}

// UpdateSchoolRequest represents the request to update a school
type UpdateSchoolRequest struct {
	Name    *string `json:"name"`
	Code    *string `json:"code"`
	Address *string `json:"address"`
	Email   *string `json:"email"`
	Phone   *string `json:"phone"`
	City    *string `json:"city"`
}

// ListParams represents query parameters for listing schools
type ListParams struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
}
