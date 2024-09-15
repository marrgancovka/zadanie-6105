package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	//logger := logrus.New()
	//logger.SetFormatter(&logrus.TextFormatter{
	//	TimestampFormat: time.DateTime,
	//	FullTimestamp:   true,
	//})
	//application := app.NewApp(logger)
	//if err := application.Start(); err != nil {
	//	logger.Fatal(err)
	//}
	var happiness = 0.04
	var isHappy = true
	if happiness >= 0.5 || isHappy {
		fmt.Println("Happy ")
	}
	fmt.Println(":)")
}
