package pubsub

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules"
	"github.com/loadimpact/k6/lib"
)

func init() {
	modules.Register("k6/x/pubsub", new(PubSub))
}

// PubSub is the k6 extension for a Google Pub/Sub client.
// See https://cloud.google.com/pubsub/docs/overview
type PubSub struct{}

// Publisher is the basic wrapper for Google Pub/Sub publisher. It uses watermill
// as a client. See https://github.com/ThreeDotsLabs/watermill/
type Publisher struct {
	client *googlecloud.Publisher
}

// XPublisher represents the Publisher constructor and creates an instance
// of googlecloud.Publisher with provided settings.
func (pb *PubSub) XPublisher(ctx *context.Context, projectID string, publishTimeout int) interface{} {
	client, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: projectID,
			Marshaler: googlecloud.DefaultMarshalerUnmarshaler{},
			PublishTimeout: time.Second * time.Duration(publishTimeout),
		},
		watermill.NewStdLogger(true, true),
	)

	if err != nil {
		log.Fatalf("xk6-pubsub: unable to init extension: %v", err)
	}

	rt := common.GetRuntime(*ctx)
	return common.Bind(rt, &Publisher{client: client}, ctx)
}

// Publish publish a message to the provided topic.
func (p *Publisher) Publish(ctx context.Context, topic, msg string) error {
	state := lib.GetState(ctx)

	if state == nil {
		err := errors.New("xk6-pubsub: state is nil")
		ReportError(err, "cannot determine state")
		return err
	}

	err := p.client.Publish(
		topic,
		message.NewMessage(watermill.NewShortUUID(), []byte(msg)),
	)

	if err != nil {
		ReportError(err, "xk6-pubsub: unable to publish message")
		return err
	}

	return nil
}
