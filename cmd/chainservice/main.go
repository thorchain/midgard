package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	server "gitlab.com/thorchain/bepswap/chain-service/server"
)

func main() {

	cfgFile := flag.StringP("cfg", "c", "config", "configuration file with extension")
	flag.Parse()

	s, err := server.New(cfgFile)
	if err != nil {
		log.Fatal("fail to create service: ", err)
	}

	if err := s.Start(); err != nil {
		log.Fatal("fail to start server: ", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	s.Log().Info().Msg("stop signal received")
	if err := s.Stop(); nil != err {
		s.Log().Fatal().Err(err).Msg("fail to stop chain service")
	}
}
