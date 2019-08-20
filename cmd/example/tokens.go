package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

var (
	ErrNotFound = errors.New("not found")
)

type TokenStore interface {
	List(ctx context.Context) ([]*Token, error)
	Get(ctx context.Context, sym string) (*Token, error)
	Set(ctx context.Context, t *Token) error
}

type StubTokenStore struct {
	ListFn func(ctx context.Context) ([]*Token, error)
	GetFn  func(ctx context.Context, sym string) (*Token, error)
	SetFn  func(ctx context.Context, t *Token) error
}

func (s *StubTokenStore) List(ctx context.Context) ([]*Token, error) {
	return s.ListFn(ctx)
}

func (s *StubTokenStore) Get(ctx context.Context, sym string) (*Token, error) {
	return s.GetFn(ctx, sym)
}

func (s *StubTokenStore) Set(ctx context.Context, t *Token) error {
	return s.SetFn(ctx, t)
}

func listTokens(tokenStore TokenStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		ts, err := tokenStore.List(ctx)
		if err != nil {
			return &apiError{err, "failed to list tokens", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(ts); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}

func getToken(tokenStore TokenStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		vars := mux.Vars(r)
		sym := vars["sym"]

		t, err := tokenStore.Get(ctx, sym)
		if err == ErrNotFound {
			return &apiError{err, "not found", http.StatusNotFound}
		}
		if err != nil {
			return &apiError{err, "failed to get token", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(t); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}

type Price struct {
	Symbol string `json:"symbol"`
	Ticker string `json:"ticker"`
	Price  int64  `json:"price"`
}

type PriceStore interface {
	List(ctx context.Context) ([]*Price, error)
	Get(ctx context.Context, sym string) (*Price, error)
	Set(ctx context.Context, p *Price) error
}

type StubPriceStore struct {
	ListFn func(ctx context.Context) ([]*Price, error)
	GetFn  func(ctx context.Context, sym string) (*Price, error)
	SetFn  func(ctx context.Context, p *Price) error
}

func (s *StubPriceStore) List(ctx context.Context) ([]*Price, error) {
	return s.ListFn(ctx)
}

func (s *StubPriceStore) Get(ctx context.Context, sym string) (*Price, error) {
	return s.GetFn(ctx, sym)
}

func (s *StubPriceStore) Set(ctx context.Context, p *Price) error {
	return s.SetFn(ctx, p)
}

func listPrices(priceStore PriceStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		ps, err := priceStore.List(ctx)
		if err != nil {
			return &apiError{err, "failed to list prices", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(ps); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}

func getPrice(priceStore PriceStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		vars := mux.Vars(r)
		sym := vars["sym"]

		p, err := priceStore.Get(ctx, sym)
		if err == ErrNotFound {
			return &apiError{err, "not found", http.StatusNotFound}
		}
		if err != nil {
			return &apiError{err, "failed to get price", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(p); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}
