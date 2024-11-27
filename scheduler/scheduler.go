package scheduler

import (
	"log"
	"net/http"
	"github.com/go-co-op/gocron"
	"time"
)


func StartScheduler() {
	
	s := gocron.NewScheduler(time.UTC)

	time.AfterFunc(20*time.Minute,func() {
		s.Every(20).Minute().Do(func() {
		resp, err := http.Get("http://localhost:8080/send-email")
		if err != nil {
			log.Printf("Error hitting /send-email route: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Hit /send-email route. Status: %s", resp.Status)
	})

	s.StartAsync()
	})
}
