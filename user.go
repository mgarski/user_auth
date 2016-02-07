package main

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	_ "github.com/lib/pq"
)

type User struct {
	Name, Email, Password string
	Id int
}

func (u User) Register() (int, string) {
	salt := generateSalt()
	if(salt == nil) {
		return 500, "Internal error"
	}

	hash := generateHash(u.Password, salt)

	db, err := sql.Open("postgres", config.DbConnection)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(email,name,hash,salt) VALUES($1,$2,$3,$4);")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return 500, "Internal Error"
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Email, u.Name, string(hash), string(salt))
	if err != nil {
		fmt.Printf("ERROR execute: %s\n", err.Error())
		return 500, "Internal Error"
	}
	return 200, ""
}

func (u User) ValidateCredentials() int  {

	db, err := sql.Open("postgres", config.DbConnection)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, salt, hash FROM users WHERE email = $1;")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return -1
	}
	row := stmt.QueryRow(u.Email)
	var id int
	var dbHash string
	var dbSalt string

	err = row.Scan(&id, &dbSalt, &dbHash)
	if err != nil {
		fmt.Println(err.Error())
	}

	var hashedPass []byte = []byte(dbHash)
	var salt []byte = []byte(dbSalt)

	incoming := generateHash(u.Password, salt)
	if len(hashedPass) != len(incoming) {
		return -1
	}
	for i := 0; i < len(hashedPass); i++ {
		if hashedPass[i] != incoming[i] {
			return -1
		}
	}
	return id
}

func (u User) Update() (int, string) {
	salt := generateSalt()
	if(salt == nil) {
		return 500, "Internal error"
	}

	hash := generateHash(u.Password, salt)

	db, err := sql.Open("postgres", config.DbConnection)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users set email = $1, name = $2, hash = $3, salt = $4 where id = $5;")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return 500, "Internal Error"
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Email, u.Name, string(hash), string(salt), u.Id)
	if err != nil {
		fmt.Printf("ERROR execute: %s\n", err.Error())
		return 500, "Internal Error"
	}
	return 200, "updated"
}

func (u User) Delete() (int, string) {
	db, err := sql.Open("postgres", config.DbConnection)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1;")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return 500, "Internal Error"
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Id)
	if err != nil {
		fmt.Printf("ERROR execute: %s\n", err.Error())
		return 500, "Internal Error"
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("ERROR affected: %s\n", err.Error())
		return 500, "Internal Error"
	}
	if affect > 0 {
		return 200, "deleted"
	} else {
		return 404, "User not found"
	}

}

func generateSalt() []byte {
	length := 10
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	} else {
		return salt
	}
}

func generateHash(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)
}