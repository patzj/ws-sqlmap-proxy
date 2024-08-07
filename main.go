package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "websocket address")
var path = flag.String("path", "/", "websocket endpoint")

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	log.Println("Listening to port 8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err.Error())
				break
			} else {
				log.Printf("Recv: %s\n", msg)
				w.Write(msg)
				break
			}
		}
	}()

	var pl map[string]any
	if err := json.NewDecoder(r.Body).Decode(&pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
	}

	err = c.WriteJSON(pl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
	log.Printf("Sent: %s\n", pl)

	<-done
}
