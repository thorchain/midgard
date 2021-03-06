package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/gorilla/mux"
)

var SmokeTestDataEnabled *bool

const eventPageSize = 100

func main() {
	SmokeTestDataEnabled = flag.BoolP("smoke", "s", false, "event the use of the last know smoke-test data suite")
	flag.Parse()
	fmt.Println("SmokeTestDataEnabled: ", *SmokeTestDataEnabled)

	addr := ":8081"
	router := mux.NewRouter()

	// router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/thorchain/events/{id}/{chain}", eventsMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/events/tx/{id}", eventsTxMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/pool_addresses", poolAddressesMockedEndpoint).Methods("GET")
	router.HandleFunc("/thorchain/vaults/asgard", asgardVaultsMockedEndpoint).Methods("GET")

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
	if strings.ToUpper(vars["chain"]) != "BNB" {
		fmt.Fprintf(writer, "[]")
		return
	}
	offset, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var input *os.File
	if *SmokeTestDataEnabled {
		input, err = os.Open("./thorchain/events/smoke-test/events.json")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		input, err = os.Open("./thorchain/events/events.json")
		if err != nil {
			log.Fatal(err)
		}
	}
	dec := json.NewDecoder(input)
	events := make([]map[string]interface{}, 0)
	dec.Token()
	for dec.More() {
		var event map[string]interface{}
		err := dec.Decode(&event)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := strconv.ParseInt(event["id"].(string), 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if id > offset+eventPageSize {
			break
		}
		if id < offset {
			continue
		}
		events = append(events, event)
	}
	content, err := json.Marshal(events)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func poolAddressesMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
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
	// vars := mux.Vars(request)
	// id := vars["id"]

	content, err := ioutil.ReadFile("./thorchain/events/tx/tx.json")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(content))
}

func asgardVaultsMockedEndpoint(writer http.ResponseWriter, request *http.Request) {
	log.Println("asgardVaultsMockedEndpoint Hit!")
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "[{\"chains\":[\"BNB\"]}]")
}
