package db

import (
	"fmt"
	blogs "jamesmukumu/maasaimaratripsportal/models/Blogs"
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	enquiries "jamesmukumu/maasaimaratripsportal/models/Enquiries"
	packagescategory "jamesmukumu/maasaimaratripsportal/models/PackagesCategory"
	"jamesmukumu/maasaimaratripsportal/models/adventures"
	"jamesmukumu/maasaimaratripsportal/models/attractions"
	"jamesmukumu/maasaimaratripsportal/models/awards"
	"jamesmukumu/maasaimaratripsportal/models/bulkmails"
	customizedpackages "jamesmukumu/maasaimaratripsportal/models/customizedPackages"
	"jamesmukumu/maasaimaratripsportal/models/deposits"
	"jamesmukumu/maasaimaratripsportal/models/drafts"
	"jamesmukumu/maasaimaratripsportal/models/email"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/logins"
	"jamesmukumu/maasaimaratripsportal/models/packages"
	"jamesmukumu/maasaimaratripsportal/models/payments"
	"jamesmukumu/maasaimaratripsportal/models/rates"
	"jamesmukumu/maasaimaratripsportal/models/resets"
	"jamesmukumu/maasaimaratripsportal/models/rooms"
	"jamesmukumu/maasaimaratripsportal/models/sitemaps"
	"jamesmukumu/maasaimaratripsportal/models/users"
	"jamesmukumu/maasaimaratripsportal/models/wallets"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
);
var Db_Connection *gorm.DB;
var user users.User;
var award awards.Award;
var wallet wallets.Wallet
var deposit deposits.Deposits
var Dest destinations.Destination
var packageCategory packagescategory.PackageCategory 
var Package packages.Package
var room rooms.Rooms
var hotel hotels.Hotels
var blogsCategory blogs.BlogsCategory
var bb blogs.Blogs
var enquiry enquiries.Enquiries
var attraction attractions.Attraction
var emailTemplates email.EmailTemplates
var newsletter email.Newsletter
var bulks bulkmails.BulkMails
var adventureCategory adventures.AdventuresCategory
var adventure adventures.Adventure
var initializedPayments payments.InitializedPayments
var completedPayments payments.CompletedPayment
var draft drafts.Draft
var sitemap sitemaps.SiteMap
var customPackages customizedpackages.CustomizedPackage
var Rates rates.Rates
var Log logins.Logins
var reset resets.Resets

func DBConnection(){
godotenv.Load()


	dsn := os.Getenv("Dbconnection")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{});
	if err != nil {
	 log.Fatal("Error Connecting to database"+err.Error());
	}
   Db_Connection = db;
   fmt.Println("Connected to DB")    
   Db_Connection.AutoMigrate(&user,&Log,&award,&wallet,&deposit,&Dest,&packageCategory,&Package,&hotel,&room,&blogsCategory,&bb,&enquiry,&attraction,&emailTemplates,&newsletter,&bulks,&adventureCategory,&adventure,&initializedPayments,&completedPayments,&draft,&sitemap,&customPackages,&Rates,&reset)       
 
}  