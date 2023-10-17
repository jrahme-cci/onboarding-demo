package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"net"
	"net/http"
	"os"
)

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /greeting request\n")
	recordRequest(r.RemoteAddr, r.UserAgent())
	io.WriteString(w, "Greetings World!\n")
}

func recordRequest(ip, userAgent string) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	pass := os.Getenv("DB_PASS")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Printf("oops coulnd't connect to db %s\n", err)
		return
	}
	defer db.Close()

	insert := `INSERT INTO "tracker"("ip","useragent") values($1, $2)`
	_, err = db.Exec(insert, ip, userAgent)
	if err != nil {
		fmt.Printf("oops couldn't insert to db: %s\n", err)
	}
}

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", greetingHandler)
	server := &http.Server{
		Addr:    ":1234",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, "0.0.0.0", l.Addr().String())
			return ctx
		},
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", err)
	}
}
