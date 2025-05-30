package main

import (
	"fmt"
	"goTechReady/initializer"
	"goTechReady/pubsub"
)

func main() {
	initializer.LoadConfig()
	initializer.ConnectCloudSQL() // Only call this once to initialize the global DB
	fmt.Println("db:", initializer.GetDB())
	pubsub.ListenToOrders() // Start listening to Pub/Sub messages
}
