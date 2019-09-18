package coingecko

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	coingecko "github.com/superoo7/go-gecko/v3"
)

type runner interface {
	Run(timeDuration time.Duration)
}

type PriceServiceResponse struct {
	CoinName     string
	CurrencyName string
	Price        float32
}

type PriceService struct {
	cache          *Cache
	logger         zerolog.Logger
	cgClient       *coingecko.Client
	id, vsCurrency string
}

func NewPriceService(cache *Cache, id, vsCurrency string) *PriceService {
	hc := &http.Client{
		Timeout: time.Second * 5,
	}

	return &PriceService{
		cache:      cache,
		cgClient:   coingecko.NewClient(hc),
		id:         id,
		vsCurrency: vsCurrency,
		logger:     log.With().Str("module", "priceservice").Logger(),
	}
}

func (ps *PriceService) GetPrice() (*PriceServiceResponse, error) {
	res, err := ps.cache.Get(ps.id)
	if err != nil {
		return nil, err
	}
	result, ok := res.(*PriceServiceResponse)
	if !ok {
		return nil, fmt.Errorf("id %s not found", ps.id)
	}
	return result, nil
}

func (ps *PriceService) Run(timeDuration time.Duration) {
	err := ps.setPrice()
	if err != nil {
		ps.logger.Log().Str("on start price service error", err.Error())
	}
	for tick := range time.Tick(timeDuration) {
		err := ps.setPrice()
		if err != nil {
			ps.logger.Log().Str(tick.String()+"price service error", err.Error())
			continue
		}
		ps.logger.Log().Str(tick.String()+"price service", "updated")
	}
}

func (ps *PriceService) setPrice() error {
	res, err := ps.cgClient.SimplePrice([]string{ps.id}, []string{ps.vsCurrency})
	result := (*res)[ps.id]

	ps.cache.Set(ps.id, &PriceServiceResponse{
		CoinName:     ps.id,
		CurrencyName: ps.vsCurrency,
		Price:        result[ps.vsCurrency],
	})
	return err
}

type Cache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data:  make(map[string]interface{}),
		mutex: sync.RWMutex{},
	}
}

func (cp *Cache) Set(key string, val interface{}) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.data[key] = val
}

func (cp *Cache) Get(key string) (interface{}, error) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	var res interface{}
	res, ok := cp.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return res, nil
}
