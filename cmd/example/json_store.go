package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
)

func NewPoolStoreFromJSON(r io.Reader) (PoolStore, error) {
	ps := []*Pool{}
	psMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&ps); err != nil {
		return nil, err
	}

	return &StubPoolStore{
		ListFn: func(ctx context.Context) ([]*Pool, error) {
			psMu.Lock()
			defer psMu.Unlock()

			return ps, nil
		},
		GetFn: func(ctx context.Context, sym string) (*Pool, error) {
			psMu.Lock()
			defer psMu.Unlock()

			for _, p := range ps {
				if p.Symbol == sym {
					return p, nil
				}
			}

			return nil, ErrNotFound
		},
		SetFn: func(ctx context.Context, p *Pool) error {
			return errors.New("NOT IMPLEMENTED")
		},
	}, nil
}

func NewTokenStoreFromJSON(r io.Reader) (TokenStore, error) {
	ts := []*Token{}
	tsMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&ts); err != nil {
		return nil, err
	}

	return &StubTokenStore{
		ListFn: func(ctx context.Context) ([]*Token, error) {
			tsMu.Lock()
			defer tsMu.Unlock()

			return ts, nil
		},
		GetFn: func(ctx context.Context, sym string) (*Token, error) {
			tsMu.Lock()
			defer tsMu.Unlock()

			for _, t := range ts {
				if t.Symbol == sym {
					return t, nil
				}
			}

			return nil, ErrNotFound
		},
		SetFn: func(ctx context.Context, t *Token) error {
			return errors.New("NOT IMPLEMENTED")
		},
	}, nil
}

func NewPriceStoreFromJSON(r io.Reader) (PriceStore, error) {
	ps := []*Price{}
	psMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&ps); err != nil {
		return nil, err
	}

	return &StubPriceStore{
		ListFn: func(ctx context.Context) ([]*Price, error) {
			psMu.Lock()
			defer psMu.Unlock()

			return ps, nil
		},
		GetFn: func(ctx context.Context, sym string) (*Price, error) {
			psMu.Lock()
			defer psMu.Unlock()

			for _, p := range ps {
				if p.Symbol == sym {
					return p, nil
				}
			}

			return nil, ErrNotFound
		},
		SetFn: func(ctx context.Context, p *Price) error {
			return errors.New("NOT IMPLEMENTED")
		},
	}, nil
}

func NewSwapStoreFromJSON(r io.Reader) (SwapStore, error) {
	ss := []*Swap{}
	ssMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&ss); err != nil {
		return nil, err
	}

	return &StubSwapStore{
		ListFn: func(ctx context.Context) ([]*Swap, error) {
			ssMu.Lock()
			defer ssMu.Unlock()

			return ss, nil
		},
		GetFn: func(ctx context.Context, sym string) (*Swap, error) {
			ssMu.Lock()
			defer ssMu.Unlock()

			for _, s := range ss {
				if s.Symbol == sym {
					return s, nil
				}
			}

			return nil, ErrNotFound
		},
		SetFn: func(ctx context.Context, s *Swap) error {
			return errors.New("NOT IMPLEMENTED")
		},
	}, nil
}

func NewStakeStoreFromJSON(r io.Reader) (StakeStore, error) {
	ss := []*Stake{}
	ssMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&ss); err != nil {
		return nil, err
	}

	return &StubStakeStore{
		ListFn: func(ctx context.Context, addr string) ([]*Stake, error) {
			ssMu.Lock()
			defer ssMu.Unlock()

			return ss, nil
		},
	}, nil
}

func NewNetworkStoreFromJSON(r io.Reader) (NetworkStore, error) {
	n := &Network{}
	nMu := &sync.Mutex{}

	if err := json.NewDecoder(r).Decode(&n); err != nil {
		return nil, err
	}

	return &StubNetworkStore{
		GetFn: func(ctx context.Context) (*Network, error) {
			nMu.Lock()
			defer nMu.Unlock()

			return n, nil
		},
		SetFn: func(ctx context.Context, s *Network) error {
			return errors.New("NOT IMPLEMENTED")
		},
	}, nil
}
