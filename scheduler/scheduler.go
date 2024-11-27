package scheduler

import (
	"log"
	"net/http"
	"github.com/go-co-op/gocron"
	"time"
)


func StartScheduler() {
	
	s := gocron.NewScheduler(time.UTC)

	s.Every(2).Minute().Do(func() {
		resp, err := http.Get("http://localhost:8080/send-email")
		if err != nil {
			log.Printf("Error hitting /getAllUsers route: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Hit /getAllUsers route. Status: %s", resp.Status)
	})

	s.StartAsync()
}
