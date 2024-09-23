// simple http server: https://dev.to/andyjessop/building-a-basic-http-server-in-go-a-step-by-step-tutorial-ma4

package main

import (
	"log"

	"main.go/db"
	"main.go/internal/user"
	"main.go/internal/ws"
	"main.go/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}
	defer dbConn.Close()

	err = dbConn.MigrateDB("UP")
	if err != nil {
		log.Fatalf("Could not migrate the database UP: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")

	// err = dbConn.MigrateDB("DOWN")
	// if err != nil {
	// 	log.Fatalf("Could not migrate the database DOWN: %s", err)
	// }

}
