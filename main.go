package main

import (
	"i9sched/actions"
	"i9sched/setup"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sendgrid/sendgrid-go"
)

func main() {

	apiKey := os.Getenv("SENDGRID_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable is not set")
	}

	sendclient := sendgrid.NewSendClient(apiKey)

	auth := setup.InitFirebase()

	scheduler := gocron.NewScheduler(time.UTC)

	client, database, err := setup.ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer setup.DisConnectDB(client)

	_, err = scheduler.Every(1).Hour().Do(actions.DoneCancels, sendclient, database, auth)
	if err != nil {
		log.Fatalf("Error scheduling done cancels: %s\n", err.Error())
	}
	scheduler.StartAsync()

	select {}
}
