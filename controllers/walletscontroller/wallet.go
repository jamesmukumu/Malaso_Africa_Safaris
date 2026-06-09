package walletscontroller

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jamesmukumu/maasaimaratripsportal/db"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"jamesmukumu/maasaimaratripsportal/models/deposits"
	"jamesmukumu/maasaimaratripsportal/models/wallets"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
)

const baseApiUrl string = "https://payment.intasend.com/api/v1"

type PayloadDeposit struct {
	Amount      int    `json:"amount"`
	PhoneNumber string `json:"phoneNumber"`
}

type ExpressLoad struct {
	Amount        string      `json:"amount"`
	PhoneNumber   string      `json:"phone_number"`
	ApiRef        string      `json:"api_ref"`
	WalletID      string      `json:"wallet_id"`
	MobileTarriff interface{} `json:"mobile_tariff"`
}

func GenerateSecureRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func PhoneValidator(phone string) bool {
	if !strings.HasPrefix(phone, "254") && len(phone) != 12 {
		return false
	} else {
		return true
	}
}

// success deposit
// add token and also a helper channel tommorow
// also payment validation
var user_id uint

func InitiateMpesaDeposit(res http.ResponseWriter, req *http.Request) {

	time.Sleep(time.Nanosecond * 100)
	user_id = <-helpers.UserIDChannel
	fmt.Println(user_id)
	env.Load()
	client := http.Client{}
	var responseInterface map[string]interface{}
	var payload PayloadDeposit
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		panic(err.Error())
	}
	if !PhoneValidator(payload.PhoneNumber) {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Invalid Phone Number",
		})
	}

	amt := strconv.Itoa(payload.Amount)
	var PayloadIntasend = ExpressLoad{
		Amount:        amt,
		PhoneNumber:   payload.PhoneNumber,
		ApiRef:        "",
		WalletID:      "",
		MobileTarriff: "",
	}
	payloadBytes, err2 := json.Marshal(PayloadIntasend)
	if err2 != nil {
		panic(err2.Error())
	}
	reader := bytes.NewReader(payloadBytes)

	Request, errRequest := http.NewRequest("POST", baseApiUrl+"/payment/mpesa-stk-push/", reader)
	if errRequest != nil {
		panic(errRequest.Error())
	}

	Request.Header.Add("Content-Type", "application/json")
	Request.Header.Add("Authorization", "Bearer "+os.Getenv("INTASENDSECRET"))

	Response, errResponse := client.Do(Request)
	if errResponse != nil {
		panic(errResponse.Error())
	}

	resp, err1 := ioutil.ReadAll(Response.Body)
	if err1 != nil {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "error " + err1.Error(),
		})

	}
	// fmt.Println(string(resp))
	json.Unmarshal(resp, &responseInterface)
	var invoice = responseInterface["invoice"].(map[string]interface{})
	invoice_id := invoice["invoice_id"].(string)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"invoice_id": invoice_id,
		"id":         user_id,
		"exp":        time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("PaymentSecret")))
	res.Header().Add("Access-Control-Expose-Headers", "Authentication")
	res.Header().Add("Authentication", "Bearer "+tokenString)

	json.NewEncoder(res).Encode(map[string]string{
		"message": "Success",
		"id":      invoice_id,
	})

}

func DepositValidationMpesa(res http.ResponseWriter, req *http.Request) {
	env.Load()
	time.Sleep(time.Nanosecond * 100)
	user_id = <-helpers.UserIDChannel
	var responseInterface map[string]interface{}
	var client = http.Client{}
	var tokenHeader string = req.Header.Get("Authentication")
	var tokenString, _ = strings.CutPrefix(tokenHeader, "Bearer ")
	jwtTokenPointer, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("PaymentSecret")), nil
	})
	if err != nil {
		panic(err.Error())
	}
	claims := jwtTokenPointer.Claims.(jwt.MapClaims)
	var invoice_id string = claims["invoice_id"].(string)
	var payloadMap = map[string]string{
		"invoice_id": invoice_id,
	}
	databytes, _ := json.Marshal(payloadMap)
	bytesReader := bytes.NewReader(databytes)
	Request, err1 := http.NewRequest("POST", baseApiUrl+"/payment/status/", bytesReader)
	if err1 != nil {
		panic(err1.Error())
	}
	Request.Header.Add("Authorization", "Bearer "+os.Getenv("INTASENDSECRET"))
	Request.Header.Add("Content-Type", "application/json")
	Response, err2 := client.Do(Request)
	if err2 != nil {
		panic(err2.Error())
	}
	defer Response.Body.Close()
	responseBytes, errResp := ioutil.ReadAll(Response.Body)
	if errResp != nil {
		panic(errResp.Error())
	}
	fmt.Println(string(responseBytes))
	errUnmarshal := json.Unmarshal(responseBytes, &responseInterface)
	if errUnmarshal != nil {
		panic(errUnmarshal.Error())
	}
	var invoice map[string]interface{} = responseInterface["invoice"].(map[string]interface{})

	if invoice["state"] == "COMPLETE" {
		var meta map[string]interface{} = responseInterface["meta"].(map[string]interface{})
		var customer = meta["customer"].(map[string]interface{})
		var Deposit = deposits.Deposits{
			Invoice_Id:           invoice_id,
			Amount:               invoice["net_amount"].(string),
			UserID:               user_id,
			CustomersName:        customer["first_name"].(string) + " " + customer["last_name"].(string),
			CustomersPhoneNumber: customer["phone_number"].(string),
		}
		op := db.Db_Connection.Create(&Deposit)
		if op.RowsAffected > 0 && op.Error == nil {
			json.NewEncoder(res).Encode(map[string]string{
				"message":    "Wallet Updated",
				"invoice_id": Deposit.Invoice_Id,
			})
		} else if op.Error != nil {
			panic(op.Error.Error())
		}
	} else {
		json.NewEncoder(res).Encode(map[string]string{
			"message": "Payment Failed",
		})
	}
}

func FetchDepositTransaction(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	invoice_id := vars["invoice_id"]
	var deposit deposits.Deposits
	op := db.Db_Connection.Where("invoice_id=?", invoice_id).Find(&deposit)
	if op.Error == nil && op.RowsAffected > 0 {
		json.NewEncoder(res).Encode(map[string]interface{}{
			"data": deposit,
		})
	} else {
		http.Error(res, "Something went wrong"+op.Error.Error(), 500)
	}
}

func FetchMyDesposits(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Nanosecond * 6)
	user_id := <-helpers.UserIDChannel
	var wallet wallets.Wallet
	var desposits []deposits.Deposits
	op := db.Db_Connection.Where("user_id =?", user_id).Find(&desposits)
	op_wallets := db.Db_Connection.Where("user_id =?", user_id).Find(&wallet)
	if op.RowsAffected > 0 && op.Error == nil && op_wallets.RowsAffected > 0 && op_wallets.Error == nil {
		json.NewEncoder(res).Encode(map[string]any{
			"wallet":   wallet,
			"deposits": desposits,
		})
	} else {
		http.Error(res, "Something went wrong", 500)
	}

}
