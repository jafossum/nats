# NATS Playground

![](./resources/network.png)

## Environment

    docker compose up -d
    docker compose down

## NATS Context

    nats context add demo-cloud --server nats://acc:acc@localhost:4222 --description "NATS Demo Cloud"
    nats context add demo-edge --server nats://acc:acc@localhost:4224 --description "NATS Demo Edge" --select
    nats context add demo-mid --server nats://acc:acc@localhost:4223 --description "NATS Demo Mid"

    nats context ls

## Streams - Config

    nats stream add timeseries --js-domain edge --context demo-edge
    nats stream add --js-domain mid --context demo-mid --source timeseries
    nats stream add --js-domain cloud --context demo-cloud --source timeseries-mid


    nats stream info --js-domain cloud timeseries -j > cloud/cloud-stream.json
    nats stream info --js-domain mid timeseries-mid -j > mid/mid-stream.json
    nats stream info --js-domain edge timeseries -j > edge/edge-stream.json


    nats stream add --config=./edge/edge-stream.json --js-domain edge --context demo-edge
    nats stream add --config=./mid/mid-stream.json --js-domain mid --context demo-mid
    nats stream add --config=./cloud/cloud-stream.json --js-domain cloud --context demo-cloud


    watch -n 1 nats stream report --context demo-cloud
    watch -n 1 nats stream report --context demo-mid
    watch -n 1 nats stream report --context demo-edge

## Publish and Subscribe to Messages

    nats pub ts.yess --count=1000 --sleep 10ms "publication #{{Count}} @ {{TimeStamp}}"

    nats pub ts.yess --count=10000 "publication #{{Count}} @ {{TimeStamp}}"

    nats sub --context demo-cloud "ts.*"

## Run Example GOlang Publisher

    cd services/publisher/

    go run .

## Publish MQTT From Edge - Receive in Cloud

    docker compose exec -it mosquitto-cloud /bin/sh
    mosquitto_sub -h cloud -u acc -P acc -t ts/mqtt

    docker compose exec -it mosquitto-edge /bin/sh
    mosquitto_pub -h edge -u acc -P acc -t ts/mqtt -q 1 -m hello

## Stop and Start Servers

    docker compose stop mid
    docker compose start mid

    docker compose stop cloud
    docker compose start cloud
