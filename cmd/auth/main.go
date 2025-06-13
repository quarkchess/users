package main

import (
	"log"
	"net/http"
	"os"

	auth "github.com/stanekondrej/quarkchess/auth/internal/app/auth"
)

func main() {
	log.Println("Starting auth server")

	listenAddress, ok := os.LookupEnv("LISTEN_ADDRESS")
	if !ok {
		log.Fatalln("No listen address specified")
	}

	connstring, ok := os.LookupEnv("DB_CONNSTRING")
	if !ok {
		log.Fatalln("No database connstring specified")
	}

	h, err := auth.NewHandler(connstring)
	if err != nil {
		log.Fatalln("Failed to create handler:", err)
	}

	http.HandleFunc("/", h.GetVersion)
	http.HandleFunc("/get_user", h.GetUser)
	http.HandleFunc("POST /create_user", h.CreateUser)

	http.ListenAndServe(listenAddress, nil)
}
