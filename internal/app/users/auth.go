package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SIGNING_METHOD jwt.SigningMethod = jwt.SigningMethodHS512

const TOKEN_VALID_DURATION = time.Hour * 3

func createToken(secret []byte, claims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(SIGNING_METHOD, claims).SignedString(secret)
}

func createUserToken(secret []byte, username string) (string, error) {
	d := time.Now().UTC()

	return createToken(secret, jwt.MapClaims{
		"sub": username,
		"iat": d.Unix(),
		"exp": d.Add(TOKEN_VALID_DURATION).Unix(),
	})
}

func createAnonToken(secret []byte) (string, error) {
	d := time.Now().UTC()

	return createToken(secret, jwt.MapClaims{
		"iat": d.Unix(),
		"exp": d.Add(TOKEN_VALID_DURATION).Unix(),
	})
}

// Logs the user in and gives them a JWT they can then use.
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Creating new JWT")

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
	h.logger.Infoln("Creating new anonymous JWT")

	token, err := createAnonToken([]byte(h.secret))
	if err != nil {
		h.logger.Errorf("Failed to create anon JWT: %s\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))
}

func decodeToken(secret []byte, tokenString string) (*jwt.Token, error) {
	// FIXME: set valid methods
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if time.Now().After(exp.Time) {
		return nil, errors.New("Expiry time has already passed")
	}

	return token, nil
}

func extractTokenFromHeader(header string) string {
	split := strings.SplitN(header, "Bearer", 2)
	if len(split) < 2 {
		return ""
	}

	return strings.TrimSpace(split[1])
}

// Verifies that the token is valid
func (h *handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Verifying existing JWT")

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No auth header"))

		return
	}

	token := extractTokenFromHeader(authHeader)
	if len(token) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No token supplied"))
		return
	}

	if _, err := decodeToken([]byte(h.secret), token); err != nil {
		h.logger.Errorf("Unable to parse token: %s\n", err)

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// verifies a token's validity and returns its claims if it is valid
func (h *handler) VerifyTokenAndGetClaims(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Decoding JWT token")

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No auth header"))

		return
	}

	t := extractTokenFromHeader(authHeader)
	token, err := decodeToken([]byte(h.secret), t)
	if err != nil {
		h.logger.Infoln("User is not authorized")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	j, err := json.Marshal(token.Claims)
	if err != nil {
		h.logger.Errorf("Unable to serialize token claims: %s\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
