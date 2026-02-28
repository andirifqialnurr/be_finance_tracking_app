package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsService services.StatisticsService
}

func NewStatisticsHandler(statisticsService services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statisticsService: statisticsService}
}

// GetMonthlyStats godoc
// @Summary Get monthly income/expense statistics for a year
// @Tags statistics
// @Produce json
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} models.Response
// @Router /statistics/monthly [get]
func (h *StatisticsHandler) GetMonthlyStats(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	stats, err := h.statisticsService.GetMonthlyStats(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Monthly statistics retrieved successfully",
		Data:    stats,
	})
}

// GetOverview godoc
// @Summary Get financial overview for a specific month
// @Tags statistics
// @Produce json
// @Param month query int false "Month (1-12, default: current month)"
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} models.Response
// @Router /statistics/overview [get]
func (h *StatisticsHandler) GetOverview(c *gin.Context) {
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	overview, err := h.statisticsService.GetOverview(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Financial overview retrieved successfully",
		Data:    overview,
	})
}
