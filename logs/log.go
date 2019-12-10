package logs

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Animal struct {
	Name string
	age int
}

func Fty_Logs() {
	//log.SetFormatter(&log.JSONFormatter{})
	a := Animal{"dog", 22}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true})
	log.WithFields(log.Fields{
		"event": "ne",
		"topic": "title",
		"key":   "my key",
	}).Info("hello", a)

	log.Error("hello world")
	for {
		time.Sleep(time.Second)
		log.Printf("i am ok %s", "dock")
	}
	log.Fatal("kill ")
}