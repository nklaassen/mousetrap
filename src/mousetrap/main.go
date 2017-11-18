package main

import (
    "log"
    "net/http"
	"encoding/json"
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

type scroll struct {
	Dy int
}

func xdo(args ...string) {
	cmd := exec.Command("xdotool", args...)
	env := os.Environ()
	env = append(env, "DISPLAY=:0.0")
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		log.Printf("error with xdotool: %v\n", err)
	}
}

func handleMoveMouse(w http.ResponseWriter, r *http.Request) {
	var d delta
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		panic(err)
	}
	xdo("mousemove_relative", "--", strconv.Itoa(d.Dx), strconv.Itoa(d.Dy))
}

func handleScrollUp(w http.ResponseWriter, r *http.Request) {
	xdo("click", "4");
}

func handleScrollDown(w http.ResponseWriter, r *http.Request) {
	xdo("click", "5");
}

func handleClickMouse(w http.ResponseWriter, r *http.Request) {
	xdo("click", "1")
}

func handleInputText(w http.ResponseWriter, r *http.Request) {
	var s str
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&s)
	if err != nil {
		panic(err)
	}
	if s.Text == "" {
		xdo("key", "Return")
	} else {
		xdo("type", s.Text)
	}
}

func main() {
    http.HandleFunc("/movemouse.go", handleMoveMouse);
    http.HandleFunc("/clickmouse.go", handleClickMouse);
    http.HandleFunc("/inputtext.go", handleInputText);
    http.HandleFunc("/scrollup.go", handleScrollUp);
    http.HandleFunc("/scrolldown.go", handleScrollDown);

    log.Fatal(http.ListenAndServe(":8080", nil))
}
