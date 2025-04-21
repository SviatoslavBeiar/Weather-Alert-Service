package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"myapp/pkg/models"
	"myapp/pkg/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ResponseDTO struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type SubscriptionController struct {
	Svc    *services.SubscriptionService
	Logger *zap.Logger
}

func NewSubscriptionController(svc *services.SubscriptionService, logger *zap.Logger) *SubscriptionController {
	return &SubscriptionController{Svc: svc, Logger: logger}
}

func (h *SubscriptionController) CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Svc.Create(&sub); err != nil {
		switch {

		case errors.Is(err, services.ErrCityNotFound):
			h.errorResponse(c, http.StatusNotFound, err.Error())

		case errors.Is(err, services.ErrDuplicateSubscription), strings.Contains(err.Error(), "Duplicate entry"):
			h.errorResponse(c, http.StatusConflict, "subscription already exists")

		default:
			h.logError("CreateSubscription failed", zap.Error(err))
			h.errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.Header("Location", fmt.Sprintf("/subscriptions/%d", sub.ID))
	c.JSON(http.StatusCreated, ResponseDTO{
		Status: "success",
		Data: gin.H{
			"message":         "Check your email and click on the confirmation link.",
			"subscription_id": sub.ID,
		},
	})
}

func (h *SubscriptionController) ConfirmSubscription(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		h.errorResponse(c, http.StatusBadRequest, "token required")
		return
	}

	confirmedSub, err := h.Svc.Confirm(token)
	if err != nil {
		switch {

		case errors.Is(err, services.ErrTokenNotFound):
			h.errorResponse(c, http.StatusNotFound, "invalid token")

		case errors.Is(err, services.ErrTokenExpired):
			h.errorResponse(c, http.StatusGone, "token expired")

		default:
			h.logError("ConfirmSubscription failed", zap.Error(err))
			h.errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, ResponseDTO{
		Status: "success",
		Data: gin.H{
			"message":         "Email verified",
			"subscription_id": confirmedSub.ID,
		},
	})
}

func (h *SubscriptionController) errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, ResponseDTO{Status: "error", Error: msg})
}

func (h *SubscriptionController) logError(msg string, fields ...zap.Field) {
	if h.Logger != nil {
		h.Logger.Error(msg, fields...)
	}
}
