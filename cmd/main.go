package main

import (
	"abuse/server"
	"flag"
)

func main() {
	var envType string
	flag.StringVar(&envType, "env", "local", "Environment type: production, staging, dev, local")
	flag.Parse()
	server.Run(&envType)
}
