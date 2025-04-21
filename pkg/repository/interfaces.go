package repository

import (
	models2 "myapp/pkg/models"
)

// WeatherRepository описує операції з моделлю Weather
type WeatherRepository interface {
	GetByCity(city string) (models2.Weather, error)
	Save(w *models2.Weather) error
	UpdateWeather(city string, updates map[string]interface{}) error
}

// SubscriptionRepository описує операції з моделлю Subscription
type SubscriptionRepository interface {
	Create(sub *models2.Subscription) error
	FindAllVerified() ([]models2.Subscription, error)
	FindByToken(token string) (models2.Subscription, error)
	UpdateSubscription(sub *models2.Subscription) error
}
