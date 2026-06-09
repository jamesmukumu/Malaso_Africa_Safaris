package scopes

import (
	"jamesmukumu/maasaimaratripsportal/db"
	blogs "jamesmukumu/maasaimaratripsportal/models/Blogs"
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/packages"
	"strings"

	"gorm.io/gorm"
)

func Paginator(query *gorm.DB,page int) *gorm.DB{
if page < 1 {
page = 1
}
limit := 6

offset := (page - 1) * limit
return query.Offset(offset).Limit(limit).Order("created_at DESC")
}

//Similar to this i want a scope that i will receive a  query param string called searchTerm
// with ths search Term i want to search 3 tables using %i% case insenstive term
// then return the results in paginated form as above then display like blogs:results,packages:results
// that is a json map [string]any
//table to be searched
//1 hotels columns hotel_name,hotel_overview,hotel_description,location_description,hotel_contact_email
//2 packages columns package_title,package_overview
//3 blogs columns blog_title,blog_overview,blog_description
// optimize the query to make the search not take long


var Hotel hotels.Hotels
var pack packages.Package
var destination destinations.Destination
var Blog blogs.Blogs
func QueryFinder(searchTerm string,page int)map[string]any{
var results= make(map[string]any,0)
termQuery := "%"+strings.TrimSpace(searchTerm) +"%"

var hotels []hotels.Hotels
var pcks []packages.Package
var Blogs []blogs.Blogs
var Dests []destinations.Destination
// this is where i query the hoteks based on term
hotelQuery := db.Db_Connection.Model(&Hotel)
hotelQuery = hotelQuery.Where("hotel_name ILIKE ? OR hotel_overview ILIKE ? OR hotel_description ILIKE ? OR location_description ILIKE ? OR hotel_contact_email ILIKE ?",termQuery,termQuery,termQuery,termQuery,termQuery).Where("published =?",true)
Paginator(hotelQuery,page).Find(&hotels)


blogsQuery := db.Db_Connection.Model(&Blog)
blogsQuery = blogsQuery.Where("blog_title ILIKE ? OR blog_overview ILIKE ? OR blog_description ILIKE ?",termQuery,termQuery,termQuery).Where("published =?",true)
Paginator(blogsQuery,page).Find(&Blogs)


packagesQuery := db.Db_Connection.Model(&pcks)
packagesQuery = packagesQuery.Where("package_title ILIKE ? OR package_overview ILIKE ? ",termQuery,termQuery).Where("published =?",true)
Paginator(packagesQuery,page).Find(&pcks)


destinationsQuery := db.Db_Connection.Model(&Dests)
destinationsQuery = destinationsQuery.Where("destination_title ILIKE ? OR destination_overview ILIKE ? OR destination_description ILIKE ? ",termQuery,termQuery,termQuery).Where("published =?",true)
Paginator(destinationsQuery,page).Find(&Dests)

results["hotels"] = hotels
results["blogs"] = Blogs 
results["destinations"] = Dests
results["packages"] = pcks
count := len(hotels) + len(Blogs) + len(pcks) + len(Dests)
results["count"] = count

return results
}
