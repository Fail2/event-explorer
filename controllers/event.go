package controllers

import (
	"event-explorer/validators"

	beego "github.com/beego/beego/v2/server/web"
)

type EventController struct {
	beego.Controller
}

func (c *EventController) GetEvents() {
	query := c.Ctx.Request.URL.Query()

	if err := validators.EventsValidator(query); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		return
	}

}
