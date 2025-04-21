package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"myapp/pkg/models"
	"myapp/pkg/repository"
	"myapp/pkg/utils"
	"time"
)

type SubscriptionService struct {
	SubRepo     repository.SubscriptionRepository
	WeatherRepo repository.WeatherRepository
}

func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	weatherRepo repository.WeatherRepository,
) *SubscriptionService {
	return &SubscriptionService{
		SubRepo:     subRepo,
		WeatherRepo: weatherRepo,
	}
}

func (s *SubscriptionService) Create(sub *models.Subscription) error {
	log.Printf("Create: start subscription for email=%s, city=%s", sub.Email, sub.City)

	// 1) Перевіряємо наявність міста в БД
	if _, err := s.WeatherRepo.GetByCity(sub.City); err != nil {
		log.Printf("Create: city not found=%s, err=%v", sub.City, err)
		return fmt.Errorf("місто %q не знайдено", sub.City)
	}

	// 2) Генеруємо токен
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Create: token generation failed, err=%v", err)
		return err
	}
	token := hex.EncodeToString(b)
	expires := time.Now().Add(24 * time.Hour)

	sub.Verified = false
	sub.VerificationToken = token
	sub.TokenExpiresAt = &expires

	// 3) Зберігаємо підписку у репозиторій
	if err := s.SubRepo.Create(sub); err != nil {
		log.Printf("Create: failed to save subscription, err=%v", err)
		return err
	}
	log.Printf("Create: subscription saved, id=%d", sub.ID)

	// 4) Відправляємо лист для підтвердження
	link := fmt.Sprintf("http://localhost:8080/subscriptions/confirm?token=%s", token)
	subject := "Please confirm your subscription"
	body := fmt.Sprintf("Click to confirm: %s\nExpires at: %s", link, expires.Format(time.RFC1123))
	if err := utils.SendEmail(sub.Email, subject, body); err != nil {
		log.Printf("Create: failed to send email to %s, err=%v", sub.Email, err)
		return err
	}
	log.Printf("Create: confirmation email sent to %s", sub.Email)

	return nil
}

// Confirm підтверджує підписку за токеном
func (s *SubscriptionService) Confirm(token string) (*models.Subscription, error) {
	log.Printf("Confirm: start confirm for token=%s", token)

	sub, err := s.SubRepo.FindByToken(token)
	if err != nil {
		log.Printf("Confirm: token not found err=%v", err)
		return nil, err
	}
	if sub.TokenExpiresAt == nil || time.Now().After(*sub.TokenExpiresAt) {
		log.Printf("Confirm: token expired for email=%s", sub.Email)
		return nil, fmt.Errorf("token expired")
	}

	sub.Verified = true
	sub.VerificationToken = ""
	sub.TokenExpiresAt = nil

	//  Оновлюємо статус у репозиторії
	if err := s.SubRepo.UpdateSubscription(&sub); err != nil {
		log.Printf("Confirm: failed to update subscription, err=%v", err)
		return nil, err
	}
	log.Printf("Confirm: subscription confirmed for email=%s", sub.Email)
	return &sub, nil
}

// ListVerified повертає всі підтверджені підписки
func (s *SubscriptionService) ListVerified() ([]models.Subscription, error) {
	log.Printf("ListVerified: fetching all verified subscriptions")
	return s.SubRepo.FindAllVerified()
}
