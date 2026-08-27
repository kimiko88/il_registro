package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/colloqui"
	"registro-backend/internal/communications"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Colloqui_CreateSlot_PastDate_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockColloquiRepo{}
	svc := colloqui.NewService(mockRepo)
	handler := colloqui.NewHandler(svc)

	r := gin.New()
	r.POST("/colloqui/slots", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.CreateSlot(c)
	})

	body := colloqui.CreateSlotRequest{
		Date:        "2020-01-01",
		StartTime:   "10:00",
		EndTime:     "11:00",
		MaxBookings: 5,
		Type:        "ordinario",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/colloqui/slots", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIntegration_Communications_SendMessage_MissingSubject_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockCommsRepo{}
	mockUsers := &mockUserRepoForComms{}
	svc := communications.NewService(mockRepo, mockUsers)
	handler := communications.NewHandler(svc)

	r := gin.New()
	r.POST("/communications", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Send(c)
	})

	body := communications.CreateMessageRequest{
		Subject:    "",
		Body:       "Test body",
		Type:       "circolare",
		Recipients: []string{"user-2"},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/communications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIntegration_Communications_SSRF_Protection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockCommsRepo{}
	mockUsers := &mockUserRepoForComms{}
	svc := communications.NewService(mockRepo, mockUsers)
	handler := communications.NewHandler(svc)

	r := gin.New()
	r.POST("/communications", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Send(c)
	})

	invalidURL := "http://169.254.169.254/latest/meta-data/"
	body := communications.CreateMessageRequest{
		Subject:       "Oggetto",
		Body:          "Test body",
		Type:          "circolare",
		Recipients:    []string{"user-2"},
		AttachmentURL: &invalidURL,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/communications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type mockColloquiRepo struct{}

func (m *mockColloquiRepo) CreateSlot(ctx context.Context, slot *colloqui.ColloquioSlot) error {
	slot.ID = "slot-1"
	return nil
}
func (m *mockColloquiRepo) GetSlotByID(ctx context.Context, id string) (*colloqui.ColloquioSlot, error) {
	return &colloqui.ColloquioSlot{ID: id, SchoolID: "school-1"}, nil
}
func (m *mockColloquiRepo) ListSlots(ctx context.Context, filter colloqui.SlotFilter) ([]*colloqui.ColloquioSlot, error) {
	return []*colloqui.ColloquioSlot{}, nil
}
func (m *mockColloquiRepo) PatchSlot(ctx context.Context, slotID string, startTime, endTime string) error {
	return nil
}
func (m *mockColloquiRepo) CancelSlot(ctx context.Context, id string) error { return nil }
func (m *mockColloquiRepo) CreateBooking(ctx context.Context, b *colloqui.ColloquioBooking) error {
	return nil
}
func (m *mockColloquiRepo) GetBookingByID(ctx context.Context, id string) (*colloqui.ColloquioBooking, error) {
	return &colloqui.ColloquioBooking{ID: id}, nil
}
func (m *mockColloquiRepo) ListUserBookings(ctx context.Context, userID, schoolID string) ([]*colloqui.ColloquioBooking, error) {
	return []*colloqui.ColloquioBooking{}, nil
}
func (m *mockColloquiRepo) ListSlotBookings(ctx context.Context, slotID string) ([]*colloqui.ColloquioBooking, error) {
	return []*colloqui.ColloquioBooking{}, nil
}
func (m *mockColloquiRepo) UpdateBookingStatus(ctx context.Context, bookingID string, status colloqui.BookingStatus, changedBy string, reason string) error {
	return nil
}
func (m *mockColloquiRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return "teacher-prof-1", nil
}
func (m *mockColloquiRepo) GetParentProfileID(ctx context.Context, userID string) (string, error) {
	return "parent-prof-1", nil
}
func (m *mockColloquiRepo) GetStudentProfileID(ctx context.Context, userID string) (string, error) {
	return "student-prof-1", nil
}
func (m *mockColloquiRepo) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	return true, nil
}
func (m *mockColloquiRepo) ExistsOverlappingSlot(ctx context.Context, teacherID, date, start, end string) (bool, error) {
	return false, nil
}
func (m *mockColloquiRepo) CreateGeneralMeeting(ctx context.Context, gm *colloqui.GeneralParentMeeting, teacherIDs []string) error {
	return nil
}
func (m *mockColloquiRepo) ListGeneralMeetings(ctx context.Context, schoolID string) ([]colloqui.GeneralParentMeeting, error) {
	return []colloqui.GeneralParentMeeting{}, nil
}
func (m *mockColloquiRepo) GetGeneralMeeting(ctx context.Context, id string) (*colloqui.GeneralParentMeeting, error) {
	return &colloqui.GeneralParentMeeting{ID: id}, nil
}
func (m *mockColloquiRepo) BookQueueTicket(ctx context.Context, t *colloqui.GeneralMeetingQueueTicket) error {
	return nil
}
func (m *mockColloquiRepo) ListQueueTickets(ctx context.Context, meetingID, teacherID, parentID string) ([]colloqui.GeneralMeetingQueueTicket, error) {
	return []colloqui.GeneralMeetingQueueTicket{}, nil
}
func (m *mockColloquiRepo) UpdateTicketStatus(ctx context.Context, id, status, notes string) error {
	return nil
}

type mockCommsRepo struct{}

func (m *mockCommsRepo) Create(ctx context.Context, msg *communications.Message) error {
	msg.ID = "msg-1"
	return nil
}
func (m *mockCommsRepo) List(ctx context.Context, userID, schoolID string) ([]*communications.Message, error) {
	return []*communications.Message{}, nil
}
func (m *mockCommsRepo) ListBacheca(ctx context.Context, schoolID, userID string) ([]*communications.Message, error) {
	return []*communications.Message{}, nil
}
func (m *mockCommsRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockCommsRepo) Sign(ctx context.Context, communicationID string, userID string) error {
	return nil
}
func (m *mockCommsRepo) SignWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error {
	return nil
}
func (m *mockCommsRepo) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	return []string{}, nil
}
func (m *mockCommsRepo) GetSignatureReport(ctx context.Context, communicationID string) (*communications.SignatureReportResponse, error) {
	return &communications.SignatureReportResponse{}, nil
}
func (m *mockCommsRepo) Get(ctx context.Context, id string) (*communications.Message, error) {
	return &communications.Message{ID: id}, nil
}
func (m *mockCommsRepo) Update(ctx context.Context, id string, subject, body string) error {
	return nil
}
func (m *mockCommsRepo) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	return nil
}
func (m *mockCommsRepo) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	return []string{}, nil
}
func (m *mockCommsRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *mockCommsRepo) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*communications.Message, error) {
	return []*communications.Message{}, nil
}
func (m *mockCommsRepo) Ack(ctx context.Context, communicationID, userID string) error {
	return nil
}

type mockUserRepoForComms struct{}

func (m *mockUserRepoForComms) Create(ctx context.Context, user *users.User) error { return nil }
func (m *mockUserRepoForComms) GetByID(ctx context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{ID: id, SchoolID: &schoolID}, nil
}
func (m *mockUserRepoForComms) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) Update(ctx context.Context, user *users.User) error { return nil }
func (m *mockUserRepoForComms) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	return nil
}
func (m *mockUserRepoForComms) Delete(ctx context.Context, id string) error  { return nil }
func (m *mockUserRepoForComms) Restore(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForComms) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForComms) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	schoolID := "school-1"
	var res []users.User
	for _, id := range ids {
		res = append(res, users.User{ID: id, SchoolID: &schoolID})
	}
	return res, nil
}
func (m *mockUserRepoForComms) LogAudit(ctx context.Context, log *users.AuditLog) error { return nil }
func (m *mockUserRepoForComms) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForComms) BulkCreate(ctx context.Context, users []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *mockUserRepoForComms) BulkDelete(ctx context.Context, ids []string) (int, error) {
	return 0, nil
}
func (m *mockUserRepoForComms) HardDelete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForComms) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return nil
}
func (m *mockUserRepoForComms) ClearTempMFASecret(ctx context.Context, userID string) error {
	return nil
}
func (m *mockUserRepoForComms) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForComms) GetParentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForComms) AddGuardian(ctx context.Context, studentProfileID, parentProfileID, relationship string) error {
	return nil
}
func (m *mockUserRepoForComms) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	return nil
}
func (m *mockUserRepoForComms) GetGuardians(ctx context.Context, studentProfileID string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *mockUserRepoForComms) IsActive(ctx context.Context, id string) (bool, error) {
	return true, nil
}
func (m *mockUserRepoForComms) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	return true, nil
}
