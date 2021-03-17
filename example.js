import { check } from 'k6';
import { publisher, publish } from 'k6/x/pubsub';

export default function () {
    // Creates a new publisher for ProjectID with a timeout of 2 seconds
    // with debug and trace mod enabled
    const client = publisher('', 2, true, true)
    let error = publish(client, 'topic_name', 'message data');

    check(error, {
        "is sent": err => err === null
    });

    client.close()
}
