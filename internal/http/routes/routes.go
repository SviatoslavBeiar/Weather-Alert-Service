package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	controllers2 "myapp/internal/http/controllers"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"myapp/pkg/validation"
)

func NewRouter(
	cfg config.Config,
	db *gorm.DB,
	wc *controllers2.WeatherController,
	sc *controllers2.SubscriptionController,
) *gin.Engine {
	database.DB = db

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterConditionValidator(v)
	}

	r := gin.Default()
	controllers2.Register(r, wc, sc)
	return r
}
