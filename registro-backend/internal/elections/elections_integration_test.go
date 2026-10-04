package elections

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockElectionRepo struct {
	elections map[string]*SchoolElection
	lists     map[string][]ElectionList
	voters    map[string]bool
	ballots   []CastVotePayload
}

func newMockElectionRepo() *mockElectionRepo {
	return &mockElectionRepo{
		elections: make(map[string]*SchoolElection),
		lists:     make(map[string][]ElectionList),
		voters:    make(map[string]bool),
	}
}

func (m *mockElectionRepo) CreateElection(ctx context.Context, e *SchoolElection) error {
	e.ID = "elec-1"
	e.CreatedAt = time.Now()
	m.elections[e.ID] = e
	return nil
}

func (m *mockElectionRepo) GetElection(ctx context.Context, id string) (*SchoolElection, error) {
	if e, ok := m.elections[id]; ok {
		return e, nil
	}
	return nil, nil
}

func (m *mockElectionRepo) ListElections(ctx context.Context, schoolID string) ([]SchoolElection, error) {
	var list []SchoolElection
	for _, e := range m.elections {
		list = append(list, *e)
	}
	return list, nil
}

func (m *mockElectionRepo) AddList(ctx context.Context, l *ElectionList) error {
	l.ID = "list-" + l.Motto
	m.lists[l.ElectionID] = append(m.lists[l.ElectionID], *l)
	return nil
}

func (m *mockElectionRepo) AddCandidate(ctx context.Context, c *ElectionCandidate) error {
	return nil
}

func (m *mockElectionRepo) GetListsWithCandidates(ctx context.Context, electionID string) ([]ElectionList, error) {
	return m.lists[electionID], nil
}

func (m *mockElectionRepo) HasVoted(ctx context.Context, electionID, voterID string) (bool, error) {
	return m.voters[electionID+"_"+voterID], nil
}

func (m *mockElectionRepo) CastVote(ctx context.Context, electionID, voterID, receiptToken string, payload CastVotePayload) error {
	key := electionID + "_" + voterID
	if m.voters[key] {
		return nil
	}
	m.voters[key] = true
	m.ballots = append(m.ballots, payload)
	return nil
}

func (m *mockElectionRepo) ComputeScrutiny(ctx context.Context, electionID string, totalSeats int) (*ScrutinyResult, error) {
	lists := []ListVoteCount{
		{ListID: "list-1", ListNumber: 1, Motto: "Lista Studenti", TotalVotes: 10},
	}
	allocated := AllocateSeatsDhondt(lists, totalSeats)
	return &ScrutinyResult{
		ElectionID:     electionID,
		TotalVoters:    len(m.voters),
		TotalVotesCast: len(m.ballots),
		ListsResults:   allocated,
	}, nil
}

func (m *mockElectionRepo) CloseElection(ctx context.Context, electionID string) error {
	if e, ok := m.elections[electionID]; ok {
		e.IsClosed = true
	}
	return nil
}

func setupElectionRouter(repo Repository, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "student")
		c.Set("school_id", "school-1")
		c.Set("user_id", userID)
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestElectionWorkflow(t *testing.T) {
	repo := newMockElectionRepo()
	router := setupElectionRouter(repo, "voter-student-1")

	// Step 1: Pre-populate election in repo
	election := &SchoolElection{
		ID:             "elec-1",
		SchoolID:       "school-1",
		Title:          "Elezioni Rappresentanti di Classe",
		ElectionTier:   "classe",
		TargetRole:     "student",
		StartTime:      time.Now().Add(-1 * time.Hour),
		EndTime:        time.Now().Add(24 * time.Hour),
		MaxPreferences: 1,
		IsClosed:       false,
	}
	repo.elections[election.ID] = election
	listID := "list-1"
	repo.lists[election.ID] = []ElectionList{
		{ID: listID, ElectionID: election.ID, ListNumber: 1, Motto: "Lista Studenti"},
	}

	// Step 2: Cast vote
	votePayload := map[string]interface{}{
		"election_id":   election.ID,
		"list_id":       listID,
		"candidate_ids": []string{},
		"is_blank":      false,
	}
	body, _ := json.Marshal(votePayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/elections/elec-1/vote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on cast vote, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Receipt VoteReceipt `json:"receipt"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Receipt.ReceiptToken) != 32 {
		t.Errorf("expected 32-char anonymous receipt token, got %s", res.Receipt.ReceiptToken)
	}

	// Step 3: Second vote by same voter should fail
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/elections/elec-1/vote", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on duplicate vote, got %d", w2.Code)
	}

	// Step 4: Scrutiny
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/elections/elec-1/scrutiny?seats=2", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200 on scrutiny, got %d", w3.Code)
	}
}
