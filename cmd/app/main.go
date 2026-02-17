package main

import (
	"github.com/nurgal1ev/yotabo-go/internal/database/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1"
)

func main() {
	postgres.NewDatabaseConnection()
	httpv1.StartServer()
}
