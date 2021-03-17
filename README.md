# xk6-pubsub[WIP:heavy_exclamation_mark:]

This is a [k6](https://github.com/loadimpact/k6) extension using the [xk6](https://github.com/k6io/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Install `xk6` framework for extending `k6`:
```shell
go get -u github.com/k6io/xk6/cmd/xk6
```

2. Build the binary:
```shell
xk6 build v0.31.0 --with github.com/olvod/xk6-pubsub
```

3. Setup Google Pub/Sub configuration via environment variables:
```shell
export PUBSUB_PROJECT_ID=<project_id>
export PUBSUB_CREDENTIALS=<credentials>
```

Or use [PubSub emulator](https://cloud.google.com/pubsub/docs/emulator#linux-macos) for local development 
`PUBSUB_EMULATOR_HOST` environment variable must be present.
```shell
export PUBSUB_EMULATOR_HOST=<emulator_host>
```

## Example

```javascript
import check from 'k6';

import pubsub from 'k6/x/pubsub';

const publisher = new pubsub.Publisher({
    ProjectID: 'ProjectID',
    DoNotCreateTopicIfMissing: true,
});

export default function () {
    let error = publisher.publish('topic_name', 'message data');

    check(error, {
        "is sent": err => err === undefined
    });
}
```

Result output:

:warning: **You will receive an error**: `could not check if topic <topic_name> exists: context deadline exceeded`

```
$ ./k6 run example.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: test.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

running (00m05.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m05.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=5s min=5s med=5s max=5s p(90)=5s p(95)=5s
     iterations...........: 1   0.199914/s
     vus..................: 1   min=1 max=1
     vus_max..............: 1   min=1 max=1
```
