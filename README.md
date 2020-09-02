# tyk-kafka-plugin

Example golang plugin for kafka

## Start Kafka and Zookeeper (Assumes brew installation on Mac)
`zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties`

## Create kafka topic
`kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test`

## Set topic retention time to 30s (default retention time is 24 hours)
`kafka-configs --zookeeper localhost:2181 --alter --entity-type topics --entity-name test --add-config retention.ms=30000`

## Import API definition
Import `./apidefinitions/kafka-test.json` into Tyk Dashboard as a new API

## Copy plugin to directory on Gateway 
`docker cp kafka.so tyk-pro-docker-demo_tyk-gateway_1:/opt/tyk-gateway/middleware`

## Restart Gateway
`docker restart tyk-pro-docker-demo_tyk-dashboard_1`

## See it working!

### Produce sample json msg on kafka topic
`jq -rc . sample.json | kafka-console-producer --broker-list localhost:9092 --topic test`
    
### Send request to service on Tyk
`$ curl http://localhost:8080/kafka-test/get
{
    "id": "123",
    "status": "Active"
}`