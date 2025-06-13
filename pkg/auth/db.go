package auth

import (
	"database/sql"
	"errors"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stanekondrej/quarkchess/auth/pkg/auth/util"
)

type User struct {
	Username string
	Elo      uint
}

type Database struct {
	inner *sql.DB
}

func NewDatabase(connstring string) (Database, error) {
	db, err := sql.Open("sqlite", connstring)
	if err != nil {
		return Database{}, err
	}
	util.Assert(db != nil)

	if err := db.Ping(); err != nil {
		return Database{}, errors.New("Failed to connect to database")
	}

	return Database{
		inner: db,
	}, nil
}

func (d *Database) GetUser(username string) (User, error) {
	row, err := d.inner.Query("SELECT * FROM users WHERE username = ? LIMIT 1;", username)
	if err != nil {
		return User{}, err
	}

	var u User
	err = row.Scan(&u)
	if err != nil {
		panic("User struct doesn't have the correct amount of fields")
	}

	return u, nil
}
