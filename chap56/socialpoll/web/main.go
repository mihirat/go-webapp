package main

import (
	"flag"
	"net/http"
	"log"
)
func main(){
	var addr = flag.String("addr", ":8081", "website address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/",
	http.FileServer(http.Dir("public"))))
	log.Println("website address", *addr)
	http.ListenAndServe(*addr, mux)
}
