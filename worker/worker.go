package worker

import (
	"content_autogen/drivers"
	"content_autogen/job"
	"content_autogen/parser"
	"content_autogen/receiver"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Consume() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		osCall := <-c
		drivers.Logger.Infof("interrupt caught osCall=%v", osCall)
		drivers.SubCancel()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	eventParser := parser.NewParser()
	eventReceiver := receiver.NewPubSubReceiver(drivers.ContentAutogenPubsubSub)
	job := job.NewHamsaInsertionJob(ctx, eventReceiver, eventParser)

	maxWorkers := 1
	var wg sync.WaitGroup
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			job.Run(ctx)
		}(w)
	}

	wg.Wait()

	cancel()
}
