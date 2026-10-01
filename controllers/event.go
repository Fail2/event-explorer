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

	if musicHit {
		c.Data["MusicEvents"] = cachedMusic
		c.Data["MusicCacheStatus"] = "Cache Hit"
	}

	if sportsHit {
		c.Data["SportsEvents"] = cachedSports
		c.Data["SportsCacheStatus"] = "Cache Hit"
	}

	if musicHit && sportsHit {
		c.TplName = "listing.tpl"
		return
	}

	results, errors := services.GetConcurrentEvents(city, countryCode)

	if !musicHit {
		if err, exists := errors["music"]; exists {
			c.Data["MusicError"] = fmt.Sprintf("Music events are temporarily unavailable: %v", err)
			c.Data["MusicCacheStatus"] = "Expired Data"
		} else {
			c.Data["MusicEvents"] = results["music"]
			cache.Set(city, countryCode, "Music", results["music"])
			c.Data["MusicCacheStatus"] = "Fresh Data"
		}
	}

	if !sportsHit {
		if err, exists := errors["sports"]; exists {
			c.Data["SportsError"] = fmt.Sprintf("Sports events are temporarily unavailable: %v", err)
			c.Data["SportsCacheStatus"] = "Expired Data"
		} else {
			c.Data["SportsEvents"] = results["sports"]
			cache.Set(city, countryCode, "Sports", results["sports"])
			c.Data["SportsCacheStatus"] = "Fresh Data"
		}
	}

	c.TplName = "listing.tpl"
}

func (c *EventController) GetEventDetails() {
	eventID := c.Ctx.Input.Param(":id")

	if eventID == "" {
		c.Redirect("/", 302)
		return
	}

	event, err := services.FetchTicketmasterEventDetails(eventID)
	if err != nil {
		c.Data["Error"] = fmt.Sprintf(
			"Event details are temporarily unavailable: %v",
			err,
		)
		c.TplName = "error.tpl"
		return
	}

	c.Data["Event"] = event
	c.TplName = "details.tpl"
}
