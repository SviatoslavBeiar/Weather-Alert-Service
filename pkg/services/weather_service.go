package services

import (
	"log"
	"myapp/pkg/models"
	"myapp/pkg/repository"
)

type WeatherService struct {
	Repo repository.WeatherRepository
}

func NewWeatherService(r repository.WeatherRepository) *WeatherService {
	log.Println("Initializing WeatherService")
	return &WeatherService{Repo: r}
}

func (s *WeatherService) GetCurrentWeather(city string) (models.Weather, error) {
	log.Printf("GetCurrentWeather called with city=%q", city)
	w, err := s.Repo.GetByCity(city)
	if err != nil {
		log.Printf("GetCurrentWeather error for city=%q: %v", city, err)
		return models.Weather{}, err
	}
	log.Printf("GetCurrentWeather success for city=%q: %+v", city, w)
	return w, nil
}

func (s *WeatherService) SaveWeather(w *models.Weather) error {
	log.Printf("SaveWeather called for city=%q: %+v", w.City, w)
	err := s.Repo.Save(w)
	if err != nil {
		log.Printf("SaveWeather error for city=%q: %v", w.City, err)
		return err
	}
	log.Printf("SaveWeather success for city=%q", w.City)
	return nil
}

func (s *WeatherService) UpdateWeather(city string, inp UpdateInput) (models.Weather, error) {
	log.Printf("UpdateWeather called for city=%q with updates=%+v", city, inp)
	updates := map[string]interface{}{
		"temperature": inp.Temperature,
		"humidity":    inp.Humidity,
		"condition":   inp.Condition,
	}
	err := s.Repo.UpdateWeather(city, updates)
	if err != nil {
		log.Printf("UpdateWeather error for city=%q: %v", city, err)
		return models.Weather{}, err
	}
	w, err := s.Repo.GetByCity(city)
	if err != nil {
		log.Printf("Fetch after UpdateWeather error for city=%q: %v", city, err)
		return models.Weather{}, err
	}
	log.Printf("UpdateWeather success for city=%q: %+v", city, w)
	return w, nil
}

type UpdateInput struct {
	Temperature float64 `json:"temperature" binding:"required"`
	Humidity    int     `json:"humidity"    binding:"required,gte=0,lte=100"`
	Condition   string  `json:"condition"   binding:"required"`
}
