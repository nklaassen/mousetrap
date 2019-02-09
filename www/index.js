var ws, lastX, lastY, touchStart, doubleTouchStart, lastTouch, lastScrollY, timerId

function Drag(dx, dy) {
	return {
		drag: {
			dx: dx,
			dy: dy
		}
	}
}
function Scroll(dy) {
	return {
		scroll: {
			dy: dy
		}
	}
}
function Type(text) {
	return {
		type: {
			text: text
		}
	}
}
function Click(button) {
	return {
		click: {
			button: button
		}
	}
}

function startup() {
	timerId = 0
	let loc = window.location
	let ws_uri = 'ws://' + loc.host + '/ws'
	resetWebSocket(ws_uri)
	let el = document.getElementById("touch")
	el.addEventListener("touchstart", handleTouchStart, false)
	el.addEventListener("touchend", handleTouchEnd, false)
	el.addEventListener("touchmove", handleTouchMove, false)
	el = document.getElementById("scroll")
	el.addEventListener("touchstart", handleScrollStart, false)
	el.addEventListener("touchmove", handleScrollMove, false)
	el = document.getElementById("textForm")
	el.onsubmit = handleFormSubmit
}

function resetWebSocket(ws_uri) {
	ws = new WebSocket(ws_uri)
	ws.onopen = () => {
		if (timerId) {
			clearInterval(timerId)
			timerId = 0
		}
	}
	ws.onclose = () => {
		if (!timerId) {
			timerId = setInterval(() => resetWebSocket(ws_uri), 3000)
		}
	}
}

function handleFormSubmit(e) {
	e.preventDefault()
	let text = Type(e.target[0].value)
	e.target.reset()
	ws.send(JSON.stringify(text))
}

function handleScrollStart(e) {
	lastScrollY = Math.floor(e.touches[0].screenY)
}
function handleScrollMove(e) {
	e.preventDefault()
	let y = Math.floor(e.touches[0].screenY)
	let dy = y - lastScrollY
	if(Math.abs(dy) > 10) {
		ws.send(JSON.stringify(Scroll(dy > 0 ? 1 : -1)))
		lastScrollY = y
	}
}

function handleTouchStart(e) {
	lastX = Math.floor(e.touches[0].screenX)
	lastY = Math.floor(e.touches[0].screenY)
	touchStart = Date.now()
	lastTouch = touchStart
}
function handleTouchEnd(e) {
	if(Date.now() - touchStart < 100) {
		ws.send(JSON.stringify(Click('left')))
	}
}
function handleTouchMove(e) {
	e.preventDefault()

	let time = Date.now()
	let dt = time - lastTouch
	lastTouch = time

	let x = e.touches[0].screenX
	let y = e.touches[0].screenY

	let dx = (x - lastX)
	dx = 4 * dx * Math.abs(dx) / dt
	dx = dx > 0 ? Math.ceil(dx) : Math.floor(dx)

	let dy = (y - lastY)
	dy = 4 * dy * Math.abs(dy) / dt
	dy = dy > 0 ? Math.ceil(dy) : Math.floor(dy)

	ws.send(JSON.stringify(Drag(dx, dy)))

	lastX = x
	lastY = y
}
