package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
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

	quotation, err := getQuotation(ctx)
	if err != nil {
		log.Println("Error during getQuotation: ", err)
		w.WriteHeader(http.StatusBadRequest)
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
