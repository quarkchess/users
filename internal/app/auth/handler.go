package auth

import (
	"encoding/json"
	"net/http"

	"github.com/stanekondrej/quarkchess/auth/pkg/auth"
)

type handler struct {
	db auth.Database
}

func NewHandler(connstring string) (handler, error) {
	db, err := auth.NewDatabase(connstring)
	if err != nil {
		return handler{}, err
	}

	return handler{
		db,
	}, nil
}

func (h *handler) GetVersion(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	qUsername := r.URL.Query().Get("username")
	if len(qUsername) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No username specified"))

		return
	}

	user, err := h.db.GetUser(qUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))

		return
	}

	d, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(d)
}
