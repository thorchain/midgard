package timescale

import (
	"math/rand"

	"gitlab.com/thorchain/midgard/internal/store"
	. "gopkg.in/check.v1"
)

type BenchmarkSuite struct {
	Store     *Client
	generator *store.RandEventGenerator
}

var _ = Suite(&BenchmarkSuite{})

func (s *BenchmarkSuite) SetUpSuite(c *C) {
	var err error
	s.Store, err = NewTestStore(c)
	if err != nil {
		c.Fatal(err.Error())
	}

	s.generator = store.NewRandEventGenerator(&store.RandEventGeneratorConfig{
		Source:      rand.NewSource(1878939228537408224),
		Pools:       10,
		Stakers:     20,
		Swappers:    20,
		Blocks:      10000,
		AddEvents:   100,
		StakeEvents: 100,
		SwapEvents:  2000,
	})
	err = s.generator.GenerateEvents(s.Store)
	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *BenchmarkSuite) TearDownSuite(c *C) {
	err := s.Store.MigrationsDown()
	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *BenchmarkSuite) BenchmarkGetPoolData(c *C) {
	for i := 0; i < c.N; i++ {
		pool := s.generator.Pools[i%len(s.generator.Pools)]
		_, err := s.Store.GetPoolData(pool)
		c.Assert(err, IsNil)
	}
}
