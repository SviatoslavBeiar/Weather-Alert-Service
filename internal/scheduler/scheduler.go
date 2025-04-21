package scheduler

import (
	"log"
	"myapp/pkg/database"
	"myapp/pkg/repository"
	services2 "myapp/pkg/services"
	"time"

	"github.com/robfig/cron/v3"
)

func Start() {

	repo := repository.NewGormRepo()
	ws := services2.NewWeatherService(repo)
	ss := services2.NewSubscriptionService(repo, repo)

	c := cron.New()
	job := cron.FuncJob(func() {
		now := time.Now()

		subs, err := ss.ListVerified()
		if err != nil {
			log.Println("subscription fetch error:", err)
			return
		}

		for _, sub := range subs {

			if sub.LastSent != nil && now.Sub(*sub.LastSent) < 1*time.Minute {
				continue
			}

			w, err := ws.GetCurrentWeather(sub.City)
			if err != nil {
				log.Println("weather fetch error:", err)
				continue
			}

			sent, err := services2.EvaluateAndNotify(sub, w)
			if err != nil {
				log.Println("notify error:", err)
				continue
			}

			if sent {
				database.DB.Model(&sub).Update("LastSent", now)
			}
		}
	})

	c.Schedule(cron.Every(1*time.Minute), job)
	c.Start()
}
