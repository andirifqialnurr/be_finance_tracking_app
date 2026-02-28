package repositories

import (
	"finance-tracking-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id uuid.UUID) (*models.Account, error)
	FindAll(accountType string, activeOnly bool) ([]models.Account, error)
	FindActive() ([]models.Account, error)
	Update(account *models.Account) error
	UpdateBalance(id uuid.UUID, delta float64) error
	Archive(id uuid.UUID) error
	GetBalanceSummary(id uuid.UUID, month, year int) (*models.AccountBalanceSummary, error)
	GetSavingsProgress(id uuid.UUID) (*models.SavingsProgress, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) FindByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindAll(accountType string, activeOnly bool) ([]models.Account, error) {
	var accounts []models.Account
	q := r.db.Model(&models.Account{})
	if accountType != "" {
		q = q.Where("type = ?", accountType)
	}
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	err := q.Order("created_at ASC").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindActive() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("is_active = ?", true).Order("created_at ASC").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) UpdateBalance(id uuid.UUID, delta float64) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", delta)).Error
}

func (r *accountRepository) Archive(id uuid.UUID) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *accountRepository) GetBalanceSummary(id uuid.UUID, month, year int) (*models.AccountBalanceSummary, error) {
	summary := &models.AccountBalanceSummary{}

	// Get current balance
	var account models.Account
	if err := r.db.Select("balance").First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	summary.CurrentBalance = account.Balance

	// Total income this month
	r.db.Model(&models.Income{}).
		Where("account_id = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", id, month, year).
		Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalIncomeThisMonth)

	// Total spent this month
	r.db.Model(&models.Expense{}).
		Where("account_id = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", id, month, year).
		Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalSpentThisMonth)

	// Total transferred out this month
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

	r.db.Model(&models.AccountTransfer{}).
		Where("from_account_id = ? AND transfer_date >= ? AND transfer_date < ?", id, start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalTransferredOut)

	r.db.Model(&models.AccountTransfer{}).
		Where("to_account_id = ? AND transfer_date >= ? AND transfer_date < ?", id, start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalTransferredIn)

	return summary, nil
}

func (r *accountRepository) GetSavingsProgress(id uuid.UUID) (*models.SavingsProgress, error) {
	var account models.Account
	if err := r.db.First(&account, "id = ? AND type = ?", id, models.AccountTypeSavings).Error; err != nil {
		return nil, err
	}

	progress := &models.SavingsProgress{}
	if account.GoalLabel != nil {
		progress.GoalLabel = *account.GoalLabel
	}
	if account.GoalAmount != nil && *account.GoalAmount > 0 {
		progress.GoalAmount = *account.GoalAmount
		progress.ProgressPercentage = (account.Balance / *account.GoalAmount) * 100
		if progress.ProgressPercentage > 100 {
			progress.ProgressPercentage = 100
		}
		progress.RemainingToGoal = *account.GoalAmount - account.Balance
		if progress.RemainingToGoal < 0 {
			progress.RemainingToGoal = 0
		}
	}
	return progress, nil
}
