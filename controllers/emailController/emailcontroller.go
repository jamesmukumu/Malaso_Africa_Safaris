package emailcontroller

import (
	"bytes"
	
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/bulkmails"
	"jamesmukumu/maasaimaratripsportal/models/email"
	"log"
	"strconv"
	"encoding/base64"

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	// "gopkg.in/gomail.v2"
	"github.com/resend/resend-go/v2"
)

type EmailData struct {
	Title string
	Body  template.HTML
}

func SaveEmailTemplate(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	var emailTemplate = email.EmailTemplates{
		Title:   req.FormValue("title"),
		Body:    req.FormValue("body"),
		UserID:  user_id,
		Subject: req.FormValue("subject"),
	}
	op := db.Db_Connection.Create(&emailTemplate)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Email Template Saved",
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Email Template Saved",
		})
	} else {
		http.Error(res, op.Error.Error(), 500)
	}
}


func EditEmailTemplate(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)

	userID := <-helpers.UserIDChannel
	id := req.URL.Query().Get("id")

	if id == "" {
		http.Error(res, "Missing template ID", http.StatusBadRequest)
		return
	}

	title := req.FormValue("title")
	body := req.FormValue("body")
	subject := req.FormValue("subject")

	// 1. Check if template exists
	var existingTemplate email.EmailTemplates
	if err := db.Db_Connection.First(&existingTemplate, "id = ?", id).Error; err != nil {
		http.Error(res, "Email template not found", http.StatusNotFound)
		return
	}
	err := db.Db_Connection.
		Model(&email.EmailTemplates{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"title":   title,
			"body":    body,
			"subject": subject,
			"user_id": userID,
		}).Error

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Success response
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"message": "Email Template Updated",
	})
}


func SaveEmailNewsletter(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 3)
	user_id := <-helpers.UserIDChannel
	bodyBytes, _ := json.Marshal(req.FormValue("body"))
	var emailTemplate = email.Newsletter{
		Title:   req.FormValue("title"),
		Body:    bodyBytes,
		UserID:  user_id,
		Subject: req.FormValue("subject"),
	}
	op := db.Db_Connection.Create(&emailTemplate)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Email Template Saved",
		})
	} else if op.RowsAffected == 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Email Template Saved",
		})
	} else {
		http.Error(res, op.Error.Error(), 500)
	}
}






func SendEmail(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()

	// Parse HTML template
	tmpl, err := template.ParseFiles("templates/mail.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	emailData := EmailData{
		Title: req.FormValue("title"),
		Body:  template.HTML(req.FormValue("body")),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, &emailData); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse multipart form (20MB max)
	if err := req.ParseMultipartForm(20 << 20); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read uploaded attachments into memory
	var attachments []*resend.Attachment

	files := req.MultipartForm.File["attachment"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := io.ReadAll(file)
		file.Close()

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		attachments = append(attachments, &resend.Attachment{
			Filename: fileHeader.Filename,
			Content:  []byte(base64.StdEncoding.EncodeToString(data)),
		})
	}

	// Create Resend client
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From: fmt.Sprintf("Malaso Africa Safaris <%s>", os.Getenv("MAILADDRESS")),
		To: []string{
			req.FormValue("recepient"),
		},
		Subject:     req.FormValue("subject"),
		Html:        body.String(),
		Attachments: attachments,
	}

	email, err := client.Emails.Send(params)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]any{
		"message":  "Email Delivered",
		"email_id": email.Id,
	})
}





func SaveBulkEmail(res http.ResponseWriter, req *http.Request) {
	var bulkEmail bulkmails.BulkMails
	err := json.NewDecoder(req.Body).Decode(&bulkEmail)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(res, err.Error(), 500)
		return
	}
	op := db.Db_Connection.Create(&bulkEmail)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Bulk Email Saved",
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}
}

func FetchEmailBulks(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)
	var bulkMails []bulkmails.BulkMails
	op := db.Db_Connection.Find(&bulkMails)
	if op.RowsAffected > 0 && op.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"data": bulkMails,
		})
	} else if op.Error == nil {
		http.Error(res, "Something went wrong", 500)
	}
}


func PropagateEmailsBulk(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()

	time.Sleep(time.Nanosecond * 3)

	user_id := <-helpers.UserIDChannel
	fmt.Println(user_id)

	var usersMap []map[string]any

	t := mux.Vars(req)["t"]

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	json.Unmarshal([]byte(req.FormValue("recepients")), &usersMap)

	if t == "email_templates" {
		var matchingEmail email.EmailTemplates
		db.Db_Connection.Where("id=?", id).Find(&matchingEmail)

		tmpl, err := template.ParseFiles("templates/mail.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		emailData := EmailData{
			Title: matchingEmail.Title,
			Body:  template.HTML(matchingEmail.Body),
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, &emailData); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create Resend client once
		client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

		for _, user := range usersMap {
			params := &resend.SendEmailRequest{
				From: fmt.Sprintf("PJ Safaris <%s>", os.Getenv("MAILADDRESS")),
				To: []string{
					user["email"].(string),
				},
				Subject: matchingEmail.Subject,
				Html:    body.String(),
			}

			_, err := client.Emails.Send(params)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				continue
			}
		}

		json.NewEncoder(res).Encode(map[string]string{
			"message": "Emails Sent",
		})
		return
	} else {

	}
}


func FetchTemplates(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 5)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var emails []email.EmailTemplates
var newsLetter []email.Newsletter
db.Db_Connection.Find(&emails)
db.Db_Connection.Find(&newsLetter)
json.NewEncoder(res).Encode(map[string]any{
"emails":emails,
"newsletters":newsLetter,
})
}


func SingleTemplate(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 3)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
id := req.URL.Query().Get("id")
idInt,_ := strconv.Atoi(id)
var email email.EmailTemplates
db.Db_Connection.Where("id =?",idInt).Find(&email)
json.NewEncoder(res).Encode(map[string]any{
"data":email,
})
}