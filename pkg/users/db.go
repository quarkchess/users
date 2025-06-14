package users

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stanekondrej/quarkchess/users/pkg/users/util"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Elo          uint   `json:"elo"`
}

type Database struct {
	inner  *sql.DB
	logger *util.Logger
}

// FIXME: REMOVE THIS AS SOON AS POSSIBLE
func initDb(db *sql.DB, logger *util.Logger) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT UNIQUE, password_hash TEXT, elo INTEGER);")
	if err != nil {
		logger.Fatalln("Unable to init database")
	}
}

func NewDatabase(connstring string) (Database, error) {
	logger := util.NewLogger("DB")

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

func (d *Database) GetUser(username string) (User, error) {
	d.logger.Infof("Getting user %s\n", username)

	row := d.inner.QueryRow("SELECT username, password_hash, elo FROM users WHERE username = ? LIMIT 1;", username)

	var u User
	err := row.Scan(&u.Username, &u.PasswordHash, &u.Elo)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func hashPassword(password string) string {
	sum := sha512.Sum512([]byte(password))
	return hex.EncodeToString(sum[:])
}

const DEFAULT_ELO uint = 1000

func (d *Database) CreateUser(username string, password string) (User, error) {
	d.logger.Infof("Creating user %s\n", username)

	hash := hashPassword(password)
	_, err := d.inner.Exec("INSERT INTO users (username, password_hash, elo) VALUES (?, ?, ?);", username, hash, DEFAULT_ELO)
	if err != nil {
		return User{}, err
	}

	// TODO: query the row from the database (?)
	return User{
		username,
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
