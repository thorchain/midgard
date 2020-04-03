[![pipeline status](https://gitlab.com/thorchain/midgard/badges/master/pipeline.svg)](https://gitlab.com/thorchain/midgard/commits/master)

Midgard API 

****

> **Mirror**
>
> This repo mirrors from THORChain Gitlab to Github. 
> To contribute, please contact the team and commit to the Gitlab repo:
>
> https://gitlab.com/thorchain/midgard


****

=============

### Run Midgard API
To run Midgard you will need two terminal windows or tabs. In the
first tab, run...
```bash
make pg
# create an user (if you have not already done it before)
make create-user
# create a database (if you have not already done it before)
make create-database
```

In the second tab, run...
```bash
make install run
```

### Run mock server
To use a mock server run everything as described in `Run Midgard API`. After that, run following command in another terminal:

```bash
make run-thormock
```

### Run generated specs locally
First, run everything as described in `Run chain service` and `Run mock server` by using different terminals.

Open  http://127.0.0.1:8080/v1/doc in your browser.



### Testing
```bash
make test
```

For rapid testing, in one terminal tab...
```bash
make pg
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
timescale running on top of postgres
```bash
make test-short
```

### Linting
```bash
make lint
```
