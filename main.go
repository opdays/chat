package main

import (
	_ "push/routers"
	"github.com/astaxie/beego"
	"push/controllers"
)

func init() {
	/*tk1 := toolbox.NewTask("tk1", "0/5 * * * * *", func() error {
		controllers.PushMessage()
		fmt.Println("tk1");
		return nil
	})
	toolbox.AddTask("tk1", tk1)
	toolbox.StartTask()*/
}
func main() {
	go controllers.PushMessage()
	beego.Run()
}
