package destinationcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/scopes"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/packages"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var keepFileName bool = false
func EditDestination(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Millisecond * 5)
     user_id := <-helpers.UserIDChannel
	 fmt.Println(user_id)
	// Parse form
	err := req.ParseMultipartForm(50 << 20) // 50MB
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
		http.Error(res, "Cloudinary error", 500)
		return
	}

	ctx := context.Background()

	// --- Find destination in DB ---
	id := req.URL.Query().Get("id")
	var destination destinations.Destination
	if err := db.Db_Connection.First(&destination, "id = ?", id).Error; err != nil {
		http.Error(res, "Destination not found", 404)
		return
	}

	// ----------- UPDATE FIELDS ONLY IF THEY EXIST ------------

	// Title
	if title := req.FormValue("destinationTitle"); title != "" {
		destination.DestinationTitle = title
		destination.DestinationSlug = slug.SlugMaker(req.FormValue("destinationTitle"))
	}

	// Overview
	if overview := req.FormValue("destinationOverview"); overview != "" {
		destination.DestinationOverview = overview
	}

	if desc := req.FormValue("destinationDescription"); desc != "" {
		destination.DestinationDescription = desc
	}

	
	file, _, errPhoto := req.FormFile("destinationPhoto")
	if errPhoto == nil {
	
		oldPublicID := extractPublicID(destination.DestinationPhoto)
		if oldPublicID != "" {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: oldPublicID})
		}

		// Upload new photo
		uploadRes, errUp := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "destinations",
		UseFilename: &keepFileName,
		})
		if errUp == nil {
			destination.DestinationPhoto = uploadRes.SecureURL
		}
	}
	
	multiFiles := req.MultipartForm.File["destinationPhotos"]

	if len(multiFiles) > 0 {
		
		var oldPhotos []string
		json.Unmarshal(destination.DestinationPhotos, &oldPhotos)

		for _, photo := range oldPhotos {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: extractPublicID(photo),
			})
		}

		// Upload new photos
		var newPhotos []string
		for _, header := range multiFiles {
			f, err := header.Open()
			if err != nil {
				continue
			}
			defer f.Close()

			upRes, err2 := cloud.Upload.Upload(ctx, f, uploader.UploadParams{
			Folder: "destinations",
			UseFilename: &keepFileName,
			})
			if err2 == nil {
				newPhotos = append(newPhotos, upRes.SecureURL)
			}
		}

		// Save new array
		newJSON, _ := json.Marshal(newPhotos)
		destination.DestinationPhotos = datatypes.JSON(newJSON)
	}
	if err := db.Db_Connection.Save(&destination).Error; err != nil {
		http.Error(res, "Update failed", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Destination Updated Successfully",
	})
}

func extractPublicID(url string) string {
	if url == "" {
		return ""
	}

	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]           
	publicID := strings.TrimSuffix(last, filepath.Ext(last))
	return publicID
}

func CreateDestination(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 9)
	userid := <-helpers.UserIDChannel
	err := req.ParseMultipartForm(30 << 45)
	if err != nil {
		log.Fatal(err.Error())
	}
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}
	ctx := context.Background()
	file, _, errFile := req.FormFile("destinationPhoto")
	if errFile != nil {
		log.Fatal(errFile.Error())
	}
	result, errUpload := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
	Folder: "destinations",
	UseFilename: &keepFileName,
	})
	if errUpload != nil {
		panic(errUpload.Error())
	}
	destinationFilesPhoto := req.MultipartForm.File["destinationPhotos"]
	var photoDestinations []string
	for _, header := range destinationFilesPhoto {
		f, err1 := header.Open()
		if err1 != nil {
			http.Error(res, err1.Error(), 500)
			continue
		}

		defer f.Close()
		uploadResult, err2 := cloud.Upload.Upload(ctx, f, uploader.UploadParams{})
		if err2 != nil {
			panic(err2.Error())
		}
		photoDestinations = append(photoDestinations, uploadResult.SecureURL)
	}
	photoJSON, _ := json.Marshal(photoDestinations)
	var destination = destinations.Destination{
		DestinationTitle:       req.FormValue("destinationTitle"),
		DestinationOverview:    req.FormValue("destinationOverview"),
		DestinationPhoto:       result.SecureURL,
		DestinationDescription: req.FormValue("destinationDescription"),
		DestinationPhotos:      datatypes.JSON(photoJSON),
		UserID:                 userid,
		DestinationSlug: slug.SlugMaker(req.FormValue("destinationTitle")),
	}
	op := db.Db_Connection.Create(&destination)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Destination Saved",
		})
	} else if op.Error != nil {
		http.Error(res, "Something Went Wrong", 500)
	}
}



func FetchDisplayDestinations(res http.ResponseWriter, req *http.Request) {
	var destinations []destinations.Destination
	pageString := req.URL.Query().Get("pageString")
    page := 1

    if pageString != ""{
    intPage,_:= strconv.Atoi(pageString)
     page = intPage
	} 
	var toursCount int64
	var destinationSlice []any
	baseQuery := db.Db_Connection.Where("published=?", true)
	paginatedQuery := scopes.Paginator(baseQuery,page)
 op := paginatedQuery.Find(&destinations)
    for _, d := range destinations {
	db.Db_Connection.Table("packages").Where("destination_id=?",d.ID).Count(&toursCount)
destinationSlice = append(destinationSlice, map[string]any{
"data":d,
"count":toursCount,
})	
}
 
	if op.RowsAffected > 0 && op.Error == nil {
		var totalDestinationsCount int64
		db.Db_Connection.Table("destinations").Where("published =? ",true).Count(&totalDestinationsCount)
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Destination Fetched",
			"data":    destinationSlice,
			"count":totalDestinationsCount,
		})
	} else {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Destination Fetched",
			"data":    destinations,
		})
	}
}

func FetchMyDestinations(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 4)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var destinations []destinations.Destination
	op := db.Db_Connection.Find(&destinations)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Destinations Fetched",
			"data":    destinations,
		})
	} else if op.Error == nil && op.RowsAffected == 0 {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Empty Destinations",
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}
}

func PublishDestination(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	dest_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(dest_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("destinations").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("destinations").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	}
}



func FetchDestination(res http.ResponseWriter,req *http.Request){
var dest destinations.Destination
var packages []packages.Package
var hotels []hotels.Hotels
var slug = mux.Vars(req)["slug"]
db.Db_Connection.Where("published=?",true).Where("destination_slug=?",slug).Find(&dest)
db.Db_Connection.Where("published=?",true).Where("destination_id=?",dest.ID).Order("created_at DESC").Limit(4).Find(&hotels)
db.Db_Connection.Where("published=?",true).Where("destination_id=?",dest.ID).Find(&packages)
json.NewEncoder(res).Encode(map[string]any{
"message":"Destination Found",
"data":dest,
"packages":packages,
"hotels":hotels,
})
}


func DestinationSlugs(res http.ResponseWriter,req *http.Request){
var destinations []destinations.Destination
db.Db_Connection.Select("destination_slug").Find(&destinations)
json.NewEncoder(res).Encode(map[string]any{
"data":destinations,
})
}