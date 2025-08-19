package engine

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"time"
)

func getOrderByClause(params map[string]string) string {
	orderClause := " ORDER BY created_at DESC"
	if params["order_by"] != "" {
		orderbyName := params["order_by"]
		// Solo si orderbyName esta en la lista de campos permitidos de "record_time", "created_at", "ticker", "company", "brokerage", "action", "rating_from", "rating_to", "target_from", "target_to"
		allowedFields := []string{"record_time", "created_at", "ticker", "company", "brokerage", "action", "rating_from", "rating_to", "target_from", "target_to"}
		if !slices.Contains(allowedFields, orderbyName) {
			return orderClause
		}
		ascOrDesc := params["asc"]
		if ascOrDesc == "1" {
			orderClause = " ORDER BY " + orderbyName + " ASC"
		} else {
			orderClause = " ORDER BY " + orderbyName + " DESC"
		}
	}
	return orderClause
}

func getWhereClause(params map[string]string) string {

	// Construir WHERE clause
	whereClause := " WHERE 1=1"

	if params["ticker"] != "" {
		whereClause += fmt.Sprintf(" AND ticker like '%%%s%%'", params["ticker"])
	}

	if params["brokerage"] != "" {
		whereClause += fmt.Sprintf(" AND brokerage like '%%%s%%'", params["brokerage"])
	}

	if params["action"] != "" {
		whereClause += fmt.Sprintf(" AND action like '%%%s%%'", params["action"])
	}

	if params["rating_from"] != "" {
		whereClause += fmt.Sprintf(" AND rating_from = '%s'", params["rating_from"])
	}

	if params["rating_to"] != "" {
		whereClause += fmt.Sprintf(" AND rating_to = '%s'", params["rating_to"])
	}

	return whereClause
}

func GetStocks(params map[string]string) (PaginatedStocksResponse, error) {

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	whereClause := getWhereClause(params)

	orderClause := getOrderByClause(params)

	// Paginación
	var page int
	perPage := 20
	if params["page"] == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(params["page"])
		if err != nil {
			log.Printf("page error: %v", err)
			return PaginatedStocksResponse{}, err
		}
	}

	// Consulta para obtener el total
	countQuery := "SELECT COUNT(*) FROM stocks" + whereClause
	var total int
	err = db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		log.Printf("count query error: %v", err)
		return PaginatedStocksResponse{}, err
	}

	// Consulta principal con paginación
	query := "SELECT code, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, record_time, created_at, updated_at FROM stocks"
	query += whereClause
	query += orderClause
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, (page-1)*perPage)

	fmt.Println(query)
	var stocks []Stock

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Printf("query error: %v", err)
		return PaginatedStocksResponse{}, err
	}
	defer rows.Close()

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
		return PaginatedStocksResponse{}, err
	}

	// Calcular página siguiente
	var nextPage *int
	if page*perPage < total {
		next := page + 1
		nextPage = &next
	}

	response := PaginatedStocksResponse{
		Stocks:      stocks,
		CurrentPage: page,
		NextPage:    nextPage,
		Total:       total,
		PerPage:     perPage,
	}

	return response, nil
}
