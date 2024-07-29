package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "stock-task", worker.Options{})
	w.RegisterWorkflow(stockFetcherWorkflow)
	w.RegisterActivity(FetchStockDataActivity)

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
