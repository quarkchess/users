package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stanekondrej/quarkchess/auth/pkg/auth"
	"github.com/stanekondrej/quarkchess/auth/pkg/auth/util"
)

type handler struct {
	db     auth.Database
	logger *log.Logger
}

func NewHandler(connstring string) (handler, error) {
	logger := util.NewLogger("HANDLER")
	logger.Println("Initializing handler")

	db, err := auth.NewDatabase(connstring)
	if err != nil {
		return handler{}, err
	}

	return handler{
		db,
		logger,
	}, nil
}

func (h *handler) GetVersion(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting version")

	d := struct {
		Version string `json:"version"`
	}{Version: "1.0"}

	b, err := json.Marshal(d)
	if err != nil {
		panic("version code is wrong")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
