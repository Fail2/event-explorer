package main

import (
	_ "event-explorer/routers"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if web.BConfig.RunMode == "dev"{
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/static"] = "static"
	}
	beego.Run()
}

