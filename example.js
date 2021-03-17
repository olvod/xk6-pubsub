import { check } from 'k6';
import pubsub from 'k6/x/pubsub';

const publisher = new pubsub.Publisher('ProjectID', 2);

export default function () {
    let error = publisher.publish('topic_name', 'message data');

    check(error, {
        "is sent": err => err === undefined
    });
}
