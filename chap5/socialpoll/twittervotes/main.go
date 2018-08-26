package main

import (
	"log"
	"sync"
	"syscall"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/bitly/go-nsq"
)

var db *mgo.Session

const (
	databaseName       = "ballots"
	databaseCollection = "polls"
)

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB(databaseName).C(databaseCollection).Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Clone()
	log.Println("closed db connection")
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			// publish votes
			pub.Publish("votes", []byte(vote))
		}
		log.Println("publisher: stopping")
		pub.Stop()
		log.Println("publisher: stopped")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {
	// 2 goroutines use same "stop". to avoid conflict
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	if err := dialdb(); err != nil{
		log.Fatalln("failed to dial mongo db:", err)
	}
	defer closedb()
	go func(){
		// waiting signal and followings start once received
		<-signalChan
		stop = false
		stoplock.Unlock()
		log.Println("stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	//start processing
	
	// channel for votes result
	votes := make(chan string)
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func(){
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	// after stopChan received, this chan closes
	<-twitterStoppedChan
	close(votes)
	// closing vote chan send signal to stop publish chan
	<-publisherStoppedChan
}
