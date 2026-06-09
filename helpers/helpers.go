package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesmukumu/maasaimaratripsportal/db"
	redisdb "jamesmukumu/maasaimaratripsportal/db/redisDB"
	"jamesmukumu/maasaimaratripsportal/models/users"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var UserIDChannel = make(chan uint,1)
func PrevalidatePayment(next http.HandlerFunc)http.HandlerFunc {
return func(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	header := req.Header.Get("Authorization")
	token,_ := strings.CutPrefix(header,"Bearer ")
	if token == ""{
	res.WriteHeader(401)
	json.NewEncoder(res).Encode(map[string]string{
	"message":"Unauthorized",
	})
	return
	}
    tokString,_ :=  redisdb.ClientRedis.HGet(context.TODO(),token,"token").Result()
	tokenPointer,err := jwt.Parse(tokString,func(t *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWTSECRET")),nil
	})

	if err != nil {
	res.WriteHeader(401)
	json.NewEncoder(res).Encode(map[string]string{
	"message":"Unauthorized",
	})
	return
	}
	
	var claims = tokenPointer.Claims.(jwt.MapClaims)

	user_id := claims["id"].(float64)
	fmt.Println(user_id)
	UserIDChannel<-uint(user_id)
	next.ServeHTTP(res,req)   
}}





func AdminHelper(next http.HandlerFunc)http.HandlerFunc {
return func(res http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	header := req.Header.Get("Authorization")
	token,_ := strings.CutPrefix(header,"Bearer ")
	if token == ""{
	res.WriteHeader(401)
	json.NewEncoder(res).Encode(map[string]string{
	"message":"Unauthorized",
	})
	return
	}

	tokenPointer,err := jwt.Parse(token,func(t *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWTSECRET")),nil
	})

	if err != nil {
	res.WriteHeader(401)
	json.NewEncoder(res).Encode(map[string]string{
	"message":"Unauthorized",
	})
	return
	}
	var user users.User
	var claims = tokenPointer.Claims.(jwt.MapClaims)

	user_id := claims["id"].(float64)
	op := db.Db_Connection.First(&user,user_id)
	if op.Error != nil || !user.SuperUser {
	http.Error(res,"Something Went Wrong",500)
	return
	}

	
	UserIDChannel<-uint(user_id)
	next.ServeHTTP(res,req)   
}}



func Helper(next http.HandlerFunc)http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		godotenv.Load()
		header := req.Header.Get("Authorization")
		token,_ := strings.CutPrefix(header,"Bearer ")
		if token == ""{
		res.WriteHeader(401)
		json.NewEncoder(res).Encode(map[string]string{
		"message":"Unauthorized",
		})
		return
		}
	
		tokenPointer,err := jwt.Parse(token,func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSECRET")),nil
		})
	
		if err != nil {
		res.WriteHeader(401)
		json.NewEncoder(res).Encode(map[string]string{
		"message":"Unauthorized",
		})
		return
		}
		
		var claims = tokenPointer.Claims.(jwt.MapClaims)
	
		user_id := claims["id"].(float64)
		
	
		
		UserIDChannel<-uint(user_id)
		next.ServeHTTP(res,req)   
	}}