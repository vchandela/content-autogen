package job

import (
	"content_autogen/drivers"
	"content_autogen/parser"
	"content_autogen/receiver"
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"gopkg.in/matryer/try.v1"
)

type HamsaInsertionJob struct {
	messageChan chan *pubsub.Message
	quitChan    chan bool
	receiver    receiver.Receiver
	parser      parser.Parser
}

const (
	BufferedFeaturesChannelLen = 10000
)

func NewHamsaInsertionJob(ctx context.Context, receiver receiver.Receiver,
	parser parser.Parser) *HamsaInsertionJob {

	job := &HamsaInsertionJob{
		messageChan: make(chan *pubsub.Message, BufferedFeaturesChannelLen),
		quitChan:    make(chan bool),
		receiver:    receiver,
		parser:      parser,
	}

	go func() {
		maxRetries := 100
		err := try.Do(func(attempt int) (bool, error) {
			err := receiver.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				job.messageChan <- msg
			})

			if err != nil {
				drivers.Logger.Errorf("error in pulling messages from pubusb, retrying, err=%v", err)
				time.Sleep(2 * time.Second)
			}

			return attempt < maxRetries, err
		})

		if err != nil {
			drivers.Logger.Errorf("error in pulling from pubsub, err=%v", err)
			job.Cancel()
		}
	}()

	return job
}

func (job *HamsaInsertionJob) Cancel() {
	job.quitChan <- true
}

func (job *HamsaInsertionJob) Run(ctx context.Context) {
	for {
		select {
		case msg := <-job.messageChan:
			{
				err := job.processMessage(ctx, msg)
				if err != nil {
					drivers.Logger.Errorf("error while processing message, err=%v", err)
					msg.Nack()
				} else {
					msg.Ack()
				}
			}
		case <-job.quitChan:
			{
				drivers.Logger.Errorf("cancelling job")
				return
			}
		}
	}
}

func (job *HamsaInsertionJob) processMessage(ctx context.Context, msg *pubsub.Message) error {
	eventData, err := job.parser.Parse(msg.Data)
	if err != nil {
		return nil
	}

	fmt.Println((eventData.Maintext))

	return nil
}
