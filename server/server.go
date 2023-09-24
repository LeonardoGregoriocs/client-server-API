package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Usdbrl struct {
	Usdbrl Quotation `json:"USDBRL"`
}

type Quotation struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/quotation", quotationHandler)
	http.ListenAndServe(":8080", nil)
}

func quotationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	connection, err := ConnectionDB()
	if err != nil {
		log.Fatal("Error connection DB: ", err)
	}

	quotation, err := getQuotation(ctx)
	if err != nil {
		log.Println("Error during getQuotation: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	err = Insert(connection, quotation)
	if err != nil {
		log.Println("Error insert Data in database: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation.Usdbrl.Bid)
}

func getQuotation(ctx context.Context) (*Usdbrl, error) {
	URL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Println("Error during create request")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error get request")
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error during read body")
		return nil, err
	}

	var quotation Usdbrl
	err = json.Unmarshal(body, &quotation)
	if err != nil {
		log.Println("Error Unmarshal")
		return nil, err
	}

	return &quotation, nil
}

func ConnectionDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../quotation.db")
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS quotation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code       TEXT,
			codein     TEXT,
			name       TEXT,
			high       TEXT,
			low        TEXT,
			varBid     TEXT,
			pctChange  TEXT,
			bid        TEXT,
			ask        TEXT,
			timestamp  TEXT,
			createDate TEXT
		);
	`)

	if err != nil {
		return err
	}
	return nil
}

func Insert(db *sql.DB, quotation *Usdbrl) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()
	stmt, err := db.Prepare("INSERT INTO quotation (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.ExecContext(ctx,
		quotation.Usdbrl.Code,
		quotation.Usdbrl.Codein,
		quotation.Usdbrl.Name,
		quotation.Usdbrl.High,
		quotation.Usdbrl.Low,
		quotation.Usdbrl.VarBid,
		quotation.Usdbrl.PctChange,
		quotation.Usdbrl.Bid,
		quotation.Usdbrl.Ask,
		quotation.Usdbrl.Timestamp,
		quotation.Usdbrl.CreateDate,
	)

	if err != nil {
		return err
	}

	log.Println("Data inserted in database")
	return nil
}
