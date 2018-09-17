package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matryer/filedb"
	"github.com/mihirat/go-webapp/chap8/backup"
)

// it appears twice, but not so big to create abstruct class
// so declare twice.
// and it might have different member variables
type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()
	var (
		interval = flag.Int("interval", 10, "check interval(second)")
		archive  = flag.String("archive", "archive", "directory to save archived file")
		dbpath   = flag.String("db", "./db", "path to filedb database")
	)
	flag.Parse()

	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.ZIP,
		Paths:       make(map[string]string),
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}
	var path path
	col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}
		m.Paths[path.Path] = path.Hash
		return false // continues
	})
	if fatalErr != nil {
		fatalErr = err
		return
	}
	if len(m.Paths) < 1 {
		fatalErr = errors.New("no path. add new with 'backup' tool")
	}

	check(m, col)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
Loop:
	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m, col)
		case <-signalChan:
			fmt.Println()
			log.Printf("finishing...")
			break Loop
		}
	}
}

func check(m *backup.Monitor, col *filedb.C) {
	log.Println("start checking...")
	counter, err := m.Now()
	if err != nil {
		log.Panicln("failed to backup:", err)
	}
	if counter > 0 {
		log.Printf(" archived %d directories\n", counter)
		// update hash
		var path path
		col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
			if err := json.Unmarshal(data, &path); err != nil {
				log.Println("failed to load json data"+"go to next item: ", err)
				return true, data, false
			}
			path.Hash, _ = m.Paths[path.Path]
			newdata, err := json.Marshal(&path)
			if err != nil {
				log.Println("failed to write json data, "+"go to next item: ", err)
				return true, data, false
			}
			return true, newdata, false
		})
	} else {
		log.Println("no new diff.")
	}
}
