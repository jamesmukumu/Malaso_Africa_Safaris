package blogscontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/helpers/scopes"
	"jamesmukumu/maasaimaratripsportal/helpers/slug"
	blogs "jamesmukumu/maasaimaratripsportal/models/Blogs"
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
func CreateBlogCategory(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	var blogsCategory = blogs.BlogsCategory{
		Title:       req.FormValue("title"),
		Description: req.FormValue("description"),
		UserID:      user_id,
	}
	op := db.Db_Connection.Create(&blogsCategory)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Blogs Category Saved",
		})
	} else {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Blogs Not Saved",
		})
	}
}
func EditBlog(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Millisecond * 5)

	// Parse form
	err := req.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(res, "Invalid form data", 400)
		return
	}

	// Cloudinary init
	cloud, errCloud := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloud != nil {
		http.Error(res, "Cloudinary initialization failed", 500)
		return
	}

	ctx := context.Background()

	// --------------- FIND BLOG ----------------
	blogID := req.URL.Query().Get("id")
	var blog blogs.Blogs

	if err := db.Db_Connection.First(&blog, "id = ?", blogID).Error; err != nil {
		http.Error(res, "Blog not found", 404)
		return
	}

	// ---------------- UPDATE TEXT FIELDS -----------------

	if v := req.FormValue("title"); v != "" {
		blog.BlogTitle = v
		blog.Slug = slug.SlugMaker(req.FormValue("title"))
	}

	if v := req.FormValue("overview"); v != "" {
		blog.BlogOverview = v
	}

	if v := req.FormValue("description"); v != "" {
		blog.BlogDescription = v
	}

	if v := req.FormValue("blogCategoryID"); v != "" {
		id, _ := strconv.Atoi(v)
		blog.BlogsCategoryID = uint(id)
	}


	file, _, errPhoto := req.FormFile("blogPhoto")
	if errPhoto == nil {
	
		oldPublicID := extractPublicID(blog.BlogPhoto)
		if oldPublicID != "" {
			cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: oldPublicID})
		}

		// Upload new image
		up, errUp := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder: "blogs",
			UseFilename: &keepFilename,
		})
		if errUp == nil {
			blog.BlogPhoto = up.SecureURL
		}
	}

	if err := db.Db_Connection.Save(&blog).Error; err != nil {
		http.Error(res, "Failed to update blog", 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]any{
		"message": "Blog Updated Successfully",
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

func CreateBlog(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 4)
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

	blogsFile, _, _ := req.FormFile("blogPhoto")
	result, errBlog := cloud.Upload.Upload(ctx, blogsFile, uploader.UploadParams{
		Folder: "blogs",
		UseFilename: &keepFilename,
	})
	if errBlog != nil {
		log.Fatal(errBlog.Error())
		http.Error(res, "Something went wrong", 500)
		return
	}
	slug := slug.SlugMaker(req.FormValue("title"))
	blog_id, _ := strconv.Atoi(req.FormValue("blogCategoryID"))

	var blog = blogs.Blogs{
		BlogTitle:       req.FormValue("title"),
		BlogOverview:    req.FormValue("overview"),
		BlogDescription: req.FormValue("description"),
		BlogPhoto:       result.SecureURL,
		Slug:            slug,
		UserID:          user_id,
		BlogsCategoryID: uint(blog_id),
	}

	op := db.Db_Connection.Create(&blog)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Blog Added",
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Blog Not Added",
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}
}

func FetchMyBlogs(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 4)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var blogs []blogs.Blogs
	op := db.Db_Connection.Find(&blogs)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Blogs Fetched",
			"data":    blogs,
		})
	} else {
		json.NewEncoder(res).Encode(map[string]any{
			"message": "Blogs Empty",
		})
	}

}

func PublishBlog(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	dest_id := req.URL.Query().Get("id")
	id, _ := strconv.Atoi(dest_id)
	var args = mux.Vars(req)["method"]
	if args == "publish" {
		db.Db_Connection.Table("blogs").Where("user_id=?", user_id).Where("id=?", id).Update("published", true)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})
	} else if args == "unpublish" {
		db.Db_Connection.Table("blogs").Where("user_id=?", user_id).Where("id=?", id).Update("published", false)
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Update Success",
		})

	}
}



func FetchDisplayBlogs(res http.ResponseWriter,req *http.Request){
var blogs []blogs.Blogs
pageString := req.URL.Query().Get("pageString")
page := 1

if pageString != ""{
intPage,_:= strconv.Atoi(pageString)
 page = intPage
} 
baseQuery := db.Db_Connection.Where("published=?",true)
paginatedQuery := scopes.Paginator(baseQuery,page)
op := paginatedQuery.Find(&blogs)
if op.RowsAffected > 0 && op.Error == nil {
	var totalBlogsCount int64
	db.Db_Connection.Table("blogs").Where("published =?",true).Count(&totalBlogsCount)
	json.NewEncoder(res).Encode(map[string]any{
		"data":blogs,
		"count":totalBlogsCount,
		})
}
}

func FetchBlog(res http.ResponseWriter,req *http.Request){
var blog blogs.Blogs
var relatedBlogs []blogs.Blogs
slug := mux.Vars(req)["slug"]
db.Db_Connection.Where("slug=?",slug).Where("published=?",true).Find(&blog)
db.Db_Connection.Where("published = ? AND blogs_category_id = ? AND slug <> ?", true, blog.BlogsCategoryID, slug).Limit(4).Find(&relatedBlogs)
//
json.NewEncoder(res).Encode(map[string]any{
"data":blog,
"relatedBlogs":relatedBlogs,
})
}


func FetchSlugs(res http.ResponseWriter,req *http.Request){
var blogs []blogs.Blogs
db.Db_Connection.Select("slug").Find(&blogs)
json.NewEncoder(res).Encode(map[string]any{
"data":blogs,
})
}