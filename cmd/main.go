package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tempo/internal/service"
	"tempo/internal/types"
)

func main() {
	fmt.Println("🚀 Starting Tempo Scheduler Test...")

	scheduler := service.NewScheduler()

	job := types.Job{
		ID:       "test",
		URL:      "https://httpbin.org/get",
		CronExpr: "*/5 * * * * *",
	}

	scheduler.AddJob(job)
	scheduler.Start()

	fmt.Println("Press Ctrl+C to stop")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	scheduler.Stop()
}
