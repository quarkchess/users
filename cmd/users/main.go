package main

import (
	"net/http"
	"os"

	"github.com/quarkchess/users/internal/app/users"
	"github.com/stanekondrej/logger"
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

	http.HandleFunc("GET /version", h.GetVersion)

	http.HandleFunc("GET /token/verify", h.VerifyToken)
	http.HandleFunc("GET /token/claims", h.VerifyTokenAndGetClaims)

	http.HandleFunc("GET /user/get", h.GetUser)
	http.HandleFunc("POST /user/auth/login", h.Login)
	http.HandleFunc("POST /user/auth/login_anon", h.LoginAnon)
	http.HandleFunc("POST /user/create", h.CreateUser)

	logger.Fatalln(http.ListenAndServe(listenAddress, nil))
}
