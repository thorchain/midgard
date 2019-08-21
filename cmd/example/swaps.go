package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Swap struct {
	Symbol       string `json:"symbol"`
	AvgTxToken   int64  `json:"aveTxTkn"`
	AvgTxRune    int64  `json:"aveTxRune"`
	AvgSlipToken int64  `json:"aveSlipTkn"`
	AvgSlipRune  int64  `json:"aveSlipRune"`
	NumTxToken   int64  `json:"numTxTkn"`
	NumTxRune    int64  `json:"numTxRune"`
	AvgFeeToken  int64  `json:"aveFeeTkn"`
	AvgFeeRune   int64  `json:"aveFeeRune"`
}

type SwapStore interface {
	List(ctx context.Context) ([]*Swap, error)
	Get(ctx context.Context, sym string) (*Swap, error)
	Set(ctx context.Context, s *Swap) error
}

type StubSwapStore struct {
	ListFn func(ctx context.Context) ([]*Swap, error)
	GetFn  func(ctx context.Context, sym string) (*Swap, error)
	SetFn  func(ctx context.Context, s *Swap) error
}

func (s *StubSwapStore) List(ctx context.Context) ([]*Swap, error) {
	return s.ListFn(ctx)
}

func (s *StubSwapStore) Get(ctx context.Context, sym string) (*Swap, error) {
	return s.GetFn(ctx, sym)
}

func (s *StubSwapStore) Set(ctx context.Context, sw *Swap) error {
	return s.SetFn(ctx, sw)
}

func getPoolSwap(swapStore SwapStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		vars := mux.Vars(r)
		sym := vars["symbol"]

		s, err := swapStore.Get(ctx, sym)
		if err == ErrNotFound {
			return &apiError{err, "not found", http.StatusNotFound}
		}
		if err != nil {
			return &apiError{err, "failed to get swap", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(s); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}
