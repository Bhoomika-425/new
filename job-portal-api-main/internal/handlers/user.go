package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"project/internal/middleware"
	"project/internal/models"
	service "project/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type handler struct {
	service service.UserService
}

//go:generate mockgen -source=user.go -destination=user_mock.go -package=handlers
type UserHandler interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	ViewCompany(c *gin.Context)
	ViewAllCompanies(c *gin.Context)
	AddCompany(c *gin.Context)
	JobByID(c *gin.Context)
	AllJobs(c *gin.Context)
	Jobs(c *gin.Context)
	CreateJobs(c *gin.Context)
}
func Newhandler(s service.UserService) (UserHandler, error) {
	if s == nil {
		return nil, errors.New("service cannot be null")
	}
	return &handler{
		service: s,
	}, nil
}


func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	var userData models.NewUser

	err := json.NewDecoder(c.Request.Body).Decode(&userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid email and password",
		})
		return
	}

	token, err := h.service.UserLogin(ctx, userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func (h *handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	traceid, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	var userData models.NewUser

	err := json.NewDecoder(c.Request.Body).Decode(&userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid username, email and password",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid username, email and password",
		})
		return
	}

	userDetails, err := h.service.UserSignup(ctx, userData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userDetails)

}
