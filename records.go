package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strings"
)

type Record struct {
	Id    int
	Name  string
	Phone string
}

func Insert(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Некорректные данные"))
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	if name == "" || phone == "" {
		w.WriteHeader(500)
		w.Write([]byte("Имя и телефон не могут быть пустыми"))
		return
	}
	record := Record{
		Name:  name,
		Phone: phone,
	}
	insertSQL := `INSERT INTO contact_book (name, phone) VALUES ($1, $2);`
	_, err := db.Exec(context.Background(), insertSQL, record.Name, record.Phone)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка при добавлении записи в бд"))
		return
	}
	w.WriteHeader(200)
}

func SelectAll(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(),
		"SELECT name, phone FROM contact_book")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка при чтении данных"))
		return
	}
	defer rows.Close()
	var contacts []string
	for rows.Next() {
		var rec Record
		err := rows.Scan(&rec.Name, &rec.Phone)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Ошибка при чтении данных"))
			return
		}
		contacts = append(contacts, fmt.Sprintf("Имя: %-8s Телефон: %s", rec.Name, rec.Phone))
	}

	if len(contacts) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Список контактов пуст"))
		return
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка при обработке результатов"))
		return
	}

	responseString := strings.Join(contacts, "\n")
	w.Write([]byte(responseString))

}

func Select(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(400)
		w.Write([]byte("Параметр 'name' отсутствует"))
		return
	}
	row := db.QueryRow(context.Background(),
		"SELECT name, phone FROM contact_book WHERE name = $1;",
		name)
	var rec Record
	err := row.Scan(&rec.Name, &rec.Phone)
	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(404)
		w.Write([]byte("Имя не найдено"))
		return
	}
	responseString := fmt.Sprintf("Имя: %s Телефон: %s", rec.Name, rec.Phone)
	w.Write([]byte(responseString))
}

func Update(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Некорректные данные"))
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	if name == "" || phone == "" {
		w.WriteHeader(500)
		w.Write([]byte("Имя и телефон не могут быть пустыми"))
		return
	}
	record := Record{
		Name:  name,
		Phone: phone,
	}
	updateSQL := `UPDATE contact_book SET name = $1, phone = $2	WHERE name = $1;`
	commandTag, err := db.Exec(context.Background(), updateSQL, record.Name, record.Phone)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка при изменении записи в бд"))
		return
	}
	if commandTag.RowsAffected() == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Запись не найдена"))
		return
	}
	w.WriteHeader(200)
}

func Delete(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(400)
		w.Write([]byte("Параметр 'name' отсутствует"))
		return
	}
	deleteSQL := `DELETE FROM contact_book WHERE name = $1;`
	commandTag, err := db.Exec(context.Background(), deleteSQL, name)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка при изменении записи в бд"))
		return
	}
	if commandTag.RowsAffected() == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Запись не найдена"))
		return
	}
	w.WriteHeader(200)
}
