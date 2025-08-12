package engine

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func InsertStocks(stockResponse StockResponse) error {

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	fmt.Println("Iniciando inserción de stocks")
	for _, stock := range stockResponse.Items {
		code := uuid.New().String()
		targetFromCleaned := strings.ReplaceAll(stock.TargetFrom, "$", "")
		targetFromCleaned = strings.ReplaceAll(targetFromCleaned, ",", "")
		targetFrom, err1 := strconv.ParseFloat(targetFromCleaned, 64)
		if err1 != nil {
			log.Fatalf("error parsing target from: %v", err1)
		}
		targetToCleaned := strings.ReplaceAll(stock.TargetTo, "$", "")
		targetToCleaned = strings.ReplaceAll(targetToCleaned, ",", "")
		targetTo, err2 := strconv.ParseFloat(targetToCleaned, 64)
		if err2 != nil {
			log.Fatalf("error parsing target to: %v", err2)
		}
		createdAt := time.Now()
		updatedAt := time.Now()
		_, err := db.Exec(ctx, "INSERT INTO stocks (code, ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, record_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", code, stock.Ticker, targetFrom, targetTo, stock.Company, stock.Action, stock.Brokerage, stock.RatingFrom, stock.RatingTo, stock.Time, createdAt, updatedAt)
		if err != nil {
			log.Fatalf("error inserting stock: %v", err)
		}
	}

	defer db.Close(ctx)

	return nil
}
