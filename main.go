package main

import (
	"log"

	rest "github.com/nedpals/sulatcms/rest_server"
	"github.com/nedpals/sulatcms/sulat"
	_ "modernc.org/sqlite"
)

func main() {
	rootInst, err := sulat.NewInstance("sulat.db")
	if err != nil {
		log.Fatalf("failed to initialize: %s\n", err)
	}

	rootInst.RegisterDataSourceProvider(&sulat.FileDataSourceProvider{})

	if err := rest.StartServer(rootInst, "3000"); err != nil {
		log.Fatalf("failed to start server: %s\n", err)
	}
}
