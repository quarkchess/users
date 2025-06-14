package main

import (
	"net/http"
	"os"

	"github.com/stanekondrej/logger"
	"github.com/stanekondrej/quarkchess/users/internal/app/users"
)

func main() {
	logger := logger.NewLogger("MAIN")
	logger.Infoln("Starting auth server")

	listenAddress, ok := os.LookupEnv("LISTEN_ADDRESS")
	if !ok {
		logger.Fatalln("LISTEN_ADDRESS not specified")
	}

	connstring, ok := os.LookupEnv("DB_CONNSTRING")
	if !ok {
		logger.Fatalln("DB_CONNSTRING not specified")
	}

	h, err := users.NewHandler(connstring)
	if err != nil {
		logger.Fatalln("Failed to create handler:", err)
	}

	http.HandleFunc("GET /", h.GetVersion)
	http.HandleFunc("GET /verify_token", h.VerifyToken)
	http.HandleFunc("GET /get_user", h.GetUser)
	http.HandleFunc("POST /login", h.Login)
	http.HandleFunc("POST /login_anon", h.LoginAnon)
	http.HandleFunc("POST /create_user", h.CreateUser)

	logger.Fatalln(http.ListenAndServe(listenAddress, nil))
}
