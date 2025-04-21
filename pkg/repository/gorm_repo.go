package repository

import (
	"myapp/pkg/database"
	models2 "myapp/pkg/models"
)

type GormRepo struct{}

func NewGormRepo() *GormRepo {
	return &GormRepo{}
}

// --- Weather ---
func (r *GormRepo) GetByCity(city string) (models2.Weather, error) {
	var w models2.Weather
	err := database.DB.First(&w, "city = ?", city).Error
	return w, err
}

func (r *GormRepo) Save(w *models2.Weather) error {
	return database.DB.Save(w).Error
}

func (r *GormRepo) UpdateWeather(city string, updates map[string]interface{}) error {
	return database.DB.
		Model(&models2.Weather{}).
		Where("city = ?", city).
		Updates(updates).
		Error
}

// --- Subscription ---
func (r *GormRepo) Create(sub *models2.Subscription) error {
	return database.DB.Create(sub).Error
}

func (r *GormRepo) FindAllVerified() ([]models2.Subscription, error) {
	var subs []models2.Subscription
	err := database.DB.Where("verified = ?", true).Find(&subs).Error
	return subs, err
}

func (r *GormRepo) FindByToken(token string) (models2.Subscription, error) {
	var sub models2.Subscription
	err := database.DB.Where("verification_token = ?", token).First(&sub).Error
	return sub, err
}

func (r *GormRepo) UpdateSubscription(sub *models2.Subscription) error {
	return database.DB.Save(sub).Error
}
