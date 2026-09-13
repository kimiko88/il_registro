package strike

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStrikeRepo struct {
	mock.Mock
}

func (m *mockStrikeRepo) CreateNotice(ctx context.Context, n *StrikeNotice) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *mockStrikeRepo) GetNoticeByID(ctx context.Context, id, schoolID string) (*StrikeNotice, error) {
	args := m.Called(ctx, id, schoolID)
	if res := args.Get(0); res != nil {
		return res.(*StrikeNotice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStrikeRepo) ListNotices(ctx context.Context, schoolID string) ([]*StrikeNotice, error) {
	args := m.Called(ctx, schoolID)
	if res := args.Get(0); res != nil {
		return res.([]*StrikeNotice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStrikeRepo) DeleteNotice(ctx context.Context, id, schoolID string) error {
	args := m.Called(ctx, id, schoolID)
	return args.Error(0)
}

func (m *mockStrikeRepo) UpsertDeclaration(ctx context.Context, d *StrikeDeclaration) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func (m *mockStrikeRepo) GetDeclaration(ctx context.Context, noticeID, userID string) (*StrikeDeclaration, error) {
	args := m.Called(ctx, noticeID, userID)
	if res := args.Get(0); res != nil {
		return res.(*StrikeDeclaration), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStrikeRepo) GetNoticeSummary(ctx context.Context, noticeID, schoolID string) (*StrikeNoticeSummaryResponse, error) {
	args := m.Called(ctx, noticeID, schoolID)
	if res := args.Get(0); res != nil {
		return res.(*StrikeNoticeSummaryResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStrikeRepo) CreateBachecaCommunication(ctx context.Context, schoolID, senderID, title, content string, deadline time.Time) (string, error) {
	args := m.Called(ctx, schoolID, senderID, title, content, deadline)
	return args.String(0), args.Error(1)
}

func (m *mockStrikeRepo) ResolveSchoolID(ctx context.Context, userID string) string {
	args := m.Called(ctx, userID)
	return args.String(0)
}

func TestStrikeService_CreateNotice(t *testing.T) {
	repo := new(mockStrikeRepo)
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("Unauthorized role cannot create strike notice", func(t *testing.T) {
		req := CreateStrikeNoticeRequest{
			Title:               "Sciopero",
			ProclaimedBy:        "Sindacato",
			StrikeDate:          "2026-10-20",
			DeclarationDeadline: "2026-10-18T12:00:00Z",
			Content:             "Test",
		}
		_, err := svc.CreateNotice(ctx, "user-1", "teacher", "school-1", req)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("DSGA can create strike notice", func(t *testing.T) {
		req := CreateStrikeNoticeRequest{
			Title:               "Sciopero Nazionale",
			ProclaimedBy:        "FLC CGIL",
			StrikeDate:          "2026-10-25",
			DeclarationDeadline: "2026-10-23T14:00:00Z",
			Content:             "Comunicazione ufficiale",
			PublishToBacheca:    true,
		}

		repo.On("CreateBachecaCommunication", ctx, "school-1", "user-dsga", req.Title, req.Content, mock.Anything).Return("comm-123", nil).Once()
		repo.On("CreateNotice", ctx, mock.MatchedBy(func(n *StrikeNotice) bool {
			return n.Title == "Sciopero Nazionale" && n.SchoolID == "school-1" && *n.CommunicationID == "comm-123"
		})).Return(nil).Once()

		notice, err := svc.CreateNotice(ctx, "user-dsga", "dsga", "school-1", req)
		assert.NoError(t, err)
		assert.NotNil(t, notice)
		assert.Equal(t, "Sciopero Nazionale", notice.Title)
		assert.Equal(t, "comm-123", *notice.CommunicationID)
	})
}

func TestStrikeService_SubmitDeclaration(t *testing.T) {
	repo := new(mockStrikeRepo)
	svc := NewService(repo)
	ctx := context.Background()

	t.Run("Invalid intention fails", func(t *testing.T) {
		req := SubmitDeclarationRequest{
			Intention: "maybe",
		}
		_, err := svc.SubmitDeclaration(ctx, "notice-1", "school-1", "user-teacher", "127.0.0.1", req)
		assert.ErrorIs(t, err, ErrInvalidIntention)
	})

	t.Run("Declaration blocked after deadline", func(t *testing.T) {
		pastNotice := &StrikeNotice{
			ID:                  "notice-past",
			SchoolID:            "school-1",
			DeclarationDeadline: time.Now().Add(-1 * time.Hour), // expired 1 hour ago
		}
		repo.On("GetNoticeByID", ctx, "notice-past", "school-1").Return(pastNotice, nil).Once()

		req := SubmitDeclarationRequest{
			Intention: IntentionParticipates,
		}
		_, err := svc.SubmitDeclaration(ctx, "notice-past", "school-1", "user-teacher", "127.0.0.1", req)
		assert.ErrorIs(t, err, ErrDeclarationDeadlinePassed)
	})

	t.Run("Valid declaration before deadline succeeds", func(t *testing.T) {
		activeNotice := &StrikeNotice{
			ID:                  "notice-active",
			SchoolID:            "school-1",
			DeclarationDeadline: time.Now().Add(48 * time.Hour),
		}
		repo.On("GetNoticeByID", ctx, "notice-active", "school-1").Return(activeNotice, nil).Once()
		repo.On("UpsertDeclaration", ctx, mock.MatchedBy(func(d *StrikeDeclaration) bool {
			return d.StrikeNoticeID == "notice-active" && d.UserID == "user-teacher" && d.Intention == IntentionParticipates
		})).Return(nil).Once()

		req := SubmitDeclarationRequest{
			Intention: IntentionParticipates,
			Notes:     "Aderisco",
		}
		decl, err := svc.SubmitDeclaration(ctx, "notice-active", "school-1", "user-teacher", "192.168.1.50", req)
		assert.NoError(t, err)
		assert.NotNil(t, decl)
		assert.Equal(t, IntentionParticipates, decl.Intention)
	})
}
