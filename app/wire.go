//go:build wireinject
// +build wireinject

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	controllers2 "myapp/internal/http/controllers"
	"myapp/internal/http/routes"
	"myapp/pkg/config"
	"myapp/pkg/database"
	repository2 "myapp/pkg/repository"
	services2 "myapp/pkg/services"
)

func InitializeApp() (*gin.Engine, error) {
	wire.Build(
		// Конфігурація
		config.NewConfig,
		database.Connect,

		// Репозиторії
		repository2.NewGormRepo,
		wire.Bind(new(repository2.WeatherRepository), new(*repository2.GormRepo)),
		wire.Bind(new(repository2.SubscriptionRepository), new(*repository2.GormRepo)),

		// Сервіси
		services2.NewWeatherService,
		services2.NewSubscriptionService,

		// Порожній набір опцій для zap
		wire.Value([]zap.Option{}),

		// Логер
		zap.NewProduction,

		// Контролери тепер з логером
		controllers2.NewWeatherController,
		controllers2.NewSubscriptionController,

		// Роутер
		routes.NewRouter,
	)
	return &gin.Engine{}, nil
}
