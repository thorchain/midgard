package main

import "net/http"

type Token struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Links       map[string][]string `json:"links"`
	Logo        string              `json:"logo"`
	Price       float64             `json:"price"`
}

type ListTokensResponse []Token

func listTokens() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}

type GetTokenResponse Token

func getToken() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}
