package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

var actionMap = map[string]func(*json.RawMessage){
	"drag": func(msg *json.RawMessage) {
		type delta struct {
			Dx int
			Dy int
		}
		var d delta
		json.Unmarshal(*msg, &d)
		x, y := robotgo.GetMousePos()
		x += d.Dx
		y += d.Dy
		robotgo.DragMouse(x, y)
	},
	"scroll": func(msg *json.RawMessage) {
		type delta struct {
			Dy int
		}
		var d delta
		json.Unmarshal(*msg, &d)
		if d.Dy > 0 {
			robotgo.ScrollMouse(1, "up")
		} else {
			robotgo.ScrollMouse(1, "down")
		}
	},
	"click": func(msg *json.RawMessage) {
		type click struct {
			Button string
		}
		var c click
		json.Unmarshal(*msg, &c)
		robotgo.MouseClick(c.Button)
	},
	"type": func(msg *json.RawMessage) {
		type typer struct {
			Text string
		}
		var t typer
		json.Unmarshal(*msg, &t)
		s := t.Text
		switch s {
		case "":
			robotgo.KeyTap("enter")
		default:
			robotgo.TypeString(s)
		}
	},
}

func main() {
	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 0,
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer ws.Close()
		log.Println("opened new websocket")

		for {
			var jsonMap map[string]*json.RawMessage
			err := ws.ReadJSON(&jsonMap)
			if err != nil {
				log.Printf("websocket read error: %v\n", err)
				return
			}
			for command, msg := range jsonMap {
				if action, ok := actionMap[command]; ok {
					action(msg)
				}
			}
		}
	})

	portNum := "8080"
	if len(os.Args) > 1 {
		portNum = os.Args[1]
	}

	log.Printf("starting mousetrap on port %v\n", portNum)
	log.Fatal(http.ListenAndServe(":"+portNum, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
