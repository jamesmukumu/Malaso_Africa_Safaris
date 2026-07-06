package adventurecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	"jamesmukumu/maasaimaratripsportal/models/adventures"
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
var keepFileName bool = false
func GetPublicIDFromURL(url string) string {
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func CreateAdventureCategory(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	ctx := context.Background()
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}
	packageImage, _, _ := req.FormFile("adventureImage")
	uploadResult, err1 := cloud.Upload.Upload(ctx, packageImage, uploader.UploadParams{
	Folder: "adventures",
	UseFilename: &keepFileName,
	})
	if err1 != nil {
		log.Fatal(err1.Error())
		http.Error(res, err1.Error(), 500)
		return
	}
	var adventure_category = adventures.AdventuresCategory{
	AdventureCategoryName: req.FormValue("name"),
	AventureCategoryDescription: req.FormValue("description"),
	UserID: user_id,
	AdventureSlug: strings.ReplaceAll(req.FormValue("name")," ","-"),
	AdventureCategoryImage: uploadResult.SecureURL,
	}
	op := db.Db_Connection.Create(&adventure_category)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Category Saved",
		})
	} else if op.Error != nil {
		http.Error(res, op.Error.Error(), 500)
	}
}



func UpdateAdventureCategory(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)


	user_id := <-helpers.UserIDChannel


	id := req.URL.Query().Get("ID")
	if id == "" {
		http.Error(res, "Category ID is required", http.StatusBadRequest)
		return
	}


	var category adventures.AdventuresCategory
	if err := db.Db_Connection.First(&category, id).Error; err != nil {
		http.Error(res, "Category not found", http.StatusNotFound)
		return
	}


	if err := req.ParseMultipartForm(50 << 20); err != nil {
		http.Error(res, "Invalid form data", http.StatusBadRequest)
		return
	}


	updates := map[string]any{}

	
	if name := strings.TrimSpace(req.FormValue("name")); name != "" {
		updates["adventure_category_name"] = name
		updates["adventure_slug"] = strings.ReplaceAll(strings.ToLower(name), " ", "-")
	}


	if description := strings.TrimSpace(req.FormValue("description")); description != "" {
		updates["aventure_category_description"] = description
	}

	
	file, _, err := req.FormFile("adventureImage")
	if err == nil {
		defer file.Close()

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
var keepFileName bool = true

		uploadResult, errUpload := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder: "adventures",
			UseFilename: &keepFileName,
		})
		if errUpload != nil {
			http.Error(res, errUpload.Error(), 500)
			return
		}

		updates["adventure_category_image"] = uploadResult.SecureURL
	}
fmt.Println(user_id)

	
	if len(updates) == 0 {
		http.Error(res, "No fields provided for update", http.StatusBadRequest)
		return
	}


	if err := db.Db_Connection.Model(&category).Updates(updates).Error; err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Category updated successfully",
	})
}


func FetchShowAdventureCategories(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 5)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var adv_categories []adventures.AdventuresCategory
op := db.Db_Connection.Find(&adv_categories)
if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"data":adv_categories,
"message":"Categories Found",
})
}else{
json.NewEncoder(res).Encode(map[string]any{
"message":"Categories Not Found",
})	
}}



func AdventureCat(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 5)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var adv adventures.AdventuresCategory
slug := mux.Vars(req)["slug"]
op := db.Db_Connection.Where("adventure_slug =?",slug).Find(&adv)
if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"message":"Adventure Found",
"data":adv,
})
}else{
	json.NewEncoder(res).Encode(map[string]any{
		"message":"Adventure Found",
		"data":adv,
		})	
}
}
func EditAdventure(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel

	err := req.ParseMultipartForm(30 << 45)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	// Adventure ID
	adventureID := req.URL.Query().Get("ID")
	if adventureID == "" {
		http.Error(res, "adventureID is required", 400)
		return
	}

	var adventure adventures.Adventure
	if err := db.Db_Connection.First(&adventure, adventureID).Error; err != nil {
		http.Error(res, "Adventure not found", 404)
		return
	}

	// Cloudinary setup
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

	// ================
	//  UPDATE FIELDS
	// ================
	if v := req.FormValue("adventureTitle"); v != "" {
		adventure.Name = v
		adventure.Slug = slug.SlugMaker(v)
	}

	if v := req.FormValue("inclusions"); v != "" {
		adventure.AdventureInclusions = v
	}

	if v := req.FormValue("overview"); v != "" {
		adventure.AdventureOverview = v
	}

	if v := req.FormValue("description"); v != "" {
		adventure.AdventureDescription = v
	}

	if v := req.FormValue("exclusions"); v != "" {
		adventure.AdventureExclusions = v
	}

	if v := req.FormValue("specialNotes"); v != "" {
		adventure.SpecialNotes = v
	}

	if v := req.FormValue("startValidity"); v != "" {
		adventure.StartValidity = v
	}

	if v := req.FormValue("endValidity"); v != "" {
		adventure.EndValidty = v
	}

	if v := req.FormValue("chargeCurrency"); v != "" {
		adventure.AdventureChargeCurrency = v
	}

	if v := req.FormValue("modeTransport"); v != "" {
		adventure.ModeOfTransport = v
	}

	if v := req.FormValue("destinationID"); v != "" {
		destID, _ := strconv.Atoi(v)
		adventure.DestinationID = uint(destID)
	}

	if v := req.FormValue("adventureCategoryID"); v != "" {
		catID, _ := strconv.Atoi(v)
		adventure.AdventuresCategoryID = uint(catID)
	}

	if v := req.FormValue("Amount"); v != "" {
		amount, _ := strconv.Atoi(v)
		adventure.AdventureCharges = amount
	}

	if v := req.FormValue("maxCapacity"); v != "" {
		c, _ := strconv.Atoi(v)
		adventure.MaximumCapacity = c
	}


	mainImage, _, errMain := req.FormFile("adventureImage")
	if errMain == nil {
	
		if adventure.AdventurePhoto != "" {
			publicID := GetPublicIDFromURL(adventure.AdventurePhoto)
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		}

		uploaded, errUpload := cloud.Upload.Upload(ctx, mainImage, uploader.UploadParams{
			Folder: "adventures",
			UseFilename: &keepFileName,
		})
		if errUpload == nil {
			adventure.AdventurePhoto = uploaded.SecureURL
		}
	}


	newPhotos := req.MultipartForm.File["adventureImages"]

	if len(newPhotos) > 0 {
	
		var oldPhotos []string
		json.Unmarshal(adventure.AdventurePhotos, &oldPhotos)

		for _, url := range oldPhotos {
			publicID := GetPublicIDFromURL(url)
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
		}

	
		var uploadedPhotos []string
		for _, head := range newPhotos {
			file, errOpen := head.Open()
			if errOpen != nil {
				continue
			}
			defer file.Close()

			up, errUp := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
				Folder: "adventures",
				UseFilename: &keepFileName,
			})
			if errUp == nil {
				uploadedPhotos = append(uploadedPhotos, up.SecureURL)
			}
		}

		// Save
		photosBytes, _ := json.Marshal(uploadedPhotos)
		adventure.AdventurePhotos = photosBytes
	}

	adventure.UserID = user_id

	// Save changes
	if err := db.Db_Connection.Save(&adventure).Error; err != nil {
		http.Error(res, "Adventure Not Updated", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Adventure Updated Successfully",
	})
}



func CreateAdventure(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	err := req.ParseMultipartForm(30 << 45)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(res,err.Error(),500)
		return
	}
	// itenerary is map[string]interface{}{
	// "title":"string",
	// "description":"string",
	// "mealPlan":"string",
	// "coordinates":map[int]int
	// }
	ctx := context.Background()
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}
	packageImage, _, _ := req.FormFile("adventureImage")
	uploadResult, err1 := cloud.Upload.Upload(ctx, packageImage, uploader.UploadParams{
	Folder: "adventures",
	UseFilename: &keepFileName,
	})
	if err1 != nil {
		log.Fatal(err1.Error())
		http.Error(res, err1.Error(), 500)
		return
	}
	packageImages := req.MultipartForm.File["adventureImages"]
	var photosPackage []string
	for _, head := range packageImages {
		file, errOpenFile := head.Open()
		defer file.Close()
		if errOpenFile != nil {
			http.Error(res, errOpenFile.Error(), 500)
			continue
		}
		packageImagesResult, _ := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "adventures",
		UseFilename: &keepFileName,
		})
		photosPackage = append(photosPackage, packageImagesResult.SecureURL)
	}
	slug := slug.SlugMaker(req.FormValue("adventureTitle"))
	photosBytes, _ := json.Marshal(photosPackage)
	adventureCategory, _ := strconv.Atoi(req.FormValue("adventureCategoryID"))
	destinationID, _ := strconv.Atoi(req.FormValue("destinationID"))
	chargeAmount, _ := strconv.Atoi(req.FormValue("Amount"))
	maxCapacity, _ := strconv.Atoi(req.FormValue("maxCapacity"))
	var adventureSave = adventures.Adventure{
	Name: req.FormValue("adventureTitle"),
	Slug: slug,
	AdventureInclusions: req.FormValue("inclusions"),
    AdventureOverview: req.FormValue("overview"),
    AdventureDescription: req.FormValue("description"),
    DestinationID: uint(destinationID),
	UserID: user_id,
	AdventuresCategoryID: uint(adventureCategory),
    AdventurePhoto: uploadResult.SecureURL,
	AdventureExclusions: req.FormValue("exclusions"),
	AdventureCharges: chargeAmount,
	AdventureChargeCurrency: req.FormValue("chargeCurrency"),
	AdventurePhotos: photosBytes,
	SpecialNotes: req.FormValue("specialNotes"),
	EndValidty: req.FormValue("endValidity"),
    StartValidity: req.FormValue("startValidity"),
    ModeOfTransport: req.FormValue("modeTransport"),
    MaximumCapacity: maxCapacity,
	}
	op := db.Db_Connection.Create(&adventureSave)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Adventure Saved",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Adventure Not Saved",
		})
	}
}
//
func FetchAdventures(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 3)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var adeventures []adventures.Adventure
op := db.Db_Connection.Find(&adeventures)
if op.RowsAffected >0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"data":adeventures,
"message":"Adventures Found",
})
}else if op.RowsAffected == 0  {
json.NewEncoder(res).Encode(map[string]any{
"data":adeventures,
"message":"Adventures Not Found",
})
}else if op.Error != nil {
http.Error(res,"Something Went Wrong",500)
}
}



func PublishAdventures(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	dest_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(dest_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("adventures").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("adventures").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})

	}
}





func FetchAdventureCategorys(res http.ResponseWriter,req *http.Request){
	var adventureCategories []adventures.AdventuresCategory
	var individualCount int64
	var advCategory []any
	db.Db_Connection.Find(&adventureCategories)
	for _, v := range adventureCategories{
	db.Db_Connection.Table("adventures").Where("published=?",true).Where("adventures_category_id =?",v.ID).Count(&individualCount)	
	advCategory = append(advCategory, map[string]any{
	"data":v,
	"count":individualCount,
	})
	}
	json.NewEncoder(res).Encode(map[string]any{
	"data":advCategory,
	})
	}
	


	
func AdventureCategoryAdventures(res http.ResponseWriter,req *http.Request){
	name := mux.Vars(req)["name"]

	var adventureCategory adventures.AdventuresCategory
	var advs []adventures.Adventure
	db.Db_Connection.Where("adventure_slug = ?",name).Find(&adventureCategory)
	db.Db_Connection.Where("published=?",true).Where("adventures_category_id = ?",adventureCategory.ID).Find(&advs)
	json.NewEncoder(res).Encode(map[string]any{
	"data":map[string]any{
	"category":adventureCategory,
	"adventures":advs,
	},
	})
	}


	func FetchSingleAdventure(res http.ResponseWriter,req *http.Request){
		slug := mux.Vars(req)["slug"]
		var adventure adventures.Adventure
		var adventures []adventures.Adventure
		db.Db_Connection.Where("published=?",true).Where("slug=?",slug).Find(&adventure)
		db.Db_Connection.Where("published=?",true).Where("adventures_category_id=?",adventure.AdventuresCategoryID).Where("slug<>?",slug).Find(&adventures)
		
		json.NewEncoder(res).Encode(map[string]any{
		"data":adventure,
		"relatedAdventures":adventures,
		})
		}