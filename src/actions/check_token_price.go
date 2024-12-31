package actions

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"token_price/src/config"
	"token_price/src/utils"
)

type TokenPrice struct {
	Symbol    string    `bson:"symbol"`
	Currency  string    `bson:"currency"`
	Price     float64   `bson:"price"`
	Timestamp time.Time `bson:"timestamp"`
}

var tokenTargetMap = make(map[string]float64)

func CheckPrice(token string) {

	alchemy := utils.GetValue("ALCHEMY_KEY", "")
	url := alchemy + token
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	price := insert_token_price(body, token)

	checkValue(price, token)

}

/**
 * insert data
 */
func insert_token_price(body []byte, token string) TokenPrice {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// 提取 "data" 部分（只取第一个元素）
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		log.Fatalf("Invalid or missing 'data' field")
	}

	// 获取 data[0]
	firstData, ok := data[0].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid format for first data element")
	}

	// 提取 symbol
	symbol, ok := firstData["symbol"].(string)
	if !ok {
		log.Fatalf("Missing or invalid 'symbol'")
	}

	// 提取 prices（假设只有一个 price）
	prices, ok := firstData["prices"].([]interface{})
	if !ok || len(prices) == 0 {
		log.Fatalf("Missing or invalid 'prices'")
	}

	priceData, ok := prices[0].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid format for price element")
	}

	// 提取 currency, value, lastUpdatedAt
	currency, _ := priceData["currency"].(string)
	value, _ := priceData["value"].(string)
	lastUpdatedAt, _ := priceData["lastUpdatedAt"].(string)

	parsedTime, err := time.Parse(time.RFC3339, lastUpdatedAt)
	if err != nil {
		log.Fatalf("Failed to parse LastUpdatedAt: %v", err)
	}

	parsedPrice, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Failed to parse price: %v", err)
	}

	price := TokenPrice{
		Symbol:    symbol,
		Currency:  currency,
		Price:     parsedPrice,
		Timestamp: parsedTime,
	}

	client := utils.MongoClient
	db_name := utils.GetValue("DB", "")
	col_name := utils.GetValue("COLLECTION", "")
	collection := client.Database(db_name).Collection(col_name)
	_, err = collection.InsertOne(context.TODO(), price)
	if err != nil {
		log.Fatalf("Failed to insert price: %v", err)
	}

	log.Printf("Success insert [%s] price to MongoDB!", token)

	return price
}

func checkValue(price TokenPrice, token string) {
	if value, exists := tokenTargetMap[token]; exists {
		if value <= price.Price {
			//发邮件
			Send(price.Price, token)

			tokenTargetMap[token] = price.Price * config.TIMES
		}
	} else {
		tokenTargetMap[token] = price.Price * config.TIMES
	}

}
