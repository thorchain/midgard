package main

import "net/http"

type Pool struct {
	Token                    string  `json:"token"`
	Depth                    float64 `json:"depth"`
	Volume24Hour             float64 `json:"volume_24h"`
	AverageTransactionAmount float64 `json:"avg_tx_amount"`
	AverageLiquidityFee      float64 `json:"avg_liquidity_fee"`
	HistoricalROI            float64 `json:"historical_roi"`
}

func listPools() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}

func getPool() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request) *apiError {
		return nil
	}
}
