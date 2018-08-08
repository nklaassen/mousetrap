package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

type delta struct {
	Dx int `json:"dx"`
	Dy int `json:"dy"`
}
type data struct {
	Delta  *delta
	Text   *string `json:"text"`
	Scroll *int    `json:"scroll"`
	Click  bool    `json:"click"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 16,
}

func doMoveMouse(d delta) {
	x, y := robotgo.GetMousePos()
	x += d.Dx
	y += d.Dy
	robotgo.MoveMouse(x, y)
}

func doScroll(scroll int) {
	if scroll > 0 {
		robotgo.ScrollMouse(scroll, "up")
	} else {
		robotgo.ScrollMouse(scroll, "down")
	}
}

func doClick() {
	robotgo.MouseClick("left", true)
}

func doInputText(s string) {
	switch s {
	case "":
		robotgo.KeyTap("enter")
	default:
		robotgo.TypeString(s)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v\n", err)
		return
	}
	go func() {
		defer conn.Close()
		var d data
		for {
			d = data{}
			err := conn.ReadJSON(&d)
			if err != nil {
				log.Printf("websocket read error: %v\n", err)
				return
			}
			log.Println(d)
			if d.Delta != nil {
				doMoveMouse(*d.Delta)
			}
			if d.Scroll != nil {
				doScroll(*d.Scroll)
			}
			if d.Text != nil {
				doInputText(*d.Text)
			}
			if d.Click {
				doClick()
			}
		}
	}()
}

func main() {
	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleWebSocket)

	portNum := "8080"
	if len(os.Args) > 1 {
		portNum = os.Args[1]
	}

	log.Fatal(http.ListenAndServe(":"+portNum, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
