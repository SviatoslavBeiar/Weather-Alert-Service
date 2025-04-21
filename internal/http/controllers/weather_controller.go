package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"myapp/pkg/models"
	"myapp/pkg/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WeatherController struct {
	Svc    *services.WeatherService
	Logger *zap.Logger
}

func NewWeatherController(svc *services.WeatherService, logger *zap.Logger) *WeatherController {
	return &WeatherController{Svc: svc, Logger: logger}
}

func (h *WeatherController) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		h.errorResponse(c, http.StatusBadRequest, "city required")
		return
	}

	w, err := h.Svc.GetCurrentWeather(city)
	if err != nil {
		if errors.Is(err, services.ErrCityNotFound) {
			h.errorResponse(c, http.StatusNotFound, "city not found")
		} else {
			h.logError("GetWeather failed", zap.Error(err))
			h.errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, ResponseDTO{Status: "success", Data: w})
}

func (h *WeatherController) PostWeather(c *gin.Context) {
	var w models.Weather
	if err := c.ShouldBindJSON(&w); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Svc.SaveWeather(&w); err != nil {
		if errors.Is(err, services.ErrCityNotFound) {
			h.errorResponse(c, http.StatusNotFound, "city not found")
		} else {
			h.logError("SaveWeather failed", zap.Error(err))
			h.errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.Header("Location", fmt.Sprintf("/weather/%s", w.City))
	c.JSON(http.StatusCreated, ResponseDTO{Status: "success", Data: w})
}

func (h *WeatherController) UpdateWeather(c *gin.Context) {
	city := c.Param("city")
	var inp services.UpdateInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	w, err := h.Svc.UpdateWeather(city, inp)
	if err != nil {
		if errors.Is(err, services.ErrCityNotFound) {
			h.errorResponse(c, http.StatusNotFound, "city not found")
		} else {
			h.logError("UpdateWeather failed", zap.Error(err))
			h.errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, ResponseDTO{Status: "success", Data: w})
}

func (h *WeatherController) errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, ResponseDTO{Status: "error", Error: msg})
}

func (h *WeatherController) logError(msg string, fields ...zap.Field) {
	if h.Logger != nil {
		h.Logger.Error(msg, fields...)
	}
}
