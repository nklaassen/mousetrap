package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/gorilla/handlers"
)

type delta struct {
	Dx int `json:"dx"`
	Dy int `json:"dy"`
}
type data struct {
	Delta  delta  `json:"delta"`
	Scroll int    `json:"scroll"`
	Text   string `json:"text"`
	Click  bool   `json:"click"`
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
		xdo("click", "4")
	case 1:
		xdo("click", "5")
	}
}

func doClick(click bool) {
	if click {
		xdo("click", "1")
	}
}

func doInputText(text string) {
	if text == "" {
		xdo("key", "Return")
	} else {
		xdo("type", text)
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
			var d data
			err := conn.ReadJSON(&d)
			if err != nil {
				log.Printf("websocket read error: %v\n", err)
				return
			}
			log.Println(d)
			doMoveMouse(d.Delta)
			doScroll(d.Scroll)
			doInputText(d.Text)
			doClick(d.Click)
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
