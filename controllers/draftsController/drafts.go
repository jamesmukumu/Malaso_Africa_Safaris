package draftscontroller

import (
	
	"encoding/json"

	"jamesmukumu/maasaimaratripsportal/db"
	
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/drafts"
	
	"net/http"
	"time"
)

func SaveDraft(res http.ResponseWriter,req *http.Request) {
time.Sleep(time.Nanosecond * 5)
user_id := <-helpers.UserIDChannel
var existingDraft drafts.Draft
op1 := db.Db_Connection.Where("draft_table =?",req.FormValue("draftName")).Find(&existingDraft)
if op1.RowsAffected > 0 {
db.Db_Connection.Where("draft_table =?",req.FormValue("draftName")).Delete(&existingDraft)
}
var drafts = drafts.Draft{
UserID: user_id,
DraftTable: req.FormValue("draftName"),
DraftContent: req.FormValue("draftValue"),
}
drafts.DraftExpiry()
op := db.Db_Connection.Create(&drafts)
if op.RowsAffected > 0 && op.Error == nil { 
json.NewEncoder(res).Encode(map[string]string{
"message":"Draft Added",
})
}else if op.Error != nil {
http.Error(res,op.Error.Error(),500)
}
}





// func DeleteOutDated(res http.ResponseWriter, req *http.Request){
// today := time.Now()

// //how can i query drafts where DraftCloses is today and perform a soft delete
// }
