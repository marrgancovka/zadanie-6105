package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
	"zadanie-6105/internal/app"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})
	application := app.NewApp(logger)
	if err := application.Start(); err != nil {
		logger.Fatal(err)
	}
}
