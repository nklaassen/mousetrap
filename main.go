package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

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

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 16,
	}
	env = append(os.Environ(), "DISPLAY=:0.0")
)

func xdo(args ...string) {
	cmd := exec.Command("xdotool", args...)
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		log.Printf("xdotool error: %v\n", err)
		return
	}
}

func doMoveMouse(d delta) {
	xdo("mousemove_relative", "--", strconv.Itoa(d.Dx), strconv.Itoa(d.Dy))
}

func doScroll(scroll int) {
	switch scroll {
	case -1:
		xdo("click", "5")
	case 1:
		xdo("click", "4")
	}
}

func doClick() {
	xdo("click", "1")
}

func doInputText(s string) {
	switch s {
	case "":
		xdo("key", "Return")
	default:
		xdo("type", s)
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
		for {
			d := new(data)
			err := conn.ReadJSON(d)
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
