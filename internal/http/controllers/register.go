package controllers

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine,
	wc *WeatherController,
	sc *SubscriptionController,
) {
	// Weather
	r.GET("/weather", wc.GetWeather)
	r.POST("/weather", wc.PostWeather)
	r.PUT("/weather/:city", wc.UpdateWeather)

	// Subscriptions
	r.POST("/subscriptions", sc.CreateSubscription)
	r.GET("/subscriptions/confirm", sc.ConfirmSubscription)
}
