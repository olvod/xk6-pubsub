# xk6-pubsub

This is a [k6](https://go.k6.io/k6) extension using the [xk6](https://github.com/k6io/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Install `xk6` framework for extending `k6`:
```shell
go install github.com/k6io/xk6/cmd/xk6@latest
```

2. Build the binary:
```shell
#Required by Pub/Sub client
export CGO_ENABLED=1
```
```shell
xk6 build --with github.com/olvod/xk6-pubsub@latest
```

xk6 build --with github.com/k6io/xk6-redis="/Users/avpretty/pr/xk6-pubsub"

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
```shell
./k6 run example.js
```

Result output:
```
          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: example.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

[watermill] 2021/03/17 22:17:04.450779 publisher.go:116: 	level=TRACE msg="Sending message to Google PubSub" message_uuid=Y3mriTgx4SuWo2ZxMgg8FF topic=topic_name 
[watermill] 2021/03/17 22:17:04.464942 publisher.go:131: 	level=TRACE msg="Message published to Google PubSub" message_uuid=Y3mriTgx4SuWo2ZxMgg8FF topic=topic_name 
[watermill] 2021/03/17 22:17:04.465082 publisher.go:139: 	level=INFO  msg="Closing Google PubSub publisher" 
[watermill] 2021/03/17 22:17:04.465128 publisher.go:153: 	level=INFO  msg="Google PubSub publisher closed" 

running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

     ✓ is sent

     checks...............: 100.00% ✓ 1 ✗ 0
     data_received........: 0 B     0 B/s
     data_sent............: 0 B     0 B/s
     iteration_duration...: avg=27.14ms min=27.14ms med=27.14ms max=27.14ms p(90)=27.14ms p(95)=27.14ms
     iterations...........: 1       35.356928/s
```
