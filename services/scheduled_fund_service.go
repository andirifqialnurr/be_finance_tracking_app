package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduledFundService interface {
	CreateScheduledFund(sf *models.ScheduledFund) (*models.ScheduledFund, error)
	GetAll(activeOnly bool) ([]models.ScheduledFund, error)
	GetByID(id uuid.UUID) (*models.ScheduledFund, error)
	UpdateScheduledFund(id uuid.UUID, updates *models.ScheduledFund) (*models.ScheduledFund, error)
	DeleteScheduledFund(id uuid.UUID) error
	ExecuteDueFunds() error
}

type scheduledFundService struct {
	sfRepo       repositories.ScheduledFundRepository
	accountRepo  repositories.AccountRepository
	transferRepo repositories.TransferRepository
	incomeRepo   repositories.IncomeRepository
	db           *gorm.DB
}

func NewScheduledFundService(
	sfRepo repositories.ScheduledFundRepository,
	accountRepo repositories.AccountRepository,
	transferRepo repositories.TransferRepository,
	incomeRepo repositories.IncomeRepository,
	db *gorm.DB,
) ScheduledFundService {
	return &scheduledFundService{
		sfRepo:       sfRepo,
		accountRepo:  accountRepo,
		transferRepo: transferRepo,
		incomeRepo:   incomeRepo,
		db:           db,
	}
}

func (s *scheduledFundService) CreateScheduledFund(sf *models.ScheduledFund) (*models.ScheduledFund, error) {
	if sf.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if sf.DayOfMonth < 0 || sf.DayOfMonth > 31 {
		return nil, errors.New("day_of_month must be between 0 (last day) and 31")
	}

	validTypes := map[string]bool{
		models.ScheduledFundTypeTopUp:    true,
		models.ScheduledFundTypeTransfer: true,
	}
	if !validTypes[sf.ScheduleType] {
		return nil, errors.New("schedule_type must be TOP_UP or TRANSFER")
	}

	if sf.ScheduleType == models.ScheduledFundTypeTransfer && sf.FromAccountID == nil {
		return nil, errors.New("from_account_id is required for TRANSFER type")
	}

	// Validate destination account
	if _, err := s.accountRepo.FindByID(sf.AccountID); err != nil {
		return nil, errors.New("destination account not found")
	}

	// Validate source account for TRANSFER
	if sf.FromAccountID != nil {
		if _, err := s.accountRepo.FindByID(*sf.FromAccountID); err != nil {
			return nil, errors.New("source account not found")
		}
	}

	// Compute first execution date
	sf.NextExecuteAt = computeNextExecuteAt(sf.DayOfMonth, time.Now())
	sf.IsActive = true

	return sf, s.sfRepo.Create(sf)
}

func (s *scheduledFundService) GetAll(activeOnly bool) ([]models.ScheduledFund, error) {
	return s.sfRepo.FindAll(activeOnly)
}

func (s *scheduledFundService) GetByID(id uuid.UUID) (*models.ScheduledFund, error) {
	return s.sfRepo.FindByID(id)
}

func (s *scheduledFundService) UpdateScheduledFund(id uuid.UUID, updates *models.ScheduledFund) (*models.ScheduledFund, error) {
	sf, err := s.sfRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("scheduled fund not found")
	}

	if updates.Amount > 0 {
		sf.Amount = updates.Amount
	}
	if updates.Description != "" {
		sf.Description = updates.Description
	}
	if updates.DayOfMonth >= 0 && updates.DayOfMonth <= 31 {
		sf.DayOfMonth = updates.DayOfMonth
		sf.NextExecuteAt = computeNextExecuteAt(sf.DayOfMonth, time.Now())
	}
	sf.IsActive = updates.IsActive

	return sf, s.sfRepo.Update(sf)
}

func (s *scheduledFundService) DeleteScheduledFund(id uuid.UUID) error {
	if _, err := s.sfRepo.FindByID(id); err != nil {
		return errors.New("scheduled fund not found")
	}
	return s.sfRepo.Delete(id)
}

// ExecuteDueFunds runs as a cron job; executes all due scheduled funds
func (s *scheduledFundService) ExecuteDueFunds() error {
	now := time.Now()
	dueFunds, err := s.sfRepo.FindDueToday(now)
	if err != nil {
		return err
	}

	for _, sf := range dueFunds {
		_ = s.executeSingleFund(sf, now)
	}
	return nil
}

func (s *scheduledFundService) executeSingleFund(sf models.ScheduledFund, now time.Time) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		switch sf.ScheduleType {
		case models.ScheduledFundTypeTopUp:
			// Add income to destination account
			account, err := s.accountRepo.FindByID(sf.AccountID)
			if err != nil || !account.IsActive {
				return errors.New("destination account unavailable")
			}
			income := &models.Income{
				AccountID:   sf.AccountID,
				Source:      "Scheduled Top-Up: " + sf.Description,
				Amount:      sf.Amount,
				Date:        now,
				Description: sf.Description,
			}
			if err := s.incomeRepo.Create(income); err != nil {
				return err
			}
			if err := s.accountRepo.UpdateBalance(sf.AccountID, sf.Amount); err != nil {
				return err
			}

		case models.ScheduledFundTypeTransfer:
			if sf.FromAccountID == nil {
				return errors.New("missing from_account_id for TRANSFER")
			}
			from, err := s.accountRepo.FindByID(*sf.FromAccountID)
			if err != nil || !from.IsActive {
				return errors.New("source account unavailable")
			}
			if from.Balance < sf.Amount {
				return errors.New("insufficient balance for scheduled transfer")
			}
			to, err := s.accountRepo.FindByID(sf.AccountID)
			if err != nil || !to.IsActive {
				return errors.New("destination account unavailable")
			}

			transfer := &models.AccountTransfer{
				FromAccountID: *sf.FromAccountID,
				ToAccountID:   sf.AccountID,
				Amount:        sf.Amount,
				Note:          "Scheduled: " + sf.Description,
				TransferDate:  now,
			}
			if err := s.transferRepo.Create(transfer); err != nil {
				return err
			}
			if err := s.accountRepo.UpdateBalance(*sf.FromAccountID, -sf.Amount); err != nil {
				return err
			}
			if err := s.accountRepo.UpdateBalance(sf.AccountID, sf.Amount); err != nil {
				return err
			}
		}

		// Update execution timestamps
		next := computeNextExecuteAt(sf.DayOfMonth, now)
		return s.sfRepo.UpdateExecution(sf.ID, now, next)
	})
}

// computeNextExecuteAt calculates the next execution date for a given dayOfMonth
// dayOfMonth=0 means last day of month
func computeNextExecuteAt(dayOfMonth int, from time.Time) time.Time {
	// Start from beginning of next month
	y, m, _ := from.Date()
	nextMonth := time.Date(y, m+1, 1, 8, 0, 0, 0, time.UTC)

	if dayOfMonth == 0 {
		// Last day of next month
		lastDay := time.Date(nextMonth.Year(), nextMonth.Month()+1, 0, 8, 0, 0, 0, time.UTC)
		return lastDay
	}

	// Target day of next month; clamp to last day if month is shorter
	lastDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := dayOfMonth
	if day > lastDayOfNextMonth {
		day = lastDayOfNextMonth
	}
	return time.Date(nextMonth.Year(), nextMonth.Month(), day, 8, 0, 0, 0, time.UTC)
}
