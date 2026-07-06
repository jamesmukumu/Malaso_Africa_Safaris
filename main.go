package main

import (
	cronsjobscontroller "jamesmukumu/maasaimaratripsportal/controllers/cronsJobsController"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	"jamesmukumu/maasaimaratripsportal/router"

	"github.com/robfig/cron/v3"
)

func main() {
	db.DBConnection()
	redisdb.ConnectRedis()
	cronsjobscontroller.HydrateCacheLogins()

	var c *cron.Cron = cron.New()
	c.AddFunc("@every 15m", cronsjobscontroller.HydrateCacheLogins)
	c.AddFunc("@weekly", cronsjobscontroller.SyncEmailToBulk)
	c.Start()
	defer c.Stop()

	router.Router()

}
