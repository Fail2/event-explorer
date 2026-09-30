package controllers

import (
	"event-explorer/services"
	"event-explorer/validators"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type EventController struct {
	beego.Controller
}

func (c *EventController) GetEvents() {
	query := c.Ctx.Request.URL.Query()

	if err := validators.EventsValidator(query); err != nil {
		c.Redirect("/", 302)
		return
	}

	city := query.Get("city")
	countryCode := query.Get("countryCode")

	c.Data["City"] = city
	c.Data["CountryCode"] = countryCode

	cache := services.GetCacheInstance()

	cachedMusic, musicHit := cache.Get(city, countryCode, "Music")
	cachedSports, sportsHit := cache.Get(city, countryCode, "Sports")

	if musicHit && sportsHit {
		c.Data["MusicEvents"] = cachedMusic
		c.Data["SportsEvents"] = cachedSports
		c.TplName = "listing.tpl"
		return
	}

	results, errors := services.GetConcurrentEvents(city, countryCode)

	if err, exists := errors["music"]; exists {
		c.Data["MusicError"] = fmt.Sprintf("Music events are temporarily unavailable: %v", err)
	} else {
		c.Data["MusicEvents"] = results["music"]
		cache.Set(city, countryCode, "Music", results["music"])
	}

	if err, exists := errors["sports"]; exists {
		c.Data["SportsError"] = fmt.Sprintf("Sports events are temporarily unavailable: %v", err)
	} else {
		c.Data["SportsEvents"] = results["sports"]
		cache.Set(city, countryCode, "Sports", results["sports"])
	}

	c.TplName = "listing.tpl"
}
