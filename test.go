package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var db *sql.DB

func createGreetings() {
	_, err := db.Exec(
		"create table if not exists greetings (greeting text)",
		nil,
	)
	check(err)
}

func insertGreeting(greeting string) {
	transaction, err := db.Begin()
	check(err)
	statement, err := transaction.Prepare(
		"insert into greetings (greeting) values (?)",
	)
	check(err)
	defer statement.Close()
	_, err = statement.Exec(greeting)
	check(err)
	transaction.Commit()
}

func selectGreetings() []string {
	rows, err := db.Query("select greeting from greetings")
	check(err)
	defer rows.Close()
	var greetings []string
	for rows.Next() {
		var greeting string
		check(rows.Scan(&greeting))
		greetings = append(greetings, greeting)
	}
	check(rows.Err())
	return greetings
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "test.db")
	check(err)
	createGreetings()
	insertGreeting(time.Now().UTC().Format(time.RFC3339))
	for _, greeting := range selectGreetings() {
		fmt.Println(greeting)
	}
}
