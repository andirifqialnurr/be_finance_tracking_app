package repositories

import (
	"finance-tracking-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferRepository interface {
	Create(transfer *models.AccountTransfer) error
	FindByID(id uuid.UUID) (*models.AccountTransfer, error)
	FindAll(limit, offset int) ([]models.AccountTransfer, error)
	FindByAccount(accountID uuid.UUID, month, year int) ([]models.AccountTransfer, error)
	Delete(id uuid.UUID) error
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Create(transfer *models.AccountTransfer) error {
	return r.db.Create(transfer).Error
}

func (r *transferRepository) FindByID(id uuid.UUID) (*models.AccountTransfer, error) {
	var transfer models.AccountTransfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		First(&transfer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *transferRepository) FindAll(limit, offset int) ([]models.AccountTransfer, error) {
	var transfers []models.AccountTransfer
	err := r.db.Preload("FromAccount").Preload("ToAccount").
		Order("transfer_date DESC").
		Limit(limit).Offset(offset).
		Find(&transfers).Error
	return transfers, err
}

func (r *transferRepository) FindByAccount(accountID uuid.UUID, month, year int) ([]models.AccountTransfer, error) {
	var transfers []models.AccountTransfer
	q := r.db.Preload("FromAccount").Preload("ToAccount").
		Where("from_account_id = ? OR to_account_id = ?", accountID, accountID)
	if month > 0 && year > 0 {
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		var endMonth time.Month
		endYear := year
		if month == 12 {
			endMonth = 1
			endYear = year + 1
		} else {
			endMonth = time.Month(month + 1)
		}
		end := time.Date(endYear, endMonth, 1, 0, 0, 0, 0, time.UTC)
		q = q.Where("transfer_date >= ? AND transfer_date < ?", start, end)
	}
	err := q.Order("transfer_date DESC").Find(&transfers).Error
	return transfers, err
}

func (r *transferRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AccountTransfer{}, "id = ?", id).Error
}
