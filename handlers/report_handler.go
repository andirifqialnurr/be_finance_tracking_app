package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"finance-tracking-app/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetMonthlyReport godoc
// @Summary Get monthly financial report
// @Description Get comprehensive financial report for a specific month
// @Tags reports
// @Produce json
// @Param month query int false "Month (1-12, default: current month)"
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/monthly [get]
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid month parameter",
		})
		return
	}

	if year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	report, err := h.reportService.GetMonthlyReport(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Monthly report retrieved successfully",
		Data:    report,
	})
}

// GetYearlyReport godoc
// @Summary Get yearly financial report
// @Description Get comprehensive financial report for an entire year
// @Tags reports
// @Produce json
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/yearly [get]
func (h *ReportHandler) GetYearlyReport(c *gin.Context) {
	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	if year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	report, err := h.reportService.GetYearlyReport(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Yearly report retrieved successfully",
		Data:    report,
	})
}

// ExportMonthlyReport godoc
// @Summary Export monthly financial report
// @Description Export comprehensive financial report for a specific month as PDF or Excel
// @Tags reports
// @Produce application/pdf,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param month query int false "Month (1-12, default: current month)"
// @Param year query int false "Year (default: current year)"
// @Param format query string true "Export format: pdf or excel"
// @Success 200 {file} binary
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/monthly/export [get]
func (h *ReportHandler) ExportMonthlyReport(c *gin.Context) {
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	format := c.DefaultQuery("format", "pdf")

	if month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid month parameter",
		})
		return
	}

	if year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	if format != "pdf" && format != "excel" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid format parameter, must be 'pdf' or 'excel'",
		})
		return
	}

	report, err := h.reportService.GetMonthlyReport(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var contentType string
	var filename string

	if format == "pdf" {
		pdfBuf, err := utils.ExportMonthlyReportToPDF(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		contentType = "application/pdf"
		filename = fmt.Sprintf("monthly_report_%s_%d.pdf", time.Month(month).String(), year)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, contentType, pdfBuf.Bytes())
	} else {
		excelBuf, err := utils.ExportMonthlyReportToExcel(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = fmt.Sprintf("monthly_report_%s_%d.xlsx", time.Month(month).String(), year)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, contentType, excelBuf.Bytes())
	}
}

// ExportYearlyReport godoc
// @Summary Export yearly financial report
// @Description Export comprehensive financial report for an entire year as PDF or Excel
// @Tags reports
// @Produce application/pdf,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param year query int false "Year (default: current year)"
// @Param format query string true "Export format: pdf or excel"
// @Success 200 {file} binary
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/yearly/export [get]
func (h *ReportHandler) ExportYearlyReport(c *gin.Context) {
	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	format := c.DefaultQuery("format", "pdf")

	if year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	if format != "pdf" && format != "excel" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid format parameter, must be 'pdf' or 'excel'",
		})
		return
	}

	report, err := h.reportService.GetYearlyReport(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var contentType string
	var filename string

	if format == "pdf" {
		pdfBuf, err := utils.ExportYearlyReportToPDF(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		contentType = "application/pdf"
		filename = fmt.Sprintf("yearly_report_%d.pdf", year)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, contentType, pdfBuf.Bytes())
	} else {
		excelBuf, err := utils.ExportYearlyReportToExcel(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = fmt.Sprintf("yearly_report_%d.xlsx", year)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, contentType, excelBuf.Bytes())
	}
}
