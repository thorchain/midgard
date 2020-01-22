package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	SmokeTestDataEnabled *bool
)

func main() {

	SmokeTestDataEnabled = flag.BoolP("smoke", "s", false, "event the use of the last know smoke-test data suite")
	flag.Parse()
	fmt.Println("SmokeTestDataEnabled: ", *SmokeTestDataEnabled)

	addr := ":8081"
	router := mux.NewRouter()

	// router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/genesis", genesisMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/events/{id}", eventsMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/events/tx/{id}", eventsTxMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/pool_addresses", pool_addresses).Methods("GET")

	// used to debug incorrect dynamically generated requests
	router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Request: ", request.URL)
	})
	// setup server
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	fmt.Println("Running mocked endpoints: ", addr)
	log.Fatal(srv.ListenAndServe())
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

	var content []byte
	var err error
	if *SmokeTestDataEnabled {
		content, err = ioutil.ReadFile("./thorchain/events/smoke-test/events.json")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		content, err = ioutil.ReadFile("./thorchain/events/events.json")
		if err != nil {
			log.Fatal(err)
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func genesisMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("genesisMockedEndpoint Hit!")

	content, err := ioutil.ReadFile("./thorchain/genesis/genesis.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func pool_addresses(writer http.ResponseWriter, request *http.Request) {
	log.Println("pool_addresses Hit!")

	content, err := ioutil.ReadFile("./thorchain/pool_addresses/pool_addresses.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func eventsTxMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("eventsTxMockedEndpoint Hit!")
	//vars := mux.Vars(request)
	//id := vars["id"]

	content, err := ioutil.ReadFile("./thorchain/events/tx/tx.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}
