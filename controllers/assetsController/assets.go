package assetscontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"log"

	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func FetchAssets(res http.ResponseWriter,req *http.Request) {
	godotenv.Load()
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	cloud, errCloudinary := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloudinary != nil {
		panic(errCloudinary.Error())
	}

folder := mux.Vars(req)["folder"]
result,err := cloud.Admin.Assets(context.Background(),admin.AssetsParams{
Prefix: folder,
MaxResults: 100,
DeliveryType: "upload",
})

if err != nil{
http.Error(res,err.Error(),500)
return
}
json.NewEncoder(res).Encode(map[string]any{
"data":result.Assets,
})
}

func SaveAssets(res http.ResponseWriter, req *http.Request) {
    time.Sleep(time.Nanosecond * 3)
    user_id := <-helpers.UserIDChannel
    fmt.Println(user_id)

   
    err := req.ParseMultipartForm(20 << 20) 
    if err != nil {
        http.Error(res, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
        return
    }

    folderName := mux.Vars(req)["folder"]

    cloud, errCloudinary := cloudinary.NewFromParams(
        os.Getenv("CloudinaryName"),
        os.Getenv("CloudinarySecret"),
        os.Getenv("CloudinaryPublic"),
    )
    if errCloudinary != nil {
        http.Error(res, "Cloudinary init error: "+errCloudinary.Error(), http.StatusInternalServerError)
        return
    }

    
    if req.MultipartForm == nil || req.MultipartForm.File == nil {
        http.Error(res, "No files uploaded", http.StatusBadRequest)
        return
    }

    files := req.MultipartForm.File["files"]

    if len(files) == 0 {
        http.Error(res, "No files received", http.StatusBadRequest)
        return
    }

    for _, f := range files {
        openedFile, err := f.Open()
        if err != nil {
            fmt.Println("Error opening file:", err)
            continue
        }

        _, uploadErr := cloud.Upload.Upload(
            context.Background(),
            openedFile,
            uploader.UploadParams{ Folder: folderName },
        )

        openedFile.Close()

        if uploadErr != nil {
            fmt.Println("Upload error:", uploadErr)
            continue
        }
    }

    json.NewEncoder(res).Encode(map[string]string{
        "message": "Assets Saved",
    })
}
// 


func DeleteAssets(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 3)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
cloud, errCloudinary := cloudinary.NewFromParams(
os.Getenv("CloudinaryName"),
os.Getenv("CloudinarySecret"),
os.Getenv("CloudinaryPublic"),
)
if errCloudinary != nil {
panic(errCloudinary.Error())
}

var assets []string
err1 := json.Unmarshal([]byte(req.FormValue("assets")),&assets)
if err1 != nil{
panic(err1.Error())
}

for _, i := range assets {
_,err := cloud.Admin.DeleteAssets(context.Background(),admin.DeleteAssetsParams{
DeliveryType: "upload",
PublicIDs: []string{i},
})
if err != nil {
log.Fatal("Asset Not Deleted")
continue
}        
}
json.NewEncoder(res).Encode(map[string]string{
"message":"Assets Deleted",
})
}