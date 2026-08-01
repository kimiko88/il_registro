package communications

import (
	"context"
	"errors"
	"fmt"
	"time"

	"registro-backend/internal/users"
)

type Service struct {
	repo     Repository
	userRepo users.Repository
}

func NewService(repo Repository, uRepo ...users.Repository) *Service {
	if repo == nil {
		panic("communications.NewService: repo must not be nil")
	}
	var userRepo users.Repository
	if len(uRepo) > 0 {
		userRepo = uRepo[0]
	}
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) SendMessage(ctx context.Context, actorRole, schoolID, senderID string, req CreateMessageRequest) (*Message, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" && actorRole != "secretary" && actorRole != "principal" && actorRole != "vice_principal" {
		return nil, errors.New("unauthorized: non hai i permessi per inviare comunicazioni")
	}
	if req.Subject == "" || req.Body == "" {
		return nil, errors.New("subject and body are required")
	}
	if req.Type != "bacheca" && len(req.Recipients) == 0 {
		return nil, errors.New("recipients are required for targeted messages")
	}
	if req.Type == "bacheca" && req.RequiresSignature && len(req.Recipients) == 0 {
		return nil, errors.New("cannot set RequiresSignature on a board message without specific recipients")
	}

	targetSchoolID := schoolID
	if req.SchoolID != nil && *req.SchoolID != "" {
		if actorRole != "superadmin" && *req.SchoolID != schoolID {
			return nil, errors.New("unauthorized: non puoi inviare messaggi a nome di un'altra scuola")
		}
		targetSchoolID = *req.SchoolID
	}

	if s.userRepo != nil && len(req.Recipients) > 0 && targetSchoolID != "" {
		for _, recipientID := range req.Recipients {
			u, err := s.userRepo.GetByID(ctx, recipientID)
			if err != nil || u == nil || u.SchoolID == nil || *u.SchoolID != targetSchoolID {
				return nil, fmt.Errorf("unauthorized: destinatario %s non appartiene alla scuola", recipientID)
			}
		}
	}

	msg := &Message{
		SchoolID:          req.SchoolID,
		SenderID:          senderID,
		ReceiverIDs:       req.Recipients,
		Subject:           req.Subject,
		Body:              req.Body,
		AttachmentURL:     req.AttachmentURL,
		Type:              req.Type,
		RequiresSignature: req.RequiresSignature,
	}

	if req.SignatureDeadline != nil && *req.SignatureDeadline != "" {
		if d, err := time.Parse("2006-01-02", *req.SignatureDeadline); err == nil {
			msg.SignatureDeadline = &d
		} else if d, err := time.Parse(time.RFC3339, *req.SignatureDeadline); err == nil {
			msg.SignatureDeadline = &d
		} else {
			return nil, errors.New("formato data scadenza firma non valido (richiesto YYYY-MM-DD o RFC3339)")
		}
	}

	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *Service) ListMessages(ctx context.Context, userID string) ([]*Message, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) ListBacheca(ctx context.Context, schoolID, userID string) ([]*Message, error) {
	return s.repo.ListBacheca(ctx, schoolID, userID)
}

func (s *Service) DeleteMessage(ctx context.Context, actorID, actorRole, schoolID, id string) error {
	msg, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if msg.SenderID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: cannot delete message of another user")
	}
	if actorRole != "superadmin" {
		if msg.SchoolID == nil || schoolID == "" || *msg.SchoolID != schoolID {
			return errors.New("unauthorized: cannot delete message of another school")
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) SignMessage(ctx context.Context, userID, messageID, ipAddress string) error {
	msg, err := s.repo.Get(ctx, messageID)
	if err != nil {
		return fmt.Errorf("messaggio non trovato: %w", err)
	}

	isRecipient := false
	for _, rid := range msg.ReceiverIDs {
		if rid == userID {
			isRecipient = true
			break
		}
	}
	if !isRecipient {
		return errors.New("unauthorized: solo i destinatari possono firmare questo messaggio")
	}

	return s.repo.SignWithIP(ctx, messageID, userID, ipAddress)
}

func (s *Service) SignMessageWithIP(ctx context.Context, messageID, userID, ipAddress string) error {
	return s.SignMessage(ctx, userID, messageID, ipAddress)
}

func (s *Service) GetMessageSignatures(ctx context.Context, communicationID string) ([]string, error) {
	return s.repo.GetSignatures(ctx, communicationID)
}

func (s *Service) GetSignatureReport(ctx context.Context, actorRole, communicationID string) (*SignatureReportResponse, error) {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" {
		return nil, errors.New("forbidden: signature reports are restricted to staff")
	}
	return s.repo.GetSignatureReport(ctx, communicationID)
}

// GetMessageByID recupera un messaggio verificando i permessi di accesso.
// Bug 112: la verifica schoolID è obbligatoria indipendentemente dal tipo di messaggio.
func (s *Service) GetMessageByID(ctx context.Context, userID, schoolID, userRole, messageID string) (*Message, error) {
	msg, err := s.repo.Get(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// Verifica appartenenza alla scuola per TUTTI i tipi di messaggio (incluso bacheca)
	if userRole != "superadmin" {
		if msg.SchoolID == nil || *msg.SchoolID == "" || *msg.SchoolID != schoolID {
			return nil, errors.New("unauthorized: messaggio non appartiene alla tua scuola")
		}
	}

	// Admin e superadmin possono vedere tutti i messaggi della loro scuola
	if userRole == "admin" || userRole == "superadmin" {
		return msg, nil
	}

	// Verifica accesso individuale: deve essere il mittente o un destinatario
	if msg.SenderID == userID {
		return msg, nil
	}
	for _, rid := range msg.ReceiverIDs {
		if rid == userID {
			return msg, nil
		}
	}

	// Per messaggi di tipo bacheca senza destinatari specifici, controlla schoolID (già verificato sopra)
	if msg.Type == "bacheca" && len(msg.ReceiverIDs) == 0 {
		return msg, nil
	}

	return nil, errors.New("unauthorized: non hai accesso a questo messaggio")
}

// Bug 113: UpdateMessage già blocca i messaggi firmati.
// NOTA: il cambio di tipo di messaggio dopo la consegna non è bloccato a livello di servizio.
func (s *Service) UpdateMessage(ctx context.Context, actorID, actorRole, schoolID, id, subject, body string) error {
	if subject == "" || body == "" {
		return errors.New("subject and body are required")
	}
	msg, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if msg.SenderID != actorID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("unauthorized: cannot edit message of another user")
	}
	if actorRole != "superadmin" {
		if msg.SchoolID == nil || schoolID == "" || *msg.SchoolID != schoolID {
			return errors.New("unauthorized: cannot edit message of another school")
		}
	}
	sigs, err := s.repo.GetSignatures(ctx, id)
	if err == nil && len(sigs) > 0 {
		return errors.New("impossibile modificare un messaggio che contiene già firme digitali")
	}
	return s.repo.Update(ctx, id, subject, body)
}

// MarkAsRead marca un messaggio come letto.
// Bug 114: verifica che l'utente sia un destinatario prima di aggiornare lo stato.
func (s *Service) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	msg, err := s.repo.Get(ctx, communicationID)
	if err != nil {
		return fmt.Errorf("messaggio non trovato: %w", err)
	}

	if msg.Type == "bacheca" && len(msg.ReceiverIDs) == 0 {
		return s.repo.MarkAsRead(ctx, communicationID, userID, ipAddress)
	}

	for _, rid := range msg.ReceiverIDs {
		if rid == userID {
			return s.repo.MarkAsRead(ctx, communicationID, userID, ipAddress)
		}
	}
	return errors.New("unauthorized: non sei un destinatario di questo messaggio")
}

func (s *Service) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	return s.repo.GetUnreadUsers(ctx, communicationID)
}

func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *Service) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*Message, error) {
	return s.repo.ListCircolari(ctx, schoolID, userID, year)
}
