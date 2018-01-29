package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type delta struct {
	Dx int
	Dy int
}

type str struct {
	Text string
}

func xdo(args ...string) {
	cmd := exec.Command("xdotool", args...)
	env := os.Environ()
	env = append(env, "DISPLAY=:0.0")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		log.Printf("xdotool error: %v\n", err)
		return
	}
}

func handleMoveMouse(w http.ResponseWriter, r *http.Request) {
	var d delta
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&d); err != nil {
		log.Printf("json error: %v\n", err)
		return
	}
	xdo("mousemove_relative", "--", strconv.Itoa(d.Dx), strconv.Itoa(d.Dy))
}

func handleScrollUp(w http.ResponseWriter, r *http.Request) {
	xdo("click", "4")
}

func handleScrollDown(w http.ResponseWriter, r *http.Request) {
	xdo("click", "5")
}

func handleClickMouse(w http.ResponseWriter, r *http.Request) {
	xdo("click", "1")
}

func handleInputText(w http.ResponseWriter, r *http.Request) {
	var s str
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&s); err != nil {
		log.Printf("json error: %v\n", err)
		return
	}
	xdo("type", s.Text)
}

func main() {
	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)
	http.HandleFunc("/movemouse", handleMoveMouse)
	http.HandleFunc("/clickmouse", handleClickMouse)
	http.HandleFunc("/inputtext", handleInputText)
	http.HandleFunc("/scrollup", handleScrollUp)
	http.HandleFunc("/scrolldown", handleScrollDown)

	portNum := "8080"
	if len(os.Args) > 1 {
		portNum = os.Args[1]
	}

	log.Fatal(http.ListenAndServe(":"+portNum, nil))
}
