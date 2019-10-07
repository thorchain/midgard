package main

import (
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	"gitlab.com/thorchain/bepswap/chain-service/server/echo"
)

func main() {

	cfgFile := flag.StringP("cfg", "c", "config", "configuration file with extension")
	flag.Parse()

	s, err := echo.New(cfgFile)
	if err != nil {

	}

	if err := s.Start(); err != nil {

	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	// log.Info().Msg("stop signal received")
	if err := s.Stop(); nil != err {
		// log.Fatal().Err(err).Msg("fail to stop chain service")
	}
}
