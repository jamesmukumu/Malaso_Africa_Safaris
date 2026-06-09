package usercontrollers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	template "html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	jwt "github.com/golang-jwt/jwt/v5"

	// "fmt"
	"jamesmukumu/maasaimaratripsportal/controllers/walletscontroller"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	"jamesmukumu/maasaimaratripsportal/helpers"

	"jamesmukumu/maasaimaratripsportal/models/logins"
	"jamesmukumu/maasaimaratripsportal/models/users"
	"jamesmukumu/maasaimaratripsportal/models/wallets"
	"net/http"

	"github.com/google/uuid"

	"github.com/joho/godotenv"
	env "github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	mail "gopkg.in/gomail.v2"

	bcrypt "golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	Password   string `json:"password"`
	Credential string `json:"credential"`
}

var someEmailTemplate string = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"  
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" xmlns="http://www.w3.org/1999/xhtml" lang="en">
<head>
  <meta charset="UTF-8">
  <meta content="width=device-width, initial-scale=1" name="viewport">
  <meta name="x-apple-disable-message-reformatting">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta content="telephone=no" name="format-detection">
  <title>Email Send</title>
  <link href="https://fonts.googleapis.com/css?family=Oswald:300,700&display=swap" rel="stylesheet">
  <style type="text/css">
    body {width:100%;height:100%;margin:0;padding:0;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%;}
    h1, p {margin:0; padding:0;}
    .logo-container {text-align:center; width:100%;}
    .email-body {text-align:left; padding:20px; font-family:'Open Sans', sans-serif; line-height:1.6; font-size:16px; color:#262626; max-width:600px; margin:0 auto;}
    .dear-section {padding:20px 0; line-height:1.8; font-size:16px;}
    .contact-section {margin-top:40px; padding:20px; border-radius:8px;}
    .contact-item {margin-bottom:10px; font-size:16px; line-height:1.4;}
    .contact-label {font-weight:bold; color:#303233;}
    @media only screen and (max-width:600px) {
      .email-body {padding:10px;}
      h1 {font-size:24px!important;}
      p {font-size:16px!important;}
    }
  </style>
</head>
<body style="background-color:#F5F5F5; margin:0; padding:0;">
  <table style="margin:auto; max-width:600px; padding:10px;" cellpadding="0" cellspacing="0" style="background-color:#F5F5F5; width:100%;">
    <tr>
      <td align="center">

        <!-- Header (Full Width Background) -->
        <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#303233;">
          <tr>
            <td style="padding:25px 10px 40px 10px; text-align:center;">
              <img src="https://res.cloudinary.com/dasrniwpk/image/upload/v1746617656/logoMaasai_vw03id_r3n9yf.png" 
                   alt="Maasai Mara Trips Logo" style="height:60px;width:150px;display:block;margin:0 auto;">
            </td>
          </tr>
        </table>

        <!-- Body (Centered 600px) -->
        <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#ffffff; margin-top:0;">
          <tr>
            <td>
              <div class="email-body">
                <div class="dear-section">
                  <p>Dear {{ .name }} ,</p>
                  <p>Thank you for reaching out to Masai Mara Trips. We’re excited to help you plan a memorable safari experience.</p>
                  <p>We’ve received your request and one of our dedicated travel advisors will review it right away. You can expect a personalized response within 8 hours, complete with:</p>
                  <ul>
                  <li>A curated itinerary designed around your travel dates</li>
                  <li>Expert recommendations on accommodations, game drives, and activities</li>
                  <li>Transparent pricing with no hidden costs</li>
                  </ul>
                  <!-- Enquiry Details Section -->
                    <table width="100%" cellpadding="0" cellspacing="0" style="font-family: Arial, sans-serif; border-collapse: collapse; margin-top: 20px;">
                    
                      <!-- Section Title -->
                      <tr>
                        <td colspan="2" style="background-color: #F4E9D8; color: #4A2E19; font-weight: bold; padding: 12px; font-size: 16px; border: 1px solid #E5D2B8;">
                          Your Enquiry Details
                        </td>
                      </tr>
                    
                      <!-- Name Row -->
                      <tr>
                        <td style="width: 30%; padding: 10px; border: 1px solid #E5D2B8; font-weight: bold; color: #4A2E19;">Name</td>
                        <td style="width: 70%; padding: 10px; border: 1px solid #E5D2B8;">{{ .name }}</td>
                      </tr>
                    
                      <!-- Email Row -->
                      <tr>
                        <td style="width: 30%; padding: 10px; border: 1px solid #E5D2B8; font-weight: bold; color: #4A2E19;">Email</td>
                        <td style="width: 70%; padding: 10px; border: 1px solid #E5D2B8;">jj</td>
                      </tr>
                    
                      <!-- Phone Row -->
                      <tr>
                        <td style="width: 30%; padding: 10px; border: 1px solid #E5D2B8; font-weight: bold; color: #4A2E19;">Phone</td>
                        <td style="width: 70%; padding: 10px; border: 1px solid #E5D2B8;">dd</td>
                      </tr>
                    
                      <!-- Preferred Contact Row -->
                      <tr>
                        <td style="width: 30%; padding: 10px; border: 1px solid #E5D2B8; font-weight: bold; color: #4A2E19;">Preferred Contact</td>
                        <td style="width: 70%; padding: 10px; border: 1px solid #E5D2B8;">jj</td>
                      </tr>
                    
                      <!-- Message Row -->
                      <tr>
                        <td style="width: 30%; padding: 10px; border: 1px solid #E5D2B8; font-weight: bold; color: #4A2E19;">Message</td>
                        <td style="width: 70%; padding: 10px; border: 1px solid #E5D2B8;">kk</td>
                      </tr>
                    
                    </table>
                  <p>If you have any additional details to share – such as preferred travel dates, group size, or special interests – simply reply to this email and we’ll include them in your itinerary.
                  </p>
                  <p>
                    We look forward to supporting your trip.
                    </p>
                    <br/>
                    <p>
                        Warm Regards,
                    <br/>Masai Mara Trips Team
                  <br><span class="contact-label">Phone:</span> +254 705 769 896
                  <br><span class="contact-label">Email:</span> info@masaimaratrips.com
                  </p>
                </div>

                  
               
              </div>
            </td>
          </tr>
        </table>

        <!-- Footer (Full Width Background) -->
        <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#303233; margin-top:0;">
          <tr>
            <td style="padding:20px; text-align:center; color:#ffffff; font-family:'Open Sans', sans-serif; font-size:14px;">
              &copy; 2025 Maasai Mara Trips. All rights reserved.
            </td>
          </tr>
        </table>

      </td>
    </tr>
  </table>
</body>
</html>

`
var responseMap = make(map[string]interface{}, 0)

// register user
func CreateUser(res http.ResponseWriter, req *http.Request) {
	// time.Sleep(time.Nanosecond * 3)
	// user_id := <-helpers.UserIDChannel
    // fmt.Println(user_id)
	godotenv.Load()


	cloud, errCloud := cloudinary.NewFromParams(
		os.Getenv("CloudinaryName"),
		os.Getenv("CloudinarySecret"),
		os.Getenv("CloudinaryPublic"),
	)
	if errCloud != nil {
		log.Fatal(errCloud.Error())
		return
	}
	ctx := context.TODO()
	file, _, _ := req.FormFile("orgPhoto")

	result, errUpload := cloud.Upload.Upload(ctx, file, uploader.UploadParams{})
	if errUpload != nil {
		panic(errUpload.Error())
	}

	hashedPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), 13)
	if errBcrypt != nil {
		panic(errBcrypt.Error())
	}
	var real_user = users.User{
		Name:             req.FormValue("name"),
		Email:            req.FormValue("email"),
		PhoneNumber:      req.FormValue("phoneNumber"),
		Password:         string(hashedPassword),
		Role:             "normal user",
		Org:              true,
		Verified:         false,
		OrgBanner:        result.SecureURL,
		OrgPrimaryColors: req.FormValue("orgColor"),
	}
	op := db.Db_Connection.Create(&real_user)

	if op.RowsAffected > 0 && op.Error == nil {
		responseMap["message"] = "User Created Successfully"
		responseMap["status"] = 1
		responseBody, _ := json.Marshal(responseMap)
		res.Write(responseBody)
		sendAlertEmail(real_user.Email, real_user.Name)
		var wallet_label, _ = walletscontroller.GenerateSecureRandomString(20)
		var walletBody = wallets.Wallet{
			WalletLabel:   wallet_label,
			WalletBalance: 0,
			UserID:        real_user.ID,
		}
		db.Db_Connection.Create(&walletBody)
	} else if strings.Contains(op.Error.Error(), "ERROR: duplicate key value violates unique constraint \"uni_users_email\" (SQLSTATE 23505)") {
		responseMap["message"] = "Email Or Phone Number already Exists"
		responseMap["status"] = -1
		responseBody, _ := json.Marshal(responseMap)
		res.Write(responseBody)
	} else {
		responseMap["message"] = "Something Went Wrong"
		responseBody, _ := json.Marshal(responseMap)
		res.WriteHeader(500)
		res.Write(responseBody)
	}

}


func sendAlertEmail(emailAddress string, name string) {
	env.Load()
	m := mail.NewMessage()

	temp, _ := template.New("welcomeMail").Parse(someEmailTemplate)
	var body bytes.Buffer
	var nameMap = map[string]string{
		"name": name,
	}
	err_send := temp.Execute(&body, nameMap)
	if err_send != nil {
		fmt.Println("Error exucting data")
	}

	m.SetHeader("From", os.Getenv("MAILADDRESS"))
	m.SetHeader("To", emailAddress)
	m.SetHeader("Subject", "Welcome")
	m.SetBody("text/html", body.String())
	dialer := gomail.NewDialer("mail.privateemail.com", 465, os.Getenv("MAILADDRESS"), os.Getenv("MAILPASSWORD"))
	dialer.DialAndSend(m)
}

func Login(res http.ResponseWriter, req *http.Request) {
	env.Load()

	var user users.User
	var response = make(map[string]string, 0)
	var credent LoginCredentials
	err := json.NewDecoder(req.Body).Decode(&credent)
	if err != nil {
		panic(err.Error())
	}
	op := db.Db_Connection.Table("users").Where("Email=? OR Name =?", credent.Credential, credent.Credential).Find(&user)
	if op.RowsAffected == 0 {
		response["message"] = "Account Non existent"
		databytes, _ := json.Marshal(response)
		res.Header().Add("Content-Type","application/json")
		res.Write(databytes)
	} else {
		errHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credent.Password))
		if errHash != nil {
			response["message"] = "Invalid Credentials"
			databytes, _ := json.Marshal(response)
				res.Header().Add("Content-Type","application/json")
			res.Write(databytes)
		} else {

			jwtTok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 3).Unix(),
			})

			jwtString, errJwt := jwtTok.SignedString([]byte(os.Getenv("JWTSECRET")))
			if errJwt != nil {
				log.Fatal(errJwt.Error())
			}
			uuidstring := uuid.NewString()
			redis_rep := redisdb.ClientRedis.HSet(context.TODO(), uuidstring, "token", jwtString)
			if redis_rep.Err() != nil {
				log.Fatal(redis_rep.Err().Error())
				return
			}
			redisdb.ClientRedis.Expire(context.TODO(), uuidstring, time.Hour*3)
			response["message"] = "Success Login"
			var log logins.Logins
			log.TokenString = jwtString
			log.UserID =  user.ID
			log.Preset()
			db.Db_Connection.Create(&log)
			databytes, _ := json.Marshal(response)
			res.Header().Add("Authorization", uuidstring)
			res.Header().Add("Access-Control-Expose-Headers", "Authorization")
			res.Header().Add("Auth", jwtString)
			res.Header().Add("Access-Control-Expose-Headers", "Auth")
			res.Header().Add("Content-Type","application/json")
			res.Write(databytes)
		}

	}

}


func FetchAdmins(res http.ResponseWriter,req *http.Request){
time.Sleep(time.Nanosecond * 4)
user_id := <-helpers.UserIDChannel
fmt.Println(user_id)
var users []users.User
op := db.Db_Connection.Select("name","email","phone_number","id","super_user").Find(&users)
if op.RowsAffected > 0 && op.Error == nil {
json.NewEncoder(res).Encode(map[string]any{
"message":"Admins Found",
"data":users,
})
}else{
json.NewEncoder(res).Encode(map[string]any{
"message":"Admins Not Found",
"data":users,
})
}

}



func EditUser(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	userID := <-helpers.UserIDChannel
	fmt.Println(userID)
	var payload users.User

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	idParam := req.URL.Query().Get("id")
	ID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(res, "Invalid user id", http.StatusBadRequest)
		return
	}

	// Fetch existing user
	var existing users.User
	if err := db.Db_Connection.First(&existing, ID).Error; err != nil {
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	updates := make(map[string]any)

	// Name
	if payload.Name != "" && payload.Name != existing.Name {
		updates["name"] = payload.Name
	}

	// Email (check uniqueness only if changed)
	if payload.Email != "" && payload.Email != existing.Email {
		var count int64
		db.Db_Connection.
			Model(&users.User{}).
			Where("email = ? AND id != ?", payload.Email, ID).
			Count(&count)

		if count > 0 {
			http.Error(res, "Email already exists", http.StatusConflict)
			return
		}
		updates["email"] = payload.Email
	}

	// Phone number
	if payload.PhoneNumber != "" && payload.PhoneNumber != existing.PhoneNumber {
		updates["phone_number"] = payload.PhoneNumber
	}

	// Password
	if payload.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 13)
		updates["password"] = string(hashed)
	}

	// Role
	if payload.Role != "" && payload.Role != existing.Role {
		updates["role"] = payload.Role
	}

	// Booleans (always safe to update explicitly)
	updates["super_user"] = payload.SuperUser
	updates["org"] = payload.Org
	updates["verified"] = payload.Verified

	// Optional fields
	if payload.OrgBanner != "" {
		updates["org_banner"] = payload.OrgBanner
	}
	if payload.OrgPrimaryColors != "" {
		updates["org_primary_colors"] = payload.OrgPrimaryColors
	}

	if len(updates) == 0 {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "No changes detected",
		})
		return
	}

	if err := db.Db_Connection.Model(&existing).Updates(updates).Error; err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{
		"message": "User updated successfully",
	})
}


func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 5)
	userID := <-helpers.UserIDChannel
	fmt.Println(userID)
	var user users.User
	op := db.Db_Connection.First(&user, userID)
	if op.Error != nil {
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(res).Encode(map[string]any{
		"message": "User profile found",
		"data":    user,
	})
}
