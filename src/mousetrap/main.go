package main

import (
    "log"
    "net/http"
	"encoding/json"
	"os/exec"
	"strconv"
)

type delta struct {
	Dx int
	Dy int
}

func handleMoveMouse(w http.ResponseWriter, r *http.Request) {
	var d delta
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("xdotool", "mousemove_relative", "--", strconv.Itoa(d.Dx), strconv.Itoa(d.Dy))
	err = cmd.Run()
	if err != nil {
		log.Printf("error with xdotool: %v", err)
	}
}

func handleClickMouse(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("xdotool", "click", "1")
	err := cmd.Run()
	if err != nil {
		log.Printf("error with xdotool: %v", err)
	}
}

func handleInputText(w http.ResponseWriter, r *http.Request) {
}

func main() {
    http.HandleFunc("/movemouse.go", handleMoveMouse);
    http.HandleFunc("/clickmouse.go", handleClickMouse);
    http.HandleFunc("/inputtext.go", handleInputText);

    log.Fatal(http.ListenAndServe(":8080", nil))
}
