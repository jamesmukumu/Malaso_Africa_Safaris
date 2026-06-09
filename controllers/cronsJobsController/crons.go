package cronsjobscontroller

import (
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/models/bulkmails"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	

)

func SyncEmailToBulk() {
var hotels []hotels.Hotels
op := db.Db_Connection.Find(&hotels)
if op.RowsAffected > 0 && op.Error == nil {
for _, hotel := range hotels {
var bulkEmail  = bulkmails.BulkMails{
UserName: hotel.HotelName,
Email: hotel.HotelContactEmail,
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




