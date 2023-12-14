package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

func initHandler(db *pgx.Conn) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/contacts",
		func(w http.ResponseWriter, r *http.Request) {
			SelectAll(db, w, r)
		}).Methods("GET")

	r.HandleFunc("/contact",
		func(w http.ResponseWriter, r *http.Request) {
			Select(db, w, r)
		}).Methods("GET")

	r.HandleFunc("/contact",
		func(w http.ResponseWriter, r *http.Request) {
			Insert(db, w, r)
		}).Methods("POST")

	r.HandleFunc("/contact",
		func(w http.ResponseWriter, r *http.Request) {
			Delete(db, w, r)
		}).Methods("DELETE")

	r.HandleFunc("/contact",
		func(w http.ResponseWriter, r *http.Request) {
			Update(db, w, r)
		}).Methods("PUT")
	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	dbLogin := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	hostPort := strings.Split(os.Getenv("POSTGRES_PORT"), ":")[0]
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", dbLogin, dbPassword, hostPort, dbName)
	conn, connectErr := pgx.Connect(context.Background(), dbURL)
	if connectErr != nil {
		log.Fatalf("Ошибка соединения с базой: %v\n", connectErr)
	}
	log.Println("Соединение с базой прошло успешно")
	defer conn.Close(context.Background())

	createTableSQL := `CREATE TABLE IF NOT EXISTS contact_book (
		contact_id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		phone VARCHAR(10)
	);`
	_, createErr := conn.Exec(context.Background(), createTableSQL)
	if createErr != nil {
		log.Fatalf("Не удалось создать таблицу: %v\n", createErr)
	}
	log.Println("Таблица создана")

	http.Handle("/", initHandler(conn))
	contactHttpErr := http.ListenAndServe(":8080", nil)
	if contactHttpErr != nil {
		log.Fatalf("Ошибка запуска сервера: %v\n", contactHttpErr)
	}
}
