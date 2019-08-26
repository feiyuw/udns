package main

import (
	"udns/logger"
)

func main() {
	addr := ":53"
	logger.Infof("Listen to UDP '%s'", addr)
}
