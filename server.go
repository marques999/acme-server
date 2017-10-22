package main

import (
	_ "github.com/lib/pq"
	"github.com/marques999/acme-server/application"
)

func main() {
	application.Run()
}
