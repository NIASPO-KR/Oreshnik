package main

import (
	"flag"
	"log"

	"oreshnik/config"
	"oreshnik/internal/server"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ReadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
