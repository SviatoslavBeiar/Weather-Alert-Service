package scheduler

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"myapp/pkg/database"
	"myapp/pkg/repository"
	"myapp/pkg/services"
)

func Start() {

	spec := os.Getenv("CRON_SCHEDULE")
	if spec == "" {
		spec = "@daily"
	}

	repo := repository.NewGormRepo()
	ws := services.NewWeatherService(repo)
	ss := services.NewSubscriptionService(repo, repo)

	c := cron.New(cron.WithSeconds())

	job := func() {
		now := time.Now()
		subs, err := ss.ListVerified()
		if err != nil {
			log.Println("subscription fetch error:", err)
			return
		}
		for _, sub := range subs {

			if sub.LastSent != nil && now.Sub(*sub.LastSent) < 24*time.Hour {
				continue
			}

			w, err := ws.GetCurrentWeather(sub.City)
			if err != nil {
				log.Println("weather fetch error:", err)
				continue
			}

			sent, err := services.EvaluateAndNotify(sub, w)
			if err != nil {
				log.Println("notify error:", err)
				continue
			}
			if sent {

				database.DB.Model(&sub).Update("LastSent", now)
			}
		}
	}

	if _, err := c.AddFunc(spec, job); err != nil {
		log.Fatalf("invalid CRON_SCHEDULE %q: %v", spec, err)
	}

	c.Start()
}
