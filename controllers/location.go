package controllers

import (
	"event-explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type LocationController struct {
	beego.Controller
}

func (c *LocationController) Autocomplete() {
	c.EnableRender = false

	input := c.GetString("input")
	sessionToken := c.GetString("sessionToken")
	
	output, err := services.FetchAutocompleteData(input, sessionToken)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch autocomplete predictions"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = output
	c.ServeJSON()
}

func (c *LocationController) GetPlaceDetails() {
	c.EnableRender = false

	placeID := c.Ctx.Input.Param(":placeId")
	sessionToken := c.GetString("sessionToken")

	output, err := services.FetchPlaceDetailsData(placeID, sessionToken)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = output
	c.ServeJSON()
}
