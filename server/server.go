package main

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var names = []string{
	"adam",
	"bob",
	"chris",
	"diana",
	"eve",
	"fiona",
}

func getRandomName() string {
	num := len(names)
	idx := rand.IntN(num)
	return names[idx]
}

var (
	db *sql.DB
)

func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := "testcontexts"
	dsn := fmt.Sprintf("%s:%s@/%s",
		mysqlUser,
		mysqlPassword,
		dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully connected to DB", dbName)
}

func longResponse(w http.ResponseWriter, req *http.Request) {
	fmt.Println("server: longResponse started")
	defer fmt.Println("server: longResponse ended")

	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "slept 5 seconds\n")
}

// longResponseChecksContext attempts to sleep 5 seconds,
// and will end early if it detects the context has been canceled.
// Much of it is borrowed from https://gobyexample.com/context
func longResponseChecksContext(w http.ResponseWriter, req *http.Request) {
	fmt.Println("server: longResponseChecksContext started")
	defer fmt.Println("server: longResponseChecksContext ended")

	ctx := req.Context()
	select {
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case <-time.After(5 * time.Second):
		fmt.Fprintf(w, "slept 5 seconds\n")
	}
}

func longResponseDB(w http.ResponseWriter, req *http.Request) {
	fmt.Println("server: longResponseDB started")
	defer fmt.Println("server: longResponseDB ended")

	ctx := req.Context()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.ExecContext(ctx, "DO SLEEP(5)")
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	randomName := getRandomName()
	_, err = tx.ExecContext(ctx, "INSERT INTO `students` (name) VALUES (?)",
		randomName)
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "slept 5 seconds and inserted record %s\n", randomName)
}

func longResponseDBNoTx(w http.ResponseWriter, req *http.Request) {
	fmt.Println("server: longResponseDBNoTx started")
	defer fmt.Println("server: longResponseDBNoTx ended")

	ctx := req.Context()

	_, err := db.ExecContext(ctx, "DO SLEEP(5)")
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	randomName := getRandomName()
	_, err = db.ExecContext(ctx, "INSERT INTO `students` (name) VALUES (?)",
		randomName)
	if err != nil {
		fmt.Println("server:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "slept 5 seconds and inserted record %s\n", randomName)
}

func main() {
	fmt.Println("starting server")
	defer db.Close()

	http.HandleFunc("/longResponse", longResponse)
	http.HandleFunc("/longResponseChecksContext", longResponseChecksContext)
	http.HandleFunc("/longResponseDB", longResponseDB)
	http.HandleFunc("/longResponseDBNoTx", longResponseDBNoTx)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
