package main

import (
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1"
)

func main() {
	config.Load()
	postgres.NewDatabaseConnection()
	httpv1.StartServer()
}
