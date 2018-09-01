package main

import (
	"net/http"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
}

func handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePollsGet(w, r)
		return
	case "POST":
		handlePollsPost(w, r)
		return
	case "DELETE":
		handlePollsDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func handlePollsGet(w http.ResponseWriter, r *http.Request) {
	log.Println("get request.")
	db := GetVar(r, databaseSessionKey).(*mgo.Database)
	c := db.C(databaseTableName)
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		q = c.Find(nil)
	}
	var result []*poll
	// all results on memory. avoid if getting larger.
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}
func handlePollsPost(w http.ResponseWriter, r *http.Request) {
	log.Println("post request.")
	db := GetVar(r, databaseSessionKey).(*mgo.Database)
	c := db.C(databaseTableName)
	var p poll
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "cannot read poll title from request", err)
		return
	}
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to restore a new poll", err)
		return
	}
	w.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}
func handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	db := GetVar(r, databaseSessionKey).(*mgo.Database)
	c := db.C(databaseTableName)
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "cannot delete all title")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delet", err)
		return
	}
	respond(w, r, http.StatusOK, nil)
}
