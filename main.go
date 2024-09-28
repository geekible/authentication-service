package main

import (
	"authservice/src/config"
	"authservice/src/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	serviceCfg := config.InitServiceConfig()

	dbMigration := config.InitDatabaseMigration(serviceCfg.Db)
	if err := dbMigration.DoMigration(); err != nil {
		log.Fatalf("error performing migration: %v", err)
	}

	// register routes
	routes.InitAuthenticationRoutes(serviceCfg).Register()
	routes.InitClaimRoutes(*serviceCfg).Register()

	log.Printf("starting service on port: %d\n", serviceCfg.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serviceCfg.Port), serviceCfg.Mux); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}
