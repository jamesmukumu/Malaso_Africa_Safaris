package attractioncontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/scopes"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	"jamesmukumu/maasaimaratripsportal/models/attractions"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/packages"
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
)
var keepFilename bool = false
func SaveAttraction(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
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
	attractionPhoto, _, _ := req.FormFile("attractionPhoto")
	result, _ := cloud.Upload.Upload(ctx, attractionPhoto, uploader.UploadParams{
		Folder: "attractions",
		UseFilename: &keepFilename,
	})
	var attractionPhotos []string
	attractionPhotosMap := req.MultipartForm.File["attractionPhotos"]

	for _, photo := range attractionPhotosMap {
		photoFile, errPhoto := photo.Open()
		if errPhoto != nil {
			log.Fatal(errPhoto.Error())
			continue
		}
		rslt, _ := cloud.Upload.Upload(ctx, photoFile, uploader.UploadParams{
			Folder: "attractions",
			UseFilename: &keepFilename,
		})
		attractionPhotos = append(attractionPhotos, rslt.SecureURL)
	}
	var slug = slug.SlugMaker(req.FormValue("name"))
	destination_id, _ := strconv.Atoi(req.FormValue("destinationID"))
	bytesPhotos, _ := json.Marshal(attractionPhotos)
	var attraction = attractions.Attraction{
		AttractionName:        req.FormValue("name"),
		AttractionSlug:        slug,
		AttractionOverview:    req.FormValue("overview"),
		AttractionDescription: req.FormValue("description"),
		AttractionPhoto:       result.SecureURL,
		DestinationID:         uint(destination_id),
		AttractionPhotos:      bytesPhotos,
		UserID:                user_id,
	}
	op := db.Db_Connection.Create(&attraction)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Attraction Saved",
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Attraction Not Saved",
		})
	} else {
		http.Error(res, "Something Went Wrong", 500)
	}

}



func EditAttraction(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel

	err := req.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	// REQUIRED
	attractionID := req.URL.Query().Get("ID")
	if attractionID == "" {
		http.Error(res, "attractionID is required", 400)
		return
	}

	// Fetch existing attraction
	var attraction attractions.Attraction
	if err := db.Db_Connection.First(&attraction, attractionID).Error; err != nil {
		http.Error(res, "Attraction Not Found", 404)
		return
	}

	// Cloudinary
	ctx := context.Background()
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		http.Error(res, errCloudinary.Error(), 500)
		return
	}

	// ======================
	// UPDATE BASIC FIELDS
	// ======================

	if v := req.FormValue("name"); v != "" {
		attraction.AttractionName = v
		attraction.AttractionSlug = slug.SlugMaker(v)
	}

	if v := req.FormValue("overview"); v != "" {
		attraction.AttractionOverview = v
	}

	if v := req.FormValue("description"); v != "" {
		attraction.AttractionDescription = v
	}

	if v := req.FormValue("destinationID"); v != "" {
		destID, _ := strconv.Atoi(v)
		attraction.DestinationID = uint(destID)
	}

	attraction.UserID = user_id

	// ======================
	// UPDATE MAIN PHOTO
	// ======================
	mainPhoto, _, errMainPhoto := req.FormFile("attractionPhoto")
	if errMainPhoto == nil {

		// DELETE OLD IMAGE
		if attraction.AttractionPhoto != "" {
			publicID := GetPublicIDFromURL(attraction.AttractionPhoto)
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		}

		// UPLOAD NEW
		upload, errUpload := cloud.Upload.Upload(ctx, mainPhoto, uploader.UploadParams{
			Folder: "attractions",
			UseFilename: &keepFilename,
		})
		if errUpload == nil {
			attraction.AttractionPhoto = upload.SecureURL
		}
	}

	// ======================
	// UPDATE MULTIPLE PHOTOS
	// ======================
	newPhotos := req.MultipartForm.File["attractionPhotos"]

	if len(newPhotos) > 0 {

		// Delete existing photos
		var oldPhotos []string
		json.Unmarshal(attraction.AttractionPhotos, &oldPhotos)

		for _, url := range oldPhotos {
			publicID := GetPublicIDFromURL(url)
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		}

		// Upload new photos
		var uploadedPhotos []string
		for _, fileHeader := range newPhotos {
			file, errOpen := fileHeader.Open()
			if errOpen != nil {
				continue
			}
			defer file.Close()

			up, errUp := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
				Folder: "attractions",
				UseFilename: &keepFilename,
			})
			if errUp == nil {
				uploadedPhotos = append(uploadedPhotos, up.SecureURL)
			}
		}

		// Save new photos array
		bytesPhotos, _ := json.Marshal(uploadedPhotos)
		attraction.AttractionPhotos = bytesPhotos
	}

	// ======================
	// SAVE CHANGES
	// ======================
	if err := db.Db_Connection.Save(&attraction).Error; err != nil {
		http.Error(res, "Attraction Not Updated", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Attraction Updated Successfully",
	})
}


func GetPublicIDFromURL(url string) string {
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func FetchMyAttractions(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var attractions []attractions.Attraction
	op := db.Db_Connection.Find(&attractions)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"data":    attractions,
			"message": "Attractions Found",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]any{
			"data":    attractions,
			"message": "Attractions Found",
		})
	}

}

func PublishAttraction(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	dest_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(dest_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("attractions").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("attractions").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})

	}
}


func FetchSingularAttraction(res http.ResponseWriter,req *http.Request){
var attraction attractions.Attraction
var hotels []hotels.Hotels
var packages []packages.Package
var attractionSlug string = mux.Vars(req)["slug"]
op := db.Db_Connection.Where("attraction_slug =? ",attractionSlug).Preload("Destination").Find(&attraction)
db.Db_Connection.Where("published=?",true).Where("destination_id =?",attraction.DestinationID).Limit(6).Find(&hotels)
db.Db_Connection.Where("published=?",true).Where("destination_id =?",attraction.DestinationID).Find(&packages)

if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"message":"Attraction Found",
"data":attraction,
"hotels":hotels,
"packages":packages,
})
}else{
json.NewEncoder(res).Encode(map[string]any{
		"message":"Attraction Not Found",
		"data":attraction,
		})	
}


}


func FetchActiveAttractions(res http.ResponseWriter,req *http.Request){
var attractions []attractions.Attraction

pageString := req.URL.Query().Get("pageString")
    page := 1

    if pageString != ""{
    intPage,_:= strconv.Atoi(pageString)
     page = intPage
	} 

baseQuery:= db.Db_Connection.Where("published =?",true).Preload("Destination")
paginatedQuery := scopes.Paginator(baseQuery,page)
op := paginatedQuery.Find(&attractions)
if op.RowsAffected > 0 && op.Error == nil {
var totalCount int64
db.Db_Connection.Table("attractions").Where("published =?",true).Count(&totalCount)
json.NewEncoder(res).Encode(map[string]any{
"data":attractions,
"message":"Attraction Fetched",
"count":totalCount,
})  //
}else{
	json.NewEncoder(res).Encode(map[string]any{
		"data":attractions,
		"message":"Attraction Not Fetched",
		})	
}

}