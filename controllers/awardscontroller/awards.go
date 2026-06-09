package awardscontroller

import (
	"context"
	"encoding/json"
	"io"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/awards"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func CreateAward(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	cloud, _ := cld.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	ctx := context.TODO()
	time.Sleep(time.Nanosecond * 9)
	userid := <-helpers.UserIDChannel
	// parse like json.decode multipart form
	// then file upload i want to save them in the gorilla mux server
	req.ParseMultipartForm(20 << 30)

	awardThumb, handler, err := req.FormFile("awardPhoto")
	if err != nil {

		http.Error(res, err.Error(), 500)
	}
	defer awardThumb.Close()
	os.MkdirAll("files", os.ModePerm)
	localPath := filepath.Join("files", handler.Filename)
	fileAward, err1 := os.Create(localPath)
	if err1 != nil {
		http.Error(res, err1.Error(), 500)
		return
	}
	defer fileAward.Close()
	_, err2 := fileAward.ReadFrom(awardThumb)
	if err2 != nil {
		http.Error(res, err2.Error(), 500)
		return
	}

	result, err3 := cloud.Upload.Upload(ctx, localPath, uploader.UploadParams{
		Folder: "Awards",
	})

	if err3 != nil {
		panic(err3.Error())
	}
	excl, _ := json.Marshal(req.FormValue("excl"))
	incl, _ := json.Marshal(req.FormValue("incl"))
	slug := strings.ReplaceAll(req.FormValue("title"), " ", "-")
	pts, _ := strconv.Atoi(req.FormValue("points"))
	var award = awards.Award{
		AwardThumbnailPhoto: result.SecureURL,
		Title:               req.FormValue("title"),
		Exclusions:          excl,
		Slug:                slug,
		Inclusions:          incl,
		UserID:              userid,
		PointsNeeded:        pts,
		Published:           false,
		AwardDescription:    req.FormValue("awardDescription"),
	}
	op := db.Db_Connection.Create(&award)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Award Saved",
		})
		return
	} else if op.Error != nil {
		http.Error(res, op.Error.Error(), 500)
		return
	}
}

func UpdateAward(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	cloud, err := cld.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if err != nil {
		http.Error(res, "cloud init error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	userid := <-helpers.UserIDChannel

	if err := req.ParseMultipartForm(20 << 30); err != nil && err != http.ErrNotMultipart {
		http.Error(res, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get award ID from query
	id := req.URL.Query().Get("awardID")
	if id == "" {
		http.Error(res, "missing awardID", http.StatusBadRequest)
		return
	}

	// Fetch existing award so we can preserve fields (thumbnail) if no new file is provided
	var existing awards.Award
	if err := db.Db_Connection.First(&existing, id).Error; err != nil {
		http.Error(res, "award not found: "+err.Error(), http.StatusNotFound)
		return
	}

	var uploadResult *uploader.UploadResult // keep nil if no new upload

	// Try reading file; if missing, that's OK — we'll preserve existing thumbnail
	awardThumb, handler, err := req.FormFile("awardPhoto")
	if err == nil {
		// file was provided
		defer awardThumb.Close()

		// ensure files dir
		if err := os.MkdirAll("files", os.ModePerm); err != nil {
			http.Error(res, "failed to create files dir: "+err.Error(), http.StatusInternalServerError)
			return
		}

		localPath := filepath.Join("files", handler.Filename)
		fileAward, err := os.Create(localPath)
		if err != nil {
			http.Error(res, "failed to create local file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer fileAward.Close()

		// copy from uploaded file to local file
		if _, err := io.Copy(fileAward, awardThumb); err != nil {
			http.Error(res, "failed to save uploaded file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// upload to cloudinary
		result, err := cloud.Upload.Upload(ctx, localPath, uploader.UploadParams{
			Folder: "Awards",
		})
		if err != nil {
			http.Error(res, "cloud upload error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if result == nil {
			http.Error(res, "cloud upload returned nil result", http.StatusInternalServerError)
			return
		}

		uploadResult = result
	} else {
		// If FormFile returned an error other than missing file, handle it.
		if err != http.ErrMissingFile && err != http.ErrNotMultipart {
			http.Error(res, "error reading uploaded file: "+err.Error(), http.StatusBadRequest)
			return
		}
		// No file provided — we'll keep existing.AwardThumbnailPhoto
	}

	// Prepare fields (keep similar to your original logic)
	excl, _ := json.Marshal(req.FormValue("excl"))
	incl, _ := json.Marshal(req.FormValue("incl"))
	slug := strings.ReplaceAll(req.FormValue("title"), " ", "-")

	pts := 0
	if p := req.FormValue("points"); p != "" {
		if tmp, err := strconv.Atoi(p); err == nil {
			pts = tmp
		} else {
			// if conversion fails, return error (or keep 0 depending on desired behavior)
			http.Error(res, "invalid points value: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Build an updates map so we don't overwrite fields with zero values accidentally
	updates := map[string]interface{}{
		"title":         req.FormValue("title"),
		"exclusions":    excl,
		"inclusions":    incl,
		"slug":          slug,
		"user_id":       userid,
		"points_needed": pts,
	}

	if uploadResult != nil && uploadResult.SecureURL != "" {
		updates["award_thumbnail_photo"] = uploadResult.SecureURL
	} else {
		// preserve existing thumbnail (explicitly set it to existing value to be safe)
		updates["award_thumbnail_photo"] = existing.AwardThumbnailPhoto
	}

	op := db.Db_Connection.Model(&existing).Updates(updates)
	if op.Error != nil {
		http.Error(res, "db update error: "+op.Error.Error(), http.StatusInternalServerError)
		return
	}
	if op.RowsAffected == 0 {
		http.Error(res, "no rows updated", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"message": "Award Updated",
	})
}

func DeleteAward(res http.ResponseWriter, req *http.Request) {
	// how can i do req.params as i do in laravel /{someParam}
	var award awards.Award
	vars := mux.Vars(req)
	awardId := vars["id"]
	op := db.Db_Connection.Delete(&award, awardId)
	if op.RowsAffected > 0 && op.Error == nil {
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Deletion Done",
		})
		return
	} else {
		http.Error(res, op.Error.Error(), 500)
	}
}

func FetchAwards(res http.ResponseWriter, req *http.Request) {
	var awards []awards.Award
	err := db.Db_Connection.Where("published=?", true).Select([]string{"Slug", "Title", "Inclusions", "Exclusions", "AwardThumbnailPhoto", "AwardDescription"}).Find(&awards).Error
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	} else {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Awards Found",
			"data":    awards,
		})
	}
}

func FetchAwardsAdmin(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 4)
	user_id := <-helpers.UserIDChannel
	var awards []awards.Award
	err := db.Db_Connection.Where("user_id=?", user_id).Select([]string{"Slug", "Title", "Inclusions", "Exclusions", "AwardThumbnailPhoto", "AwardDescription"}).Find(&awards).Error
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	} else {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Awards Found",
			"data":    awards,
		})
	}
}

func FetchSlugwise(res http.ResponseWriter, req *http.Request) {
	slug := req.URL.Query().Get("slug")
	var award awards.Award
	op := db.Db_Connection.Where("slug=?", slug).First(&award)
	if op.RowsAffected > 0 {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Award found",
			"data":    award,
		})
	} else {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message": "Award Not found",
		})
	}

}
