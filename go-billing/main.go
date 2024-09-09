package main

import (
	"embed"
	"flag"

	"github.com/sanderdescamps/go-billing-api/api"
)

//go:embed swagger/*
var swaggerDir embed.FS

func main() {
	dataDir := flag.String("data", "./data", "dir to store json db files")
	httpPort := flag.Int("port", 8080, "http port")
	host := flag.String("host", "0.0.0.0", "http host")
	flag.Parse()

	api.InitDB(*dataDir)
	api.InitSwagger(swaggerDir)
	api.Run(*host, *httpPort)
}
