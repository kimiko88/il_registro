package elections

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateElection(ctx context.Context, e *SchoolElection) error {
	if e.Title == "" || e.SchoolID == "" {
		return errors.New("titolo ed istituto scolastico sono obbligatori")
	}
	if e.StartTime.IsZero() {
		e.StartTime = time.Now()
	}
	if e.EndTime.IsZero() {
		e.EndTime = e.StartTime.Add(48 * time.Hour)
	}
	if e.MaxPreferences <= 0 {
		e.MaxPreferences = 1
	}
	return s.repo.CreateElection(ctx, e)
}

func (s *Service) ListElections(ctx context.Context, schoolID string) ([]SchoolElection, error) {
	if schoolID == "" {
		return nil, errors.New("school_id obbligatorio")
	}
	return s.repo.ListElections(ctx, schoolID)
}

func (s *Service) GetElectionDetails(ctx context.Context, electionID string) (*SchoolElection, []ElectionList, error) {
	e, err := s.repo.GetElection(ctx, electionID)
	if err != nil {
		return nil, nil, err
	}
	lists, err := s.repo.GetListsWithCandidates(ctx, electionID)
	if err != nil {
		return nil, nil, err
	}
	return e, lists, nil
}

func (s *Service) CastVote(ctx context.Context, electionID, voterID string, payload CastVotePayload) (*VoteReceipt, error) {
	if electionID == "" || voterID == "" {
		return nil, errors.New("elezione e identificativo elettore sono obbligatori")
	}

	election, err := s.repo.GetElection(ctx, electionID)
	if err != nil {
		return nil, errors.New("elezione non trovata")
	}

	if election.IsClosed {
		return nil, errors.New("il seggio per questa elezione è chiuso")
	}

	now := time.Now()
	if now.Before(election.StartTime) || now.After(election.EndTime) {
		return nil, errors.New("votazione non consentita al di fuori degli orari di apertura del seggio")
	}

	hasVoted, err := s.repo.HasVoted(ctx, electionID, voterID)
	if err != nil {
		return nil, err
	}
	if hasVoted {
		return nil, errors.New("hai già espresso il tuo voto per questa elezione")
	}

	if !payload.IsBlank && len(payload.CandidateIDs) > election.MaxPreferences {
		return nil, errors.New("numero di preferenze espresse superiore al limite consentito")
	}

	receiptToken := GenerateReceiptToken()
	payload.ElectionID = electionID

	if err := s.repo.CastVote(ctx, electionID, voterID, receiptToken, payload); err != nil {
		return nil, err
	}

	return &VoteReceipt{
		ReceiptToken: receiptToken,
		ElectionID:   electionID,
		VotedAt:      now,
	}, nil
}

func (s *Service) GetScrutiny(ctx context.Context, electionID string, totalSeats int) (*ScrutinyResult, error) {
	if electionID == "" {
		return nil, errors.New("election_id obbligatorio")
	}
	if totalSeats <= 0 {
		totalSeats = 4 // Default seggi consiglio o classe
	}
	return s.repo.ComputeScrutiny(ctx, electionID, totalSeats)
}

func (s *Service) CloseElection(ctx context.Context, electionID string) error {
	return s.repo.CloseElection(ctx, electionID)
}
