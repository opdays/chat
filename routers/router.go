package routers

import (
	"push/controllers"
	"github.com/astaxie/beego"
)


func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Handler("/ws",  &controllers.MainController{})
	beego.Router("/DescribeBaseMetrics",  &controllers.ServerControllers{})
}
