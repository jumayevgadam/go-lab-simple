package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Products is
type Products struct {
	ProductID   int     `json:"productID"`
	ProductName string  `json:"productName"`
	RetailPrice float64 `json:"retailPrice"`
}

// JsonResponse is
type JsonResponse struct {
	Data   []Products `json:"data"`
	Source string     `json:"source"`
}

func getProducts() (*JsonResponse, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	cachedProducts, err := redisClient.Get(ctx, "products").Bytes()
	response := JsonResponse{}

	if err != nil {
		dbProducts, err := fetchFromDB()
		if err != nil {
			return nil, err
		}

		cachedProducts, err = json.Marshal(dbProducts)
		if err != nil {
			return nil, err
		}

		err = redisClient.Set(ctx, "products", cachedProducts, 1*time.Minute).Err()
		if err != nil {
			return nil, err
		}

		response = JsonResponse{Data: dbProducts, Source: "PostgreSQL"}

		return &response, nil
	}

	products := []Products{}
	err = json.Unmarshal(cachedProducts, &products)
	if err != nil {
		return nil, err
	}

	response = JsonResponse{Data: products, Source: "RedisCache"}
	return &response, nil
}

func fetchFromDB() ([]Products, error) {
	dbUser := "postgres"
	dbPassword := "12345"
	dbName := "sample_company"

	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	queryString := `select
                      product_id,
                      product_name,
                      retail_price
                  from products`

	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}

	var records []Products
	for rows.Next() {
		var p Products

		err = rows.Scan(&p.ProductID, &p.ProductName, &p.RetailPrice)
		records = append(records, p)

		if err != nil {
			return nil, err
		}
	}

	return records, nil
}
