package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/stretchr/graceful"
	"gopkg.in/mgo.v2"
)

const (
	APIKey             = "abc123"
	databaseName       = "ballots"
	databaseSessionKey = "db"
	databaseTableName  = "polls"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "127.0.0.1:27017", "mongodb address")
	)
	flag.Parse()
	log.Println("connect to MongoDB", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect MongoDB", err)
	}
	defer db.Close()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db,
		withAPIKey(handlePolls)))))
	log.Println("start web server:", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("stopping...")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key.")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == APIKey
}

func withData(d *mgo.Session, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		// to share db session among all handler
		SetVar(r, databaseSessionKey, thisDb.DB(databaseName))
		fn(w, r)
	}
}

// to safely and easily use GetVar and SetVar
// this kind of handler wrapper enables to share common procedure
func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

// just to clarify what to do for CORS.
// normally github has better lib.
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
	}
}
