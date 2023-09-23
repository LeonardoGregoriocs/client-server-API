package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/quotation", nil)
	if err != nil {
		log.Println("Error during create the request: ", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error during request: ", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error during read body: ", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("Error during create file: ", err)
	}

	test := fmt.Sprint("Dólar: ", string(body))

	_, err = file.Write([]byte(test))
	if err != nil {
		log.Println("Error during write in file: ", err)
	}

}
