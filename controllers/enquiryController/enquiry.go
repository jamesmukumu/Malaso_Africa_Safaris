package enquirycontroller

import (
	"bytes"
	"log"

	// "crypto/tls"

	"encoding/json"
	"fmt"
	"html/template"
	"jamesmukumu/maasaimaratripsportal/db"
	enquiries "jamesmukumu/maasaimaratripsportal/models/Enquiries"

	// "log"
	"net/http"
	"os"
	"strconv"

	// "github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func EnquiryHandler(res http.ResponseWriter, req *http.Request) {
	FirstName := req.FormValue("firstName")
	LastName := req.FormValue("lastName")
	Email := req.FormValue("Email")
	Phonenumber := req.FormValue("Phonenumber")
	AdultsCount := req.FormValue("AdultsCount")
	KidsAges := req.FormValue("KidsAges")
	RoomsCount := req.FormValue("RoomsCount")
	var kidsAges int8
	var roomCount int8
	var adultsCount int8
	if FirstName == "" || LastName == "" || Email == "" || Phonenumber == "" {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Ensure All Required Fields Are Filled",
		})
		return
	}
	if AdultsCount != "" {
		adults, _ := strconv.Atoi(AdultsCount)
		adultsCount = int8(adults)
	} else if KidsAges != "" {
		kids, _ := strconv.Atoi(KidsAges)
		kidsAges = int8(kids)
	} else if RoomsCount != "" {
		rooms, _ := strconv.Atoi(RoomsCount)
		roomCount = int8(rooms)
	}

	var enquiry = enquiries.Enquiries{
		FirstName:          FirstName,
		LastName:           LastName,
		Email:              Email,
		PhoneNumber:        Phonenumber,
		ChildrenCount:      kidsAges,
		AdultsCount:        adultsCount,
		RoomsCount:         roomCount,
		EnquiryDescription: req.FormValue("EnquiryDescription"),
		StartStayDate:      req.FormValue("StartStayDate"),
		EndStayDate:        req.FormValue("EndStayDate"),
		Nationality:        req.FormValue("Nationality"),
	}
	op := db.Db_Connection.Create(&enquiry)
	if op.RowsAffected > 0 && op.Error == nil {
		if err := SendAcknowledgment(enquiry); err != nil {
			log.Println("email error:", err)
		}

		json.NewEncoder(res).Encode(map[string]any{
			"message": "Enquiry Saved",
			"ID":      enquiry.ID,
		})
		return
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Enquiry Not Saved",
		})
	}
}

func SendAcknowledgment(enquiry enquiries.Enquiries) error {
	// DO NOT load .env in production
	// godotenv.Load() // REMOVE THIS

	temp, err := template.ParseFiles("templates/enquiry.html")
	if err != nil {
		return err
	}

	var emailData bytes.Buffer
	if err := temp.Execute(&emailData, &enquiry); err != nil {
		return err
	}

	m := gomail.NewMessage()
	from := os.Getenv("MAILADDRESS")
	pass := os.Getenv("MAILPASSWORD")

	if from == "" || pass == "" {
		return fmt.Errorf("missing MAILADDRESS or MAILPASSWORD env vars")
	}

	m.SetHeader("From", fmt.Sprintf("PJ Safaris <%s>", from))
	m.SetHeader("To", enquiry.Email)
	m.SetHeader("Subject", "Enquiry Acknowledgment")
	m.SetBody("text/html", emailData.String())

	dialer := gomail.NewDialer("mail.privateemail.com", 587, from, pass)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func FetchEnquiries(res http.ResponseWriter, req *http.Request) {
	var enquiries []enquiries.Enquiries
	db.Db_Connection.Find(&enquiries)
	json.NewEncoder(res).Encode(map[string]any{
		"data": enquiries,
	})
}

func FetchEnquiry(res http.ResponseWriter, req *http.Request) {
	var enquiry enquiries.Enquiries
	ID := req.URL.Query().Get("ID")
	IDInt, _ := strconv.Atoi(ID)
	db.Db_Connection.Where("ID = ?", IDInt).Find(&enquiry)
	json.NewEncoder(res).Encode(map[string]any{
		"data": enquiry,
	})
}
