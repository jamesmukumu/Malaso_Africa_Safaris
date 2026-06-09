package hotelscontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/scopes"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/rooms"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)
var keepFilename bool = false

func EditHotel(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Millisecond * 5)
   user_id := <-helpers.UserIDChannel
   fmt.Println(user_id)
	// Parse multipart form
	err := req.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(res, "Invalid form data", 400)
		return
	}

	// Cloudinary setup
	cloud, errCloud := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloud != nil {
		http.Error(res, "Cloudinary init failed", 500)
		return
	}

	ctx := context.Background()

	// ---------------- FIND HOTEL ----------------
	hotelID := req.URL.Query().Get("id")
	var hotel hotels.Hotels

	if err := db.Db_Connection.First(&hotel, "id = ?", hotelID).Error; err != nil {
		http.Error(res, "Hotel not found", 404)
		return
	}

	// -------------- UPDATE BASIC FIELDS --------------
	if v := req.FormValue("hotelName"); v != "" {
		hotel.HotelName = v
		hotel.Slug = slug.SlugMaker(req.FormValue("hotelName"))
	}

	if v := req.FormValue("hotelOverview"); v != "" {
		hotel.HotelOverview = v
	}

	if v := req.FormValue("hotelDescription"); v != "" {
		hotel.HotelDescription = v
	}

	if v := req.FormValue("hotelEmail"); v != "" {
		hotel.HotelContactEmail = v
	}

	if v := req.FormValue("hotelPhonenumber"); v != "" {
		hotel.HotelContactPhonenumber = v
	}

	if v := req.FormValue("contactPerson"); v != "" {
		hotel.ContactPerson = v
	}

	if v := req.FormValue("cancellationPolicy"); v != "" {
		hotel.CancellationPolicy = v
	}

	if v := req.FormValue("locationDescription"); v != "" {
		hotel.LocationDescription = v
	}

	// Ratings
	if v := req.FormValue("ratings"); v != "" {
		r, _ := strconv.Atoi(v)
		hotel.Ratings = int8(r)
	}

	// Coordinates
	if v := req.FormValue("hotelCoordinates"); v != "" {
		b, _ := json.Marshal(v)
		hotel.HotelCoordinates = b
	}

	// Hotel Rates
	if v := req.FormValue("hotelRates"); v != "" {
		b, _ := json.Marshal(v)
		hotel.HotelRates = b
	}

	// Contacts
	if v := req.FormValue("hotelContacts"); v != "" {
		b, _ := json.Marshal(v)
		hotel.HotelContacts = b
	}

	// Destination
	if v := req.FormValue("destinationID"); v != "" {
		n, _ := strconv.Atoi(v)
		hotel.DestinationID = uint(n)
	}

	// =====================================================
	// HANDLE HOTEL THUMBNAIL (Single Image)
	// =====================================================
	thumbFile, _, errThumb := req.FormFile("hotelPhoto")
	if errThumb == nil {
		// Delete old thumbnail
		oldID := extractPublicID(hotel.HotelThumbnail)
		if oldID != "" {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: oldID})
		}

		// Upload new one
		up, errUpload := cloud.Upload.Upload(ctx, thumbFile, uploader.UploadParams{
			Folder: "hotels",
			UseFilename: &keepFilename,
		})
		if errUpload == nil {
			hotel.HotelThumbnail = up.SecureURL
		}
	}
	// If no file → keep old thumbnail

	// =====================================================
	// HANDLE HOTEL IMAGES (Multiple)
	// =====================================================
	images := req.MultipartForm.File["hotelImages"]

	if len(images) > 0 {
		// Delete old hotel images
		var old []string
		json.Unmarshal(hotel.HotelPhotos, &old)

		for _, url := range old {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: extractPublicID(url),
			})
		}

		// Upload new images
		var newImgs []string
		for _, hdr := range images {
			f, err := hdr.Open()
			if err != nil {
				continue
			}
			defer f.Close()

			up, err2 := cloud.Upload.Upload(ctx, f, uploader.UploadParams{})
			if err2 == nil {
				newImgs = append(newImgs, up.SecureURL)
			}
		}

		// Save array
		b, _ := json.Marshal(newImgs)
		hotel.HotelPhotos = b
	}

	// ---------------- SAVE CHANGES ----------------
	if err := db.Db_Connection.Save(&hotel).Error; err != nil {
		http.Error(res, "Update failed", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Hotel Updated Successfully",
	})
}

func extractPublicID(url string) string {
	if url == "" {
		return ""
	}

	parts := strings.Split(url, "/")
	last := parts[len(parts)-1] // example: "photo_abc123.jpg"
	publicID := strings.TrimSuffix(last, filepath.Ext(last))
	return publicID
}

func AddHotels(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 9)
	user_id := <-helpers.UserIDChannel
	ctx := context.Background()
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}
	hotelPhotoFile, _, _ := req.FormFile("hotelPhoto")
	hotelPhoto, errPhoto := cloud.Upload.Upload(ctx, hotelPhotoFile, uploader.UploadParams{
		Folder: "hotels",
		UseFilename: &keepFilename,
	})
	if errPhoto != nil {
		log.Fatal(errPhoto.Error())
		return
	}

	hotelImages := req.MultipartForm.File["hotelImages"]
	var hotelImagesSlice []string
	for _, head := range hotelImages {
		file, errOpenFile := head.Open()
		defer file.Close()
		if errOpenFile != nil {
			http.Error(res, errOpenFile.Error(), 500)
			continue
		}
		hotelImagesResult, _ := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder: "hotels",
			UseFilename: &keepFilename,
		})
		hotelImagesSlice = append(hotelImagesSlice, hotelImagesResult.SecureURL)
	}

	dest := req.FormValue("destinationID")
	dest_id, _ := strconv.Atoi(dest)
	ratings := req.FormValue("ratings")
	ratings_id, _ := strconv.Atoi(ratings)
	hotelCoordinates := req.FormValue("hotelCoordinates")
	hotelBytes, _ := json.Marshal(hotelCoordinates)
	hotelImagesbytes, _ := json.Marshal(hotelImagesSlice)

	slug := slug.SlugMaker(req.FormValue("hotelName"))
	hotelRates, _ := json.Marshal(req.FormValue("hotelRates"))
	hotelContacts, _ := json.Marshal(req.FormValue("hotelContacts"))
	var hotel = hotels.Hotels{
		HotelName:               req.FormValue("hotelName"),
		HotelContacts:           hotelContacts,
		HotelOverview:           req.FormValue("hotelOverview"),
		HotelDescription:        req.FormValue("hotelDescription"),
		HotelContactEmail:       req.FormValue("hotelEmail"),
		HotelContactPhonenumber: req.FormValue("hotelPhonenumber"),
		DestinationID:           uint(dest_id),
		ContactPerson:           req.FormValue("contactPerson"),
		CancellationPolicy:      req.FormValue("cancellationPolicy"),
		Ratings:                 int8(ratings_id),
		LocationDescription:     req.FormValue("locationDescription"),
		UserID:                  user_id,
		HotelCoordinates:        hotelBytes,
		HotelThumbnail:          hotelPhoto.SecureURL,
		HotelPhotos:             hotelImagesbytes,
		Slug:                    slug,
		HotelRates:              hotelRates,
	}

	op := db.Db_Connection.Create(&hotel)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Hotel Created",
		})
	} else if op.RowsAffected <= 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Hotel Created",
		})
	} else {
		http.Error(res, op.Error.Error(), 500)
	}
}



func FetchMyHotels(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	var hotels []hotels.Hotels
     pageString := req.URL.Query().Get("page")
     var page int = 1
     if pageString != ""{
	  pageInt,_ := strconv.Atoi(pageString)
	  page  = pageInt
	 }

	fmt.Println(user_id)
	baseQuery := db.Db_Connection.Table("hotels")
	paginatedQuery := scopes.Paginator(baseQuery,page)
	op := paginatedQuery.Find(&hotels)
	if op.RowsAffected > 0 && op.Error == nil {
		var hotelsCount int64
		db.Db_Connection.Table("hotels").Count(&hotelsCount)
		json.NewEncoder(res).Encode(map[string]any{
			"data":    hotels,
			"message": "Hotels Fetched",
			"count":hotelsCount,
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}

}

func PublishHotel(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	dest_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(dest_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("hotels").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("hotels").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})

	}
}

func FetchDisplayHotels(res http.ResponseWriter, req *http.Request) {
	var hotels []hotels.Hotels

	pageString := req.URL.Query().Get("pageString")
    page := 1

    if pageString != ""{
    intPage,_:= strconv.Atoi(pageString)
     page = intPage
	} 
	
	baseQuery := db.Db_Connection.Where("published=?", true)
	paginatedQuery := scopes.Paginator(baseQuery,page)
	op:= paginatedQuery.Find(&hotels)
	if op.RowsAffected > 0 && op.Error == nil {
		var totalCount int64
	
   db.Db_Connection.Table("hotels").Where("published =?",true).Count(&totalCount)
		json.NewEncoder(res).Encode(map[string]any{
			"data":    hotels,
			"message": "Hotels Found",
			"count":totalCount,
		})
	} else {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Hotels Found",
		})
	}

}



func FetchHotel(res http.ResponseWriter,req *http.Request){
var hotel hotels.Hotels
var rooms []rooms.Rooms
slug := mux.Vars(req)["slug"]
db.Db_Connection.Where("published=?",true).Where("slug =?",slug).Find(&hotel)
db.Db_Connection.Where("hotels_id = ?",hotel.ID).Find(&rooms)
json.NewEncoder(res).Encode(map[string]any{
"data":hotel,
"rooms":rooms,
})
}

func FetchHotelSlugs(res http.ResponseWriter,req *http.Request){
var hotels []hotels.Hotels
db.Db_Connection.Select("slug").Find(&hotels)
json.NewEncoder(res).Encode(map[string]any{
"data":hotels,
})
}



