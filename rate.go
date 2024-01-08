package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CoinGeckoResponse struct {
	MarketData struct {
		CurrentPrice struct {
			UAH float64 `json:"uah"`
		} `json:"current_price"`
	} `json:"market_data"`
}

// HandleRateRequest оброблює запит /rate,
// витягує актуальний курс бікоїну з відкритого джерела, та повертає його у вигляді string
func HandleRateRequest(w http.ResponseWriter) string {
	// Отримання актуального курсу
	price, err := http.Get("https://api.coingecko.com/api/v3/coins/bitcoin")
	if err != nil {
		return ""
	}
	defer price.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(price.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Серіалізація структури у формат JSON
	jsonData, _ := json.Marshal(data.MarketData.CurrentPrice.UAH)

	// Конвертація рядка JSON в рядок string
	jsonStr := string(jsonData)

	// Виводимо результат на localhost
	fmt.Fprintf(w, "The current Bitcoin price in UAH is $%s", jsonStr)

	return jsonStr
}
