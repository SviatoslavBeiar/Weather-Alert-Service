// pkg/services/subscription_service_test.go
package services_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"myapp/pkg/models"
	"myapp/pkg/services"
	"myapp/pkg/utils"
)

// mockSubRepo збирає аргументи викликів і повертає помилки за налаштуванням
type mockSubRepo struct {
	lastCreated  *models.Subscription
	createErr    error
	findByToken  models.Subscription
	findErr      error
	lastUpdated  *models.Subscription
	updateErr    error
	verifiedList []models.Subscription
	listErr      error
}

func (m *mockSubRepo) Create(sub *models.Subscription) error {
	// Зберігаємо лише основні поля
	m.lastCreated = &models.Subscription{
		Email:             sub.Email,
		City:              sub.City,
		VerificationToken: sub.VerificationToken,
		TokenExpiresAt:    sub.TokenExpiresAt,
		Verified:          sub.Verified,
	}
	return m.createErr
}
func (m *mockSubRepo) FindByToken(token string) (models.Subscription, error) {
	return m.findByToken, m.findErr
}
func (m *mockSubRepo) UpdateSubscription(sub *models.Subscription) error {
	m.lastUpdated = &models.Subscription{
		Email:             sub.Email,
		City:              sub.City,
		Verified:          sub.Verified,
		VerificationToken: sub.VerificationToken,
		TokenExpiresAt:    sub.TokenExpiresAt,
	}
	return m.updateErr
}
func (m *mockSubRepo) FindAllVerified() ([]models.Subscription, error) {
	return m.verifiedList, m.listErr
}

// mockWeatherRepo перевіряє наявність міста
type mockWeatherRepo struct {
	exists bool
	err    error
}

func (m *mockWeatherRepo) GetByCity(city string) (models.Weather, error) {
	if !m.exists {
		return models.Weather{}, m.err
	}
	return models.Weather{}, nil
}
func (m *mockWeatherRepo) Save(w *models.Weather) error { return nil }
func (m *mockWeatherRepo) UpdateWeather(city string, updates map[string]interface{}) error {
	return nil
}

func TestSubscriptionService_Create(t *testing.T) {
	orig := utils.SendEmail
	defer func() { utils.SendEmail = orig }()

	cases := []struct {
		name      string
		exists    bool
		createErr error
		emailErr  error
		wantErr   bool
		check     func(t *testing.T, m *mockSubRepo)
	}{
		{"CityNotFound", false, nil, nil, true, nil},
		{"RepoError", true, errors.New("db err"), nil, true, nil},
		{"EmailError", true, nil, errors.New("send err"), true, nil},
		{"Success", true, nil, nil, false, func(t *testing.T, m *mockSubRepo) {
			sub := m.lastCreated
			if sub == nil {
				t.Fatal("Expected Create to be called")
			}
			if sub.Email != "e@e" || sub.City != "C" {
				t.Errorf("Unexpected Email/City: %+v", sub)
			}
			// Токен має бути 32 hex-символи
			if match, _ := regexp.MatchString("^[0-9a-f]{32}$", sub.VerificationToken); !match {
				t.Errorf("Invalid token format: %q", sub.VerificationToken)
			}
			// Термін дії має бути за межами зараз() + [23h,25h]
			diff := sub.TokenExpiresAt.Sub(time.Now())
			if diff < 23*time.Hour || diff > 25*time.Hour {
				t.Errorf("Unexpected TokenExpiresAt: %v", sub.TokenExpiresAt)
			}
			if sub.Verified {
				t.Error("Expected Verified to be false")
			}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			utils.SendEmail = func(_, _, _ string) error { return tc.emailErr }
			mSub := &mockSubRepo{createErr: tc.createErr}
			mW := &mockWeatherRepo{exists: tc.exists, err: errors.New("not found")}
			svc := services.NewSubscriptionService(mSub, mW)
			sub := &models.Subscription{Email: "e@e", City: "C"}

			err := svc.Create(sub)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got %v", tc.wantErr, err)
			}
			if tc.check != nil {
				tc.check(t, mSub)
			}
		})
	}
}

func TestSubscriptionService_Confirm(t *testing.T) {
	now := time.Now()
	valid := now.Add(time.Hour)
	expired := now.Add(-time.Hour)

	cases := []struct {
		name      string
		repoSub   models.Subscription
		repoErr   error
		updateErr error
		wantErr   bool
		check     func(t *testing.T, m *mockSubRepo)
	}{
		{"NotFound", models.Subscription{}, errors.New("no token"), nil, true, nil},
		{"Expired", models.Subscription{TokenExpiresAt: &expired}, nil, nil, true, nil},
		{"UpdateError", models.Subscription{TokenExpiresAt: &valid}, nil, errors.New("upd err"), true, nil},
		{"Success", models.Subscription{TokenExpiresAt: &valid}, nil, nil, false, func(t *testing.T, m *mockSubRepo) {
			upd := m.lastUpdated
			if upd == nil {
				t.Fatal("Expected UpdateSubscription to be called")
			}
			if !upd.Verified {
				t.Error("Expected Verified=true")
			}
			if upd.VerificationToken != "" {
				t.Error("Expected VerificationToken cleared")
			}
			if upd.TokenExpiresAt != nil {
				t.Error("Expected TokenExpiresAt cleared")
			}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mSub := &mockSubRepo{findByToken: tc.repoSub, findErr: tc.repoErr, updateErr: tc.updateErr}
			svc := services.NewSubscriptionService(mSub, nil)
			_, err := svc.Confirm("tok")
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got %v", tc.wantErr, err)
			}
			if tc.check != nil {
				tc.check(t, mSub)
			}
		})
	}
}

func TestSubscriptionService_ListVerified(t *testing.T) {
	expected := []models.Subscription{{Email: "a"}, {Email: "b"}}
	mSub := &mockSubRepo{verifiedList: expected}
	svc := services.NewSubscriptionService(mSub, nil)

	out, err := svc.ListVerified()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(expected) {
		t.Fatalf("expected %d, got %d", len(expected), len(out))
	}
}
