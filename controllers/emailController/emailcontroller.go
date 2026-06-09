package emailcontroller

import (
	"bytes"
	"crypto/tls"
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

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
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

	// 2. Check if title already exists for another record
	// var count int64
	// db.Db_Connection.
	// 	Model(&email.EmailTemplates{}).
	// 	Where("title = ? AND id <> ?", title, id).
	// 	Count(&count)

	// if count > 0 {
	// 	http.Error(res, "Email template title already exists", http.StatusConflict)
	// 	return
	// }

	// 3. Update the existing template
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

	// Parse template
	tmpl, err := template.ParseFiles("templates/mail.html")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	// Prepare email data
	emailData := EmailData{
		Title: req.FormValue("title"),
		Body:  template.HTML(req.FormValue("body")),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, &emailData); err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	// Parse multipart form (IMPORTANT for multiple files)
	err = req.ParseMultipartForm(20 << 20) // 20MB limit
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	os.MkdirAll("attachments/files", os.ModePerm)

	// Slice to hold paths of saved files
	var attachmentPaths []string

	// Retrieve multiple files
	files := req.MultipartForm.File["attachment"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer file.Close()

		// Save file to disk
		savePath := "attachments/files/" + fileHeader.Filename

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		attachmentPaths = append(attachmentPaths, savePath)
	}

	// SEND EMAIL
	m := gomail.NewMessage()

	m.SetHeader("From", fmt.Sprintf("PJ Safaris <%s>", os.Getenv("MAILADDRESS")))
	m.SetHeader("To", req.FormValue("recepient"))
	m.SetHeader("Subject", req.FormValue("subject"))
	m.SetBody("text/html", body.String())

	// Attach all files
	for _, p := range attachmentPaths {
		m.Attach(p)
	}

	dialer := gomail.NewDialer("mail.privateemail.com", 465,
		os.Getenv("MAILADDRESS"),
		os.Getenv("MAILPASSWORD"),
	)
dialer.SSL = true

	if err := dialer.DialAndSend(m); err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Email Delivered",
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
			http.Error(res, err.Error(), 500)
			return
		}
		emailData := EmailData{
			Title: matchingEmail.Title,
			Body:  template.HTML(matchingEmail.Body),
		}
		var body bytes.Buffer
		if err := tmpl.Execute(&body, &emailData); err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		m := gomail.NewMessage()
		for _, user := range usersMap {
			m.SetHeader("From", fmt.Sprintf("PJ Safaris <%s>", os.Getenv("MAILADDRESS")))
			m.SetHeader("To", user["email"].(string))
			m.SetHeader("Subject", matchingEmail.Subject)
			m.SetBody("text/html", body.String())
			dialer := gomail.NewDialer("mail.privateemail.com", 587, os.Getenv("MAILADDRESS"), os.Getenv("MAILPASSWORD"))
			dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: "mail.privateemail.com"}
			if err := dialer.DialAndSend(m); err != nil {
				http.Error(res, err.Error(), 500)
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