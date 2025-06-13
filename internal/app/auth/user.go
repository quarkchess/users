package auth

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting user")

	qUsername := r.URL.Query().Get("username")
	if len(qUsername) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No username specified"))

		return
	}

	user, err := h.db.GetUser(qUsername)
	if err != nil {
		h.logger.Println("Requested user not found")
		h.logger.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))

		return
	}

	d, err := json.Marshal(user)
	if err != nil {
		h.logger.Printf("Unable to marshal user into JSON: %+v", user)
		h.logger.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Creating user")

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Println("Unable to read body")
		h.logger.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := struct {
		Username string
		Password string
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		h.logger.Println("Body isn't valid json:", b)
		h.logger.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))

		return
	}

	u, err := h.db.CreateUser(body.Username, body.Password)
	if err != nil {
		h.logger.Println("Unable to create user:", err)
		h.logger.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to create user"))

		return
	}

	h.logger.Printf("Created user %s\n", u.Username)

	j, err := json.Marshal(u)
	if err != nil {
		h.logger.Printf("User (%+v) is not valid json\n", u)
		h.logger.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
