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

type CotationResponse struct {
	Usdbrl struct {
		Bid    string `json:"bid"`
		Codein string `json:"codein"`
	} `json:"USDBRL"`
}

type Cotation struct {
	Bid  string `json:"bid"`
	Code string `json:"-"`
}

func NewCotation(cotation, code string) *Cotation {
	return &Cotation{Bid: cotation, Code: code}
}

func main() {
	http.HandleFunc("/cotacao", CotationHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CotationHandler(w http.ResponseWriter, r *http.Request) {
	cotation, err := GetCotation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if err := Insert(cotation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotation)
}

func GetCotation() (*Cotation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotationResponse CotationResponse
	if err := json.Unmarshal(body, &cotationResponse); err != nil {
		return nil, err
	}

	return NewCotation(cotationResponse.Usdbrl.Bid, cotationResponse.Usdbrl.Codein), nil
}

func Insert(cotation *Cotation) error {
	db, err := sql.Open("sqlite3", "./desafio.db")
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	stmt, err := db.Prepare("insert into cotacoes(code, cotation) values($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotation.Code, cotation.Bid)
	if err != nil {
		return err
	}
	return nil
}
