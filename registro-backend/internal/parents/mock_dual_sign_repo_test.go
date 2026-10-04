package parents

import (
	"context"
	"errors"
)

type mockDualSignRepo struct {
	auths   map[string]*DualParentalAuthorization
	custody map[string]*StudentCustodyInfo
	pins    map[string]string
}

func newMockDualSignRepo() *mockDualSignRepo {
	r := &mockDualSignRepo{
		auths:   make(map[string]*DualParentalAuthorization),
		custody: make(map[string]*StudentCustodyInfo),
		pins:    make(map[string]string),
	}

	// Preset test PINs
	r.pins["parent-user-1"] = "1234"
	r.pins["parent-user-2"] = "5678"

	// Preset custody
	r.custody["parent-user-1:student-1"] = &StudentCustodyInfo{
		ParentID:               "parent-user-1",
		StudentID:              "student-1",
		CustodyType:            CustodyShared,
		CanAuthorizeActivities: true,
		IsMirrorNotified:       true,
	}

	r.custody["parent-restricted:student-1"] = &StudentCustodyInfo{
		ParentID:               "parent-restricted",
		StudentID:              "student-1",
		CustodyType:            CustodyRestricted,
		CourtOrderDetails:      "Sentenza Trib. Minorenni n. 123/2025",
		CanAuthorizeActivities: false,
		IsMirrorNotified:       false,
	}

	return r
}

func (m *mockDualSignRepo) CreateAuthorization(ctx context.Context, auth *DualParentalAuthorization) error {
	m.auths[auth.ID] = auth
	return nil
}

func (m *mockDualSignRepo) GetAuthorizationByID(ctx context.Context, id string) (*DualParentalAuthorization, error) {
	if a, ok := m.auths[id]; ok {
		return a, nil
	}
	return nil, errors.New("not found")
}

func (m *mockDualSignRepo) UpdateAuthorization(ctx context.Context, auth *DualParentalAuthorization) error {
	m.auths[auth.ID] = auth
	return nil
}

func (m *mockDualSignRepo) ListAuthorizationsForParent(ctx context.Context, parentID string) ([]DualParentalAuthorization, error) {
	var list []DualParentalAuthorization
	for _, a := range m.auths {
		if a.Parent1ID == parentID || a.Parent2ID == parentID {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockDualSignRepo) GetCustodyInfo(ctx context.Context, parentID, studentID string) (*StudentCustodyInfo, error) {
	key := parentID + ":" + studentID
	if c, ok := m.custody[key]; ok {
		return c, nil
	}
	// Default to shared if unknown
	return &StudentCustodyInfo{
		ParentID:               parentID,
		StudentID:              studentID,
		CustodyType:            CustodyShared,
		CanAuthorizeActivities: true,
		IsMirrorNotified:       true,
	}, nil
}

func (m *mockDualSignRepo) ValidateParentPIN(ctx context.Context, parentID, pin string) (bool, error) {
	expectedPin, ok := m.pins[parentID]
	if !ok {
		return pin == "1234", nil // default PIN if not specified
	}
	return expectedPin == pin, nil
}
