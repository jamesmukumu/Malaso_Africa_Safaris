package sitemapscontroller

import (
	"encoding/json"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/models/sitemaps"
	"net/http"
)

func SaveSiteMap(res http.ResponseWriter,req *http.Request) {
var siteMap = sitemaps.SiteMap{
Filename: req.FormValue("sitename"),
FilePath: req.FormValue("sitepath"),
}
siteMap.Validity()
db.Db_Connection.Create(&siteMap)
json.NewEncoder(res).Encode(map[string]any{
"message":"Sitemap saved",
})
}