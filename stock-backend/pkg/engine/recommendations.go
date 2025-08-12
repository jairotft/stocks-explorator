package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetDBRecommendations() (Recommendation, error) {

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	query := "SELECT code, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, record_time, created_at, updated_at FROM stocks WHERE target_to > target_from order by record_time desc limit 50"

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Printf("query error: %v", err)
		return Recommendation{}, err
	}
	defer rows.Close()

	var stocks []Stock
	for rows.Next() {
		var stock Stock
		err := rows.Scan(
			&stock.Code,
			&stock.Ticker,
			&stock.Company,
			&stock.Brokerage,
			&stock.Action,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.RecordTime,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		log.Printf("rows error: %v", err)
		return Recommendation{}, err
	}

	return Recommendation{
		Stocks: stocks,
	}, nil
}

func GetOpenAIRecommendations(stocks []Stock) (Recommendation, error) {

	var messages []OpenAIMessagePayload = make([]OpenAIMessagePayload, 0)
	messages = append(messages, OpenAIMessagePayload{
		Role: "system",
		Content: `Eres un asistente de inversiones en mercados bursatiles.
		Vas a recibir un listado de stocks preseleccionados y recomendarás tu top 3 de mejores inversiones tomando el cuenta las variables:
			* action
			* rating from
			* rating to
			* target from
			* target to
			* record_time
		Siempre empieza tu respuesta con "Estas son las acciones que creemos te pueden interesar:"
		Justifica cada una de tus recomendaciones.
		Finaliza tu respuesta indicando que es solo una recomendación considerando los datos que se tiene y no un consejo de inversión.
		`,
	})

	jsonStocks, err := json.Marshal(stocks)
	if err != nil {
		return Recommendation{}, err
	}

	messages = append(messages, OpenAIMessagePayload{
		Role:    `user`,
		Content: "Recomiendame los 3 mejores stocks para invertir en base en los datos que te voy a proporcionar: " + string(jsonStocks),
	})

	fmt.Println(messages)
	payload := OpenAIPayload{
		MaxTokens:        800,
		Temperature:      0.0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		TopP:             0.95,
		Stop:             nil,
		Messages:         messages,
	}

	recomedation := Recommendation{}

	llmResponse, err := CreateChat(payload)
	if err != nil {
		return recomedation, err
	}
	var message string
	if len(llmResponse.Choices) == 0 {
		message = "No se pudo obtener una respuesta"
	} else {
		message = llmResponse.Choices[0].Message.Content
	}

	recomedation.Stocks = stocks
	recomedation.Message = &message
	return recomedation, nil

}
