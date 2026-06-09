package main

import (
	cronsjobscontroller "jamesmukumu/maasaimaratripsportal/controllers/cronsJobsController"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	"jamesmukumu/maasaimaratripsportal/router"
	"sync"

	"github.com/robfig/cron/v3"
)

var wg = &sync.WaitGroup{}

func main() {
var c *cron.Cron = cron.New()
c.AddFunc("@weekly",cronsjobscontroller.SyncEmailToBulk)
c.Start()
defer c.Stop()
wg.Add(3)
defer wg.Wait()
go func(){
defer wg.Done()
router.Router()
}()
go func (){
defer wg.Done();
db.DBConnection()
}()
go func (){
defer wg.Done();
redisdb.ConnectRedis()
}()
	



}
