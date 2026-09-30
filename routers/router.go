package routers

import (
	"event-explorer/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/locations/autocomplete", &controllers.LocationController{}, "get:Autocomplete")
	beego.Router("/api/locations/:placeId", &controllers.LocationController{}, "get:GetPlaceDetails")
	beego.Router("/events", &controllers.EventController{}, "get:GetEvents")
	beego.Router("/events/:id", &controllers.EventController{}, "get:GetEventDetails")
}
