package redisdb

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
)
var ClientRedis *redis.Client
func ConnectRedis(){
godotenv.Load()
redis_string := os.Getenv("redisUrl")
opt,err := redis.ParseURL(redis_string)
if err != nil {
log.Fatal(err.Error())
return
}
client := redis.NewClient(opt)
ClientRedis = client

fmt.Println("Connected to redis Successfully")
}