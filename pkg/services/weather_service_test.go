// pkg/services/weather_service_test.go
package services_test

import (
	"errors"
	"testing"

	"myapp/pkg/models"
	"myapp/pkg/services"
)

// spyRepo збирає аргументи викликів repository.WeatherRepository
// і зберігає lastUpdates для перевірки полів
type spyRepo struct {
	lastCity      string
	lastUpdates   map[string]interface{}
	returnWeather models.Weather
	returnErr     error
	savedWeather  *models.Weather
	updateErr     error
}

func (s *spyRepo) GetByCity(city string) (models.Weather, error) {
	s.lastCity = city
	return s.returnWeather, s.returnErr
}
func (s *spyRepo) Save(w *models.Weather) error {
	s.savedWeather = w
	return s.returnErr
}
func (s *spyRepo) UpdateWeather(city string, updates map[string]interface{}) error {
	s.lastCity = city
	s.lastUpdates = updates
	return s.updateErr
}

func TestGetCurrentWeather_Success(t *testing.T) {
	spy := &spyRepo{returnWeather: models.Weather{Temperature: 1.23, Humidity: 45, Condition: "Fog"}}
	svc := services.NewWeatherService(spy)

	w, err := svc.GetCurrentWeather("CityX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if spy.lastCity != "CityX" {
		t.Errorf("expected GetByCity called with CityX, got %s", spy.lastCity)
	}
	// порівнюємо поля окремо
	if w.Temperature != spy.returnWeather.Temperature {
		t.Errorf("expected temperature %.2f, got %.2f", spy.returnWeather.Temperature, w.Temperature)
	}
}

func TestGetCurrentWeather_Error(t *testing.T) {
	spy := &spyRepo{returnErr: errors.New("db fail")}
	svc := services.NewWeatherService(spy)

	_, err := svc.GetCurrentWeather("CityY")
	if err == nil || err.Error() != "db fail" {
		t.Fatalf("expected db fail, got %v", err)
	}
}

func TestSaveWeather_Success(t *testing.T) {
	spy := &spyRepo{returnErr: nil}
	svc := services.NewWeatherService(spy)
	in := &models.Weather{Temperature: 5.5, Humidity: 30, Condition: "Sunny"}

	if err := svc.SaveWeather(in); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if spy.savedWeather == nil || spy.savedWeather.Temperature != in.Temperature || spy.savedWeather.Humidity != in.Humidity {
		t.Errorf("expected saved %+v, got %+v", in, spy.savedWeather)
	}
}

func TestSaveWeather_Error(t *testing.T) {
	spy := &spyRepo{returnErr: errors.New("save fail")}
	svc := services.NewWeatherService(spy)

	if err := svc.SaveWeather(&models.Weather{}); err == nil || err.Error() != "save fail" {
		t.Fatalf("expected save fail, got %v", err)
	}
}

func TestUpdateWeather_Success(t *testing.T) {
	expected := models.Weather{Temperature: 9.99, Humidity: 10, Condition: "Sun"}
	spy := &spyRepo{returnWeather: expected, updateErr: nil}
	svc := services.NewWeatherService(spy)
	in := services.UpdateInput{Temperature: 9.99, Humidity: 10, Condition: "Sun"}

	out, err := svc.UpdateWeather("CityZ", in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// перевіряємо виклик UpdateWeather
	if spy.lastCity != "CityZ" {
		t.Errorf("expected update city CityZ, got %s", spy.lastCity)
	}
	// перевірка переданих updates
	if val, ok := spy.lastUpdates["temperature"]; !ok || val.(float64) != in.Temperature {
		t.Errorf("expected temperature %v, got %v", in.Temperature, val)
	}
	if out.Temperature != expected.Temperature {
		t.Errorf("expected returned temp %.2f, got %.2f", expected.Temperature, out.Temperature)
	}
}

func TestUpdateWeather_ErrorOnUpdate(t *testing.T) {
	spy := &spyRepo{updateErr: errors.New("upd fail")}
	svc := services.NewWeatherService(spy)

	_, err := svc.UpdateWeather("CityA", services.UpdateInput{})
	if err == nil || err.Error() != "upd fail" {
		t.Fatalf("expected upd fail, got %v", err)
	}
}

func TestUpdateWeather_ErrorOnGetAfterUpdate(t *testing.T) {
	spy := &spyRepo{updateErr: nil, returnErr: errors.New("get fail")}
	svc := services.NewWeatherService(spy)

	_, err := svc.UpdateWeather("CityB", services.UpdateInput{})
	if err == nil || err.Error() != "get fail" {
		t.Fatalf("expected get fail, got %v", err)
	}
}
