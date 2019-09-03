[![pipeline status](https://gitlab.com/thorchain/bepswap/chain-service/badges/master/pipeline.svg)](https://gitlab.com/thorchain/bepswap/chain-service/commits/master)
[![coverage report](https://gitlab.com/thorchain/bepswap/chain-service/badges/master/coverage.svg)](https://gitlab.com/thorchain/bepswap/chain-service/commits/master)

Chain Service 
=============

### Run chain service
To run the chain service you will need two terminal windows or tabs. In the
first tab, run...
```bash
make influxdb
```

In the second tab, run...
```bash
make install
```


### Testing
```bash
make test
```

For rapid testing, in one terminal tab...
```bash
make influxdb
```

In another tab, run...
```bash
make test-internal
```

If you'd like to run tests everytime there is a change to a go file...
```bash
make test-watch
```

#### Short Testing
You can run unit tests and omit the ones that require a running instance of
influxdb
```bash
make test-short
```

### Linting
```bash
make lint
```
