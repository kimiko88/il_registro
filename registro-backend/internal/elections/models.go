package elections

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"time"
)

type SchoolElection struct {
	ID             string    `json:"id"`
	SchoolID       string    `json:"school_id"`
	Title          string    `json:"title"`
	ElectionTier   string    `json:"election_tier"` // classe, consiglio_istituto, consulta_studenti
	TargetRole     string    `json:"target_role"`   // parent, student, teacher, ata
	ClassID        *string   `json:"class_id,omitempty"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	MaxPreferences int       `json:"max_preferences"`
	IsClosed       bool      `json:"is_closed"`
	CreatedAt      time.Time `json:"created_at"`
}

type ElectionList struct {
	ID         string              `json:"id"`
	ElectionID string              `json:"election_id"`
	ListNumber int                 `json:"list_number"`
	Motto      string              `json:"motto"`
	Candidates []ElectionCandidate `json:"candidates,omitempty"`
}

type ElectionCandidate struct {
	ID             string `json:"id"`
	ListID         string `json:"list_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	CandidateOrder int    `json:"candidate_order"`
	VotesCount     int    `json:"votes_count,omitempty"`
}

type CastVotePayload struct {
	ElectionID   string   `json:"election_id"`
	ListID       *string  `json:"list_id,omitempty"`
	CandidateIDs []string `json:"candidate_ids,omitempty"`
	IsBlank      bool     `json:"is_blank"`
}

type VoteReceipt struct {
	ReceiptToken string    `json:"receipt_token"`
	ElectionID   string    `json:"election_id"`
	VotedAt      time.Time `json:"voted_at"`
}

type ListVoteCount struct {
	ListID     string `json:"list_id"`
	ListNumber int    `json:"list_number"`
	Motto      string `json:"motto"`
	TotalVotes int    `json:"total_votes"`
	SeatsWon   int    `json:"seats_won"`
}

type ScrutinyResult struct {
	ElectionID     string          `json:"election_id"`
	TotalVoters    int             `json:"total_voters"`
	TotalVotesCast int             `json:"total_votes_cast"`
	BlankVotes     int             `json:"blank_votes"`
	ListsResults   []ListVoteCount `json:"lists_results"`
	ElectedMembers []CandidateSeat `json:"elected_members"`
}

type CandidateSeat struct {
	CandidateID string `json:"candidate_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	ListID      string `json:"list_id"`
	ListMotto   string `json:"list_motto"`
	Votes       int    `json:"votes"`
}

func GenerateReceiptToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type quotientEntry struct {
	listID   string
	quotient float64
}

// AllocateSeatsDhondt allocates seats proportionally to election lists using the d'Hondt method
func AllocateSeatsDhondt(lists []ListVoteCount, totalSeats int) []ListVoteCount {
	if totalSeats <= 0 || len(lists) == 0 {
		return lists
	}

	resultMap := make(map[string]*ListVoteCount)
	for i := range lists {
		l := lists[i]
		l.SeatsWon = 0
		resultMap[l.ListID] = &l
	}

	var allQuotients []quotientEntry
	for _, l := range lists {
		if l.TotalVotes <= 0 {
			continue
		}
		for s := 1; s <= totalSeats; s++ {
			allQuotients = append(allQuotients, quotientEntry{
				listID:   l.ListID,
				quotient: float64(l.TotalVotes) / float64(s),
			})
		}
	}

	// Sort quotients descending
	sort.Slice(allQuotients, func(i, j int) bool {
		return allQuotients[i].quotient > allQuotients[j].quotient
	})

	seatsToAssign := totalSeats
	if len(allQuotients) < seatsToAssign {
		seatsToAssign = len(allQuotients)
	}

	for i := 0; i < seatsToAssign; i++ {
		listID := allQuotients[i].listID
		if item, ok := resultMap[listID]; ok {
			item.SeatsWon++
		}
	}

	var out []ListVoteCount
	for _, l := range lists {
		if item, ok := resultMap[l.ListID]; ok {
			out = append(out, *item)
		}
	}
	return out
}
