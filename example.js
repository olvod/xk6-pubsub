import { check } from 'k6';
import pubsub from 'k6/x/pubsub';

export default function () {
    /**
     * Creates a new publisher client that will be used by publisher.
     * publishTimeout value is 5 seconds by default.
     * debug and trace are disabled by default.
     */
    const client = pubsub.publisher({
        projectID: "",
        publishTimeout: 2,
        debug: true,
        trace: true,
        doNotCreateTopicIfMissing: true
    });
    let error = pubsub.publish(client, 'topic_name', 'message data');

    check(error, {
        "is sent": err => err === null
    });

    /**
     * The publisher client must be closed after a message was provided.
     * It will trigger the sending process.
     */
    client.close()
}
