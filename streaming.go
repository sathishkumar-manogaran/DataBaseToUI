package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/sse"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

func insertRandomDataToDatabase() {
	for {
		productionDatabase, err := gorm.Open(postgres.Open(databaseConnection), &gorm.Config{})
		productionDb, _ := productionDatabase.DB()
		if err != nil {
			fmt.Println("Error opening database looks like it does not exist")
			productionDb.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		randomNumber := strconv.Itoa(rand.Intn(100-0) + 0)
		fmt.Println("Inserting random number :: " + randomNumber)
		data := Data{
			Data: randomNumber,
		}
		productionDatabase.Save(&data)
		productionDb.Close()
		time.Sleep(2 * time.Second)

	}
}

func streamDataToWebPage(streamer *sse.Streamer) {
	_, err := sql.Open("postgres", databaseConnection)
	if err != nil {
		fmt.Println("Error opening database :: ", err.Error())
		return
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(databaseConnection, 10*time.Second, time.Minute, reportProblem)
	listener.Listen("events")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		waitForNotification(listener, streamer)
	}
}

func waitForNotification(listener *pq.Listener, streamer *sse.Streamer) {
	for {
		select {
		case n := <-listener.Notify:
			streamer.SendString("data", "data", n.Extra)
			return
		}
	}
}
