package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferService interface {
	CreateTransfer(fromID, toID uuid.UUID, amount float64, note string, date time.Time) (*models.AccountTransfer, error)
	GetAllTransfers(limit, offset int) ([]models.AccountTransfer, error)
	GetTransferByID(id uuid.UUID) (*models.AccountTransfer, error)
	GetTransfersByAccount(accountID uuid.UUID, month, year int) ([]models.AccountTransfer, error)
	CancelTransfer(id uuid.UUID) error
}

type transferService struct {
	transferRepo repositories.TransferRepository
	accountRepo  repositories.AccountRepository
	db           *gorm.DB
}

func NewTransferService(
	transferRepo repositories.TransferRepository,
	accountRepo repositories.AccountRepository,
	db *gorm.DB,
) TransferService {
	return &transferService{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		db:           db,
	}
}

func (s *transferService) CreateTransfer(fromID, toID uuid.UUID, amount float64, note string, date time.Time) (*models.AccountTransfer, error) {
	if amount <= 0 {
		return nil, errors.New("transfer amount must be greater than 0")
	}
	if fromID == toID {
		return nil, errors.New("cannot transfer to the same account")
	}

	var transfer models.AccountTransfer

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Load from account
		from, err := s.accountRepo.FindByID(fromID)
		if err != nil {
			return errors.New("source account not found")
		}
		if !from.IsActive {
			return errors.New("source account is archived")
		}
		if from.Balance < amount {
			return errors.New("insufficient balance in source account")
		}

		// Load to account
		to, err := s.accountRepo.FindByID(toID)
		if err != nil {
			return errors.New("destination account not found")
		}
		if !to.IsActive {
			return errors.New("destination account is archived")
		}

		// Deduct from source
		if err := s.accountRepo.UpdateBalance(fromID, -amount); err != nil {
			return err
		}
		// Add to destination
		if err := s.accountRepo.UpdateBalance(toID, amount); err != nil {
			return err
		}

		// Create transfer record
		transfer = models.AccountTransfer{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
			Note:          note,
			TransferDate:  date,
		}
		return s.transferRepo.Create(&transfer)
	})

	if err != nil {
		return nil, err
	}

	// Preload relations for response
	result, _ := s.transferRepo.FindByID(transfer.ID)
	if result != nil {
		return result, nil
	}
	return &transfer, nil
}

func (s *transferService) GetAllTransfers(limit, offset int) ([]models.AccountTransfer, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.transferRepo.FindAll(limit, offset)
}

func (s *transferService) GetTransferByID(id uuid.UUID) (*models.AccountTransfer, error) {
	return s.transferRepo.FindByID(id)
}

func (s *transferService) GetTransfersByAccount(accountID uuid.UUID, month, year int) ([]models.AccountTransfer, error) {
	return s.transferRepo.FindByAccount(accountID, month, year)
}

func (s *transferService) CancelTransfer(id uuid.UUID) error {
	transfer, err := s.transferRepo.FindByID(id)
	if err != nil {
		return errors.New("transfer not found")
	}

	// Only allow cancel within 24 hours
	if time.Since(transfer.CreatedAt) > 24*time.Hour {
		return errors.New("transfer can only be cancelled within 24 hours of creation")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Reverse: add back to source, deduct from destination
		if err := s.accountRepo.UpdateBalance(transfer.FromAccountID, transfer.Amount); err != nil {
			return err
		}
		if err := s.accountRepo.UpdateBalance(transfer.ToAccountID, -transfer.Amount); err != nil {
			return err
		}
		// Verify destination still has enough (they might have spent some)
		dest, err := s.accountRepo.FindByID(transfer.ToAccountID)
		if err != nil {
			return err
		}
		if dest.Balance < 0 {
			return errors.New("destination account has insufficient funds to cancel (balance would go negative)")
		}
		return s.transferRepo.Delete(id)
	})
}
