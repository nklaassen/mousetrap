var ws, lastX, lastY, touchStartTime, lastScrollY, moving, scrolling

function Delta(dx, dy) {
	return {
		delta: {
			dx: dx,
			dy: dy,
		}
	}
}
function Scroll(scroll) {
	return {
		scroll: scroll
	}
}
function Text(text) {
	return {
		text: text
	}
}
function Click() {
	return {
		click: true
	}
}

function startup() {
	moving = false
	scrolling = false

	let loc = window.location
	let ws_uri = 'ws://' + loc.host + '/ws'
	ws = new WebSocket(ws_uri)
	ws.onopen = () => console.log('ws opened')
	ws.onclose = () => console.log('ws closed')
	ws.onerror = (err) => console.log('ws error: ' + err)
	ws.onmessage = (msg) => console.log('ws: ' + msg.data)


	let el = document.getElementById("touch")
	el.addEventListener("mousedown", handleTouchStart, false)
	el.addEventListener("mouseup", handleTouchEnd, false)
	el.addEventListener("mouseleave", handleTouchEnd, false)
	el.addEventListener("mousemove", handleTouchMove, false)
	el = document.getElementById("scroll")
	el.addEventListener("mousedown", handleScrollStart, false)
	el.addEventListener("mousemove", handleScrollMove, false)
	el.addEventListener("mouseup", handleScrollEnd, false)
	el.addEventListener("mouseleave", handleScrollEnd, false)
	el = document.getElementById("textForm")
	el.onsubmit = formSubmit
}

function formSubmit(e) {
	e.preventDefault()
	let text = Text(e.target[0].value)
	e.target.reset()
	ws.send(JSON.stringify(text))
}

function handleScrollStart(e) {
	scrolling = true
	lastScrollY = Math.floor(e.screenY)
}
function handleScrollMove(e) {
	e.preventDefault()
	if (!scrolling)
		return
	let y = Math.floor(e.screenY)
	let dy = y - lastScrollY
	if(Math.abs(dy) > 5) {
		ws.send(JSON.stringify(Scroll(dy > 0 ? 1 : -1)))
		lastScrollY = y
	}
}
function handleScrollEnd(e) {
	scrolling = false
}

function handleTouchStart(e) {
	moving = true
	console.log('mousedown')
	lastX = Math.floor(e.screenX)
	lastY = Math.floor(e.screenY)
	touchStartTime = Date.now()
}
function handleTouchEnd(e) {
	moving = false
	console.log('mouseup')
	if(Date.now() - touchStartTime < 100) {
		ws.send(JSON.stringify(Click()))
	}
}
function handleTouchMove(e) {
	e.preventDefault()
	if (!moving)
		return
	let x = e.screenX
	let y = e.screenY

	let dx = (x - lastX)
	dx = dx * Math.abs(dx) / 4
	dx = dx > 0 ? Math.ceil(dx) : Math.floor(dx)

	let dy = (y - lastY)
	dy = dy * Math.abs(dy) / 4
	dy = dy > 0 ? Math.ceil(dy) : Math.floor(dy)

	ws.send(JSON.stringify(Delta(dx, dy)))

	lastX = x
	lastY = y
}

function send(url, msg) {
	let xhr = new XMLHttpRequest()
	xhr.open("post", url, true)
	xhr.send(msg)
}
