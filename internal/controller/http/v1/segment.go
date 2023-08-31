package v1

import (
	"avito-internship/internal/apperror"
	"avito-internship/internal/entity"
	"avito-internship/internal/service"
	"avito-internship/pkg/logging"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type segmentRoutes struct {
	segmentService service.Segment
	l              *logging.Logger
}

func newSegmentRoutes(h *gin.RouterGroup, segmentService service.Segment, l *logging.Logger) {
	r := &segmentRoutes{segmentService, l}

	{
		h.POST("/create", r.create)
		h.DELETE("/delete", r.delete)
	}
}

// @Summary Create segment
// @Tags segment
// @Accept json
// @Produce json
// @Param request body entity.SegmentRequest true "request"
// @Success 201
// @Router /segment/create [post]
func (r *segmentRoutes) create(c *gin.Context) {
	var request entity.SegmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(apperror.ErrBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	err := r.segmentService.CreateSegment(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		if errors.Is(err, apperror.ErrWrongPercent) {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrWrongPercent)

			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// @Summary Delete segment
// @Tags segment
// @Accept json
// @Produce json
// @Param request body entity.SegmentRequest true "request"
// @Success 200
// @Router /segment/delete [delete]
func (r *segmentRoutes) delete(c *gin.Context) {
	var request entity.SegmentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(apperror.ErrBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	err := r.segmentService.DeleteSegment(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
