package main

import (
    "flag"
    "log"
    "net/http"
    "fmt"
    "time"
    "github.com/gorilla/mux"

)

func get_podcast(w http.ResponseWriter, r *http.Request) {
    
}

func main() {
    port := flag.String("p", "8100", "port to serve on")
    flag.Parse()

    r := mux.NewRouter()
    r.HandleFunc("/podcast/youtube/{yt_channel}", get_podcast)

    server := &http.Server{
        Handler:      r,
        Addr:         "0.0.0.0:" + *port,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    fmt.Printf("Listening on %s\n", *port)
    log.Fatal(server.ListenAndServe())
}
