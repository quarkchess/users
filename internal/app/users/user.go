package users

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Getting user")

	qUsername := r.URL.Query().Get("username")
	if len(qUsername) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No username specified"))

		return
	}

	user, err := h.db.GetUser(qUsername)
	if err != nil {
		h.logger.Errorln("Requested user not found")
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))

		return
	}

	d, err := json.Marshal(user)
	if err != nil {
		h.logger.Errorf("Unable to marshal user into JSON: %+v\n", user)
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("Creating user")

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorln("Unable to read body")
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := struct {
		Username string
		Password string
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		h.logger.Errorf("Body isn't valid json: %s\n", b)
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))

		return
	}

	u, err := h.db.CreateUser(body.Username, body.Password)
	if err != nil {
		h.logger.Errorf("Unable to create user: %s\n", err)
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to create user"))

		return
	}

	h.logger.Infof("Created user %s\n", u.Username)

	j, err := json.Marshal(u)
	if err != nil {
		h.logger.Errorf("User (%+v) is not valid json\n", u)
		h.logger.Errorln(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
