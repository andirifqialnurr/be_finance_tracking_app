package services

import (
"errors"
"finance-tracking-app/models"
"finance-tracking-app/repositories"

"github.com/google/uuid"
"gorm.io/gorm"
)

// IncomeCreateResult wraps the result of CreateIncome to include optional warning
type IncomeCreateResult struct {
Income                  *models.Income
Allocations             []models.BudgetAllocation
AlreadyAllocatedWarning bool
}

type IncomeService interface {
CreateIncome(income *models.Income) (*IncomeCreateResult, error)
GetIncomeByID(id uuid.UUID) (*models.Income, error)
GetAllIncomes(limit, offset int) ([]models.Income, error)
GetIncomesByMonthYear(month, year int) ([]models.Income, error)
UpdateIncome(id uuid.UUID, income *models.Income) (*IncomeCreateResult, error)
DeleteIncome(id uuid.UUID) error
}

type incomeService struct {
incomeRepo     repositories.IncomeRepository
accountRepo    repositories.AccountRepository
categoryRepo   repositories.ExpenseCategoryRepository
budgetRepo     repositories.CategoryBudgetRepository
allocationRepo repositories.BudgetAllocationRepository
db             *gorm.DB
}

func NewIncomeService(
incomeRepo repositories.IncomeRepository,
accountRepo repositories.AccountRepository,
categoryRepo repositories.ExpenseCategoryRepository,
budgetRepo repositories.CategoryBudgetRepository,
allocationRepo repositories.BudgetAllocationRepository,
db *gorm.DB,
) IncomeService {
return &incomeService{
incomeRepo:     incomeRepo,
accountRepo:    accountRepo,
categoryRepo:   categoryRepo,
budgetRepo:     budgetRepo,
allocationRepo: allocationRepo,
db:             db,
}
}

func (s *incomeService) CreateIncome(income *models.Income) (*IncomeCreateResult, error) {
if income.Amount <= 0 {
return nil, errors.New("income amount must be greater than 0")
}

account, err := s.accountRepo.FindByID(income.AccountID)
if err != nil {
return nil, errors.New("account not found")
}
if !account.IsActive {
return nil, errors.New("account is archived and cannot receive funds")
}

result := &IncomeCreateResult{}

err = s.db.Transaction(func(tx *gorm.DB) error {
if err := s.incomeRepo.Create(income); err != nil {
return err
}
if err := s.accountRepo.UpdateBalance(income.AccountID, income.Amount); err != nil {
return err
}

if account.Type == models.AccountTypeCard &&
account.IncomeType != nil &&
*account.IncomeType == models.CardIncomeTypeSalary {

month := int(income.Date.Month())
year := income.Date.Year()
existingAllocs, _ := s.allocationRepo.FindByMonthYear(month, year)
if len(existingAllocs) > 0 {
result.AlreadyAllocatedWarning = true
} else {
allocs, allocErr := s.runAllocation(income, month, year)
if allocErr != nil {
return allocErr
}
result.Allocations = allocs
}
}
return nil
})

if err != nil {
return nil, err
}
result.Income = income
return result, nil
}

func (s *incomeService) GetIncomeByID(id uuid.UUID) (*models.Income, error) {
return s.incomeRepo.FindByID(id)
}

func (s *incomeService) GetAllIncomes(limit, offset int) ([]models.Income, error) {
if limit <= 0 {
limit = 20
}
return s.incomeRepo.FindAll(limit, offset)
}

func (s *incomeService) GetIncomesByMonthYear(month, year int) ([]models.Income, error) {
if month < 1 || month > 12 {
return nil, errors.New("invalid month")
}
if year < 2000 {
return nil, errors.New("invalid year")
}
return s.incomeRepo.FindByMonthYear(month, year)
}

func (s *incomeService) UpdateIncome(id uuid.UUID, income *models.Income) (*IncomeCreateResult, error) {
if income.Amount <= 0 {
return nil, errors.New("income amount must be greater than 0")
}

result := &IncomeCreateResult{}

err := s.db.Transaction(func(tx *gorm.DB) error {
existing, err := s.incomeRepo.FindByID(id)
if err != nil {
return errors.New("income not found")
}

oldAccountID := existing.AccountID
oldAmount := existing.Amount
oldMonth := int(existing.Date.Month())
oldYear := existing.Date.Year()
newMonth := int(income.Date.Month())
newYear := income.Date.Year()

newAccount, err := s.accountRepo.FindByID(income.AccountID)
if err != nil {
return errors.New("account not found")
}
if !newAccount.IsActive {
return errors.New("account is archived")
}

if oldAccountID == income.AccountID {
delta := income.Amount - oldAmount
if err := s.accountRepo.UpdateBalance(oldAccountID, delta); err != nil {
return err
}
} else {
if err := s.accountRepo.UpdateBalance(oldAccountID, -oldAmount); err != nil {
return err
}
if err := s.accountRepo.UpdateBalance(income.AccountID, income.Amount); err != nil {
return err
}
}

oldAllocations, _ := s.allocationRepo.FindByIncome(id)
for _, oldAlloc := range oldAllocations {
budget, err := s.budgetRepo.FindByCategoryAndMonth(oldAlloc.CategoryID, oldMonth, oldYear)
if err == nil {
budget.AllocatedAmount -= oldAlloc.AllocatedAmount
budget.RemainingAmount -= oldAlloc.AllocatedAmount
_ = s.budgetRepo.Update(budget)
}
_ = s.allocationRepo.Delete(oldAlloc.ID)
}

income.ID = id
if err := s.incomeRepo.Update(income); err != nil {
return err
}

if newAccount.Type == models.AccountTypeCard &&
newAccount.IncomeType != nil &&
*newAccount.IncomeType == models.CardIncomeTypeSalary {

existingAllocs, _ := s.allocationRepo.FindByMonthYear(newMonth, newYear)
if len(existingAllocs) > 0 {
result.AlreadyAllocatedWarning = true
} else {
allocs, allocErr := s.runAllocation(income, newMonth, newYear)
if allocErr != nil {
return allocErr
}
result.Allocations = allocs
}
}
return nil
})

if err != nil {
return nil, err
}
result.Income = income
return result, nil
}

func (s *incomeService) DeleteIncome(id uuid.UUID) error {
return s.db.Transaction(func(tx *gorm.DB) error {
income, err := s.incomeRepo.FindByID(id)
if err != nil {
return errors.New("income not found")
}

month := int(income.Date.Month())
year := income.Date.Year()

allocations, _ := s.allocationRepo.FindByIncome(id)
for _, alloc := range allocations {
budget, err := s.budgetRepo.FindByCategoryAndMonth(alloc.CategoryID, month, year)
if err == nil {
budget.AllocatedAmount -= alloc.AllocatedAmount
budget.RemainingAmount -= alloc.AllocatedAmount
_ = s.budgetRepo.Update(budget)
}
_ = s.allocationRepo.Delete(alloc.ID)
}

if err := s.accountRepo.UpdateBalance(income.AccountID, -income.Amount); err != nil {
return err
}

return s.incomeRepo.Delete(id)
})
}

func (s *incomeService) runAllocation(income *models.Income, month, year int) ([]models.BudgetAllocation, error) {
categories, err := s.categoryRepo.FindAllActive()
if err != nil {
return nil, err
}
if len(categories) == 0 {
return nil, nil
}

effectiveBudgets := make(map[uuid.UUID]float64)
totalBudget := 0.0
for _, cat := range categories {
var eb float64
if cat.Type == models.CategoryTypeDailyContinuous && cat.DailyAmount != nil {
eb = *cat.DailyAmount * float64(daysInMonth(month, year))
} else {
eb = cat.MonthlyBudget
}
effectiveBudgets[cat.ID] = eb
totalBudget += eb
}

if totalBudget == 0 {
return nil, nil
}

var allocations []models.BudgetAllocation
for _, cat := range categories {
eb := effectiveBudgets[cat.ID]
allocationAmount := (eb / totalBudget) * income.Amount
if allocationAmount <= 0 {
continue
}

days := daysInMonth(month, year)
budget, err := s.budgetRepo.FindByCategoryAndMonth(cat.ID, month, year)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
var effectiveDaily *float64
var daysPtr *int
if cat.Type == models.CategoryTypeDailyContinuous && cat.DailyAmount != nil {
effectiveDaily = cat.DailyAmount
daysPtr = &days
}
budget = &models.CategoryBudget{
CategoryID:           cat.ID,
Month:                month,
Year:                 year,
AllocatedAmount:      allocationAmount,
SpentAmount:          0,
RemainingAmount:      allocationAmount,
EffectiveDailyAmount: effectiveDaily,
DaysInMonth:          daysPtr,
}
if err := s.budgetRepo.Create(budget); err != nil {
return nil, err
}
} else {
return nil, err
}
} else {
budget.AllocatedAmount += allocationAmount
budget.RemainingAmount += allocationAmount
_ = s.budgetRepo.Update(budget)
}

alloc := models.BudgetAllocation{
IncomeID:        income.ID,
CategoryID:      cat.ID,
AllocatedAmount: allocationAmount,
Month:           month,
Year:            year,
}
if err := s.allocationRepo.Create(&alloc); err != nil {
return nil, err
}
alloc.Category = cat
allocations = append(allocations, alloc)
}
return allocations, nil
}
