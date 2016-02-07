package main

import (
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"database/sql"
	_ "github.com/lib/pq"
)

var secret = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

func ValidateToken(token string) bool {
	fmt.Println("validating token " + token)
	parsed, err := jws.ParseJWT([]byte(token))
	if(err != nil) {
		return false
	}

	id := int(parsed.Claims().Get("id").(float64))
	if id > 0 {
		fmt.Printf("parsed: %d : %v\n", id, parsed)

		db, err := sql.Open("postgres", dbConn)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()

		stmt, err := db.Prepare("SELECT token FROM tokens where user_id = $1;")
		if err != nil {
			fmt.Printf("ERROR prepare: %s\n", err.Error())
			return false
		}
		defer stmt.Close()

		row := stmt.QueryRow(id)
		var dbToken string
		err = row.Scan(&dbToken)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		return token == dbToken
	} else {
		return false
	}
}

func GenerateToken(id int) string {
	var claims = jws.Claims{
		"id": id,
	}
	j := jws.NewJWT(claims, crypto.SigningMethodHS256)
	t, err := j.Serialize(secret)
	if err != nil {
		//TODO: handle error
	}
	token := string(t)

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tokens(user_id,token) VALUES($1,$2);")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return ""
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, token)
	if err != nil {
		fmt.Printf("ERROR execute: %s\n", err.Error())
		return ""
	}
	return token
}

func FlushToken(id int) bool {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM tokens where user_id = $1;")
	if err != nil {
		fmt.Printf("ERROR prepare: %s\n", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Printf("ERROR execute: %s\n", err.Error())
		return false
	}
	return true
}