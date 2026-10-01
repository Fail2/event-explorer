package controllers

import (
	"event-explorer/services"
	"event-explorer/validators"

	beego "github.com/beego/beego/v2/server/web"
)

type CacheController struct {
	beego.Controller
}

func (c *CacheController) Purge() {
	c.EnableRender = false

	if err := c.Ctx.Request.ParseForm(); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "failed to parse request",
			"count":   0,
		}
		c.ServeJSON()
		return
	}

	query := c.Ctx.Request.Form

	if err := validators.CachePurgeValidator(query); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
			"count":   0,
		}
		c.ServeJSON()
		return
	}

	city := query.Get("city")
	countryCode := query.Get("countryCode")
	category := query.Get("category")

	cache := services.GetCacheInstance()

	if city == "" && countryCode == "" && category == "" {
		count, message := cache.ClearData()

		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": message,
			"count":   count,
		}
		c.ServeJSON()
		return
	}

	count, message := cache.DeleteData(city, countryCode, category)

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": message,
		"count":   count,
	}

	c.ServeJSON()
}
