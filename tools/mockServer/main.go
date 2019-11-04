package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func poolsMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("poolsMockedEndpoint Hit!")

	content, err := ioutil.ReadFile("./test/mocks/thorNode/pools.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func eventsMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("eventsMockedEndpoint Hit!")
	vars := mux.Vars(request)

	id := vars["id"]

	if id != "1" {
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(writer, "[]")
		return
	}

	content, err := ioutil.ReadFile("./test/mocks/thorNode/seed_events.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func main() {
	addr := "127.0.0.1:8081"
	router := mux.NewRouter()

	router.HandleFunc("/swapservice/events/{id}", eventsMockedEndpoint).Methods("GET")
	router.HandleFunc("/swapservice/pools", poolsMockedEndpoint).Methods("GET")

	// setup server
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	fmt.Println("Running mocked endpoints: ", addr)
	log.Fatal(srv.ListenAndServe())
}
