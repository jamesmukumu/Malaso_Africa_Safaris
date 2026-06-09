package ratescontroller

import (
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/rates"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"time"
)

func CreateRates(res http.ResponseWriter,req *http.Request) {
time.Sleep(time.Nanosecond * 4)
user_id := <-helpers.UserIDChannel
req.ParseMultipartForm(20 << 30)
file, header, err := req.FormFile("ratesFile")
if err != nil {
	http.Error(res, err.Error(), http.StatusInternalServerError)
	return
}
defer file.Close()


if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
	http.Error(res, "Only ZIP files are allowed", http.StatusBadRequest)
	return
}

buf := make([]byte, 512)
_, err = file.Read(buf)
if err != nil {
	http.Error(res, "Failed to read file", http.StatusBadRequest)
	return
}


file.Seek(0, 0)

mimeType := http.DetectContentType(buf)
if mimeType != "application/zip" && mimeType != "application/octet-stream" {
	http.Error(res, "Invalid ZIP file", http.StatusBadRequest)
	return
}
// create a local directory 
err1 := os.MkdirAll("rates",os.ModePerm)

if err1 != nil {
log.Fatal(err1.Error())
http.Error(res,err1.Error(),500)
return
}
localPath := filepath.Join("rates",header.Filename)
rateFile,err2 := os.Create(localPath)
if err2 != nil {
log.Fatal(err2.Error())
http.Error(res,err2.Error(),500)
return
}
defer rateFile.Close()
hotel_id,_ :=strconv.Atoi(req.FormValue("hotelID"))

var rts = rates.Rates{
Title: req.FormValue("title"),
Description: req.FormValue("description"),
FilePath: localPath,
UserID: user_id,
HotelsID: uint(hotel_id) ,
RatesYear: req.FormValue("year"),
}
op := db.Db_Connection.Create(&rts)
if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]string{
"message":"Rates Added",
})
}else{
http.Error(res,op.Error.Error(),500)
}
}




func RatesChecker(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 4)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var hotels []hotels.Hotels
op := db.Db_Connection.Table("hotels").Joins("LEFT JOIN rates ON rates.hotels_id = hotels.id").Where("rates.id IS NULL").Find(&hotels)
if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"message":"Records Found",
"data":hotels,
"count":len(hotels),
})
}else{
json.NewEncoder(res).Encode(map[string]any{
"message":"Records Not Found",
"data":hotels,
})	
}
}



