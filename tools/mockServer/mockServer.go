package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func eventsMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("eventsMockedEndpoint Hit!")
	vars := mux.Vars(request)

	id := vars["id"]

	if id != "1" {
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(writer, "[]")
		return
	}

	content, err := ioutil.ReadFile("./events/events.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func genesisMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("genesisMockedEndpoint Hit!")

	content, err := ioutil.ReadFile("./test/mocks/thorNode/genesis.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func welcome(writer http.ResponseWriter, request *http.Request) {
	log.Println("Welcome Hit!")
	fmt.Fprintf(writer, "Welcome to thorMock")
}

func main() {
	addr := ":8081"
	router := mux.NewRouter()

	router.HandleFunc("/thorchain/events/{id}", eventsMockedEndpoint).Methods("GET")
	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/genesis", genesisMockedEndpoint).Methods("GET")

	// setup server
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
	}

	fmt.Println("Running mocked endpoints: ", addr)
	log.Fatal(srv.ListenAndServe())
}
