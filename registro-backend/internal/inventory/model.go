package inventory

import (
	"time"
)

// Accounting categories conforming to D.I. 129/2018
const (
	CategoryScientifico = "Materiale Scientifico e Laboratori"
	CategoryInformatica = "Hardware, LIM e Dispositivi Mobili"
	CategoryArredi      = "Arredi e Mobili d'Ufficio"
	CategoryPNRR        = "PNRR Scuola 4.0 - Laboratori per le Professioni Digitali"
	CategoryBiblioteca  = "Beni Librari e Materiale Didattico"
)

type AssetItem struct {
	ID                 string     `json:"id" db:"id"`
	SchoolID           string     `json:"school_id" db:"school_id"`
	InventoryNumber    string     `json:"inventory_number" db:"inventory_number"` // e.g. "INV-2026-0042"
	Category           string     `json:"category" db:"category"`
	Description        string     `json:"description" db:"description"`
	SerialNumber       string     `json:"serial_number" db:"serial_number"`
	BuildingLocation   string     `json:"building_location" db:"building_location"`
	RoomLocation       string     `json:"room_location" db:"room_location"`
	AcquisitionDate    time.Time  `json:"acquisition_date" db:"acquisition_date"`
	InitialValue       float64    `json:"initial_value" db:"initial_value"`
	DepreciationRate   float64    `json:"depreciation_rate" db:"depreciation_rate"` // e.g. 20.0%
	CurrentValue       float64    `json:"current_value" db:"current_value"`
	BarcodeData        string     `json:"barcode_data" db:"barcode_data"`
	QRData             string     `json:"qr_data" db:"qr_data"`
	Status             string     `json:"status" db:"status"` // "in_uso", "in_manutenzione", "scaricato"
	LastInspectionDate *time.Time `json:"last_inspection_date,omitempty" db:"last_inspection_date"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

type CreateAssetRequest struct {
	Category         string  `json:"category" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	SerialNumber     string  `json:"serial_number"`
	BuildingLocation string  `json:"building_location" binding:"required"`
	RoomLocation     string  `json:"room_location" binding:"required"`
	AcquisitionDate  string  `json:"acquisition_date" binding:"required"` // YYYY-MM-DD
	InitialValue     float64 `json:"initial_value" binding:"required,gt=0"`
	DepreciationRate float64 `json:"depreciation_rate"` // Defaults to 20%
}

type AssetLabel struct {
	InventoryNumber  string `json:"inventory_number"`
	Description      string `json:"description"`
	BarcodeData      string `json:"barcode_data"`
	QRCodeData       string `json:"qr_code_data"`
	BuildingLocation string `json:"building_location"`
	RoomLocation     string `json:"room_location"`
}

type LoanContract struct {
	ID                   string     `json:"id" db:"id"`
	SchoolID             string     `json:"school_id" db:"school_id"`
	AssetID              string     `json:"asset_id" db:"asset_id"`
	StudentID            string     `json:"student_id" db:"student_id"`
	ParentUserID         string     `json:"parent_user_id" db:"parent_user_id"`
	ContractNumber       string     `json:"contract_number" db:"contract_number"`
	HandoverDate         time.Time  `json:"handover_date" db:"handover_date"`
	ExpectedReturnDate   time.Time  `json:"expected_return_date" db:"expected_return_date"`
	ActualReturnDate     *time.Time `json:"actual_return_date,omitempty" db:"actual_return_date"`
	DeviceConditionNotes string     `json:"device_condition_notes" db:"device_condition_notes"`
	ParentSignatureToken string     `json:"parent_signature_token" db:"parent_signature_token"`
	SignedAt             *time.Time `json:"signed_at,omitempty" db:"signed_at"`
	Status               string     `json:"status" db:"status"` // "attivo", "restituito", "danneggiato", "non_restituito"
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

type CreateLoanRequest struct {
	AssetID              string `json:"asset_id" binding:"required"`
	StudentID            string `json:"student_id" binding:"required"`
	ParentUserID         string `json:"parent_user_id" binding:"required"`
	ExpectedReturnDate   string `json:"expected_return_date" binding:"required"` // YYYY-MM-DD
	DeviceConditionNotes string `json:"device_condition_notes"`
}
