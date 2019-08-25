package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// apiError is used to simplify api error handling
// it allows a handler to return an error and further instructions
// in the form of a desired message and status code
type apiError struct {
	Error      error
	Message    string
	StatusCode int
}

func (e apiError) String() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Message, e.Error)
}

// handlerWithError is an extension of http.HandlerFunc
// which expects the handler to potentially return an apiError
type handlerWithError func(http.ResponseWriter, *http.Request) *apiError

// Define a base json error response so we can return errors to clients
type errResponse struct {
	Error string `json:"error"`
}

// newErrorResponse creates a serialized errResponse from a given error string
func newErrorResponse(err string) string {
	// This shouldn't ever fail
	bs, _ := json.Marshal(&errResponse{Error: err})
	return string(bs)
}

// mwJSONError wraps a handlerWithError so that if it returned an apiError
// it will be returned to the client in JSON format
func mwJSONError(hn handlerWithError) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		// Call the handler
		err := hn(w, r)

		// Check if there was an error
		if err != nil {
			// Respond to client with json wrapped error
			http.Error(w, newErrorResponse(err.Message), err.StatusCode)
		}

		return err
	}
}

func mwDiscardError(hn handlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call handler, but discard error
		_ = hn(w, r)
	}
}
