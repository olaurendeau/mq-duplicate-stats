package main

import (
    "fmt"
    "os"
    "log"
    "io"
    "crypto/sha1"

    redis "gopkg.in/redis.v4"
    consumer "github.com/olaurendeau/tailmq/consumer"
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    _, err := client.Ping().Result()
    failOnError(err, "Fail connecting to redis")

    sortedSet := os.Getenv("RABBITMQ_EXCHANGE")

    c := consumer.New(os.Getenv("RABBITMQ_URL"), os.Getenv("RABBITMQ_EXCHANGE"))
    go c.Start()
    defer c.Stop()

    for {
        select {
            case d := <-c.Deliveries:

                routingKey, lookUp := os.LookupEnv("RABBITMQ_ROUTING_KEY")
                if (lookUp && routingKey != d.RoutingKey) {
                    continue
                }

                h := sha1.New()
                io.WriteString(h, d.RoutingKey)
                io.WriteString(h, fmt.Sprintf("%s", d.Body))
                fmt.Printf("%s ", d.RoutingKey)
                fmt.Printf("%x \n\n", h.Sum(nil))

                member := fmt.Sprintf("%x", h.Sum(nil))

                score, err := client.ZIncrBy(sortedSet, 1, member).Result()
                failOnError(err, "Fail incrementing member in redis")

                if (score >= 2) {
                    err := client.Set(member, fmt.Sprintf("%s", d), 0).Err()
                    failOnError(err, "Fail setting value in redis")
                }

            case err := <-c.Err:
                failOnError(err, "Fail consuming")
        }
    }
}


func failOnError(err error, msg string) {
  if err != nil {
  log.Fatalf("%s: %s", msg, err)
  }
}