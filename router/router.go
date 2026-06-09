package router

import (
	"fmt"
	adventurecontroller "jamesmukumu/maasaimaratripsportal/controllers/adventureController"
	assetscontroller "jamesmukumu/maasaimaratripsportal/controllers/assetsController"
	"jamesmukumu/maasaimaratripsportal/controllers/attractioncontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/awardscontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/blogscontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/destinationcontroller"
	draftscontroller "jamesmukumu/maasaimaratripsportal/controllers/draftsController"
	emailcontroller "jamesmukumu/maasaimaratripsportal/controllers/emailController"
	enquirycontroller "jamesmukumu/maasaimaratripsportal/controllers/enquiryController"
	"jamesmukumu/maasaimaratripsportal/controllers/hotelscontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/packagecontroller"
	paymentcontroller "jamesmukumu/maasaimaratripsportal/controllers/paymentController"
	querycontroller "jamesmukumu/maasaimaratripsportal/controllers/queryController"
	ratescontroller "jamesmukumu/maasaimaratripsportal/controllers/ratesController"
	"jamesmukumu/maasaimaratripsportal/controllers/roomscontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/sitemapscontroller"
	"jamesmukumu/maasaimaratripsportal/controllers/usercontrollers"
	"jamesmukumu/maasaimaratripsportal/controllers/walletscontroller"
	"jamesmukumu/maasaimaratripsportal/helpers"
	"log"
	"net/http"
	"os"

	env "github.com/joho/godotenv"

	mux "github.com/gorilla/mux"
	cors "github.com/rs/cors"
)

func Router() {
	env.Load()
	var port = os.Getenv("PORT")
	var router = mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./files"))
router.PathPrefix("/files/").
    Handler(http.StripPrefix("/files/", fileServer))


	// rates server
	ratesServer := http.FileServer(http.Dir("./rates"))
	router.PathPrefix("/rates/").
	Handler(http.StripPrefix("/rates/", ratesServer))
	var corsOptions = cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
		AllowedHeaders: []string{"*"},
	}
	corsHandler := cors.New(corsOptions)
	router.PathPrefix("/rates/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := "." + r.URL.Path // ./rates/filename.zip
	
		// Check file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	

		
		// Force download
		w.Header().Set("Content-Disposition", "attachment")
		w.Header().Set("Content-Type", "application/zip")
	
		http.ServeFile(w, r, filePath)
	})
	
	var actualHandler = corsHandler.Handler(router)
	router.HandleFunc("/api/create/new/user", usercontrollers.CreateUser).Methods("POST")
	router.HandleFunc("/api/login/user", usercontrollers.Login).Methods("POST")
    router.HandleFunc("/api/fetch/admins",helpers.AdminHelper(usercontrollers.FetchAdmins)).Methods("GET")
    router.HandleFunc("/api/edit/user",helpers.AdminHelper(usercontrollers.EditUser)).Methods("POST")
    router.HandleFunc("/api/get/user/profile",helpers.Helper(usercontrollers.GetUserProfile)).Methods("GET")


	// awards api`s`
	router.HandleFunc("/api/create/award", helpers.AdminHelper(awardscontroller.CreateAward)).Methods("POST")
	router.HandleFunc("/api/update/award", helpers.AdminHelper(awardscontroller.UpdateAward)).Methods("PUT")
	router.HandleFunc("/api/delete/award/{id}", helpers.AdminHelper(awardscontroller.DeleteAward)).Methods("DELETE")
	router.HandleFunc("/api/fetch/display/awards", awardscontroller.FetchAwards).Methods("GET")
	router.HandleFunc("/api/fetch/award", awardscontroller.FetchSlugwise).Methods("GET")
	router.HandleFunc("/api/fetch/awards", helpers.AdminHelper(awardscontroller.FetchAwardsAdmin)).Methods("GET")

	// Deposits api
	
	router.HandleFunc("/api/fetch/deposit/{invoice_id}", walletscontroller.FetchDepositTransaction).Methods("GET")
	router.HandleFunc("/api/deposit/stk", helpers.PrevalidatePayment(walletscontroller.InitiateMpesaDeposit))
	router.HandleFunc("/api/validate/deposit", helpers.Helper(walletscontroller.DepositValidationMpesa)).Methods("POST")
	router.HandleFunc("/api/fetch/wallet", helpers.Helper(walletscontroller.FetchMyDesposits)).Methods("GET")

	// Destination Controller
	router.HandleFunc("/api/add/destination", helpers.Helper(destinationcontroller.CreateDestination)).Methods("POST")
	router.HandleFunc("/api/fetch/my/destinations", helpers.Helper(destinationcontroller.FetchMyDestinations)).Methods("GET")
	router.HandleFunc("/api/fetch/destinations", destinationcontroller.FetchDisplayDestinations).Methods("GET")
	router.HandleFunc("/api/update/destination/{method}", helpers.Helper(destinationcontroller.PublishDestination)).Methods("PATCH")
    router.HandleFunc("/api/fetch/destination/{slug}",destinationcontroller.FetchDestination).Methods("GET")
    router.HandleFunc("/api/destination/slugs",destinationcontroller.DestinationSlugs).Methods("GET")
router.HandleFunc("/api/edit/destination",helpers.Helper(destinationcontroller.EditDestination)).Methods("POST")


	// Package Controller  
	router.HandleFunc("/api/packages/slug",packagecontroller.PackageSlugs).Methods("GET")
	router.HandleFunc("/api/create/package/category", helpers.Helper(packagecontroller.CreatePackageCategory)).Methods("POST")
	router.HandleFunc("/api/create/package", helpers.Helper(packagecontroller.CreatePackage)).Methods("POST")
	router.HandleFunc("/api/create/custom/package", helpers.Helper(packagecontroller.CreatePackage)).Methods("POST")
	router.HandleFunc("/api/pre/package", helpers.Helper(packagecontroller.Prepackage)).Methods("GET")
	router.HandleFunc("/api/fetch/package/active", packagecontroller.FetchPublishedPackages).Methods("GET")
	router.HandleFunc("/api/query/package", packagecontroller.FilterPackages).Methods("GET")
	router.HandleFunc("/api/fetch/orgs/packages", helpers.Helper(packagecontroller.FetchOrgsPackages)).Methods("GET")
	router.HandleFunc("/api/update/package/{method}", helpers.Helper(packagecontroller.PublishPackage)).Methods("PATCH")
	router.HandleFunc("/api/delete/package", helpers.Helper(packagecontroller.DeletePackage)).Methods("DELETE")
    router.HandleFunc("/api/find/package/categories",packagecontroller.FetchPackageCategorys).Methods("GET")
    router.HandleFunc("/api/category/packages/{name}",packagecontroller.PackageCategoryPackages).Methods("GET")
    router.HandleFunc("/api/fetch/single/package/{slug}",packagecontroller.FetchSinglePackage).Methods("GET")
  router.HandleFunc("/api/edit/package",helpers.Helper(packagecontroller.EditPackage)).Methods("POST")
router.HandleFunc("/api/edit/package/category",helpers.Helper(packagecontroller.EditPackageCategory)).Methods("POST")
router.HandleFunc("/api/fetch/pck/categories",helpers.Helper(packagecontroller.FetchPackageCategories))
router.HandleFunc("/api/fetch/category/{id}",helpers.Helper(packagecontroller.FetchSingularPackageCategory)).Methods("GET")
  // hotels apis
	router.HandleFunc("/api/create/hotel", helpers.Helper(hotelscontroller.AddHotels)).Methods("POST")
	router.HandleFunc("/api/fetch/hotels", helpers.Helper(hotelscontroller.FetchMyHotels)).Methods("GET")
	router.HandleFunc("/api/update/hotel/{method}", helpers.Helper(hotelscontroller.PublishHotel)).Methods("PATCH")
	router.HandleFunc("/api/fetch/hotels/display", hotelscontroller.FetchDisplayHotels).Methods("GET")
    router.HandleFunc("/api/fetch/hotel/{slug}",hotelscontroller.FetchHotel).Methods("GET")
    router.HandleFunc("/api/hotel/slugs",hotelscontroller.FetchHotelSlugs).Methods("GET")
    router.HandleFunc("/api/edit/hotel",helpers.Helper(hotelscontroller.EditHotel)).Methods("POST")
    


	// rooms api
	router.HandleFunc("/api/create/room", helpers.Helper(roomscontroller.CreateRoom)).Methods("POST")
	router.HandleFunc("/api/fetch/rooms", helpers.Helper(roomscontroller.FetchMyRooms)).Methods("GET")

	// attractions api
	router.HandleFunc("/api/create/attraction", helpers.Helper(attractioncontroller.SaveAttraction)).Methods("POST")
	router.HandleFunc("/api/fetch/my/attractions", helpers.Helper(attractioncontroller.FetchMyAttractions)).Methods("GET")
	router.HandleFunc("/api/update/attraction/{method}", helpers.Helper(attractioncontroller.PublishAttraction)).Methods("PATCH")
    router.HandleFunc("/api/fetch/attraction/{slug}",attractioncontroller.FetchSingularAttraction).Methods("GET")
    router.HandleFunc("/api/edit/attraction",helpers.Helper(attractioncontroller.EditAttraction)).Methods("POST")
    router.HandleFunc("/api/fetch/active/attractions",attractioncontroller.FetchActiveAttractions).Methods("GET")
	// blogs api
	router.HandleFunc("/api/create/blog/category", helpers.Helper(blogscontroller.CreateBlogCategory)).Methods("POST")
	router.HandleFunc("/api/create/blog", helpers.Helper(blogscontroller.CreateBlog)).Methods("POST")
	router.HandleFunc("/api/fetch/my/blogs", helpers.Helper(blogscontroller.FetchMyBlogs)).Methods("GET")
	router.HandleFunc("/api/update/blog/{method}", helpers.Helper(blogscontroller.PublishBlog)).Methods("PATCH")
	router.HandleFunc("/api/fetch/display/blogs",blogscontroller.FetchDisplayBlogs).Methods("GET")
	router.HandleFunc("/api/fetch/blog/{slug}",blogscontroller.FetchBlog).Methods("GET")
	router.HandleFunc("/api/blog/slugs",blogscontroller.FetchSlugs).Methods("GET")
	router.HandleFunc("/api/edit/blog",blogscontroller.EditBlog).Methods("POST")
	
	
	// media api
	router.HandleFunc("/api/folder/{folder}",helpers.Helper(assetscontroller.FetchAssets)).Methods("GET")
	router.HandleFunc("/api/upload/{folder}",helpers.Helper(assetscontroller.SaveAssets)).Methods("POST")
	router.Handle("/api/delete/assets",helpers.Helper(assetscontroller.DeleteAssets)).Methods("POST")

// emails
router.HandleFunc("/api/send/email",emailcontroller.SendEmail).Methods("POST")
router.HandleFunc("/api/save/email/template",helpers.Helper(emailcontroller.SaveEmailTemplate)).Methods("POST")
router.HandleFunc("/api/edit/email/template",helpers.Helper(emailcontroller.EditEmailTemplate)).Methods("POST")
router.HandleFunc("/api/save/newsletter/template",helpers.Helper(emailcontroller.SaveEmailNewsletter)).Methods("POST")
router.HandleFunc("/api/save/email/bulk",emailcontroller.SaveBulkEmail).Methods("POST")
router.HandleFunc("/api/fetch/bulks",helpers.Helper(emailcontroller.FetchEmailBulks)).Methods("GET")
router.HandleFunc("/api/send/template/{t}",helpers.Helper(emailcontroller.PropagateEmailsBulk)).Methods("POST")
router.HandleFunc("/api/manage/email/templates",helpers.Helper(emailcontroller.FetchTemplates)).Methods("GET")
router.HandleFunc("/api/find/email/template",helpers.Helper(emailcontroller.SingleTemplate)).Methods("GET")

// Enquiries
router.HandleFunc("/api/fetch/enquiry",enquirycontroller.FetchEnquiry).Methods("GET")
router.HandleFunc("/api/submit/enquiry",enquirycontroller.EnquiryHandler).Methods("POST")
router.HandleFunc("/api/fetch/enquiries",enquirycontroller.FetchEnquiries).Methods("GET")
// adventures
router.HandleFunc("/api/create/adventure/category",helpers.Helper(adventurecontroller.CreateAdventureCategory)).Methods("POST")
router.HandleFunc("/api/create/adventure",helpers.Helper(adventurecontroller.CreateAdventure)).Methods("POST")
router.HandleFunc("/api/fetch/adventures",helpers.Helper(adventurecontroller.FetchAdventures)).Methods("GET")
router.HandleFunc("/api/update/adventure/{method}", helpers.Helper(adventurecontroller.PublishAdventures)).Methods("PATCH")
router.HandleFunc("/api/fetch/adventure/category",adventurecontroller.FetchAdventureCategorys).Methods("GET")
router.HandleFunc("/api/adventure/category/{name}",adventurecontroller.AdventureCategoryAdventures).Methods("GET")
router.HandleFunc("/api/fetch/adventure/{slug}",adventurecontroller.FetchSingleAdventure).Methods("GET")
router.HandleFunc("/api/edit/adventure",helpers.Helper(adventurecontroller.EditAdventure)).Methods("POST")
router.HandleFunc("/api/update/adventure/category",helpers.Helper(adventurecontroller.UpdateAdventureCategory)).Methods("POST")
router.HandleFunc("/api/adventure/categories",helpers.Helper(adventurecontroller.FetchShowAdventureCategories)).Methods("GET")
router.HandleFunc("/api/adventure/{slug}",helpers.Helper(adventurecontroller.AdventureCat)).Methods("GET")



// Drafts Api
router.HandleFunc("/api/save/draft",helpers.Helper(draftscontroller.SaveDraft)).Methods("POST")

// payments Controller
router.HandleFunc("/api/initialize/payment/{slug}",paymentcontroller.InitializePackagePayment).Methods("POST")
router.HandleFunc("/api/validate/payment/{reference}",paymentcontroller.ValidatePayment).Methods("POST")

// rates controller
router.HandleFunc("/api/save/rates",helpers.Helper(ratescontroller.CreateRates)).Methods("POST")
router.HandleFunc("/api/check/rates",helpers.Helper(ratescontroller.RatesChecker)).Methods("GET")
// Sitemap Handler
router.HandleFunc("/api/save/sitemap",sitemapscontroller.SaveSiteMap).Methods("POST")
// Query Handler
router.HandleFunc("/api/filter/term/{page}",querycontroller.QueryHandler).Methods("GET")


// 
router.HandleFunc("/api/nav/stats",packagecontroller.NavStats).Methods("GET")
fmt.Printf("Server Listening for requests at port %s", port)
log.Fatal(http.ListenAndServe(":"+port, actualHandler))
}
