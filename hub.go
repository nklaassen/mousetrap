package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 0,
}

type hub struct {
	actions map[string]action
}

func newHub() (h *hub) {
	h = &hub{
		actions: make(map[string]action),
	}
	h.register(newMouseMover())
	h.register(&typer{})
	h.register(&scroller{})
	h.register(&clicker{})
	return h
}

func (h *hub) register(a action) {
	h.actions[a.path()] = a
}

func (h *hub) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("opened new websocket")
	go func() {
		defer ws.Close()
		for {
			objMap := make(map[string]*json.RawMessage)
			err := ws.ReadJSON(&objMap)
			if err != nil {
				log.Printf("websocket read error: %v\n", err)
				return
			}
			for k := range objMap {
				if a, ok := h.actions[k]; ok {
					a.do(objMap[k])
				}
			}
		}
	}()
}
