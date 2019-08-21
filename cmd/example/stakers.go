package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Stake struct {
	StakerAddress   string `json:"staker_address"`
	Symbol          string `json:"symbol"`
	StakeUnits      int64  `json:"stake_units"`
	RuneStaked      int64  `json:"rune_staked"`
	TokenStaked     int64  `json:"token_staked"`
	RuneValueStaked int64  `json:"rune_value_staked"`
	InitialStakeAt  int64  `json:"initial_stake_at"`
}

type StakeStore interface {
	List(ctx context.Context, addr string) ([]*Stake, error)
}

type StubStakeStore struct {
	ListFn func(ctx context.Context, addr string) ([]*Stake, error)
}

func (s *StubStakeStore) List(ctx context.Context, addr string) ([]*Stake, error) {
	return s.ListFn(ctx, addr)
}

func listStakes(stakeStore StakeStore) handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		ctx := r.Context()

		vars := mux.Vars(r)
		addr := vars["address"]

		ss, err := stakeStore.List(ctx, addr)
		if err != nil {
			return &apiError{err, "failed to list stakes", http.StatusInternalServerError}
		}

		if err := json.NewEncoder(w).Encode(ss); err != nil {
			return &apiError{err, "failed to encode response", http.StatusInternalServerError}
		}

		return nil
	}
}
