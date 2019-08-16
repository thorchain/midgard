package main

import "net/http"

func listNodes() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}

func listValidators() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}
