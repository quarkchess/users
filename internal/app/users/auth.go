package users

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SIGNING_METHOD jwt.SigningMethod = jwt.SigningMethodHS512

const TOKEN_VALID_DURATION = time.Hour * 3

func createToken(secret []byte, claims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(SIGNING_METHOD, claims).SignedString(secret)
}

func createUserToken(secret []byte, username string) (string, error) {
	return createToken(secret, jwt.MapClaims{
		"sub": username,
		"iat": time.Now().UTC(),
		"exp": time.Now().UTC().Add(TOKEN_VALID_DURATION),
	})
}

func createAnonToken(secret []byte) (string, error) {
	return createToken(secret, jwt.MapClaims{
		"iat": time.Now().UTC(),
		"exp": time.Now().UTC().Add(TOKEN_VALID_DURATION),
	})
}

// Logs the user in and gives them a JWT they can then use.
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorln("Unable to read request body")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body struct {
		Username string
		Password string
	}

	if err := json.Unmarshal(b, &body); err != nil {
		h.logger.Errorln("Body isn't valid JSON")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := h.db.VerifyPassword(body.Username, body.Password)
	if err != nil {
		h.logger.Errorf("Unable to verify password of user %s\n", body.Username)
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		h.logger.Errorln("User isn't authorized")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := createUserToken([]byte(h.secret), body.Username)
	if err != nil {
		h.logger.Errorf("Failed to sign JWT: %s\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))
}

// Gives the user an anonymous JWT that they can use for e.g. a persistent session thing
func (h *handler) LoginAnon(w http.ResponseWriter, r *http.Request) {
	token, err := createAnonToken([]byte(h.secret))
	if err != nil {
		h.logger.Errorf("Failed to create anon JWT: %s\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))
}
