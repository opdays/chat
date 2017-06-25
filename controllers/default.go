package controllers

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"time"
	"log"
	"strings"
	"github.com/astaxie/beego/httplib"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Client)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //代表不检查
	},
}

type MainController struct {
	beego.Controller
}
type Client struct {
	Ip      string
	Join    string
	Image   string
	Message map[string]string
	Online  int
}

var ClientList []Client

func init() {

}
func (c *MainController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	beego.Debug("into servhttp")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		beego.Error(err)
	}
	clients[ws] = true
	image, _ := httplib.Get("http://opdays.com/music/api/randompic").Bytes()
	beego.Debug(image)
	//每次进来用户 给用户获取一张图片
	//ws.WriteJSON(Client{
	//	//推送给当前连接
	//	Image:string(image),
	//})
	for online := range clients {
		//推送给所有在线 更新在线人数
		online.WriteJSON(Client{
			Online: len(clients),
		})
	}
	defer ws.Close()
	for {
		var message map[string]string
		err := ws.ReadJSON(&message)
		if err != nil {
			beego.Error(err)
			delete(clients, ws)
			for online := range clients {
				online.WriteJSON(Client{
					Online: len(clients),
				})
			}
			break
		}
		client := Client{
			Ip:      strings.Split(r.RemoteAddr, ":")[0],
			Message: message,
			Join:    time.Now().Format("2006-01-02 15:04:05"), //go 格式化时间
			Image:   string(image),
		}
		beego.Info(client)
		broadcast <- client

	}
	// Make sure we close the connection when the function returns

	//target.Terminal(ws)
	//defer ws.Close()
	//for {
	//	// Send the newly received message to the broadcast channel
	//	err:=ws.WriteMessage(1,[]byte("websocket..."))
	//	if err != nil{
	//		beego.Info("exit")
	//		break
	//	}
	//	fmt.Println("gogo")
	//	time.Sleep(time.Second * 1)
	//}

}
func PushMessage() {
	beego.Debug("into push message")
	for {
		client := <-broadcast
		for ws := range clients {
			beego.Debug(1)

			beego.Info(client);
			err := ws.WriteJSON(client)
			if err != nil {
				log.Printf("client.WriteJSON error: %v", err)
				ws.WriteJSON(Client{
					Online: len(clients),
				})
				ws.Close()
				delete(clients, ws)
			}
		}
	}

}
