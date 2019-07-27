package main

import (
	"flag"
	"persona/config"
	"persona/server"
)

func main() {
	environment := flag.String("e", config.EnvironmentDev, "")
	flag.Parse()
	config.Init(*environment)

	r := server.NewRouter()

	// TODO: provide port and set TLS/non TLS based on environment(config)
	r.Run()
}
