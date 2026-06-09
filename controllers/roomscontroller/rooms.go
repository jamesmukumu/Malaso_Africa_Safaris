package roomscontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/rooms"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CreateRoom(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	hotel_id, _ := strconv.Atoi(req.FormValue("hotelID"))
	ctx := context.Background()
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}

	roomPhoto, _, _ := req.FormFile("roomPhoto")
	result, _ := cloud.Upload.Upload(ctx, roomPhoto, uploader.UploadParams{
		Folder: "rooms",
	})
	roomPhotos := req.MultipartForm.File["roomPhotos"]
	var resultingRoomPhotos []string
	for _, head := range roomPhotos {
		roomPhotoFile, errFile := head.Open()
		defer roomPhotoFile.Close()
		if errFile != nil {
			log.Fatal(errFile.Error())
			return
		}
		resultUpload, _ := cloud.Upload.Upload(ctx, roomPhotoFile, uploader.UploadParams{})
		resultingRoomPhotos = append(resultingRoomPhotos, resultUpload.SecureURL)
	}

	photoBytes, _ := json.Marshal(resultingRoomPhotos)
	roomRates, _ := json.Marshal(req.FormValue("roomRates"))
	mealPlans, _ := json.Marshal(req.FormValue("mealPlans"))
	var room = rooms.Rooms{
		HotelsID:        uint(hotel_id),
		UserID:          user_id,
		RoomRates:       roomRates,
		RoomMeals:       mealPlans,
		RoomName:        req.FormValue("roomName"),
		RoomOverview:    req.FormValue("roomOverview"),
		RoomDescription: req.FormValue("roomDescription"),
		RoomPhoto:       result.SecureURL,
		RoomPhotos:      photoBytes,
	}

	op := db.Db_Connection.Create(&room)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Room Created",
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Room Not Created",
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}
}

func FetchMyRooms(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var rooms []rooms.Rooms
	db.Db_Connection.Preload("Hotels").Find(&rooms)
	json.NewEncoder(res).Encode(map[string]any{
		"data":    rooms,
		"message": "Rooms fetched",
	})
}
