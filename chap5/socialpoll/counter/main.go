package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
)

var fatalErr error

const (
	databaseName       = "ballots"
	databaseCollection = "polls"
	updateDuration     = 1 * time.Second
)

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}
func main() {
	var countsLock sync.Mutex
	var counts map[string]int
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()
	log.Println("connect to database ...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("close db connection...")
		db.Close()
	}()
	pollData := db.DB(databaseName).C(databaseCollection)

	log.Println("connect to NSQ...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}
	// called by every nsq message receive action
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		// to avoid dup. write from multi go routines
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))
	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}

	log.Println("waiting votes on NSQ...")
	var updater *time.Timer

	// execute func on goroutine in specified interval
	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("no new votes. skip updating database")
		} else {
			log.Println("update database...")
			log.Println(counts)
			ok := true
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if _, err := pollData.UpdateAll(sel, up); err != nil {
					log.Println("failed updating:", err)
					ok = false
					continue
				}
				counts[option] = 0
			}
			if ok {
				log.Println("completed update database")
				// reset counts not to count duplicately
				counts = nil
			}
		}
		// repeat same process
		updater.Reset(updateDuration)
	})

	// wait ctrl+c
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-termChan: // in case of ctrl+c
			updater.Stop()
			q.Stop() // loop is blocked until updater stops
		case <-q.StopChan:
			return
		}
	}
}
