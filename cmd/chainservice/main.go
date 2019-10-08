package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"gitlab.com/thorchain/bepswap/chain-service/config"
	chain_service "gitlab.com/thorchain/bepswap/chain-service/server/gin"
)

// we define version / revision here , so we could inject the version from CI pipeline if we want to
var (
	version  string
	revision string
)

const (
	serverIdentity = "chain-service"
)

func printVersion() {
	fmt.Printf("%s v%s, rev %s\n", serverIdentity, version, revision)
}

func initLog(level string, pretty bool) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Msgf("%s is not a valid log-level, falling back to 'info'", level)
	}
	var out io.Writer = os.Stdout
	if pretty {
		out = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	zerolog.SetGlobalLevel(l)
	log.Logger = log.Output(out).With().Str("service", serverIdentity).Logger()
}

func main() {
	showVersion := flag.Bool("version", false, "Shows version")
	// TODO set the default log level to info later
	logLevel := flag.StringP("log-level", "l", "debug", "Log Level")
	pretty := flag.BoolP("pretty-log", "p", false, "Enables unstructured prettified logging. This is useful for local debugging")
	cfgFile := flag.StringP("cfg", "c", "config", "configuration file with extension")
	flag.Parse()
	if *showVersion {
		printVersion()
		return
	}
	initLog(*logLevel, *pretty)
	cfg, err := config.LoadConfiguration(*cfgFile)
	if nil != err {
		log.Fatal().Err(err).Msg("fail to load chain service config ")
	}
	s, err := chain_service.NewServer(*cfg)
	if nil != err {
		log.Fatal().Err(err).Msg("fail to create chain service")
	}
	if err := s.Start(); nil != err {
		log.Fatal().Err(err).Msg("fail to start server")
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")
	if err := s.Stop(); nil != err {
		log.Fatal().Err(err).Msg("fail to stop chain service")
	}
}
