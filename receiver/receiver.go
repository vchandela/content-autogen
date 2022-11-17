package receiver

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Receiver interface {
	Receive(ctx context.Context, f func(ctx context.Context, message *pubsub.Message)) error
}

type PubSubReceiver struct {
	subscription *pubsub.Subscription
}

func NewPubSubReceiver(subscription *pubsub.Subscription) Receiver {
	return &PubSubReceiver{subscription: subscription}
}

func (p *PubSubReceiver) Receive(ctx context.Context, processFunc func(ctx context.Context, message *pubsub.Message)) error {
	err := p.subscription.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		processFunc(ctx, message)
	})

	return err
}
