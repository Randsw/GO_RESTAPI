package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/randsw/ha-postgres/records"
	log "github.com/sirupsen/logrus"
)

func initHandlers(pool *pgxpool.Pool) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/records",
		func(w http.ResponseWriter, r *http.Request) {
			records.SelectAll(pool, w, r)
		}).Methods("GET")

	r.HandleFunc("/api/v1/records/{id:[0-9]+}",
		func(w http.ResponseWriter, r *http.Request) {
			records.Select(pool, w, r)
		}).Methods("GET")

	r.HandleFunc("/api/v1/records",
		func(w http.ResponseWriter, r *http.Request) {
			records.Insert(pool, w, r)
		}).Methods("POST")

	r.HandleFunc("/api/v1/records/{id:[0-9]+}",
		func(w http.ResponseWriter, r *http.Request) {
			records.Update(pool, w, r)
		}).Methods("PUT")

	r.HandleFunc("/api/v1/records/{id:[0-9]+}",
		func(w http.ResponseWriter, r *http.Request) {
			records.Delete(pool, w, r)
		}).Methods("DELETE")
	return r
}

func CreateAndFillTable (p *pgxpool.Pool) error{
	_, err := p.Exec(context.Background(),
		"CREATE TABLE peoples(id SERIAL PRIMARY KEY, name VARCHAR(64), surname VARCHAR(64), gender VARCHAR(64), email VARCHAR(64));")
	if err != nil {
			log.Errorf("Cannot create table: %v", err)
			return err
	}
	// Create variable with data
	TestPeople := Init()
	var id uint64
	for _, person := range TestPeople{
		row := p.QueryRow(context.Background(),
			"INSERT INTO peoples (name, surname, gender, email) VALUES ($1, $2, $3, $4) RETURNING id",
			person.Name, person.Surname, person.Gender, person.Email)
		err = row.Scan(&id)
		if err != nil {
			log.Errorf("Unable to INSERT: %v", err)
			return err
		}
	}
	return nil
}

func main() {
	type DBparams struct {
		host     string
		port     int
		user     string
		password string
		dbname   string
	}
	var id uint64
	var Name, Surname, Gender, Email string

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.InfoLevel)
	//log.SetOutput()

	//db := DBparams{"localhost", 5432, "postgres", "password", "peoples"}
	db := DBparams{"192.168.0.214", 5000, "admin", "admin", "postgres"}
	//db.password = os.Getenv("PG_PASSWORD")
	//db.host = os.Getenv("PG_HOST")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.host, db.port, db.user, db.password, db.dbname)
	var conn *pgxpool.Pool
	var err error
	//conn, err := pgx.Connect(context.Background(), psqlconn)
	operation := func() error {
		conn, err = pgxpool.Connect(context.Background(), psqlconn)
		if err != nil {
			log.Errorf("Unable to connect to database: %v\n", err)
			//fmt.Fprintf(os.Stdout, "Unable to connect to database: %v\n", err)
			return err
		}
		return nil
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	log.Infof("Connect to postgres base: %s:%d  succesfull\n", db.host, db.port)
	defer conn.Close()

	row := conn.QueryRow(context.Background(), "SELECT * FROM peoples WHERE id=1")
	err = row.Scan(&id, &Name, &Surname, &Gender, &Email)
	if err != nil {
		log.Warnf("Table is empty. Add test data to table")
		err := CreateAndFillTable(conn)
		{
			if err != nil {
				log.Errorf("Unable to INSERT test data: %v", err)
				os.Exit(1)
			}
		}
	}
	listenAddr := ":8080"
	http.Handle("/", initHandlers(conn))
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

}
