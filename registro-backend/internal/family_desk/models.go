package family_desk

import (
	"errors"
	"fmt"
	"time"
)

const (
	TypeDelegaRitiro                   = "delega_ritiro"
	TypeUscitaAutonomaUnder14          = "uscita_autonoma_under14"
	TypeEsoneroMotoria                 = "esonero_motoria"
	TypeSomministrazioneFarmaci        = "somministrazione_farmaci"
	TypeNullaOsta                      = "nulla_osta"
	TypeCertificatoIscrizioneFrequenza = "certificato_iscrizione_frequenza"
)

const (
	StatusSubmitted     = "submitted"
	StatusInIstruttoria = "in_istruttoria"
	StatusApproved      = "approved"
	StatusRejected      = "rejected"
)

type FamilyRequest struct {
	ID              string                 `json:"id"`
	SchoolID        string                 `json:"school_id"`
	StudentID       string                 `json:"student_id"`
	ParentID        string                 `json:"parent_id"`
	RequestType     string                 `json:"request_type"`
	FormData        map[string]interface{} `json:"form_data"`
	AttachmentURLs  []string               `json:"attachment_urls"`
	Status          string                 `json:"status"`
	RejectionReason string                 `json:"rejection_reason,omitempty"`
	ProtocolNumber  string                 `json:"protocol_number,omitempty"`
	ReviewedBy      *string                `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time             `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type PermanentDelegate struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	SchoolID          string    `json:"school_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	TaxCode           string    `json:"tax_code"`
	Relationship      string    `json:"relationship"`
	Phone             string    `json:"phone"`
	IDCardDetails     string    `json:"id_card_details"`
	IDCardFileURL     string    `json:"id_card_file_url,omitempty"`
	IsValid           bool      `json:"is_valid"`
	ApprovedRequestID *string   `json:"approved_request_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

func IsValidRequestType(rt string) bool {
	switch rt {
	case TypeDelegaRitiro, TypeUscitaAutonomaUnder14, TypeEsoneroMotoria,
		TypeSomministrazioneFarmaci, TypeNullaOsta, TypeCertificatoIscrizioneFrequenza:
		return true
	default:
		return false
	}
}

func ValidateRequestFormData(req *FamilyRequest) error {
	if req.RequestType == TypeUscitaAutonomaUnder14 {
		consent, ok := req.FormData["consent_given"].(bool)
		if !ok || !consent {
			return errors.New("è necessario fornire esplicito consenso informato per l'uscita autonoma ai sensi della L. 172/2017")
		}
	}
	return nil
}

func BuildDelegateFromRequest(req *FamilyRequest) (*PermanentDelegate, error) {
	if req.RequestType != TypeDelegaRitiro {
		return nil, fmt.Errorf("la richiesta non è di tipo delega_ritiro ma %s", req.RequestType)
	}

	getString := func(key string) string {
		if val, ok := req.FormData[key].(string); ok {
			return val
		}
		return ""
	}

	fn := getString("first_name")
	ln := getString("last_name")
	tc := getString("tax_code")
	rel := getString("relationship")
	phone := getString("phone")
	card := getString("id_card_details")
	cardURL := getString("id_card_file_url")

	if fn == "" || ln == "" || tc == "" {
		return nil, errors.New("dati delegato incompleti: nome, cognome e codice fiscale sono obbligatori")
	}

	reqID := req.ID
	return &PermanentDelegate{
		StudentID:         req.StudentID,
		SchoolID:          req.SchoolID,
		FirstName:         fn,
		LastName:          ln,
		TaxCode:           tc,
		Relationship:      rel,
		Phone:             phone,
		IDCardDetails:     card,
		IDCardFileURL:     cardURL,
		IsValid:           true,
		ApprovedRequestID: &reqID,
	}, nil
}
