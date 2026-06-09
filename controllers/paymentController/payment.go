package paymentcontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	enquiries "jamesmukumu/maasaimaratripsportal/models/Enquiries"
	"jamesmukumu/maasaimaratripsportal/models/packages"
	"jamesmukumu/maasaimaratripsportal/models/payments"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)
var baseUrl = "https://api.paystack.co";


var responseChannel = make(chan any,1)
func InitializePackagePayment(res http.ResponseWriter,req *http.Request) {
godotenv.Load()
client := http.Client{}
packageSlug := mux.Vars(req)["slug"]
var matchingPackage packages.Package
var matchingEnquiry enquiries.Enquiries
EnquiryID,_ := strconv.Atoi(req.URL.Query().Get("ID"))


db.Db_Connection.Where("package_slug =? ",packageSlug).Find(&matchingPackage)
db.Db_Connection.Where("ID =? ",EnquiryID).Find(&matchingEnquiry)

chargePrice := strconv.Itoa(matchingPackage.PackageCharge * 100)
Currency := matchingPackage.PackageChargeCurrency

var payload map[string]any = map[string]any{
"email":matchingEnquiry.Email,
"amount":chargePrice,
"currency":Currency,
"channels":[]string{"card", "bank", "apple_pay", "ussd", "qr", "mobile_money", "bank_transfer", "eft", "payattitude"},
}

payloadBytes,err := json.Marshal(&payload)
if err != nil {
log.Fatal(err.Error())
http.Error(res,"Something Went Wrong",500)
}

bytesReader := bytes.NewReader(payloadBytes)
request,err2 := http.NewRequest("POST",baseUrl+"/transaction/initialize",bytesReader)
if err2 != nil {
log.Fatal(err2.Error())
http.Error(res,err2.Error(),500)
return
}
request.Header.Add("Authorization", "Bearer "+os.Getenv("PaystackSecret"))
defer request.Body.Close()

response,err3 := client.Do(request)
if err3 != nil {
log.Fatal(err3.Error())
http.Error(res,err3.Error(),500)
return
}
responseBytes, _ := ioutil.ReadAll(response.Body)
var responsePayment map[string]any
err5 := json.Unmarshal(responseBytes,&responsePayment)
if err5 != nil {
log.Fatal(err5.Error())
http.Error(res,err5.Error(),500)
return
}
fmt.Println(responsePayment)
if responsePayment["message"] == "Authorization URL created"{
var dataResponse = responsePayment["data"].(map[string]any) 
var initializePayment = payments.InitializedPayments{
Reference: dataResponse["reference"].(string),
EnquiriesID: matchingEnquiry.ID,
Email: matchingEnquiry.Email,
AuthorizationUrl: dataResponse["authorization_url"].(string),
AccessCode: dataResponse["access_code"].(string),
}
db.Db_Connection.Create(&initializePayment)

token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
"access_code":dataResponse["access_code"],
"reference":dataResponse["reference"],
"enquiryID":EnquiryID,
"initializedPaymentID":initializePayment.ID,

})

tokenString,err := token.SignedString([]byte(os.Getenv("PAYMENTSECRETPAYSTACK")))
if err != nil {
log.Fatal(err.Error())
return
}
uuidstring := uuid.NewString()
op := redisdb.ClientRedis.HSet(context.TODO(),uuidstring,"payments",tokenString)
if op.Err() != nil{
log.Fatal(op.Err().Error())
return
}
redisdb.ClientRedis.Expire(context.TODO(),uuidstring,time.Minute * 45)
res.Header().Add("Payment",uuidstring)
res.Header().Add("Access-Control-Expose-Headers", "Payment")
json.NewEncoder(res).Encode(map[string]string{
"message":"Payment Initialized",
"data":dataResponse["authorization_url"].(string),
})
}else{
json.NewEncoder(res).Encode(map[string]string{
"message":"Payment Failed To Initialize",

})	
}
}




func ValidatePayment(res http.ResponseWriter,req *http.Request){
client := http.Client{}
referenceSlug := mux.Vars(req)["reference"]
token,err1 := redisdb.ClientRedis.HGet(context.TODO(),referenceSlug,"payments").Result()
if err1 != nil {
log.Fatal(err1.Error())
http.Error(res,err1.Error(),500)
return
}
jwtToken,err2 := jwt.Parse(token,func(t *jwt.Token) (any, error) {
return []byte(os.Getenv("PAYMENTSECRETPAYSTACK")),nil
})
if err2 != nil {
log.Fatal(err2.Error())
http.Error(res,err2.Error(),500)
return
}

claims := jwtToken.Claims.(jwt.MapClaims)
paystackReference := claims["reference"].(string)
initializedPaymentID := claims["initializedPaymentID"]
fmt.Print(claims)
fmt.Println(paystackReference)
request,err3 := http.NewRequest("GET",baseUrl+"/transaction/verify/"+paystackReference,nil)
if err3 != nil {
log.Fatal(err3.Error())
http.Error(res,err3.Error(),500)
return
}

request.Header.Add("Authorization", "Bearer "+os.Getenv("PaystackSecret"))
// defer request.Body.Close()
response,err6 := client.Do(request)
if err6 != nil {
log.Fatal(err6.Error())
http.Error(res,err6.Error(),500)
return
}
defer response.Body.Close()
var responseMap map[string]any
responseBytes,_ := ioutil.ReadAll(response.Body)
err5 := json.Unmarshal(responseBytes,&responseMap)
fmt.Println(responseMap)
if err5 != nil {
log.Fatal(err5.Error())
http.Error(res,err5.Error(),500)
return
}
var paymentStatus = responseMap["data"].(map[string]any)
if paymentStatus["status"].(string) == "success"{

var completedPayment = payments.CompletedPayment{
TransactionID: uint(paymentStatus["id"].(float64)),
Reference: paystackReference,
Amount:int(paymentStatus["amount"].(float64)),
Currency: paymentStatus["currency"].(string),
GateWayResponse: paymentStatus["gateway_response"].(string),
InitializedPaymentsID:uint(initializedPaymentID.(float64)),
} 
db.Db_Connection.Create(&completedPayment)
json.NewEncoder(res).Encode(map[string]string{
"message":"Payment Completed",
})
}else if paymentStatus["status"] == "ongoing"{
json.NewEncoder(res).Encode(map[string]string{
"message":"Payment Ongoing",
})	
}else if paymentStatus["status"] == "abandoned" || paymentStatus["status"] == "failed"{
json.NewEncoder(res).Encode(map[string]string{
"message":"Payment Aborted",
})	
}



}