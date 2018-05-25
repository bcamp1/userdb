package userdb

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

type DBError string

func (err DBError) Error() string {
	return string(err)
}

// List of possible DBErrors
const (
	INCORRECT_PASSWORD      DBError = "Incorrect Password"
	USERNAME_EXISTS         DBError = "Username already exists"
	USERNAME_DOES_NOT_EXIST DBError = "Username doesn't exist"
)

type User struct {
	name string
	salt string
	hash string
}

type DB map[string]User

func (d DB) Save(fname string) error {
	var contents string
	for _, v := range d {
		contents += v.name + "\n"
		contents += v.salt + "\n"
		contents += v.hash + "\n"
	}
	return ioutil.WriteFile(fname, []byte(contents), 0644)
}

func NewDB(fname string) *DB {
	var db DB
	db.Load(fname)
	return &db
}

func (db *DB) Load(fname string) error {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		ioutil.WriteFile(fname, []byte(""), 0644)
	}
	m := make(map[string]User)
	contents, err := ioutil.ReadFile(fname)
	lines := strings.Split(string(contents), "\n")
	for i := 0; i < len(lines)-3; i += 3 {
		name := lines[i]
		salt := lines[i+1]
		hash := lines[i+2]
		u := User{name, salt, hash}
		m[name] = u
	}
	*db = DB(m)
	return err
}

func (db *DB) CreateUser(username, password string) error {
	users := *db
	_, exists := users[username]
	if exists {
		return USERNAME_EXISTS
	}
	salt := genSalt(50)
	hash := genHash(salt + password)
	u := User{username, salt, hash}
	users[username] = u
	db = &users
	return nil
}

func (db DB) validateLogin(username, password string) error {
	if !db.UserExists(username) {
		return USERNAME_DOES_NOT_EXIST
	}
	user := db[username]
	potentialHash := genHash(user.salt + password)
	if potentialHash != user.hash {
		return INCORRECT_PASSWORD
	}
	return nil
}

func (db DB) String() string {
	s := fmt.Sprintln("----------USERDB----------")

	for _, v := range map[string]User(db) {
		s += fmt.Sprintln("----------------------------------------------------------------------")
		s += fmt.Sprintf("NAME: %v\nSALT: %v\nHASH: %v\n", v.name, v.salt, v.hash)
	}
	s += fmt.Sprintln("----------------------------------------------------------------------")
	return s
}

func (db DB) UserExists(username string) bool {
	_, exists := db[username]
	return exists
}

func genHash(input string) string {
	h := sha256.New()
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hashFile(fname string) string {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	input := string(contents)
	return genHash(input)
}

func genSalt(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(33, 126))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
