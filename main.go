package main

import (
	"log"

	"github.com/nedpals/sulatcms/server"
	"github.com/nedpals/sulatcms/sulat"
	_ "modernc.org/sqlite"
)

func main() {
	rootInst, err := sulat.NewInstance("sulat.db")
	if err != nil {
		log.Fatalf("failed to initialize: %s\n", err)
	}

	rootInst.RegisterDataSourceProvider(&sulat.FileDataSourceProvider{})

	if err := server.Start(rootInst, "3000"); err != nil {
		log.Fatalf("failed to start server: %s\n", err)
	}
}
