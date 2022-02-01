package records

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Record struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
}

func SelectAll(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var str string
	var rs = make([]Record, 0)
	var rec Record
	var rows pgx.Rows
	if len(r.URL.RawQuery) > 0 {
		str = r.URL.Query().Get("name")
		if str == "" {
			w.WriteHeader(400)
			return
		}
	}
	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	if str != "" {
		rows, err = conn.Query(context.Background(),
			"SELECT id, name, surname, gender, email FROM peoples WHERE name LIKE $1 ORDER BY id", "%"+str+"%")

	} else {
		rows, err = conn.Query(context.Background(), "SELECT * FROM peoples ORDER BY id")
	}
	if err != nil {
		log.Errorf("Unable to SELECT ALL: %v\n", err)
		w.WriteHeader(500)
		return
	}

	for rows.Next() {
	err = rows.Scan(&rec.Id, &rec.Name, &rec.Surname, &rec.Gender, &rec.Email)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		w.WriteHeader(500)
		return
	}
	rs = append(rs, rec)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(rs); err != nil {
		w.WriteHeader(500)
	}
}

func Select(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"SELECT id, name, surname, gender, email FROM peoples WHERE id = $1",
		id)

	var rec Record
	err = row.Scan(&rec.Id, &rec.Name, &rec.Surname, &rec.Gender, &rec.Email)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		return
	}

	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(rec)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func Insert(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var rec Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"INSERT INTO peoples (name, surname, gender, email) VALUES ($1, $2, $3, $4) RETURNING id",
		rec.Name, rec.Surname, rec.Gender, rec.Email)
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		w.WriteHeader(500)
		return
	}

	resp := make(map[string]string, 1)
	resp["id"] = strconv.FormatUint(id, 10)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		w.WriteHeader(500)
		return
	}
}

func Update(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	var rec Record
	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		"UPDATE peoples SET name = $2, surname = $3, gender = $4, email = $5 WHERE id = $1",
		id, rec.Name, rec.Surname, rec.Gender, rec.Email)
	if err != nil {
		log.Errorf("Unable to UPDATE: %v\n", err)
		w.WriteHeader(500)
		return
	}

	if ct.RowsAffected() == 0 {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
}

func Delete(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil { // bad request
		w.WriteHeader(400)
		return
	}

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), "DELETE FROM peoples WHERE id = $1", id)
	if err != nil {
		log.Errorf("Unable to DELETE: %v", err)
		w.WriteHeader(500)
		return
	}

	if ct.RowsAffected() == 0 {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
}
