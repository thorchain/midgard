package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Pool struct {
	Symbol                   string  `json:"symbol"`
	Depth                    float64 `json:"depth"`
	Volume24Hour             float64 `json:"volume_24h"`
	AverageTransactionAmount float64 `json:"avg_tx_amount"`
	AverageLiquidityFee      float64 `json:"avg_liquidity_fee"`
	HistoricalROI            float64 `json:"historical_roi"`
}

type PoolStore interface {
	List(ctx context.Context) ([]*Pool, error)
	Get(ctx context.Context, sym string) (*Pool, error)
	Set(ctx context.Context, p *Pool) error
}

type StubPoolStore struct {
	ListFn func(ctx context.Context) ([]*Pool, error)
	GetFn  func(ctx context.Context, sym string) (*Pool, error)
	SetFn  func(ctx context.Context, p *Pool) error
}

func (s *StubPoolStore) List(ctx context.Context) ([]*Pool, error) {
	return s.ListFn(ctx)
}

func (s *StubPoolStore) Get(ctx context.Context, sym string) (*Pool, error) {
	return s.GetFn(ctx, sym)
}

func (s *StubPoolStore) Set(ctx context.Context, p *Pool) error {
	return s.SetFn(ctx, p)
}

func listPools(poolStore PoolStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		ps, err := poolStore.List(ctx)
		if err != nil {
			return &apiError{err, "failed to list pools", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(ps); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}

func getPool(poolStore PoolStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		vars := mux.Vars(r)
		sym := vars["sym"]

		p, err := poolStore.Get(ctx, sym)
		if err == ErrNotFound {
			return &apiError{err, "not found", http.StatusNotFound}
		}
		if err != nil {
			return &apiError{err, "failed to get pool", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(p); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}
