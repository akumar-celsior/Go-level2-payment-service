package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"goTechReady/initializer"
	"goTechReady/model"
	"log"

	"cloud.google.com/go/pubsub"
)

func ListenToOrders() {
	db := initializer.GetDB()
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "linear-outcome-456809-t1")
	if err != nil {
		log.Fatal(err)
	}

	sub := client.Subscription("order-events-sub")
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Println("Invalid message data")
			msg.Nack()
			return
		}

		// Simulate payment processing
		payment := model.Payment{
			OrderID: order.ID,
			Amount:  order.Amount,
			Status:  "PAID",
		}
		db.Create(&payment)

		fmt.Printf("Payment processed for Order ID %s\n", order.ID)
		msg.Ack()
	})

	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}
}
