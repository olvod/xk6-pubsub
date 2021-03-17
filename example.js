import { check } from 'k6';
import pubsub from 'k6/x/pubsub';

// Creates a new publisher for ProjectID with a timeout of 2 seconds for the publisher
const publisher = new pubsub.Publisher('ProjectID', 2);

export default function () {
    let error = publisher.publish('topic_name', 'message data');

    check(error, {
        "is sent": err => err === undefined
    });
}
