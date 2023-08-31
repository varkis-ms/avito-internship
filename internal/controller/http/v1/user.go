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

type userRoutes struct {
	userService service.User
	l           *logging.Logger
}

func newUserRoutes(h *gin.RouterGroup, userService service.User, l *logging.Logger) {
	r := &userRoutes{userService, l}

	{
		h.POST("/add", r.add)
		h.DELETE("/remove", r.remove)
		h.GET("/get", r.get)
	}
}

// @Summary Add user to segment
// @Tags user
// @Accept json
// @Produce json
// @Param request body entity.UserAddToSegmentRequest true "request"
// @Success 200
// @Router /user/add [post]
func (r *userRoutes) add(c *gin.Context) {
	var request entity.UserAddToSegmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(apperror.ErrBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	if request.Ttl < 0 {
		r.l.Error(apperror.ErrWrongTtl)
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrWrongTtl)

		return
	}

	err := r.userService.AddSegment(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		if errors.Is(err, apperror.ErrNoSegment) {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrNoSegment)

			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// @Summary Remove user from segment
// @Tags user
// @Accept json
// @Produce json
// @Param request body entity.UserRemoveFromSegmentRequest true "request"
// @Success 200
// @Router /user/remove [delete]
func (r *userRoutes) remove(c *gin.Context) {
	var request entity.UserRemoveFromSegmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(apperror.ErrBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	err := r.userService.RemoveSegment(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err)
		if errors.Is(err, apperror.ErrNoUser) {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrNoUser)

			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// @Summary Get active user's segments
// @Tags user
// @Produce json
// @Param user_id query string true "user_id"
// @Success 200 {object} map[string][]string
// @Router /user/get [get]
func (r *userRoutes) get(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrBadRequest)

		return
	}

	request := entity.UserActiveSegmentRequest{UserId: userId}
	segments, err := r.userService.GetActiveSegments(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, apperror.ErrNoUser) {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperror.ErrNoUser)

			return
		}
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.SystemError(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"segment": segments})
}
