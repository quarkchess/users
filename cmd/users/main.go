package main

import (
	"net/http"
	"os"

	"github.com/stanekondrej/quarkchess/users/internal/app/users"
	"github.com/stanekondrej/quarkchess/users/pkg/users/util"
)

func main() {
	logger := util.NewLogger("MAIN")
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

	http.HandleFunc("/", h.GetVersion)
	http.HandleFunc("POST /login", h.Login)
	http.HandleFunc("POST /login_anon", h.LoginAnon)
	http.HandleFunc("/get_user", h.GetUser)
	http.HandleFunc("POST /create_user", h.CreateUser)

	logger.Fatalln(http.ListenAndServe(listenAddress, nil))
}
