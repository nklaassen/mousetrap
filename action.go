package main

import (
	"encoding/json"
	"github.com/go-vgo/robotgo"
)

type action interface {
	path() string
	do(*json.RawMessage)
}

type delta struct {
	Dx int `json:"dx"`
	Dy int `json:"dy"`
}
type mouseMover struct {
	deltas chan delta
}

func newMouseMover() (m *mouseMover) {
	m = &mouseMover{
		make(chan delta, 128),
	}
	/* we can get a lot of mousemove events really fast, and moveMouse can be slow,
	   so start a goroutine that reads incoming deltas from a channel, sums all that
	   are ready, and batches them into a single moveMouse call
	*/
	go func() {
		for {
			d := <-m.deltas
		readChannel:
			for {
				select {
				case n := <-m.deltas:
					d = d.add(n)
				default:
					break readChannel
				}
			}
			moveMouse(d)
		}
	}()
	return m
}
func (m *mouseMover) do(msg *json.RawMessage) {
	var d delta
	json.Unmarshal(*msg, &d)
	m.deltas <- d
}
func (m *mouseMover) path() string {
	return "drag"
}
func moveMouse(d delta) {
	x, y := robotgo.GetMousePos()
	x += d.Dx
	y += d.Dy
	robotgo.DragMouse(x, y)
}
func (this delta) add(other delta) delta {
	return delta{
		this.Dx + other.Dx,
		this.Dy + other.Dy,
	}
}

type scroller struct {
	Dy int `json:"dy"`
}

func (s scroller) path() string {
	return "scroll"
}
func (s scroller) do(msg *json.RawMessage) {
	json.Unmarshal(*msg, &s)
	if s.Dy > 0 {
		robotgo.ScrollMouse(1, "up")
	} else {
		robotgo.ScrollMouse(1, "down")
	}
}

type clicker struct {
	Button string `json:"button"`
}

func (c clicker) path() string {
	return "click"
}
func (c clicker) do(msg *json.RawMessage) {
	json.Unmarshal(*msg, c)
	robotgo.MouseClick(c.Button)
}

type typer struct {
	Text string `json:"text"`
}

func (t typer) path() string {
	return "type"
}
func (t typer) do(msg *json.RawMessage) {
	json.Unmarshal(*msg, &t)
	s := t.Text
	switch s {
	case "":
		robotgo.KeyTap("enter")
	default:
		robotgo.TypeString(s)
	}
}
