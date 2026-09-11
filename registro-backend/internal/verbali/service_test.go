package verbali

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateMeeting(ctx context.Context, meeting *CouncilMeeting) error {
	args := m.Called(ctx, meeting)
	return args.Error(0)
}
func (m *MockRepository) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	args := m.Called(ctx, schoolID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*CouncilMeeting), args.Error(1)
}
func (m *MockRepository) GetMeetingByID(ctx context.Context, id string) (*CouncilMeeting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CouncilMeeting), args.Error(1)
}
func (m *MockRepository) GetMeetingWithCoordinator(ctx context.Context, meetingID string) (*CouncilMeeting, string, error) {
	args := m.Called(ctx, meetingID)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*CouncilMeeting), args.String(1), args.Error(2)
}
func (m *MockRepository) CreateVerbale(ctx context.Context, v *MeetingVerbale) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}
func (m *MockRepository) GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MeetingVerbale), args.Error(1)
}
func (m *MockRepository) UpdateVerbale(ctx context.Context, v *MeetingVerbale) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}
func (m *MockRepository) DeleteVerbale(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) ListVerbali(ctx context.Context, meetingID string, userID string, onlySigned bool) ([]*MeetingVerbale, error) {
	args := m.Called(ctx, meetingID, userID, onlySigned)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MeetingVerbale), args.Error(1)
}
func (m *MockRepository) ListAllVerbali(ctx context.Context, schoolID, classID, userID string, onlySigned bool) ([]*MeetingVerbale, error) {
	args := m.Called(ctx, schoolID, classID, userID, onlySigned)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MeetingVerbale), args.Error(1)
}
func (m *MockRepository) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	args := m.Called(ctx, verbaleID, userID, ipAddress)
	return args.Error(0)
}
func (m *MockRepository) MarkVerbaleSigned(ctx context.Context, verbaleID string) error {
	args := m.Called(ctx, verbaleID)
	return args.Error(0)
}
func (m *MockRepository) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	args := m.Called(ctx, verbaleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]VerbaleSignature), args.Error(1)
}
func (m *MockRepository) ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	args := m.Called(ctx, classID, schoolID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepository) CreateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockRepository) ListTemplates(ctx context.Context, schoolID, meetingType string) ([]*MeetingVerbaleTemplate, error) {
	args := m.Called(ctx, schoolID, meetingType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MeetingVerbaleTemplate), args.Error(1)
}
func (m *MockRepository) GetTemplateByID(ctx context.Context, id string) (*MeetingVerbaleTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MeetingVerbaleTemplate), args.Error(1)
}
func (m *MockRepository) UpdateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockRepository) DeleteTemplate(ctx context.Context, id, schoolID string) error {
	args := m.Called(ctx, id, schoolID)
	return args.Error(0)
}

func TestCreateMeetingAndVerbale(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	reqMeeting := CreateMeetingRequest{
		ClassID:     "class-1",
		MeetingType: "consiglio_classe",
		Title:       "Consiglio di Classe di Novembre",
		Date:        "2025-11-10",
		StartTime:   "16:00",
		EndTime:     "17:30",
		Agenda:      "Approvazione piano didattico",
	}

	mockRepo.On("ClassBelongsToSchool", mock.Anything, "class-1", "school-1").Return(true, nil).Once()
	mockRepo.On("CreateMeeting", mock.Anything, mock.MatchedBy(func(m *CouncilMeeting) bool {
		return m.Title == "Consiglio di Classe di Novembre"
	})).Return(nil).Once()

	m, err := svc.CreateMeeting(context.Background(), "t-1", "school-1", reqMeeting)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	mockRepo.On("GetMeetingWithCoordinator", mock.Anything, "m-1").Return(m, "coord-1", nil).Once()
	mockRepo.On("CreateVerbale", mock.Anything, mock.MatchedBy(func(v *MeetingVerbale) bool {
		return v.Title == "Verbale n.1"
	})).Return(nil).Once()

	secID := "t-1"
	v, err := svc.CreateVerbale(context.Background(), "t-1", "teacher", CreateVerbaleRequest{
		MeetingID:   "m-1",
		Title:       "Verbale n.1",
		Content:     "Si approva il piano didattico all'unanimità.",
		SecretaryID: &secID,
		IsPublished: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.False(t, v.IsSigned)
	assert.Equal(t, "draft", v.Status)
	mockRepo.AssertExpectations(t)
}

func TestVerbaleLifecycle_DraftEditByCoordinatorAndSecretary(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	secID := "sec-1"
	coordID := "coord-1"
	otherDocenteID := "docente-other"

	meeting := &CouncilMeeting{ID: "m-1", ClassID: "c-1"}
	draftVerbale := &MeetingVerbale{
		ID:          "v-1",
		MeetingID:   "m-1",
		Title:       "Bozza Iniziale",
		Content:     "Testo iniziale",
		SecretaryID: &secID,
		IsSigned:    false,
		Status:      "draft",
	}

	// 1. Verbalista can edit
	mockRepo.On("GetVerbaleByID", ctx, "v-1", secID).Return(draftVerbale, nil).Once()
	mockRepo.On("GetMeetingWithCoordinator", ctx, "m-1").Return(meeting, coordID, nil).Once()
	mockRepo.On("UpdateVerbale", ctx, mock.Anything).Return(nil).Once()

	updated, err := svc.UpdateVerbale(ctx, secID, "teacher", "v-1", UpdateVerbaleRequest{
		Title:       "Bozza Aggiornata dal Verbalista",
		Content:     "Nuovo testo redatto dal verbalista",
		SecretaryID: &secID,
	})
	assert.NoError(t, err)
	assert.True(t, updated.CanEdit)

	// 2. Coordinator can edit
	mockRepo.On("GetVerbaleByID", ctx, "v-1", coordID).Return(draftVerbale, nil).Once()
	mockRepo.On("GetMeetingWithCoordinator", ctx, "m-1").Return(meeting, coordID, nil).Once()
	mockRepo.On("UpdateVerbale", ctx, mock.Anything).Return(nil).Once()

	updated, err = svc.UpdateVerbale(ctx, coordID, "teacher", "v-1", UpdateVerbaleRequest{
		Title:       "Bozza Revisionata dal Coordinatore",
		Content:     "Note del coordinatore",
		SecretaryID: &secID,
	})
	assert.NoError(t, err)
	assert.True(t, updated.CanEdit)

	// 3. Other teacher cannot edit (Forbidden)
	mockRepo.On("GetVerbaleByID", ctx, "v-1", otherDocenteID).Return(draftVerbale, nil).Once()
	mockRepo.On("GetMeetingWithCoordinator", ctx, "m-1").Return(meeting, coordID, nil).Once()

	_, err = svc.UpdateVerbale(ctx, otherDocenteID, "teacher", "v-1", UpdateVerbaleRequest{
		Title:   "Tentativo non autorizzato",
		Content: "Modifica non permessa",
	})
	assert.ErrorIs(t, err, ErrOnlyCoordinatorOrSecretaryCanEdit)

	mockRepo.AssertExpectations(t)
}

func TestVerbaleLifecycle_PostSigningLock(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	secID := "sec-1"
	signedVerbale := &MeetingVerbale{
		ID:          "v-1",
		MeetingID:   "m-1",
		Title:       "Verbale Firmato",
		Content:     "Documento ufficiale",
		SecretaryID: &secID,
		IsSigned:    true,
		Status:      "signed",
	}

	// 1. Attempt to update signed verbale by secretary -> Locked!
	mockRepo.On("GetVerbaleByID", ctx, "v-1", secID).Return(signedVerbale, nil).Once()
	_, err := svc.UpdateVerbale(ctx, secID, "teacher", "v-1", UpdateVerbaleRequest{
		Title:   "Tentativo modifica post-firma",
		Content: "Non deve passare",
	})
	assert.ErrorIs(t, err, ErrVerbaleLocked)

	// 2. Attempt to update signed verbale by admin -> Locked!
	mockRepo.On("GetVerbaleByID", ctx, "v-1", "admin-1").Return(signedVerbale, nil).Once()
	_, err = svc.UpdateVerbale(ctx, "admin-1", "admin", "v-1", UpdateVerbaleRequest{
		Title:   "Tentativo admin post-firma",
		Content: "Non deve passare",
	})
	assert.ErrorIs(t, err, ErrVerbaleLocked)

	// 3. Attempt to delete signed verbale -> Locked!
	mockRepo.On("GetVerbaleByID", ctx, "v-1", secID).Return(signedVerbale, nil).Once()
	err = svc.DeleteVerbale(ctx, secID, "teacher", "v-1")
	assert.ErrorIs(t, err, ErrVerbaleLocked)

	mockRepo.AssertExpectations(t)
}

func TestVerbaleLifecycle_PrincipalVisibility(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	principalID := "ds-1"
	secID := "sec-1"

	draftVerbale := &MeetingVerbale{
		ID:          "v-draft",
		MeetingID:   "m-1",
		Title:       "Bozza Riservata",
		SecretaryID: &secID,
		IsSigned:    false,
		Status:      "draft",
	}

	signedVerbale := &MeetingVerbale{
		ID:          "v-signed",
		MeetingID:   "m-1",
		Title:       "Verbale Ufficiale Firmato",
		SecretaryID: &secID,
		IsSigned:    true,
		Status:      "signed",
	}

	// 1. Principal attempts to view draft -> Hidden!
	mockRepo.On("GetVerbaleByID", ctx, "v-draft", principalID).Return(draftVerbale, nil).Once()
	_, err := svc.GetVerbale(ctx, "v-draft", principalID, "principal")
	assert.ErrorIs(t, err, ErrDraftHiddenFromPrincipal)

	// 2. Principal lists verbali -> only signed verbali are queried
	mockRepo.On("ListVerbali", ctx, "m-1", principalID, true).Return([]*MeetingVerbale{signedVerbale}, nil).Once()
	mockRepo.On("GetMeetingWithCoordinator", ctx, "m-1").Return(&CouncilMeeting{ID: "m-1"}, "coord-1", nil).Once()
	list, err := svc.ListVerbali(ctx, "m-1", principalID, "principal")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "v-signed", list[0].ID)
	assert.False(t, list[0].CanEdit) // Read-only for everyone post-signing

	// 3. Principal views signed verbale -> Allowed in read-only!
	mockRepo.On("GetVerbaleByID", ctx, "v-signed", principalID).Return(signedVerbale, nil).Once()
	v, err := svc.GetVerbale(ctx, "v-signed", principalID, "principal")
	assert.NoError(t, err)
	assert.Equal(t, "v-signed", v.ID)
	assert.False(t, v.CanEdit)

	mockRepo.AssertExpectations(t)
}

func TestSignVerbaleAndLock(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	secID := "u-1"
	meeting := &CouncilMeeting{ID: "m-1"}
	verbale := &MeetingVerbale{ID: "v-1", MeetingID: "m-1", SecretaryID: &secID, IsPublished: true}

	mockRepo.On("GetVerbaleByID", mock.Anything, "v-1", "u-1").Return(verbale, nil).Once()
	mockRepo.On("GetMeetingWithCoordinator", mock.Anything, "m-1").Return(meeting, "coord-1", nil).Once()
	mockRepo.On("SignVerbale", mock.Anything, "v-1", "u-1", "10.0.0.1").Return(nil).Once()
	mockRepo.On("MarkVerbaleSigned", mock.Anything, "v-1").Return(nil).Once()

	err := svc.SignVerbale(ctx, "v-1", "u-1", "10.0.0.1")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTemplateManagement_PrincipalOnly(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	req := CreateTemplateRequest{
		Title:           "Modello Consiglio di Classe",
		MeetingType:     "consiglio_classe",
		Description:     "Schema standard",
		DefaultAgenda:   "1. Andamento didattico\n2. Varie ed eventuali",
		TemplateContent: "In data odierna si riunisce il consiglio...",
	}

	// 1. Teacher attempts to create template -> Unauthorized
	_, err := svc.CreateTemplate(ctx, "docente-1", "teacher", "school-1", req)
	assert.ErrorIs(t, err, ErrUnauthorized)

	// 2. Principal creates template -> Allowed
	mockRepo.On("CreateTemplate", ctx, mock.MatchedBy(func(t *MeetingVerbaleTemplate) bool {
		return t.Title == "Modello Consiglio di Classe" && t.SchoolID == "school-1"
	})).Return(nil).Once()

	tpl, err := svc.CreateTemplate(ctx, "principal-1", "principal", "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, tpl)

	// 3. List templates
	mockRepo.On("ListTemplates", ctx, "school-1", "consiglio_classe").Return([]*MeetingVerbaleTemplate{tpl}, nil).Once()
	templates, err := svc.ListTemplates(ctx, "school-1", "consiglio_classe")
	assert.NoError(t, err)
	assert.Len(t, templates, 1)

	mockRepo.AssertExpectations(t)
}
