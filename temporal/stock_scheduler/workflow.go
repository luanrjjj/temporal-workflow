package stock_scheduler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type AmazingPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// func main() {
// 	c, err := client.Dial(client.Options{
// 		HostPort: client.DefaultHostPort,
// 	})
// 	if err != nil {
// 		log.Fatalln("Unable to create client", err)
// 	}
// 	defer c.Close()

// 	// This workflow ID can be user business logic identifier as well.
// 	workflowID := "cron_" + uuid.New()
// 	workflowOptions := client.StartWorkflowOptions{
// 		ID:           workflowID,
// 		TaskQueue:    "cron",
// 		CronSchedule: "* * * * *",
// 	}
// 	myPayload := AmazingPayload{Name: "aezo", Age: 12}
// 	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "say-hello-workflow", myPayload)
// 	if err != nil {
// 		log.Fatalln("Unable to execute workflow", err)
// 	}
// 	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
// }

var stockSymbols = []string{
	"itub4",
	"petrr3",
	"vale3",
	"petr4",
	"b3sa3",
	// Add more symbols as needed
}

func FetchStockData(symbol string) (string, error) {
	fmt.Println("asudhuashduas", symbol)

	url := fmt.Sprintf("https://investidor10.com.br/api/cotacoes/acao/chart/%s/1825/true/real", symbol)
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// var data json
	// err = json.Unmarshal(body, &data)
	// fmt.Println("kkkkkkk", string(body))

	if err != nil {
		return "", err
	}
	// funs := make(map[string]string)
	// funs[symbol] = string(body)

	return symbol + " " + string(body), nil
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

	// scheduleHandle, err := c.ScheduleClient().Create(ctx, c.ScheduleOptions{
	// 	ID: scheduleID,
	// 	Spec: c.ScheduleSpec{},
	// 	Action: &c.ScheduleWorkflowAction{
	// 		ID: workflowID,
	// 		Workflow: stockFetcherWorkflow,
	// 		TaskQueue: "shcedule",
	// 	},
	// })

	for _, symbol := range stockSymbols {
		err = workflow.ExecuteActivity(ctx, FetchStockDataActivity, scheduledById, startTime, symbol).Get(ctx, &result)

		if err != nil {
			return result, nil
		}

		// err = workflow.ExecuteActivity(ctx, saveStockDataActivityOnDatabase, result).Get(ctx, &stockData)

		// if err != nil {
		// 	return result, err
		// }

		workflow.GetLogger(ctx).Info("Result for symbol", "symbol", symbol, "result", result)

		return result, nil
	}
	return "finished", nil
	// final := fmt.Sprintf("Finish workflow %s, %s", result, "oi")
}

// func main() {
// 	c, err := client.Dial(client.Options{
// 		HostPort: client.DefaultHostPort,
// 	})
// 	if err != nil {
// 		log.Fatalln("Unable to create client", err)
// 	}
// 	defer c.Close()

// 	w := worker.New(c, "stock-task", worker.Options{})
// 	w.RegisterWorkflow(StockFetcherWorkflow)
// 	w.RegisterActivity(FetchStockDataActivity)

// 	err = w.Run(worker.InterruptCh())

// 	// we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, stockFetcherWorkflow)
// 	// if err != nil {
// 	// 	log.Fatalln("Unable to execute workflow", err)
// 	// }
// 	if err != nil {
// 		log.Fatalln("Unable to start worker", err)
// 	}
// 	// log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
// }
