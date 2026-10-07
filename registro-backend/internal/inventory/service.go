package inventory

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	ErrAssetNotFound     = errors.New("cespite/bene inventariale non trovato")
	ErrContractNotFound  = errors.New("contratto di comodato non trovato")
	ErrAssetNotAvailable = errors.New("il cespite non è disponibile per il comodato d'uso")
)

// CalculateDepreciation computes linear asset depreciation according to D.I. 129/2018.
func CalculateDepreciation(initialValue, rate float64, acquisitionDate, atDate time.Time) float64 {
	if initialValue <= 0 {
		return 0
	}
	if rate <= 0 {
		return initialValue
	}

	days := atDate.Sub(acquisitionDate).Hours() / 24
	if days <= 0 {
		return initialValue
	}

	years := days / 365.25
	depreciatedAmount := initialValue * (rate / 100.0) * years
	residual := initialValue - depreciatedAmount
	if residual < 0 {
		residual = 0
	}

	return math.Round(residual*100) / 100
}

type Repository interface {
	SaveAsset(ctx context.Context, item *AssetItem) error
	GetAssetByID(ctx context.Context, id string) (*AssetItem, error)
	ListAssets(ctx context.Context, schoolID, category, location string) ([]*AssetItem, error)
	GetNextInventoryNumber(ctx context.Context, schoolID string, year int) (int, error)
	SaveLoan(ctx context.Context, loan *LoanContract) error
	GetLoanByID(ctx context.Context, id string) (*LoanContract, error)
	ListLoansByStudent(ctx context.Context, studentID string) ([]*LoanContract, error)
}

type MemoryRepository struct {
	mu     sync.RWMutex
	assets map[string]*AssetItem
	loans  map[string]*LoanContract
	invSeq map[string]int // key: schoolID:year
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		assets: make(map[string]*AssetItem),
		loans:  make(map[string]*LoanContract),
		invSeq: make(map[string]int),
	}
}

func (m *MemoryRepository) SaveAsset(_ context.Context, item *AssetItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.assets[item.ID] = item
	return nil
}

func (m *MemoryRepository) GetAssetByID(_ context.Context, id string) (*AssetItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.assets[id]
	if !ok {
		return nil, ErrAssetNotFound
	}
	return a, nil
}

func (m *MemoryRepository) ListAssets(_ context.Context, schoolID, category, location string) ([]*AssetItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*AssetItem
	for _, a := range m.assets {
		if schoolID != "" && a.SchoolID != schoolID {
			continue
		}
		if category != "" && !strings.EqualFold(a.Category, category) {
			continue
		}
		if location != "" && !strings.Contains(strings.ToLower(a.RoomLocation), strings.ToLower(location)) {
			continue
		}
		list = append(list, a)
	}
	return list, nil
}

func (m *MemoryRepository) GetNextInventoryNumber(_ context.Context, schoolID string, year int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%d", schoolID, year)
	m.invSeq[key]++
	return m.invSeq[key], nil
}

func (m *MemoryRepository) SaveLoan(_ context.Context, loan *LoanContract) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loans[loan.ID] = loan
	return nil
}

func (m *MemoryRepository) GetLoanByID(_ context.Context, id string) (*LoanContract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.loans[id]
	if !ok {
		return nil, ErrContractNotFound
	}
	return l, nil
}

func (m *MemoryRepository) ListLoansByStudent(_ context.Context, studentID string) ([]*LoanContract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*LoanContract
	for _, l := range m.loans {
		if l.StudentID == studentID {
			list = append(list, l)
		}
	}
	return list, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		repo = NewMemoryRepository()
	}
	return &Service{repo: repo}
}

func (s *Service) CreateAsset(ctx context.Context, schoolID string, req CreateAssetRequest) (*AssetItem, error) {
	if schoolID == "" {
		return nil, errors.New("school_id is required")
	}

	acqDate, err := time.Parse("2006-01-02", req.AcquisitionDate)
	if err != nil {
		acqDate = time.Now()
	}

	year := acqDate.Year()
	nextNum, err := s.repo.GetNextInventoryNumber(ctx, schoolID, year)
	if err != nil {
		return nil, err
	}

	rate := req.DepreciationRate
	if rate <= 0 {
		rate = 20.0 // Standard 20% D.I. 129/2018
	}

	invNumber := fmt.Sprintf("INV-%d-%04d", year, nextNum)
	currentVal := CalculateDepreciation(req.InitialValue, rate, acqDate, time.Now())

	barcode := fmt.Sprintf("INV%04d%04d", year, nextNum)
	qr := fmt.Sprintf("INVENTARIO|%s|%s|%s|%s", invNumber, req.SerialNumber, req.BuildingLocation, req.RoomLocation)

	item := &AssetItem{
		ID:               fmt.Sprintf("asset-%s", invNumber),
		SchoolID:         schoolID,
		InventoryNumber:  invNumber,
		Category:         req.Category,
		Description:      req.Description,
		SerialNumber:     req.SerialNumber,
		BuildingLocation: req.BuildingLocation,
		RoomLocation:     req.RoomLocation,
		AcquisitionDate:  acqDate,
		InitialValue:     req.InitialValue,
		DepreciationRate: rate,
		CurrentValue:     currentVal,
		BarcodeData:      barcode,
		QRData:           qr,
		Status:           "in_uso",
		CreatedAt:        time.Now(),
	}

	if err := s.repo.SaveAsset(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetLabel(ctx context.Context, assetID string) (*AssetLabel, error) {
	a, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return &AssetLabel{
		InventoryNumber:  a.InventoryNumber,
		Description:      a.Description,
		BarcodeData:      a.BarcodeData,
		QRCodeData:       a.QRData,
		BuildingLocation: a.BuildingLocation,
		RoomLocation:     a.RoomLocation,
	}, nil
}

func (s *Service) ListAssets(ctx context.Context, schoolID, category, location string) ([]*AssetItem, error) {
	return s.repo.ListAssets(ctx, schoolID, category, location)
}

func (s *Service) CreateLoanContract(ctx context.Context, schoolID string, req CreateLoanRequest) (*LoanContract, error) {
	asset, err := s.repo.GetAssetByID(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}
	if asset.Status != "in_uso" {
		return nil, ErrAssetNotAvailable
	}

	returnDate, err := time.Parse("2006-01-02", req.ExpectedReturnDate)
	if err != nil {
		returnDate = time.Now().Add(180 * 24 * time.Hour) // 6 months
	}

	contractNum := fmt.Sprintf("COMODATO-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	now := time.Now()

	contract := &LoanContract{
		ID:                   fmt.Sprintf("loan-%d", now.UnixNano()),
		SchoolID:             schoolID,
		AssetID:              req.AssetID,
		StudentID:            req.StudentID,
		ParentUserID:         req.ParentUserID,
		ContractNumber:       contractNum,
		HandoverDate:         now,
		ExpectedReturnDate:   returnDate,
		DeviceConditionNotes: req.DeviceConditionNotes,
		ParentSignatureToken: fmt.Sprintf("SIGN-PARENT-%s", contractNum),
		SignedAt:             &now,
		Status:               "attivo",
		CreatedAt:            now,
	}

	if err := s.repo.SaveLoan(ctx, contract); err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *Service) ReturnLoanDevice(ctx context.Context, contractID string, conditionNotes string) (*LoanContract, error) {
	contract, err := s.repo.GetLoanByID(ctx, contractID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	contract.ActualReturnDate = &now
	contract.Status = "restituito"
	if conditionNotes != "" {
		contract.DeviceConditionNotes = fmt.Sprintf("%s | Restituzione: %s", contract.DeviceConditionNotes, conditionNotes)
	}

	if err := s.repo.SaveLoan(ctx, contract); err != nil {
		return nil, err
	}
	return contract, nil
}
