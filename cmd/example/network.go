package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Network struct {
	BlockHeight           int64 `json:"block_height"`
	TransactionsPerSecond int64 `json:"tps"`
}

type NetworkStore interface {
	Get(ctx context.Context) (*Network, error)
	Set(ctx context.Context, n *Network) error
}

type StubNetworkStore struct {
	GetFn func(ctx context.Context) (*Network, error)
	SetFn func(ctx context.Context, n *Network) error
}

func (s *StubNetworkStore) Get(ctx context.Context) (*Network, error) {
	return s.GetFn(ctx)
}

func (s *StubNetworkStore) Set(ctx context.Context, n *Network) error {
	return s.SetFn(ctx, n)
}

func getNetwork(networkStore NetworkStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		n, err := networkStore.Get(ctx)
		if err == ErrNotFound {
			return &apiError{err, "not found", http.StatusNotFound}
		}
		if err != nil {
			return &apiError{err, "failed to get network", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(n); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}

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
