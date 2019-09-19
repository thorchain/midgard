package coingecko

import (
	"testing"
	"time"
)

func TestNewPriceService(t *testing.T) {
	ps := NewPriceService(NewCache(), "thorchain", "usd")
	ch := make(chan struct{}, 1)
	go ps.Run(time.Second*3, ch)
	time.Sleep(time.Second * 4)
	resp, err := ps.GetPrice()
	if err != nil {
		t.Error(err)
	}
	if resp.CoinName == "" || resp.Price == 0 || resp.CurrencyName == "" {
		t.Error("empty data")
	}
	close(ch)
}
