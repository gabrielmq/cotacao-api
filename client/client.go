package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Cotation struct {
	Bid string `json:"bid"`
}

func main() {
	cotation, err := GetCotation()
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./cotacao.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Arquivo criado com sucesso!")
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotation.Bid))
	if err != nil {
		panic(err)
	}
}

func GetCotation() (*Cotation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
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

	var cotation Cotation
	if err := json.Unmarshal(body, &cotation); err != nil {
		return nil, err
	}
	return &cotation, nil
}
