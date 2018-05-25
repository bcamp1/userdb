package userdb

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const testfile = "usersDB_test.txt"

func TestGenHash(t *testing.T) {
	str := "Hello World"
	expectedHash := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
	if genHash(str) != expectedHash {
		t.Errorf("Expected %v, got %v\n", expectedHash, genHash(str))
	}
}

func TestHashFile(t *testing.T) {
	fname := "test-file-for-hash.txt"
	contents := "this is a test"
	ioutil.WriteFile(fname, []byte(contents), 0644)
	expectedHash := "2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c"
	gotHash := hashFile(fname)

	if gotHash != expectedHash {
		t.Errorf("Expected %v, got %v\n", expectedHash, gotHash)
	}
}

func TestSave(t *testing.T) {
	m := map[string]User{
		"bcamp1": User{
			"bcamp1", "salt1", "hash1",
		},
		"bcamp2": User{
			"bcamp2", "salt2", "hash2",
		},
	}
	expected := "bcamp1\nsalt1\nhash1\nbcamp2\nsalt2\nhash2\n"

	db := DB(m)
	err := db.Save(testfile)
	contents, err := ioutil.ReadFile(testfile)
	if err != nil {
		panic(err)
	}

	result := string(contents)

	if result != expected {
		t.Errorf("Expected: %vGot: %v", expected, result)
	}
}

func TestLoad(t *testing.T) {
	expected := map[string]User{
		"bcamp1": User{
			"bcamp1", "salt1", "hash1",
		},
		"bcamp2": User{
			"bcamp2", "salt2", "hash2",
		},
	}

	var db DB
	err := db.Load(testfile)
	if err != nil {
		panic(err)
	}

	result := map[string]User(db)
	if fmt.Sprint(result) != fmt.Sprint(expected) {
		t.Errorf("Expected:\n%vGot:\n%v", expected, result)
	}
}

func TestGenSalt(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(genSalt(50))
	}
}

func TestCreateUser(t *testing.T) {
	uname := "bcamp3"
	pw := "mypass"
	var db DB
	err := db.Load("test")
	t.Log(db)
	if err != nil {
		panic(err)
	}
	err = db.CreateUser(uname, pw)
	if err != nil {
		t.Log(err)
	}
	db.Save("test")
	t.Log(db)
}

func TestValidateLogin(t *testing.T) {
	var db DB
	db.Load("test")
	err := db.validateLogin("bcamp2345", "pw12345")
	if err != USERNAME_DOES_NOT_EXIST {
		t.Fail()
	}
	err = db.validateLogin("bcamp3", "wrongpassword")
	if err != INCORRECT_PASSWORD {
		t.Fail()
	}
	err = db.validateLogin("bcamp3", "mypass")
	if err != nil {
		t.Fail()
	}
}

func TestNewDB(t *testing.T) {
	db := NewDB("test")
	if !db.UserExists("bcamp3") {
		t.Fail()
	}

	db2 := NewDB("testnew")
	t.Log(db2)
}
