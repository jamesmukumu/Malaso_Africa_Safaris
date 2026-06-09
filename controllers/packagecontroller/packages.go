package packagecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	blogs "jamesmukumu/maasaimaratripsportal/models/Blogs"
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	packagescategory "jamesmukumu/maasaimaratripsportal/models/PackagesCategory"
	"jamesmukumu/maasaimaratripsportal/models/adventures"
	"jamesmukumu/maasaimaratripsportal/models/attractions"
	customizedpackages "jamesmukumu/maasaimaratripsportal/models/customizedPackages"
	"jamesmukumu/maasaimaratripsportal/models/email"
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
	"github.com/joho/godotenv"
)
var keepFilename bool = true
func CreatePackageCategory(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)

	err := req.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(res, "Invalid form data", 400)
		return
	}

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
	file,_,_ := req.FormFile("categoryImage")
	cld,errUpload := cloud.Upload.Upload(ctx,file,uploader.UploadParams{
	Folder: "PackageCategories",
	})
	if errUpload != nil {
	log.Fatal(errUpload.Error())
	http.Error(res,errUpload.Error(),500)
	return
	}
	var package_category = packagescategory.PackageCategory{
		PackageCategoryName: req.FormValue("categoryName"),
		PackageDescription:  req.FormValue("categoryDescription"),
		UserID:              user_id,
		PackageCategoryImage: cld.SecureURL,
	}
	op := db.Db_Connection.Create(&package_category)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Category Saved",
		})
	} else if op.Error != nil {
		http.Error(res, op.Error.Error(), 500)
	}
}



func EditPackageCategory(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	userID := <-helpers.UserIDChannel
	// if err := req.ParseMultipartForm(50 << 20); err != nil {
	// 	http.Error(res, "Invalid form data", http.StatusBadRequest)
	// 	return
	// }
	categoryID := req.URL.Query().Get("ID")
	if categoryID == "" {
		http.Error(res, "Category ID is required", http.StatusBadRequest)
		return
	}


	var category packagescategory.PackageCategory
	if err := db.Db_Connection.First(&category, categoryID).Error; err != nil {
		http.Error(res, "Category not found", http.StatusNotFound)
		return
	}


	updates := map[string]interface{}{
		"user_id": userID,
	}


	if name := req.FormValue("categoryName"); name != "" {
		updates["package_category_name"] = name
	}


	if desc := req.FormValue("categoryDescription"); desc != "" {
		updates["package_description"] = desc
	}

	file, _, err := req.FormFile("categoryImage")
	if err == nil {
		defer file.Close()

		cloud, err := cloudinary.NewFromParams(
			os.Getenv("CloudinaryName"),
			os.Getenv("CloudinarySecret"),
			os.Getenv("CloudinaryPublic"),
		)
		if err != nil {
			http.Error(res, "Cloudinary init failed", http.StatusInternalServerError)
			return
		}
      
		uploadRes, err := cloud.Upload.Upload(
			context.Background(),
			file,
			uploader.UploadParams{
				Folder: "PackageCategories",
				UseFilename: &keepFilename,
				
			},
		)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		updates["package_category_image"] = uploadRes.SecureURL
	}

	if len(updates) <= 1 { 
		http.Error(res, "No fields provided to update", http.StatusBadRequest)
		return
	}

	if err := db.Db_Connection.Model(&category).Updates(updates).Error; err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Category updated successfully",
	})
}

func FetchPackageCategories(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 5)
user_id :=<-helpers.UserIDChannel
fmt.Println(user_id) 
var packageCategories []packagescategory.PackageCategory
db.Db_Connection.Find(&packageCategories)
json.NewEncoder(res).Encode(map[string]any{
"data":packageCategories,
})
}

func FetchSingularPackageCategory(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 3)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var packageCategory packagescategory.PackageCategory
idString := mux.Vars(req)["id"]
ID,_ := strconv.Atoi(idString)
db.Db_Connection.Where("ID =?",ID).Find(&packageCategory)
json.NewEncoder(res).Encode(map[string]any{
"data":packageCategory,
})
}


func EditPackage(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Millisecond * 5)

	err := req.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(res, "Invalid form data", 400)
		return
	}

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

	// ---------------- FIND PACKAGE ----------------
	packageID := req.URL.Query().Get("id")
	var pkg packages.Package

	if err := db.Db_Connection.First(&pkg, "id = ?", packageID).Error; err != nil {
		http.Error(res, "Package not found", 404)
		return
	}

	// ---------------- UPDATE STRING FIELDS ----------------

	if v := req.FormValue("packageTitle"); v != "" {
		pkg.PackageTitle = v
		pkg.PackageSlug = slug.SlugMaker(v)
	}

	if v := req.FormValue("packageOverview"); v != "" {
		pkg.PackageOverview = v
	}

	if v := req.FormValue("specialNotes"); v != "" {
		pkg.SpecialNotes = v
	}

	if v := req.FormValue("meanTransport"); v != "" {
		pkg.MeansOfTransport = v
	}

	if v := req.FormValue("packageChargeCurrency"); v != "" {
		pkg.PackageChargeCurrency = v
	}

	if v := req.FormValue("startValidity"); v != "" {
		pkg.StartValidity = v
	}

	if v := req.FormValue("endValidity"); v != "" {
		pkg.EndValidity = v
	}

	// ---------------- COMPLEX FIELD UPDATES -----------------

	if v := req.FormValue("inclusions"); v != "" {
		b, _ := json.Marshal(v)
		pkg.Inclusions = b
	}

	if v := req.FormValue("exclusions"); v != "" {
		b, _ := json.Marshal(v)
		pkg.Exclusions = b
	}

	if v := req.FormValue("itenerary"); v != "" {
		b, _ := json.Marshal(v)
		pkg.PackageItenerary = b
	}

	if v := req.FormValue("packageChargeAmount"); v != "" {
		n, _ := strconv.Atoi(v)
		pkg.PackageCharge = n
	}

	if v := req.FormValue("packageCategoryID"); v != "" {
		n, _ := strconv.Atoi(v)
		pkg.PackageCategoryID = uint(n)
	}

	if v := req.FormValue("destinationID"); v != "" {
		n, _ := strconv.Atoi(v)
		pkg.DestinationID = uint(n)
	}

	// =========================================================
	// UPDATE MAIN IMAGE: packageImage
	// =========================================================

	file, _, errImage := req.FormFile("packageImage")
	if errImage == nil {

		// Delete old image
		oldID := extractPublicID(pkg.PackagePhoto)
		if oldID != "" {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: oldID})
		}

		// Upload new one
		up, errUp := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder: "packages",
		})
		if errUp == nil {
			pkg.PackagePhoto = up.SecureURL
		}
	}

	// =========================================================
	// UPDATE MULTIPLE IMAGES: packageImages
	// =========================================================

	files := req.MultipartForm.File["packageImages"]
	if len(files) > 0 {

		// Delete old photos
		var oldPhotos []string
		json.Unmarshal(pkg.PackagePhotos, &oldPhotos)

		for _, url := range oldPhotos {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: extractPublicID(url),
			})
		}

		// Upload new photos
		var newPhotos []string

		for _, hdr := range files {
			f, err := hdr.Open()
			if err != nil {
				continue
			}
			defer f.Close()

			up, err2 := cloud.Upload.Upload(ctx, f, uploader.UploadParams{
				Folder: "packages",
			})
			if err2 == nil {
				newPhotos = append(newPhotos, up.SecureURL)
			}
		}

		b, _ := json.Marshal(newPhotos)
		pkg.PackagePhotos = b
	}

	// ---------------- SAVE UPDATES ----------------
	if err := db.Db_Connection.Save(&pkg).Error; err != nil {
		http.Error(res, "Package update failed", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Package Updated Successfully",
	})
}
func extractPublicID(url string) string {
	if url == "" {
		return ""
	}
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, filepath.Ext(last))
}

func CreatePackage(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	err := req.ParseMultipartForm(30 << 45)
	if err != nil {
		log.Fatal(err.Error())
		json.NewEncoder(res).Encode(map[string]string{
		"message":"File Sizes exceed limit",
		})
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
	packageImage, _, _ := req.FormFile("packageImage")
	uploadResult, err1 := cloud.Upload.Upload(ctx, packageImage, uploader.UploadParams{
	Folder: "packages",
	UseFilename: &keepFilename,
	})
	if err1 != nil {
		log.Fatal(err1.Error())
		http.Error(res, err1.Error(), 500)
		return
	}
	packageImages := req.MultipartForm.File["packageImages"]
	var photosPackage []string
	for _, head := range packageImages {
		file, errOpenFile := head.Open()
		defer file.Close()
		if errOpenFile != nil {
			http.Error(res, errOpenFile.Error(), 500)
			continue
		}
		packageImagesResult, _ := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "packages",
		UseFilename: &keepFilename,
		})
		photosPackage = append(photosPackage, packageImagesResult.SecureURL)
	}
	slug := slug.SlugMaker(req.FormValue("packageTitle"))
	incl, _ := json.Marshal(req.FormValue("inclusions"))
	excl, _ := json.Marshal(req.FormValue("exclusions"))
	itenerary, _ := json.Marshal(req.FormValue("itenerary"))
	photosBytes, _ := json.Marshal(photosPackage)
	packageCategory, _ := strconv.Atoi(req.FormValue("packageCategoryID"))
	destinationID, _ := strconv.Atoi(req.FormValue("destinationID"))
	chargeAmount, _ := strconv.Atoi(req.FormValue("packageChargeAmount"))
	var packageSave = packages.Package{
		PackageTitle:          req.FormValue("packageTitle"),
		PackageItenerary:      itenerary,
		Inclusions:            incl,
		Exclusions:            excl,
		PackageSlug:           slug,
		PackageCharge:         chargeAmount,
		PackageChargeCurrency: req.FormValue("packageChargeCurrency"),
		PackageOverview:       req.FormValue("packageOverview"),
		PackageCategoryID:     uint(packageCategory),
		DestinationID:         uint(destinationID),
		SpecialNotes:          req.FormValue("specialNotes"),
		PackagePhotos:         photosBytes,
		PackagePhoto:          uploadResult.SecureURL,
		UserID:                user_id,
		MeansOfTransport:      req.FormValue("meanTransport"),
		StartValidity:         req.FormValue("startValidity"),
		EndValidity:           req.FormValue("endValidity"),
	}
	op := db.Db_Connection.Create(&packageSave)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Package Saved",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Package Not Saved",
		})
	}
}



func Prepackage(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var destination []destinations.Destination
	var package_category []packagescategory.PackageCategory
	var hotels []hotels.Hotels
	var blogs_category []blogs.BlogsCategory
    var emailTemplates []email.EmailTemplates
    var newsLetters []email.Newsletter
    var adventure_category []adventures.AdventuresCategory 

	db.Db_Connection.Find(&destination)
	db.Db_Connection.Find(&package_category)
	db.Db_Connection.Find(&hotels)
	db.Db_Connection.Find(&blogs_category)
   db.Db_Connection.Find(&emailTemplates)
    db.Db_Connection.Find(&newsLetters)
	db.Db_Connection.Find(&adventure_category)
	var destinations_final []interface{}
	var pcks_final []any
	var hotels_final []any
	var blogs_final []any
    var email_final []any
    var newsletters_final []any
    var adv_final []any
	for _, d := range destination {
		var destinationMap = make(map[string]interface{}, 0)
		destinationMap["label"] = d.DestinationTitle
		destinationMap["value"] = d.ID
		destinations_final = append(destinations_final, destinationMap)
	}
	for _, h := range hotels {
		var hotelsMap = make(map[string]interface{}, 0)
		hotelsMap["label"] = h.HotelName
		hotelsMap["value"] = h.ID
		hotels_final = append(hotels_final, hotelsMap)
	}
	for _, p := range package_category {
		var packageMap = make(map[string]interface{}, 0)
		packageMap["label"] = p.PackageCategoryName
		packageMap["value"] = p.ID
		pcks_final = append(pcks_final, packageMap)
	}
	for _, b := range blogs_category {
		var blogsMap = make(map[string]interface{}, 0)
		blogsMap["label"] = b.Title
		blogsMap["value"] = b.ID
		blogs_final = append(blogs_final, blogsMap)
	}

	for _, e := range emailTemplates {
		var eMap = make(map[string]interface{}, 0)
		eMap["label"] =e.Title
		eMap["value"] = e.ID
		email_final = append(email_final, eMap)
	}
	
	for _, n := range newsLetters {
		var nMap = make(map[string]interface{}, 0)
		nMap["label"] = n.Title
		nMap["value"] = n.ID
	 newsletters_final= append(email_final, nMap)
	}
	for _, a := range adventure_category{
		var aMap = make(map[string]interface{}, 0)
		aMap["label"] = a.AdventureCategoryName
		aMap["value"] = a.ID
	 adv_final= append(adv_final, aMap)
	}
	json.NewEncoder(res).Encode(map[string]interface{}{
		"pcks":   pcks_final,
		"dest":   destinations_final,
		"hotels": hotels_final,
		"blogs":  blogs_final,
		"email":email_final,
		"newsletters":newsletters_final,
		"adventures":adv_final,
	})
}

func FetchPublishedPackages(res http.ResponseWriter, req *http.Request) {
	var packages []packages.Package
	op := db.Db_Connection.Where("published=?", true).Preload("User").Limit(12).Find(&packages)
	if op.RowsAffected > 0 && op.Error == nil {
		res.Header().Add("Content-Type","application/json")
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Packages fetched",
			"data":    packages,
		})
	}



}
func FilterPackages(res http.ResponseWriter, req *http.Request) {
	// filter can happen by price, currency, and destination
	var resultPackages []packages.Package

	// Read query params
	startRange := req.URL.Query().Get("startRange")
	endRange := req.URL.Query().Get("endRange")
	currency := req.URL.Query().Get("currency")
	destination := req.URL.Query().Get("destination") // optional extra if needed

	// Start building query
	query := db.Db_Connection.Model(&packages.Package{})

	// Apply filters only if params exist
	if startRange != "" {
		if startRangeInt, err := strconv.Atoi(startRange); err == nil {
			query = query.Where("package_charge >= ?", startRangeInt)
		}
	}

	if endRange != "" {
		if endRangeInt, err := strconv.Atoi(endRange); err == nil {
			query = query.Where("package_charge <= ?", endRangeInt)
		}
	}

	if currency != "" {
		query = query.Where("package_charge_currency = ?", currency)
	}

	if destination != "" {
		query = query.Where("destination_id = ?", destination)
	}

	// Execute query
	if err := query.Find(&resultPackages).Error; err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resultPackages) == 0 {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "No matching Results",
		})
	} else {
		json.NewEncoder(res).Encode(resultPackages)
	}

}

func FetchOrgsPackages(res http.ResponseWriter, req *http.Request) {
	var pcks []packages.Package
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	op := db.Db_Connection.Find(&pcks)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"data": pcks,
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "No Active package",
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}
}

func PublishPackage(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)

	user_id := <-helpers.UserIDChannel
	package_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(package_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("packages").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("packages").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	}
}

func DeletePackage(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	var pck packages.Package
	package_id := req.URL.Query().Get("id")
	packageID, _ := strconv.Atoi(package_id)
	op := db.Db_Connection.Table("packages").Where("user_id=?", user_id).Where("id=?", packageID).Delete(&pck)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Deletion Success",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Something Went wrong",
		})
	}

}

func FetchPackageSlugs(res http.ResponseWriter,req *http.Request){
var packages []packages.Package
db.Db_Connection.Select("package_slug").Find(&packages)
json.NewEncoder(res).Encode(map[string]any{
"data":packages,
})

}


func FetchPackageCategorys(res http.ResponseWriter,req *http.Request){
var packageCategorys []packagescategory.PackageCategory
var individualCount int64
var pcksCategory []any
db.Db_Connection.Find(&packageCategorys)
for _, v := range packageCategorys{
db.Db_Connection.Table("packages").Where("published=?",true).Where("package_category_id =?",v.ID).Count(&individualCount)	
pcksCategory = append(pcksCategory, map[string]any{
"data":v,
"count":individualCount,
})
}
json.NewEncoder(res).Encode(map[string]any{
"data":pcksCategory,
})
}

func PackageCategoryPackages(res http.ResponseWriter,req *http.Request){
name := mux.Vars(req)["name"]
newName := strings.ReplaceAll(name,"-"," ")
var packageCategory packagescategory.PackageCategory
var packages []packages.Package
db.Db_Connection.Where("package_category_name = ?",newName).Find(&packageCategory)
db.Db_Connection.Where("published=?",true).Where("package_category_id = ?",packageCategory.ID).Find(&packages)
json.NewEncoder(res).Encode(map[string]any{
"data":map[string]any{
"category":packageCategory,
"packages":packages,
},
})
}



func FetchSinglePackage(res http.ResponseWriter,req *http.Request){
slug := mux.Vars(req)["slug"]
var pck packages.Package
var relatedPackages []packages.Package
db.Db_Connection.Where("published=?",true).Where("package_slug=?",slug).Find(&pck)
db.Db_Connection.Where("published=?",true).Where("package_category_id=?",pck.PackageCategoryID).Where("package_slug<>?",slug).Find(&relatedPackages)

json.NewEncoder(res).Encode(map[string]any{
"data":pck,
"relatedPackages":relatedPackages,
})
}


func PackageSlugs(res http.ResponseWriter,req *http.Request){
var packageSlugs []packages.Package
var packageCategorySlugs []packagescategory.PackageCategory
var Adventures []adventures.Adventure
var AdventuresCqtegories []adventures.AdventuresCategory

db.Db_Connection.Find(&packageSlugs)
var packageSlugsString []string
var categoriesSlug []string
for _,pck := range packageSlugs{
packageSlugsString = append(packageSlugsString, pck.PackageSlug)
}
db.Db_Connection.Find(&packageCategorySlugs)
for _,category := range packageCategorySlugs{
categoriesSlug = append(categoriesSlug, strings.ReplaceAll(category.PackageCategoryName," ","-"))
}
var adv []string
db.Db_Connection.Find(&Adventures)

for _,Adventure := range Adventures {
adv = append(adv, Adventure.Slug)
}

var advCategories []string
db.Db_Connection.Find(&AdventuresCqtegories)

for _,AdventureCat := range AdventuresCqtegories {
advCategories = append(advCategories,AdventureCat.AdventureSlug)
}

var attr []string
var attractions []attractions.Attraction
db.Db_Connection.Find(&attractions)
for _,attraction := range attractions{
attr = append(attr, attraction.AttractionSlug)
}
json.NewEncoder(res).Encode(map[string]any{
"packageSlugs":packageSlugsString,
"categoriesSlug":categoriesSlug,
"adventures":adv,
"adventureCategories":advCategories,
"attractions":attr,
})
}



func CreateCustomPackage(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	err := req.ParseMultipartForm(30 << 45)
	if err != nil {
		log.Fatal(err.Error())
		json.NewEncoder(res).Encode(map[string]string{
		"message":"File Sizes exceed limit",
		})
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
	packageImage, _, _ := req.FormFile("packageImage")
	uploadResult, err1 := cloud.Upload.Upload(ctx, packageImage, uploader.UploadParams{
	Folder: "packages",
	UseFilename: &keepFilename,
	})
	if err1 != nil {
		log.Fatal(err1.Error())
		http.Error(res, err1.Error(), 500)
		return
	}
	packageImages := req.MultipartForm.File["packageImages"]
	var photosPackage []string
	for _, head := range packageImages {
		file, errOpenFile := head.Open()
		defer file.Close()
		if errOpenFile != nil {
			http.Error(res, errOpenFile.Error(), 500)
			continue
		}
		packageImagesResult, _ := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "packages",
		UseFilename: &keepFilename,
		})
		photosPackage = append(photosPackage, packageImagesResult.SecureURL)
	}
	slug := strings.ReplaceAll(req.FormValue("packageTitle"), " ", "-")
	incl, _ := json.Marshal(req.FormValue("inclusions"))
	excl, _ := json.Marshal(req.FormValue("exclusions"))
	itenerary, _ := json.Marshal(req.FormValue("itenerary"))
	photosBytes, _ := json.Marshal(photosPackage)
	packageCategory, _ := strconv.Atoi(req.FormValue("packageCategoryID"))
	
	chargeAmount, _ := strconv.Atoi(req.FormValue("packageChargeAmount"))
	var client = customizedpackages.Client{
		FirstName: req.FormValue("firstName"),
		LastName: req.FormValue("lastName"),
		Email: req.FormValue("email"),
		PhoneNumber: req.FormValue("phoneNumber"),
	}
	clientBytes,_ := json.Marshal(client)
	var packageSave = customizedpackages.CustomizedPackage{
		PackageTitle:          req.FormValue("packageTitle"),
		PackageItenerary:      itenerary,
		Inclusions:            incl,
		Exclusions:            excl,
		PackageSlug:           slug,
		PackageCharge:         chargeAmount,
		PackageChargeCurrency: req.FormValue("packageChargeCurrency"),
		PackageOverview:       req.FormValue("packageOverview"),
		PackageCategoryID:     uint(packageCategory),
		
		SpecialNotes:          req.FormValue("specialNotes"),
		PackagePhotos:         photosBytes,
		PackagePhoto:          uploadResult.SecureURL,
		UserID:                user_id,
		MeansOfTransport:      req.FormValue("meanTransport"),
		ClientsDetail: clientBytes,
	}
	packageSave.ExpirySetter()
	op := db.Db_Connection.Create(&packageSave)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Package Saved",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Package Not Saved",
		})
	}
}



func NavStats(res http.ResponseWriter,req *http.Request){
var destinations int64
var hotels int64
var pcks int64
var blogs int64
db.Db_Connection.Table("destinations").Count(&destinations)
db.Db_Connection.Table("blogs").Count(&blogs)
db.Db_Connection.Table("hotels").Count(&hotels)
db.Db_Connection.Table("packages").Count(&pcks)
var recentsPackages []packages.Package
db.Db_Connection.Table("packages").Order("created_at ASC").Limit(6).Find(&recentsPackages)
json.NewEncoder(res).Encode(map[string]any{
"overview":map[string]any{
"destinations":destinations,
"hotels":hotels,
"packages":pcks,
"blogs":blogs,
},
"packages":recentsPackages,
})
}