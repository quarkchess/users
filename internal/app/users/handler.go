package users

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/stanekondrej/logger"
	"github.com/stanekondrej/quarkchess/users/pkg/users"
)

type handler struct {
	db     users.Database
	logger *logger.Logger

	secret string
}

func genSecret() string {
	const BUF_SIZE = 64

	buf := make([]byte, BUF_SIZE)
	rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func NewHandler(connstring string) (handler, error) {
	logger := logger.NewLogger("HANDLER")
	logger.Infoln("Initializing handler")

	db, err := users.NewDatabase(connstring)
	if err != nil {
		return handler{}, err
	}

	secret := genSecret()
	logger.Infof("Secret is set to %s\n", secret)

	return handler{
		db,
		&logger,
		secret,
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
