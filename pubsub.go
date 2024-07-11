package pubsub

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/pubsub".
func init() {
	modules.Register("k6/x/pubsub", new(PubSub))
}

// PubSub is the k6 extension for a Google Pub/Sub client.
// See https://cloud.google.com/pubsub/docs/overview
type PubSub struct{}

// publisherConf provides a Pub/Sub publisher client configuration. This configuration
// structure can be used on a client side. All parameters are optional.
type publisherConf struct {
	ProjectID                 string
	Credentials               string
	PublishTimeout            int
	Debug                     bool
	Trace                     bool
	DoNotCreateTopicIfMissing bool
}

// Publisher is the basic wrapper for Google Pub/Sub publisher and uses
// watermill as a client. See https://github.com/ThreeDotsLabs/watermill/
//
// Publisher represents the constructor and creates an instance of
// googlecloud.Publisher with provided projectID and publishTimeout.
// Publisher uses watermill StdLoggerAdapter logger.
func (ps *PubSub) Publisher(config map[string]interface{}) *googlecloud.Publisher {
	cnf := &publisherConf{}
	err := mapstructure.Decode(config, cnf)
	if err != nil {
		log.Fatalf("xk6-pubsub: unable to read publisher config: %v", err)
	}

	if cnf.PublishTimeout < 1 {
		cnf.PublishTimeout = 5
	}

	client, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID:                 cnf.ProjectID,
			Marshaler:                 googlecloud.DefaultMarshalerUnmarshaler{},
			PublishTimeout:            time.Second * time.Duration(cnf.PublishTimeout),
			DoNotCreateTopicIfMissing: cnf.DoNotCreateTopicIfMissing,
			ClientOptions:             withCredentials(cnf.Credentials),
		},
		watermill.NewStdLogger(cnf.Debug, cnf.Trace),
	)

	if err != nil {
		log.Fatalf("xk6-pubsub: unable to init publisher: %v", err)
	}

	return client
}

// Publish publishes a message to the provided topic using provided
// googlecloud.Publisher. The msg value must be passed as string
// and will be converted to bytes sequence before publishing.
func (ps *PubSub) Publish(ctx context.Context, p *googlecloud.Publisher, topic, msg string, attributes map[string]string) error {
	state := lib.GetState(ctx)

	if state == nil {
		err := errors.New("xk6-pubsub: state is nil")
		ReportError(err, "cannot determine state")
		return err
	}

	m := message.NewMessage(watermill.NewShortUUID(), []byte(msg))
	m.Metadata = attributes
	err := p.Publish(
		topic,
		m,
	)

	if err != nil {
		ReportError(err, "xk6-pubsub: unable to publish message")
		return err
	}

	return nil
}

// withCredentials explicitly setup Pub/Sub credentials as option.ClientOption.
func withCredentials(credentials string) []option.ClientOption {
	var opt []option.ClientOption

	if len(credentials) > 0 {
		opt = append(opt, option.WithCredentialsJSON([]byte(credentials)))
	}

	return opt
}
