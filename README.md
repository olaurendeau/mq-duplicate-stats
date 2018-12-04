# mq-duplicate-stats
Store hash of messages and provide duplication stats

# Usage

Gather data
```
docker-compose up -d redis
docker-compose run -e RABBITMQ_URL=amqp://guest:guest@localhost// -e RABBITMQ_EXCHANGE=exchange -e RABBITMQ_ROUTING_KEY=routing_key go go run main.go
```

Explore duplicates
```
docker-compose run redis redis-cli -h redis ZREVRANGEBYSCORE exchange 1000 2 WITHSCORES
docker-compose run redis redis-cli -h redis GET {key}
```