package auth

import (
	"encoding/json"
	"io"
	"log"
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
		log.Println("Requested user not found")
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))

		return
	}

	d, err := json.Marshal(user)
	if err != nil {
		log.Printf("Unable to marshal user into JSON: %+v", user)
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read body")
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := struct {
		Username string
		Password string
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		log.Println("Body isn't valid json:", b)
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))

		return
	}

	u, err := h.db.CreateUser(body.Username, body.Password)
	if err != nil {
		log.Println("Unable to create user:", err)
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to create user"))

		return
	}

	j, err := json.Marshal(u)
	if err != nil {
		log.Printf("User (%+v) is not valid json\n", u)
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
