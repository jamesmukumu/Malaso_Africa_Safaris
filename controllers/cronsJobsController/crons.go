package cronsjobscontroller

import (
	"context"
	"fmt"
	"log"
	"time"
	"crypto/rand"
	"math/big"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	"jamesmukumu/maasaimaratripsportal/models/bulkmails"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/logins"
)

func SyncEmailToBulk() {
	var hotels []hotels.Hotels
	op := db.Db_Connection.Find(&hotels)
	if op.RowsAffected > 0 && op.Error == nil {
		for _, hotel := range hotels {
			var bulkEmail = bulkmails.BulkMails{
				UserName:    hotel.HotelName,
				Email:       hotel.HotelContactEmail,
				PhoneNumber: hotel.HotelContactPhonenumber,
			}
			OP := db.Db_Connection.Create(&bulkEmail)
			if OP.Error != nil {
				fmt.Println(OP.Error)
				continue
			}
		}

		fmt.Println("Emails SuccessFully Synced")
	}

}

func GenerateHydrationKey() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	key := "mala"

	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		key += string(chars[n.Int64()])
	}

	return key
}

func HydrateCacheLogins() {
	if redisdb.ClientRedis == nil || db.Db_Connection == nil {
		log.Println("redis or db is not ready for login hydration")
		return
	}

	ctx := context.TODO()
	if err := redisdb.ClientRedis.Ping(ctx).Err(); err != nil {
		log.Println("redis ping failed:", err)
		return
	}

	var cachedLogins []logins.Logins
	if err := db.Db_Connection.
		Where("expiry_time > ?", time.Now().UTC()).
		Find(&cachedLogins).Error; err != nil {
		log.Println("failed to fetch login sessions:", err)
		return
	}

	for _, session := range cachedLogins {
		if session.LoginID == "" || session.TokenString == "" {
			continue
		}

		if err := redisdb.ClientRedis.HSet(ctx, session.LoginID, "token", session.TokenString).Err(); err != nil {
			log.Println("failed to hydrate login session:", err)
			continue
		}

		if err := redisdb.ClientRedis.ExpireAt(ctx, session.LoginID, session.ExpiryTime).Err(); err != nil {
			log.Println("failed to set login expiry:", err)
		}
	}

	log.Printf("hydrated %d login sessions into redis", len(cachedLogins))
	//generate some unique key to store in redis to indicate that the hydration has been done
	key := GenerateHydrationKey()
	redisdb.ClientRedis.Set(context.Background(),key,"done",0)

}
