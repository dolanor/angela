package main

import "flag"

type config struct {
	port      int
	logFormat string
}

func loadConfig(args []string) (config, error) {
	cfg := config{}

	fs := flag.NewFlagSet("serverConfig", flag.ContinueOnError)
	fs.IntVar(&cfg.port, "port", 7777, "port of the HTTP server")
	fs.StringVar(&cfg.logFormat, "log-format", "text", "log format (eg: values, json)")

	err := fs.Parse(args)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
