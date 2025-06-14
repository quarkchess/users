package users

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stanekondrej/logger"
	"github.com/stanekondrej/quarkchess/users/pkg/users/util"
)

const DEFAULT_ELO uint = 1000

type Role uint8

const (
	Admin   Role = iota
	Regular      = iota
	Anon         = iota
)

type User struct {
	Username     string `json:"username"`
	Role         Role   `json:"role"`
	PasswordHash string `json:"password_hash"`
	Elo          uint   `json:"elo"`
}

type Database struct {
	inner  *sql.DB
	logger *logger.Logger
}

// FIXME: REMOVE THIS AS SOON AS POSSIBLE
func initDb(db *sql.DB, logger *logger.Logger) {
	// FIXME: this is crazy
	_, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS users (
			username      TEXT UNIQUE NOT NULL,
			role          INTEGER NOT NULL DEFAULT %d,
			password_hash TEXT NOT NULL,
			elo           INTEGER NOT NULL DEFAULT %d
		);`,
		Regular,
		DEFAULT_ELO),
	)
	if err != nil {
		logger.Fatalf("Unable to init database: %s\n", err)
	}
}

func NewDatabase(connstring string) (Database, error) {
	logger := logger.NewLogger("DB")

	logger.Infoln("Connecting to the database")
	db, err := sql.Open("sqlite", connstring)
	if err != nil {
		return Database{}, err
	}
	util.Assert(db != nil)

	if err := db.Ping(); err != nil {
		return Database{}, errors.New("Failed to connect to database")
	}

	initDb(db, &logger)
	logger.Infoln("Initialized the database")

	return Database{
		inner:  db,
		logger: &logger,
	}, nil
}

func (d *Database) GetUser(username string) (string, error) {
	d.logger.Infof("Getting user %s\n", username)

	row := d.inner.QueryRow("SELECT username, role, elo FROM users WHERE username = ? LIMIT 1;", username)

	// var u User
	var u struct {
		Username string `json:"username"`
		Role     Role   `json:"role"`
		Elo      uint   `json:"elo"`
	}
	err := row.Scan(&u.Username, &u.Role, &u.Elo)
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(u)
	if err != nil {
		return "", nil
	}

	return string(j), nil
}

func hashPassword(password string) string {
	sum := sha512.Sum512([]byte(password))
	return hex.EncodeToString(sum[:])
}

func (d *Database) CreateUser(username string, password string) (User, error) {
	d.logger.Infof("Creating user %s\n", username)

	hash := hashPassword(password)
	_, err := d.inner.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?);", username, hash)
	if err != nil {
		return User{}, err
	}

	// TODO: query the row from the database (?)
	return User{
		username,
		Regular,
		hash,
		DEFAULT_ELO,
	}, nil
}

func (d *Database) VerifyPassword(username, password string) (bool, error) {
	d.logger.Infof("Verifying password of user %s\n", username)

	computedHash := hashPassword(password)
	row := d.inner.QueryRow("SELECT password_hash FROM users WHERE username = ? LIMIT 1;", username)

	var expectedHash string
	if err := row.Scan(&expectedHash); err != nil {
		return false, err
	}

	return computedHash == expectedHash, nil
}
