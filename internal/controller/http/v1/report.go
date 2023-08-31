package v1

import (
	"avito-internship/internal/apperror"
	"avito-internship/internal/entity"
	"avito-internship/internal/service"
	"avito-internship/pkg/logging"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type reportRoutes struct {
	reportService service.Report
	l             *logging.Logger
}

func newReportRoutes(h *gin.RouterGroup, reportService service.Report, l *logging.Logger) {
	r := &reportRoutes{reportService, l}

	{
		h.GET("/", r.getHistory)
		h.GET("/link", r.getReportLink)
		h.GET("/file", r.getReportFile)
	}
}

// @Summary Get history JSON
// @Tags report
// @Produce json
// @Param month query string true "month"
// @Param year query string true "year"
// @Success 200 {object} []entity.ReportUserHistory
// @Router /report/ [get]
func (r *reportRoutes) getHistory(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	request := entity.ReportRequest{Month: month, Year: year}
	userHistory, err := r.reportService.GetUserHistory(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusOK, userHistory)
}

// @Summary Get report file
// @Tags report
// @Produce json
// @Param month query string true "month"
// @Param year query string true "year"
// @Success 200 {object} map[string]string
// @Router /report/link [get]
func (r *reportRoutes) getReportLink(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	request := entity.ReportRequest{Month: month, Year: year}
	link, err := r.reportService.MakeReportLink(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		if errors.Is(err, apperror.ErrGDriveNotAvailable) {
			c.AbortWithStatusJSON(http.StatusOK, apperror.ErrGDriveNotAvailable)

			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"Link": link})
}

// @Summary Get report file
// @Tags report
// @Produce text/csv
// @Param month query string true "month"
// @Param year query string true "year"
// @Success 200 {object} []byte
// @Router /report/file [get]
func (r *reportRoutes) getReportFile(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	request := entity.ReportRequest{Month: month, Year: year}
	file, err := r.reportService.MakeReportFile(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.Data(http.StatusOK, "text/csv", file)
}
