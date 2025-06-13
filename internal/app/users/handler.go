package users

import (
	"encoding/json"
	"net/http"

	"github.com/stanekondrej/quarkchess/users/pkg/users"
	"github.com/stanekondrej/quarkchess/users/pkg/users/util"
)

type handler struct {
	db     users.Database
	logger *util.Logger
}

func NewHandler(connstring string) (handler, error) {
	logger := util.NewLogger("HANDLER")
	logger.Infoln("Initializing handler")

	db, err := users.NewDatabase(connstring)
	if err != nil {
		return handler{}, err
	}

	return handler{
		db,
		&logger,
	}, nil
}

func (h *handler) GetVersion(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Getting version")

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
