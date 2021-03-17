package pubsub

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/loadimpact/k6/js/modules"
	"github.com/loadimpact/k6/lib"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/pubsub".
func init() {
	modules.Register("k6/x/pubsub", new(PubSub))
}

// PubSub is the k6 extension for a Google Pub/Sub client.
// See https://cloud.google.com/pubsub/docs/overview
type PubSub struct{}

// Publisher is the basic wrapper for Google Pub/Sub publisher and uses
// watermill as a client. See https://github.com/ThreeDotsLabs/watermill/
//
// Publisher represents the constructor and creates an instance of
// googlecloud.Publisher with provided projectID and publishTimeout.
// Publisher uses watermill StdLoggerAdapter logger.
func (ps *PubSub) Publisher(projectID string, publishTimeout int, debug, trace bool) *googlecloud.Publisher {
	if publishTimeout < 1 {
		publishTimeout = 5
	}

	client, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: projectID,
			Marshaler: googlecloud.DefaultMarshalerUnmarshaler{},
			PublishTimeout: time.Second * time.Duration(publishTimeout),
		},
		watermill.NewStdLogger(debug, trace),
	)

	if err != nil {
		log.Fatalf("xk6-pubsub: unable to init publisher: %v", err)
	}

	return client
}

// Publish publishes a message to the provided topic using provided
// googlecloud.Publisher. The msg value must be passed as string
// and will be converted to bytes sequence before publishing.
func (ps *PubSub) Publish(ctx context.Context, p *googlecloud.Publisher, topic, msg string) error {
	state := lib.GetState(ctx)

	if state == nil {
		err := errors.New("xk6-pubsub: state is nil")
		ReportError(err, "cannot determine state")
		return err
	}

	err := p.Publish(
		topic,
		message.NewMessage(watermill.NewShortUUID(), []byte(msg)),
	)

	if err != nil {
		ReportError(err, "xk6-pubsub: unable to publish message")
		return err
	}

	return nil
}
