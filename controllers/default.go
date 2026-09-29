package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["MetaDescription"] = "Search and find upcoming live concert and sproting matches in your city"
	c.TplName = "home.tpl"
}
