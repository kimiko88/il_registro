package meals

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrDietNotFound        = errors.New("dieta speciale non trovata")
	ErrInsufficientBalance = errors.New("saldo insufficiente per l'addebito del pasto")
	ErrWalletBlocked       = errors.New("borsellino mensa bloccato per morosità")
)

type Repository interface {
	SaveDiet(ctx context.Context, diet *SpecialDiet) error
	GetDietByStudent(ctx context.Context, studentID string) (*SpecialDiet, error)
	ListDiets(ctx context.Context, schoolID string) ([]*SpecialDiet, error)
	SaveWallet(ctx context.Context, wallet *MealWallet) error
	GetWalletByStudent(ctx context.Context, studentID string) (*MealWallet, error)
	SaveTransaction(ctx context.Context, tx *WalletTransaction) error
	ListTransactions(ctx context.Context, walletID string) ([]*WalletTransaction, error)
	SaveCateringReport(ctx context.Context, r *CateringOrderReport) error
	GetCateringReport(ctx context.Context, schoolID, date string) (*CateringOrderReport, error)
}

type MemoryRepository struct {
	mu           sync.RWMutex
	diets        map[string]*SpecialDiet         // key: studentID
	wallets      map[string]*MealWallet          // key: studentID
	transactions map[string][]*WalletTransaction // key: walletID
	reports      map[string]*CateringOrderReport // key: schoolID:date
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		diets:        make(map[string]*SpecialDiet),
		wallets:      make(map[string]*MealWallet),
		transactions: make(map[string][]*WalletTransaction),
		reports:      make(map[string]*CateringOrderReport),
	}
}

func (m *MemoryRepository) SaveDiet(_ context.Context, diet *SpecialDiet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.diets[diet.StudentID] = diet
	return nil
}

func (m *MemoryRepository) GetDietByStudent(_ context.Context, studentID string) (*SpecialDiet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.diets[studentID]
	if !ok {
		return nil, ErrDietNotFound
	}
	return d, nil
}

func (m *MemoryRepository) ListDiets(_ context.Context, schoolID string) ([]*SpecialDiet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*SpecialDiet
	for _, d := range m.diets {
		if schoolID != "" && d.SchoolID != schoolID {
			continue
		}
		list = append(list, d)
	}
	return list, nil
}

func (m *MemoryRepository) SaveWallet(_ context.Context, wallet *MealWallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wallets[wallet.StudentID] = wallet
	return nil
}

func (m *MemoryRepository) GetWalletByStudent(_ context.Context, studentID string) (*MealWallet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w, ok := m.wallets[studentID]
	if !ok {
		// Return new default wallet
		return &MealWallet{
			ID:                  fmt.Sprintf("wal-%s", studentID),
			StudentID:           studentID,
			Balance:             0.0,
			MealPrice:           5.50,
			LowBalanceThreshold: 11.00,
			IsBlocked:           false,
			UpdatedAt:           time.Now(),
		}, nil
	}
	return w, nil
}

func (m *MemoryRepository) SaveTransaction(_ context.Context, tx *WalletTransaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[tx.WalletID] = append(m.transactions[tx.WalletID], tx)
	return nil
}

func (m *MemoryRepository) ListTransactions(_ context.Context, walletID string) ([]*WalletTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.transactions[walletID], nil
}

func (m *MemoryRepository) SaveCateringReport(_ context.Context, r *CateringOrderReport) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s:%s", r.SchoolID, r.Date.Format("2006-01-02"))
	m.reports[key] = r
	return nil
}

func (m *MemoryRepository) GetCateringReport(_ context.Context, schoolID, date string) (*CateringOrderReport, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", schoolID, date)
	r, ok := m.reports[key]
	if !ok {
		return nil, errors.New("report non trovato per la data specificata")
	}
	return r, nil
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

func (s *Service) RegisterSpecialDiet(ctx context.Context, schoolID string, req SpecialDietRequest) (*SpecialDiet, error) {
	if req.StudentID == "" {
		return nil, errors.New("student_id is required")
	}

	var expiry *time.Time
	if req.CertificateExpiryDate != "" {
		exp, err := time.Parse("2006-01-02", req.CertificateExpiryDate)
		if err == nil {
			expiry = &exp
		}
	}

	diet := &SpecialDiet{
		ID:                    fmt.Sprintf("diet-%s", req.StudentID),
		SchoolID:              schoolID,
		StudentID:             req.StudentID,
		DietCategory:          req.DietCategory,
		SpecificDiet:          req.SpecificDiet,
		MedicalCertificateURL: req.MedicalCertificateURL,
		CertificateExpiryDate: expiry,
		AllergensList:         req.AllergensList,
		IsApproved:            true,
		Notes:                 req.Notes,
		CreatedAt:             time.Now(),
	}

	if err := s.repo.SaveDiet(ctx, diet); err != nil {
		return nil, err
	}
	return diet, nil
}

func (s *Service) GetSpecialDiet(ctx context.Context, studentID string) (*SpecialDiet, error) {
	return s.repo.GetDietByStudent(ctx, studentID)
}

func (s *Service) ListSpecialDiets(ctx context.Context, schoolID string) ([]*SpecialDiet, error) {
	return s.repo.ListDiets(ctx, schoolID)
}

func (s *Service) RecordRollCallAndAggregate(ctx context.Context, schoolID string, req MealRollCallRequest) (*CateringOrderReport, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	existing, _ := s.repo.GetCateringReport(ctx, schoolID, req.Date)
	report := existing
	if report == nil {
		report = &CateringOrderReport{
			ID:               fmt.Sprintf("cat-%s-%s", schoolID, req.Date),
			SchoolID:         schoolID,
			Date:             date,
			TransmissionTime: time.Now(),
			TotalStandard:    0,
			TotalDiets:       0,
			TotalMeals:       0,
			DietsBreakdown:   make(map[string]int),
			CateringProvider: "Centro Pasti Comunale",
			Status:           "trasmesso",
		}
	}

	for _, stdID := range req.PresentIDs {
		diet, err := s.repo.GetDietByStudent(ctx, stdID)
		if err == nil && diet.IsApproved {
			report.TotalDiets++
			report.DietsBreakdown[diet.SpecificDiet]++
		} else {
			report.TotalStandard++
		}
		report.TotalMeals++

		// Automatically debit meal from wallet
		_, _ = s.ChargeMeal(ctx, stdID)
	}

	if err := s.repo.SaveCateringReport(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *Service) GetWallet(ctx context.Context, studentID string) (*MealWallet, error) {
	return s.repo.GetWalletByStudent(ctx, studentID)
}

func (s *Service) TopUpWallet(ctx context.Context, studentID string, amount float64, pagopaIUV string) (*MealWallet, error) {
	if amount <= 0 {
		return nil, errors.New("l'importo della ricarica deve essere positivo")
	}

	w, err := s.repo.GetWalletByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	w.Balance += amount
	if w.Balance >= 0 {
		w.IsBlocked = false
	}
	w.UpdatedAt = time.Now()

	if err := s.repo.SaveWallet(ctx, w); err != nil {
		return nil, err
	}

	tx := &WalletTransaction{
		ID:           fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		WalletID:     w.ID,
		MovementType: "ricarica_pagopa",
		Amount:       amount,
		BalanceAfter: w.Balance,
		PagoPAIUV:    pagopaIUV,
		Description:  fmt.Sprintf("Ricarica borsellino pasti PagoPA (IUV: %s)", pagopaIUV),
		CreatedAt:    time.Now(),
	}
	_ = s.repo.SaveTransaction(ctx, tx)

	return w, nil
}

func (s *Service) ChargeMeal(ctx context.Context, studentID string) (*MealWallet, error) {
	w, err := s.repo.GetWalletByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if w.IsBlocked {
		return nil, ErrWalletBlocked
	}

	w.Balance -= w.MealPrice
	if w.Balance < -w.LowBalanceThreshold {
		w.IsBlocked = true
	}
	w.UpdatedAt = time.Now()

	if err := s.repo.SaveWallet(ctx, w); err != nil {
		return nil, err
	}

	tx := &WalletTransaction{
		ID:           fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		WalletID:     w.ID,
		MovementType: "addebito_pasto",
		Amount:       -w.MealPrice,
		BalanceAfter: w.Balance,
		Description:  "Addebito pasto consumato in data odierna",
		CreatedAt:    time.Now(),
	}
	_ = s.repo.SaveTransaction(ctx, tx)

	return w, nil
}
