package querycontroller

import (
	"encoding/json"
	"jamesmukumu/maasaimaratripsportal/helpers/scopes"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func QueryHandler(res http.ResponseWriter,req *http.Request) {
queryTerm := req.URL.Query().Get("term")
pageString := mux.Vars(req)["page"]
page,err := strconv.Atoi(pageString)
if err != nil {
log.Fatal(err.Error())
return
}
resp := scopes.QueryFinder(queryTerm,page)
json.NewEncoder(res).Encode(map[string]any{
"data":resp,
})
}