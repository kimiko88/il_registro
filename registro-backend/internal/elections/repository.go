package elections

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Repository interface {
	CreateElection(ctx context.Context, e *SchoolElection) error
	GetElection(ctx context.Context, id string) (*SchoolElection, error)
	ListElections(ctx context.Context, schoolID string) ([]SchoolElection, error)
	AddList(ctx context.Context, l *ElectionList) error
	AddCandidate(ctx context.Context, c *ElectionCandidate) error
	GetListsWithCandidates(ctx context.Context, electionID string) ([]ElectionList, error)
	HasVoted(ctx context.Context, electionID, voterID string) (bool, error)
	CastVote(ctx context.Context, electionID, voterID, receiptToken string, payload CastVotePayload) error
	ComputeScrutiny(ctx context.Context, electionID string, totalSeats int) (*ScrutinyResult, error)
	CloseElection(ctx context.Context, electionID string) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateElection(ctx context.Context, e *SchoolElection) error {
	query := `
		INSERT INTO school_elections (
			school_id, title, election_tier, target_role, class_id, start_time, end_time, max_preferences
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at;
	`
	return r.db.QueryRowContext(ctx, query,
		e.SchoolID, e.Title, e.ElectionTier, e.TargetRole, e.ClassID, e.StartTime, e.EndTime, e.MaxPreferences,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *postgresRepo) GetElection(ctx context.Context, id string) (*SchoolElection, error) {
	query := `
		SELECT id, school_id, title, election_tier, target_role, class_id, start_time, end_time, max_preferences, is_closed, created_at
		FROM school_elections
		WHERE id = $1;
	`
	var e SchoolElection
	var classID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.SchoolID, &e.Title, &e.ElectionTier, &e.TargetRole, &classID, &e.StartTime, &e.EndTime, &e.MaxPreferences, &e.IsClosed, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if classID.Valid {
		e.ClassID = &classID.String
	}
	return &e, nil
}

func (r *postgresRepo) ListElections(ctx context.Context, schoolID string) ([]SchoolElection, error) {
	query := `
		SELECT id, school_id, title, election_tier, target_role, class_id, start_time, end_time, max_preferences, is_closed, created_at
		FROM school_elections
		WHERE school_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SchoolElection
	for rows.Next() {
		var e SchoolElection
		var classID sql.NullString
		if err := rows.Scan(&e.ID, &e.SchoolID, &e.Title, &e.ElectionTier, &e.TargetRole, &classID, &e.StartTime, &e.EndTime, &e.MaxPreferences, &e.IsClosed, &e.CreatedAt); err != nil {
			return nil, err
		}
		if classID.Valid {
			e.ClassID = &classID.String
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *postgresRepo) AddList(ctx context.Context, l *ElectionList) error {
	query := `
		INSERT INTO election_lists (election_id, list_number, motto)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	return r.db.QueryRowContext(ctx, query, l.ElectionID, l.ListNumber, l.Motto).Scan(&l.ID)
}

func (r *postgresRepo) AddCandidate(ctx context.Context, c *ElectionCandidate) error {
	query := `
		INSERT INTO election_candidates (list_id, first_name, last_name, candidate_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	return r.db.QueryRowContext(ctx, query, c.ListID, c.FirstName, c.LastName, c.CandidateOrder).Scan(&c.ID)
}

func (r *postgresRepo) GetListsWithCandidates(ctx context.Context, electionID string) ([]ElectionList, error) {
	query := `
		SELECT id, election_id, list_number, motto
		FROM election_lists
		WHERE election_id = $1
		ORDER BY list_number ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, electionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []ElectionList
	for rows.Next() {
		var l ElectionList
		if err := rows.Scan(&l.ID, &l.ElectionID, &l.ListNumber, &l.Motto); err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}

	for i := range lists {
		cQuery := `
			SELECT id, list_id, first_name, last_name, candidate_order
			FROM election_candidates
			WHERE list_id = $1
			ORDER BY candidate_order ASC;
		`
		cRows, err := r.db.QueryContext(ctx, cQuery, lists[i].ID)
		if err == nil {
			for cRows.Next() {
				var c ElectionCandidate
				if err := cRows.Scan(&c.ID, &c.ListID, &c.FirstName, &c.LastName, &c.CandidateOrder); err == nil {
					lists[i].Candidates = append(lists[i].Candidates, c)
				}
			}
			cRows.Close()
		}
	}
	return lists, nil
}

func (r *postgresRepo) HasVoted(ctx context.Context, electionID, voterID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM election_voter_registry WHERE election_id = $1 AND voter_id = $2);`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, electionID, voterID).Scan(&exists)
	return exists, err
}

func (r *postgresRepo) CastVote(ctx context.Context, electionID, voterID, receiptToken string, payload CastVotePayload) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Mark voter in registry
	voterQuery := `
		INSERT INTO election_voter_registry (election_id, voter_id, receipt_token)
		VALUES ($1, $2, $3);
	`
	if _, err := tx.ExecContext(ctx, voterQuery, electionID, voterID, receiptToken); err != nil {
		return errors.New("l'elettore ha già espresso il proprio voto per questa elezione")
	}

	// 2. Insert into anonymous ballot box
	boxQuery := `
		INSERT INTO election_ballot_box (election_id, list_id, candidate_ids, is_blank)
		VALUES ($1, $2, $3, $4);
	`
	if _, err := tx.ExecContext(ctx, boxQuery, electionID, payload.ListID, pq.Array(payload.CandidateIDs), payload.IsBlank); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *postgresRepo) ComputeScrutiny(ctx context.Context, electionID string, totalSeats int) (*ScrutinyResult, error) {
	lists, err := r.GetListsWithCandidates(ctx, electionID)
	if err != nil {
		return nil, err
	}

	var totalVoters int
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM election_voter_registry WHERE election_id = $1", electionID).Scan(&totalVoters)

	var totalVotesCast, blankVotes int
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM election_ballot_box WHERE election_id = $1", electionID).Scan(&totalVotesCast)
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM election_ballot_box WHERE election_id = $1 AND is_blank = true", electionID).Scan(&blankVotes)

	var listVoteCounts []ListVoteCount
	for _, l := range lists {
		var votes int
		_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM election_ballot_box WHERE election_id = $1 AND list_id = $2", electionID, l.ID).Scan(&votes)
		listVoteCounts = append(listVoteCounts, ListVoteCount{
			ListID:     l.ID,
			ListNumber: l.ListNumber,
			Motto:      l.Motto,
			TotalVotes: votes,
		})
	}

	allocatedLists := AllocateSeatsDhondt(listVoteCounts, totalSeats)

	var elected []CandidateSeat
	for _, l := range allocatedLists {
		if l.SeatsWon <= 0 {
			continue
		}
		// Count candidate preferences
		for _, cand := range lists {
			if cand.ID != l.ListID {
				continue
			}
			for _, c := range cand.Candidates {
				var cVotes int
				_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM election_ballot_box WHERE election_id = $1 AND $2 = ANY(candidate_ids)", electionID, c.ID).Scan(&cVotes)
				elected = append(elected, CandidateSeat{
					CandidateID: c.ID,
					FirstName:   c.FirstName,
					LastName:    c.LastName,
					ListID:      l.ListID,
					ListMotto:   l.Motto,
					Votes:       cVotes,
				})
			}
		}
	}

	return &ScrutinyResult{
		ElectionID:     electionID,
		TotalVoters:    totalVoters,
		TotalVotesCast: totalVotesCast,
		BlankVotes:     blankVotes,
		ListsResults:   allocatedLists,
		ElectedMembers: elected,
	}, nil
}

func (r *postgresRepo) CloseElection(ctx context.Context, electionID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE school_elections SET is_closed = true WHERE id = $1", electionID)
	return err
}
