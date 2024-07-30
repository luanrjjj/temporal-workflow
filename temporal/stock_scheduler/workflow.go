package stock_scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RealItem struct {
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

// Define the struct to represent the entire JSON response
type ApiResponse struct {
	Real []RealItem `json:"real"`
}

type NewJsonFormat struct {
	Code string      `json:"code"`
	Data ApiResponse `json:"data`
}

var stockSymbols = []string{
	"itub4",
	"petr3",
	"vale3",
	"petr4",
	"b3sa3",
	// Add more symbols as needed
}

func transformJson(code string, jsonData string) (string, error) {
	// Create a variable to hold the unmarshalled data
	var apiResponse ApiResponse

	// Unmarshal the JSON data into the ApiResponse struct
	err := json.Unmarshal([]byte(jsonData), &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Create a new struct to hold the transformed data
	newFormat := NewJsonFormat{
		Code: code,
		Data: apiResponse,
	}

	// Marshal the new struct into JSON
	newJsonData, err := json.MarshalIndent(newFormat, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	return string(newJsonData), nil
}

func FetchStockData(symbol string) (string, error) {
	fmt.Println("kakka", symbol)

	url := fmt.Sprintf("https://investidor10.com.br/api/cotacoes/acao/chart/%s/1825/true/real", symbol)
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var data ApiResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		return "", err
	}
	// funs := make(map[string]string)
	// funs[symbol] = string(body)

	newJson, err := transformJson(symbol, string(body))
	return newJson, nil
}

func FetchStockDataActivity(ctx context.Context, symbol string) (string, error) {
	activity.GetLogger(ctx).Info("Fetching stock data for symbol", "symbol", symbol)
	activity.RecordHeartbeat(ctx)
	return FetchStockData(symbol)
}

func SaveStockDataActivityOnDatabase(ctx context.Context, records string) (string, error) {
	fmt.Printf(records)
	return "string", nil
}

const fetchStockDataActivityName = "fetch-stock-data-activity"

func StockFetcherWorkflow(ctx workflow.Context) (string, error) {
	// Fetch stock data for each symbol
	var result string
	// var stockData string

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	info := workflow.GetInfo(ctx)

	scheduleByIdPayload := info.SearchAttributes.IndexedFields["TemporalScheduledById"]
	var scheduledById string

	err := converter.GetDefaultDataConverter().FromPayload(scheduleByIdPayload, &scheduledById)
	if err != nil {
		// return err
	}
	startTimePayload := info.SearchAttributes.IndexedFields["TemporalScheduledStartTime"]
	var startTime time.Time
	err = converter.GetDefaultDataConverter().FromPayload(startTimePayload, &startTime)
	if err != nil {
		// return err
	}

	for _, symbol := range stockSymbols {
		err = workflow.ExecuteActivity(ctx, FetchStockDataActivity, symbol, scheduledById, startTime).Get(ctx, &result)

		if err != nil {
			return result, nil
		}
		// var verify string
		// err = workflow.ExecuteActivity(ctx, SaveStockDataActivityOnDatabase, scheduledById, startTime, symbol).Get(ctx, &verify)

		// if err != nil {
		// 	return verify, nil
		// }

		workflow.GetLogger(ctx).Info("Result for symbol", "symbol", symbol, "result", result)
	}
	return result, nil
	return "finished", nil
	// final := fmt.Sprintf("Finish workflow %s, %s", result, "oi")
}
