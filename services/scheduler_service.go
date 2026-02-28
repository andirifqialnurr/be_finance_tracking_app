package services

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// SchedulerService handles scheduled tasks like monthly report generation
type SchedulerService struct {
	reportService        ReportService
	scheduledFundService ScheduledFundService
	notificationService  NotificationService
	cron                 *cron.Cron
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(
	reportService ReportService,
	scheduledFundService ScheduledFundService,
	notificationService NotificationService,
) *SchedulerService {
	c := cron.New(cron.WithSeconds())
	return &SchedulerService{
		reportService:        reportService,
		scheduledFundService: scheduledFundService,
		notificationService:  notificationService,
		cron:                 c,
	}
}

// Start begins all scheduled tasks
func (s *SchedulerService) Start() {
	log.Println("Starting scheduler service...")

	// Schedule monthly report generation - 1st day of month at 00:01:00
	_, err := s.cron.AddFunc("0 1 0 1 * *", s.GenerateMonthlyReport)
	if err != nil {
		log.Printf("Error scheduling monthly report generation: %v", err)
	} else {
		log.Println("Scheduled: Monthly report generation (1st day of month at 00:01:00)")
	}

	// Execute due scheduled funds - daily at 08:00:00
	_, err = s.cron.AddFunc("0 0 8 * * *", s.ExecuteScheduledFunds)
	if err != nil {
		log.Printf("Error scheduling fund execution: %v", err)
	} else {
		log.Println("Scheduled: Scheduled fund execution (daily at 08:00)")
	}

	// Send due push notifications - every minute
	_, err = s.cron.AddFunc("0 * * * * *", s.SendDueNotifications)
	if err != nil {
		log.Printf("Error scheduling notifications: %v", err)
	} else {
		log.Println("Scheduled: Push notification delivery (every minute)")
	}

	s.cron.Start()
	log.Println("Scheduler service started successfully")
}

// Stop stops all scheduled tasks
func (s *SchedulerService) Stop() {
	log.Println("Stopping scheduler service...")
	s.cron.Stop()
	log.Println("Scheduler service stopped")
}

// GenerateMonthlyReport generates a monthly report for the previous month
func (s *SchedulerService) GenerateMonthlyReport() {
	now := time.Now()
	// Generate report for the previous month
	lastMonth := now.AddDate(0, -1, 0)
	month := int(lastMonth.Month())
	year := lastMonth.Year()

	log.Printf("Auto-generating monthly report for %s %d...", time.Month(month).String(), year)

	report, err := s.reportService.GetMonthlyReport(month, year)
	if err != nil {
		log.Printf("Error generating monthly report for %s %d: %v", time.Month(month).String(), year, err)
		return
	}

	// Log summary
	log.Printf("Monthly Report Generated - %s %d", time.Month(month).String(), year)
	log.Printf("  Total Income: Rp %.2f", report.Income.Total)
	log.Printf("  Total Expenses: Rp %.2f", report.Expenses.Total)
	log.Printf("  Total Savings: Rp %.2f", report.Savings)
	log.Printf("  Top Category: %s (Rp %.2f)", report.TopCategory, report.TopCategoryAmount)

	// TODO: Optionally save to database (monthly_reports table)
	// TODO: Optionally send email notification
	// TODO: Optionally generate and store PDF/Excel files

	log.Printf("Monthly report generation completed successfully")
}

// ExecuteScheduledFunds runs due scheduled top-ups and transfers
func (s *SchedulerService) ExecuteScheduledFunds() {
	log.Println("Running scheduled fund execution...")
	if err := s.scheduledFundService.ExecuteDueFunds(); err != nil {
		log.Printf("Error executing scheduled funds: %v", err)
	}
}

// SendDueNotifications sends push notifications whose time has come
func (s *SchedulerService) SendDueNotifications() {
	if err := s.notificationService.CheckAndSendDue(); err != nil {
		log.Printf("Error sending notifications: %v", err)
	}
}

// ManualGenerateMonthlyReport allows manual trigger of monthly report generation
func (s *SchedulerService) ManualGenerateMonthlyReport(month, year int) error {
	log.Printf("Manually generating monthly report for %s %d...", time.Month(month).String(), year)

	report, err := s.reportService.GetMonthlyReport(month, year)
	if err != nil {
		return fmt.Errorf("error generating monthly report: %w", err)
	}

	log.Printf("Manual monthly report generated successfully for %s %d", time.Month(month).String(), year)
	log.Printf("  Total Income: Rp %.2f", report.Income.Total)
	log.Printf("  Total Expenses: Rp %.2f", report.Expenses.Total)
	log.Printf("  Total Savings: Rp %.2f", report.Savings)

	return nil
}

// GetScheduledJobs returns information about scheduled jobs
func (s *SchedulerService) GetScheduledJobs() []cron.Entry {
	return s.cron.Entries()
}
