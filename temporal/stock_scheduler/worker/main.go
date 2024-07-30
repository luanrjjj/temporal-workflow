package main

import (
	"log"

	stock_scheduler "github.com/luanrjjj/temporal-workflow/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	// temporal/stock_scheduler
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "schedule", worker.Options{})

	w.RegisterWorkflow(stock_scheduler.StockFetcherWorkflow)
	w.RegisterActivity(stock_scheduler.FetchStockDataActivity)
	w.RegisterActivity(stock_scheduler.SaveStockDataActivityOnDatabase)

	err = w.Run(worker.InterruptCh())

	// we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, stockFetcherWorkflow)
	// if err != nil {
	// 	log.Fatalln("Unable to execute workflow", err)
	// }
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
	// log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
