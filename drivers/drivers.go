package drivers

import (
	"content_autogen/config"
	"context"
	"fmt"
	"runtime"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

var LogInstance, _ = zap.NewProduction()
var Logger = LogInstance.Sugar()
var ContentAutogenPubsubSub *pubsub.Subscription
var SubCancel context.CancelFunc

func InitializeDrivers() {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.ConfigurationMap.GCPProjectId)
	if err != nil {
		fmt.Printf("error in connecting to pubsub, err=%v\n", err)
	}

	numWorkers := scaledNumWorkers(1)

	ContentAutogenPubsubSub = pubsubClient.Subscription(config.ConfigurationMap.GCPSubscriptionId)
	ContentAutogenPubsubSub.ReceiveSettings.MaxOutstandingMessages = 10 * 1000
	ContentAutogenPubsubSub.ReceiveSettings.MaxOutstandingBytes = 100000000
	ContentAutogenPubsubSub.ReceiveSettings.NumGoroutines = numWorkers
}

func scaledNumWorkers(factor int) int {
	workers := factor * runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	return workers
}
