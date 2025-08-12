package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"stock/getter/pkg/engine"

	"github.com/joho/godotenv"
)

/*
*
* Get stock list from api

Query Params:
- next_page: next page token

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Response Item:

	{
		"ticker": "ETR",
		"target_from": "$85.00",
		"target_to": "$88.00",
		"company": "Entergy",
		"action": "target raised by",
		"brokerage": "KeyCorp",
		"rating_from": "Overweight",
		"rating_to": "Overweight",
		"time": "2025-07-17T00:30:07.155596923Z"
	}

* @param nextPage: next page token
* @return StockResponse
*/
func getStock(nextPage string) (engine.StockResponse, error) {

	url := os.Getenv("STOCKS_URL")
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	fmt.Println("Consultando a:", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+os.Getenv("API_TOKEN"))
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var stockResponse engine.StockResponse
	err = json.Unmarshal(body, &stockResponse)
	if err != nil {
		fmt.Println(string(body))
		log.Fatal(err)
	}

	return stockResponse, nil
}

func main() {

	godotenv.Overload()

	startTime := time.Now()
	totalItems := 0
	counter := 0
	var nextPage string = ""
	for {
		stockResponse, err := getStock(nextPage)
		counter++
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(stockResponse)
		nextPage = stockResponse.NextPage
		totalItems += len(stockResponse.Items)
		err = engine.InsertStocks(stockResponse)
		if err != nil {
			log.Fatal(err)
		}
		if nextPage == "" {
			break
		}
	}
	fmt.Println("Total de llamadas ejecutadas:", counter)
	fmt.Println("Total de items insertados:", totalItems)
	fmt.Println("Tiempo total de ejecución:", time.Since(startTime))
}
