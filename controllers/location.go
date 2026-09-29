package controllers

import (
	"event-explorer/services"
	"fmt"

	"github.com/beego/beego/v2/server/web"
)

type LocationController struct {
	web.Controller
}

func (c *LocationController) Autocomplete() {
	c.EnableRender = false

	input := c.GetString("input")
	sessionToken := c.GetString("sessionToken")

	fmt.Println(input, sessionToken)

	output, err := services.FetchAutocompleteData(input, sessionToken)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch autocomplete predictions"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = output
	c.ServeJSON()
}
